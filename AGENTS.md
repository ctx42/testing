# Project Rules for testing

This file provides guidance to Claude Code (claude.ai/code) when
working with code in this repository.

## Project Overview

Go testing module (`github.com/ctx42/testing`) — a zero-dependency
testing toolkit providing assertions, mocks, golden file testing, and
test helpers. Requires Go 1.26 (as declared in go.mod).

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
- `pkg/testcases` — intentionally public collection of test values and
  cases for people writing custom assertions or testing their own test
  helpers with [tester.Spy]. See its package documentation for the
  intended usage patterns.
- `pkg/kit` (specifically `AddGlobalCleanup` / `RunGlobalCleanups`) —
  intentionally public global cleanup coordination helpers. These are used
  extensively by external libraries for `TestMain`-style post-test cleanup.
  The feature uses package-level state by design. See the godoc for usage.
- `internal/diff` — text diffing (Myers algorithm via
  `internal/diff/lcs`)

## Source Code Conventions

- Ignore code in `internal/diff` package.

### Formatting and Readability

When writing or editing Go code, prefer readability over compactness:

- Extract long or complex format strings (and similar literals) into a
  local variable on its own line before passing them to `fmt.Sprintf`,
  `notice.New`, or equivalent functions.
- Preferred pattern (used extensively in `pkg/mock` and `pkg/check`):

  ```go
  mHeader := "[mock] arguments: Foo(%d) is of type \"%T\" not foo"
  msg := notice.New(mHeader, idx, val)
  panic(msg)
  ```

  or

  ```go
  format := "%d: FAIL:\n    want: %s\n    have: %s"
  return fmt.Sprintf(format, i, left, right)
  ```

- This keeps lines short, makes intent clearer at a glance, and aligns
  with the project's 80-column guidance.
- Use short, idiomatic local variable names (`mHeader`, `msg`, `ops`,
  `chk`, `format`, etc.) for these intermediate values.
- Mechanical formatting is handled by `gofmt` / `gofumpt`, but higher-level
  readability decisions (such as the above) are made manually.

- Only break function calls, method chains, and long expressions across
  multiple lines when the compact single-line version would exceed 80
  columns. When breaking a call for length, use this preferred layout
  (one logical argument per line, opening `(` on the call line):
  ```go
  ret = spy.checkCallMaybeCnt(
  	"TempDir",
  	spy.wantTempDirCnt,
  	len(spy.haveTempDirs),
  )
  ```

- **Function body leading blank lines** (narrow exception only):
  - When a function or method declaration has its parameters on
    separate lines (the `{` appears on its own line as the signature
    closer), place exactly one blank line after the `{` before the
    first statement of the body. This provides visual separation only
    when the declaration itself already spans multiple lines.
  - All other declarations (single-line signatures) must start the
    body immediately on the next line with no leading blank.

    Correct multi-line signature:

    ```go
    func ChannelWillClose[C any](
    	t tester.T,
    	within any,
    	c <-chan C,
    	opts ...any,
    ) bool {

    	t.Helper()
    	...
    }
    ```

    Correct single-line signature:

    ```go
    func True(t tester.T, have bool, opts ...any) bool {
    	t.Helper()
    	if e := check.True(have, opts...); e != nil {
    		t.Error(e)
    		return false
    	}
    	return true
    }
    ```

### Writing Good Assertions

When adding or modifying assertion or check functions:

- Keep per-function godoc short and focused on what the assertion *does*.
- Document standard failure behavior, option handling, and common patterns
  once at the package level (see the package documentation of [assert] and
  [check]).
- Use cross-references liberally (`[Type]`, `[check.WithTrail]`,
  `[notice.Notice]`, etc.).
- Avoid repeating the same "marks the test as failed..." boilerplate in every
  function.

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
  `*testing.T` when testing helpers.

### Test Function Naming

Each exported function or method should have a corresponding test with a
matching name:

- Package-level functions: `Test_<FunctionName>`
- Methods on a type: `Test_<Type>_<Method>`

Examples:
- `func Equal(...)` → `func Test_Equal(t *testing.T)`
- `func (s *Spy) ExpectCleanups(...)` → `func Test_Spy_ExpectCleanups(t *testing.T)`

When a test uses table-driven (tabular) style, append `_tabular`:

- `func Test_Equal_tabular(t *testing.T)`
- `func Test_Contain_success_tabular(t *testing.T)`
- `func Test_Contain_error_tabular(t *testing.T)`

Additional descriptive suffixes (e.g. `_kind_Ptr`, `_smoke`) are
discouraged for the primary test of a function or method. Use subtests or
separate focused test functions instead when more granularity is needed.

When the same set of test cases is used to exercise multiple methods on a
type (or multiple related functions), extract the test cases into a shared
helper and create a dedicated tabular test for each method. For example,
instead of a combined `Test_Dump_Any_Value_tabular`, prefer
`Test_Dump_Any_tabular` and `Test_Dump_Value_tabular`.

### Subtest Names

When using `t.Run(name, ...)`, `b.Run(...)`, or `f.Run(...)`, avoid characters
that are illegal or problematic in filenames on common platforms
(`()*#$?<>|:"\/` and similar). Such characters frequently cause issues with
coverage profiles, test artifacts, IDE runners, and CI systems that derive
filenames from subtest names.

See the rule `subtest-names.md` in the Go quality skills for the full
rationale and recommended character set.

### Tests That Mutate Package Globals

Some tests in this module need to directly read or write package-level global
variables (for example the global cleanup registry in `pkg/kit` or the global
type registries in `pkg/dump` and `pkg/check`).

All such tests **must** start with:

```go
t.Setenv("___", "___")
```

This forces Go's test framework to serialize the test with any other test that
also calls `Setenv` on the same key. Because we use the same dummy key in every
test that touches shared global state, these tests are effectively prevented
from running in parallel with each other.

When a test needs to reset global state, prefer a small helper such as
`resetCleanupsForTest()` (see `pkg/kit/cleanup_test.go`) rather than scattering
direct assignments to globals throughout the test. The helper should acquire
any necessary locks and the calling test must still start with the
`t.Setenv("___", "___")` marker.

## License

MIT — files carry SPDX headers.
