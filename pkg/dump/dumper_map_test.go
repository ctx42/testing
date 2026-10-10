// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/goldy"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_MapDumper(t *testing.T) {
	t.Run("interface keys are sorted", func(t *testing.T) {
		// --- Given ---
		val := map[any]int{"b": 2, "a": 1, 3: 3}

		// --- When ---
		have := MapDumper(New(), 0, reflect.ValueOf(val))

		// --- Then ---
		want := "" +
			"map[any]int{\n" +
			"  3: 3,\n" +
			"  \"a\": 1,\n" +
			"  \"b\": 2,\n" +
			"}"
		affirm.Equal(t, want, have)
	})

	t.Run("values do not print their type", func(t *testing.T) {
		// --- Given ---
		val := map[string][]int{"a": {1}, "b": {2}}

		// --- When ---
		have := MapDumper(New(), 0, reflect.ValueOf(val))

		// --- Then ---
		want := "" +
			"map[string][]int{\n" +
			"  \"a\": {\n" +
			"    1,\n" +
			"  },\n" +
			"  \"b\": {\n" +
			"    2,\n" +
			"  },\n" +
			"}"
		affirm.Equal(t, want, have)
	})

	t.Run("value leading spaces are kept", func(t *testing.T) {
		// --- Given ---
		long := strings.Repeat("y", 210)
		val := map[string]string{"k": "  x\n" + long}

		// --- When ---
		have := MapDumper(New(), 0, reflect.ValueOf(val))

		// --- Then ---
		want := "map[string]string{\n  \"k\":   x\n" + long + ",\n}"
		affirm.Equal(t, want, have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := MapDumper(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}

func Test_MapDumper_tabular(t *testing.T) {
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
			"flat map[int]testcases.T1",
			New(WithFlat, WithCompact, WithTimeFormat(TimeAsUnix)),
			map[int]testcases.TRec{0: {Int: 0}, 1: {Int: 1}},
			"map[int]testcases.TRec{0:{Int:0,Rec:nil},1:{Int:1,Rec:nil}}",
		},
		{
			"default map[int]testcases.TRec",
			New(WithTimeFormat(TimeAsUnix)),
			map[int]testcases.TRec{0: {Int: 0}, 1: {Int: 1}},
			goldy.Open(t, "testdata/map_of_structs.gld").String(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := MapDumper(tc.dmp, 0, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
