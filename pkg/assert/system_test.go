// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"os"
	"os/exec"
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
	"github.com/ctx42/xtst/pkg/check"
	"github.com/ctx42/xtst/pkg/tester"
)

func Test_ExitCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "0").Run()

		// --- When ---
		have := ExitCode(tspy, 0, val)

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "99").Run()

		// --- When ---
		have := ExitCode(tspy, 77, val)

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("\ttrail: type.field\n")
		tspy.Close()

		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "99").Run()
		opt := check.WithTrail("type.field")

		// --- When ---
		have := ExitCode(tspy, 77, val, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}
