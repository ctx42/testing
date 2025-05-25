<!-- TOC -->
* [Notice Package](#notice-package)
  * [Basic Usage](#basic-usage)
    * [Create a Message](#create-a-message)
    * [Wrap Errors](#wrap-errors)
    * [Add Metadata](#add-metadata)
  * [Indenting Lines](#indenting-lines)
<!-- TOC -->

# Notice Package

The `notice` package offers a suite of utilities for crafting clear, structured
assertion messages. It simplifies the creation of readable error messages,
featuring a header and contextual rows. With a fluent interface, it enables
seamless message construction and includes helper functions for formatting and
unwrapping errors.

## Basic Usage

### Create a Message

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

```go
msg := notice.New("expected values to be equal").
    Want("%s", "abc").
    Have("%s", "xyz").
    MetaSet("key", 123)

fmt.Println(msg.MetaLookup("key"))
// Output: value true
}
```

For more examples see the [examples_test.go](examples_test.go) file.

## Indenting Lines

```go
lines := notice.Indent(4, ' ', "line1\nline2\nline3")

fmt.Println(lines)
// Output:
//     line1
//     line2
//     line3
```
