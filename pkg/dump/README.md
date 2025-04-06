<!-- TOC -->
* [Dump Package](#dump-package)
  * [Basic Usage](#basic-usage)
  * [Configuration Options](#configuration-options)
    * [Flat Output](#flat-output)
    * [Custom Time Formats](#custom-time-formats)
    * [Pointer Addresses](#pointer-addresses)
    * [Custom Dumpers](#custom-dumpers)
* [Handling Complex and Recursive Types](#handling-complex-and-recursive-types)
* [Extensibility](#extensibility)
* [Conclusion](#conclusion)
<!-- TOC -->

# Dump Package

The `dump` package is a utility that serializes any Go value into its string
representation.

The `dump` package is part of the CTX42 Testing Module, an ongoing effort to 
build a new, flexible, and developer-friendly testing framework.

The `dump` package provides a configurable way to render any Go value — whether 
it’s a simple integer, a nested struct, or a recursive data structure — into a 
human-readable string. This is particularly useful in testing, where comparing 
complex values often requires more than Go’s built-in `reflect.DeepEqual`. By 
converting values to strings, the `dump` package lets you leverage string 
comparison or diffing tools to pinpoint discrepancies quickly and accurately.

## Basic Usage

```go
val := types.TA{
    Dur: 3,
    Int: 42,
    Loc: types.WAW,
    Str: "abc",
    Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
    TAp: nil,
}

have := dump.Default().Any(val)

fmt.Println(have)
// Output:
//	{
//		Int: 42,
//		Str: "abc",
//		Tim: "2000-01-02T03:04:05Z",
//		Dur: "3ns",
//		Loc: "Europe/Warsaw",
//		TAp: nil,
//	}
```

The default dump renders the struct in a nicely formatted, multi-line string.
Fields are listed in the order they’re declared in the struct, ensuring
consistent output for reliable comparisons.

## Configuration Options

One of the `dump` package’s strengths is its configurability. You can tweak how
values are rendered to suit your needs. Here are some key options:

### Flat Output

For a compact, single-line representation, use the Flat option:

```go
val := map[string]any{
    "int": 42,
    "loc": types.WAW,
    "nil": nil,
}

have := dump.New(dump.WithFlat).Any(val)

fmt.Println(have)
// Output:
// map[string]any{"int": 42, "loc": "Europe/Warsaw", "nil": nil}
```

For maps, keys are sorted (when possible) to maintain consistency.

### Custom Time Formats

You can customize how `time.Time` values are displayed using the 
`dump.WithTimeFormat` option:

```go
val := map[time.Time]int{time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC): 42}

have := dump.New(dump.WithFlat, dump.WithTimeFormat(time.Kitchen)).Any(val)

fmt.Println(have)
// Output:
// map[time.Time]int{"3:04AM": 42}
```

### Pointer Addresses

By default, pointer addresses are hidden, but you can enable them with 
`dump.WithPtrAddr` option:

```go
val := map[string]any{
    "fn0": func() {},
    "fn1": func() {},
}

have := New(dump.WithPtrAddr).Any(val)

fmt.Println(have)
// Output:
// map[string]any{
// 	"fn0": <func>(<0x533760>),
// 	"fn1": <func>(<0x533780>),
// }
```

### Custom Dumpers

For ultimate flexibility, you can define custom dumpers for specific types.
Dumpers for types are regular functions matching `dump.Dumper` signature declared in 
the package.

```go
type Dumper func(dmp Dump, level int, val reflect.Value) string
```

For example:

```go
var i int
customIntDumper := func(dmp Dump, lvl int, val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int:
		return fmt.Sprintf("%X", val.Int())
	default:
		panic("unexpected kind")
	}
}

have := dump.New(dump.WithFlat, dump.WithCompact, dump.WithDumper(i, customIntDumper)).Any(42)

fmt.Println(have)
// Output:
// 2A
```

The above example dumps integers as hexadecimal values, showcasing how you can
tailor the output for your use case.

# Handling Complex and Recursive Types

The `dump` package shines when dealing with complicated or recursive data
structures. It includes cycle detection to prevent infinite loops. Here’s an
example with a recursive struct:

```go
type Node struct {
    Value    int
    Children []*Node
}

val := &Node{
    Value: 1,
    Children: []*Node{
        {Value: 2},
        {Value: 3, Children: []*Node{{Value: 4}}},
    },
}

have := dump.New().Any(val)
fmt.Println(have)
// Output:
// {
//		Value: 1,
//		Children: []*dump_test.Node{{
//			Value: 2,
//			Children: nil,
//		}, {
//			Value: 3,
//			Children: {{
//				Value: 4,
//				Children: nil,
//			}},
//		}},
//	}
```

```go
type Node struct {
    Value    int
    Children []*Node
}

val := &Node{
    Value: 1,
    Children: []*Node{
        {Value: 2, Children: nil},
        {
            Value: 3,
            Children: []*Node{
                {Value: 4, Children: nil},
            },
        },
    },
}

have := dump.New().Any(val)
fmt.Println(have)
// Output:
// {
// 	Value: 1,
// 	Children: []*dump_test.Node{
// 		{
// 			Value: 2,
// 			Children: nil,
// 		},
// 		{
// 			Value: 3,
// 			Children: []*dump_test.Node{
// 				{
// 					Value: 4,
// 					Children: nil,
// 				},
// 			},
// 		},
// 	},
// }
```

# Extensibility

The Dump package is built with extensibility in mind. Custom dumpers let you
define how your own types are serialized, integrating seamlessly with the
framework. This adaptability ensures the package can grow with your project’s
needs.

# Conclusion

Testing equivalence of complex data structures in Go doesn’t have to be a chore.
The `dump` package simplifies the process by converting any value into a string,
ready for comparison or diffing. Its configurability and extensibility make it 
a versatile tool for any testing scenario.
