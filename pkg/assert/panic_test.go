// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Panic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Panic(tspy, func() { panic("test") })

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Panic(tspy, func() {})

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field")
		tspy.Close()

		// --- When ---
		have := Panic(tspy, func() {}, check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NoPanic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NoPanic(tspy, func() {})

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := NoPanic(tspy, func() { panic("test") })

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("        trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := NoPanic(tspy, func() { panic("test") }, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_PanicContain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := PanicContain(tspy, "def", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := PanicContain(tspy, "xyz", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("        trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := PanicContain(tspy, "xyz", func() { panic("abc def ghi") }, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_PanicMsg(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := PanicMsg(tspy, func() { panic("abc def ghi") })

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, "abc def ghi", *have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := PanicMsg(tspy, func() {})

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field")
		tspy.Close()

		// --- When ---
		have := PanicMsg(tspy, func() {}, check.WithTrail("type.field"))

		// --- Then ---
		affirm.Nil(t, have)
	})
}
