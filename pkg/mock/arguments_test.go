// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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
		wantA := []any{"str", 42, true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: (int=42) == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("not matching two out of three", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", 42, true}
		haveA := []any{"str", 44, "false"}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 2, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: (int=42) != (int=44)`,
			`2: FAIL: (bool=true) != (string="false")`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("more want arguments", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", 42, true, "extra"}
		haveA := []any{"str", 44, false}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 3, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: (int=42) != (int=44)`,
			`2: FAIL: (bool=true) != (bool=false)`,
			`3: FAIL: (string="extra") != (Missing)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("more wnt any arguments", func(t *testing.T) {
		// --- Given ---
		wantA := []any{Any, Any, Any, Any, Any}
		haveA := []any{"A", "B", []any{"D", "E"}}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 2, cnt)
		want := []string{
			"0: PASS: (any=mock.Any) == (string=\"A\")",
			"1: PASS: (any=mock.Any) == (string=\"B\")",
			"2: PASS: (any=mock.Any) == ([]interface {}=[]interface {}{\"D\", \"E\"})",
			"3: FAIL: (any=mock.Any) != (Missing)",
			"4: FAIL: (any=mock.Any) != (Missing)",
		}
		assert.Equal(t, want, have)
	})

	t.Run("more have arguments", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", 42, true}
		haveA := []any{"str", 44, false, "extra"}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 3, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: (int=42) != (int=44)`,
			`2: FAIL: (bool=true) != (bool=false)`,
			`3: FAIL: (Missing) != (string="extra")`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching with one have argument set to Any", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", Any, true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: (any=mock.Any) == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching with all have arguments set to Any", func(t *testing.T) {
		// --- Given ---
		wantA := []any{Any, Any, Any}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			`0: PASS: (any=mock.Any) == (string="str")`,
			`1: PASS: (any=mock.Any) == (int=42)`,
			`2: PASS: (any=mock.Any) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching ArgumentMatcher", func(t *testing.T) {
		// --- Given ---
		mby := func(v int) bool { return v == 42 }
		wantA := []any{"str", MatchBy(mby), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchBy=func(int) bool] == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("not matching ArgumentMatcher - MatchBy", func(t *testing.T) {
		// --- Given ---
		mby := func(v int) bool { return v == 44 }
		wantA := []any{"str", MatchBy(mby), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchBy=func(int) bool] != (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("panicking ArgumentMatcher - MatchBy", func(t *testing.T) {
		// --- Given ---
		mby := func(v int) bool { panic("abc") }
		wantA := []any{"str", MatchBy(mby), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchBy=func(int) bool] {panic: "abc"} != (int=42)`,
			"2: PASS: (bool=true) == (bool=true)",
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching MatchOfType", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", MatchOfType("int"), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchOfType=int] == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching MatchOfType not matching other", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", MatchOfType("int"), true}
		haveA := []any{"str", 42, false}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchOfType=int] == (int=42)`,
			`2: FAIL: (bool=true) != (bool=false)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("not matching MatchOfType", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", MatchOfType("string"), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchOfType=string] != (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching MatchType", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", MatchType(0), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchType=int] == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("not matching MatchType", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", MatchType(""), true}
		haveA := []any{"str", 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: FAIL: [mock.MatchType=string] != (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("matching MatchType not matching other", func(t *testing.T) {
		// --- Given ---
		wantA := []any{"str", MatchType(0), true}
		haveA := []any{"str", 42, false}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 1, cnt)
		want := []string{
			`0: PASS: (string="str") == (string="str")`,
			`1: PASS: [mock.MatchType=int] == (int=42)`,
			`2: FAIL: (bool=true) != (bool=false)`,
		}
		assert.Equal(t, want, have)
	})

	t.Run("context arguments", func(t *testing.T) {
		// --- Given ---
		wantA := []any{context.Background(), 42, true}
		haveA := []any{context.Background(), 42, true}

		// --- When ---
		have, cnt := Arguments(wantA).Diff(haveA)

		// --- Then ---
		assert.Equal(t, 0, cnt)
		want := []string{
			"0: PASS: (context.backgroundCtx=context.backgroundCtx) == " +
				"(context.backgroundCtx=context.backgroundCtx)",
			`1: PASS: (int=42) == (int=42)`,
			`2: PASS: (bool=true) == (bool=true)`,
		}
		assert.Equal(t, want, have)
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
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *have)
	})

	t.Run("panics when argument cannot be cast to string", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.String(1) })

		// --- Then ---
		want := "[mock] arguments: String(1) is of type \"int\" not string"
		assert.Equal(t, want, *have)
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
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *have)
	})

	t.Run("panics when argument cannot be cast to int", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Int(2) })

		// --- Then ---
		want := "[mock] arguments: Int(2) is of type \"bool\" not int"
		assert.Equal(t, want, *have)
	})
}

func Test_Arguments_Float32(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", float32(3.14), true})

		// --- When ---
		have := args.Float32(1)

		// --- Then ---
		assert.Equal(t, float32(3.14), have)
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", float32(3.14), true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Float32(100) })

		// --- Then ---
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *have)
	})

	t.Run("panics when argument cannot be cast to float32", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", float32(3.14), true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Float32(2) })

		// --- Then ---
		want := "[mock] arguments: Float32(2) is of type \"bool\" not float32"
		assert.Equal(t, want, *have)
	})
}

func Test_Arguments_Float64(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 3.14159, true})

		// --- When ---
		have := args.Float64(1)

		// --- Then ---
		assert.Equal(t, 3.14159, have)
	})

	t.Run("panics for invalid index", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 3.14159, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Float64(100) })

		// --- Then ---
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *have)
	})

	t.Run("panics when argument cannot be cast to float64", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 3.14159, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Float64(2) })

		// --- Then ---
		want := "[mock] arguments: Float64(2) is of type \"bool\" not float64"
		assert.Equal(t, want, *have)
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
		assert.Same(t, err, have)
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
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *have)
	})

	t.Run("panics when argument cannot be cast to error", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		msg := assert.PanicMsg(t, func() { _ = args.Error(2) })

		// --- Then ---
		want := "[mock] arguments: Error(2) is of type \"bool\" not error"
		assert.Equal(t, want, *msg)
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
		want := "[mock] arguments: Get(100) out of range 2 max"
		assert.Equal(t, want, *have)
	})

	t.Run("panics when argument cannot be cast to bool", func(t *testing.T) {
		// --- Given ---
		args := Arguments([]any{"str", 42, true})

		// --- When ---
		have := assert.PanicMsg(t, func() { args.Bool(1) })

		// --- Then ---
		want := "[mock] arguments: Bool(1) is of type \"int\" not bool"
		assert.Equal(t, want, *have)
	})
}
