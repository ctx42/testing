// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewMatcher(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }
		want := "[mock.MyMatcher=func(int) bool]"

		// --- When ---
		am := NewMatcher(fn, want)

		// --- Then ---
		assert.Equal(t, want, am.desc)
		assert.True(t, am.fn.Call([]reflect.Value{reflect.ValueOf(42)})[0].Bool())
		assert.False(t, am.fn.Call([]reflect.Value{reflect.ValueOf(44)})[0].Bool())
	})

	t.Run("panics", func(t *testing.T) {
		// --- Given ---
		fn := func(int) {}

		// --- Then ---
		assert.Panic(t, func() { NewMatcher(fn, "") })
	})
}

func Test_Matcher_Desc(t *testing.T) {
	// --- Given ---
	want := "[mock.MyMatcher=func(int) bool]"

	// --- When ---
	am := NewMatcher(func(_ int) bool { return true }, want)

	// --- Then ---
	assert.Equal(t, want, am.Desc())
}

func Test_Matcher_Match(t *testing.T) {
	t.Run("matcher true", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(42)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("matcher false", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 44 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(42)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not assignable type", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(42.0)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("panics if type cannot be nil", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		msg := assert.PanicMsg(t, func() { am.Match(nil) })

		// --- Then ---
		want := "attempting to call matcher with nil for non-nil expected type"
		assert.Equal(t, want, *msg)
	})

	t.Run("nil-able types", func(t *testing.T) {
		// --- Given ---
		fnSlice := func(have []int) bool { return have == nil }
		fnMap := func(have map[string]int) bool { return have == nil }
		fnPtr := func(have *ExampleType) bool { return have == nil }
		fnItf := func(have ExampleItf) bool { return have == nil }
		fnFn := func(have func() int) bool { return have == nil }
		fnCh := func(have chan int) bool { return have == nil }

		// --- Then ---
		assert.True(t, NewMatcher(fnSlice, "mock.Custom").Match(nil))
		assert.True(t, NewMatcher(fnMap, "mock.Custom").Match(nil))
		assert.True(t, NewMatcher(fnPtr, "mock.Custom").Match(nil))
		assert.True(t, NewMatcher(fnItf, "mock.Custom").Match(nil))
		assert.True(t, NewMatcher(fnFn, "mock.Custom").Match(nil))
		assert.True(t, NewMatcher(fnCh, "mock.Custom").Match(nil))
	})

	t.Run("slice matcher", func(t *testing.T) {
		// --- Given ---
		fn := func(have []int) bool { return have[0] == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match([]int{42})

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("map matcher", func(t *testing.T) {
		// --- Given ---
		fn := func(have map[string]int) bool { return have["42"] == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(map[string]int{"42": 42})

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("pointer matcher", func(t *testing.T) {
		// --- Given ---
		fn := func(have *ExampleType) bool { return have.ran == true }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(&ExampleType{ran: true})

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("interface matcher", func(t *testing.T) {
		// --- Given ---
		fn := func(have ExampleItf) bool { return have.HasRan() == true }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(&ExampleType{ran: true})

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("function matcher", func(t *testing.T) {
		// --- Given ---
		fn := func(have func(v int) int) bool { return have(21) == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := am.Match(func(v int) int { return v * 2 })

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("channel matcher", func(t *testing.T) {
		// --- Given ---
		fn := func(have chan int) bool { return <-have == 42 }
		am := NewMatcher(fn, "mock.Custom")

		ch := make(chan int, 1)
		ch <- 42
		close(ch)

		// --- When ---
		have := am.Match(ch)

		// --- Then ---
		assert.True(t, have)
	})
}

func Test_matcherFunc(t *testing.T) {
	// --- Given ---
	var got int
	fn := func(a int) bool {
		got = a
		return true
	}

	// --- When ---
	val := matcherFunc(fn)

	// --- Then ---
	result := val.Call([]reflect.Value{reflect.ValueOf(42)})
	assert.True(t, result[0].Bool())
	assert.Equal(t, 42, got)
}

func Test_matcherFnValue_panic_tabular(t *testing.T) {
	tt := []struct {
		testN string

		fn  any
		exp string
	}{
		{"1", 42, "mock: \"int\" is not a match function"},
		{"2", func(a, b int) bool { return false }, " does not take exactly one argument"},
		{"3", func(a int) {}, " does not return a bool"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := assert.PanicMsg(t, func() { matcherFunc(tc.fn) })

			// --- Then ---
			assert.Contain(t, tc.exp, *have)
		})
	}
}
