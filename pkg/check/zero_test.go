// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/cases"
)

func Test_Zero(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Zero(time.Time{})

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := Zero(time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC))

		// --- Then ---
		affirm.NotNil(t, err)
	})
}

func Test_zeroError(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := zeroError(42)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected argument to be zero value:\n" +
			"  want: <zero>\n" +
			"  have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := zeroError(42, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected argument to be zero value:\n" +
			"  trail: type.field\n" +
			"   want: <zero>\n" +
			"   have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Zero_ZENValues(t *testing.T) {
	for _, tc := range cases.ZENValues() {
		t.Run("Zero "+tc.Desc, func(t *testing.T) {
			// --- When ---
			have := Zero(tc.Val)

			// --- Then ---
			if tc.IsZero && have != nil {
				format := "expected nil error:\n  have: %#v"
				t.Errorf(format, have)
			}
			if !tc.IsZero && have == nil {
				format := "expected not-nil error:\n  have: %#v"
				t.Errorf(format, have)
			}
		})
	}
}

func Test_NotZero(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := NotZero(time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := NotZero(time.Time{})

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected argument not to be zero value:\n" +
			"  want: <non-zero>\n" +
			"  have: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NotZero(time.Time{}, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected argument not to be zero value:\n" +
			"  trail: type.field\n" +
			"   want: <non-zero>\n" +
			"   have: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)"
		affirm.Equal(t, wMsg, err.Error())
	})
}
