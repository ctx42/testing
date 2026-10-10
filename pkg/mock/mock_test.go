// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/testcases"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_WithNoStack(t *testing.T) {
	// --- Given ---
	mck := &Mock{stack: true}

	// --- When ---
	WithNoStack(mck)

	// --- Then ---
	assert.False(t, mck.stack)
}

func Test_NewMock(t *testing.T) {
	t.Run("no expectations", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		mck := NewMock(tspy)

		// --- Then ---
		tspy.Finish()
		assert.Len(t, 0, mck.expected)
		assert.Len(t, 0, mck.calls)
		assert.Len(t, 0, mck.meta)
		assert.True(t, mck.stack)
		assert.False(t, mck.failed)
		assert.Same(t, tspy, mck.t)
	})

	t.Run("AssertExpectations called at the test end", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := goldy.Open(t, "testdata/mock_cleanup_assert_error.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0)

		// --- When ---
		tspy.Finish()

		// --- Then ---
		assert.True(t, mck.failed)
	})

	t.Run("with options", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		mck := NewMock(tspy, WithNoStack)

		// --- Then ---
		tspy.Finish()
		assert.False(t, mck.stack)
	})
}

func Test_Mock_MetaSetAll(t *testing.T) {
	t.Run("map key set", func(t *testing.T) {
		// --- Given ---
		mck := &Mock{}

		// --- When ---
		mck.MetaSetAll(map[string]any{"A": 1})

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, mck.meta)
	})

	t.Run("overrides", func(t *testing.T) {
		// --- Given ---
		mck := &Mock{meta: map[string]any{"A": 1}}
		meta := map[string]any{"B": 2}

		// --- When ---
		mck.MetaSetAll(meta)

		// --- Then ---
		assert.Equal(t, map[string]any{"B": 2}, mck.meta)
	})

	t.Run("stores a copy", func(t *testing.T) {
		// --- Given ---
		mck := &Mock{}
		meta := map[string]any{"A": 1}

		// --- When ---
		mck.MetaSetAll(meta)

		// --- Then ---
		meta["B"] = 2
		assert.Equal(t, map[string]any{"A": 1}, mck.meta)
	})
}

func Test_Mock_MetaAll(t *testing.T) {
	t.Run("does not return nil", func(t *testing.T) {
		// --- Given ---
		mck := &Mock{}

		// --- When ---
		have := mck.MetaAll()

		// --- Then ---
		assert.NotNil(t, have)
	})

	t.Run("returns a copy", func(t *testing.T) {
		// --- Given ---
		meta := map[string]any{"A": 1}
		mck := &Mock{meta: meta}

		// --- When ---
		have := mck.MetaAll()

		// --- Then ---
		assert.Equal(t, meta, have)

		have["B"] = 2
		assert.Equal(t, map[string]any{"A": 1}, mck.meta)
	})
}

func Test_Mock_On(t *testing.T) {
	t.Run("add expectation", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.On("Zero", "str", 42, true)

		// --- Then ---
		assert.Len(t, 1, mck.expected)
		assert.Same(t, mck.expected[0], have)
		assert.Same(t, mck, have.parent)
		assert.Equal(t, "Zero", have.Method)
		assert.Equal(t, Arguments{"str", 42, true}, have.args)
		assert.False(t, have.argsAny)
		assert.Len(t, 0, have.returns)
		assert.Len(t, 3, have.Stack)
		assert.Contain(t, "mock_test.go", have.Stack[2])
		assert.Equal(t, 0, have.wantCalls)
		assert.Equal(t, 0, have.haveCalls)
		assert.False(t, have.optional)
		assert.Nil(t, have.until)
		assert.Equal(t, time.Duration(0), have.after)
		assert.Nil(t, have.alter)
		assert.Nil(t, have.panic)
		assert.Len(t, 0, have.requires)
	})

	t.Run("error - functions are in args expectations", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)

		fn := func(string) error { return nil }

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.On("Zero", fn) })

		// --- Then ---
		assert.Equal(t, "cannot use functions in argument expectations", *have)
	})

	t.Run("chain calls", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.On("Zero", 0).Return("zero").End().
			On("One", 1).Return("one")

		// --- Then ---
		assert.Len(t, 2, mck.expected)
		assert.Same(t, mck.expected[1], have)
		assert.Equal(t, "Zero", mck.expected[0].Method)
		assert.Equal(t, Arguments{"zero"}, mck.expected[0].returns)
		assert.Equal(t, "One", mck.expected[1].Method)
		assert.Equal(t, Arguments{"one"}, mck.expected[1].returns)
	})

	t.Run("arg matcher slice", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		am := MatchBy(func(slice []bool) bool { return slice == nil })
		imp.On("MethodBoolS", am).Return(errors.New("fixture1"))

		// --- When ---
		have := imp.MethodBoolS(nil)

		// --- Then ---
		assert.ErrorEqual(t, "fixture1", have)
	})

	t.Run("error - with variadic method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_unexpected_variadic.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy, WithNoStack))

		// --- When ---
		have := imp.On("MethodIntVar", 1, 2, 3).Return(nil)

		// --- Then ---
		assert.Len(t, 1, imp.expected)
		assert.Same(t, have, imp.expected[0])
		assert.Len(t, 3, have.args)
		assert.Equal(t, Arguments{1, 2, 3}, have.args)

		assert.NoPanic(t, func() { _ = imp.MethodIntVar(1, 2, 3) })
		assert.Panic(t, func() { _ = imp.MethodIntVar(1, 2) })
	})
}

func Test_Mock_OnAny(t *testing.T) {
	t.Run("add expectation", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.OnAny("Zero")

		// --- Then ---
		assert.Len(t, 1, mck.expected)
		assert.Same(t, mck.expected[0], have)
		assert.Same(t, mck, have.parent)
		assert.Equal(t, "Zero", have.Method)
		assert.Empty(t, have.args)
		assert.True(t, have.argsAny)
		assert.Len(t, 0, have.returns)
		assert.Len(t, 3, have.Stack)
		assert.Contain(t, "mock_test.go", have.Stack[2])
		assert.Equal(t, 0, have.wantCalls)
		assert.Equal(t, 0, have.haveCalls)
		assert.False(t, have.optional)
		assert.Nil(t, have.until)
		assert.Equal(t, time.Duration(0), have.after)
		assert.Nil(t, have.alter)
		assert.Nil(t, have.panic)
		assert.Len(t, 0, have.requires)
	})

	t.Run("chain calls", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.OnAny("Zero").Return("zero").End().
			OnAny("One").Return("one")

		// --- Then ---
		assert.Len(t, 2, mck.expected)
		assert.Same(t, mck.expected[1], have)
		assert.Equal(t, "Zero", mck.expected[0].Method)
		assert.Equal(t, Arguments{"zero"}, mck.expected[0].returns)
		assert.Equal(t, "One", mck.expected[1].Method)
		assert.Equal(t, Arguments{"one"}, mck.expected[1].returns)
	})

	t.Run("variadic method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))

		// --- When ---
		have := imp.OnAny("MethodIntVar").Return(nil)

		// --- Then ---
		assert.Len(t, 1, imp.expected)
		assert.Same(t, have, imp.expected[0])
		assert.Len(t, 0, have.args)

		assert.NoPanic(t, func() { _ = imp.MethodIntVar(1, 2, 3) })
		assert.NoPanic(t, func() { _ = imp.MethodIntVar(1, 2) })
	})
}

func Test_Mock_Proxy(t *testing.T) {
	t.Run("with a default name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		ptr := &testcases.TPtr{Val: "abc"}

		// --- When ---
		have := mck.Proxy(ptr.AAA)

		// --- Then ---
		assert.Same(t, mck, have.parent)
		assert.Nil(t, have.args)
		assert.Equal(t, "AAA", have.Method)
	})

	t.Run("with a custom name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		ptr := &testcases.TPtr{Val: "abc"}

		// --- When ---
		have := mck.Proxy(ptr.AAA, "MyName")

		// --- Then ---
		assert.Same(t, mck, have.parent)
		assert.Nil(t, have.args)
		assert.Equal(t, "MyName", have.Method)
	})

	t.Run("error - not method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Proxy(123) })

		// --- Then ---
		assert.Equal(t, "Proxy requires a valid not nil method", *have)
	})

	t.Run("error - invalid", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		type T func()
		var fn T
		mck := NewMock(tspy)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Proxy(fn) })

		// --- Then ---
		assert.Equal(t, "Proxy requires a valid not nil method", *have)
	})
}

func Test_Mock_Called(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodInts", 1, 2, 3).Return(4)

		// --- When ---
		have, err := imp.MethodInts(1, 2, 3)

		// --- Then ---
		assert.ErrorEqual(t, "whoops", err)
		assert.Equal(t, 4, have)
		assert.Len(t, 1, imp.calls)
		assert.Equal(t, "MethodInts", imp.calls[0].Method)
		assert.Len(t, 5, imp.calls[0].Stack)
		assert.Contain(t, "mock_test.go", imp.calls[0].Stack[4])
	})

	t.Run("error - unexpected method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_unexpected.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy, WithNoStack))

		// --- When ---
		have := assert.PanicMsg(t, func() { imp.Called(1, 2, 3) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
	})
}

func Test_Mock_called(t *testing.T) {
	t.Run("self", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.called(0)

		// --- Then ---
		assert.Equal(t, "called", have)
	})

	t.Run("sub test", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.called(1)

		// --- Then ---
		assert.Equal(t, "<anonymous>", have)
	})

	t.Run("error - invalid skip value", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.called(100) })

		// --- Then ---
		assert.Equal(t, "could not get the caller information", *have)
	})
}

func Test_isClosureName_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name string
		want bool
	}{
		{"func literal", "func1", true},
		{"nested func literal", "2", true},
		{"func prefix", "funcX", false},
		{"func", "func", false},
		{"method", "Method", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := isClosureName(tc.name)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_Mock_Call(t *testing.T) {
	t.Run("call existing", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		call0 := mck.On("Zero", 0)

		// --- When ---
		have := mck.Call("Zero", 0)

		// --- Then ---
		assert.Len(t, 0, have)
		assert.Equal(t, 0, call0.wantCalls)
		assert.Equal(t, 1, call0.haveCalls)
		assert.False(t, mck.failed)
	})

	t.Run("call existing with limit", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		call0 := mck.On("Zero", 0).Times(2)

		// --- When ---
		have := mck.Call("Zero", 0)

		// --- Then ---
		assert.Len(t, 0, have)
		assert.Equal(t, 2, call0.wantCalls)
		assert.Equal(t, 1, call0.haveCalls)
		assert.False(t, mck.failed)
	})

	t.Run("return values", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		call0 := mck.On("Zero", 0).Return("zero")

		// --- When ---
		have := mck.Call("Zero", 0)

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "zero", have[0])
		assert.Equal(t, 0, call0.wantCalls)
		assert.Equal(t, 1, call0.haveCalls)
		assert.False(t, mck.failed)
	})

	t.Run("error - when a method is called too many times", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/mock_too_many_calls.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0, 1).Once()
		_ = mck.Call("Zero", 0, 1)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Zero", 0, 1) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("error - different arguments", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_found_args_dont_match.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		mck.On("Zero", 0, 1)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Zero", 2, 3) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("error - call to not expected method name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_not_found_with_args.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		mck.On("Zero", 0)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("One", 1) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("error - unmet requirements", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_deps_not_met.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		req0 := mck.On("Zero", 0)
		req1 := mck.On("One", 1).Return("one")
		mck.On("Two", 2).Return("two").Requires(req0, req1)

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Two", 2) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("error - when not all requirements are met", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_deps_some_not_met.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		req0 := mck.On("Zero", 0)
		req1 := mck.On("One", 1).Return("one")
		mck.On("Two", 2).Return("two").Requires(req0, req1)
		mck.Call("Zero", 0) // Satisfy one of the requirements.

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Two", 2) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("error - when time does not match exactly", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_deps_not_met_timezone.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		tim := time.Date(2020, 1, 2, 3, 4, 5, 0, testcases.WAW)

		mck := NewMock(tspy, WithNoStack)
		mck.On("Time", tim.UTC())

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Time", tim) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("error - requirement values did not match", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_deps_not_met_values.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		req0 := mck.On("Zero", 0).Return("zero 0")
		req1 := mck.On("Zero", 1).Return("zero 1")
		mck.On("Two", 2).Return("two").Requires(req0, req1)
		mck.Call("Zero", 0) // Satisfy one of the requirements.

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Two", 2) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, mck.failed)
	})

	t.Run("skips expectation with unmet requirements", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		req := mck.On("Zero", 0)
		mck.On("Two", 2).Return("two a").Requires(req)
		mck.On("Two", 2).Return("two b")

		// --- When ---
		have := mck.Call("Two", 2)

		// --- Then ---
		assert.Equal(t, Arguments{"two b"}, have)
		assert.False(t, mck.failed)
	})

	t.Run("wait for until", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		ch := time.After(500 * time.Millisecond)

		mck := NewMock(tspy)
		mck.On("MethodBool", true).Return("done").Until(ch)

		var have Arguments
		done := make(chan struct{})
		start := time.Now()

		// --- When ---
		go func() {
			have = mck.Call("MethodBool", true)
			close(done)
		}()

		// --- Then ---
		assert.ChannelWillClose(t, time.Second, done)
		assert.True(t, time.Since(start) > 490*time.Millisecond)
		assert.Equal(t, Arguments{"done"}, have)
		assert.False(t, mck.failed)
	})

	t.Run("wait for sleep", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("MethodBool", true).Return("done").After(500 * time.Millisecond)

		var have Arguments
		done := make(chan struct{})
		start := time.Now()

		// --- When ---
		go func() {
			have = mck.Call("MethodBool", true)
			close(done)
		}()

		// --- Then ---
		assert.ChannelWillClose(t, time.Second, done)
		assert.True(t, time.Since(start) > 500*time.Millisecond)
		assert.Equal(t, Arguments{"done"}, have)
		assert.False(t, mck.failed)
	})

	t.Run("until released by concurrent Call", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		ch := make(chan time.Time)
		mck := NewMock(tspy)
		mck.On("Wait").Return().Until(ch)
		mck.On("Done").Return().Alter(func(Arguments) { close(ch) })

		done := make(chan struct{})

		// --- When ---
		go func() {
			_ = mck.Call("Wait")
			close(done)
		}()
		// Give Wait a chance to block before Done runs.
		time.Sleep(20 * time.Millisecond)
		_ = mck.Call("Done")

		// --- Then ---
		assert.ChannelWillClose(t, time.Second, done)
		assert.False(t, mck.failed)
	})

	t.Run("mocked method should panic", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Panic("zero panics")

		// --- When ---
		have := assert.PanicMsg(t, func() { mck.Call("Zero", 0) })

		// --- Then ---
		assert.Equal(t, "zero panics", *have)
		assert.False(t, mck.failed)
	})

	t.Run("runs in order", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		m := map[string]any{
			"k0": "v0",
		}

		mck := NewMock(tspy)
		mck.On("Zero", MatchType(m)).
			Alter(func(args Arguments) {
				arg := args.Get(0).(map[string]any)
				arg["k1"] = "v1"
			}).
			Alter(func(args Arguments) {
				arg := args.Get(0).(map[string]any)
				arg["k1"] = "v2"
			})

		// --- When ---
		have := mck.Call("Zero", m)

		// --- Then ---
		assert.Nil(t, have)
		assert.Equal(t, "v2", m["k1"])
		assert.False(t, mck.failed)
	})

	t.Run("parallel simple", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		call := mck.On("ConcurrencyTestMethod", 1).Return(1)

		const cnt = 1000
		wg := sync.WaitGroup{}
		wg.Add(2)

		// --- When ---
		go func() {
			// Continuously call the mocked method.
			for range cnt {
				_ = mck.Call("ConcurrencyTestMethod", 1)
			}
			wg.Done()
		}()

		go func() {
			// Continuously call the mocked method.
			for range cnt {
				_ = mck.Call("ConcurrencyTestMethod", 1)
			}
			wg.Done()
		}()

		// --- Then ---
		wg.Wait()
		assert.Equal(t, cnt*2, call.haveCalls)
		assert.False(t, mck.failed)
	})

	t.Run("call proxy method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		ptr := &testcases.TPtr{Val: "b"}
		mck := NewMock(tspy)
		mck.Proxy(ptr.Wrap)

		// --- When ---
		have := mck.Call("Wrap", "a", "c")

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "abc", have[0])
		assert.False(t, mck.failed)
	})

	t.Run("call variadic proxy method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		ptr := &testcases.TPtr{Val: "b"}
		mck := NewMock(tspy)
		mck.Proxy(ptr.Variadic)

		// --- When ---
		have := mck.Call("Variadic", "a", 1, 2, 3)

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Equal(t, "b a [1 2 3]", have[0])
		assert.False(t, mck.failed)
	})
}

func Test_Mock_Callable(t *testing.T) {
	t.Run("error - too many calls registered first", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		mck.On("A", 1).Once()
		mck.On("A", 2).Optional()
		_ = mck.Call("A", 1)

		// --- When ---
		err := mck.Callable("A", 1)

		// --- Then ---
		assert.ErrorIs(t, ErrTooManyCalls, err)
	})

	t.Run("error - too many calls registered last", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)
		mck.On("A", 2).Optional()
		mck.On("A", 1).Once()
		_ = mck.Call("A", 1)

		// --- When ---
		err := mck.Callable("A", 1)

		// --- Then ---
		assert.ErrorIs(t, ErrTooManyCalls, err)
	})

	t.Run("callable", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodIntVar", 42, 44).Return(nil)

		// --- When ---
		err := imp.Callable("MethodIntVar", 42, 44)

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("error - when not existing method name given", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodIntVar", 42, 44).Return(nil)

		// --- When ---
		err := imp.Callable("NotExisting")

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		assert.Equal(t, hNotFoundCall, notice.From(err).HeaderText)
	})

	t.Run("error - when not matching arguments given", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodIntVar", 42, 44).Return(nil)

		// --- When ---
		err := imp.Callable("MethodIntVar", 7, 44)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		assert.Equal(t, hUnexpectedCall, notice.From(err).HeaderText)
	})

	t.Run("error - when called again", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodIntVar", 42, 44).Once().Return(nil)
		assert.Nil(t, imp.MethodIntVar(42, 44))

		// --- When ---
		err := imp.Callable("MethodIntVar", 42, 44)

		// --- Then ---
		assert.ErrorIs(t, ErrTooManyCalls, err)
	})

	t.Run("error - when prerequisites not satisfied", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		pre := mck.On("Init").Return()
		mck.On("Do").Return().Requires(pre)

		// --- When ---
		err := mck.Callable("Do")

		// --- Then ---
		assert.ErrorIs(t, ErrRequirements, err)
	})

	t.Run("skips expectation with unmet prerequisites", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		pre := mck.On("Init").Return()
		mck.On("Do").Return().Requires(pre)
		mck.On("Do").Return()

		// --- When ---
		err := mck.Callable("Do")

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("error - prerequisites before too many calls", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		pre := mck.On("Init").Return()
		mck.On("Do").Once().Return().addCall()
		mck.On("Do").Return().Requires(pre)

		// --- When ---
		err := mck.Callable("Do")

		// --- Then ---
		assert.ErrorIs(t, ErrRequirements, err)
	})
}

func Test_Mock_find(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Return("zero")
		mck.On("One", 1).Return("one")
		want := mck.On("One", 2).Return("two")

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("error - matching name not matching arg count", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("One", 1).Return("one")

		// --- When ---
		have, err := mck.find("One", nil, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_fail_arg_count.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("error - method name call not found", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Return("zero")

		// --- When ---
		have, err := mck.find("Two", nil, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_call_fail_no_args.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("error - not found with arguments", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Return("zero")

		// --- When ---
		have, err := mck.find("Two", []any{1, 2}, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_call_fail_with_args.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("error - matching name not matching arg type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Return("zero")
		mck.On("One", 1).Return("one")

		// --- When ---
		have, err := mck.find("One", []any{1.0}, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_call_fail_arg_type.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("respects wantCalls", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Return("zero")
		mck.On("One", 1).Once().Return("one").addCall()
		want := mck.On("One", 1).Return("two")

		// --- When ---
		have, err := mck.find("One", []any{1}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("method with variadic args", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", true)
		want := imp.On("MethodIntVar", 1, 2, 3).Return(nil)

		// --- When ---
		have, err := imp.find("MethodIntVar", []any{1, 2, 3}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("error - method with variadic args not found", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", true)
		imp.On("MethodIntVar", 1, 2, 3, 4).Return(nil)

		// --- When ---
		have, err := imp.find("MethodIntVar", []any{1, 2, 3}, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_fail_variadic_arg_count.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("with argument matcher", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))

		by := MatchBy(func(have int) bool { return have < 10 })
		want := imp.On("MethodIntVar", by, by, 10, by).Return(nil)

		// --- When ---
		have, err := imp.find("MethodIntVar", []any{1, 2, 10, 4}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("error - when matcher panics", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))

		by := MatchBy(func(have int) bool { panic("boom") })
		imp.On("MethodIntVar", 10, by).Return(nil)

		// --- When ---
		have, err := imp.find("MethodIntVar", []any{10, 4}, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_fail_panicking_matcher.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("returns calls using OnAny first", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.OnAny("Zero").Return("zero")
		want := mck.OnAny("One").Return("one")
		mck.On("One", 2).Return("two") // Will never return this one.

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("returns OnAny call again", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.OnAny("Zero").Return("zero")
		want := mck.OnAny("One").Return("one")
		mck.On("One", 2).Return("two") // Will never return this one.
		want.haveCalls++

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("respects wantCount for calls using OnAny", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.OnAny("Zero").Return("zero")
		call := mck.OnAny("One").Return("one").Once()
		want := mck.On("One", 2).Return("two")
		call.haveCalls++

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, want, have)
	})

	t.Run("proxy calls", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		ptr := &testcases.TPtr{Val: "b"}
		mck := NewMock(tspy)
		mck.On("Error").Return(nil)
		call := mck.Proxy(ptr.Wrap)

		// --- When ---
		have, err := mck.find("Wrap", nil, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, call, have)
	})

	t.Run("proxy calls with arguments", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		ptr := &testcases.TPtr{Val: "b"}
		mck := NewMock(tspy)
		mck.On("Error").Return(nil)
		mck.Proxy(ptr.Wrap)
		call := mck.Proxy(ptr.Wrap).With("C", "D")

		// --- When ---
		have, err := mck.find("Wrap", []any{"C", "D"}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, call, have)
	})

	t.Run("error - when too many calls to proxy method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Error").Return(nil)
		ptr := &testcases.TPtr{Val: "b"}
		mck.Proxy(ptr.Wrap).Once()
		_ = mck.Call("Wrap", "a", "c")

		// --- When ---
		have, err := mck.find("Wrap", nil, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrTooManyCalls, err)
		wMsg := goldy.Open(t, "testdata/find_fail_too_many_calls.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("error - with stack", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Return("zero")
		stack := []string{"a", "b", "c"}

		// --- When ---
		have, err := mck.find("Zero", nil, stack)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
		wMsg := goldy.Open(t, "testdata/find_fail_with_stack.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})
}

func Test_Mock_closest(t *testing.T) {
	t.Run("best", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", true)
		want := imp.On("MethodIntVar", 42, 44).Return(nil)
		imp.On("MethodIntVar", 7, 42, 44).Return(nil)

		// --- When ---
		hCall, hDiff := imp.closest("MethodIntVar", 42, 44)

		// --- Then ---
		assert.Same(t, want, hCall)
		wDiff := []string{
			"0: PASS: (int=42) == (int=42)",
			"1: PASS: (int=44) == (int=44)",
		}
		assert.Equal(t, wDiff, hDiff)
	})

	t.Run("not in expectations", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", true)
		imp.On("MethodBool", false)

		// --- When ---
		hCall, hDiff := imp.closest("MethodBoolean", true)

		// --- Then ---
		assert.Nil(t, hCall)
		assert.Nil(t, hDiff)
	})
}

func Test_Mock_Failed(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		mck := &Mock{t: tspy, failed: true}

		// --- When ---
		have := mck.Failed()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("false", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		mck := &Mock{t: tspy, failed: false}

		// --- When ---
		have := mck.Failed()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Mock_Unset(t *testing.T) {
	t.Run("error - nil call", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		tspy.ExpectLogContain("[mock] unsetting non-existing method")
		tspy.Close()

		mck := NewMock(tspy, WithNoStack)

		// --- When ---
		have := mck.Unset(nil)

		// --- Then ---
		assert.Same(t, mck, have)
	})

	t.Run("existing call", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		call0 := mck.On("Zero")
		call1 := mck.On("One")
		call2 := mck.On("Two")

		// --- When ---
		have := mck.Unset(call1)

		// --- Then ---
		assert.Same(t, mck, have)
		assert.Same(t, call0, mck.expected[0])
		assert.Same(t, call2, mck.expected[1])
		assert.Len(t, 2, mck.expected)
	})

	t.Run("error - not existing call", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/unset_method_not_found.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy)
		call := mck.On("Zero")
		mck.Unset(call)

		// --- When ---
		have := mck.Unset(call)

		// --- Then ---
		assert.Same(t, mck, have)
	})
}

func Test_Mock_AssertExpectations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", true).Once()
		imp.On("MethodBool", false).Once()
		imp.On("MethodBool", true).Once()
		imp.MethodBool(true)
		imp.MethodBool(false)
		imp.MethodBool(true)

		// --- When ---
		have := imp.AssertExpectations()

		// --- Then ---
		assert.True(t, have)
		assert.False(t, imp.failed)
	})

	t.Run("optional does not have to be called", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Optional()

		// --- When ---
		have := mck.AssertExpectations()

		// --- Then ---
		assert.True(t, have)
		assert.False(t, mck.failed)
	})

	t.Run("error - when one call not satisfied", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/missing_calls_one.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy, WithNoStack))
		imp.On("MethodBool", true).Once()
		imp.On("MethodBool", false).Once()
		imp.On("MethodBool", true).Once()
		imp.MethodBool(true)
		imp.MethodBool(false)

		// --- When ---
		have := imp.AssertExpectations()

		// --- Then ---
		assert.False(t, have)
		assert.True(t, imp.failed)
	})

	t.Run("error - when multiple calls not satisfied", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/missing_calls_multiple.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy, WithNoStack))
		imp.On("MethodBool", true).Once()
		imp.On("MethodBool", false).Once()
		imp.On("MethodInts", 1, 2, 3).Return(6, nil).Times(3)
		imp.On("MethodIntVar", 4, 5).Return(nil).Times(2)
		imp.MethodBool(true)
		_, _ = imp.MethodInts(1, 2, 3)

		// --- When ---
		have := imp.AssertExpectations()

		// --- Then ---
		assert.False(t, have)
		assert.True(t, imp.failed)
	})

	t.Run("error - matcher panics", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/mock_matcher_panics.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy, WithNoStack))
		am := MatchBy(func(_ bool) bool { var i int; return 1/i == 0 })
		imp.On("MethodBool", am).Return()

		// AssertExpectations must not deadlock after the matcher panicked.
		defer func() {
			ch := make(chan struct{})
			go func() {
				_ = imp.AssertExpectations()
				close(ch)
			}()
			assert.ChannelWillClose(t, time.Second, ch)
		}()

		// --- When ---
		have := assert.PanicMsg(t, func() { imp.MethodBool(false) })

		// --- Then ---
		assert.Equal(t, tester.FailNowMsg, *have)
		assert.True(t, imp.failed)
	})
}

func Test_Mock_AssertCallCount(t *testing.T) {
	t.Run("not called", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodIntVar", 42).Return(nil)

		// --- When ---
		have := imp.AssertCallCount("MethodIntVar", 0)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodIntVar", 42).Return(nil)
		assert.Nil(t, imp.MethodIntVar(42))
		assert.Nil(t, imp.MethodIntVar(42))

		// --- When ---
		have := imp.AssertCallCount("MethodIntVar", 2)

		// --- Then ---
		assert.True(t, have)
		assert.False(t, imp.failed)
	})

	t.Run("error - when method called too few times", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/assert_call_cnt_too_few.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", 42).Return(nil)

		// --- When ---
		have := imp.AssertCallCount("MethodBool", 1)

		// --- Then ---
		assert.False(t, have)
		assert.True(t, imp.failed)
	})

	t.Run("error - when method called too many times", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/assert_call_cnt_too_many.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		imp := NewExampleImpl(NewMock(tspy))
		imp.On("MethodBool", true)
		imp.MethodBool(true)
		imp.MethodBool(true)

		// --- When ---
		have := imp.AssertCallCount("MethodBool", 1)

		// --- Then ---
		assert.False(t, have)
		assert.True(t, imp.failed)
	})
}

// Benchmark_Mock_RealisticMatching exercises the core matching hot path
// (find + argument Diff + mix of exact/Any/MatchBy).
func Benchmark_Mock_RealisticMatching(b *testing.B) {
	dummyT := &testing.T{}
	mck := &Mock{t: dummyT}

	// Set up a realistic mix of expectations, bypassing NewMock for speed.
	mck.On("GetUser", 42).Return("alice", nil)
	mck.On("GetUser", Any).Return("bob", nil)
	notNil := MatchBy(func(x any) bool { return x != nil })
	mck.On("Process", notNil).Return(true)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mck.Call("GetUser", 42)
		_ = mck.Call("GetUser", 99)
		_ = mck.Call("Process", "data")
	}
}

// Benchmark_Mock_AssertExpectations exercises the full expectation
// verification path on a mock with a moderate number of calls.
func Benchmark_Mock_AssertExpectations(b *testing.B) {
	dummyT := &testing.T{}
	mck := &Mock{t: dummyT}

	for i := 0; i < 20; i++ {
		mck.On("Op", i).Return(i * 2).Once()
	}

	// Pre-satisfy them
	for i := 0; i < 20; i++ {
		_ = mck.Call("Op", i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mck.AssertExpectations()
	}
}
