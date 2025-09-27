// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_AnyString_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have any
		want bool
	}{
		{"string", "abc", true},
		{"integer", 123, false},
		{"bool", true, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := AnyString.Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_AnyInt_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have any
		want bool
	}{
		{"string", "abc", false},
		{"integer", 123, true},
		{"bool", true, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := AnyInt.Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_AnyBool_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have any
		want bool
	}{
		{"string", "abc", false},
		{"integer", 123, false},
		{"bool", true, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := AnyBool.Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_AnyCtx_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have any
		want bool
	}{
		{"context", context.Background(), true},
		{"nil context", context.Context(nil), false},
		{"nil any", nil, false},
		{"wrong type", 123, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := AnyCtx.Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_MatchSame(t *testing.T) {
	ptr0 := &testcases.TInt{}
	ptr1 := &testcases.TInt{}

	t.Run("same ptr", func(t *testing.T) {
		// --- When ---
		have := MatchSame(ptr0).Match(ptr0)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not same ptr", func(t *testing.T) {
		// --- When ---
		have := MatchSame(ptr0).Match(ptr1)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_MatchBy(t *testing.T) {
	t.Run("construct", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }

		// --- Then ---
		assert.True(t, MatchBy(fn).Match(42))
		assert.False(t, MatchBy(fn).Match(44))
	})

	t.Run("panics if more than one argument", func(t *testing.T) {
		// --- Given ---
		fn := func(a, b int) bool { return true }

		// --- Then ---
		msg := assert.PanicMsg(t, func() { MatchBy(fn) })
		assert.Contain(t, "mock: match function", *msg)
		assert.Contain(t, "does not take exactly one argument", *msg)
	})

	t.Run("panics if not returning bool", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) {}

		// --- Then ---
		msg := assert.PanicMsg(t, func() { MatchBy(fn) })
		assert.Contain(t, "mock: match function", *msg)
		assert.Contain(t, "does not return a bool", *msg)
	})

	t.Run("panics if not function", func(t *testing.T) {
		msg := assert.PanicMsg(t, func() { matcherFunc(42) })
		assert.Equal(t, *msg, "mock: \"int\" is not a match function")
	})
}

func Test_MatchBy_description_tabular(t *testing.T) {
	tt := []struct {
		testN string

		fn   any
		want string
	}{
		{"1", func(have int) bool { return false }, "[mock.MatchBy=func(int) bool]"},
		{"2", func(have bool) bool { return false }, "[mock.MatchBy=func(bool) bool]"},
		{"3", func(have ExampleItf) bool { return false }, "[mock.MatchBy=func(mock.ExampleItf) bool]"},
		{"4", func(have ExampleType) bool { return false }, "[mock.MatchBy=func(mock.ExampleType) bool]"},
		{"5", func(have *ExampleType) bool { return false }, "[mock.MatchBy=func(*mock.ExampleType) bool]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchBy(tc.fn)

			// --- Then ---
			assert.Equal(t, tc.want, have.Desc())
		})
	}
}

func Test_MatchOfType_match_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  string
		have any
		want bool
	}{
		{"1", "int", 42, true},
		{"2", "int", "42", false},
		{"3", "string", "str", true},
		{"4", "string", 42, false},
		{"5", "mock.ExampleItf", &ExampleType{}, false},
		{"6", "*mock.ExampleType", &ExampleType{}, true},
		{"7", "mock.ExampleType", ExampleType{}, true},
		{"8", "*mock.ExampleType", ExampleType{}, false},
		{"9", "mock.ExampleType", &ExampleType{}, false},
		{"10", "mock.ExampleItf", ExampleItf(&ExampleType{}), false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchOfType(tc.typ).Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_MatchOfType_description_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have string
		want string
	}{
		{"1", "int", "[mock.MatchOfType=int]"},
		{"2", "bool", "[mock.MatchOfType=bool]"},
		{"3", "mock.ExampleType", "[mock.MatchOfType=mock.ExampleType]"},
		{"4", "*mock.ExampleType", "[mock.MatchOfType=*mock.ExampleType]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchOfType(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have.Desc())
		})
	}
}

func Test_MatchType_match_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  any
		have any
		want bool
	}{
		{"1", 42, 44, true},
		{"2", 42.0, 42, false},
		{"3", "abc", "def", true},
		{"4", "abc", 42, false},
		{"5", true, true, true},
		{"6", true, false, true},
		{"7", &ExampleType{}, &ExampleType{}, true},
		{"8", ExampleType{}, ExampleType{}, true},
		{"9", ExampleType{}, &ExampleType{}, false},
		{"10", &ExampleType{}, ExampleType{}, false},
		{"11", &ExampleType{}, ExampleItf(&ExampleType{}), true},
		{"12", ExampleItf(&ExampleType{}), &ExampleType{}, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchType(tc.typ).Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_MatchType_description_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have any
		want string
	}{
		{"1", 42, "[mock.MatchType=int]"},
		{"2", true, "[mock.MatchType=bool]"},
		{"3", ExampleType{}, "[mock.MatchType=mock.ExampleType]"},
		{"4", &ExampleType{}, "[mock.MatchType=*mock.ExampleType]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchType(tc.have)

			// --- Then ---
			assert.Equal(t, tc.want, have.Desc())
		})
	}
}

func Test_MatchErrorContain_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have error
		exp  bool
	}{
		{"1", "long text", errors.New("some long text to match"), true},
		{"2", "other text", errors.New("some long text to match"), false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchErrorContain(tc.want).Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.exp, have)
		})
	}
}

func Test_MatchError_tabular(t *testing.T) {
	err0 := errors.New("test error 0")
	err1 := errors.New("test error 1")
	err3 := fmt.Errorf("wrapped: %w", err1)
	err4 := fmt.Errorf("two wrapped: %w %w", err0, err1)

	tt := []struct {
		testN string

		want any
		have error
		exp  bool
	}{
		{"match error message", "test error 0", err0, true},
		{"no match error message", "other error", err0, false},

		{"match error is", err0, err0, true},
		{"no mach error is", err0, err1, false},
		{"match wrapped error", err1, err3, true},
		{"match wrapped error multi 0", err0, err4, true},
		{"match wrapped error multi 1", err1, err4, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchError(tc.want).Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.exp, have)
		})
	}
}

func Test_MatchError(t *testing.T) {
	t.Run("panics when invalid type", func(t *testing.T) {
		// --- When ---
		msg := assert.PanicMsg(t, func() { MatchError(42) })

		// --- Then ---
		assert.Equal(t, "mock: MatchError: invalid type", *msg)
	})
}

func Test_AnySlice(t *testing.T) {
	// --- When ---
	have := AnySlice(3)

	// --- Then ---
	assert.Len(t, 3, have)
	assert.Equal(t, Any, have[0])
	assert.Equal(t, Any, have[1])
	assert.Equal(t, Any, have[2])
}
