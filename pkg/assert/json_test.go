// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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
		have := `{"hello": "world"}`

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
		have := `{"hello": "ms"}`

		// --- When ---
		got := JSON(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, false, got)
	})

	t.Run("success bytes want", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		want := []byte(` {"hello": "world"} `)
		have := `{"hello": "world"}`

		// --- When ---
		got := JSON(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("success bytes have", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		want := ` {"hello": "world"} `
		have := []byte(`{"hello": "world"}`)

		// --- When ---
		got := JSON(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("success bytes", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		want := []byte(` {"hello": "world"} `)
		have := []byte(`{"hello": "world"}`)

		// --- When ---
		got := JSON(tspy, want, have)

		// --- Then ---
		affirm.Equal(t, true, got)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		want := ` {"hello": "world"} `
		have := `{"hello": "ms"}`
		opt := check.WithTrail("type.field")

		// --- When ---
		got := JSON(tspy, want, have, opt)

		// --- Then ---
		affirm.Equal(t, false, got)
	})
}
