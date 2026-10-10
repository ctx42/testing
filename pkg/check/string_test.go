// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_EqualFold(t *testing.T) {
	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := EqualFold("ABC", "xyz", WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected strings to be equal (case-insensitive):\n" +
			"  trail: type.field\n" +
			"   want: \"ABC\"\n" +
			"   have: \"xyz\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_EqualFold_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have string
	}{
		{"abc and abc", "abc", "abc"},
		{"ABC and abc", "ABC", "abc"},
		{"abc and ABC", "abc", "ABC"},
		{"AbC and aBc", "AbC", "aBc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := EqualFold(tc.want, tc.have)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_EqualFold_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have string
	}{
		{"abc and xyz", "abc", "xyz"},
		{"ABC and xyz", "ABC", "xyz"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := EqualFold(tc.want, tc.have)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "" +
				"expected strings to be equal (case-insensitive):\n" +
				"  want: %q\n" +
				"  have: %q"
			wMsg = fmt.Sprintf(wMsg, tc.want, tc.have)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}

func Test_Contain(t *testing.T) {
	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := Contain("abc", "xyz", WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected string to contain substring:\n" +
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
		{"b and abc", "b", "abc"},
		{"bc and abc", "bc", "abc"},
		{"a and abc", "a", "abc"},
		{"abc and abc", "abc", "abc"},
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
		{"abc and xy", "abc", "xy"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Contain(tc.sub, tc.s)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "" +
				"expected string to contain substring:\n" +
				"     string: %q\n" +
				"  substring: %q"
			wMsg = fmt.Sprintf(wMsg, tc.s, tc.sub)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}

func Test_ContainFold(t *testing.T) {
	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := ContainFold("abc", "xyz", WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected string to contain substring (case-insensitive):\n" +
			"      trail: type.field\n" +
			"     string: \"xyz\"\n" +
			"  substring: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_ContainFold_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have string
	}{
		{"b and abc", "b", "abc"},
		{"B and abc", "B", "abc"},
		{"bc and ABC", "bc", "ABC"},
		{"ABC and abc", "ABC", "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := ContainFold(tc.want, tc.have)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_ContainFold_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want string
		have string
	}{
		{"xyz and abc", "xyz", "abc"},
		{"XYZ and abc", "XYZ", "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := ContainFold(tc.want, tc.have)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "" +
				"expected string to contain substring (case-insensitive):\n" +
				"     string: %q\n" +
				"  substring: %q"
			wMsg = fmt.Sprintf(wMsg, tc.have, tc.want)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}

func Test_NotContain(t *testing.T) {
	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := NotContain("abc", "abc", WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected string not to contain substring:\n" +
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
		{"abc and xy", "abc", "xy"},
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
		{"abc and b", "abc", "b"},
		{"abc and bc", "abc", "bc"},
		{"abc and a", "abc", "a"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := NotContain(tc.sub, tc.s)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "" +
				"expected string not to contain substring:\n" +
				"     string: %q\n" +
				"  substring: %q"
			wMsg = fmt.Sprintf(wMsg, tc.s, tc.sub)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}
