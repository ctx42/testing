<!-- TOC -->
* [The `check` Package](#the-check-package)
  * [Example Usage](#example-usage)
  * [Custom Assertions](#custom-assertions)
<!-- TOC -->

# The `check` Package

The `check` package is designed for performing assertions in Go tests,
particularly as a foundational layer for the `assert` package. It provides
functions that return errors instead of boolean values, allowing callers to
adjust error messages to a particular context, add more contextual information
about the check, improving assertion message comprehension.

## Example Usage

You use checks like any other function returning error.

<!-- gmdoceg:ExampleNoError_withTrail -->
```go
have := errors.New("test error")

err := check.NoError(have, check.WithTrail("type.field"))

fmt.Println(err)
// Output:
// expected the error to be nil:
//   trail: type.field
//    want: nil
//    have: "test error"
```

The main purpose of returning an error from a check, instead of true false like
it is in case of `assert` package is to allow user to customize the message
and/or add context.

<!-- gmdoceg:ExampleNoError_changeMessage -->
```go
have := errors.New("test error")

err := check.NoError(have, check.WithTrail("type.field"))

err = notice.From(err, "prefix").Append("context", "wow")

fmt.Println(err)
// Output:
// [prefix] expected the error to be nil:
//     trail: type.field
//      want: nil
//      have: "test error"
//   context: wow
```

## Custom Assertions

See example in [custom_assertions_test.go](custom_assertions_test.go) file.