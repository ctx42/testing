// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"os"
	"os/exec"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_ExitCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		val := exec.Command(os.Args[0], "--exitCode", "0").Run()

		// --- When ---
		have := ExitCode(tspy, 0, val)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := exec.Command(os.Args[0], "--exitCode", "99").Run()

		// --- When ---
		have := ExitCode(tspy, 77, val)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := exec.Command(os.Args[0], "--exitCode", "99").Run()

		// --- When ---
		have := ExitCode(tspy, 77, val, check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
