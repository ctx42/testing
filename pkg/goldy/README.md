<!-- TOC -->
  * [Features](#features)
  * [Usage](#usage)
    * [Opening Golden Files](#opening-golden-files)
    * [Golden File Format](#golden-file-format)
    * [Example](#example)
      * [Directory Structure](#directory-structure)
      * [Test](#test)
    * [Error Handling](#error-handling)
  * [Updating Golden Files](#updating-golden-files)
<!-- TOC -->

The `goldy` is a Go package designed to simplify reading content from golden
files in tests.

## Features

- Documentation is part of the golden file.
- Simple file format.
- No dependencies.

## Usage

The `goldy` package provides a single function, `goldy.New`, which reads the
content of a golden file.

### Opening Golden Files

```go
func New(t core.T, pth string) string
```

- `t core.T`: A test context (typically a `*testing.T`).
- `pth string`: The golden file path relative to the test file or absolute.

The internal `core.T` interface enables unit testing by allowing mocks,
ensuring `goldy` integrates seamlessly with test helpers.

### Golden File Format

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

It's customary for golden files to have `.gld` extension. 

### Example

Below is a practical example showing how to use `goldy` to read a golden file
and compare it with generated output in a test.

#### Directory Structure

```
project/
├── testdata/
│   └── case1.gld
├── my_test.go
```

#### Test

Content of `my_test.go` file.

```go
package project

import (
    "testing"

    "github.com/ctx42/testing/pkg/goldy"
)

func Test_Generator(t *testing.T) {
    t.Run("Generate content", func(t *testing.T) {
        // When
        have := Generate() // Returns string.

        // Then
        want := goldy.New(t, "testdata/case1.gld").String()
        if want != have {
            format := "expected values to be equal:\n  want: %q\n  have: %q"
            t.Errorf(format, want, have)
        }
    })
}
```

### Error Handling

If `goldy` encounters issues (e.g., file not found, missing `---` marker), it
reports errors via the test context using `t.Fatalf`, marking the test as
failed. This ensures clear feedback for debugging.

## Updating Golden Files

The `Goldy` struct, returned by the `New` function, provides a `Save` method to 
update a golden file. Calling `Goldy.Save` writes the modified `Comment` and 
`Content` fields to the original file path. The `Comment` field typically 
contains metadata or a description, while `Content` holds the expected test 
output, separated by the `Marker`.

Example:

```go
gld := gold.New(t, "test.gld")
gld.Comment = "Mock for TestInterface"
gld.Content = "type TestInterface struct {...}"
gld.Save()
