<!-- TOC -->
* [Tester Package](#tester-package)
* [Test Manager Interface](#test-manager-interface)
  * [Usage](#usage)
* [Spy](#spy)
  * [Testing Test Helpers](#testing-test-helpers)
  * [Setting Spy Expectations](#setting-spy-expectations)
    * [Expectations For `Helper`](#expectations-for-helper)
    * [Executing Cleanup Functions](#executing-cleanup-functions)
    * [Checking Spy State](#checking-spy-state)
    * [Get TempDir Paths](#get-tempdir-paths)
    * [Examine Log Messages](#examine-log-messages)
    * [Ignore Log Messages](#ignore-log-messages)
  * [Examples](#examples)
<!-- TOC -->

# Tester Package

If you’ve spent any time writing Go tests, you’ve probably encountered the joy
of `*testing.T`. It’s the backbone of Go’s testing framework — powerful,
flexible, and ubiquitous. But as your test suite grows, you might find yourself
repeating the same chunks of test logic across multiple test cases. Enter _test
helpers_: reusable functions that streamline your tests, improve readability,
and reduce complexity. Libraries like assert are prime examples, turning verbose
checks into concise assertions.

But here’s the catch: how do you test the test helpers themselves? After all,
these are the tools you rely on to ensure your code works as expected. If they
fail, your tests might silently lie to you. This is where the `tester` package
comes to the rescue.

Package `tester` provides interface `T` which is a subset of `testing.TB`
interface and `Spy` struct which helps with testing test helpers.

# Test Manager Interface

The goal of `T` interface is to make testing of _test helpers_ possible. We 
define _test helper_ as code which uses `*testing.T` instances as 
_test manager_. By design `T` is a subset of the `testing.TB` interface to 
allow using implementers as well as `*testing.T` as a _test helper_ argument.    

Creating test helpers is part of making tests more readable. Instead of 
repeating big blocks of code in many test cases we can create a helper, and 
delegate part of testing procedures to it. A test helper usually receives
some kind of test manager instance (usually `*testing.T`) as an argument, so 
it can log and provide test outcome to the test runner. 

Very good example of test helpers are assertion functions in `assert` package, 
which improve test readability and in many cases reduce their complexity.

## Usage

Anywhere where `*testing.T` is used you can replace it with `tester.T` interface 
as long as the test helper uses the following methods:

- Spy.Cleanup(func())
- Spy.Error(args ...any)
- Spy.Errorf(format string, args ...any)
- Spy.Fatal(args ...any)
- Spy.Fatalf(format string, args ...any)
- Spy.FailNow()
- Spy.Failed() bool
- Spy.Helper()
- Spy.Log(args ...any)
- Spy.Logf(format string, args ...any)
- Spy.Name() string
- Spy.Setenv(key, value string)
- Spy.Skip(args ...any)
- Spy.TempDir() string
- Spy.Context() context.Context

So for example a test helper

```go
// IsOdd asserts "have" is odd number. Returns true if it is, otherwise marks
// the test as failed, writes an error message to the test log and returns false.
func IsOdd(t *testing.T, have int) bool {
	t.Helper()
	if have%2 != 0 {
		t.Errorf("expected %d to be odd", have)
		return false
	}
	return true
}
```

can be refactored as

```go
// IsOdd asserts "have" is odd number. Returns true if it is, otherwise marks
// the test as failed, writes an error message to the test log and returns false.
func IsOdd(t tester.T, have int) bool {
	t.Helper()
	if have%2 != 0 {
		t.Errorf("expected %d to be odd", have)
		return false
	}
	return true
}
```

without any change to the body of the function. Once you replace `*testing.T` 
with implementer of `tester.T` (for example `Spy` instance) you can create 
tests for the helper.

# Spy

The `Spy` type was designed to be a spy for `tester.TB` interface. The spy 
allows you to define expectations how the test manager instance is used by a 
test helper. 

## Testing Test Helpers

We can test the `IsOdd` test helper created above in the following way:

```go
func Test_IsOdd(t *testing.T) {
    t.Run("error is not odd number", func(t *testing.T) {
        // --- Given ---

		// Set up the spy with expectations
		tspy := tester.New(t)
		tspy.ExpectError()                              // Expect an error.
		tspy.ExpectLogEqual("expected %d to be odd", 2) // Expect log.
		tspy.Close()                                    // No more expectations.

		// --- When ---
		success := IsOdd(tspy, 2) // Run the helper.

		// --- Then ---
		if success { // Verify the outcome.
			t.Error("expected success to be false")
		}
		tspy.AssertExpectations() // Ensure all expectations were met.
	})

	t.Run("success is odd number", func(t *testing.T) {
		// Given
		tspy := tester.New(t)
		tspy.Close()

		// When
		success := IsOdd(tspy, 3)

		// Then
		if !success {
			t.Error("expected success to be true")
		}

		// The `tspy.AssertExpectations()` is called automatically.
	})
}
```

## Setting Spy Expectations

To set expectations for the Helper Under Test (`HUT`) `Spy` instance provides 
multiple `Expect*` methods.

```go
tspy := tester.New(t)

tspy.ExpectCleanups(n)  // Expect HUT to call Cleanup exactly n times. 
tspy.ExpectError()      // Expect HUT to call one of the Error* methods at least once. 
tspy.ExpectFatal()      // Expect HUT to call one of the Fatal* methods at least once.
tspy.ExpectFail()       // Expect HUT to call one of the Error* or Fatal* at least once.  
tspy.ExpectHelpers(n)   // Expect HUT to call Helper method exactly n times. 
tspy.ExpectSetenv(k, v) // Expect HUT to call Setenv method with key, value pair.
tspy.ExpectSkipped()    // Expect HUT to skip the test.
tspy.ExpectTempDir(n)   // Expect HUT to call TempDir n times.
tspy.ExpectFail()       // Expect HUT to call one of the Error* or Fatal* methods.
tspy.ExpectedNames(n)   // Expect HUT to call Name exactly n times.

// Log message expectations: 

tspy.ExpectLog(matcher, format, args...)  // Expect logged message to match formated string.
tspy.ExpectLogEqual(format, args...)      // Expect logged message to equal to formated string. 
tspy.ExpectLogContain(format, args...)    // Expect logged message to contain formated string.
tspy.ExpectLogNotContain(format, args...) // Expect logged message not to contain formated string.
```

Since each of the methods has very good documentation we encourage you to 
explore it for more details. Here we wil just document some of the cases which
might not be so obvious at the first glance.

### Expectations For `Helper`

By default, when you instantiate the `Spy` 

```go
tspy := tester.New(t)
```

it will expect _at least one_ call to `Helper` method, but you can define 
exact number of times it should be called by adding optional argument

```go
tspy := tester.New(t)
```

Now if the HUT does not make exactly 2 calls to `Helper` it will fail the test.

### Executing Cleanup Functions

Execute `Spy.Finish()` to run all registered cleanups.

```go
func Test_Spy_Cleanups(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t, 0)
	tspy.ExpectCleanups(1)
	tspy.Close()

	// --- When ---
	var have int
	tspy.Cleanup(func() { have = 42 })

	// --- Then ---
	tspy.Finish()
	if have != 42 {
		t.Errorf("expected 42 got %d", have)
	}
}
```

### Checking Spy State

At any point in time you may call `Spy.Failed()` to check if the HUT called 
any of the `Error*`, `Fatal*` or `FailNow` methods.

### Get TempDir Paths

To get paths generated by `Spy.TempDir` use `Spy.GetTempDir(idx)` where `idx` 
is an index into array of generated paths (zero indexed).

### Examine Log Messages

Calling methods like `Spy.Error*` not only change the state of the test being 
executed but also log messages (usually to standard output). For example:

```
=== RUN   Test_relativeTo
=== RUN   Test_relativeTo/current_package
    helpers_test.go:27: expected values to be equal:
        	want: "case"
        	have: "cases"
=== RUN   Test_relativeTo/not_current_package
=== RUN   Test_relativeTo/nil_package
--- FAIL: Test_relativeTo (0.28s)
    --- FAIL: Test_relativeTo/current_package (0.00s)

    --- PASS: Test_relativeTo/not_current_package (0.00s)
    --- PASS: Test_relativeTo/nil_package (0.00s)

FAIL
```

The `Spy` provides a couple of ways to examine log messages. The 
`ExpectLog(matcher MStrategy, format string, args ...any)` where you 
provide matching strategy: 

- `tester.Equal` - the log must be exact match with given formatted string
- `tester.Contains` - the log must contain given formatted string
- `tester.NotContains` - the log must NOT contain given formatted string
- `tester.Regexp` - the log must match regexp 

Or using convenience methods:

- `ExpectLogEqual`
- `ExpectLogContain`
- `ExpectLogNotContain`

which call `ExpectLog` with given matching strategy.

By default, if the HUT logged anything, and it was not examined the test will be 
failed. To change this behaviour use `Spy.IgnoreLogs` method. 

### Ignore Log Messages

```go
func Test_Spy_IgnoreLogExamination(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t, 0)
	tspy.ExpectError()
	// Without this line Spy will report an error 
	// that it did not expect the HUT to log.
	tspy.IgnoreLogs() 
	tspy.Close()

	// --- When ---
	tspy.Error("message")

	// --- Then ---
	tspy.AssertExpectations()
}
```

## Examples

See [tester.go](../../examples/tester.go) and [tester.go](../../examples/tester_test.go)
for more examples.
