// Package selfkit enables test cases that need to exec a binary and
// assert its exit code, stdout, and stderr — by reusing the test
// binary itself. Wire it into TestMain via [New] and [Self.Run].
package selfkit

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// WithArgs is a [New] option overriding the argument slice. args must
// follow os.Args convention: args[0] is the program name and args[1:]
// are the flags to parse. Defaults to os.Args when not set.
func WithArgs(args []string) func(*Self) {
	return func(se *Self) { se.args = args }
}

// Self helps test cases where we would like to call a binary and expect the
// given exit code, standard output and/or standard error message.
//
// To do that, we can reuse the test binary itself.
//
// When go runs tests, it creates a binary with the test code pretty much in
// the same way as it compiles regular executables, then this binary is run
// with some arguments. The Self adds some additional flags that can be used to
// control the behavior of the test binary.
//
// Example:
//
//	func TestMain(m *testing.M) {
//		runTests, exitCode := selfkit.New().Run(os.Stdout, os.Stderr)
//		if runTests {
//			os.Exit(m.Run())
//		}
//		os.Exit(exitCode)
//	}
type Self struct {
	// toStdout represents the compiled test binary flag, when set it
	// will instruct the binary to print out its value to the stdout
	// and exit without running tests.
	toStdout string

	// toStderr represents the compiled test binary flag, when set it
	// will instruct the binary to print out its value to the stderr
	// and exit without running tests.
	toStderr string

	// printEnv represents the compiled test binary flag, when set it
	// will instruct the binary to print out the environment variable value
	// with the given name to the stdout and exit without running tests.
	printEnv string

	// printArgs represents the compiled test binary flag, when set it
	// will instruct the binary to print out not consumed arguments to the
	// stdout and exit without running tests. If you use this flag with other
	// flags put it as the last one with arguments to print after it.
	//
	// Example:
	//   program --toStdout abc --printArgs arg0 arg1
	printArgs string

	// printToStderr represents the compiled test binary flag redirecting
	// printouts from printEnv and printArgs flags to standard error instead of
	// standard output.
	//
	// Example:
	//   program --printToStdErr --printEnv ABC
	printToStderr bool

	// noWrap represents the compiled test binary flag turning off output
	// wrapping.
	noWrap bool

	// exitCode represents the compiled test binary flag, when set to value
	// greater or equal to 0, it will exit with that code without running tests.
	exitCode int

	fs   *flag.FlagSet
	args []string
}

// New returns new instance of [Self]. By default, [Self] uses [os.Args], but
// you can change it with the [WithArgs] option.
func New(opts ...func(*Self)) *Self {
	se := &Self{
		exitCode: -1,
		args:     os.Args,
	}
	for _, opt := range opts {
		opt(se)
	}
	if len(se.args) == 0 {
		return se
	}
	se.fs = flag.NewFlagSet(se.args[0], flag.ContinueOnError)
	se.fs.SetOutput(io.Discard)
	se.fs.StringVar(&se.toStdout, "toStdout", "", "")
	se.fs.StringVar(&se.toStderr, "toStderr", "", "")
	se.fs.StringVar(&se.printEnv, "printEnv", "", "")
	se.fs.StringVar(&se.printArgs, "printArgs", "", "")
	se.fs.BoolVar(&se.printToStderr, "printToStderr", false, "")
	se.fs.BoolVar(&se.noWrap, "noWrap", false, "")
	se.fs.IntVar(&se.exitCode, "exitCode", -1, "")
	// Unknown flags are intentionally ignored: the test runner passes
	// its own flags to the binary alongside ours.
	_ = se.fs.Parse(se.args[1:])

	return se
}

// Run executes the behavior requested by Self flags. It returns runTests
// (true when the caller should proceed with m.Run()) and exitCode (the
// value to pass to os.Exit). When runTests is false the binary should
// exit immediately without running any tests.
func (slf *Self) Run(stdout, stderr io.Writer) (bool, int) {
	// By default, run tests.
	runTests := true

	// Print to stdout.
	if slf.toStdout != "" {
		format := "|sout: %s|"
		if slf.noWrap {
			format = "%s"
		}
		_, _ = fmt.Fprintf(stdout, format, slf.toStdout)
		runTests = false
	}

	// Print to stderr.
	if slf.toStderr != "" {
		format := "|eout: %s|"
		if slf.noWrap {
			format = "%s"
		}
		_, _ = fmt.Fprintf(stderr, format, slf.toStderr)
		runTests = false
	}

	// Print to the value of an environment variable to stdout.
	if slf.printEnv != "" {
		out := stdout
		if slf.printToStderr {
			out = stderr
		}
		format := "|env: %s|"
		if slf.noWrap {
			format = "%s"
		}
		_, _ = fmt.Fprintf(out, format, os.Getenv(slf.printEnv))
		runTests = false
	}

	// Print command arguments to stdout.
	if slf.printArgs != "" {
		out := stdout
		if slf.printToStderr {
			out = stderr
		}
		format := "|args: %s %s|"
		if slf.noWrap {
			format = "%s %s"
		}
		args := strings.Join(slf.fs.Args(), ",")
		_, _ = fmt.Fprintf(out, format, slf.printArgs, args)
		runTests = false
	}

	// If specific code has not been requested, we return 0.
	exitCode := slf.exitCode
	if exitCode == -1 {
		exitCode = 0
	} else {
		runTests = false
	}

	return runTests, exitCode
}
