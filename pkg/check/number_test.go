// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Greater(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Greater(4, 2)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - equal", func(t *testing.T) {
		// --- When ---
		err := Greater(4, 4)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be greater:\n" +
			"  greater than: 4\n" +
			"          have: 4"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := Greater(2, 4)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be greater:\n" +
			"  greater than: 2\n" +
			"          have: 4"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Greater(2, 4, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be greater:\n" +
			"         trail: type.field\n" +
			"  greater than: 2\n" +
			"          have: 4"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_GreaterOrEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := GreaterOrEqual(4, 2)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - equal", func(t *testing.T) {
		// --- When ---
		err := GreaterOrEqual(4, 4)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := GreaterOrEqual(2, 4)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be greater or equal:\n" +
			"  greater or equal than: 2\n" +
			"                   have: 4"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := GreaterOrEqual(2, 4, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be greater or equal:\n" +
			"                  trail: type.field\n" +
			"  greater or equal than: 2\n" +
			"                   have: 4"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Smaller(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Smaller(2, 4)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - equal", func(t *testing.T) {
		// --- When ---
		err := Smaller(4, 4)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be smaller:\n" +
			"  smaller than: 4\n" +
			"          have: 4"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := Smaller(4, 2)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be smaller:\n" +
			"  smaller than: 4\n" +
			"          have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Smaller(4, 2, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be smaller:\n" +
			"         trail: type.field\n" +
			"  smaller than: 4\n" +
			"          have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_SmallerOrEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := SmallerOrEqual(2, 4)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - equal", func(t *testing.T) {
		// --- When ---
		err := SmallerOrEqual(4, 4)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := SmallerOrEqual(4, 2)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be smaller or equal:\n" +
			"  smaller or equal than: 4\n" +
			"                   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := SmallerOrEqual(4, 2, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be smaller or equal:\n" +
			"                  trail: type.field\n" +
			"  smaller or equal than: 4\n" +
			"                   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Delta(t *testing.T) {
	t.Run("success - delta less than expected", func(t *testing.T) {
		// --- When ---
		err := Delta(42.0, 0.11, 41.9)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - delta equal to expected", func(t *testing.T) {
		// --- When ---
		err := Delta(42, 1, 41)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("int64", func(t *testing.T) {
		// --- When ---
		err := Delta(int64(42), 6, int64(47))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := Delta(42.0, 2, 39.9)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected numbers to be within the given delta:\n" +
			"        want: 42\n" +
			"        have: 39.9\n" +
			"  want delta: 2\n" +
			"  have delta: 2.1000000000000014"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Delta(42, 0.10, 47, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected numbers to be within the given delta:\n" +
			"       trail: type.field\n" +
			"        want: 42\n" +
			"        have: 47\n" +
			"  want delta: 0.1\n" +
			"  have delta: 5"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_DeltaSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.123, 3.123}

		// --- When ---
		err := DeltaSlice(s0, 0.01, s1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - different lengths", func(t *testing.T) {
		// --- Given ---
		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.143}

		// --- When ---
		err := DeltaSlice(s0, 0.01, s1)

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
		err := DeltaSlice(s0, 0.009, s1)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected all numbers to be within the given delta respectively:\n" +
			"       trail: <[]float64>[1]\n" +
			"        want: 2.123\n" +
			"        have: 2.143\n" +
			"  want delta: 0.009\n" +
			"  have delta: 0.019999999999999574"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Epsilon(t *testing.T) {
	t.Run("success - epsilon less then expected", func(t *testing.T) {
		// --- When ---
		err := Epsilon(42.0, 0.11, 41.9)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - epsilon equal to expected", func(t *testing.T) {
		// --- When ---
		err := Epsilon(42.0, 1, 43.0)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := Epsilon(uint(42), 0.10, uint(47))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected numbers to be within the given epsilon:\n" +
			"          want: 42\n" +
			"          have: 47\n" +
			"  want epsilon: 0.1\n" +
			"  have epsilon: 0.11904761904761904"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Epsilon(42, 0.10, 47, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected numbers to be within the given epsilon:\n" +
			"         trail: type.field\n" +
			"          want: 42\n" +
			"          have: 47\n" +
			"  want epsilon: 0.1\n" +
			"  have epsilon: 0.11904761904761904"
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
		err := EpsilonSlice(s0, 0.009, s1)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected all numbers to be within the given epsilon respectively:\n" +
			"         trail: <[]float64>[1]\n" +
			"          want: 2.123\n" +
			"          have: 2.143\n" +
			"  want epsilon: 0.009\n" +
			"  have epsilon: 0.00942063118228901"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Increasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 3, 4}

		// --- When ---
		err := Increasing(seq)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - strict - previous equal current", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 2, 4}

		// --- When ---
		err := Increasing(seq)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected an increasing sequence:\n" +
			"     trail: <[]float64>[2]\n" +
			"      mode: strict\n" +
			"  previous: 2\n" +
			"   current: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - strict - not increasing", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 0.5, 3, 4}

		// --- When ---
		err := Increasing(seq)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected an increasing sequence:\n" +
			"     trail: <[]float64>[1]\n" +
			"      mode: strict\n" +
			"  previous: 1\n" +
			"   current: 0.5"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 2, 4}

		// --- When ---
		err := Increasing(seq, WithIncreasingSoft)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - empty slice", func(t *testing.T) {
		// --- When ---
		err := Increasing([]float64{})

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - nil slice", func(t *testing.T) {
		// --- When ---
		err := Increasing[int](nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 2, 4}
		opt := WithTrail("type.field")

		// --- When ---
		err := Increasing(seq, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected an increasing sequence:\n" +
			"     trail: type.field[2]\n" +
			"      mode: strict\n" +
			"  previous: 2\n" +
			"   current: 2"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NotIncreasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 2, 1}

		// --- When ---
		err := NotIncreasing(seq)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 2, 1}

		// --- When ---
		err := NotIncreasing(seq, WithIncreasingSoft)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - increasing strict", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 3, 4}

		// --- When ---
		err := NotIncreasing(seq)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a not increasing sequence:\n" +
			"  mode: strict"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - increasing soft", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 2, 4}

		// --- When ---
		err := NotIncreasing(seq, WithIncreasingSoft)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a not increasing sequence:\n" +
			"  mode: soft"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Decreasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 2, 1}

		// --- When ---
		err := Decreasing(seq)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - strict - previous equal current", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 3, 1}

		// --- When ---
		err := Decreasing(seq)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a decreasing sequence:\n" +
			"     trail: <[]float64>[2]\n" +
			"      mode: strict\n" +
			"  previous: 3\n" +
			"   current: 3"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - strict - not decreasing", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 0.5, 1}

		// --- When ---
		err := Decreasing(seq)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a decreasing sequence:\n" +
			"     trail: <[]float64>[3]\n" +
			"      mode: strict\n" +
			"  previous: 0.5\n" +
			"   current: 1"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 3, 1}

		// --- When ---
		err := Decreasing(seq, WithDecreasingSoft)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - empty slice", func(t *testing.T) {
		// --- When ---
		err := Decreasing([]float64{})

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - nil slice", func(t *testing.T) {
		// --- When ---
		err := Decreasing[int](nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 3, 1}
		opt := WithTrail("type.field")

		// --- When ---
		err := Decreasing(seq, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a decreasing sequence:\n" +
			"     trail: type.field[2]\n" +
			"      mode: strict\n" +
			"  previous: 3\n" +
			"   current: 3"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NotDecreasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 3, 4}

		// --- When ---
		err := NotDecreasing(seq)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		seq := []float64{1, 2, 3, 4}

		// --- When ---
		err := NotDecreasing(seq, WithDecreasingSoft)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - decreasing strict", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 2, 1}

		// --- When ---
		err := NotDecreasing(seq)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a not decreasing sequence:\n" +
			"  mode: strict"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - decreasing soft", func(t *testing.T) {
		// --- Given ---
		seq := []float64{4, 3, 3, 1}

		// --- When ---
		err := NotDecreasing(seq, WithDecreasingSoft)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected a not decreasing sequence:\n" +
			"  mode: soft"
		affirm.Equal(t, wMsg, err.Error())
	})
}
