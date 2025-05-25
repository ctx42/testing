// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

// Call represents a mocked method call, used for defining expectations.
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

	// Guards the fields.
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
func (c *Call) withParent(parent *Mock) *Call {
	c.parent = parent
	return c
}

// withStack sets the stack trace where the call requirement was defined.
func (c *Call) withStack(stack []string) *Call {
	c.Stack = stack
	return c
}

// With sets argument requirements for the proxied method call.
//
// Example usage:
//
//	Mock.Proxy(obj.Method).With(1, 2, 3)
func (c *Call) With(args ...any) *Call {
	if !c.proxy.IsValid() {
		panic("cannot set arguments on proxy calls")
	}
	c.args = args
	return c
}

// Return specifies the expected return arguments when the mocked method is
// invoked. If the [Call] instance represents a proxy call, this method will
// panic, as proxy calls cannot have predefined return values.
//
// Example usage:
//
//	Mock.On("Method").Return(errors.New("operation failed"))
func (c *Call) Return(args ...any) *Call {
	if c.proxy.IsValid() {
		panic("proxy calls cannot have return values")
	}
	c.returns = args
	return c
}

// Panic sets the panic value for the mock call and ensures that the method
// call will cause a panic during execution. This is useful for simulating
// failure scenarios where a panic is expected during a method call. It will
// panic if the [Call] instance represents a proxy call.
//
// Example usage:
//
//	Mock.On("MethodName").Panic("test panic")
//	Mock.On("MethodName").Panic(errors.New("test panic"))
func (c *Call) Panic(value any) *Call {
	if c.proxy.IsValid() {
		panic("cannot call panic on proxy calls")
	}
	c.panic = value
	return c
}

// Once specifies that the mocked method should only be called once with this
// set of arguments and return values. It is a helper calling [Call.Times].
//
// Example usage:
//
//	Mock.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Once()
func (c *Call) Once() *Call { return c.Times(1) }

// Times specify that the mocked method should only be called i times with
// this set of arguments and return values.
//
//	Mock.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Times(5)
func (c *Call) Times(i int) *Call {
	if c.optional {
		panic("cannot use Optional and Times in the same time")
	}
	c.wantCalls = i
	return c
}

// Until sets a channel that will block the return of the method call until the
// channel is either closed or a message is received. This is useful for
// scenarios where delaying the method's return until certain conditions are
// met is required.
//
// Example usage:
//
//	Mock.On("Method", arg1, arg2).Until(time.After(time.Second))
func (c *Call) Until(ch <-chan time.Time) *Call {
	c.until = ch
	return c
}

// After sets how long to block until the call returns.
//
//	Mock.On("Method", arg1, arg2).After(time.Second)
func (c *Call) After(d time.Duration) *Call {
	c.after = d
	return c
}

// Alter sets functions to be called on arguments received by the mocked method
// before they are returned. It can be used when mocking a method (such as an
// unmarshaler) that takes a pointer to a struct and sets properties in such
// struct.
//
// Example usage:
//
//		Mock.On("Unmarshal", MatchOfType("map[string]any")).
//	 	Return().
//	 	Alter(func(args Arguments) {
//		    arg := args.Get(0).(*map[string]any)
//		    arg["foo"] = "bar"
//		})
//
// If [Call.After] or [Call.Until] methods are used, the provided function is
// run after the blocking is done.
func (c *Call) Alter(fn ...func(Arguments)) *Call {
	c.alter = append(c.alter, fn...)
	return c
}

// Optional allows the method call to be optional, meaning that the method
// may or may not be invoked without causing test failures. This can be
// useful for scenarios where some method calls are not strictly required.
// Invoking the method after [Call.Times] will cause a panic as the two are
// mutually exclusive.
//
// Example usage:
//
//	Mock.On("MethodName").Optional()
func (c *Call) Optional() *Call {
	if c.wantCalls > 0 {
		panic("cannot use Optional and Times in the same time")
	}
	c.optional = true
	return c
}

// Requires specifies that the current call depends on other calls being
// satisfied. It's used to define prerequisite calls that need to be executed
// before this call. The referenced calls may be from the same mock instance
// and/or other mock instances.
//
// Example usage:
//
//	other := Mock.On("Init").Return(nil)
//	Mock.On("Do").Return(nil).Requires(other)
func (c *Call) Requires(calls ...*Call) *Call {
	for _, call := range calls {
		if call == nil {
			panic("a nil instance of mock.Call passed to mock.Call.Requires")
		}
	}
	c.requires = append(c.requires, calls...)
	return c
}

// CanCall provided the argument match returns nil if the method can be called
// again, otherwise it returns error explaining why method cannot be called.
func (c *Call) CanCall() error {
	err := c.satisfied(c.haveCalls + 1)
	if err == nil ||
		errors.Is(err, ErrNeverCalled) || errors.Is(err, ErrTooFewCalls) {
		return nil
	}
	return err
}

// Satisfied returns true if the method call requirements are met.
//
// A method call is considered satisfied if:
//   - It has been called the exact number of requested times.
//   - It has no maximum call limit and has been called at least once.
//   - It has not been called, but is marked as optional.
func (c *Call) Satisfied() bool {
	return c.satisfied(c.haveCalls) == nil
}

// satisfied returns nil if the call requirements are satisfied. It takes
// haveCalls instead of using instance field value, so it can be used to check
// if it is ok to call it one more time see [Call.CanCall].
func (c *Call) satisfied(haveCalls int) error {
	if c.wantCalls == haveCalls && c.wantCalls != 0 {
		return nil // Called requested number of times.
	}

	if c.wantCalls == 0 {
		if haveCalls == 0 {
			if c.optional {
				return nil // Optional and never called.
			}
			method := formatMethod(c.Method, c.args, c.returns)
			return notice.New(hNeverCalled).
				Append("method", "%s", method).
				Append("expected args", "\n%s", formatArgs(c.args)).
				Wrap(ErrNeverCalled)
		}
		return nil // Can be called any number of times.
	}

	if c.wantCalls < haveCalls {
		method := formatMethod(c.Method, c.args, c.returns)
		msg := notice.New(hTooManyCalls).
			Append("method", "%s", method)
		if len(c.args) > 0 {
			_ = msg.Append("expected args", "\n%s", formatArgs(c.args))
		}
		return msg.Append("want calls", "%d", c.wantCalls).
			Append("have calls", "%d", haveCalls).
			Wrap(ErrTooManyCalls)
	}

	if c.optional {
		return nil // Can be called up to wantCalls times.
	}

	method := formatMethod(c.Method, c.args, c.returns)
	return notice.New(hTooFewCalls).
		Append("method", "%s", method).
		Append("expected args", "\n%s", formatArgs(c.args)).
		Append("want calls", "%d", c.wantCalls).
		Append("have calls", "%d", haveCalls).
		Wrap(ErrTooFewCalls)
}

// checkReq verifies that all prerequisites for a method call are met. The
// stack parameter should contain the stack trace from where the method was
// invoked.
func (c *Call) checkReq(cs []string) (err error) {
	for _, req := range c.requires {
		if req.Satisfied() {
			continue
		}

		parent := "the same mock"
		if req.parent != c.parent {
			parent = "a different mock"
		}

		method := formatMethod(c.Method, c.args, c.returns)
		requires := formatMethod(req.Method, req.args, req.returns)

		msg := notice.New(hUnexpectedCall).Append("method", "%s", method)
		if len(c.args) > 0 {
			_ = msg.Append("expected args", "\n%s", formatArgs(c.args))
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
// configured return values.
func (c *Call) call(args ...any) Arguments {
	c.haveCalls++
	if c.until != nil {
		<-c.until
	} else {
		time.Sleep(c.after)
	}
	if c.panic != nil {
		panic(c.panic)
	}

	for _, fn := range c.alter {
		fn(args)
	}
	if c.proxy.IsValid() {
		return c.callProxy(args...)
	}
	return c.returns
}

// callProxy calls proxy method with given arguments. Panics if the proxy
// method has not been defined.
func (c *Call) callProxy(args ...any) Arguments {
	if !c.proxy.IsValid() {
		panic("proxy method not found")
	}

	argVals := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		argVals = append(argVals, reflect.ValueOf(arg))
	}

	retArgs := c.proxy.Call(argVals)
	var returns []any
	for _, ret := range retArgs {
		returns = append(returns, ret.Interface())
	}
	return returns
}

// satisfy marks the call as satisfied.
func (c *Call) satisfy() *Call {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.wantCalls == 0 {
		c.haveCalls = 1
	} else {
		c.haveCalls = c.wantCalls
	}
	return c
}

// End ends a chain of calls on the [Call] instance and returns parent [Mock].
func (c *Call) End() *Mock { return c.parent }
