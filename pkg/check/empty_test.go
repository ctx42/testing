// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Empty(t *testing.T) {
	t.Run("error - self referencing pointer", func(t *testing.T) {
		// --- Given ---
		var x any
		x = &x

		// --- When ---
		err := Empty(&x)

		// --- Then ---
		affirm.NotNil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Empty("")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := Empty("abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected argument to be empty:\n" +
			"  want: <empty>\n" +
			"  have: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Empty("abc", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected argument to be empty:\n" +
			"  trail: type.field\n" +
			"   want: <empty>\n" +
			"   have: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Empty_tabular(t *testing.T) {
	for _, tc := range testcases.ZENValues() {
		t.Run("Empty "+tc.Desc, func(t *testing.T) {
			// --- When ---
			err := Empty(tc.Val)

			// --- Then ---
			if have := err == nil; tc.IsEmpty != have {
				format := "expected Empty to accept %#v: %v\n  have: %v"
				t.Errorf(format, tc.Val, tc.IsEmpty, err)
			}
		})
	}
}

func Test_NotEmpty(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := NotEmpty("abc")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := NotEmpty("")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected non-empty value"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NotEmpty("", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected non-empty value:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}
