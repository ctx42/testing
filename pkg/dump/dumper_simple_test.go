// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_SampleDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val    any
		indent int
		level  int
		want   string
	}{
		{"bool true", true, 0, 0, "true"},
		{"bool false", false, 0, 0, "false"},
		{"int", 123, 0, 0, "123"},
		{"int8", int8(123), 0, 0, "123"},
		{"int16", int16(123), 0, 0, "123"},
		{"int32", int32(123), 0, 0, "123"},
		{"int64", int64(123), 0, 0, "123"},
		{"uint", uint(123), 0, 0, "123"},
		{"uint16", uint16(123), 0, 0, "123"},
		{"uint32", uint32(123), 0, 0, "123"},
		{"uint64", uint64(123), 0, 0, "123"},
		{"float32", float32(12.3), 0, 0, "12.3"},
		{"float64", 12.3, 0, 0, "12.3"},
		{"float64 very small", 0.00000000000003, 0, 0, "0.00000000000003"},
		{"string", "string", 0, 0, `"string"`},

		{"with indent", 123, 2, 0, "    123"},
		{"with level", 123, 0, 1, "  123"},
		{"with indent and level", 123, 2, 1, "      123"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			dmp := New(WithIndent(tc.indent))

			// --- When ---
			have := SimpleDumper(dmp, tc.level, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_SampleDumper(t *testing.T) {
	t.Run("string with Flat false and FlatStrings off", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithFlatStrings(0))

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf("str0\nstr1\n"))

		// --- Then ---
		affirm.Equal(t, "str0\nstr1\n", have)
	})

	t.Run("string with Flat false and FlatStrings default", func(t *testing.T) {
		// --- Given ---
		dmp := New()

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf("str0\nstr1\n"))

		// --- Then ---
		affirm.Equal(t, "\"str0\\nstr1\\n\"", have)
	})

	t.Run("long string with Flat false and FlatStrings default", func(t *testing.T) {
		// --- Given ---
		// The string is one character longer than the default value of
		// [Dump.FlatStrings]
		str := strings.Repeat("a", 100) + "\n" + strings.Repeat("a", 100)
		dmp := New()

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf(str))

		// --- Then ---
		want := "" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		affirm.Equal(t, want, have)
	})

	t.Run("string with Flat true and FlatStrings default", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithFlat)

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf("str0\nstr1\n"))

		// --- Then ---
		affirm.Equal(t, "\"str0\\nstr1\\n\"", have)
	})

	t.Run("string with Dump.flatStrings true", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{flatStrings: true}

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf("str0\nstr1\n"))

		// --- Then ---
		affirm.Equal(t, "\"str0\\nstr1\\n\"", have)
	})

	t.Run("default single line string format", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{}

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf("str0"))

		// --- Then ---
		affirm.Equal(t, "\"str0\"", have)
	})

	t.Run("unsupported kind", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := SimpleDumper(dmp, 2, reflect.ValueOf(struct{}{}))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})

	t.Run("with WithAlwaysMultiline", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithFlatStrings(0), WithAlwaysMultiline)

		// --- When ---
		have := SimpleDumper(dmp, 0, reflect.ValueOf("str0"))

		// --- Then ---
		affirm.Equal(t, "str0", have)
	})
}
