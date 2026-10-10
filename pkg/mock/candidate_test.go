// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_candidate_betterThan_tabular(t *testing.T) {
	tt := []struct {
		testN string

		this  candidate
		other candidate
		want  bool
	}{
		{
			"both nil",
			candidate{call: nil},
			candidate{call: nil},
			false,
		},
		{
			"this call is nil",
			candidate{call: nil},
			candidate{call: &Call{}},
			false,
		},
		{
			"other call is nil",
			candidate{call: &Call{}},
			candidate{call: nil},
			true,
		},
		{
			"method names dont match",
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 0},
			candidate{call: &Call{cStack: cStack{Method: "B"}}, diffCnt: 0},
			true,
		},
		{
			"both same indicators",
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 1},
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 1},
			true,
		},
		{
			"this diffCnt lower",
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 1},
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 2},
			true,
		},
		{
			"other diffCnt lower",
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 2},
			candidate{call: &Call{cStack: cStack{Method: "A"}}, diffCnt: 1},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := tc.this.betterThan(tc.other)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
