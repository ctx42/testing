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
checkers and dumpers can target any type or trail. The `check` package returns
plain `error` instead of calling `t.Fatal`, so checks compose naturally into
your own assertion functions. The `tester` package provides a `Spy` that
records calls to `t.Error`/`t.Fatal`, letting you write real tests for
your own test helpers.

The library supports both global customization (via `RegisterTypeChecker`,
`RegisterTypeDumper`, and package-level variables) and fine-grained per-call
control via options.

## Design

The toolkit follows a deliberate layered architecture:

- `assert` — the user-facing layer. Thin wrappers that call `t.Error` or
  `t.Fatal` on failure.
- `check` — pure functions that return `error` (or `*notice.Notice`).
  Composable, no test dependency, ideal for building your own helpers.
- `notice` — rich, structured message builder with trails, rows,
  metadata, and join support for detailed diagnostics.

This separation lets the same powerful checks be used both directly in
tests and inside custom assertion functions without pulling in testing
concerns.

The customization model supports both global defaults (via
`RegisterTypeChecker` / `RegisterTypeDumper` and package-level variables)
and fine-grained per-use control through options passed to individual
checks and assertions (see [check.DefaultOptions] and the `With*` functions
in `check` and `dump`).

# Installation

```shell
go get github.com/ctx42/testing
```

# Packages

All packages include rich package-level overviews, extensive
cross-references, and (where appropriate) executable `examples_test.go`
files whose output is wired into the READMEs.

Certain packages expose intentionally public surface for advanced or
cross-project use:
- `pkg/testcases` — battle-tested test values for writing and testing
  custom assertions/helpers (see its package documentation).
- `pkg/kit` top-level helpers such as `AddGlobalCleanup`/`RunGlobalCleanups`
  — for `TestMain`-style post-test coordination (see godoc for warnings).

## Main Packages

Packages used directly in test files.

- [assert](pkg/assert/README.md) — 72+ assertion functions with rich
  trails and custom checkers.
- [check](pkg/check/README.md) — composable checks that return `error`
  (or `*notice.Notice`), designed for both direct use and custom helpers.
- [goldy](pkg/goldy/README.md) — golden file testing with testable
  public surface via `tester.T`.
- [kit](pkg/kit/README.md) — curated collection of focused test helpers
  (I/O buffers, deterministic clocks, reflection utilities, plus a few
  top-level helpers).
- [mock](pkg/mock/README.md) — primitives for writing interface mocks
  (expectations, matchers, call recording).
- [mocker](pkg/mocker/README.md) — code generator for interface mocks
  that integrate with the `mock` package.
- [must](pkg/must/README.md) — helpers that panic on error for concise
  test setup and assertions.

## Supporting Packages

Foundational packages for building custom checks, assertions, and helpers.

- [dump](pkg/dump/README.md) — configurable renderer of any value to a
  human-readable string (used throughout for diagnostics).
- [notice](pkg/notice/README.md) — builder for rich, structured
  assertion messages with trails, rows, metadata, and join support.
- [tester](pkg/tester/README.md) — `Spy` implementation of `tester.T`
  for writing real tests of your own test helpers.

Each package has its own `README.md` with detailed usage and most
include an `examples_test.go` file with runnable, documented examples.

