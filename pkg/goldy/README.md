<!-- TOC -->
  * [Features](#features)
  * [Usage](#usage)
    * [Function Signature](#function-signature)
    * [Text Golden File Format](#text-golden-file-format)
    * [Example](#example)
      * [Directory Structure](#directory-structure)
      * [Test](#test)
    * [Error Handling](#error-handling)
<!-- TOC -->

The `goldy` is a Go package designed to simplify reading content from golden
files in tests.

## Features

- Documentation is part of the golden file.
- Simple file format.
- No dependencies.

## Usage

The `goldy` package provides a single function, `goldy.Text`, which reads the
content of a golden file starting after a mandatory `---` marker, skipping any
preceding documentation. It takes a test context and the file path as arguments,
returning the content as a string.

### Function Signature

```go
func Text(t core.T, pth string) string
```

- `t core.T`: A test context (typically a `*testing.T`).
- `pth string`: The golden file path relative to the test file or absolute.

The internal `core.T` interface enables unit testing by allowing mocks,
ensuring goldy integrates seamlessly with test helpers.

### Text Golden File Format

A golden file must include a `---` marker line. Content before this marker is
treated as documentation and ignored, while everything after it is returned as
the test data.

```text
This is multi-line documentation about the golden file’s contents. It explains
what the file contains and why it’s used. Documenting golden files helps keep
tests maintainable.
---
Content #1.
Content #2.
```

### Example

Below is a practical example showing how to use `goldy` to read a golden file
and compare it with generated output in a test.

#### Directory Structure

```project/
├── testdata/
│   └── case1.txt
├── my_test.go
```

#### Test

```go
package project

import (
    "testing"

    "github.com/ctx42/testing/pkg/goldy"
)

func Test_Generator(t *testing.T) {
    t.Run("generate content", func(t *testing.T) {
        // When
        have := Generate()

        // Then
        want := goldy.Text(t, "testdata/case1.txt")
        if want != have {
            format := "expected values to be equal:\n  want: %q\n  have: %q"
            t.Errorf(format, want, have)
        }
    })
}
```

### Error Handling

If `goldy` encounters issues (e.g., file not found, missing `---` marker), it
reports errors via the test context using `t.Errorf`, marking the test as
failed without panicking. This ensures clear feedback for debugging.
