// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
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
		assert.Equal(t, 0, len(mck.data))
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

func Test_Mock_SetData_GetData(t *testing.T) {
	t.Run("does not return nil", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		have := mck.GetData()

		// --- Then ---
		assert.NotNil(t, have)
	})

	t.Run("map key set", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		data := mck.GetData()

		// --- When ---
		data["key"] = 123

		// --- Then ---
		have := mck.GetData()
		assert.Equal(t, 123, have["key"])
	})

	t.Run("set and get", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		data := map[string]any{
			"k0": "v0",
			"k1": 123,
			"k2": true,
		}
		mck := NewMock(tspy)

		// --- When ---
		mck.SetData(data)

		// --- Then ---
		assert.Equal(t, data, mck.GetData())
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
		call0 := mck.On("Zero", "str", 42, true)

		// --- Then ---
		assert.Len(t, 1, mck.expected)
		assert.Same(t, mck.expected[0], call0)
		assert.Same(t, mck, call0.parent)
		assert.Equal(t, "Zero", call0.Method)
		assert.Equal(t, Arguments{"str", 42, true}, call0.args)
		assert.False(t, call0.argsAny)
		assert.Len(t, 0, call0.returns)
		assert.Len(t, 3, call0.cStack.Stack)
		assert.Contain(t, "mock_test.go", call0.cStack.Stack[2])
		assert.Equal(t, 0, call0.wantCalls)
		assert.Equal(t, 0, call0.haveCalls)
		assert.False(t, call0.optional)
		assert.Nil(t, call0.until)
		assert.Equal(t, time.Duration(0), call0.after)
		assert.Nil(t, call0.alter)
		assert.Nil(t, call0.panic)
		assert.Len(t, 0, call0.requires)
	})

	t.Run("panics when functions are in args expectations", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		fn := func() { mck.On("Zero", func(string) error { return nil }) }

		// --- Then ---
		msg := assert.PanicMsg(t, fn)
		assert.Equal(t, *msg, "cannot use functions in argument expectations")
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
		call := mck.On("Zero", 0).Return("zero").End().
			On("One", 1).Return("one")

		// --- Then ---
		assert.Len(t, 2, mck.expected)
		assert.Same(t, mck.expected[1], call)
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

		mck := NewExampleImpl(NewMock(tspy))

		// --- When ---
		mby0 := MatchBy(func(slice []bool) bool { return slice == nil })
		mck.On("MethodBoolS", mby0).Return(errors.New("fixture1"))

		// --- Then ---
		assert.Equal(t, mck.MethodBoolS(nil).Error(), "fixture1")
	})

	t.Run("error with variadic method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_unexpected_variadic.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy, WithNoStack))

		// --- When ---
		call := mck.On("MethodIntVar", 1, 2, 3).Return(nil)

		// --- Then ---
		assert.Len(t, 1, mck.expected)
		assert.Same(t, call, mck.expected[0])
		assert.Len(t, 3, call.args)
		assert.Equal(t, Arguments{1, 2, 3}, call.args)

		assert.NoPanic(t, func() { _ = mck.MethodIntVar(1, 2, 3) })
		assert.Panic(t, func() { _ = mck.MethodIntVar(1, 2) })
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
		call0 := mck.OnAny("Zero")

		// --- Then ---
		assert.Len(t, 1, mck.expected)
		assert.Same(t, mck.expected[0], call0)
		assert.Same(t, mck, call0.parent)
		assert.Equal(t, "Zero", call0.Method)
		assert.Empty(t, call0.args)
		assert.True(t, call0.argsAny)
		assert.Len(t, 0, call0.returns)
		assert.Len(t, 3, call0.cStack.Stack)
		assert.Contain(t, "mock_test.go", call0.cStack.Stack[2])
		assert.Equal(t, 0, call0.wantCalls)
		assert.Equal(t, 0, call0.haveCalls)
		assert.False(t, call0.optional)
		assert.Nil(t, call0.until)
		assert.Equal(t, time.Duration(0), call0.after)
		assert.Nil(t, call0.alter)
		assert.Nil(t, call0.panic)
		assert.Len(t, 0, call0.requires)
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
		call := mck.OnAny("Zero").Return("zero").End().
			OnAny("One").Return("one")

		// --- Then ---
		assert.Len(t, 2, mck.expected)
		assert.Same(t, mck.expected[1], call)
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

		mck := NewExampleImpl(NewMock(tspy))

		// --- When ---
		call := mck.OnAny("MethodIntVar").Return(nil)

		// --- Then ---
		assert.Len(t, 1, mck.expected)
		assert.Same(t, call, mck.expected[0])
		assert.Len(t, 0, call.args)

		assert.NoPanic(t, func() { _ = mck.MethodIntVar(1, 2, 3) })
		assert.NoPanic(t, func() { _ = mck.MethodIntVar(1, 2) })
	})
}

func Test_Mock_Proxy(t *testing.T) {
	t.Run("with default name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		ptr := &types.TPtr{Val: "abc"}

		// --- When ---
		call := mck.Proxy(ptr.AAA)

		// --- Then ---
		assert.Same(t, mck, call.parent)
		assert.Nil(t, call.args)
		assert.Equal(t, "AAA", call.Method)
	})

	t.Run("with custom name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		ptr := &types.TPtr{Val: "abc"}

		// --- When ---
		call := mck.Proxy(ptr.AAA, "MyName")

		// --- Then ---
		assert.Same(t, mck, call.parent)
		assert.Nil(t, call.args)
		assert.Equal(t, "MyName", call.Method)
	})

	t.Run("panics if not method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		msg := assert.PanicMsg(t, func() { mck.Proxy(123) })

		// --- Then ---
		assert.Equal(t, "Proxy requires a valid not nil method", *msg)
	})

	t.Run("panics if invalid", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		type T func()
		var fn T
		mck := NewMock(tspy)

		// --- When ---
		msg := assert.PanicMsg(t, func() { mck.Proxy(fn) })

		// --- Then ---
		assert.Equal(t, "Proxy requires a valid not nil method", *msg)
	})
}

func Test_Mock_Called(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodInts", 1, 2, 3).Return(4)

		// --- When ---
		got, err := mck.MethodInts(1, 2, 3)

		// --- Then ---
		assert.Equal(t, "whoops", err.Error())
		assert.Equal(t, 4, got)

		assert.Len(t, 1, mck.calls)
		assert.Equal(t, "MethodInts", mck.calls[0].Method)
		assert.Len(t, 5, mck.calls[0].Stack)
		assert.Contain(t, "mock_test.go", mck.calls[0].Stack[4])
	})

	t.Run("calling unexpected method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/call_unexpected.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy, WithNoStack))

		// --- When ---
		assert.Panic(t, func() { mck.Called(1, 2, 3) })
	})
}

func Test_called(t *testing.T) {
	t.Run("self", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		method := mck.called(0)

		// --- Then ---
		assert.Equal(t, "called", method)
	})

	t.Run("sub test", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- When ---
		method := mck.called(1)

		// --- Then ---
		assert.Equal(t, "<anonymous>", method)
	})

	t.Run("panics for invalid skip value", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)

		// --- Then ---
		msg := assert.PanicMsg(t, func() { mck.called(100) })
		assert.Equal(t, *msg, "could not get the caller information")
	})
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

	t.Run("error when method called too many times", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/mock_too_many_calls.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0, 1).Once()
		mck.Call("Zero", 0, 1)

		// --- When ---
		assert.Panic(t, func() { mck.Call("Zero", 0, 1) })
		assert.True(t, mck.failed)
	})

	t.Run("error when existing method called with different arguments", func(t *testing.T) {
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
		assert.Panic(t, func() { mck.Call("Zero", 2, 3) })
		assert.True(t, mck.failed)
	})

	t.Run("error call to not expected method name", func(t *testing.T) {
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
		assert.Panic(t, func() { mck.Call("One", 1) })
		assert.True(t, mck.failed)
	})

	t.Run("error when method called before required deps are met", func(t *testing.T) {
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
		assert.Panic(t, func() { mck.Call("Two", 2) })
		assert.True(t, mck.failed)
	})

	t.Run("error when not all requirements are met", func(t *testing.T) {
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
		assert.Panic(t, func() { mck.Call("Two", 2) })
		assert.True(t, mck.failed)
	})

	t.Run("error requirement values did not match", func(t *testing.T) {
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
		assert.Panic(t, func() { mck.Call("Two", 2) })
		assert.True(t, mck.failed)
	})

	t.Run("wait for until", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		ch := time.After(500 * time.Millisecond)

		mck := NewMock(tspy)
		mck.On("MethodBool", true).Return("done").Until(ch)

		// --- When ---
		var got Arguments
		done := make(chan struct{})
		start := time.Now()
		go func() {
			got = mck.Call("MethodBool", true)
			close(done)
		}()

		// --- Then ---
		select {
		case <-done:
			// Test it took roughly more than 500ms to run the Call.
			since := time.Since(start)
			if !assert.True(t, since > 490*time.Millisecond) {
				t.Logf("since: %s\n", since.String())
				return
			}

		case <-time.After(time.Second):
			t.Error("Call does not work properly")
			return
		}
		assert.Equal(t, "done", got[0])
		assert.False(t, mck.failed)
	})

	t.Run("wait for sleep", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("MethodBool", true).Return("done").After(500 * time.Millisecond)

		// --- When ---
		var got Arguments
		done := make(chan struct{})
		start := time.Now()
		go func() {
			got = mck.Call("MethodBool", true)
			close(done)
		}()

		// --- Then ---
		select {
		case <-done:
			// Test it took more than 500ms to run the Call.
			since := time.Since(start)
			if !assert.True(t, since > 500*time.Millisecond) {
				t.Logf("since: %s\n", since.String())
				return
			}

		case <-time.After(time.Second):
			t.Error("Call does not work properly")
			return
		}
		assert.Equal(t, "done", got[0])
		assert.False(t, mck.failed)
	})

	t.Run("mocked method should panic", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Zero", 0).Panic("zero panics")

		// --- Then ---
		msg := assert.PanicMsg(t, func() { mck.Call("Zero", 0) })
		assert.Equal(t, *msg, "zero panics")
		assert.False(t, mck.failed)
	})

	t.Run("run is ran in order", func(t *testing.T) {
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
		call0 := mck.On("ConcurrencyTestMethod", 1).Return(1)

		wg := sync.WaitGroup{}
		wg.Add(2)

		// --- When ---
		const cnt = 1000
		go func() {
			// Edit the call changing its return arguments.
			for i := 0; i < cnt; i++ {
				mck.Call("ConcurrencyTestMethod", 1)
			}
			wg.Done()
		}()

		go func() {
			// Continuously call the mocked method.
			for i := 0; i < cnt; i++ {
				mck.Call("ConcurrencyTestMethod", 1)
			}
			wg.Done()
		}()

		// --- Then ---
		wg.Wait()
		assert.Equal(t, cnt*2, call0.haveCalls)
		assert.False(t, mck.failed)
	})

	t.Run("call proxy method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewMock(tspy)
		ptr := &types.TPtr{Val: "b"}
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

		mck := NewMock(tspy)
		ptr := &types.TPtr{Val: "b"}
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
	t.Run("callable", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodIntVar", 42, 44).Return(nil)

		// --- When ---
		err := mck.Callable("MethodIntVar", 42, 44)

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("error when not existing method name given", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodIntVar", 42, 44).Return(nil)

		// --- When ---
		err := mck.Callable("NotExisting")

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
	})

	t.Run("error when not matching arguments given", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodIntVar", 42, 44).Return(nil)

		// --- When ---
		err := mck.Callable("MethodIntVar", 7, 44)

		// --- Then ---
		assert.ErrorIs(t, ErrNotFound, err)
	})

	t.Run("error when called again", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodIntVar", 42, 44).Once().Return(nil)
		assert.Nil(t, mck.MethodIntVar(42, 44))

		// --- When ---
		err := mck.Callable("MethodIntVar", 42, 44)

		// --- Then ---
		assert.ErrorIs(t, ErrTooManyCalls, err)
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
		exp := mck.On("One", 2).Return("two")

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, exp, have)
	})

	t.Run("error matching name not matching arg count", func(t *testing.T) {
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

	t.Run("error method name call not found", func(t *testing.T) {
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

	t.Run("error method name call not found with arguments", func(t *testing.T) {
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

	t.Run("error matching name not matching arg type", func(t *testing.T) {
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
		mck.On("One", 1).Times(-1).Return("one")
		exp := mck.On("One", 1).Return("two")

		// --- When ---
		have, err := mck.find("One", []any{1}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, exp, have)
	})

	t.Run("method with variadic args", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", true)
		exp := mck.On("MethodIntVar", 1, 2, 3).Return(nil)

		// --- When ---
		have, err := mck.find("MethodIntVar", []any{1, 2, 3}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, exp, have)
	})

	t.Run("error method with variadic args not found", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", true)
		mck.On("MethodIntVar", 1, 2, 3, 4).Return(nil)

		// --- When ---
		have, err := mck.find("MethodIntVar", []any{1, 2, 3}, nil)

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

		mck := NewExampleImpl(NewMock(tspy))

		by := MatchBy(func(have int) bool { return have < 10 })
		exp := mck.On("MethodIntVar", by, by, 10, by).Return(nil)

		// --- When ---
		have, err := mck.find("MethodIntVar", []any{1, 2, 10, 4}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, exp, have)
	})

	t.Run("error when matcher panics", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))

		by := MatchBy(func(have int) bool { panic("boom") })
		mck.On("MethodIntVar", 10, by).Return(nil)

		// --- When ---
		have, err := mck.find("MethodIntVar", []any{10, 4}, nil)

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
		exp := mck.OnAny("One").Return("one")
		mck.On("One", 2).Return("two") // Will never return this one.

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, exp, have)
		exp.haveCalls++

		// --- When ---
		have, err = mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, exp, have)
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
		expAny := mck.OnAny("One").Return("one").Once()
		expOn := mck.On("One", 2).Return("two")

		// --- When ---
		have, err := mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, expAny, have)
		expAny.haveCalls++

		// --- When ---
		have, err = mck.find("One", []any{2}, nil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Same(t, expOn, have)
	})

	t.Run("proxy calls", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		ptr := &types.TPtr{Val: "b"}
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

		ptr := &types.TPtr{Val: "b"}
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

	t.Run("error when too many calls to proxy method", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewMock(tspy)
		mck.On("Error").Return(nil)
		ptr := &types.TPtr{Val: "b"}
		mck.Proxy(ptr.Wrap).Once()
		mck.Call("Wrap", "a", "c")

		// --- When ---
		have, err := mck.find("Wrap", nil, nil)

		// --- Then ---
		assert.ErrorIs(t, ErrTooManyCalls, err)
		wMsg := goldy.Open(t, "testdata/find_fail_too_many_calls.gld")
		assert.ErrorEqual(t, wMsg.String(), err)
		assert.Nil(t, have)
	})

	t.Run("error with stack", func(t *testing.T) {
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

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", true)
		exp := mck.On("MethodIntVar", 42, 44).Return(nil)
		mck.On("MethodIntVar", 7, 42, 44).Return(nil)

		// --- When ---
		got, diff := mck.closest("MethodIntVar", 42, 44)

		// --- Then ---
		assert.Same(t, exp, got)
		wantDiff := []string{
			"0: PASS: (int=42) == (int=42)",
			"1: PASS: (int=44) == (int=44)",
		}
		assert.Equal(t, wantDiff, diff)
	})

	t.Run("not in expectations", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", true)
		mck.On("MethodBool", false)

		// --- When ---
		got, have := mck.closest("MethodBoolean", true)

		// --- Then ---
		assert.Nil(t, got)
		assert.Nil(t, have)
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

	t.Run("error not existing call", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/unset_method_not_found.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewMock(tspy)
		call := mck.On("Zero")
		mck.Unset(call) // OK

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

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", true).Once()
		mck.On("MethodBool", false).Once()
		mck.On("MethodBool", true).Once()

		// --- When ---
		mck.MethodBool(true)
		mck.MethodBool(false)
		mck.MethodBool(true)
		have := mck.AssertExpectations()

		// --- Then ---
		assert.True(t, have)
		assert.False(t, mck.failed)
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

	t.Run("error when one call not satisfied", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/missing_calls_one.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy, WithNoStack))
		mck.On("MethodBool", true).Once()
		mck.On("MethodBool", false).Once()
		mck.On("MethodBool", true).Once()

		// --- When ---
		mck.MethodBool(true)
		mck.MethodBool(false)
		have := mck.AssertExpectations()

		// --- Then ---
		assert.False(t, have)
		assert.True(t, mck.failed)
	})

	t.Run("error when multiple calls not satisfied", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/missing_calls_multiple.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy, WithNoStack))
		mck.On("MethodBool", true).Once()
		mck.On("MethodBool", false).Once()
		mck.On("MethodInts", 1, 2, 3).Return(6, nil).Times(3)
		mck.On("MethodIntVar", 4, 5).Return(nil).Times(2)

		// --- When ---
		mck.MethodBool(true)
		_, _ = mck.MethodInts(1, 2, 3)
		have := mck.AssertExpectations()

		// --- Then ---
		assert.False(t, have)
		assert.True(t, mck.failed)
	})

	t.Run("matcher panics", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/mock_matcher_panics.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy, WithNoStack))
		mby0 := MatchBy(func(_ bool) bool { var i int; return 1/i == 0 })
		mck.On("MethodBool", mby0).Return()

		// --- When ---
		defer func() {
			ch := make(chan struct{})
			go func() {
				mck.AssertExpectations()
				close(ch)
			}()

			select {
			case <-ch:
			case <-time.After(time.Second):
				t.Error("AssertExpectations() deadlocked")
			}
		}()

		// --- Then ---
		assert.Panic(t, func() { mck.MethodBool(false) })
		assert.True(t, mck.failed)
	})
}

func Test_Mock_AssertCallCount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodIntVar", 42).Return(nil)

		// --- Then ---
		assert.True(t, mck.AssertCallCount("MethodIntVar", 0))
		assert.Nil(t, mck.MethodIntVar(42))
		assert.True(t, mck.AssertCallCount("MethodIntVar", 1))
		assert.Nil(t, mck.MethodIntVar(42))
		assert.True(t, mck.AssertCallCount("MethodIntVar", 2))
		assert.False(t, mck.failed)
	})

	t.Run("error when method called too few times", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/assert_call_cnt_too_few.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", true)
		mck.MethodBool(true)
		mck.MethodBool(true)

		// --- When ---
		have := mck.AssertCallCount("MethodBool", 1)

		// --- Then ---
		assert.False(t, have)
		assert.True(t, mck.failed)
	})

	t.Run("error when method called too many times", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(1)
		tspy.ExpectFail()
		wMsg := goldy.Open(t, "testdata/assert_call_cnt_too_many.gld")
		tspy.ExpectLogEqual(wMsg.String())
		tspy.Close()

		mck := NewExampleImpl(NewMock(tspy))
		mck.On("MethodBool", 42).Return(nil)

		// --- When ---
		have := mck.AssertCallCount("MethodBool", 1)

		// --- Then ---
		assert.False(t, have)
		assert.True(t, mck.failed)
	})
}
