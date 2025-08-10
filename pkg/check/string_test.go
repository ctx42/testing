// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Contain(t *testing.T) {
	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Contain("abc", "xyz", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected string to contain substring:\n" +
			"      trail: type.field\n" +
			"     string: \"xyz\"\n" +
			"  substring: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Contain_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have string
	}{
		{"1", "b", "abc"},
		{"2", "bc", "abc"},
		{"3", "a", "abc"},
		{"4", "abc", "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Contain(tc.want, tc.have)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_Contain_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		s   string
		sub string
	}{
		{"1", "abc", "xy"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Contain(tc.sub, tc.s)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "expected string to contain substring:\n" +
				"     string: %q\n" +
				"  substring: %q"
			wMsg = fmt.Sprintf(wMsg, tc.s, tc.sub)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}

func Test_NotContain(t *testing.T) {
	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NotContain("abc", "abc", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected string not to contain substring:\n" +
			"      trail: type.field\n" +
			"     string: \"abc\"\n" +
			"  substring: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NotContain_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have string
	}{
		{"1", "abc", "xy"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := NotContain(tc.want, tc.have)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_NotContain_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		s   string
		sub string
	}{
		{"1", "abc", "b"},
		{"2", "abc", "bc"},
		{"3", "abc", "a"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := NotContain(tc.sub, tc.s)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "expected string not to contain substring:\n" +
				"     string: %q\n" +
				"  substring: %q"
			wMsg = fmt.Sprintf(wMsg, tc.s, tc.sub)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}
