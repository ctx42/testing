// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

// Compile-time checks that the real testing types satisfy our subset.
// These will fail to compile if the T interface drifts from what
// *testing.T / *testing.B / *testing.F actually provide.
var (
	_ T = (*testing.T)(nil)
	_ T = (*testing.B)(nil)
	_ T = (*testing.F)(nil)
)

func Test_find_match_tabular(t *testing.T) {
	tt := []struct {
		testN string

		strategy Strategy
		want     string
		have     string
		match    bool
	}{
		{"equal success", Equal, "abc def ghi", "abc def ghi", true},
		{"error - equal", Equal, "abc def ghi", "def", false},
		{"contains success", Contains, "def", "abc def ghi", true},
		{"error - contains", Contains, "def", "abc xyz ghi", false},
		{"not contain success", NotContains, "jkl", "abc def ghi", true},
		{"error - not contains", NotContains, "def", "abc def ghi", false},
		{"regexp success", Regexp, "[a-z ]+", "abc def ghi", true},
		{"error - regexp", Regexp, "abc def ghi", "xyz", false},
		{
			"unknown as equal success",
			"unknown",
			"abc def ghi",
			"abc def ghi",
			true,
		},
		{"error - unknown as equal", "unknown", "abc def ghi", "def", false},
		{"equal fold success", EqualFold, "ABC DEF GHI", "abc def ghi", true},
		{"error - equal fold", EqualFold, "ABC DEF GHI", "abc xyz ghi", false},
		{"contains fold success", ContainsFold, "DEF", "abc def ghi", true},
		{"error - contains fold", ContainsFold, "DEF", "abc xyz ghi", false},
		{
			"not contains fold success",
			NotContainsFold,
			"JKL",
			"abc def ghi",
			true,
		},
		{
			"error - not contains fold",
			NotContainsFold,
			"DEF",
			"abc def ghi",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			fnd := &find{strategy: tc.strategy, want: tc.want}

			// --- When ---
			have := fnd.match(tc.have)

			// --- Then ---
			affirm.Equal(t, tc.match, have)
		})
	}
}

func Test_New(t *testing.T) {
	t.Run("initial values", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		spy := New(ti)

		// --- Then ---
		affirm.Equal(t, false, spy.helperCntSet)
		affirm.Equal(t, -1, spy.wantHelperCnt)
		affirm.Equal(t, 0, spy.haveHelperCnt)
		affirm.Equal(t, 0, spy.wantNamesCnt)
		affirm.Equal(t, 0, spy.haveNamesCnt)
		affirm.Equal(t, 0, spy.wantTempDirCnt)
		affirm.Equal(t, 0, len(spy.haveTempDirs))
		affirm.Equal(t, 0, len(spy.wantEnv))
		affirm.Equal(t, 0, len(spy.haveEnv))
		affirm.Equal(t, false, spy.closed)
		affirm.Equal(t, false, spy.finished)
		affirm.Equal(t, false, spy.wantFailed)
		affirm.Equal(t, false, spy.wantError)
		affirm.Equal(t, false, spy.haveError)
		affirm.Equal(t, false, spy.wantFatal)
		affirm.Equal(t, false, spy.haveFatal)
		affirm.Equal(t, false, spy.wantSkipped)
		affirm.Equal(t, false, spy.haveSkipped)
		affirm.Equal(t, false, spy.panicked)
		affirm.Equal(t, 0, spy.haveCleanupsCnt)
		affirm.Equal(t, 0, len(spy.haveCleanups))
		affirm.Equal(t, 0, len(spy.savedMsgs))
		affirm.Equal(t, 0, len(spy.haveLogMsgs))
		affirm.Equal(t, false, spy.ignoreLog)
		affirm.Equal(t, false, spy.runningFinish)
		affirm.Equal(t, 0, spy.wantCleanupsCnt)
		affirm.Equal(t, 0, len(spy.wantLogMsgs))
		affirm.Equal(t, false, spy.asserted)
		affirm.Nil(t, spy.ctx)
		affirm.Nil(t, spy.cancelCtx)
		affirm.Equal(t, true, spy.tt == ti)
	})

	t.Run("set expected Helper calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		spy := New(ti, 1)

		// --- Then ---
		affirm.Equal(t, true, spy.helperCntSet)
		affirm.Equal(t, 1, spy.wantHelperCnt)
		affirm.Equal(t, 0, spy.haveHelperCnt)
	})

	t.Run("error - invalid Helper calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		have := affirm.Panic(t, func() { New(ti, -2) })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "ExpectHelpers cnt must be greater or equal to minus one"
		affirm.Equal(t, want, *have)
	})

	t.Run("error - more than one Helper calls value", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		have := affirm.Panic(t, func() { New(ti, 1, 2) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, "New accepts at most one expectHelpers value", *have)
	})
}

func Test_New_auto_finish(t *testing.T) {
	// --- Given ---
	var runs []int
	var spy *Spy

	// --- When ---
	t.Run("inner test", func(t *testing.T) {
		// --- Given ---
		spy = New(t, 0)
		spy.ExpectCleanups(2)
		spy.Close()

		// --- When ---
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { runs = append(runs, 1) })
	})

	// --- Then ---
	affirm.Equal(t, 2, len(runs))
	affirm.Equal(t, 1, runs[0])
	affirm.Equal(t, 0, runs[1])
	affirm.Equal(t, true, spy.tt == nil)
}

func Test_Spy_ExpectCleanups(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectCleanups(2)

		// --- Then ---
		affirm.Equal(t, 2, spy.wantCleanupsCnt)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectCleanups(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectCleanups(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - negative count panics", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectCleanups(-1) })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "ExpectCleanups cnt must be greater or equal to zero"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Cleanup(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { runs = append(runs, 1) })

		// --- Then ---
		affirm.Equal(t, 2, spy.haveCleanupsCnt)
		affirm.Equal(t, 2, len(spy.haveCleanups))
		affirm.Equal(t, 0, len(runs))

		spy.haveCleanups[0]()
		affirm.Equal(t, 1, len(runs))
		affirm.Equal(t, 0, runs[0])

		spy.haveCleanups[1]()
		affirm.Equal(t, 2, len(runs))
		affirm.Equal(t, 0, runs[0])
		affirm.Equal(t, 1, runs[1])
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Cleanup(func() {}) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Cleanup(func() {}) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectError(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectError()

		// --- Then ---
		affirm.Equal(t, true, spy.wantError)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectError() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectError() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called after ExpectFail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectError() })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "cannot use ExpectError and ExpectFail at the same time"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Error(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Error("msg", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, true, spy.haveError)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Error("msg", 0)
		spy.Error("msg", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
		affirm.Equal(t, true, spy.haveError)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Error("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Error("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Errorf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Errorf("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, true, spy.haveError)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Errorf("msg %d", 0)
		spy.Errorf("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
		affirm.Equal(t, true, spy.haveError)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Errorf("have %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Errorf("have %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectFatal(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectFatal()

		// --- Then ---
		affirm.Equal(t, true, spy.wantFatal)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFatal() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFatal() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called after ExpectFail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFatal() })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "cannot use ExpectFatal and ExpectFail at the same time"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Fatal(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Fatal("msg", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, FailNowMsg, *have)
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, true, spy.haveFatal)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		msg0 := affirm.Panic(t, func() { spy.Fatal("msg", 0) })
		msg1 := affirm.Panic(t, func() { spy.Fatal("msg", 1) })

		// --- Then ---
		affirm.NotNil(t, msg0)
		affirm.Equal(t, FailNowMsg, *msg0)
		affirm.NotNil(t, msg1)
		affirm.Equal(t, FailNowMsg, *msg1)

		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
		affirm.Equal(t, true, spy.haveFatal)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Fatal("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Fatal("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Fatalf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Fatalf("msg %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, FailNowMsg, *have)
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, true, spy.haveFatal)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		msg0 := affirm.Panic(t, func() { spy.Fatalf("msg %d", 0) })
		msg1 := affirm.Panic(t, func() { spy.Fatalf("msg %d", 1) })

		// --- Then ---
		affirm.NotNil(t, msg0)
		affirm.Equal(t, FailNowMsg, *msg0)
		affirm.NotNil(t, msg1)
		affirm.Equal(t, FailNowMsg, *msg1)

		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
		affirm.Equal(t, true, spy.haveFatal)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Fatalf("have %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Fatalf("have %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_FailNow(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.FailNow() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, FailNowMsg, *have)
		affirm.Equal(t, false, spy.panicked)
		affirm.Equal(t, true, spy.haveFatal)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.FailNow() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.FailNow() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Failed(t *testing.T) {
	t.Run("not failed", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("have error", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Error("msg", 0)

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("have fatal", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		msg := affirm.Panic(t, func() { spy.Fatal("msg", 0) })
		affirm.NotNil(t, msg)
		affirm.Equal(t, FailNowMsg, *msg)

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("have fatal and error", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		msg := affirm.Panic(t, func() { spy.Fatal("msg", 0) })
		affirm.NotNil(t, msg)
		affirm.Equal(t, FailNowMsg, *msg)
		spy.Error("msg", 1)

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Failed() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Spy_ExpectHelpers(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- When ---
		have := spy.ExpectHelpers(2)

		// --- Then ---
		affirm.Equal(t, 2, spy.wantHelperCnt)
		affirm.Equal(t, true, spy.helperCntSet)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("set at least once", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- When ---
		have := spy.ExpectHelpers(-1)

		// --- Then ---
		affirm.Equal(t, -1, spy.wantHelperCnt)
		affirm.Equal(t, true, spy.helperCntSet)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - argument less than minus one", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectHelpers(-2) })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "ExpectHelpers cnt must be greater or equal to minus one"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - set twice", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)
		spy.ExpectHelpers(1)

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectHelpers(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, "ExpectHelpers may be called only once", *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectHelpers(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectHelpers(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Helper(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Helper()

		// --- Then ---
		affirm.Equal(t, 1, spy.haveHelperCnt)
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Helper()
		spy.Helper()

		// --- Then ---
		affirm.Equal(t, 2, spy.haveHelperCnt)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Helper() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Helper() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectSetenv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectSetenv("k0", "v0")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantEnv))
		affirm.Equal(t, "v0", spy.wantEnv["k0"])
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectSetenv("k0", "v0") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectSetenv("k0", "v0") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Setenv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "k0")

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Setenv("k0", "v0")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveEnv))
		affirm.Equal(t, "v0", spy.haveEnv["k0"])
	})

	t.Run("set multiple", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "k0", "k1")

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Setenv("k0", "v0")
		spy.Setenv("k1", "v1")

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveEnv))
		affirm.Equal(t, "v0", spy.haveEnv["k0"])
		affirm.Equal(t, "v1", spy.haveEnv["k1"])
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Setenv("k0", "v0") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Setenv("k0", "v0") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectSkipped(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectSkipped()

		// --- Then ---
		affirm.Equal(t, true, spy.wantSkipped)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectSkipped() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectSkipped() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Skip(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Skip("msg", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, SkipNowMsg, *have)
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, true, spy.haveSkipped)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		msg0 := affirm.Panic(t, func() { spy.Skip("msg", 0) })
		msg1 := affirm.Panic(t, func() { spy.Skip("msg", 1) })

		// --- Then ---
		affirm.NotNil(t, msg0)
		affirm.Equal(t, SkipNowMsg, *msg0)
		affirm.NotNil(t, msg1)
		affirm.Equal(t, SkipNowMsg, *msg1)
		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
		affirm.Equal(t, true, spy.haveSkipped)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Skip("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Skip("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectTempDir(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectTempDir(1)

		// --- Then ---
		affirm.Equal(t, 1, spy.wantTempDirCnt)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectTempDir(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectTempDir(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - count less than minus one panics", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectTempDir(-2) })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "ExpectTempDir cnt must be greater or equal to minus one"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_GetTempDir(t *testing.T) {
	t.Run("error - negative index", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(1)
		spy.Close()

		// --- When ---
		have := spy.GetTempDir(-1)

		// --- Then ---
		affirm.Equal(t, "", have)
		wMsg := "temp directory with index -1 does not exist"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("get directories", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(2)
		spy.Close()
		want0 := spy.TempDir()
		want1 := spy.TempDir()
		spy.Finish()

		// --- When ---
		have0 := spy.GetTempDir(0)
		have1 := spy.GetTempDir(1)

		// --- Then ---
		affirm.Equal(t, want0, have0)
		affirm.Equal(t, want1, have1)
	})

	t.Run("get directories when count is any", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(-1)
		spy.Close()
		want0 := spy.TempDir()
		want1 := spy.TempDir()
		spy.Finish()

		// --- When ---
		have0 := spy.GetTempDir(0)
		have1 := spy.GetTempDir(1)

		// --- Then ---
		affirm.Equal(t, want0, have0)
		affirm.Equal(t, want1, have1)
	})

	t.Run("error - called before ExpectTempDir", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.GetTempDir(0)

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		msg := "ExpectTempDir method must be called before GetTempDir"
		assertSpyHasMsg(t, spy, 0, msg)
	})

	t.Run("error - not existing index", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(2)
		spy.Close()
		_ = spy.TempDir()
		_ = spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.GetTempDir(2)

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		wMsg := "temp directory with index 2 does not exist"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})
}

func Test_Spy_TempDir(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.Close()

		// --- When ---
		have := spy.TempDir()

		// --- Then ---
		fi, err := os.Lstat(have)
		affirm.Equal(t, true, err == nil)
		affirm.Equal(t, true, fi.IsDir())
		affirm.Equal(t, 1, len(spy.haveTempDirs))

		// Dir is added to cleanups and removed when test finishes.
		spy.Finish()
		_, err = os.Lstat(have)
		affirm.Equal(t, true, errors.Is(err, os.ErrNotExist))
	})

	t.Run("call multi", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.Close()

		// --- When ---
		have0 := spy.TempDir()
		have1 := spy.TempDir()

		// --- Then ---
		affirm.Equal(t, false, have0 == have1)

		fi, err := os.Lstat(have0)
		affirm.Equal(t, true, err == nil)
		affirm.Equal(t, true, fi.IsDir())

		fi, err = os.Lstat(have1)
		affirm.Equal(t, true, err == nil)
		affirm.Equal(t, true, fi.IsDir())

		affirm.Equal(t, 2, len(spy.haveTempDirs))
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.TempDir() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.TempDir() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Context(t *testing.T) {
	t.Run("first call in a cleanup is canceled", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(1)
		spy.Close()

		var have context.Context
		spy.Cleanup(func() { have = spy.Context() })

		// --- When ---
		spy.Finish()

		// --- Then ---
		affirm.Equal(t, true, errors.Is(have.Err(), context.Canceled))
	})

	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := spy.Context()

		// --- Then ---
		affirm.Equal(t, true, have != nil)
		affirm.Equal(t, true, have.Err() == nil)

		spy.Finish()
		affirm.Equal(t, true, errors.Is(have.Err(), context.Canceled))
	})

	t.Run("call multi", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have0 := spy.Context()
		have1 := spy.Context()

		// --- Then ---
		affirm.Equal(t, true, have0 == have1)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Context() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Context() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_IgnoreLogs(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.IgnoreLogs()

		// --- Then ---
		affirm.Equal(t, true, spy.ignoreLog)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.IgnoreLogs() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.IgnoreLogs() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - after an ExpectLog method", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("abc")

		// --- When ---
		have := affirm.Panic(t, func() { spy.IgnoreLogs() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errIgnoreLogsAfterExpectLog, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLog(t *testing.T) {
	t.Run("error - invalid regexp", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLog(Regexp, "[") })

		// --- Then ---
		wMsg := "" +
			"invalid ExpectLog regexp: error parsing regexp: " +
			"missing closing ]: `[`"
		affirm.Equal(t, wMsg, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLog(Contains, "msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("adds empty messages", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLog(Equal, "")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Equal, fnd.strategy)
		affirm.Equal(t, "", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("set with percentages without args", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLog(Contains, "msg %d")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, fnd.strategy)
		affirm.Equal(t, "msg %d", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLog(Contains, "msg %d", 0)
		spy.ExpectLog(Equal, "msg %d", 1)
		spy.ExpectLog(Regexp, "msg [0-9]")
		spy.ExpectLog(NotContains, "msg 6")

		// --- Then ---
		affirm.Equal(t, 4, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)

		fnd = spy.wantLogMsgs[1]
		affirm.Equal(t, Equal, fnd.strategy)
		affirm.Equal(t, "msg 1", fnd.want)

		fnd = spy.wantLogMsgs[2]
		affirm.Equal(t, Regexp, fnd.strategy)
		affirm.Equal(t, "msg [0-9]", fnd.want)

		fnd = spy.wantLogMsgs[3]
		affirm.Equal(t, NotContains, fnd.strategy)
		affirm.Equal(t, "msg 6", fnd.want)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLog(Equal, "have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLog(Equal, "have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLog(Equal, "have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLogEqual(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogEqual("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Equal, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLogEqual("msg %d", 0)
		spy.ExpectLogEqual("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Equal, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)

		fnd = spy.wantLogMsgs[1]
		affirm.Equal(t, Equal, fnd.strategy)
		affirm.Equal(t, "msg 1", fnd.want)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogEqual("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogEqual("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogEqual("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLogContain(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogContain("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLogContain("msg %d", 0)
		spy.ExpectLogContain("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)

		fnd = spy.wantLogMsgs[1]
		affirm.Equal(t, Contains, fnd.strategy)
		affirm.Equal(t, "msg 1", fnd.want)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogContain("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogContain("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogContain("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLogNotContain(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogNotContain("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, NotContains, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLogNotContain("msg %d", 0)
		spy.ExpectLogNotContain("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, NotContains, fnd.strategy)
		affirm.Equal(t, "msg 0", fnd.want)

		fnd = spy.wantLogMsgs[1]
		affirm.Equal(t, NotContains, fnd.strategy)
		affirm.Equal(t, "msg 1", fnd.want)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogNotContain("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogNotContain("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogNotContain("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLogEqualFold(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogEqualFold("MSG %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, EqualFold, fnd.strategy)
		affirm.Equal(t, "MSG 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogEqualFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogEqualFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogEqualFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLogContainFold(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogContainFold("MSG %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, ContainsFold, fnd.strategy)
		affirm.Equal(t, "MSG 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogContainFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogContainFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogContainFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectLogNotContainFold(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogNotContainFold("MSG %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMsgs))
		fnd := spy.wantLogMsgs[0]
		affirm.Equal(t, NotContainsFold, fnd.strategy)
		affirm.Equal(t, "MSG 0", fnd.want)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogNotContainFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogNotContainFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
	})

	t.Run("error - called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectLogNotContainFold("have") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectLogAfterIgnoreLogs, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Log(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Log("msg", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Log("msg", 0)
		spy.Log("msg", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Log("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Log("have", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Logf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Logf("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMsgs))
		affirm.Equal(t, "msg 0", spy.haveLogMsgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMsgs[1])
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Logf("have %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Logf("have %d", 0) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExamineLog(t *testing.T) {
	t.Run("empty log", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExamineLog()

		// --- Then ---
		affirm.Equal(t, "", have)
	})

	t.Run("with log messages", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0).Close()
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 1)

		// --- When ---
		have := spy.ExamineLog()

		// --- Then ---
		affirm.Equal(t, "msg 0\nmsg 1", have)
	})
}

func Test_Spy_ExpectNames(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectNames(2)

		// --- Then ---
		affirm.Equal(t, 2, spy.wantNamesCnt)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectNames(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectNames(1) })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - negative count panics", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectNames(-1) })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "ExpectNames cnt must be greater or equal to zero"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Name(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		_ = spy.Name()

		// --- Then ---
		affirm.Equal(t, 1, spy.haveNamesCnt)
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		_ = spy.Name()
		_ = spy.Name()

		// --- Then ---
		affirm.Equal(t, 2, spy.haveNamesCnt)
	})

	t.Run("returns test name", func(t *testing.T) {
		// --- Given ---
		spy := New(t, 0)
		spy.ExpectNames(1)
		spy.Close()

		// --- When ---
		have := spy.Name()

		// --- Then ---
		affirm.Equal(t, "Test_Spy_Name/returns_test_name", have)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := affirm.Panic(t, func() { spy.Name() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errMockOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Name() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errActionOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_ExpectFail(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectFail()

		// --- Then ---
		affirm.Equal(t, true, spy.wantFailed)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called after ExpectFatal", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectFatal()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFail() })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "cannot use ExpectFail and ExpectFatal at the same time"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called after ExpectError", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectError()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFail() })

		// --- Then ---
		affirm.NotNil(t, have)
		want := "cannot use ExpectFail and ExpectError at the same time"
		affirm.Equal(t, want, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFail() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.ExpectFail() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errExpectOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Close(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		have := spy.Close()

		// --- Then ---
		affirm.Equal(t, true, spy.closed)
		affirm.Equal(t, true, spy == have)
	})

	t.Run("error - called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Close() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errDoubleClose, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Close() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errCloseOnFinished, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_Finish(t *testing.T) {
	t.Run("cleanup calling Fatal", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { spy.Fatal("cleanup") })

		// --- When ---
		spy.Finish()

		// --- Then ---
		affirm.Equal(t, true, spy.finished)
		affirm.Equal(t, true, spy.haveFatal)
		affirm.Equal(t, 1, len(runs))
	})

	t.Run("error - cleanup panic raised last", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { panic("boom") })

		// --- When ---
		have := affirm.Panic(t, func() { spy.Finish() })

		// --- Then ---
		affirm.Equal(t, "boom", *have)
		affirm.Equal(t, true, spy.finished)
		affirm.Equal(t, 1, len(runs))
	})

	t.Run("finish test", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		spy.Finish()

		// --- Then ---
		affirm.Equal(t, true, spy.finished)
	})

	t.Run("runs cleanups", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { runs = append(runs, 1) })

		// --- When ---
		spy.Finish()

		// --- Then ---
		affirm.Equal(t, 2, len(runs))
		affirm.Equal(t, 1, runs[0])
		affirm.Equal(t, 0, runs[1])
	})

	t.Run("error - run twice", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- When ---
		have := affirm.Panic(t, func() { spy.Finish() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errDoubleFinish, *have)
		affirm.Equal(t, true, spy.panicked)
	})
}

func Test_Spy_AssertExpectations(t *testing.T) {
	t.Run("cleanup after explicit assert", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectError()
		spy.Close()
		spy.AssertExpectations()

		// --- When ---
		spy.cleanup()

		// --- Then ---
		assertSpyHasMsgs(t, spy, 1)
	})

	t.Run("error - unexpected env with matching count", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "key1")

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.Close()
		spy.Setenv("key1", "val1")

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		assertSpyHasMsgs(t, spy, 2)
		wMsg := "" +
			"expected HUT not to set environment variable:\n" +
			"\t  have key: \"key1\"\n" +
			"\thave value: \"val1\""
		assertSpyHasMsg(t, spy, 1, wMsg)
	})

	t.Run("finishes and runs cleanups", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(1)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, true, spy.finished)
		affirm.Equal(t, 1, len(runs))
		affirm.Equal(t, 0, runs[0])
	})

	t.Run("Spy with no expectations is success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- When ---
		have := affirm.Panic(t, func() { spy.AssertExpectations() })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, errAssertOnNotClosed, *have)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("error - called on panicked instance", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// Make Spy panic.
		spy.Finish()
		msg := affirm.Panic(t, func() { spy.Finish() })
		affirm.NotNil(t, msg)
		affirm.Equal(t, errDoubleFinish, *msg)
		affirm.Equal(t, true, spy.panicked)

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		assertSpyHasMsg(t, spy, 0, "invalid Spy usage")
	})

	t.Run("ExpectSetenv - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "key0", "key1")

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key0", "val0")
		spy.Setenv("key1", "val1")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectSetenv - not expected but set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "key0")

		spy := New(ti, 0)
		spy.Close()
		spy.Setenv("key0", "val0")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		wMsg := "" +
			"expected HUT not to set environment variable:\n" +
			"\t  have key: \"key0\"\n" +
			"\thave value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - ExpectSetenv - none set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 2)
		wMsg := "" +
			"expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
		wMsg = "" +
			"expected HUT to set environment variable:\n" +
			"\t  want key: \"key1\"\n" +
			"\twant value: \"val1\""
		assertSpyHasMsg(t, spy, 1, wMsg)
	})

	t.Run("error - ExpectSetenv - wrong value set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "key0", "key1")

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key0", "valX")
		spy.Setenv("key1", "valY")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 2)
		wMsg := "" +
			"expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\"\n" +
			"\thave value: \"valX\""
		assertSpyHasMsg(t, spy, 0, wMsg)
		wMsg = "" +
			"expected HUT to set environment variable:\n" +
			"\t  want key: \"key1\"\n" +
			"\twant value: \"val1\"\n" +
			"\thave value: \"valY\""
		assertSpyHasMsg(t, spy, 1, wMsg)
	})

	t.Run("error - ExpectSetenv - not all required set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "key1")

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key1", "val1")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		wMsg := "" +
			"expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - ExpectSetenv - more than expected set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		isolateEnv(t, "key0", "key1", "key2")

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key0", "val0")
		spy.Setenv("key1", "val1")
		spy.Setenv("key2", "val2")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		wMsg := "" +
			"expected HUT not to set environment variable:\n" +
			"\t  have key: \"key2\"\n" +
			"\thave value: \"val2\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectSkipped - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSkipped()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		_ = affirm.Panic(t, func() { spy.Skip("msg", 0) })
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectSkipped", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		_ = affirm.Panic(t, func() { spy.Skip("msg", 0) })
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMsgs(t, spy, 1)
		wMsg := "" +
			"expected HUT to mark test as skipped:\n" +
			"\twant: false\n" +
			"\thave: true"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectFail - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		spy.Error("msg", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, true, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectFail without failure", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "expected HUT to call the t.Error* or t.Fatal* methods"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectError - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectError()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		spy.Error("msg", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, true, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectError", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectError()
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "expected HUT to call any of the t.Error* methods"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectFatal - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFatal()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		msg := affirm.Panic(t, func() { spy.Fatal("msg", 0) })
		affirm.NotNil(t, msg)
		affirm.Equal(t, FailNowMsg, *msg)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, true, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectFatal", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFatal()
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "expected HUT to call any of the t.Fatal* methods"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectHelpers - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(2)
		spy.Close()
		spy.Helper()
		spy.Helper()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectHelpers", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(2)
		spy.Close()
		spy.Helper()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected Helper to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectNames - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectNames(2)
		spy.Close()
		_ = spy.Name()
		_ = spy.Name()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectNames", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectNames(2)
		spy.Close()
		_ = spy.Name()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected Name to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectCleanups - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(2)
		spy.Close()
		spy.Cleanup(func() {})
		spy.Cleanup(func() {})
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectCleanups", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(2)
		spy.Close()
		spy.Cleanup(func() {})
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected Cleanup to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectTempDir - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(2)
		spy.Close()
		_ = spy.TempDir()
		_ = spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - does not matter and called", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(-1)
		spy.Close()
		_ = spy.TempDir()
		_ = spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - does not matter and not called", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(-1)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectTempDir", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(2)
		spy.Close()
		_ = spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected TempDir to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectLog - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(Equal, "msg %d", 0)
		spy.Close()
		spy.Logf("msg %d", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("ExpectLogEqual - multi message success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("msg 0")
		spy.ExpectLogEqual("msg 1")
		spy.Close()
		spy.Log("msg 0")
		spy.Log("msg 1")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectLogEqual - multi message surplus", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("msg 0")
		spy.ExpectLogEqual("msg 1")
		spy.Close()
		spy.Log("msg 0")
		spy.Log("msg 1")
		spy.Log("msg 2")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log 2 messages:\n" +
			"\t want: 2\n" +
			"\t have: 3\n" +
			"\textra: \"msg 2\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - ExpectLogEqual - second mismatch", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("msg 0")
		spy.ExpectLogEqual("msg 1")
		spy.Close()
		spy.Log("msg 0")
		spy.Log("msg X")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log message 1:\n" +
			"\tmatcher: equal\n" +
			"\t   want: \"msg 1\"\n" +
			"\t   have: \"msg X\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - ExpectLogEqual - multi message too few", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("msg 0")
		spy.ExpectLogEqual("msg 1")
		spy.Close()
		spy.Log("msg 0")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log message 1:\n" +
			"\tmatcher: equal\n" +
			"\t   want: \"msg 1\"\n" +
			"\t   have: <no message>"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - ExpectLog", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(Equal, "msg %d", 0)
		spy.Close()
		spy.Log("msg", 1)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log message 0:\n" +
			"\tmatcher: equal\n" +
			"\t   want: \"msg 0\"\n" +
			"\t   have: \"msg 1\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectLog - not contains - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(NotContains, "msg %d", 1)
		spy.Close()
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 2)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - ExpectLog - not contains", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(NotContains, "msg %d", 0)
		spy.Close()
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 1)
		spy.Logf("msg %d", 2)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log message 0:\n" +
			"\tmatcher: not-contains\n" +
			"\t   want: \"msg 0\"\n" +
			"\t   have: \"msg 0\\nmsg 1\\nmsg 2\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - expected no logs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Log("msg")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log no messages but got:\n" +
			"\thave: \"msg\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - empty message when no logs expected", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Log()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log no messages but got:\n" +
			"\thave: \"\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("do not care if messages were logged", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()
		spy.Close()
		spy.Log("msg")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("Helper called at least once", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(-1)
		spy.Close()
		spy.Helper()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMsgs(t, spy, 0)
	})

	t.Run("error - Helper not called", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(-1)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected Helper to be called N times:\n" +
			"\twant: >= 1\n" +
			"\thave: 0"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("error - ExpectLog - empty message", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("")
		spy.Close()
		spy.Log("msg")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, true, ti.Failed())
		wMsg := "" +
			"expected HUT to log message 0:\n" +
			"\tmatcher: equal\n" +
			"\t   want: \"\"\n" +
			"\t   have: \"msg\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})
}

func Test_Spy_assertFailed(t *testing.T) {
	t.Run("wantFailed false", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.assertFailed()

		// --- Then ---
		affirm.Equal(t, true, have)
	})
}

func Test_Spy_checkState_tabular(t *testing.T) {
	tt := []struct {
		testN string

		spy    func(ti *testing.T) *Spy
		action func(*Spy)
		want   string
	}{
		// Expect* calls.
		{
			"error - expect called on finished",
			func(ti *testing.T) *Spy { return New(ti).Close().Finish() },
			func(spy *Spy) { spy.ExpectError() },
			errExpectOnFinished,
		},
		{
			"error - expect called on closed",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.ExpectError() },
			errExpectOnClosed,
		},
		// Close calls.
		{
			"error - call to Close when finished",
			func(ti *testing.T) *Spy { return New(ti).Finish() },
			func(spy *Spy) { spy.Close() },
			errCloseOnFinished,
		},
		{
			"error - call to Close when already Closed",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.Close() },
			errDoubleClose,
		},
		// Mocked calls.
		{
			"error - call to mocked method when finished",
			func(ti *testing.T) *Spy { return New(ti).Close().Finish() },
			func(spy *Spy) { spy.Log("msg") },
			errActionOnFinished,
		},
		{
			"error - call to mocked method when not closed",
			func(ti *testing.T) *Spy { return New(ti) },
			func(spy *Spy) { spy.Log("msg") },
			errMockOnNotClosed,
		},
		// Assert* calls.
		{
			"error - assert call when not finished",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.assertExpectations() },
			errAssertOnNotFinished,
		},
		{
			"error - assert call when not closed",
			func(ti *testing.T) *Spy { return New(ti).Finish() },
			func(spy *Spy) { spy.assertExpectations() },
			errAssertOnNotClosed,
		},
		// Failed call.
		{
			"error - Failed call when not closed",
			func(ti *testing.T) *Spy { return New(ti) },
			func(spy *Spy) { spy.Failed() },
			errMockOnNotClosed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ti := &testing.T{}

			spy := tc.spy(ti)
			affirm.Equal(t, false, ti.Failed()) // Make sure it is not failed.

			// --- When ---
			have := affirm.Panic(t, func() { tc.action(spy) })

			// --- Then ---
			affirm.NotNil(t, have)
			affirm.Equal(t, tc.want, *have)
			affirm.Equal(t, true, spy.panicked)
			affirm.Equal(t, true, ti.Failed())
		})
	}
}
