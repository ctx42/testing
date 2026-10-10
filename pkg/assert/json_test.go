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

		// --- When ---
		have := JSON(tspy, ` {"hello": "world"} `, `{"hello": "world"}`)

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
		have := JSON(tspy, ` {"hello": "world"} `, `{"hello": "ms"}`)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("success bytes want", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := JSON(tspy, []byte(` {"hello": "world"} `), `{"hello": "world"}`)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success bytes have", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := JSON(tspy, ` {"hello": "world"} `, []byte(`{"hello": "world"}`))

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("success bytes", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		hJSON := []byte(`{"hello": "world"}`)

		// --- When ---
		have := JSON(tspy, []byte(` {"hello": "world"} `), hJSON)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		wJSON := ` {"hello": "world"} `
		hJSON := `{"hello": "ms"}`
		opt := check.WithTrail("type.field")

		// --- When ---
		have := JSON(tspy, wJSON, hJSON, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
