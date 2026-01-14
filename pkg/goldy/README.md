<!-- TOC -->
  * [Features](#features)
  * [Usage](#usage)
    * [Opening Golden Files](#opening-golden-files)
    * [Golden File Format](#golden-file-format)
    * [Creating Golden File](#creating-golden-file)
    * [Example](#example)
      * [Directory Structure](#directory-structure)
      * [Test](#test)
    * [Golden file template](#golden-file-template)
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

The goldy package provides two main functions:
- `goldy.Open` – reads the content of an existing golden file,
- `goldy.Create` – creates a new golden file.

### Opening Golden Files

```go
func Open(t core.T, pth string, opts ...func(*Goldy)) *Goldy
```

- `t core.T`: A test context (typically a `*testing.T`).
- `pth string`: The golden file path relative to the test file or absolute.

### Golden File Format

A golden file must include a `---` marker line. Content before this marker is
treated as documentation and ignored, while everything after it is returned as
the test data.

```text
This is multiline documentation about the golden file’s contents. It explains
what the file contains and why it’s used. Documenting golden files helps keep
tests maintainable.
---
Content #1.
Content #2.
```

It's customary for golden files to have `.gld` extension. 

### Creating Golden File

```go
gld := goldy.Create(t, "testdata/example.gld")
gld.SetComment("Multi\nline\ncontent")
gld.SetContent("Content #1.\nContent #2.")
gld.Save()

// File contents:
// Multi
// line
// content
// ---
// Content #1.
// Content #2.
```

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
        want := goldy.Open(t, "testdata/case1.gld").String()
        if want != have {
            format := "expected values to be equal:\n  want: %q\n  have: %q"
            t.Errorf(format, want, have)
        }
    })
}
```

### Golden file template

For cases where a golden file needs dynamic content, you can use Go text
templates.

```text
Golden file template.
---
Content #{{ .first }}.
```

then use it:

```go
data := WithData(map[string]any{"first": 1})
gld := Open(tspy, "testdata/test_tpl.gld", data)
```

### Error Handling

If `goldy` encounters issues (e.g., file not found, missing `---` marker), it
reports errors via the test context using `t.Fatalf`, marking the test as
failed. This ensures clear feedback for debugging.

## Updating Golden Files

Example:

```go
gld := gold.Open(t, "test.gld")
gld.SetComment("Mock for TestInterface")
gld.SetContent("type TestInterface struct {...}")
gld.Save()
```