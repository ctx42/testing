// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_JSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		want := ` {"hello": "world"} `
		have := " { \"hello\"\t:\n\n \"world\"\n\n\t}    "

		// --- When ---
		got := JSON(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		want := ` {"hello": "world"} `
		have := " { \"hello\"\t:\n\n \"ms\"\n\n\t}    "

		// --- When ---
		got := JSON(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		want := ` {"hello": "world"} `
		have := " { \"hello\"\t:\n\n \"ms\"\n\n\t}    "
		opt := check.WithTrail("type.field")

		// --- When ---
		got := JSON(tspy, want, have, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}
