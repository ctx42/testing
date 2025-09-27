// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/testcases"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Time(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, testcases.WAW)

		// --- When ---
		got := Time(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, true, got)
		affirm.Equal(t, true, want.Equal(have))
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 6, 0, testcases.WAW)

		// --- When ---
		got := Time(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, false, got)
		affirm.Equal(t, false, want.Equal(have))
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 6, 0, testcases.WAW)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := Time(tspy, want, have, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
		affirm.Equal(t, false, want.Equal(have))
	})
}

func Test_Exact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		got := Exact(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, true, got)
		affirm.Equal(t, true, want.Equal(have))
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, testcases.WAW)

		// --- When ---
		got := Exact(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, false, got)
		affirm.Equal(t, false, want.Equal(have))
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := Exact(tspy, want, have, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
		affirm.Equal(t, false, want.Equal(have))
	})
}

func Test_Before(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		have := Before(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		have := Before(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := check.WithTrail("type.field")

		// --- When ---
		have := Before(tspy, mark, date, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_After(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		got := After(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		got := After(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := After(tspy, mark, date, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}

func Test_BeforeOrEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		got := BeforeOrEqual(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		got := BeforeOrEqual(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := BeforeOrEqual(tspy, mark, date, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}

func Test_AfterOrEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		date := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		got := AfterOrEqual(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		got := AfterOrEqual(tspy, mark, date)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := AfterOrEqual(tspy, mark, date, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}

func Test_Within(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		got := Within(tspy, want, "1s", have)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, int(500*time.Millisecond), time.UTC)

		// --- When ---
		got := Within(tspy, want, "1s", have)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("         trail: type.field\n")
		tspy.Close()

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, int(500*time.Millisecond), time.UTC)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := Within(tspy, want, "1s", have, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("want is not time.Time", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		wMsg := "[want] failed to parse time:\n  cause: not supported time type"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		have := time.Date(2000, 1, 2, 4, 4, 6, 0, testcases.WAW)

		// --- When ---
		got := Within(tspy, true, "1s", have)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("have is not time.Time", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		wMsg := "[have] failed to parse time:\n  cause: not supported time type"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		want := time.Date(2000, 1, 2, 4, 4, 6, 0, testcases.WAW)

		// --- When ---
		got := Within(tspy, want, "1s", true)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}

func Test_Recent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		have := time.Now().Add(-4 * time.Second)

		// --- When ---
		got := Recent(tspy, have)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		have := time.Now().Add(-10 * time.Second)

		// --- When ---
		got := Recent(tspy, have)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("         trail: type.field\n")
		tspy.Close()

		have := time.Now().Add(-10 * time.Second)
		opt := check.WithTrail("type.field")

		// --- When ---
		got := Recent(tspy, have, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}

func Test_Zone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Zone(tspy, time.UTC, time.UTC)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Zone(tspy, nil, testcases.WAW)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Zone(tspy, nil, testcases.WAW, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_Duration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Duration(tspy, time.Second, time.Second)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Duration(tspy, time.Second, 2*time.Second)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Duration(tspy, time.Second, 2*time.Second, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
