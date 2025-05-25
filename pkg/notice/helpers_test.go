// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Indent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "")

		// --- Then ---
		affirm.Equal(t, "", have)
	})

	t.Run("one line", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc")

		// --- Then ---
		affirm.Equal(t, "abc", have)
	})

	t.Run("multiple lines", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			" abc\n" +
			" def\n" +
			" ghi"
		affirm.Equal(t, want, have)
	})

	t.Run("empty lines are not indented", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc\ndef\n\nghi")

		// --- Then ---
		want := "" +
			" abc\n" +
			" def\n" +
			"\n" +
			" ghi"
		affirm.Equal(t, want, have)
	})

	t.Run("use tabs", func(t *testing.T) {
		// --- When ---
		have := Indent(1, '\t', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			"\tabc\n" +
			"\tdef\n" +
			"\tghi"
		affirm.Equal(t, want, have)
	})

	t.Run("no ident", func(t *testing.T) {
		// --- When ---
		have := Indent(0, ' ', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			"abc\n" +
			"def\n" +
			"ghi"
		affirm.Equal(t, want, have)
	})
}

func Test_Pad_tabular(t *testing.T) {
	tt := []struct {
		testN string

		str    string
		length int
		want   string
	}{
		{"empty string zero pad", "", 0, ""},
		{"empty string with pad", "", 3, "   "},
		{"string with pad longer than string length", "abc", 5, "  abc"},
		{"string with pad equal to string length", "abc", 3, "abc"},
		{"string with pad shorter than string length", "abc", 2, "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := Pad(tc.str, tc.length)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
