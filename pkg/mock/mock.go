// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package mock provides primitives for writing interface mocks.
//
// It is the runtime foundation used by mocks generated with
// [github.com/ctx42/testing/pkg/mocker] and by hand-written mocks. The
// package integrates with [github.com/ctx42/testing/pkg/tester] for test
// lifecycle management and [github.com/ctx42/testing/pkg/notice] for rich,
// structured failure messages.
//
// See the package [README] for the full expectation DSL, advanced usage
// (proxying, custom matchers, altering arguments, etc.), and go:generate
// patterns. See [examples/mock_test.go] (in the module root) for a complete
// working example.
//
// Key types and entry points:
//   - [NewMock] and [Mock] — the core mock controller
//   - [Mock.On], [Mock.OnAny], [Mock.Proxy] — define expectations
//   - [Call] and its chain methods (Return, Times, Until, ...)
//   - [Arguments] — typed getters for return values and call recording
//   - Matchers: [Any], [AnyString], [MatchBy], [MatchOfType], [MatchError], ...
package mock

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"
)

// Mock requirements violations.
var (
	// ErrNeverCalled is returned when a non-optional expected method was
	// never called.
	ErrNeverCalled = errors.New("method never called")

	// ErrTooFewCalls is returned when a method was called fewer times than
	// required by [Call.Times] or [Call.Once].
	ErrTooFewCalls = errors.New("method called too few times")

	// ErrTooManyCalls is returned when a method was called more times than
	// allowed by [Call.Times] or [Call.Once].
	ErrTooManyCalls = errors.New("method called too many times")

	// ErrRequirements is returned when a call's prerequisites (see
	// [Call.Requires]) were not satisfied.
	ErrRequirements = errors.New("method requirements not met")

	// ErrNotFound is returned when no matching expectation was found for a
	// call (see [Mock.Call] and [Mock.Called]).
	ErrNotFound = errors.New("method not found")
)

const (
	// Any is a sentinel used with [Arguments.Diff] and in expectation
	// definitions to indicate that the argument value should not be
	// considered during matching.
	Any = "mock.Any"
)

// Mock fatal message headers (internal).
const (
	hNeverCalled    = "[mock] method never called"
	hTooFewCalls    = "[mock] too few method calls"
	hTooManyCalls   = "[mock] too many method calls"
	hUnexpectedCall = "[mock] unexpected method call"
	hNotFoundCall   = "[mock] method call not found"
)

// dumper is the default value renderer used for diagnostic output.
var dumper = dump.New()

// Option configures a [Mock] created by [NewMock].
type Option func(*Mock)

// WithNoStack disables stack traces in the diagnostic messages produced by
// the mock when expectations are violated.
func WithNoStack(mck *Mock) { mck.stack = false }

// Mock tracks expected and actual calls on a mocked interface.
//
// A [Mock] is typically embedded in a hand-written or generated *Mock struct
// that implements the target interface. Use [NewMock] to create one; it
// automatically registers a [testing.TB.Cleanup] handler that calls
// [Mock.AssertExpectations] at the end of the test.
type Mock struct {
	// Calls expected on the mock.
	expected []*Call

	// Calls made on the mock.
	calls []cStack

	// Holds any data that might be useful for testing. The Mock ignores it,
	// allowing you to do whatever you like with it.
	meta map[string]any

	// When true, error log messages will contain stack traces.
	// It is mostly useful during tests. I don't see any reason why someone
	// would want to set it to false. Do you?
	stack bool

	// Set to true if mock is in a failed state.
	failed bool

	// Guards the Mock fields.
	mx sync.Mutex

	// Test manager.
	t tester.T
}

// NewMock creates and returns a new [Mock] bound to the provided tester.
//
// The mock registers an automatic cleanup that invokes
// [Mock.AssertExpectations] when the test completes. Use the [Option]
// functions (currently only [WithNoStack]) to customize behavior.
func NewMock(t tester.T, opts ...Option) *Mock {
	t.Helper()
	mck := &Mock{t: t, stack: true}
	for _, opt := range opts {
		opt(mck)
	}
	t.Cleanup(func() { t.Helper(); mck.AssertExpectations() })
	return mck
}

// MetaSetAll stores arbitrary data on the mock for later retrieval via
// [Mock.MetaAll]. The data is never used by the mock itself.
func (mck *Mock) MetaSetAll(data map[string]any) {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.meta = data
}

// MetaAll returns the data previously stored with [Mock.MetaSetAll], or an
// empty map if nothing was set.
func (mck *Mock) MetaAll() map[string]any {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	if mck.meta == nil {
		mck.meta = make(map[string]any)
	}
	return mck.meta
}

// On adds an expectation that the named method will be called with the
// given arguments (or matchers). Returns a [Call] for further configuration
// ([Return], [Times], [Until], etc.).
func (mck *Mock) On(method string, args ...any) *Call {
	mck.t.Helper()
	for _, arg := range args {
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Func {
			panic("cannot use functions in argument expectations")
		}
	}

	mck.mx.Lock()
	defer mck.mx.Unlock()
	call := newCall(method, args...).withParent(mck).withStack(callStack())
	mck.expected = append(mck.expected, call)
	return call
}

// OnAny adds an expectation that the named method may be called with any
// arguments (values and count are ignored). Returns a [Call] for further
// configuration. Useful when you only care that the method was invoked.
func (mck *Mock) OnAny(method string) *Call {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()

	call := newCall(method).withParent(mck).withStack(callStack())
	call.argsAny = true
	mck.expected = append(mck.expected, call)
	return call
}

// Proxy configures the mock to forward calls for the given method to the
// provided real implementation (name can override the detected method name).
//
// It panics if met is not a valid non-nil method or function.
//
// See [Call.With] for argument validation on proxied calls.
func (mck *Mock) Proxy(met any, name ...string) *Call {
	mck.t.Helper()

	typ := reflect.TypeOf(met)
	val := reflect.ValueOf(met)
	if typ.Kind() != reflect.Func || val.IsNil() {
		panic("Proxy requires a valid not nil method")
	}

	call := newProxy(val, name...).withParent(mck).withStack(callStack())
	mck.expected = append(mck.expected, call)
	return call
}

// Called records a method invocation and returns the configured return
// values. It uses runtime caller information to determine the method name
// and panics if no matching expectation was found.
//
// This is the method normally called from inside generated or hand-written
// mock implementations.
//
// If the matching expectation uses [Call.Until] or [Call.After], the call
// blocks until the condition is met.
func (mck *Mock) Called(args ...any) Arguments {
	mck.t.Helper()

	method := mck.called(2)
	return mck.Call(method, args...)
}

// called returns method name. The argument skip is the number of stack frames
// to ascend, with 0 identifying the caller of Caller.
func (mck *Mock) called(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		panic("could not get the caller information")
	}
	name := runtime.FuncForPC(pc).Name()
	if strings.Contains(name, ".func") {
		return "<anonymous>"
	}
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

// Call invokes the named method with the given arguments and returns the
// configured [Arguments] return values. It panics (via t.Fatal) if no
// matching expectation exists or if prerequisites are not met.
//
// Call records a method invocation by name and returns the configured
// return values. Useful for advanced or dynamic scenarios (most users
// go through generated wrappers that call [Mock.Called]).
//
// The call blocks if the matching expectation uses [Call.Until] or [Call.After].
func (mck *Mock) Call(method string, args ...any) Arguments {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()

	var cs []string
	if mck.stack {
		cs = callStack()
	}

	call, err := mck.find(method, args, cs)
	if err != nil {
		mck.failed = true
		mck.t.Fatal(err)
	}

	if err = call.checkReq(cs); err != nil {
		mck.failed = true
		mck.t.Fatal(err)
	}

	mck.calls = append(mck.calls, cStack{Method: method, Stack: cs})
	return call.call(args...)
}

// Callable reports whether a method with the given name and arguments can be
// called right now without violating expectations or prerequisites. It
// returns nil when a matching callable [Call] is found, otherwise a
// descriptive error (one of the Err* sentinels or a richer [notice.Notice]).
//
// This is useful for introspection or custom test logic; normal usage goes
// through [Mock.Called] / [Mock.Call].
func (mck *Mock) Callable(method string, args ...any) error {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()
	_, err := mck.find(method, args, nil)
	return err
}

// find finds a callable method with given name and matching arguments. When
// found it returns it, otherwise it returns an error describing the reason.
// Note that there may be more methods in the expected slice matching the
// criteria in which case the first one is returned.
//
// A callable method is one that returns no error from [Call.CanCall] method,
// and has matching arguments.
//
// nolint: cyclop
func (mck *Mock) find(method string, args []any, cs []string) (*Call, error) {
	var err error

	// Find a method (including proxies) with matching arguments.
	for _, call := range mck.expected {
		if call.Method != method {
			continue
		}
		err = call.CanCall()
		if errors.Is(err, ErrTooManyCalls) {
			continue
		}
		if call.argsAny && err == nil {
			return call, nil
		}
		if _, cnt := call.args.Diff(args); cnt == 0 {
			if err == nil {
				return call, nil
			}
		}
	}

	// Find a proxy method which was added without a need for matching
	// arguments.
	for _, call := range mck.expected {
		if call.Method != method {
			continue
		}
		if call.proxy.IsValid() && len(call.args) == 0 {
			if err = call.CanCall(); err == nil {
				return call, nil
			}
		}
	}

	if err != nil {
		return nil, err
	}

	var msg *notice.Notice

	// Try to find a method that is most similar to the one we are processing.
	if closest, diff := mck.closest(method, args...); closest != nil {
		// Similar method found.
		desc := formatMethod(closest.Method, closest.args, closest.returns)
		msg = notice.New(hUnexpectedCall).
			Append("closest", "%s", desc).
			Append("argument match", "\n%s", strings.Join(diff, "\n")).
			Wrap(ErrNotFound)
	} else {
		// This is totally unexpected method call.
		msg = notice.New(hNotFoundCall).
			Append("method", "%s", formatMethod(method, args, nil)).
			Wrap(ErrNotFound)
		if len(args) > 0 {
			_ = msg.Append("with args", "\n%s", formatArgs(args))
		}
	}

	if len(cs) > 0 {
		_ = msg.Append("stack", "\n%s", strings.Join(cs, "\n"))
	}
	return nil, msg
}

// closest returns expected call with the most similar (closest) arguments.
// It returns the found call and the argument difference between the provided
// arguments and the found one. If the method is not found it will return nil
func (mck *Mock) closest(method string, args ...any) (*Call, []string) {
	var best candidate
	for _, call := range mck.expected {
		if call.Method != method {
			continue
		}
		diff, cnt := call.args.Diff(args)
		current := candidate{
			call:    call,
			diff:    diff,
			diffCnt: cnt,
		}
		if current.betterThan(best) {
			best = current
		}
	}
	return best.call, best.diff
}

// Failed reports whether the mock has entered a failed state (unexpected
// call, unsatisfied expectation, etc.).
func (mck *Mock) Failed() bool {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()
	return mck.failed
}

// Unset removes a previously registered [Call] expectation. If the call is
// not found it records a test error.
func (mck *Mock) Unset(remove *Call) *Mock {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()

	var found bool
	var expected []*Call
	for _, have := range mck.expected {
		if remove == have {
			found = true
			continue
		}
		expected = append(expected, have)
	}
	mck.expected = expected

	if !found {
		method := formatMethod(remove.Method, remove.args, nil)
		msg := notice.New("[mock] unsetting non-existing method").
			Append("method", "%s", method).
			Wrap(ErrNotFound)
		mck.t.Error(msg)
	}
	return mck
}

// AssertExpectations verifies that all non-optional expectations defined via
// [Mock.On], [Mock.OnAny], etc. have been satisfied (correct call counts and
// prerequisites). It may be called manually, but [NewMock] registers it as a
// test cleanup so it runs automatically.
//
// Returns true when all expectations are met. On failure it records the
// problem via t.Error (or t.Fatal for unexpected calls) and returns false.
func (mck *Mock) AssertExpectations() bool {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()

	if mck.failed {
		return false
	}

	var names, whys []string
	for _, call := range mck.expected {
		if call.Satisfied() {
			continue
		}
		wCls := "calls"
		if call.wantCalls == 1 {
			wCls = "call"
		}
		hCls := "calls"
		if call.haveCalls == 1 {
			hCls = "call"
		}

		names = append(names, call.Method)
		var why string
		if call.wantCalls == 0 {
			format := "expected at least one call received %d %s"
			why = fmt.Sprintf(format, call.haveCalls, hCls)
		} else {
			format := "expected %d %s received %d %s"
			why = fmt.Sprintf(format, call.wantCalls, wCls, call.haveCalls, hCls)
		}
		whys = append(whys, why)
	}
	if len(names) == 0 {
		mck.failed = false
		return true
	}
	twoColumns(names, whys)
	msg := notice.New(hTooFewCalls).
		Append("missing calls", "\n%s", strings.Join(names, "\n")).
		Wrap(ErrTooFewCalls)

	mck.failed = true
	mck.t.Error(msg)
	return false
}

// AssertCallCount asserts that the named method was invoked exactly "want"
// times. Useful when you only care about call count rather than full
// expectation configuration.
func (mck *Mock) AssertCallCount(method string, want int) bool {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()
	var have int
	for _, call := range mck.calls {
		if call.Method == method {
			have++
		}
	}
	if have == want {
		return true
	}

	var msg *notice.Notice
	if want < have {
		msg = notice.New(hTooFewCalls).Wrap(ErrTooFewCalls)
	} else {
		msg = notice.New(hTooManyCalls).Wrap(ErrTooManyCalls)
	}
	_ = msg.Append("method", "%s", method).
		Append("want calls", "%d", want).
		Append("have calls", "%d", have)
	mck.t.Error(msg)
	mck.failed = true
	return false
}
