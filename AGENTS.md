# Project Rules for testing

This file provides guidance to Claude Code (claude.ai/code) when
working with code in this repository.

## Project Overview

Go testing module (`github.com/ctx42/testing`) — a zero-dependency
testing toolkit providing assertions, mocks, golden file testing, and
test helpers. Requires Go 1.24+.

## Commands

```bash
# Run all tests
go test ./...

# Run all tests with race detection (this is what CI runs)
go test -v -race ./...

# Run a single test
go test -run TestName ./pkg/assert/

# Run tests for a specific package
go test ./pkg/check/
```

## Architecture

The module follows a layered design where each layer builds on the
one below:

**`assert`** → **`check`** → **`notice`** (assertion message formatting)

- `check` functions return `error` on failure (composable, no test
  dependency)
- `assert` functions wrap `check`, calling `t.Error()`/`t.Fatal()` on
  failure
- Both accept variadic `opts ...any` for configuration (see
  `check.DefaultOptions`)

**Key interfaces:**

- `tester.T` (`pkg/tester/t.go`) — subset of `testing.TB` used
  throughout instead of `*testing.T` directly, enabling mock/spy usage
  in tests
- `core.T` (`internal/core/spy.go`) — minimal internal test interface

**Testing test helpers:**

- `pkg/tester` provides `Spy` for testing assertion functions
  themselves — verifying that helpers correctly call `t.Error`/`t.Fatal`
- `internal/core` provides `core.NewSpy()` used in internal tests
- `internal/affirm` provides simple affirmation functions used to test
  the testing toolkit itself (avoids circular dependencies)

**Other packages:**

- `pkg/mock` — primitives for writing interface mocks (matchers,
  candidates)
- `pkg/mocker` — interface mock code generator
- `pkg/goldy` — golden file test support
- `pkg/dump` — configurable any-type-to-string renderer
- `pkg/kit` — assorted test helpers (sub-packages: `iokit`, `timekit`)
- `pkg/must` — helpers that panic on error
- `internal/diff` — text diffing (Myers algorithm via
  `internal/diff/lcs`)

## Source Code Conventions

- Ignore code in `internal/diff` package.

### Editor Configuration

This project uses a minimal `.editorconfig` for basic cross-editor
consistency:

- UTF-8, LF line endings, final newline required, no trailing whitespace.
- Go: tabs (actual formatting controlled by `gofmt`/`gofumpt` + linters).
- YAML: 2-space indent.
- Markdown (including `AGENTS.md`): 80-column guidance.

See the root `.editorconfig` for details.

## Test Conventions

- Test spies (`core.NewSpy()`, `tester.Spy`) are used instead of
  `*testing.T` when testing helpers

## License

MIT — files carry SPDX headers.
