// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Equal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Equal(tspy, 42, 42)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Equal(tspy, 42, 44)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Equal(tspy, 42, 44, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NotEqual(tspy, 42, 44)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := NotEqual(tspy, 42, 42)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := NotEqual(tspy, 42, 42, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
