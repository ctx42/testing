<!-- TOC -->
* [The `assert` package](#the-assert-package)
  * [Assertions](#assertions)
    * [Asserting Structures](#asserting-structures)
    * [Asserting Recursive Structures](#asserting-recursive-structures)
    * [Asserting Maps, Arrays and Slices](#asserting-maps-arrays-and-slices)
      * [Asserting Time](#asserting-time)
      * [Asserting JSON Strings](#asserting-json-strings)
      * [Worthy mentions](#worthy-mentions)
  * [Advanced usage](#advanced-usage)
    * [Custom Checkers](#custom-checkers)
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

Most of the assertions are self-explanatory and I encourage you to see your
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

### Asserting Maps, Arrays and Slices

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

See the [documentation](https://pkg.go.dev/github.com/ctx42/testing) for the 
full list.

## Advanced usage

### Custom Checkers

You can define custom checkers for any trail define a function matching 
`check.Check`. See example below. 

```go
type T struct {
    Str string
    Any []any
}

chk := func(want, have any, opts ...check.Option) error {
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

The trail uniquely identifies the struct field, slice or array element, or map
key the assertion visits. The assert package tracks the tails for all composite 
types. To see all visited trails do this:

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

### Skipping Fields, Elements, or Indexes

You can ask certain trials to be skipped when asserting.

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