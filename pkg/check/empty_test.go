// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/cases"
)

func Test_Empty(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Empty("")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
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

	t.Run("additional message rows added", func(t *testing.T) {
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

func Test_Empty_ZENValues(t *testing.T) {
	for _, tc := range cases.ZENValues() {
		t.Run("Empty "+tc.Desc, func(t *testing.T) {
			// --- When ---
			have := Empty(tc.Val)

			// --- Then ---
			if tc.IsEmpty && have != nil {
				format := "expected nil error:\n  have: %#v"
				t.Errorf(format, have)
			}
			if !tc.IsEmpty && have == nil {
				format := "expected not-nil error:\n  have: %#v"
				t.Errorf(format, have)
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

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := NotEmpty("")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected non-empty value"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
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
