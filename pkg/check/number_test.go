// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Epsilon(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		// --- When ---
		err := Epsilon(42.0, 0.11, 41.9)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("int", func(t *testing.T) {
		// --- When ---
		err := Epsilon(42, 1, 41)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("int64", func(t *testing.T) {
		// --- When ---
		err := Epsilon(int64(42), int64(5), int64(47))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - float64", func(t *testing.T) {
		// --- When ---
		err := Epsilon(42.0, 0.11, 39.9)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected numbers to be within given epsilon:\n" +
			"     want: 42\n" +
			"     have: 39.9\n" +
			"  epsilon: 0.11\n" +
			"     diff: 2.1000000000000014"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - uint", func(t *testing.T) {
		// --- When ---
		err := Epsilon(uint(42), uint(4), uint(47))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected numbers to be within given epsilon:\n" +
			"     want: 42\n" +
			"     have: 47\n" +
			"  epsilon: 4\n" +
			"     diff: 5"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Epsilon(42, 4, 47, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected numbers to be within given epsilon:\n" +
			"    trail: type.field\n" +
			"     want: 42\n" +
			"     have: 47\n" +
			"  epsilon: 4\n" +
			"     diff: 5"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_EpsilonSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.123, 3.123}

		// --- When ---
		err := EpsilonSlice(s0, 0.01, s1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - different lengths", func(t *testing.T) {
		// --- Given ---
		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.143}

		// --- When ---
		err := EpsilonSlice(s0, 0.01, s1)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected []float64 length:\n" +
			"  want: 3\n" +
			"  have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.143, 3.123}

		// --- When ---
		err := EpsilonSlice(s0, 0.01, s1)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected all numbers in a slice to be within given epsilon respectively:\n" +
			"    trail: <[]float64>[1]\n" +
			"     want: 2.123\n" +
			"     have: 2.143\n" +
			"  epsilon: 0.01\n" +
			"     diff: 0.019999999999999574"
		affirm.Equal(t, wMsg, err.Error())
	})
}
