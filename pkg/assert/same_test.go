// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_Same(t *testing.T) {
	t.Run("pointers", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		ptr0 := &types.TPtr{Val: "0"}

		// --- When ---
		have := Same(tspy, ptr0, ptr0)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - want is value", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		w := types.TPtr{Val: "0"}
		h := &types.TPtr{Val: "0"}

		// --- When ---
		have := Same(tspy, w, h)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - have is value", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		w := &types.TPtr{Val: "0"}
		h := types.TPtr{Val: "0"}

		// --- When ---
		have := Same(tspy, w, h)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("error - not same pointers", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		ptr0 := &types.TPtr{Val: "0"}
		ptr1 := &types.TPtr{Val: "1"}

		// --- When ---
		have := Same(tspy, ptr0, ptr1)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		ptr0 := &types.TPtr{Val: "0"}
		ptr1 := &types.TPtr{Val: "1"}

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Same(tspy, ptr0, ptr1, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}

func Test_NotSame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		ptr0 := &types.TPtr{Val: "0"}
		ptr1 := &types.TPtr{Val: "1"}

		// --- When ---
		have := NotSame(tspy, ptr0, ptr1)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		ptr0 := &types.TPtr{Val: "0"}

		// --- When ---
		have := NotSame(tspy, ptr0, ptr0)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		ptr0 := &types.TPtr{Val: "0"}

		opt := check.WithTrail("type.field")

		// --- When ---
		have := NotSame(tspy, ptr0, ptr0, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
