// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"flag"
	"os"
	"testing"
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
