// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

	t.Run("not matching exit code", func(t *testing.T) {
		// --- Given ---
		cmd := os.Args[0]
		val := exec.Command(cmd, "--exitCode", "99").Run()

		// --- When ---
		err := ExitCode(77, val)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected exit code:\n" +
			"  want: 77\n" +
			"  have: 99"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		err := ExitCode(1, nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected *exec.ExitError got nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("no exec.ExitError in the chain", func(t *testing.T) {
		// --- When ---
		err := ExitCode(1, errors.New("test"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected err to have \"*exec.ExitError\" in its chain"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := ExitCode(1, errors.New("test"), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected err to have \"*exec.ExitError\" in its chain:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}
