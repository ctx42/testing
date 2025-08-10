// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Zero(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Zero(tspy, time.Time{})

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
		have := Zero(tspy, time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC))

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
		have := Zero(tspy, time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC), opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotZero(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NotZero(tspy, time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC))

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
		have := NotZero(tspy, time.Time{})

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
		have := NotZero(tspy, time.Time{}, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
