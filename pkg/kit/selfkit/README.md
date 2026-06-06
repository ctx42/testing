<!-- TOC -->
* [The `selfkit` package](#the-selfkit-package)
  * [Setup](#setup)
  * [Writing a subprocess test](#writing-a-subprocess-test)
  * [Flags](#flags)
    * [--toStdout and --toStderr](#--tostdout-and---tostderr)
    * [--exitCode](#--exitcode)
    * [--printEnv](#--printenv)
    * [--printArgs](#--printargs)
    * [--noWrap](#--nowrap)
    * [--printToStderr](#--printtostderr)
<!-- TOC -->

# The `selfkit` package

The `selfkit` package lets you test code that calls `os.Exit` — or
any code you need to run in a subprocess — by re-invoking the test
binary itself. Rather than building a separate helper binary, each
`go test` binary can act as its own controlled subprocess target.

`selfkit` works through a small set of CLI flags. When the test binary
is re-invoked with one of those flags, `Self.Run` handles the request
(prints output, sets the exit code) and tells `TestMain` to exit
without running tests.

## Setup

Wire `selfkit` into `TestMain` once per package that needs it:

```go
func TestMain(m *testing.M) {
    runTests, exitCode := selfkit.New().Run(os.Stdout, os.Stderr)
    if runTests {
        os.Exit(m.Run())
    }
    os.Exit(exitCode)
}
```

`selfkit.New()` reads `os.Args` by default. When the binary is
launched normally by the test runner, no selfkit flags are present,
so `Run` returns `(true, 0)` and tests proceed as usual.

## Writing a subprocess test

Use `exec.Command(os.Args[0], ...)` to re-invoke the test binary.
Pass `-test.run=^$` to suppress any actual test execution in the
subprocess, then add the selfkit flag that controls its behaviour:

```go
func Test_MyCode_exits(t *testing.T) {
    var stdout, stderr bytes.Buffer

    cmd := exec.Command(os.Args[0],
        "-test.run=^$",
        "--exitCode", "42",
    )
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()

    var exitErr *exec.ExitError
    if !errors.As(err, &exitErr) {
        t.Fatalf("expected exit error, got %v", err)
    }
    assert.Equal(t, 42, exitErr.ExitCode())
}
```

## Flags

### --toStdout and --toStderr

Print a value to stdout or stderr and exit without running tests.
The output is wrapped in a sentinel so the caller can distinguish
it from other binary noise:

```
--toStdout VALUE  →  |sout: VALUE|
--toStderr VALUE  →  |eout: VALUE|
```

```go
cmd := exec.Command(os.Args[0], "-test.run=^$",
    "--toStdout", "hello",
    "--toStderr", "world",
)
// stdout → "|sout: hello|"
// stderr → "|eout: world|"
```

Both flags may be combined in one invocation.

### --exitCode

Exit with the given code without running tests. Any integer value
including 0 is accepted; 0 explicitly means "exit success without
running tests".

```go
cmd := exec.Command(os.Args[0], "-test.run=^$",
    "--exitCode", "1",
)
// exits with code 1
```

### --printEnv

Read an environment variable by name and print its value to stdout
(wrapped as `|env: VALUE|`). Useful for asserting that a subprocess
inherits or sets a particular environment variable.

```go
cmd := exec.Command(os.Args[0], "-test.run=^$",
    "--printEnv", "MY_VAR",
)
cmd.Env = append(os.Environ(), "MY_VAR=secret")
// stdout → "|env: secret|"
```

### --printArgs

Print a label followed by any remaining (non-flag) arguments to
stdout, wrapped as `|args: LABEL REST1,REST2|`. Put `--printArgs`
last so trailing positional arguments are captured correctly.

```go
cmd := exec.Command(os.Args[0], "-test.run=^$",
    "--printArgs", "label", "arg1", "arg2",
)
// stdout → "|args: label arg1,arg2|"
```

### --noWrap

Remove the `|prefix: ...|` sentinel wrappers from all output. Use
this when the caller needs the raw value without delimiters.

```go
cmd := exec.Command(os.Args[0], "-test.run=^$",
    "--noWrap", "--toStdout", "hello",
)
// stdout → "hello"
```

### --printToStderr

Redirect the output of `--printEnv` and `--printArgs` to stderr
instead of stdout.

```go
cmd := exec.Command(os.Args[0], "-test.run=^$",
    "--printToStderr", "--printEnv", "MY_VAR",
)
// stderr → "|env: VALUE|"
// stdout → ""
```
