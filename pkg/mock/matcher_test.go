// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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
		have := NewMatcher(fn, want)

		// --- Then ---
		assert.Equal(t, want, have.desc)
		assert.True(t, have.Match(42))
		assert.False(t, have.Match(44))
	})

	t.Run("error - not returning bool", func(t *testing.T) {
		// --- Given ---
		fn := func(int) {}

		// --- When ---
		have := assert.PanicMsg(t, func() { NewMatcher(fn, "") })

		// --- Then ---
		assert.Contain(t, " does not return a bool", *have)
	})
}

func Test_Matcher_Desc(t *testing.T) {
	// --- Given ---
	want := "[mock.MyMatcher=func(int) bool]"
	am := NewMatcher(func(_ int) bool { return true }, want)

	// --- When ---
	have := am.Desc()

	// --- Then ---
	assert.Equal(t, want, have)
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

	t.Run("error - type cannot be nil", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }
		am := NewMatcher(fn, "mock.Custom")

		// --- When ---
		have := assert.PanicMsg(t, func() { am.Match(nil) })

		// --- Then ---
		want := "attempting to call matcher with nil for non-nil expected type"
		assert.Equal(t, want, *have)
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

func Test_Matcher_Match_nil_tabular(t *testing.T) {
	tt := []struct {
		testN string

		fn any
	}{
		{"slice", func(have []int) bool { return have == nil }},
		{"map", func(have map[string]int) bool { return have == nil }},
		{"pointer", func(have *ExampleType) bool { return have == nil }},
		{"interface", func(have ExampleItf) bool { return have == nil }},
		{"func", func(have func() int) bool { return have == nil }},
		{"chan", func(have chan int) bool { return have == nil }},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			am := NewMatcher(tc.fn, "mock.Custom")

			// --- When ---
			have := am.Match(nil)

			// --- Then ---
			assert.True(t, have)
		})
	}
}

func Test_matcherFunc(t *testing.T) {
	// --- Given ---
	var arg int
	fn := func(a int) bool {
		arg = a
		return true
	}

	// --- When ---
	have := matcherFunc(fn)

	// --- Then ---
	result := have.Call([]reflect.Value{reflect.ValueOf(42)})
	assert.True(t, result[0].Bool())
	assert.Equal(t, 42, arg)
}

func Test_matcherFunc_panic_tabular(t *testing.T) {
	tt := []struct {
		testN string

		fn   any
		want string
	}{
		{"error - not a function", 42, "mock: \"int\" is not a match function"},
		{
			"error - two arguments",
			func(a, b int) bool { return false },
			" does not take exactly one argument",
		},
		{
			"error - not returning bool",
			func(a int) {},
			" does not return a bool",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := assert.PanicMsg(t, func() { matcherFunc(tc.fn) })

			// --- Then ---
			assert.Contain(t, tc.want, *have)
		})
	}
}
