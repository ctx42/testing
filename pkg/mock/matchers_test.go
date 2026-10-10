// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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
		{"context", t.Context(), true},
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
	t.Run("error - not a function", func(t *testing.T) {
		// --- When ---
		have := assert.PanicMsg(t, func() { MatchBy(42) })

		// --- Then ---
		assert.Equal(t, `mock: "int" is not a match function`, *have)
	})

	t.Run("error - nil", func(t *testing.T) {
		// --- When ---
		have := assert.PanicMsg(t, func() { MatchBy(nil) })

		// --- Then ---
		assert.Equal(t, `mock: "<nil>" is not a match function`, *have)
	})

	t.Run("construct", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) bool { return have == 42 }

		// --- Then ---
		assert.True(t, MatchBy(fn).Match(42))
		assert.False(t, MatchBy(fn).Match(44))
	})

	t.Run("error - more than one argument", func(t *testing.T) {
		// --- Given ---
		fn := func(a, b int) bool { return true }

		// --- Then ---
		have := assert.PanicMsg(t, func() { MatchBy(fn) })
		assert.Contain(t, "mock: match function", *have)
		assert.Contain(t, "does not take exactly one argument", *have)
	})

	t.Run("error - not returning bool", func(t *testing.T) {
		// --- Given ---
		fn := func(have int) {}

		// --- Then ---
		msg := assert.PanicMsg(t, func() { MatchBy(fn) })
		assert.Contain(t, "mock: match function", *msg)
		assert.Contain(t, "does not return a bool", *msg)
	})

	t.Run("error - not function", func(t *testing.T) {
		have := assert.PanicMsg(t, func() { matcherFunc(42) })
		assert.Equal(t, *have, "mock: \"int\" is not a match function")
	})
}

func Test_MatchBy_tabular(t *testing.T) {
	tt := []struct {
		testN string

		fn   any
		want string
	}{
		{
			"int",
			func(have int) bool { return false },
			"[mock.MatchBy=func(int) bool]",
		},
		{
			"bool",
			func(have bool) bool { return false },
			"[mock.MatchBy=func(bool) bool]",
		},
		{
			"interface",
			func(have ExampleItf) bool { return false },
			"[mock.MatchBy=func(mock.ExampleItf) bool]",
		},
		{
			"struct",
			func(have ExampleType) bool { return false },
			"[mock.MatchBy=func(mock.ExampleType) bool]",
		},
		{
			"struct pointer",
			func(have *ExampleType) bool { return false },
			"[mock.MatchBy=func(*mock.ExampleType) bool]",
		},
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

func Test_MatchOfType_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  string
		have any
		want bool
	}{
		{"int", "int", 42, true},
		{"int string value", "int", "42", false},
		{"string", "string", "str", true},
		{"string int value", "string", 42, false},
		{"interface name", "mock.ExampleItf", &ExampleType{}, false},
		{"struct pointer", "*mock.ExampleType", &ExampleType{}, true},
		{"struct", "mock.ExampleType", ExampleType{}, true},
		{
			"struct pointer name and value",
			"*mock.ExampleType",
			ExampleType{},
			false,
		},
		{"struct name and pointer", "mock.ExampleType", &ExampleType{}, false},
		{
			"interface value",
			"mock.ExampleItf",
			ExampleItf(&ExampleType{}),
			false,
		},
		{"nil have", "int", nil, false},
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

func Test_MatchOfType_desc_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have string
		want string
	}{
		{"int", "int", "[mock.MatchOfType=int]"},
		{"bool", "bool", "[mock.MatchOfType=bool]"},
		{"struct", "mock.ExampleType", "[mock.MatchOfType=mock.ExampleType]"},
		{
			"struct pointer",
			"*mock.ExampleType",
			"[mock.MatchOfType=*mock.ExampleType]",
		},
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

func Test_MatchType(t *testing.T) {
	t.Run("same name different type", func(t *testing.T) {
		// --- Given ---
		val0 := func() any { type T struct{}; return T{} }()
		val1 := func() any { type T struct{}; return T{} }()
		am := MatchType(val0)

		// --- When ---
		have := am.Match(val1)

		// --- Then ---
		assert.False(t, have)
		assert.True(t, am.Match(val0))
	})

	t.Run("error - nil type", func(t *testing.T) {
		// --- When ---
		have := assert.PanicMsg(t, func() { MatchType(nil) })

		// --- Then ---
		assert.Equal(t, "mock: MatchType: typ must not be nil", *have)
	})
}

func Test_MatchType_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  any
		have any
		want bool
	}{
		{"same int", 42, 44, true},
		{"float and int", 42.0, 42, false},
		{"same string", "abc", "def", true},
		{"string and int", "abc", 42, false},
		{"same bool", true, true, true},
		{"bool values", true, false, true},
		{"struct pointers", &ExampleType{}, &ExampleType{}, true},
		{"structs", ExampleType{}, ExampleType{}, true},
		{"struct and pointer", ExampleType{}, &ExampleType{}, false},
		{"pointer and struct", &ExampleType{}, ExampleType{}, false},
		{
			"pointer and interface",
			&ExampleType{},
			ExampleItf(&ExampleType{}),
			true,
		},
		{
			"interface and pointer",
			ExampleItf(&ExampleType{}),
			&ExampleType{},
			true,
		},
		{"nil have", 42, nil, false},
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

func Test_MatchType_desc_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have any
		want string
	}{
		{"int", 42, "[mock.MatchType=int]"},
		{"bool", true, "[mock.MatchType=bool]"},
		{"struct", ExampleType{}, "[mock.MatchType=mock.ExampleType]"},
		{
			"struct pointer",
			&ExampleType{},
			"[mock.MatchType=*mock.ExampleType]",
		},
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

		want  string
		have  error
		match bool
	}{
		{"contains", "long text", errors.New("some long text to match"), true},
		{
			"does not contain",
			"other text",
			errors.New("some long text to match"),
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchErrorContain(tc.want).Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.match, have)
		})
	}
}

func Test_MatchError(t *testing.T) {
	t.Run("error - invalid type", func(t *testing.T) {
		// --- When ---
		have := assert.PanicMsg(t, func() { MatchError(42) })

		// --- Then ---
		assert.Equal(t, "mock: MatchError: invalid type", *have)
	})
}

func Test_MatchError_tabular(t *testing.T) {
	err0 := errors.New("test error 0")
	err1 := errors.New("test error 1")
	err3 := fmt.Errorf("wrapped: %w", err1)
	err4 := fmt.Errorf("two wrapped: %w %w", err0, err1)

	tt := []struct {
		testN string

		want  any
		have  error
		match bool
	}{
		{"match error message", "test error 0", err0, true},
		{"no match error message", "other error", err0, false},

		{"match error is", err0, err0, true},
		{"no match error is", err0, err1, false},
		{"match wrapped error", err1, err3, true},
		{"match wrapped error multi 0", err0, err4, true},
		{"match wrapped error multi 1", err1, err4, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MatchError(tc.want).Match(tc.have)

			// --- Then ---
			assert.Equal(t, tc.match, have)
		})
	}
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
