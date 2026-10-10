// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_True(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := True(true)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := True(false)

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "expected value to be true", err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := True(false, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be true:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with comment", func(t *testing.T) {
		// --- When ---
		err := True(false, "A%d", 42)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be true:\n" +
			"  comment: A42"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail and comment", func(t *testing.T) {
		// --- When ---
		err := True(false, "A%d", 42, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be true:\n" +
			"    trail: type.field\n" +
			"  comment: A42"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_False(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := False(false)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := False(true)

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "expected value to be false", err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := False(true, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected value to be false:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}
