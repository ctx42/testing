// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Wait(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		fn := func() bool { return true }

		// --- When ---
		have := Wait(tspy, "100ms", fn)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - duration timeout", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		fn := func() bool { return true }

		// --- When ---
		have := Wait(tspy, 100*time.Millisecond, fn)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		fn := func() bool { return false }

		// --- When ---
		have := Wait(tspy, "100ms", fn)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		fn := func() bool { return false }

		// --- When ---
		have := Wait(tspy, "100ms", fn, check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
