// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Nil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Nil(nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := Nil(42)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be nil:\n" +
			"  want: nil\n" +
			"  have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := Nil(42, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be nil:\n" +
			"  trail: type.field\n" +
			"   want: nil\n" +
			"   have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NotNil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := NotNil(42)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := NotNil(nil)

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "expected non-nil value", err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := NotNil(nil, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected non-nil value:\n  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}
