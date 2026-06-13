<!-- TOC -->
* [Notice Package](#notice-package)
  * [Basic Usage](#basic-usage)
    * [Create a Message](#create-a-message)
    * [Wrap Errors](#wrap-errors)
    * [Add Metadata](#add-metadata)
  * [Indenting Lines](#indenting-lines)
  * [Chaining and Error Inspection](#chaining-and-error-inspection)
<!-- TOC -->

# Notice Package

The `notice` package offers a suite of utilities for crafting clear, structured
assertion messages. It simplifies the creation of readable error messages,
featuring a Header and contextual rows. With a fluent interface, it enables
seamless message construction and includes helper functions for formatting and
unwrapping errors.

## Basic Usage

### Create a Message

<!-- gmdoceg:ExampleNew -->
```go
msg := notice.New("expected values to be equal").
	Want("%s", "abc").
	Have("%s", "xyz")

fmt.Println(msg)
// Output:
// expected values to be equal:
//   want: abc
//   have: xyz
```

### Wrap Errors

<!-- gmdoceg:ExampleNotice_Wrap -->
```go
ErrMy := errors.New("my error")

msg := notice.New("expected values to be equal").
	Want("%s", "abc").
	Have("%s", "xyz").
	Wrap(ErrMy)

is := errors.Is(msg, ErrMy)
fmt.Println(is)
// Output: true
```

### Add Metadata

<!-- gmdoceg:ExampleNotice_MetaSet -->
```go
msg := notice.New("expected values to be equal").
	Want("%s", "abc").
	Have("%s", "xyz").
	MetaSet("key", 123)

fmt.Println(msg.MetaLookup("key"))
// Output: 123 true
```

For more examples see the [examples_test.go](examples_test.go) file.

## Indenting Lines

<!-- gmdoceg:ExampleIndent -->
```go
lines := notice.Indent(4, ' ', "line1\nline2\nline3")

fmt.Println(lines)
// Output:
//     line1
//     line2
//     line3
```

## Chaining and Error Inspection

Notices implement [error] and can be chained into a linked list. This is
useful when multiple independent expectations fail in one assertion.

```go
err := notice.Join(
    notice.New("first failure"),
    notice.New("second failure"),
)
```

The result is a single error whose [notice.Notice.All] returns the full
list. Walk the chain with:

- [notice.Notice.Head] — first element
- [notice.Notice.Next] / [notice.Notice.Prev] — traversal
- [notice.Join] — the builder used above

Each notice delegates [errors.Is] and [errors.As] to its base error
(default [notice.ErrNotice]). Change the base with [notice.Notice.Wrap]:

```go
myErr := errors.New("root cause")
n := notice.New("something failed").Wrap(myErr)
fmt.Println(errors.Is(n, myErr)) // true
```

Chains are mutable (see [notice.Notice.Chain]). Prefer [Join] when
building from multiple values.
