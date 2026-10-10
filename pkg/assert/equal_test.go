// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

	t.Run("error - fails", func(t *testing.T) {
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

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		// --- When ---
		have := Equal(tspy, 42, 44, check.WithTrail("type.field"))

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

	t.Run("error - fails", func(t *testing.T) {
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

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		// --- When ---
		have := NotEqual(tspy, 42, 42, check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
