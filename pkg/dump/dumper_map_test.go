// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"math"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/goldy"
	"github.com/ctx42/testing/pkg/testcases"
)

// tiedKey is a map key with a String method that makes all keys print the same
// through fmt.
type tiedKey struct{ V int }

// String implements [fmt.Stringer].
func (tiedKey) String() string { return "same" }

func Test_mapEntry_keyStr(t *testing.T) {
	// --- Given ---
	type T struct{ N int }

	var count int
	dpr := func(Dump, int, reflect.Value) string {
		count++
		return ""
	}
	dmp := New(WithDumper(T{}, dpr))
	ent := &mapEntry{key: reflect.ValueOf(T{})}

	// --- When ---
	have := ent.keyStr(dmp, 0)

	// --- Then ---
	affirm.Equal(t, "", have)
	affirm.Equal(t, "", ent.keyStr(dmp, 0))
	affirm.Equal(t, 1, count)
}

func Test_mapEntry_valStr(t *testing.T) {
	t.Run("value is dumped once", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }

		var count int
		dpr := func(Dump, int, reflect.Value) string {
			count++
			return "T"
		}
		dmp := New(WithDumper(T{}, dpr))
		ent := &mapEntry{val: reflect.ValueOf(T{})}

		// --- When ---
		have := ent.valStr(dmp, 0)

		// --- Then ---
		affirm.Equal(t, "T", have)
		affirm.Equal(t, "T", ent.valStr(dmp, 0))
		affirm.Equal(t, 1, count)
	})

	t.Run("empty dump is cached", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }

		var count int
		dpr := func(Dump, int, reflect.Value) string {
			count++
			return ""
		}
		dmp := New(WithDumper(T{}, dpr))
		ent := &mapEntry{val: reflect.ValueOf(T{})}

		// --- When ---
		have := ent.valStr(dmp, 0)

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.Equal(t, "", ent.valStr(dmp, 0))
		affirm.Equal(t, 1, count)
	})
}

func Test_mapEntry_takeValStr(t *testing.T) {
	// --- Given ---
	ent := &mapEntry{val: reflect.ValueOf(42)}

	// --- When ---
	have := ent.takeValStr(New(), 0)

	// --- Then ---
	affirm.Equal(t, "42", have)
	affirm.Equal(t, "", ent.str)
}

func Test_MapDumper(t *testing.T) {
	t.Run("cycle under pointer keys", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		defer debug.SetMaxStack(debug.SetMaxStack(64 << 20))

		type T struct{ M map[string]any }
		cyc := make(map[string]any)
		cyc["self"] = cyc
		val := map[*T]int{{M: cyc}: 1, {M: cyc}: 2}

		// --- When ---
		have := MapDumper(New(WithMaxDepth(3)), 0, reflect.ValueOf(val))

		// --- Then ---
		first := strings.Index(have, "}: 1,")
		second := strings.Index(have, "}: 2,")
		affirm.Equal(t, true, first >= 0)
		affirm.Equal(t, true, first < second)
	})

	t.Run("keys that print the same are ordered by dump", func(t *testing.T) {
		// --- Given ---
		val := map[tiedKey]int{{3}: 1, {1}: 1, {2}: 1, {5}: 1, {4}: 1, {6}: 1}

		// --- When ---
		have := MapDumper(New(WithFlat), 0, reflect.ValueOf(val))

		// --- Then ---
		want := "" +
			"map[dump.tiedKey]int{" +
			"{V: 1}: 1, {V: 2}: 1, {V: 3}: 1, " +
			"{V: 4}: 1, {V: 5}: 1, {V: 6}: 1}"
		affirm.Equal(t, want, have)
	})

	t.Run("tied keys are dumped once", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }
		val := make(map[T]int)
		for i := range 6 {
			val[T{N: i}] = i
		}

		var count int
		dpr := func(_ Dump, _ int, val reflect.Value) string {
			count++
			return strconv.Itoa(int(val.Field(0).Int()))
		}

		// --- When ---
		_ = MapDumper(New(WithDumper(T{}, dpr)), 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.Equal(t, 6, count)
	})

	t.Run("NaN key prints its value", func(t *testing.T) {
		// --- Given ---
		val := map[float64]int{math.NaN(): 1}

		// --- When ---
		have := MapDumper(New(), 0, reflect.ValueOf(val))

		// --- Then ---
		want := "" +
			"map[float64]int{\n" +
			"  NaN: 1,\n" +
			"}"
		affirm.Equal(t, want, have)
	})

	t.Run("values with tied keys are dumped once", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }
		val := make(map[float64]T)
		for i := range 6 {
			val[math.NaN()] = T{N: i}
		}

		var count int
		dpr := func(Dump, int, reflect.Value) string {
			count++
			return "T"
		}

		// --- When ---
		_ = MapDumper(New(WithDumper(T{}, dpr)), 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.Equal(t, 6, count)
	})

	t.Run("distinct keys are dumped when printed", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }
		val := make(map[int]T)
		for i := range 6 {
			val[i] = T{N: i}
		}

		order := make([]int, 0, 6)
		dpr := func(_ Dump, _ int, val reflect.Value) string {
			order = append(order, int(val.Field(0).Int()))
			return "T"
		}

		// --- When ---
		_ = MapDumper(New(WithDumper(T{}, dpr)), 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.DeepEqual(t, []int{0, 1, 2, 3, 4, 5}, order)
	})

	t.Run("NaN keys are ordered by value", func(t *testing.T) {
		// --- Given ---
		val := make(map[float64]int)
		for i := range 6 {
			val[math.NaN()] = i
		}

		// --- When ---
		have := MapDumper(New(), 0, reflect.ValueOf(val))

		// --- Then ---
		want := "" +
			"map[float64]int{\n" +
			"  NaN: 0,\n" +
			"  NaN: 1,\n" +
			"  NaN: 2,\n" +
			"  NaN: 3,\n" +
			"  NaN: 4,\n" +
			"  NaN: 5,\n" +
			"}"
		affirm.Equal(t, want, have)
	})

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
			"default map int int",
			New(),
			map[int]int{1: 10, 2: 20},
			"map[int]int{\n  1: 10,\n  2: 20,\n}",
		},
		{
			"default map int int with indent",
			New(WithIndent(2)),
			map[int]int{1: 10, 2: 20},
			"    map[int]int{\n      1: 10,\n      2: 20,\n    }",
		},
		{
			"flat map int int",
			New(WithFlat),
			map[int]int{1: 10, 2: 20},
			"map[int]int{1: 10, 2: 20}",
		},
		{
			"flat and compact map int int",
			New(WithFlat, WithCompact),
			map[int]int{1: 10, 2: 20},
			"map[int]int{1:10,2:20}",
		},
		{
			"flat map int testcases TRec",
			New(WithFlat, WithCompact, WithTimeFormat(TimeAsUnix)),
			map[int]testcases.TRec{0: {Int: 0}, 1: {Int: 1}},
			"map[int]testcases.TRec{0:{Int:0,Rec:nil},1:{Int:1,Rec:nil}}",
		},
		{
			"default map int testcases TRec",
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
