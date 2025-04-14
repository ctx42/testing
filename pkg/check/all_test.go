package check

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
)

// Flags for compiled test binary.
//
// When go runs tests it creates the binary with the test code pretty much in
// the same way it does compile the regular binaries, then this binary is run,
// and as every other binary can take flags.
//
// Below are flags used to trigger specific behaviours when test binary is run.
//
// If any of the flags are used the test binary does not run tests.
var (
	// exitCodeFlag represents compiled test binary flag, when set to value
	// greater or equal to 0 it will exit with that code without running tests.
	// If any of the above flags are set the binary will print values and then
	// exit.
	exitCodeFlag int
)

func init() {
	flag.IntVar(&exitCodeFlag, "exitCode", -1, "")
}

func TestMain(m *testing.M) {
	flag.Parse()
	// Exit with given code.
	if exitCodeFlag != -1 {
		os.Exit(exitCodeFlag)
	}
	os.Exit(m.Run())
}

// WithNow used to set custom function returning current time.
func WithNow(fn func() time.Time) Option {
	return func(ops Options) Options {
		ops.now = fn
		return ops
	}
}

func Test_WithNow(t *testing.T) {
	// --- Given ---
	ops := Options{}
	now := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// --- When ---
	have := WithNow(now)(ops)

	// --- Then ---
	affirm.Equal(t, true, core.Same(now, have.now))
}
