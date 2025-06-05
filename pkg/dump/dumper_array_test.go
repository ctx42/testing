// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_ArrayDumper(t *testing.T) {
	t.Run("error - invalid kind", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := ArrayDumper(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}

func Test_ArrayDumper_tabular(t *testing.T) {
	var nilArr [2]int

	tt := []struct {
		testN string

		dmp  Dump
		val  any
		want string
	}{
		{
			"default",
			New(),
			[2]int{0, 1},
			"[2]int{\n  0,\n  1,\n}",
		},
		{
			"nil array",
			New(),
			nilArr,
			"[2]int{\n  0,\n  0,\n}",
		},
		{
			"default with indent",
			New(WithIndent(2)),
			[2]int{0, 1},
			"    [2]int{\n      0,\n      1,\n    }",
		},
		{
			"flat array",
			New(WithFlat),
			[2]int{0, 1},
			"[2]int{0, 1}",
		},
		{
			"flat and compact array",
			New(WithFlat, WithCompact),
			[2]int{0, 1},
			"[2]int{0,1}",
		},
		{
			"flat array empty int",
			New(WithFlat),
			[2]int{},
			"[2]int{0, 0}",
		},
		{
			"flat slice empty",
			New(WithFlat),
			[]int{},
			"[]int{}",
		},
		{
			"flat array empty any",
			New(WithFlat),
			[2]any{},
			"[2]any{nil, nil}",
		},
		{
			"flat array empty any",
			(func() Dump {
				dmp := New(WithFlat)
				dmp.UseAny = false
				return dmp
			})(),
			[2]any{},
			"[2]interface {}{nil, nil}",
		},
		{
			"flat array of map[string]int",
			New(WithFlat),
			[...]map[string]int{
				{"A": 1},
				{"b": 2},
			},
			`[2]map[string]int{{"A": 1}, {"b": 2}}`,
		},
		{
			"array of map[int]int",
			New(),
			[...]map[int]int{
				{1: 10},
				{2: 20},
			},
			"[2]map[int]int{\n  {\n    1: 10,\n  },\n  {\n    2: 20,\n  },\n}",
		},
		{
			"array of structs",
			New(),
			[]types.T1{{Int: 1}, {Int: 2}},
			"[]types.T1{\n" +
				"  {\n    Int: 1,\n    T1: nil,\n  },\n" +
				"  {\n    Int: 2,\n    T1: nil,\n  },\n" +
				"}",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := ArrayDumper(tc.dmp, 0, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
