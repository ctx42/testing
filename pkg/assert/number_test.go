// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Epsilon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Epsilon(tspy, 42.0, 0.11, 41.9)

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
		have := Epsilon(tspy, 42.0, 0.11, 39.9)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Epsilon(tspy, 42.0, 0.11, 39.9, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_EpsilonSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.123, 3.123}

		// --- When ---
		have := EpsilonSlice(tspy, s0, 0.01, s1)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.143, 3.123}

		// --- When ---
		have := EpsilonSlice(tspy, s0, 0.01, s1)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field[1]\n")
		tspy.Close()

		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.143, 3.123}

		opt := check.WithTrail("type.field")

		// --- When ---
		have := EpsilonSlice(tspy, s0, 0.01, s1, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
