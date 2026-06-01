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

See the Design section in the root README for the overall layered architecture
(assert built on check built on notice) and the customization model.

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

<!-- gmdoceg:ExampleEqual_structs -->
```go
type T struct {
	Int int
	Str string
}

have := T{Int: 1, Str: "abc"}
want := T{Int: 2, Str: "xyz"}

// assert.Equal logs the error via t.Error; the message is identical to
// what check.Equal returns.
err := check.Equal(want, have)
fmt.Println(err)

// Output:
// multiple expectations violated:
//   error: expected values to be equal
//   trail: T.Int
//    want: 2
//    have: 1
//       ---
//   error: expected values to be equal
//   trail: T.Str
//    want: "xyz"
//    have: "abc"
```

### Asserting Recursive Structures

<!-- gmdoceg:ExampleEqual_recursiveStructs -->
```go
type T struct {
	Int  int
	Next *T
}

have := T{1, &T{2, &T{3, &T{42, nil}}}}
want := T{1, &T{2, &T{3, &T{4, nil}}}}

err := check.Equal(want, have)

fmt.Println(err)
// Output:
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

<!-- gmdoceg:ExampleEqual_slices -->
```go
want := []int{1, 2, 3}
have := []int{1, 2, 3, 4}

err := check.Equal(want, have)

fmt.Println(err)
// Output:
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
//       diff:
//             @@ -2,5 +2,4 @@
//                1,
//                2,
//             -  3,
//             -  4,
//             +  3,
//              }
```

#### Asserting Time

<!-- gmdoceg:ExampleTime -->
```go
want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
have := time.Date(2025, 1, 1, 0, 1, 1, 0, time.UTC)

err := check.Time(want, have)

fmt.Println(err)
// Output:
//  expected equal dates:
//   want: 2025-01-01T00:00:00Z
//   have: 2025-01-01T00:01:01Z
//   diff: -1m1s
```

#### Asserting JSON Strings

<!-- gmdoceg:ExampleJSON -->
```go
want := `{"A": 1, "B": 2}`
have := `{"A": 1, "B": 3}`

err := check.JSON(want, have)

fmt.Println(err)
// Output:
// expected JSON strings to be equal:
//   want: {"A":1,"B":2}
//   have: {"A":1,"B":3}
```

#### Worthy mentions

- `Epsilon` - assert floating point numbers within given ε.
- `ChannelWillClose` - assert channel will be closed within given time.
- `MapSubset` - checks the "want" is a subset "have".
- `Wait` - Wait waits for "fn" to return true but no longer than
  the given timeout.

See the [documentation](https://pkg.go.dev/github.com/ctx42/testing) for the
full list.

## Advanced usage

### Custom Checkers

Custom checkers allow you to define specialized comparison logic for any type
or field trail in your Go tests. A custom checker is a function that matches
the `check.Check` signature, enabling fine-grained control over assertions.
Below is an example demonstrating how to create and use a custom checker.

<!-- gmdoceg:ExampleEqual_customTrailChecker -->
```go
type T struct {
	Str string
	Any []any
}

chk := func(want, have any, opts ...any) error {
	wVal := want.(float64)
	hVal := have.(float64)
	return check.Epsilon(wVal, 0.01, hVal, opts...)
}
opt := check.WithTrailChecker("T.Any[1]", chk)

want := T{Str: "abc", Any: []any{1, 2.123, "abc"}}
have := T{Str: "abc", Any: []any{1, 2.124, "abc"}}

err := check.Equal(want, have, opt)

fmt.Println(err)
// Output:
//  <nil>
```

In this example, the custom checker `chk` compares float64 values at the trail
`T.Any[1]` with a tolerance of 0.01. The assertion passes because 2.123 and
2.124 are within the specified epsilon.

Also see the complete custom helper example in the [examples] package:
[custom_assertions_test.go](../../examples/custom_assertions_test.go).

### Understanding Trails

A trail uniquely identifies a struct field, slice or array element, or map key
visited during an assertion. The `assert` package automatically tracks trails
for composite types, enabling precise targeting for custom checkers. To inspect
all visited trails, use the check.WithTrailLog option, as shown below:

<!-- gmdoceg:ExampleEqual_listVisitedTrails -->
```go
type T struct {
	Int  int
	Next *T
}

have := T{1, &T{2, &T{3, &T{42, nil}}}}
want := T{1, &T{2, &T{3, &T{42, nil}}}}
trails := make([]string, 0)

err := check.Equal(want, have, check.WithTrailLog(&trails))

fmt.Println(err)
fmt.Println(strings.Join(trails, "\n"))
// Output:
// <nil>
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

<!-- gmdoceg:ExampleEqual_customTypeChecker -->
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
types with non-exported fields or custom comparison requirements.

Use [`check.RegisterTypeChecker`] to register a global checker for a type.
Once registered, the checker is used automatically by `assert.Equal` (and
`check.Equal`) for all values of that type.

#### Key Points

- **Registration**: Call `RegisterTypeChecker` once, typically from `TestMain`
  or an `init` function in a `_test.go` file.
- **Panic on duplicate**: `RegisterTypeChecker` panics if a checker for the
  same type is already registered. This prevents accidental conflicting
  registrations.
- **Overwriting via options**: Using the `check.WithTypeChecker` option for a
  type that already has a global checker will log a warning (it does **not**
  replace the global checker).
- **Non-exported fields**: Global checkers are the recommended way to provide
  comparison logic for types with unexported fields.
- **Thread safety**: The registered checker may be called from multiple
  goroutines during parallel tests, so it must be safe for concurrent use.

#### Recommended Registration Pattern

The `TestMain` approach is preferred for centralized setup. It guarantees the
checker is registered before any tests run in the package.

```go
func TestMain(m *testing.M) {
	check.RegisterTypeChecker(MyType{}, myChecker)
	os.Exit(m.Run())
}
```

An `init` function works but is less explicit about test setup ordering:

```go
func init() {
	check.RegisterTypeChecker(MyType{}, myChecker)
}
```

#### Log Messages

When a global checker is successfully registered you will see a line like this
in the test output (the exact path and line will vary):

```log
*** CHECK /path/to/your/package/all_test.go:42: Registering type checker for: pkg.MyType
```

When `check.WithTypeChecker` is used for a type that already has a global
checker, the following warning is logged (the global checker is **not**
replaced):

```log
*** CHECK /path/to/your/test/file_test.go:17: Overwriting the global type checker for: pkg.MyType
```

These log messages are emitted via an internal package logger and are useful
for debugging unexpected checker behavior during test runs.

> **Note**: This section documents setup patterns and side-effect logging
> rather than pure assertion usage. The examples are therefore shown as
> illustrative code rather than executable `Example*` functions.
```

### Skipping Fields, Elements, or Indexes

You can ask for certain trails to be skipped when asserting.

<!-- gmdoceg:ExampleEqual_skipTrails -->
```go
type T struct {
	Int  int
	Next *T
}

have := T{1, &T{2, &T{3, &T{42, nil}}}}
want := T{1, &T{2, &T{8, &T{42, nil}}}}
trails := make([]string, 0)

err := check.Equal(
	want,
	have,
	check.WithTrailLog(&trails),
	check.WithSkipTrail("T.Next.Next.Int"),
)

fmt.Println(err)
fmt.Println(strings.Join(trails, "\n"))
// Output:
// <nil>
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
module requires from a developer either to explicitly specify fields to skip
during comparison or enable a mode that ignores all unexported fields, as
supported by the testing framework.

<!-- gmdoceg:ExampleEqual_skipAllUnexportedFields -->
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