// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Nil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Nil(tspy, nil)

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
		have := Nil(tspy, 42)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - with option", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		// --- When ---
		have := Nil(tspy, 42, check.WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotNil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NotNil(tspy, 42)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		msg := affirm.Panic(t, func() { NotNil(tspy, nil) })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
	})

	t.Run("error - with option", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		affirm.Panic(t, func() { NotNil(tspy, nil, opt) })
	})
}
