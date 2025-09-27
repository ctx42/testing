// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/testcases"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Count(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Count(tspy, 2, "ab", "ab cd ab")

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Count(tspy, 1, 123, "ab cd ef")

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Count(tspy, 1, 123, "ab cd ef", opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_SameType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := SameType(tspy, true, true)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		msg := affirm.Panic(t, func() { SameType(tspy, 1, uint(1)) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})

	t.Run("log message with trails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		msg := affirm.Panic(t, func() { SameType(tspy, 1, uint(1), opt) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})
}

func Test_NotSameType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NotSameType(tspy, true, 42)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		msg := affirm.Panic(t, func() { NotSameType(tspy, 42, 42) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})

	t.Run("log message with trails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		msg := affirm.Panic(t, func() { NotSameType(tspy, 42, 42, opt) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})
}

func Test_Type(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		var target int

		// --- When ---
		have := Type(tspy, &target, 42)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		var target int

		// --- When ---
		msg := affirm.Panic(t, func() { Type(tspy, &target, uint(1)) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})

	t.Run("log message with trails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		var target int
		opt := check.WithTrail("type.field")

		// --- When ---
		msg := affirm.Panic(t, func() { Type(tspy, &target, uint(1), opt) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})
}

func Test_Fields(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Fields(tspy, 7, testcases.TA{})

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Fields(tspy, 1, &testcases.TA{})

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Fields(tspy, 1, &testcases.TA{}, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
