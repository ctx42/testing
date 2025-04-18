// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package mock provides helpers for creating and testing with interface mocks.
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
	// ErrNeverCalled signals mocked non-optional method was never called.
	ErrNeverCalled = errors.New("method never called")

	// ErrTooFewCalls signals mocked method was called too few times.
	ErrTooFewCalls = errors.New("method called to few times")

	// ErrTooManyCalls signals mocked method was called too many times.
	ErrTooManyCalls = errors.New("method called too many times")

	// ErrRequirements signals mocked method requirements were not met.
	ErrRequirements = errors.New("method requirements not met")

	// ErrNotFound signals mocked method has not been found.
	ErrNotFound = errors.New("method not found")
)

const (
	// Any is used in Diff and Assert when the argument being tested
	// shouldn't be taken into consideration.
	Any = "mock.Any"
)

// Mock fatal message headers.
const (
	hNeverCalled    = "[mock] method never called"
	hTooFewCalls    = "[mock] too few method calls"
	hTooManyCalls   = "[mock] too many method calls"
	hUnexpectedCall = "[mock] unexpected method call"
	hNotFoundCall   = "[mock] method call not found"
)

// dumper represents default value dumper.
var dumper = dump.New()

// Option represents a [NewMock] option.
type Option func(*Mock)

// WithNoStack is option for [NewMock] turning off displaying stack traces in
// error log messages.
func WithNoStack(mck *Mock) { mck.stack = false }

// Mock tracks activity on a mocked interface.
type Mock struct {
	// Calls expected on the mock.
	expected []*Call

	// Calls made on the mock.
	calls []cStack

	// Data holds any data that might be useful for testing. Mock ignores it
	// allowing you to do whatever you like with it.
	data map[string]any

	// When true error log messages will contain stack traces.
	// It is mostly useful during tests. I don't see any reason why someone
	// would want to set it to false. Do you?
	stack bool

	// Set to true if mock is in failed state.
	failed bool

	// Guards the Mock fields.
	mx sync.Mutex

	// Test manager.
	t tester.T
}

// NewMock returns new instance of Mock.
func NewMock(t tester.T, opts ...Option) *Mock {
	t.Helper()
	mck := &Mock{t: t, stack: true}
	for _, opt := range opts {
		opt(mck)
	}
	t.Cleanup(func() { t.Helper(); mck.AssertExpectations() })
	return mck
}

// SetData sets data that might be useful for testing. The Mock ignores it. To
// get it back call GetData method.
func (mck *Mock) SetData(data map[string]any) {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.data = data
}

// GetData returns data that might be useful for testing (see SetData method).
func (mck *Mock) GetData() map[string]any {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()
	if mck.data == nil {
		mck.data = make(map[string]any)
	}
	return mck.data
}

// On adds method call expectation for the interface being mocked.
//
// Example usage:
//
//	Mock.On("Method", 1).Return(nil)
//	Mock.On("MyOtherMethod", 'a', 'b', 'c').Return(errors.New("Some Error"))
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

// OnAny adds a method call expectation for the mocked interface. Unlike
// [Mock.On], where specific arguments are expected, this allows the method to
// be called with any combination of arguments, regardless of their number or
// type.
//
// Example usage:
//
//	Mock.OnAny("Method").Return(nil)
func (mck *Mock) OnAny(method string) *Call {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()

	call := newCall(method).withParent(mck).withStack(callStack())
	call.argsAny = true
	mck.expected = append(mck.expected, call)
	return call
}

// Proxy uses passed method as a proxy for calls to its "name". If "name"
// argument is not empty the first value from the slice will be used as the
// proxied method name. It panics if "met" is not a method or function.
//
// Example:
//
//	obj := &types.TPtr{}
//	mck.Proxy(obj.Method)
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

// Called tells the mock object that a method has been called with given
// arguments, and returns an array of arguments to return. Panics if the call
// is unexpected (i.e. not preceded by appropriate [Mock.On] or [Mock.OnAny]
// calls).
//
// If [Call.Until] or [Call.After] is set, it blocks.

// Called records that a method was invoked with the given arguments and
// returns the configured return values as a [Arguments] slice. It panics if
// the call is unexpected, meaning no matching [Mock.On] or [Mock.OnAny]
// expectation was set.
//
// If the expectation uses [Call.Until] or [Call.After], this method blocks
// until the specified condition is met.
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

// Call calls method on the mock with arguments and returns the mocked method
// [Arguments]. Panics if the call is unexpected (i.e. not preceded by
// appropriate [Mock.On] calls). Blocks before returning if [Call.Until] or
// [Call.After] were used.
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

// Callable finds a callable method with given name and matching arguments.
// When found it returns it, otherwise it returns an error describing the
// reason. Note that there may be more methods in the expected slice matching
// the criteria in which case the first one is returned.
//
// A callable method is one that returns no error from [Call.CanCall] method,
// and has matching arguments.
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
		if call.argsAny && err == nil {
			return call, nil
		}
		if _, cnt := call.args.Diff(args); cnt == 0 {
			if err == nil {
				return call, nil
			}
		}
	}

	// Find a proxy method which was added without need for matching arguments.
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

	// Try to find method that is most similar to the one we are processing.
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

// Failed returns true if [Mock] instance is in filed state, false otherwise.
func (mck *Mock) Failed() bool {
	mck.mx.Lock()
	defer mck.mx.Unlock()
	mck.t.Helper()
	return mck.failed
}

// Unset removes [Call] instance from expected [Mock] calls. If the instance
// doesn't exist, it will trigger a test failure.
//
// Example usage:
//
//	call := Mock.On("Method", mock.Any)
//	Mock.Unset(call).Unset()
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

// AssertExpectations asserts that everything specified with [Mock.On] and
// [Call.Return] was in fact called as expected. Calls may have occurred in any
// order.
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

		name := formatMethod(call.Method, call.args, call.returns)
		names = append(names, name)

		format := "expected %d %s received %d %s"
		why := fmt.Sprintf(format, call.wantCalls, wCls, call.haveCalls, hCls)
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

// AssertCallCount asserts the method was called "want" number of times.
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
