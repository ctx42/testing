![gopher.png](doc/gopher.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/testing)](https://goreportcard.com/report/github.com/ctx42/testing)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/testing)
![Tests](https://github.com/ctx42/testing/actions/workflows/go.yml/badge.svg?branch=master)

---

**Zero-dependency testing toolkit for Go.** 72+ assertion functions,
golden files, mocks, and a spy for testing your own test helpers.

When an assertion fails, you see exactly where and why:

```go
type Order struct {
    ID    int
    Email string
    Total float64
}

want := Order{ID: 1, Email: "alice@example.com", Total: 9.99}
have := Order{ID: 1, Email: "alice@wrong.com",   Total: 9.99}

assert.Equal(t, want, have)
// Test log:
//
// expected values to be equal:
//   trail: Order.Email
//    want: "alice@example.com"
//    have: "alice@wrong.com"
```

Trails work through nested structs, maps, slices, and pointers. Custom
checkers can target any type or trail. The `check` package returns plain
`error` instead of calling `t.Fatal`, so checks compose naturally into
your own assertion functions. The `tester` package provides a `Spy` that
records calls to `t.Error`/`t.Fatal`, letting you write real tests for
your own test helpers.

# Installation

```shell
go get github.com/ctx42/testing
```

# Packages

## Main Packages

Packages used in test files.

- Package [assert](pkg/assert/README.md) — 72+ assertion functions.
- Package [check](pkg/check/README.md) — composable equality checks; returns `error`, not panics.
- Package [goldy](pkg/goldy/README.md) — golden file testing.
- Package [kit](pkg/kit/README.md) — test helpers (buffers, clocks, cleanup).
- Package [mock](pkg/mock/README.md) — primitives for writing interface mocks.
- Package [mocker](pkg/mocker/README.md) — interface mock code generator.
- Package [must](pkg/must/README.md) — helpers that panic on error.

## Supporting Packages

Packages for building custom checks, assertions, and helpers.

- Package [dump](pkg/dump/README.md) — configurable renderer of any type to a string.
- Package [notice](pkg/notice/README.md) — formatted assertion message builder.
- Package [tester](pkg/tester/README.md) — `Spy` for testing test helpers.

Each package has its own `README.md` and most include an `examples_test.go`
file with runnable examples.

