// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_wideRunes(t *testing.T) {
	t.Run("R16 sorted without overlaps", func(t *testing.T) {
		var prev uint16
		for i, rng := range wideRunes.R16 {
			assert.Equal(t, uint16(1), rng.Stride)
			assert.True(t, rng.Lo <= rng.Hi)
			assert.True(t, i == 0 || prev < rng.Lo)
			prev = rng.Hi
		}
	})

	t.Run("R32 sorted without overlaps", func(t *testing.T) {
		var prev uint32
		for i, rng := range wideRunes.R32 {
			assert.Equal(t, uint32(1), rng.Stride)
			assert.True(t, rng.Lo <= rng.Hi)
			assert.True(t, i == 0 || prev < rng.Lo)
			prev = rng.Hi
		}
		assert.True(t, wideRunes.R32[0].Lo > 0xFFFF)
	})
}

func Test_lineWidth_tabular(t *testing.T) {
	tt := []struct {
		testN string

		line string
		want int
	}{
		{"empty", "", 0},
		{"ascii", "abc", 3},
		{"multibyte rune", "ąb", 2},
		{"leading tab", "\tab", 6},
		{"tab to next stop", "a\tb", 5},
		{"two tabs", "\t\tx", 9},
		{"wide ideograph", "名a", 3},
		{"wide hangul", "한", 2},
		{"fullwidth letter", "Ａ", 2},
		{"wide runes before tab", "名名名\tx", 9},
		{"combining mark", "e\u0301", 1},
		{"wide tangut", "\U00017000", 2},
		{"wide supplementary ideograph", "\U00020000", 2},
		{"wide emoji", "\U0001F680", 2},
		{"ambiguous width", "\u3248", 1},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := lineWidth(tc.line)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
