// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"context"
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Arguments_Get(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		var args = Arguments{"str", 42, true}

		// --- Then ---
		assert.Equal(t, "str", args.Get(0).(string))
		assert.Equal(t, 42, args.Get(1).(int))
		assert.Equal(t, true, args.Get(2).(bool))
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		var args = Arguments{"str", 42, true}

		// --- When ---
		msg := assert.PanicMsg(t, func() { args.Get(100) })

		// --- Then ---
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *msg)
	})
}

func Test_Arguments_Equal(t *testing.T) {
	t.Run("match", func(t *testing.T) {
		// --- Given ---
		var args = Arguments{"str", 42, true}

		// --- Then ---
		assert.True(t, args.Equal("str", 42, true))
	})

	t.Run("no match", func(t *testing.T) {
		// --- Given ---
		var args = Arguments{"str", 42, true}

		// --- Then ---
		assert.False(t, args.Equal("wrong", 456, false))
	})

	t.Run("panics for invalid indexes", func(t *testing.T) {
		// --- Given ---
		var args = Arguments{"str", 42, true}

		// --- When ---
		msg := assert.PanicMsg(t, func() { args.Equal("str", 42) })

		// --- Then ---
		want := "[must] arguments: argument lengths do not match 3 != 2"
		assert.Equal(t, want, *msg)
	})
}

func Test_Arguments_Diff(t *testing.T) {
	t.Run("matches with no differences", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", 42, true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: (int=42) == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("not matching two out of three", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", 42, true}
		have := []any{"str", 44, "false"}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 2, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: (int=42) != (int=44)`,
			`2: FAIL: (bool=true) != (string="false")`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("more want arguments", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", 42, true, "extra"}
		have := []any{"str", 44, false}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 3, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: (int=42) != (int=44)`,
			`2: FAIL: (bool=true) != (bool=false)`,
			`3: FAIL: (string="extra") != (Missing)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("more wnt any arguments", func(t *testing.T) {
		// --- Given ---
		want := []any{Any, Any, Any, Any, Any}
		have := []any{"A", "B", []any{"D", "E"}}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 2, cnt)
		exp := []string{
			"0: PASS: (any=mock.Any) == (string=\"A\")",
			"1: PASS: (any=mock.Any) == (string=\"B\")",
			"2: PASS: (any=mock.Any) == ([]interface {}=[]interface {}{\"D\", \"E\"})",
			"3: FAIL: (any=mock.Any) != (Missing)",
			"4: FAIL: (any=mock.Any) != (Missing)",
		}
		assert.Equal(t, exp, got)
	})

	t.Run("more have arguments", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", 42, true}
		have := []any{"str", 44, false, "extra"}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 3, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: (int=42) != (int=44)`,
			`2: FAIL: (bool=true) != (bool=false)`,
			`3: FAIL: (Missing) != (string="extra")`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching with one have argument set to Any", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", Any, true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: (any=mock.Any) == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching with all have arguments set to Any", func(t *testing.T) {
		// --- Given ---
		want := []any{Any, Any, Any}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			`0: PASS: (any=mock.Any) == (string="str")`,
			`1: PASS: (any=mock.Any) == (int=42)`,
			`2: PASS: (any=mock.Any) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching ArgumentMatcher", func(t *testing.T) {
		// --- Given ---
		mby := func(v int) bool { return v == 42 }
		want := []any{"str", MatchBy(mby), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchBy=func(int) bool] == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("not matching ArgumentMatcher - MatchBy", func(t *testing.T) {
		// --- Given ---
		mby := func(v int) bool { return v == 44 }
		want := []any{"str", MatchBy(mby), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchBy=func(int) bool] != (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("panicking ArgumentMatcher - MatchBy", func(t *testing.T) {
		// --- Given ---
		mby := func(v int) bool { panic("abc") }
		want := []any{"str", MatchBy(mby), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchBy=func(int) bool] {panic: "abc"} != (int=42)`,
			"2: PASS: (bool=true) == (bool=true)",
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching MatchOfType", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", MatchOfType("int"), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchOfType=int] == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching MatchOfType not matching other", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", MatchOfType("int"), true}
		have := []any{"str", 42, false}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchOfType=int] == (int=42)`,
			`2: FAIL: (bool=true) != (bool=false)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("not matching MatchOfType", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", MatchOfType("string"), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchOfType=string] != (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching MatchType", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", MatchType(0), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchType=int] == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("not matching MatchType", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", MatchType(""), true}
		have := []any{"str", 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchType=string] != (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("matching MatchType not matching other", func(t *testing.T) {
		// --- Given ---
		want := []any{"str", MatchType(0), true}
		have := []any{"str", 42, false}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		exp := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchType=int] == (int=42)`,
			`2: FAIL: (bool=true) != (bool=false)`,
		}
		assert.Equal(t, exp, got)
	})

	t.Run("context arguments", func(t *testing.T) {
		// --- Given ---
		want := []any{context.Background(), 42, true}
		have := []any{context.Background(), 42, true}

		// --- When ---
		got, cnt := Arguments(want).Diff(have)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		exp := []string{
			"0: PASS: (context.backgroundCtx=context.backgroundCtx) == " +
				"(context.backgroundCtx=context.backgroundCtx)",
			`1: PASS: (int=42) == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, exp, got)
	})
}

func Test_Arguments_String(t *testing.T) {
	t.Run("types representation", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := args.String(-1)

		// --- Then ---
		assert.Equal(t, `string, int, bool`, have)
	})

	t.Run("string by index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := args.String(0)

		// --- Then ---
		assert.Equal(t, "str", have)
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.String(100) })

		// --- Then ---
		exp := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, exp, *have)
	})

	t.Run("panics when argument cannot be cast to string", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.String(1) })

		// --- Then ---
		exp := "[mock] arguments: String(1) is of type \"int\" not string"
		assert.Equal(t, exp, *have)
	})
}

func Test_Arguments_Int(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := args.Int(1)

		// --- Then ---
		assert.Equal(t, 42, have)
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Int(100) })

		// --- Then ---
		exp := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, exp, *have)
	})

	t.Run("panics when argument cannot be cast to int", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Int(2) })

		// --- Then ---
		exp := "[mock] arguments: Int(2) is of type \"bool\" not int"
		assert.Equal(t, exp, *have)
	})
}

func Test_Arguments_Error(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		// --- Given ---
		err := errors.New("an error")
		args := Arguments([]any{"str", 42, true, err})

		// --- When ---
		have := args.Error(3)

		// --- Then ---
		assert.Equal(t, err, have)
	})

	t.Run("nil error", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true, nil})

		// --- When ---
		have := args.Error(3)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { _ = args.Error(100) })

		// --- Then ---
		exp := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, exp, *have)
	})

	t.Run("panics when argument cannot be cast to error", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		msg := assert.PanicMsg(t, func() { _ = args.Error(2) })

		// --- Then ---
		exp := "[mock] arguments: Error(2) is of type \"bool\" not error"
		assert.Equal(t, exp, *msg)
	})
}

func Test_Arguments_Bool(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := args.Bool(2)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Bool(100) })

		// --- Then ---
		exp := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, exp, *have)
	})

	t.Run("panics when argument cannot be cast to bool", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Bool(1) })

		// --- Then ---
		exp := "[mock] arguments: Bool(1) is of type \"int\" not bool"
		assert.Equal(t, exp, *have)
	})
}
