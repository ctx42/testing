// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_NewRow(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		// --- When ---
		have := NewRow("name", "format")

		// --- Then ---
		affirm.Equal(t, "name", have.Name)
		affirm.Equal(t, "format", have.Format)
		affirm.Nil(t, have.Args)
		affirm.Equal(t, 0, len(have.Args))
	})

	t.Run("with args", func(t *testing.T) {
		// --- When ---
		have := NewRow("name", "format", "a", "b", "c")

		// --- Then ---
		affirm.Equal(t, "name", have.Name)
		affirm.Equal(t, "format", have.Format)
		affirm.DeepEqual(t, []any{"a", "b", "c"}, have.Args)
	})
}

func Test_Row_String_tabular(t *testing.T) {
	tt := []struct {
		testN string

		row  Row
		want string
	}{
		{
			"missing arguments",
			Row{Name: "name", Format: "%s"},
			"%!s(MISSING)",
		},
		{
			"single argument",
			Row{Name: "name", Format: "this %s", Args: []any{"value"}},
			"this value",
		},
		{
			"multiple arguments",
			Row{Name: "name", Format: "this %d %s", Args: []any{1, "that"}},
			"this 1 that",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := tc.row.String()

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_PaddedName_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name   string
		length int
		want   string
	}{
		{"length equal to name length", "name", 4, "name"},
		{"length less than name length", "name", 3, "name"},
		{"length zero", "name", 0, "name"},
		{"length greater than name length", "name", 6, "  name"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			row := Row{Name: "name"}

			// --- When ---
			have := row.PadName(tc.length)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
