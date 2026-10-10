// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Empty(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Empty(tspy, "")

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
		have := Empty(tspy, "abc")

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
		have := Empty(tspy, "abc", check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotEmpty(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NotEmpty(tspy, "abc")

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
		have := NotEmpty(tspy, "")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field")
		tspy.Close()

		// --- When ---
		have := NotEmpty(tspy, "", check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
