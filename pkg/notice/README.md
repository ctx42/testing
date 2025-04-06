<!-- TOC -->
* [Notice Package](#notice-package)
  * [Basic Usage](#basic-usage)
    * [Create a message](#create-a-message)
  * [Formatting Lines](#formatting-lines)
<!-- TOC -->

# Notice Package

The `notice` package provides a set of utilities for building structured 
assertion messages. It's designed to create easy to read and understand error 
messages with a header and contextual rows. The package supports fluent 
interfaces for building messages and includes helper functions for formatting 
and unwrapping errors.

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

### Add Metadata

```go
msg := notice.New("expected values to be equal").
    Want("%s", "abc").
    Have("%s", "xyz").
    SetData("key", 123)

fmt.Println(msg.GetData("key"))
// Output: value true
}
```

For more examples see [examples_test.go](examples_test.go) file.

## Indenting Lines

```go
lines := notice.Indent(4, ' ', "line1\nline2\nline3")

fmt.Println(lines)
// Output:
//     line1
//     line2
//     line3
```
