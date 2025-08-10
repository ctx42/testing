// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Len(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Len(tspy, 2, []int{0, 1})

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("fatal when want is greater than actual length", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		msg := affirm.Panic(t, func() { Len(tspy, 3, []int{0, 1}) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})

	t.Run("error - when want is less than the actual length", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Len(tspy, 1, []int{0, 1})

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
		have := Len(tspy, 1, []int{0, 1}, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Cap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Cap(tspy, 2, []int{0, 1})

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("fatal when want is greater than actual capacity", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		msg := affirm.Panic(t, func() { Cap(tspy, 3, []int{0, 1}) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})

	t.Run("error - when want is less than the actual capacity", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		s := make([]int, 0, 3)

		// --- When ---
		have := Cap(tspy, 2, s)

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
		have := Cap(tspy, 1, []int{0, 1}, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Has(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		val := []int{1, 2, 3}

		// --- When ---
		have := Has(tspy, 2, val)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := []int{1, 2, 3}

		// --- When ---
		have := Has(tspy, 42, val)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := []int{1, 2, 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := Has(tspy, 42, val, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_HasNo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		val := []int{1, 2, 3}

		// --- When ---
		have := HasNo(tspy, 4, val)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := []int{1, 2, 3}

		// --- When ---
		have := HasNo(tspy, 2, val)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := []int{1, 2, 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := HasNo(tspy, 2, val, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_HasKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		haveValue, haveHas := HasKey(tspy, "B", val)

		// --- Then ---
		affirm.Equal(t, 2, haveValue)
		affirm.Equal(t, true, haveHas)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		haveValue, haveHas := HasKey(tspy, "X", val)

		// --- Then ---
		affirm.Equal(t, 0, haveValue)
		affirm.Equal(t, false, haveHas)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		haveValue, haveHas := HasKey(tspy, "X", val, opt)

		// --- Then ---
		affirm.Equal(t, 0, haveValue)
		affirm.Equal(t, false, haveHas)
	})
}

func Test_HasNoKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		have := HasNoKey(tspy, "D", val)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		have := HasNoKey(tspy, "B", val)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := HasNoKey(tspy, "B", val, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_HasKeyValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		have := HasKeyValue(tspy, "B", 2, val)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		have := HasKeyValue(tspy, "B", 100, val)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := HasKeyValue(tspy, "B", 100, val, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_SliceSubset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		s0 := []string{"A", "B", "C"}
		s1 := []string{"C", "B", "A"}

		// --- When ---
		have := SliceSubset(tspy, s0, s1)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		sWant := []string{"X", "Y", "A", "B", "C"}
		sHave := []string{"C", "B", "A"}

		// --- When ---
		have := SliceSubset(tspy, sWant, sHave)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("           trail: type.field\n")
		tspy.Close()

		s0 := []string{"X", "Y", "A", "B", "C"}
		s1 := []string{"C", "B", "A"}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := SliceSubset(tspy, s0, s1, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_MapSubset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		mWant := map[string]string{
			"KEY0": "VAL0",
		}
		mHave := map[string]string{
			"KEY0": "VAL0",
			"KEY1": "VAL1",
		}

		// --- When ---
		have := MapSubset(tspy, mWant, mHave)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		m0 := map[string]string{
			"KEY0": "VAL0",
			"KEY1": "VAL1",
			"KEY2": "VAL2",
		}
		m1 := map[string]string{
			"KEY0": "VAL0",
			"KEY1": "VAL1",
		}

		// --- When ---
		have := MapSubset(tspy, m0, m1)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		m0 := map[string]string{
			"KEY0": "VAL0",
			"KEY1": "VAL1",
		}
		m1 := map[string]string{
			"KEY0": "VAL0",
		}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := MapSubset(tspy, m0, m1, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_MapsSubset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		w0 := []map[string]string{
			{"KEY0": "VAL0"},
			{"KEY0": "VAL0", "KEY1": "VAL1"},
		}
		w1 := []map[string]string{
			{"KEY0": "VAL0", "KEY1": "VAL1"},
			{"KEY0": "VAL0", "KEY1": "VAL1"},
		}

		// --- When ---
		have := MapsSubset(tspy, w0, w1)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		w0 := []map[string]string{
			{"KEY0": "VAL0", "KEY1": "VAL1", "KEY2": "VAL2"},
		}
		w1 := []map[string]string{
			{"KEY0": "VAL0", "KEY1": "VAL1"},
		}

		// --- When ---
		have := MapsSubset(tspy, w0, w1)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: <slice>[0]map[2]\n")
		tspy.Close()

		w0 := []map[int]int{
			{1: 10, 2: 20},
		}
		w1 := []map[int]int{
			{1: 10, 2: 200},
		}

		// --- When ---
		have := MapsSubset(tspy, w0, w1)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
