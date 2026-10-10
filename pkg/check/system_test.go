// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_ExitCode(t *testing.T) {
	t.Run("error - nil error with trail", func(t *testing.T) {
		// --- When ---
		err := ExitCode(1, nil, WithTrail("type.field"))

		// --- Then ---
		wMsg := "" +
			"expected *exec.ExitError got nil:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("exit code 0", func(t *testing.T) {
		// --- Given ---
		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "0").Run()

		// --- When ---
		err := ExitCode(0, val)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("exit code 99", func(t *testing.T) {
		// --- Given ---
		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "99").Run()

		// --- When ---
		err := ExitCode(99, val)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not matching exit code", func(t *testing.T) {
		// --- Given ---
		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "99").Run()

		// --- When ---
		err := ExitCode(77, val)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected exit code:\n" +
			"  want: 77\n" +
			"  have: 99"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- When ---
		err := ExitCode(1, nil)

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "expected *exec.ExitError got nil", err.Error())
	})

	t.Run("error - no exec ExitError in the chain", func(t *testing.T) {
		// --- When ---
		err := ExitCode(1, errors.New("test"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected err to have \"*exec.ExitError\" in its chain"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := ExitCode(1, errors.New("test"), WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected err to have \"*exec.ExitError\" in its chain:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}
