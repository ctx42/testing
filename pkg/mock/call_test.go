// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
)

func Test_newCall(t *testing.T) {
	// --- When ---
	call := newCall("Method", "arg0", "arg1")

	// --- Then ---
	assert.Nil(t, call.Stack)
	assert.Nil(t, call.parent)
	assert.Equal(t, "Method", call.Method)
	assert.Equal(t, Arguments{"arg0", "arg1"}, call.args)
	assert.False(t, call.argsAny)
	assert.Len(t, 0, call.returns)
	assert.Equal(t, 0, call.wantCalls)
	assert.Equal(t, 0, call.haveCalls)
	assert.False(t, call.optional)
	assert.Nil(t, call.until)
	assert.Duration(t, 0, call.after)
	assert.Nil(t, call.alter)
	assert.Nil(t, call.panic)
	assert.Nil(t, call.requires)
	assert.False(t, call.proxy.IsValid())
}

func Test_newProxy(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{}
		met := reflect.ValueOf(ptr.AAA)

		// --- When ---
		call := newProxy(met)

		// --- Then ---
		assert.Nil(t, call.Stack)
		assert.Nil(t, call.parent)
		assert.Equal(t, "AAA", call.Method)
		assert.Nil(t, call.args)
		assert.False(t, call.argsAny)
		assert.Len(t, 0, call.returns)
		assert.Equal(t, 0, call.wantCalls)
		assert.Equal(t, 0, call.haveCalls)
		assert.False(t, call.optional)
		assert.Nil(t, call.until)
		assert.Duration(t, 0, call.after)
		assert.Nil(t, call.alter)
		assert.Nil(t, call.panic)
		assert.Nil(t, call.requires)
		assert.Same(t, ptr.AAA, call.proxy.Interface())
	})

	t.Run("variadic", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{}
		met := reflect.ValueOf(ptr.Variadic)

		// --- When ---
		call := newProxy(met)

		// --- Then ---
		assert.Nil(t, call.Stack)
		assert.Nil(t, call.parent)
		assert.Equal(t, "Variadic", call.Method)
		assert.Nil(t, call.args)
		assert.False(t, call.argsAny)
		assert.Len(t, 0, call.returns)
		assert.Equal(t, 0, call.wantCalls)
		assert.Equal(t, 0, call.haveCalls)
		assert.False(t, call.optional)
		assert.Nil(t, call.until)
		assert.Duration(t, 0, call.after)
		assert.Nil(t, call.alter)
		assert.Nil(t, call.panic)
		assert.Nil(t, call.requires)
		assert.Same(t, ptr.Variadic, call.proxy.Interface())
	})

	t.Run("custom name", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{}
		met := reflect.ValueOf(ptr.Variadic)

		// --- When ---
		call := newProxy(met, "MyName")

		// --- Then ---
		assert.Equal(t, "MyName", call.Method)
	})
}

func Test_Call_withParent(t *testing.T) {
	// --- Given ---
	parent := &Mock{}
	call := newCall("Method")

	// --- When ---
	have := call.withParent(parent)

	// --- Then ---
	assert.Same(t, call, have)
	assert.Same(t, parent, call.parent)
}

func Test_Call_withStack(t *testing.T) {
	// --- Given ---
	stack := []string{"a", "b"}
	call := newCall("Method")

	// --- When ---
	have := call.withStack(stack)

	// --- Then ---
	assert.Same(t, call, have)
	assert.Equal(t, []string{"a", "b"}, call.Stack)
}

func Test_Call_With(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "c"}
		prx := reflect.ValueOf(ptr.AAA)
		call := newProxy(prx)

		// --- When ---
		have := call.With("a", "b")

		// --- Then ---
		assert.Same(t, call, have)
		assert.Equal(t, Arguments{"a", "b"}, have.args)
	})

	t.Run("panics if not proxied call", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		msg := assert.PanicMsg(t, func() { call.With("a", "b") })

		// --- Then ---
		assert.Contain(t, "cannot set arguments on proxy calls", *msg)
	})
}

func Test_Call_Return(t *testing.T) {
	t.Run("simple values", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.Return("str", 42, true)

		// --- Then ---
		assert.Same(t, call, have)
		assert.Equal(t, Arguments{"str", 42, true}, have.returns)
	})

	t.Run("with slices", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.Return("str", []any{42, true})

		// --- Then ---
		assert.Same(t, call, have)
		assert.Equal(t, Arguments{"str", []any{42, true}}, have.returns)
	})

	t.Run("nothing", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.Return()

		// --- Then ---
		assert.Same(t, call, have)
		assert.Nil(t, have.returns)
	})

	t.Run("panics if proxy call", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "c"}
		prx := reflect.ValueOf(ptr.AAA)
		call := newProxy(prx)

		// --- When ---
		have := assert.PanicMsg(t, func() { call.Return() })

		// --- Then ---
		assert.Equal(t, *have, "proxy calls cannot have return values")
	})
}

func Test_Call_Panic(t *testing.T) {
	t.Run("success with string", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.Panic("message")

		// --- Then ---
		assert.Same(t, call, have)
		assert.Equal(t, "message", have.panic)
	})

	t.Run("success with error", func(t *testing.T) {
		// --- Given ---
		ErrTst := errors.New("test")
		call := newCall("Zero")

		// --- When ---
		have := call.Panic(ErrTst)

		// --- Then ---
		assert.Same(t, call, have)
		assert.ErrorIs(t, ErrTst, have.panic.(error))
	})

	t.Run("panics if proxy call", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "c"}
		prx := reflect.ValueOf(ptr.AAA)
		call := newProxy(prx)

		// --- When ---
		have := assert.PanicMsg(t, func() { call.Panic("message") })

		// --- Then ---
		assert.Equal(t, *have, "cannot call panic on proxy calls")
	})
}

func Test_Call_Once(t *testing.T) {
	// --- Given ---
	call := newCall("Zero")

	// --- When ---
	have := call.Once()

	// --- Then ---
	assert.Same(t, call, have)
	assert.Equal(t, 1, have.wantCalls)
}

func Test_Call_Times(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.Times(3)

		// --- Then ---
		assert.Same(t, call, have)
		assert.Equal(t, 3, have.wantCalls)
	})

	t.Run("panics when used with Optional", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero").Optional()

		// --- When ---
		msg := assert.PanicMsg(t, func() { call.Times(3) })

		// --- Then ---
		assert.NotNil(t, msg)
		wMsg := "cannot use Optional and Times in the same time"
		assert.Equal(t, wMsg, *msg)
	})
}

func Test_Call_Until(t *testing.T) {
	// --- Given ---
	call := newCall("Zero")
	ch := time.After(50 * time.Millisecond)
	defer func() { <-ch }()

	// --- When ---
	have := call.Until(ch)

	// --- Then ---
	assert.Same(t, call, have)
	assert.Equal(t, ch, have.until)
}

func Test_Call_After(t *testing.T) {
	// --- Given ---
	call := newCall("Zero")

	// --- When ---
	have := call.After(100 * time.Millisecond)

	// --- Then ---
	assert.Same(t, call, have)
	assert.Equal(t, 100*time.Millisecond, have.after)
}

func Test_Call_Alter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		alter0 := func(_ Arguments) {}
		alter1 := func(_ Arguments) {}
		have := call.Alter(alter0, alter1)

		// --- Then ---
		assert.Same(t, call, have)
		assert.Same(t, alter0, have.alter[0])
		assert.Same(t, alter1, have.alter[1])
	})
}

func Test_Call_Optional(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.Optional()

		// --- Then ---
		assert.Same(t, call, have)
		assert.True(t, have.optional)
	})

	t.Run("panics when used with Times", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero").Times(3)

		// --- When ---
		msg := assert.PanicMsg(t, func() { call.Optional() })

		// --- Then ---
		assert.NotNil(t, msg)
		wMsg := "cannot use Optional and Times in the same time"
		assert.Equal(t, wMsg, *msg)
	})
}

func Test_Call_Requires(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		mck0 := &Mock{}
		call00 := newCall("Zero0").withParent(mck0)

		mck1 := &Mock{}
		call10 := newCall("Zero1").withParent(mck1)
		call11 := newCall("One1").withParent(mck1)

		// --- When ---
		have := call00.Requires(call10, call11)

		// --- Then ---
		assert.Same(t, call00, have)
		assert.Len(t, 2, have.requires)
		assert.Same(t, call10, have.requires[0])
		assert.Same(t, call11, have.requires[1])
	})

	t.Run("panics when nil is one of the arguments", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- Then ---
		msg := assert.PanicMsg(t, func() { call.Requires(nil) })
		assert.Contain(t, "nil instance", *msg)
	})
}

func Test_Call_CanCall_tabular(t *testing.T) {
	tt := []struct {
		testN string

		max      int
		cnt      int
		optional bool
		want     error
	}{
		{"no max never called not optional", 0, 0, false, nil},
		{"no max never called optional", 0, 0, true, nil},
		{"no max called not optional", 0, 1, false, nil},
		{"no max called optional", 0, 1, true, nil},
		{"called fewer times than expected not optional", 2, 0, false, nil},
		{"called fewer times than expected optional", 2, 0, true, nil},
		{"called requested number of times not optional", 1, 1, false, ErrTooManyCalls},
		{"called requested number of times optional", 1, 1, true, ErrTooManyCalls},
		{"called more than requested number of times not optional", 1, 2, false, ErrTooManyCalls},
		{"called more than requested number of times optional", 1, 2, false, ErrTooManyCalls},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			call := &Call{
				wantCalls: tc.max,
				haveCalls: tc.cnt,
			}

			// --- When ---
			err := call.CanCall()

			// --- Then ---
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func Test_Call_Satisfied(t *testing.T) {
	t.Run("satisfied", func(t *testing.T) {
		// --- Given ---
		call := &Call{
			wantCalls: 1,
			haveCalls: 1,
			optional:  false,
		}

		// --- When ---
		have := call.Satisfied()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not satisfied", func(t *testing.T) {
		// --- Given ---
		call := &Call{
			wantCalls: 1,
			haveCalls: 0,
			optional:  false,
		}

		// --- When ---
		have := call.Satisfied()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Call_satisfied_tabular(t *testing.T) {
	tt := []struct {
		testN string

		wantCalls int
		haveCalls int
		optional  bool
		err       error
	}{
		{"no max never called not optional", 0, 0, false, ErrNeverCalled},
		{"no max never called optional", 0, 0, true, nil},
		{"no max called not optional", 0, 1, false, nil},
		{"no max called not optional", 0, 1, true, nil},
		{"no max called multiple times optional", 0, 10, true, nil},
		{"no max called multiple times not optional", 0, 10, false, nil},
		{"called fewer times than max not optional", 2, 1, false, ErrTooFewCalls},
		{"called fewer times than max optional", 2, 1, true, nil},
		{"called requested number of times and optional", 5, 5, false, nil},
		{"called requested number of times optional", 5, 5, true, nil},
		{"called more times than max not optional", 1, 10, false, ErrTooManyCalls},
		{"called more times than max optional", 1, 10, true, ErrTooManyCalls},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			call := &Call{
				wantCalls: tc.wantCalls,
				// haveCalls: Method doesn't consider it.
				optional: tc.optional,
				args:     []any{1},
			}

			// --- When ---
			err := call.satisfied(tc.haveCalls)

			// --- Then ---
			assert.True(t, errors.Is(err, tc.err))
		})
	}
}

func Test_Call_satisfied(t *testing.T) {
	t.Run("not satisfied when method never called", func(t *testing.T) {
		// --- Given ---
		call := &Call{
			cStack:    cStack{Method: "Method"},
			wantCalls: 0,
			optional:  false,
			args:      []any{1},
			returns:   []any{2},
		}

		// --- When ---
		err := call.satisfied(0)

		// --- Then ---
		want := goldy.New(t, "testdata/satisfied_never.gld")
		assert.ErrorEqual(t, want.String(), err)
	})

	t.Run("not satisfied when the method is called too few times", func(t *testing.T) {
		// --- Given ---
		call := &Call{
			cStack:    cStack{Method: "Method"},
			wantCalls: 2,
			optional:  false,
			args:      []any{1},
			returns:   []any{2},
		}

		// --- When ---
		err := call.satisfied(1)

		// --- Then ---
		want := goldy.New(t, "testdata/satisfied_too_few.gld")
		assert.ErrorEqual(t, want.String(), err)
	})

	t.Run("not satisfied when the method is called too many times", func(t *testing.T) {
		// --- Given ---
		call := &Call{
			cStack:    cStack{Method: "Method"},
			wantCalls: 2,
			optional:  false,
			args:      []any{1},
			returns:   []any{2},
		}

		// --- When ---
		err := call.satisfied(3)

		// --- Then ---
		want := goldy.New(t, "testdata/satisfied_too_many.gld")
		assert.ErrorEqual(t, want.String(), err)
	})
}

func Test_Call_call(t *testing.T) {
	t.Run("without returns", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero")

		// --- When ---
		have := call.call()

		// --- Then ---
		assert.Nil(t, have)
		assert.Equal(t, 1, call.haveCalls)
	})

	t.Run("with until", func(t *testing.T) {
		// --- Given ---
		ch := time.After(50 * time.Millisecond)
		call := newCall("Zero").Until(ch)
		now := time.Now()

		// --- When ---
		have := call.call()

		// --- Then ---
		assert.True(t, time.Since(now) > 50*time.Millisecond)
		assert.Nil(t, have)
		assert.Equal(t, 1, call.haveCalls)
	})

	t.Run("with sleep", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero").After(50 * time.Millisecond)
		now := time.Now()

		// --- When ---
		have := call.call()

		// --- Then ---
		assert.True(t, time.Since(now) > 50*time.Millisecond)
		assert.Nil(t, have)
		assert.Equal(t, 1, call.haveCalls)
	})

	t.Run("with panic", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero").Panic("test panic")

		// --- When ---
		have := assert.PanicMsg(t, func() { call.call() })

		// --- Then ---
		assert.Equal(t, "test panic", *have)
		assert.Equal(t, 1, call.haveCalls)
	})

	t.Run("with panic after given time", func(t *testing.T) {
		// --- Given ---
		call := newCall("Zero").
			After(50 * time.Millisecond).
			Panic("test panic")
		now := time.Now()

		// --- When ---
		have := assert.PanicMsg(t, func() { call.call() })

		// --- Then ---
		assert.True(t, time.Since(now) > 50*time.Millisecond)
		assert.Equal(t, "test panic", *have)
		assert.Equal(t, 1, call.haveCalls)
	})

	t.Run("alter arguments", func(t *testing.T) {
		// --- Given ---
		arg0, arg1 := 0, 1

		alter0 := func(args Arguments) { *(args.Get(0).(*int))++ }
		alter1 := func(args Arguments) { *(args.Get(1).(*int)) += 2 }
		call := newCall("Zero").Alter(alter0, alter1)

		// --- When ---
		have := call.call(&arg0, &arg1)

		// --- Then ---
		assert.Nil(t, have)
		assert.Equal(t, 1, call.haveCalls)
		assert.Equal(t, 1, arg0)
		assert.Equal(t, 3, arg1)
	})

	t.Run("calls proxy method", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "c"}
		prx := reflect.ValueOf(ptr.Variadic)
		call := newProxy(prx)

		// --- When ---
		have := call.call("abc")

		// --- Then ---
		assert.Equal(t, Arguments{"c abc []"}, have)
		assert.Equal(t, 1, call.haveCalls)
	})

	t.Run("alter called before proxied call", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "c"}
		prx := reflect.ValueOf(ptr.Identity)
		arg := "abc"
		alter := func(args Arguments) { *(args.Get(0).(*string)) += " xyz" }
		call := newProxy(prx).Alter(alter)

		// --- When ---
		have := call.call(&arg)

		// --- Then ---
		want := "abc xyz"
		assert.Equal(t, Arguments{&want}, have)
		assert.Equal(t, 1, call.haveCalls)
	})
}

func Test_Call_proxyCall(t *testing.T) {
	t.Run("regular call", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "c"}
		prx := reflect.ValueOf(ptr.AAA)
		call := newProxy(prx)

		// --- When ---
		have := call.callProxy()

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "c", have[0])
		assert.Len(t, 0, call.returns)
	})

	t.Run("regular call with arguments", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "|"}
		prx := reflect.ValueOf(ptr.Wrap)
		call := newProxy(prx)

		// --- When ---
		have := call.callProxy("a", "b")

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "a|b", have[0])
	})

	t.Run("variadic proxy without variadic arguments", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "v"}
		prx := reflect.ValueOf(ptr.Variadic)
		call := newProxy(prx)

		// --- When ---
		have := call.callProxy("abc")

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "v abc []", have[0])
	})

	t.Run("variadic proxy with one variadic arguments", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "v"}
		prx := reflect.ValueOf(ptr.Variadic)
		call := newProxy(prx)

		// --- When ---
		have := call.callProxy("abc", 1, 2, 3)

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "v abc [1 2 3]", have[0])
	})

	t.Run("variadic proxy with multiple variadic arguments", func(t *testing.T) {
		// --- Given ---
		ptr := &types.TPtr{Val: "v"}
		prx := reflect.ValueOf(ptr.Variadic)
		call := newProxy(prx)

		// --- When ---
		have := call.callProxy("abc", 1)

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "v abc [1]", have[0])
	})

	t.Run("panics if proxy not defined", func(t *testing.T) {
		// --- Given ---
		call := newCall("Method")

		// --- When ---
		msg := assert.PanicMsg(t, func() { call.callProxy() })

		// --- Then ---
		assert.Equal(t, "proxy method not found", *msg)
	})
}

func Test_Call_checkReq(t *testing.T) {
	t.Run("no prerequisites", func(t *testing.T) {
		// --- Given ---
		call := newCall("Method")

		// --- When ---
		have := call.checkReq(nil)

		// --- Then ---
		assert.NoError(t, have)
	})

	t.Run("satisfied prerequisites", func(t *testing.T) {
		// --- Given ---
		pre0 := newCall("Pre0").satisfy()
		pre1 := newCall("Pre1").satisfy()
		call := newCall("Method").Requires(pre0, pre1)

		// --- When ---
		have := call.checkReq(nil)

		// --- Then ---
		assert.NoError(t, have)
	})

	t.Run("optional prerequisites are considered", func(t *testing.T) {
		// --- Given ---
		pre0 := newCall("Pre0").Optional()
		pre1 := newCall("Pre1").Optional()
		call := newCall("Method").Requires(pre0, pre1)

		// --- When ---
		have := call.checkReq(nil)

		// --- Then ---
		assert.NoError(t, have)
	})

	t.Run("error missing one from the same mock", func(t *testing.T) {
		// --- Given ---
		pre1 := newCall("Pre1", 1)
		stk := []string{"line0", "line1", "line2"}
		call := newCall("Method", "abc").Return(1, "abc").Requires(pre1)

		// --- When ---
		have := call.checkReq(stk)

		// --- Then ---
		want := goldy.New(t, "testdata/check_req_mock_same.gld")
		assert.ErrorEqual(t, want.String(), have)
		assert.ErrorIs(t, have, ErrRequirements)
	})

	t.Run("error missing one of many from the same mock", func(t *testing.T) {
		// --- Given ---
		pre0 := newCall("Pre0", 1).satisfy()
		pre1 := newCall("Pre1", 1)
		stk := []string{"line0", "line1", "line2"}
		call := newCall("Method", "abc").Return(1, "abc").Requires(pre0, pre1)

		// --- When ---
		have := call.checkReq(stk)

		// --- Then ---
		want := goldy.New(t, "testdata/check_req_mock_same.gld")
		assert.ErrorEqual(t, want.String(), have)
		assert.ErrorIs(t, have, ErrRequirements)
	})

	t.Run("error missing many from the same mock", func(t *testing.T) {
		// --- Given ---
		pre0 := newCall("Pre0", 1)
		pre1 := newCall("Pre1", 1)
		stk := []string{"line0", "line1", "line2"}
		call := newCall("Method", "abc").Return(1, "abc").Requires(pre0, pre1)

		// --- When ---
		have := call.checkReq(stk)

		// --- Then ---
		want := goldy.New(t, "testdata/check_req_many.gld")
		assert.ErrorEqual(t, want.String(), have)
		assert.ErrorIs(t, have, ErrRequirements)
	})

	t.Run("error missing from different mock", func(t *testing.T) {
		// --- Given ---
		mck0 := &Mock{}
		pre01 := newCall("Pre01").withParent(mck0)

		mck1 := &Mock{}
		newCall("Pre11").withParent(mck1)
		call := newCall("Method").Requires(pre01)

		// --- When ---
		have := call.checkReq(nil)

		// --- Then ---
		want := goldy.New(t, "testdata/check_req_mock_other.gld")
		assert.ErrorEqual(t, want.String(), have)
		assert.ErrorIs(t, have, ErrRequirements)
	})
}

func Test_Call_satisfy(t *testing.T) {
	t.Run("number of times not defined", func(t *testing.T) {
		// --- Given ---
		call := &Call{wantCalls: 0, haveCalls: 0}

		// --- When ---
		have := call.satisfy()

		// --- Then ---
		assert.Same(t, call, have)
		assert.True(t, call.Satisfied())
	})

	t.Run("number of times defined", func(t *testing.T) {
		// --- Given ---
		call := &Call{wantCalls: 5, haveCalls: 0}

		// --- When ---
		have := call.satisfy()

		// --- Then ---
		assert.Same(t, call, have)
		assert.True(t, call.Satisfied())
	})
}

func Test_Call_End(t *testing.T) {
	// --- Given ---
	mck := &Mock{}
	call := newCall("Zero").withParent(mck)

	// --- When ---
	have := call.End()

	// --- Then ---
	assert.Same(t, mck, have)
}
