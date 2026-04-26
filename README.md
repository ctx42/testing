![gopher.png](doc/gopher.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/testing)](https://goreportcard.com/report/github.com/ctx42/testing)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/testing)
![Tests](https://github.com/ctx42/testing/actions/workflows/go.yml/badge.svg?branch=master)

---

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Packages](#packages)
  * [Main Packages](#main-packages)
  * [Supporting Packages](#supporting-packages)
<!-- TOC -->

Zero-dependency testing toolkit for Go:

- `assert` — more than 72 hand-picked assertion functions.
- `mock` and `mocker` — test doubles and mock code generation.
- `goldy` — golden file testing.
- `kit` — assorted test helpers.
- `dump` — configurable any-type-to-string renderer.
- `check` — composable checks for building custom assertions.
- `tester` — tools for testing your own test helpers.

# Overview

Focused set of testing tools with zero external dependencies. Each package
targets a specific aspect of testing — use only what you need. The `check`
package returns plain errors, making it easy to compose custom assertions.
The `tester` package lets you test those helpers against a spy `T`
implementation. The API is chainable where it makes sense, and error
messages are descriptive.

# Installation

```shell
go get github.com/ctx42/testing
```

# Packages

## Main Packages

Packages used in test files.

- Package [assert](pkg/assert/README.md) — assertion toolkit (72 functions).
- Package [check](pkg/check/README.md) — equality checks used by `assert`; returns `error` for composability.
- Package [goldy](pkg/goldy/README.md) — golden file support.
- Package [kit](pkg/kit/README.md) — test helpers that are not assertions.
- Package [mock](pkg/mock/README.md) — primitives for writing interface mocks.
- Package [mocker](pkg/mocker/README.md) — interface mock code generator.
- Package [must](pkg/must/README.md) — helpers that panic on error.

## Supporting Packages

Packages for building custom checks, assertions, and helpers.

- Package [dump](pkg/dump/README.md) — configurable renderer of any type to a string.
- Package [notice](pkg/notice/README.md) — formatted assertion message builder.
- Package [tester](pkg/tester/README.md) — facilities for testing test helpers.

Each package has its own README, and most include an `examples_test.go` file
with usage examples.
