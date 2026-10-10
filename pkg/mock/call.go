// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ctx42/testing/pkg/notice"
)

// stack represents a method call stack.
type cStack struct {
	// The name of the method that was or will be called.
	Method string

	// The call stack of the method. What kind of stack it is context-dependent.
	Stack []string
}

// Call represents a single expected method invocation and its configuration.
//
// Call instances are returned by [Mock.On], [Mock.OnAny] and [Mock.Proxy] and
// form a fluent builder for expectations:
//
//	Mock.On("Method", 1, mock.Any).Return("x").Times(3).Requires(otherCall)
//
// The setters are not synchronized: configure a Call before the mock is used
// from other goroutines.
type Call struct {
	// Method name and call stack where the method call requirement was defined.
	cStack

	// Mock instance the call belongs to.
	// It is set during construction time and never changed.
	parent *Mock

	// Expected method arguments.
	args Arguments

	// Lets a method be called with any number of arguments as long as the
	// method can be called - [Call.CanCall] returns true.
	argsAny bool

	// Arguments to return when this method is called.
	returns Arguments

	// Maximum number of times the method can be called. Zero means no
	// restrictions.
	wantCalls int

	// Number of times the method has been called.
	haveCalls int

	// Call to this method is optional.
	optional bool

	// Will block returning from [Mock.Call] until it either
	// receives a message or is closed. If nil, it returns immediately.
	until <-chan time.Time

	// Set this to a non-zero value to block returning from [Mock.Call]
	// for a given period of time.
	after time.Duration

	// Change arguments passed to the mocked method during its execution. The
	// functions are called on the arguments right before returning.
	alter []func(Arguments)

	// If it's set to a non-nil value, the [Call] will panic with the given
	// value just before returning.
	panic any

	// Calls which must be satisfied before this one. Satisfied calls are the
	// ones that return true from the [Call.Satisfied] method.
	requires []*Call

	// The actual method to call.
	proxy reflect.Value

	// Guards haveCalls, which is read by other mocks checking [Call.Requires]
	// prerequisites without holding this call's parent mock lock.
	mx sync.Mutex
}

// newCall returns a new Call instance for a method with the specified name and
// expected arguments.
func newCall(method string, args ...any) *Call {
	return &Call{
		cStack: cStack{Method: method},
		args:   args,
	}
}

// newProxy returns a Call instance representing a proxy call. Optionally, if
// name is non-empty, its first value will be used as the custom proxied method
// name.
func newProxy(proxy reflect.Value, name ...string) *Call {
	var metName string
	if len(name) == 0 {
		metName = methodName(proxy)
	} else {
		metName = name[0]
	}

	call := newCall(metName)
	call.proxy = proxy
	return call
}

// withParent sets [Mock] instance the call belongs to.
func (call *Call) withParent(parent *Mock) *Call {
	call.parent = parent
	return call
}

// withStack sets the stack trace where the call requirement was defined.
func (call *Call) withStack(stack []string) *Call {
	call.Stack = stack
	return call
}

// With sets the expected arguments for a proxied call (see [Mock.Proxy]).
// Without With, proxy calls accept any arguments.
func (call *Call) With(args ...any) *Call {
	if !call.proxy.IsValid() {
		panic("cannot set arguments on non-proxy calls")
	}
	call.args = args
	return call
}

// Return sets the values that will be returned when the mocked method is
// later invoked. It panics if called on a proxy call (proxies forward real
// return values).
func (call *Call) Return(args ...any) *Call {
	if call.proxy.IsValid() {
		panic("proxy calls cannot have return values")
	}
	call.returns = args
	return call
}

// Panic arranges for the mocked method to panic with the given value when
// invoked. Useful for testing error paths that expect panics. It panics if
// called on a proxy call.
func (call *Call) Panic(value any) *Call {
	if call.proxy.IsValid() {
		panic("cannot call panic on proxy calls")
	}
	call.panic = value
	return call
}

// Once is a shorthand for [Call.Times](1): the expectation must be satisfied
// exactly once. Mutually exclusive with [Call.Optional].
func (call *Call) Once() *Call { return call.Times(1) }

// Times sets the exact number of times this expectation must be met.
// Mutually exclusive with [Call.Optional].
func (call *Call) Times(i int) *Call {
	if call.optional {
		panic("cannot use Optional and Times at the same time")
	}
	call.wantCalls = i
	return call
}

// Until makes the mocked method block until the channel is closed or
// receives a value. Useful for testing timeouts, cancellation, or ordering.
func (call *Call) Until(ch <-chan time.Time) *Call {
	call.until = ch
	return call
}

// After makes the mocked method sleep for the given duration before
// returning (simpler fixed-delay alternative to [Call.Until]).
func (call *Call) After(d time.Duration) *Call {
	call.after = d
	return call
}

// Alter registers functions to be called with the received arguments
// immediately before the mock returns (or after any delay). Commonly used
// to mutate pointer arguments before the real implementation sees them.
func (call *Call) Alter(fn ...func(Arguments)) *Call {
	call.alter = append(call.alter, fn...)
	return call
}

// Optional marks the expectation as optional: zero or more calls are
// acceptable and will not cause AssertExpectations to fail. Mutually
// exclusive with [Call.Times] / [Call.Once].
func (call *Call) Optional() *Call {
	if call.wantCalls > 0 {
		panic("cannot use Optional and Times at the same time")
	}
	call.optional = true
	return call
}

// Requires declares that the listed calls must be satisfied before this
// expectation can be met. Prerequisites can be on the same or different mocks.
func (call *Call) Requires(calls ...*Call) *Call {
	for _, req := range calls {
		if req == nil {
			panic("a nil instance of mock.Call passed to mock.Call.Requires")
		}
	}
	call.requires = append(call.requires, calls...)
	return call
}

// CanCall reports whether this expectation can be satisfied by one more call
// right now. Returns nil if allowed, otherwise one of the Err* sentinels or
// a richer [notice.Notice] explaining the violation.
func (call *Call) CanCall() error {
	err := call.satisfied(call.calls() + 1)
	if err == nil ||
		errors.Is(err, ErrNeverCalled) || errors.Is(err, ErrTooFewCalls) {
		return nil
	}
	return err
}

// Satisfied reports whether this expectation has been fully met (correct
// call count + all prerequisites satisfied).
func (call *Call) Satisfied() bool {
	return call.satisfied(call.calls()) == nil
}

// calls returns the number of recorded calls.
func (call *Call) calls() int {
	call.mx.Lock()
	defer call.mx.Unlock()
	return call.haveCalls
}

// addCall records one call.
func (call *Call) addCall() {
	call.mx.Lock()
	defer call.mx.Unlock()
	call.haveCalls++
}

// satisfied returns nil if the call requirements are satisfied. It takes
// haveCalls instead of using instance field value, so it can be used to check
// if it is ok to call it one more time see [Call.CanCall].
func (call *Call) satisfied(haveCalls int) error {
	if call.wantCalls == haveCalls && call.wantCalls != 0 {
		return nil // Called requested number of times.
	}

	if call.wantCalls == 0 {
		if haveCalls == 0 {
			if call.optional {
				return nil // Optional and never called.
			}
			method := formatMethod(call.Method, call.args, call.returns)
			return notice.New(hNeverCalled).
				Append("method", "%s", method).
				Append("expected args", "\n%s", formatArgs(call.args)).
				Wrap(ErrNeverCalled)
		}
		return nil // Can be called any number of times.
	}

	if call.wantCalls < haveCalls {
		method := formatMethod(call.Method, call.args, call.returns)
		msg := notice.New(hTooManyCalls).
			Append("method", "%s", method)
		if len(call.args) > 0 {
			_ = msg.Append("expected args", "\n%s", formatArgs(call.args))
		}
		return msg.Append("want calls", "%d", call.wantCalls).
			Append("have calls", "%d", haveCalls).
			Wrap(ErrTooManyCalls)
	}

	if call.optional {
		return nil // Can be called up to wantCalls times.
	}

	method := formatMethod(call.Method, call.args, call.returns)
	return notice.New(hTooFewCalls).
		Append("method", "%s", method).
		Append("expected args", "\n%s", formatArgs(call.args)).
		Append("want calls", "%d", call.wantCalls).
		Append("have calls", "%d", haveCalls).
		Wrap(ErrTooFewCalls)
}

// checkReq verifies that all prerequisites for a method call are met. The
// stack parameter should contain the stack trace from where the method was
// invoked.
func (call *Call) checkReq(cs []string) (err error) {
	for _, req := range call.requires {
		if req.Satisfied() {
			continue
		}

		parent := "the same mock"
		if req.parent != call.parent {
			parent = "a different mock"
		}

		method := formatMethod(call.Method, call.args, call.returns)
		requires := formatMethod(req.Method, req.args, req.returns)

		msg := notice.New(hUnexpectedCall).Append("method", "%s", method)
		if len(call.args) > 0 {
			_ = msg.Append("expected args", "\n%s", formatArgs(call.args))
		}
		_ = msg.Append("requires", "%s", requires)
		if len(req.args) > 0 {
			_ = msg.Append(" expected args", "\n%s", formatArgs(req.args))
		}
		_ = msg.Append("from", "%s", parent).Wrap(ErrRequirements)
		if len(cs) > 0 {
			_ = msg.Append("stack", "%s", strings.Join(cs, "\n"))
		}
		err = notice.Join(err, msg)
	}
	return err
}

// call represents a call to the mocked method with arguments. Returns
// configured return values. [Mock.Call] increments haveCalls under its mutex
// before invoking this so concurrent callers do not race on the counter.
func (call *Call) call(args ...any) Arguments {
	if call.until != nil {
		<-call.until
	} else {
		time.Sleep(call.after)
	}
	if call.panic != nil {
		panic(call.panic)
	}

	for _, fn := range call.alter {
		fn(args)
	}
	if call.proxy.IsValid() {
		return call.callProxy(args...)
	}
	return call.returns
}

// callProxy calls proxy method with given arguments. Panics if the proxy
// method has not been defined.
//
// A nil interface argument is converted with [reflect.Zero] of the
// corresponding parameter type so nil error, pointer, and similar args work.
func (call *Call) callProxy(args ...any) Arguments {
	if !call.proxy.IsValid() {
		panic("proxy method not found")
	}

	typ := call.proxy.Type()
	numIn := typ.NumIn()
	variadic := typ.IsVariadic()
	argVals := make([]reflect.Value, 0, len(args))
	for i, arg := range args {
		if arg != nil {
			argVals = append(argVals, reflect.ValueOf(arg))
			continue
		}
		var inTyp reflect.Type
		if variadic && i >= numIn-1 {
			inTyp = typ.In(numIn - 1).Elem()
		} else {
			inTyp = typ.In(i)
		}
		argVals = append(argVals, reflect.Zero(inTyp))
	}

	retArgs := call.proxy.Call(argVals)
	var returns []any
	for _, ret := range retArgs {
		returns = append(returns, ret.Interface())
	}
	return returns
}

// satisfy marks the call as satisfied.
func (call *Call) satisfy() *Call {
	call.mx.Lock()
	defer call.mx.Unlock()
	if call.wantCalls == 0 {
		call.haveCalls = 1
	} else {
		call.haveCalls = call.wantCalls
	}
	return call
}

// End terminates a fluent expectation chain and returns the parent [Mock].
// Useful when you want to continue configuring the same mock after setting
// up one call:
//
//	mck.On("A").Return(1).End().On("B").Return(2)
func (call *Call) End() *Mock { return call.parent }
