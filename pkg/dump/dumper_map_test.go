// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/goldy"
)

func Test_mapDumper_tabular(t *testing.T) {
	var nilMap map[string]int

	tt := []struct {
		testN string

		dmp  Dump
		val  any
		want string
	}{
		{
			"empty map",
			New(WithFlat),
			map[string]int{},
			`map[string]int{}`,
		},
		{
			"nil map",
			New(),
			nilMap,
			`map[string]int(nil)`,
		},
		{
			"default map[int]int",
			New(),
			map[int]int{1: 10, 2: 20},
			"map[int]int{\n  1: 10,\n  2: 20,\n}",
		},
		{
			"default map[int]int ith indent",
			New(WithIndent(2)),
			map[int]int{1: 10, 2: 20},
			"    map[int]int{\n      1: 10,\n      2: 20,\n    }",
		},
		{
			"flat map[int]int",
			New(WithFlat),
			map[int]int{1: 10, 2: 20},
			"map[int]int{1: 10, 2: 20}",
		},
		{
			"flat and compact map[int]int",
			New(WithFlat, WithCompact),
			map[int]int{1: 10, 2: 20},
			"map[int]int{1:10,2:20}",
		},
		{
			"flat map[int]types.T1",
			New(WithFlat, WithCompact, WithTimeFormat(TimeAsUnix)),
			map[int]types.T1{0: {Int: 0}, 1: {Int: 1}},
			"map[int]types.T1{0:{Int:0,T1:nil},1:{Int:1,T1:nil}}",
		},
		{
			"default map[int]types.T1",
			New(WithTimeFormat(TimeAsUnix)),
			map[int]types.T1{0: {Int: 0}, 1: {Int: 1}},
			goldy.New(t, "testdata/map_of_structs.gld").String(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := mapDumper(tc.dmp, 0, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
