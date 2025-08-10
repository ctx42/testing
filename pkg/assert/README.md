<!-- TOC -->
* [The `assert` package](#the-assert-package)
  * [Assertions](#assertions)
    * [Asserting Structures](#asserting-structures)
    * [Asserting Recursive Structures](#asserting-recursive-structures)
    * [Asserting Maps, Arrays, and Slices](#asserting-maps-arrays-and-slices)
      * [Asserting Time](#asserting-time)
      * [Asserting JSON Strings](#asserting-json-strings)
      * [Worthy mentions](#worthy-mentions)
  * [Advanced usage](#advanced-usage)
    * [Custom Checkers](#custom-checkers)
    * [Understanding Trails](#understanding-trails)
    * [Registering Custom Type Checkers](#registering-custom-type-checkers)
    * [Registering Global Type Checkers](#registering-global-type-checkers)
    * [Skipping Fields, Elements, or Indexes](#skipping-fields-elements-or-indexes)
    * [Skipping unexported fields](#skipping-unexported-fields)
<!-- TOC -->

# The `assert` package

The `assert` package is a toolkit for Go testing that offers common assertions,
integrating well with the standard library. When writing tests, developers often
face a choice between using Go's standard `testing` package or packages like
`assert`. The standard library requires verbose `if` statements for assertions,
which can make tests harder to read. This package, on the other hand, provides
one-line asserts, such as `assert.NoError`, which are more concise and clear.
This simplicity helps quickly grasp the intent of each test, enhancing
readability.

By making tests easier to write and read, this package hopes to encourage
developers to invest more time in testing. Features like immediate feedback
with easily readable output and a wide range of assertion functions lower the
barrier to writing comprehensive tests. This can lead to better code coverage,
as developers are more likely to write and maintain tests when the process is
straightforward and rewarding.

## Assertions

Most of the assertions are self-explanatory, and I encourage you to see your
online [documentation](https://pkg.go.dev/github.com/ctx42/testing). Here we
will highlight only the ones that we feel are interesting.

### Asserting Structures

```go
type T struct {
    Int int
    Str string
}

have := T{Int: 1, Str: "abc"}
want := T{Int: 2, Str: "xyz"}

assert.Equal(want, have)
// Test Log:
//
// expected values to be equal:
//   trail: T.Int
//    want: 2
//    have: 1
//  ---
//   trail: T.Str
//    want: "xyz"
//    have: "abc"
```

### Asserting Recursive Structures

```go
type T struct {
    Int  int
    Next *T
}

have := T{1, &T{2, &T{3, &T{42, nil}}}}
want := T{1, &T{2, &T{3, &T{4, nil}}}}

assert.Equal(want, have)

// Test Log:
//
// expected values to be equal:
//   trail: T.Next.Next.Next.Int
//    want: 4
//    have: 42
```

### Asserting Maps, Arrays, and Slices

Maps

```go
type T struct {
    Str string
}

want := map[int]T{1: {Str: "abc"}, 2: {Str: "xyz"}}
have := map[int]T{1: {Str: "abc"}, 3: {Str: "xyz"}}

assert.Equal(want, have)

// Test Log:
//
// expected values to be equal:
//       trail: map[2]
//        want:
//              map[int]T{
//                1: {
//                  Str: "abc",
//                },
//                3: {
//                  Str: "xyz",
//                },
//              }
//        have: nil
//   want type: map[int]T
//   have type: <nil>
```

Slices and arrays

```go
want := []int{1, 2, 3}
have := []int{1, 2, 3, 4}

assert.Equal(want, have)

// Test Log:
//
// expected values to be equal:
//   want len: 3
//   have len: 4
//       want:
//             []int{
//               1,
//               2,
//               3,
//             }
//       have:
//             []int{
//               1,
//               2,
//               3,
//               4,
//             }
```

#### Asserting Time

```go
want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
have := time.Date(2025, 1, 1, 0, 1, 1, 0, time.UTC)

assert.Time(want, have)

// Test Log:
//
//  expected equal dates:
//   want: 2025-01-01T00:00:00Z
//   have: 2025-01-01T00:01:01Z
//   diff: -1m1s
```

#### Asserting JSON Strings

```go
want := `{"A": 1, "B": 2}`
have := `{"A": 1, "B": 3}`

assert.JSON(want, have)

// Test Log:
//
// expected JSON strings to be equal:
//   want: {"A":1,"B":2}
//   have: {"A":1,"B":3}
```

#### Worthy mentions

- `Epsilon` - assert floating point numbers within given ε.
- `ChannelWillClose` - assert channel will be closed within given time.
- `MapSubset` - checks the "want" is a subset "have".
- `Wait` - Wait waits for "fn" to return true but no longer then given timeout.

See the [documentation](https://pkg.go.dev/github.com/ctx42/testing) for the
full list.

## Advanced usage

### Custom Checkers

Custom checkers allow you to define specialized comparison logic for any type
or field trail in your Go tests. A custom checker is a function that matches
the `check.Check` signature, enabling fine-grained control over assertions.
Below is an example demonstrating how to create and use a custom checker.

```go
type T struct {
    Str string
    Any []any
}

chk := func(want, have any, opts ...any) error {
    wVal := want.(float64)
    hVal := want.(float64)
    return check.Epsilon(wVal, 0.01, hVal, opts...)
}
opt := check.WithTrailChecker("T.Any[1]", chk)

want := T{Str: "abc", Any: []any{1, 2.123, "abc"}}
have := T{Str: "abc", Any: []any{1, 2.124, "abc"}}

assert.Equal(want, have, opt)

// Test Log:
//
//  <nil>
```

In this example, the custom checker `chk` compares float64 values at the trail
`T.Any[1]` with a tolerance of 0.01. The assertion passes because 2.123 and
2.124 are within the specified epsilon.

Also, see the example in [custom_assertion_test.go](custom_assertions_test.go).

### Understanding Trails

A trail uniquely identifies a struct field, slice or array element, or map key
visited during an assertion. The `assert` package automatically tracks trails
for composite types, enabling precise targeting for custom checkers. To inspect
all visited trails, use the check.WithTrailLog option, as shown below:

```go
type T struct {
    Int  int
    Next *T
}

have := T{1, &T{2, &T{3, &T{42, nil}}}}
want := T{1, &T{2, &T{3, &T{42, nil}}}}
trails := make([]string, 0)

assert.Equal(want, have, check.WithTrailLog(&trails))

fmt.Println(strings.Join(trails, "\n"))
// Output:
// T.Int
// T.Next.Int
// T.Next.Next.Int
// T.Next.Next.Next.Int
// T.Next.Next.Next.Next
```

This output shows the hierarchical paths visited, including fields of nested
structs. Trails are essential for registering checkers at specific points in a
complex type.

### Registering Custom Type Checkers

You can register custom checkers for entire types using the
`check.WithTypeChecker` option. This is useful for types with complex
comparison logic, such as those requiring deep equality checks or custom
tolerances. The process is similar to `check.WithTrailChecker`, but applies to
all instances of a type rather than a specific trail.

Custom type checkers are also invaluable when working with types that have
non-exported fields. In Go, non-exported fields are inaccessible outside their
defining package, which prevents the `assert` package from directly comparing
them. By defining a custom type checker, you can implement comparison logic
that accesses these private fields within the same package, ensuring accurate
assertions for such types.

```go
type T struct{ value float64 }

chk := func(want, have any, opts ...any) error {
    w := want.(T)
    h := have.(T)
    return check.Epsilon(w.value, h.value, 0.001, opts...)
}

opt := check.WithTypeChecker(T{}, chk)

want := T{value: 1.2345}
have := T{value: 1.2346}
err := check.Equal(want, have, opt)

fmt.Println(err)
// Output:
//  <nil>
```

### Registering Global Type Checkers

Global checkers provide a convenient way to apply custom comparison logic
across all assertions for a specific type, without needing to specify the
checker in each `assert.Equal` call. This is particularly useful for complex
types with non-exported fields or custom comparison requirements. Use the
`check.RegisterTypeChecker` function to register a global checker.

Key Points for Global Checkers:

- **Registration**: call `check.RegisterTypeChecker` once, typically during
  package initialization or in a `TestMain` function, to ensure the checker is
  available for all tests.
- **Scope**: the checker applies to all assertions involving the registered
  type, streamlining test code.
- **Non-Exported Fields**: global checkers are ideal for types with
  non-exported fields, as they allow you to define comparison logic that
  accesses private data.
- **Thread Safety**: ensure the checker function is thread-safe, as it may be
  called concurrently in tests.

There are two suggested ways to register a global type checker. Either using
`TestMain` function or init function in one of your `_test.go` files.

To register a global type checker, you can use either the `TestMain` function
or an `init` function in a `_test.go` file. The `TestMain` approach is
preferred for centralized test setup, ensuring the checker is registered before
any tests run. An `init` function in a test file is suitable for
package-specific checkers, automatically executed when the test package is
loaded.

```go
func TestMain(m *testing.M) {
    check.RegisterTypeChecker(LocalType{}, checker)
    os.Exit(m.Run())
}

// or

func init() {
    check.RegisterTypeChecker(LocalType{}, checker)
}
```

Every time the new global type checker is registered, you will also see the
below line in the test log:

```log
*** CHECK /path/to/registration/call/all_test.go:20: Registering type checker for: mocker.goimp
```

Every time a check overrides the global type checker with
`checker.WithTypeChecker` option the log is written:

```log
*** CHECK /path/to/option/call/file_test.go:20: Overwriting the global type checker for: mocker.goimp
```

### Skipping Fields, Elements, or Indexes

You can ask for certain trials to be skipped when asserting.

```go
type T struct {
    Int  int
    Next *T
}

have := T{1, &T{2, &T{3, &T{42, nil}}}}
want := T{1, &T{2, &T{8, &T{42, nil}}}}
trails := make([]string, 0)

assert.Equal(
    want,
    have,
    check.WithTrailLog(&trails),
    check.WithSkipTrail("T.Next.Next.Int"),
)

fmt.Println(strings.Join(trails, "\n"))
// Test Log:
//
// T.Int
// T.Next.Int
// T.Next.Next.Int <skipped>
// T.Next.Next.Next.Int
// T.Next.Next.Next.Next
```

Notice that the requested trail was skipped from assertion even though the
values were not equal `3 != 8`. The skipped paths are always marked with
` <skipped>` tag.

### Skipping unexported fields

The `assert.Equal` will fail the test if the compared values (structs) have
unexported fields. This happens by design to make sure the equality check
doesn't silently ignore unexported fields. In cases like this the testing
module requires from a developer either explicitly specify fields to skip
during comparison or enable a mode that ignores all unexported fields, as
supported by the testing framework.

```go
type T struct {
    Int  int
    prv  int
    Next *T
}

have := T{1, -1, &T{2, -2, &T{3, -3, &T{42, -4, nil}}}}
want := T{1, -7, &T{2, -7, &T{3, -7, &T{42, -7, nil}}}}
trails := make([]string, 0)

err := check.Equal(
    want,
    have,
    check.WithTrailLog(&trails),
    check.WithSkipUnexported,
)

fmt.Println(err)
fmt.Println(strings.Join(trails, "\n"))
// Output:
// <nil>
// T.Int
// T.prv <skipped>
// T.Next.Int
// T.Next.prv <skipped>
// T.Next.Next.Int
// T.Next.Next.prv <skipped>
// T.Next.Next.Next.Int
// T.Next.Next.Next.prv <skipped>
// T.Next.Next.Next.Next
```