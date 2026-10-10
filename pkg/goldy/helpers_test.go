// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package goldy

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_crlfToLF_tabular(t *testing.T) {
	tt := []struct {
		testN string

		line string
		want string
	}{
		{"crlf", "abc\r\n", "abc\n"},
		{"lf", "abc\n", "abc\n"},
		{"no line ending", "abc", "abc"},
		{"inner cr", "a\rb\n", "a\rb\n"},
		{"empty", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := crlfToLF([]byte(tc.line))

			// --- Then ---
			affirm.Equal(t, tc.want, string(have))
		})
	}
}
