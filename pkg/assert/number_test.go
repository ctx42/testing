// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Greater(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Greater(tspy, 42, 44)

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
		have := Greater(tspy, 42, 42)

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
		have := Greater(tspy, 44, 42, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_GreaterOrEqual(t *testing.T) {
	t.Run("success - greater", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := GreaterOrEqual(tspy, 42, 44)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := GreaterOrEqual(tspy, 44, 44)

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
		have := GreaterOrEqual(tspy, 44, 42)

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
		have := GreaterOrEqual(tspy, 44, 42, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Smaller(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Smaller(tspy, 44, 42)

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
		have := Smaller(tspy, 42, 44)

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
		have := Smaller(tspy, 42, 44, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_SmallerOrEqual(t *testing.T) {
	t.Run("success - smaller", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := SmallerOrEqual(tspy, 44, 42)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := SmallerOrEqual(tspy, 44, 44)

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
		have := SmallerOrEqual(tspy, 42, 44)

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
		have := SmallerOrEqual(tspy, 42, 44, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Delta(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Delta(tspy, 42.0, 0.11, 41.9)

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
		have := Delta(tspy, 42.0, 0.01, 39.9)

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
		have := Delta(tspy, 42.0, 0.01, 39.9, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_DeltaSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		s0 := []float64{1.123, 2.123, 3.123}
		s1 := []float64{1.123, 2.123, 3.123}

		// --- When ---
		have := DeltaSlice(tspy, s0, 0.01, s1)

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
		have := DeltaSlice(tspy, s0, 0.009, s1)

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
		have := DeltaSlice(tspy, s0, 0.009, s1, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

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
		have := Epsilon(tspy, 42.0, 0.01, 39.9)

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
		have := Epsilon(tspy, 42.0, 0.01, 39.9, opt)

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
		have := EpsilonSlice(tspy, s0, 0.009, s1)

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
		have := EpsilonSlice(tspy, s0, 0.009, s1, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Increasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{1, 2, 3, 4}

		// --- When ---
		have := Increasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{1, 2, 2, 4}

		// --- When ---
		have := Increasing(tspy, seq, check.WithIncreasingSoft)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		seq := []float64{1, 2, 1, 4}

		// --- When ---
		have := Increasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field[2]\n")
		tspy.Close()

		seq := []float64{1, 2, 1, 4}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := Increasing(tspy, seq, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotIncreasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{4, 3, 2, 1}

		// --- When ---
		have := NotIncreasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{4, 3, 2, 1}

		// --- When ---
		have := NotIncreasing(tspy, seq, check.WithIncreasingSoft)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - increasing strict", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		seq := []float64{1, 2, 3, 4}

		// --- When ---
		have := NotIncreasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - increasing soft", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		seq := []float64{1, 2, 2, 4}

		// --- When ---
		have := NotIncreasing(tspy, seq, check.WithIncreasingSoft)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Decreasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{4, 3, 2, 1}

		// --- When ---
		have := Decreasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{4, 3, 3, 1}

		// --- When ---
		have := Decreasing(tspy, seq, check.WithDecreasingSoft)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		seq := []float64{4, 3, 4, 1}

		// --- When ---
		have := Decreasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field[2]\n")
		tspy.Close()

		seq := []float64{4, 3, 4, 1}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := Decreasing(tspy, seq, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotDecreasing(t *testing.T) {
	t.Run("success - strict", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{1, 2, 3, 4}

		// --- When ---
		have := NotDecreasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success - soft", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		seq := []float64{1, 2, 3, 4}

		// --- When ---
		have := NotDecreasing(tspy, seq, check.WithDecreasingSoft)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - decreasing strict", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		seq := []float64{4, 3, 2, 1}

		// --- When ---
		have := NotDecreasing(tspy, seq)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - decreasing soft", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		seq := []float64{4, 3, 3, 1}

		// --- When ---
		have := NotDecreasing(tspy, seq, check.WithDecreasingSoft)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
