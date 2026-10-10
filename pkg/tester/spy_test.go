// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func Test_find_match(t *testing.T) {
	t.Run("equal success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Equal, want: "abc def ghi"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - equal", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Equal, want: "abc def ghi"}

		// --- When ---
		have := ent.match("def")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("contains success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Contains, want: "def"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - contains", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Contains, want: "def"}

		// --- When ---
		have := ent.match("abc xyz ghi")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("not contain success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: NotContains, want: "jkl"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - not contains", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: NotContains, want: "def"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("regexp success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Regexp, want: "[a-z ]+"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - regexp", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Regexp, want: "abc def ghi"}

		// --- When ---
		have := ent.match("xyz")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("unknown as equal success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: "unknown", want: "abc def ghi"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - unknown as equal", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: "unknown", want: "abc def ghi"}

		// --- When ---
		have := ent.match("def")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("equal fold success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: EqualFold, want: "ABC DEF GHI"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - equal fold", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: EqualFold, want: "ABC DEF GHI"}

		// --- When ---
		have := ent.match("abc xyz ghi")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("contains fold success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: ContainsFold, want: "DEF"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - contains fold", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: ContainsFold, want: "DEF"}

		// --- When ---
		have := ent.match("abc xyz ghi")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("not contains fold success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: NotContainsFold, want: "JKL"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - not contains fold", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: NotContainsFold, want: "DEF"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.Equal(t, false, have)
	})
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

	t.Run("set invalid expected Helper calls panics", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- Then ---
		msg := affirm.Panic(t, func() { New(ti, -2) })
		affirm.NotNil(t, msg)
		want := "ExpectHelpers cnt must be greater or equal to minus one"
		affirm.Equal(t, want, *msg)
	})

	t.Run("more than one Helper calls value panics", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		msg := affirm.Panic(t, func() { New(ti, 1, 2) })

		// --- Then ---
		affirm.NotNil(t, msg)
		affirm.Equal(t, "New accepts at most one expectHelpers value", *msg)
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
		msg := affirm.Panic(t, func() { spy.ExpectCleanups(-1) })

		// --- Then ---
		affirm.NotNil(t, msg)
		want := "ExpectCleanups cnt must be greater or equal to zero"
		affirm.Equal(t, want, *msg)
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
		msg := affirm.Panic(t, func() { spy.Fatal("msg", 0) })

		// --- Then ---
		affirm.NotNil(t, msg)
		affirm.Equal(t, FailNowMsg, *msg)

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
		msg := affirm.Panic(t, func() { spy.Fatalf("msg %d", 0) })

		// --- Then ---
		affirm.NotNil(t, msg)
		affirm.Equal(t, FailNowMsg, *msg)
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

		// --- Then ---
		msg := affirm.Panic(t, func() { spy.FailNow() })
		affirm.NotNil(t, msg)
		affirm.Equal(t, FailNowMsg, *msg)
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

	t.Run("does not panic when called on finished Spy", func(t *testing.T) {
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

	t.Run("argument must be greater or equal minus one", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- Then ---
		msg := affirm.Panic(t, func() { spy.ExpectHelpers(-2) })
		affirm.NotNil(t, msg)
		want := "ExpectHelpers cnt must be greater or equal to minus one"
		affirm.Equal(t, want, *msg)
		affirm.Equal(t, true, spy.panicked)
	})

	t.Run("may be set only once", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)
		spy.ExpectHelpers(1)

		// --- Then ---
		msg := affirm.Panic(t, func() { spy.ExpectHelpers(1) })
		affirm.NotNil(t, msg)
		want := "ExpectHelpers may be called only once"
		affirm.Equal(t, want, *msg)
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
		msg := affirm.Panic(t, func() { spy.Skip("msg", 0) })

		// --- Then ---
		affirm.NotNil(t, msg)
		affirm.Equal(t, SkipNowMsg, *msg)
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
		msg := affirm.Panic(t, func() { spy.ExpectTempDir(-2) })

		// --- Then ---
		affirm.NotNil(t, msg)
		want := "ExpectTempDir cnt must be greater or equal to minus one"
		affirm.Equal(t, want, *msg)
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

	t.Run("GetTempDir called before ExpectTempDir", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 1)
		msg := "ExpectTempDir method must be called before GetTempDir"
		assertSpyHasMsg(t, spy, 0, msg)
	})

	t.Run("get not existing index", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(2)
		spy.Close()
		spy.TempDir()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.GetTempDir(2)

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, true, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		msg := "temp directory with index 2 does not exist"
		assertSpyHasMsg(t, spy, 0, msg)
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
	t.Run("panics on invalid regexp", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti, 0)

		// --- When ---
		msg := affirm.Panic(t, func() { spy.ExpectLog(Regexp, "[") })

		// --- Then ---
		wMsg := "" +
			"invalid ExpectLog regexp: error parsing regexp: " +
			"missing closing ]: `[`"
		affirm.Equal(t, wMsg, *msg)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg %d", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMsgs[1]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)

		ent = spy.wantLogMsgs[2]
		affirm.Equal(t, Regexp, ent.strategy)
		affirm.Equal(t, "msg [0-9]", ent.want)

		ent = spy.wantLogMsgs[3]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 6", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMsgs[1]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMsgs[1]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMsgs[1]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, EqualFold, ent.strategy)
		affirm.Equal(t, "MSG 0", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, ContainsFold, ent.strategy)
		affirm.Equal(t, "MSG 0", ent.want)
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
		ent := spy.wantLogMsgs[0]
		affirm.Equal(t, NotContainsFold, ent.strategy)
		affirm.Equal(t, "MSG 0", ent.want)
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
		msg := affirm.Panic(t, func() { spy.ExpectNames(-1) })

		// --- Then ---
		affirm.NotNil(t, msg)
		want := "ExpectNames cnt must be greater or equal to zero"
		affirm.Equal(t, want, *msg)
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
		spy.Name()

		// --- Then ---
		affirm.Equal(t, 1, spy.haveNamesCnt)
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Name()
		spy.Name()

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

	t.Run("cleanup panic is raised after other cleanups", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { panic("boom") })

		// --- When ---
		msg := affirm.Panic(t, func() { spy.Finish() })

		// --- Then ---
		affirm.Equal(t, "boom", *msg)
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
		assertSpyHasMgs(t, spy, 1)
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
		assertSpyHasMgs(t, spy, 2)
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
		assertSpyHasMgs(t, spy, 0)
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

	t.Run("called on panicked instance", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 1)
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectSetenv - error - not expected but set", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 1)
		wMsg := "" +
			"expected HUT not to set environment variable:\n" +
			"\t  have key: \"key0\"\n" +
			"\thave value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectSetenv - error - none set", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 2)
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

	t.Run("ExpectSetenv - error - wrong value set", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 2)
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

	t.Run("ExpectSetenv - error - not all required set", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 1)
		wMsg := "" +
			"expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectSetenv - error - more than expected set", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 1)
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
		assertSpyHasMgs(t, spy, 0)
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
		assertSpyHasMgs(t, spy, 1)
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
		assertSpyHasMgs(t, spy, 0)
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectError - fail", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectFatal - fail", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectHelpers - fail", func(t *testing.T) {
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
		spy.Name()
		spy.Name()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectNames - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectNames(2)
		spy.Close()
		spy.Name()
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectCleanups - fail", func(t *testing.T) {
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
		spy.TempDir()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - does not matter and called", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(-1)
		spy.Close()
		spy.TempDir()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, spy.Failed())
		affirm.Equal(t, false, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		removeTempDirs(t, spy)
		spy.ExpectTempDir(2)
		spy.Close()
		spy.TempDir()
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
		assertSpyHasMgs(t, spy, 0)
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectLogEqual - multi message surplus fails", func(t *testing.T) {
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

	t.Run("ExpectLogEqual - multi message second mismatch", func(t *testing.T) {
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

	t.Run("ExpectLogEqual - multi message too few", func(t *testing.T) {
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

	t.Run("ExpectLog - fail", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectLog - not contains - fail", func(t *testing.T) {
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

	t.Run("expected no logs - fail", func(t *testing.T) {
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
		assertSpyHasMgs(t, spy, 0)
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
		assertSpyHasMgs(t, spy, 0)
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

	t.Run("ExpectLog - empty message - fail", func(t *testing.T) {
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
	t.Run("returns true if wantFailed false", func(t *testing.T) {
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
			"expect called on finished",
			func(ti *testing.T) *Spy { return New(ti).Close().Finish() },
			func(spy *Spy) { spy.ExpectError() },
			errExpectOnFinished,
		},
		{
			"expect called on closed",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.ExpectError() },
			errExpectOnClosed,
		},
		// Close calls.
		{
			"call to Close when finished",
			func(ti *testing.T) *Spy { return New(ti).Finish() },
			func(spy *Spy) { spy.Close() },
			errCloseOnFinished,
		},
		{
			"call to Close when already Closed",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.Close() },
			errDoubleClose,
		},
		// Mocked calls.
		{
			"call to mocked method when finished",
			func(ti *testing.T) *Spy { return New(ti).Close().Finish() },
			func(spy *Spy) { spy.Log("msg") },
			errActionOnFinished,
		},
		{
			"call to mocked method when not closed",
			func(ti *testing.T) *Spy { return New(ti) },
			func(spy *Spy) { spy.Log("msg") },
			errMockOnNotClosed,
		},
		// Assert* calls.
		{
			"assert* call when not finished",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.assertExpectations() },
			errAssertOnNotFinished,
		},
		{
			"assert* call when not closed",
			func(ti *testing.T) *Spy { return New(ti).Finish() },
			func(spy *Spy) { spy.assertExpectations() },
			errAssertOnNotClosed,
		},
		// Failed call.
		{
			"Failed call when not closed",
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

			// --- Then ---
			msg := affirm.Panic(t, func() { tc.action(spy) })
			affirm.NotNil(t, msg)
			affirm.Equal(t, tc.want, *msg)
			affirm.Equal(t, true, spy.panicked)
			affirm.Equal(t, true, ti.Failed())
		})
	}
}

// =============================================================================
// ================================== HELPERS ==================================
// =============================================================================

// assertSpyHasMgs asserts spy has saved a "want" number of messages. On
// failure, it marks the test as failed, returns false, but continues execution.
func assertSpyHasMgs(t *testing.T, spy *Spy, want int) bool {
	t.Helper()
	have := len(spy.savedMsgs)
	if want != have {
		var msg strings.Builder
		_, _ = fmt.Fprintf(&msg, "expected Spy to have saved messages:\n"+
			"\twant: %d\n"+
			"\thave: %d", want, have)
		if have > 0 {
			msg.WriteString("\n\tmessages:\n")
			for idx, sm := range spy.savedMsgs {
				_, _ = fmt.Fprintf(&msg, "\t\t%d: %q\n", idx, sm)
			}
		}
		t.Error(msg.String())
		return false
	}
	return true
}

func Test_assertSpyHasMgs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMsgs = []string{"abc", "def"}
		have := assertSpyHasMgs(ti, spy, 2)
		if !have || ti.Failed() {
			t.Error("expected success")
		}
	})
	t.Run("error", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMsgs = []string{"abc", "def"}
		have := assertSpyHasMgs(ti, spy, 3)
		if have || !ti.Failed() {
			t.Error("expected failure")
		}
	})
}

// assertSpyHasMsg asserts Spy has saved a message with given idx equal to
// "want". On failure, it marks the test as failed, returns false, but
// continues execution.
func assertSpyHasMsg(t *testing.T, spy *Spy, idx int, want string) bool {
	t.Helper()
	if idx >= len(spy.savedMsgs) {
		t.Errorf("expected Spy saved message with index %d to exist", idx)
		return false
	}
	have := spy.savedMsgs[idx]
	if want != have {
		format := "" +
			"expected Spy saved message with index %d to be:\n" +
			"\twant: %q\n" +
			"\thave: %q"
		t.Errorf(format, idx, want, have)
		return false
	}
	return true
}

func Test_assertSpyHasMsg(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMsgs = []string{"abc", "def"}
		have := assertSpyHasMsg(ti, spy, 1, "def")
		if !have || ti.Failed() {
			t.Error("expected success")
		}
	})
	t.Run("error", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMsgs = []string{"abc", "def"}
		have := assertSpyHasMsg(ti, spy, 1, "xyz")
		if have || !ti.Failed() {
			t.Error("expected failure")
		}
	})
}

// isolateEnv sets keys to an empty value on t so they are restored when t
// ends. A Spy bound to a zero testing.T calls its Setenv, whose restore
// never runs.
func isolateEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		t.Setenv(key, "")
	}
}

// removeTempDirs removes, when t ends, the parent directories testing creates
// for the temporary directories returned by spy. A Spy bound to a zero
// testing.T calls its TempDir, whose cleanup never runs.
func removeTempDirs(t *testing.T, spy *Spy) {
	t.Helper()
	t.Cleanup(func() {
		for _, dir := range spy.haveTempDirs {
			_ = os.RemoveAll(filepath.Dir(dir))
		}
	})
}
