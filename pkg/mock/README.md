Package `mock` provides structures for developing interface mocks.

<!-- TOC -->
* [Introduction](#introduction)
* [Basic Usage](#basic-usage)
* [Defining Expectations](#defining-expectations)
  * [Argument Expectations](#argument-expectations)
    * [Matching Anything](#matching-anything)
    * [Matching Type](#matching-type)
    * [Predefined Matchers](#predefined-matchers)
    * [Custom Matchers](#custom-matchers)
  * [Return Value Expectations](#return-value-expectations)
  * [Panicking](#panicking)
  * [Expecting Number Of Calls](#expecting-number-of-calls)
  * [Block Return](#block-return)
    * [Using Timout](#using-timout)
    * [Using Channel](#using-channel)
  * [Alter Passed Arguments](#alter-passed-arguments)
  * [Optional Calls](#optional-calls)
<!-- TOC -->

# Introduction

Mocking is a process used in unit testing when the unit being tested has
external dependencies. The purpose of mocking is to isolate and focus on
the code being tested and not on the behavior or state of external
dependencies. In mocking, the dependencies are replaced by closely controlled
replacements objects that simulate the behavior of the real ones. There are
three main possible types of replacement objects - fakes, stubs and mocks.

> **Fakes**: A Fake is an object that will replace the actual code by
implementing the same interface but without interacting with other objects.
Usually the Fake is hard-coded to return fixed results. To test for different
use cases, a lot of Fakes must be introduced. The problem introduced by using
Fakes is that when an interface has been modified, all fakes implementing this
interface should be modified as well.

> **Stubs**: A Stub is an object that will return a specific result based on a
specific set of inputs and usually won’t respond to anything outside what is
programed for the test. With `mock` you can create a Stub in a test with
a minimal amount of code, making it clear how the dependency will respond and
how the tested system should behave.

> **Mocks**: A Mock is a much more sophisticated version of a Stub. It will 
still turn values like a Stub, but it can also be programmed with expectations 
in terms of how many times each method should be called, in which order and 
with what data. With `mock` you can create a Mock with just few lines of code,
which makes the test more understandable.

# Basic Usage

Let's start with a [basic example](../../examples/mock_test.go) where we have 
simple interface we would like to mock.

```go
package examples

import (
	"testing"

    "github.com/ctx42/testing/pkg/mock"
)

// Adder is the simple interface we would like to mock.
type Adder interface {
	Add(a, b float64) float64
}

// AdderMock implements Adder interface. By convention the name of the type
// should be the name of the interface it is mocking with "Mock" suffix.
type AdderMock struct {
	*mock.Mock // Embedded instance.

	// If you need additional fields you can declare them here.
}

// NewAdderMock is the constructor function we will use to instantiate the mock.
// By convention, it should be the name of the type prefixed with "New".
func NewAdderMock(t *testing.T) *AdderMock {
	return &AdderMock{mock.NewMock(t)}
}

// Add mocks the method from Adder interface.
func (_mck *AdderMock) Add(a, b float64) float64 {
	// Inform the mock the method was called with given arguments.
	// The call returns [mock.Arguments] representing return values.
	args := _mck.Called(a, b)

	// Here you can do additional logic if needed.

	// Get the first return value, cast it to expected type, and return it.
	return args.Get(0).(float64)
}

// Test_Adder_Add is ena example test case using the AdderMock.
func Test_Adder_Add(t *testing.T) {
	// --- Given ---
	mck := NewAdderMock(t) // Instantiate the mock.
	mck.
		On("Add", 1.0, 2.0). // Define method and argument expectations.
		Return(3.0)          // Define return values expectations.

	// --- When ---

	// In real example the mock created above would be used in code requiring
	// Adder interface, which in turn would call Add method like below.
	have := mck.Add(1.0, 2.0)

	// --- Then ---
	// Below line will print: Result: 3.000000
	t.Logf("Result: %f", have)
}
```

# Defining Expectations

After instantiating the mock you usually want to add expectations how the 
Code Under Test (CUT) will use the mock. You do that with `Mock.On` method on 
the mock, like in the example above. For argument expectations in the `Add` 
method you can either use exact values or argument matchers. The `Mock.On` 
method returns instance of `mock.Call` on which is used to define expectations 
for this specific method call.

## Argument Expectations

In the basic example we defined that method `Add` should be called with two 
values _1.0_ and _2.0_.

```go
Mock.On("Add", 1.0, 2.0)
```

Call to this mocked method by CUT with any other values will trigger a test 
error.

### Matching Anything

If you don't know the exact value, or you do not care about it, you can use 
special argument matcher `mock.Any`.

```go
Mock.On("Add", 1.0, mock.Any)
```

This will make the mock to ignore its value when validating the call to the 
mocked method.

### Matching Type

To only match argument type use `mock.MatchOfType` matcher.

```go
Mock.On("MethodName", mock.MatchOfType("int"))
Mock.On("MethodName", mock.MatchOfType("*int"))
Mock.On("MethodName", mock.MatchOfType("*http.Request"))
```

### Predefined Matchers

To make your life easier the `mock` package defines few most used matchers:

- `mock.AnyString` - matches any string.
- `mock.AnyInt` - matches any integer value.
- `mock.AnyBool` - matches any boolean value.
- `mock.MatchOfType` - matches argument type.
- `mock.MatchType` - matches argument type.
- `mock.MatchErrorContain` - matches argument is an error which contains given string.
- `mock.MatchErr` - matches argument is an error and has sentinel error in its chain.
- `mock.MatchError` - matches argument is an error and its message is equal given string.

### Custom Matchers

If you need more sophisticated matchers you can easily define your own using 
`mock.MatchBy` function. In fact come of the matchers listed in 
[Predefined Matchers](#predefined-matchers) are defined using it. For example
`mock.MatchErrorContain` is defined like this:

```go
// MatchErrorContain constructs an argument matcher (ArgMatcher) instance
// which ensures argument is a non nil error with given message.
func MatchErrorContain(want string) *ArgMatcher {
    return MatchBy(func(err error) bool {
        msg := err.Error()
        return strings.Contain(msg, want)
    })
}
```

See `mock.MatchBy` documentation for more details.

## Return Value Expectations

If the mocked method has returns we need to define them. In that case each call
to `On` method need to call `Return` to define return values.

```go
Mock.On("Add", 1.0, 2.0).Return(3.0)
Mock.On("Add", 2.0, mock.Any).Return(4.0)
```

## Panicking

You may ask the mock to panic with given message.

```go
Mock.On("Method").Panic("test panic")
```

## Expecting Number Of Calls

By default, if you create an expectation for a method it can be called by CUT
unrestricted number of times. To restrict the number of calls use `Times` 
method.

```go
Mock.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Times(5)
```

For convenience there is also a method `Once` which is equivalent to `Times(1)`.

```go
Mock.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Once()
```

## Block Return

### Using Timout

You can block the mocked function for given length of time before it returns. 

```go
Mock.On("Method", arg1, arg2).After(time.Second)
```

Using this feature you can test latencies.

### Using Channel

This is similar to [Using Timout](#using-timout) but some outside code can make 
decision when to unlock the mocked method.

```go
ch := make(chan struct{})
Mock.On("Method", arg1, arg2).Until(ch)
```

## Alter Passed Arguments

You can run a function before returning from the mock. It's very useful if you 
need to modify arguments in some way before returning.

If After or Until methods are used the provided function is run after the 
blocking is done.

```go
Mock.
    On("Unmarshal", MatchOfType("map[string]any")).
        Return().
        Alter(func(args Arguments) {
            arg := args.Get(0).(*map[string]any)
            arg["foo"] = "bar"
        })
```

## Optional Calls

You may set the mocked method call as _optional_. Optional methods when not 
called do not produce an error.

```go
Mock.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Optional()
```