<!-- TOC -->
* [Introduction](#introduction)
* [Basic Usage](#basic-usage)
* [Defining Expectations](#defining-expectations)
  * [Argument Matchers](#argument-matchers)
  * [Matching Any Value](#matching-any-value)
  * [Predefined Argument Matchers](#predefined-argument-matchers)
  * [Return Values](#return-values)
  * [Delaying Returns](#delaying-returns)
    * [Using a Timeout](#using-a-timeout)
    * [Using a Channel](#using-a-channel)
  * [Panicking](#panicking)
  * [Expecting Number of Calls](#expecting-number-of-calls)
  * [Modifying Arguments](#modifying-arguments)
  * [Optional Calls](#optional-calls)
* [Advanced Topics](#advanced-topics)
  * [Proxying Calls](#proxying-calls)
  * [Argument Matchers for Proxied Methods](#argument-matchers-for-proxied-methods)
  * [Custom Matchers](#custom-matchers)
<!-- TOC -->

# Introduction

Mocking is a unit testing technique that isolates the code under test from
external dependencies. By replacing dependencies with controlled objects that
simulate real behavior, mocks ensure tests focus solely on the logic being
tested, not on external systems. The mock package supports three types of
replacement objects: fakes, stubs, and mocks.

>Fakes: Fakes implement the same interface as the real dependency but operate
independently, returning hardcoded results. While simple, they are
inflexible—interface changes require updating all related fakes, which can be
cumbersome.

> Stubs: Stubs return predefined results for specific inputs and ignore
unexpected calls. With the `mock` package, you can define stubs concisely,
clearly specifying dependency behavior and expected system responses.

> Mocks: Mocks are advanced stubs that support detailed expectations, such as
call counts, call order, and argument values. The `mock` package enables 
creating mocks with minimal code, improving test readability and precision.

# Basic Usage

Consider a [basic example](../../examples/mock_test.go) with a simple interface
to mock:

```go
// Adder is a simple interface to mock.
type Adder interface {
	Add(a, b float64) float64
}

// AdderMock implements the Adder interface. By convention, mock types are
// named after the interface with a "Mock" suffix.
type AdderMock struct {
	*mock.Mock // Embedded mock instance.

	// Add custom fields here if needed.
}

// NewAdderMock creates a new [AdderMock] instance. By convention, constructor
// functions are prefixed with "New". More complex mocks may accept additional
// parameters.
func NewAdderMock(t *testing.T) *AdderMock {
	return &AdderMock{mock.NewMock(t)}
}

// Add implements the Add method from the [Adder] interface.
func (mck *AdderMock) Add(a, b float64) float64 {
	// Record the method call with its arguments, returning [mock.Arguments]
	// containing the defined return values.
	rets := mck.Called(a, b)

	// Add custom logic here if needed.

	// Extract and return the first return value, cast to float64.
	return rets.Get(0).(float64)
}

// Test_Adder_Add demonstrates using AdderMock in a test.
func Test_Adder_Add(t *testing.T) {
	// --- Given ---
	mck := NewAdderMock(t) // Create the mock.
	mck.
		On("Add", 1.0, 2.0). // Specify expected method and arguments.
		Return(3.0)          // Define the return value.

	// --- When ---
	// In a real test, the mock would be passed to code requiring the Adder
	// interface, which would invoke Add as shown.
	have := mck.Add(1.0, 2.0)

	// --- Then ---
	// Prints: Result: 3.000000
	t.Logf("Result: %f", have)
}
```

This creates a fully functional mock for use in tests:

```go
// Test_Adder_Add demonstrates using AdderMock in a test.
func Test_Adder_Add(t *testing.T) {
	// --- Given ---
	mck := NewAdderMock(t) // Create the mock.
	mck.
		On("Add", 1.0, 2.0). // Specify expected method and arguments.
		Return(3.0)          // Define the return value.

	// --- When ---
	// In a real test, the mock would be passed to code requiring the Adder
	// interface, which would invoke Add as shown.
	have := mck.Add(1.0, 2.0)

	// --- Then ---
	// Prints: Result: 3.000000
	t.Logf("Result: %f", have)
}
```

# Defining Expectations

After creating a mock, you define how the Code Under Test (CUT) interacts with
it using `Mock.On`, `Mock.OnAny`, or `Mock.Proxy`. These methods configure
expectations for method calls, including arguments and return values. Each
returns a `mock.Call` instance, which you use to refine the expectation.

## Argument Matchers

In the example, we specified that `Add` expects exact argument values `1.0` and 
`2.0`:

```go
mck.On("Add", 1.0, 2.0)
```

Calls with different values cause the test to fail.

## Matching Any Value

To ignore an argument’s value, use `mock.Any`:

```go
mck.On("Add", 1.0, mock.Any)
```

This accepts any value of the second argument during validation.

## Predefined Argument Matchers

The `mock` package provides common matchers for convenience:

- `mock.AnyString` – Matches any string.
- `mock.AnyInt` – Matches any integer.
- `mock.AnyBool` – Matches any boolean.
- `mock.MatchOfType` – Matches an argument’s type as string (e.g., "int", "*http.Request").
- `mock.MatchType` – Matches an argument’s type using `reflect` package.
- `mock.MatchError` – Matches a non-nil error with a specific message or error (via `errors.Is`).
- `mock.MatchErrorContain` – Matches a non-nil error containing a given substring.

## Return Values

For methods with return values, use `Call.Return`:

```go
mck.On("Add", 1.0, 2.0).Return(3.0)
mck.On("Add", 2.0, mock.Any).Return(4.0)
```

## Delaying Returns

### Using a Timeout

To delay a method’s return, use `Call.After`. This is useful for testing
latencies:

```go
mck.On("Method", arg1, arg2).After(time.Second).Return(1)
```

The method pauses for one second before returning.

### Using a Channel

To block until an external signal, use `Call.Until` with a channel:

```go
ch := make(chan struct{})
mck.On("Method").Until(ch).Return(1)
```

External code can close or send on ch to unblock the method.

## Panicking

To make a mock panic, use `Call.Panic`:

```go
mck.On("Method").Panic("test panic")
```

## Expecting Number of Calls

By default, expected methods can be called any number of times. To limit calls,
use `Call.Times`:

```go
mck.On("Method").Return("abc").Times(5)
```

For single calls, use `Call.Once` (equivalent to `Call.Times(1)`):

```go
mck.On("Method").Return(1).Once()
```

## Modifying Arguments

To modify arguments before returning, use `Call.Alter`:

```go
mck.On("Unmarshal", mock.MatchOfType("*map[string]any")).
	Return().
	Alter(func(args mock.Arguments) {
        arg := args.Get(0).(*map[string]any)
        (*arg)["foo"] = "bar"
    })
```

If `Call.After` or `Call.Until` is used, `Call.Alter` runs after the delay.

## Optional Calls

To mark a method call as optional (no error if uncalled), use `Call.Optional`:

```go
mck.On("Method").Return(1).Optional()
```

# Advanced Topics

## Proxying Calls

You can configure a mock to proxy method calls to an underlying implementation
using `Mock.Proxy`:

```go
obj := &types.TPtr{}

mck := NewItfMock(t)
mck.Proxy(obj.Method)
```

Here, calls to the mocked `Method` are forwarded to `obj.Method`. The mock
automatically detects the method’s name. To override the method name, provide 
an optional name argument:

```go
mck.Proxy(obj.Method, "OtherName")
```

## Argument Matchers for Proxied Methods

To validate arguments for proxied calls, use `Call.With`:

```go
obj := &types.TPtr{}

mck := NewItfMock(t)
mck.Proxy(obj.Method).With(1, "abc", 3)
```

This ensures the proxied method is called with the specified arguments.

## Custom Matchers

For advanced matching, define custom matchers with `mock.MatchBy`. Some 
predefined matchers, like `mock.MatchErrorContain`, are built this way:

```go
// MatchErrorContain creates a [Matcher] that verifies the argument is a
// non-nil error containing the specified substring in its message.
func MatchErrorContain(want string) *Matcher {
    return MatchBy(func(err error) bool {
        return strings.Contains(err.Error(), want)
    })
}
```

See the `mock.MatchBy` documentation for details.
