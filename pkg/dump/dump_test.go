// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/goldy"
)

func Test_WithFlat(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithFlat(dmp)

	// --- Then ---
	affirm.Equal(t, true, dmp.Flat)
}

func Test_WithFlatStrings(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithFlatStrings(123)(dmp)

	// --- Then ---
	affirm.Equal(t, 123, dmp.FlatStrings)
}

func Test_WithCompact(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithCompact(dmp)

	// --- Then ---
	affirm.Equal(t, true, dmp.Compact)
}

func Test_WithPtrAddr(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithPtrAddr(dmp)

	// --- Then ---
	affirm.Equal(t, true, dmp.PtrAddr)
}

func Test_WithTimeFormat(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	opt := WithTimeFormat(TimeAsUnix)

	// --- Then ---
	opt(dmp)
	affirm.Equal(t, TimeAsUnix, dmp.TimeFormat)
}

func Test_WithMaxDepth(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	opt := WithMaxDepth(10)

	// --- Then ---
	opt(dmp)
	affirm.Equal(t, 10, dmp.MaxDepth)
}

func Test_WithIndent(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	opt := WithIndent(10)

	// --- Then ---
	opt(dmp)
	affirm.Equal(t, 10, dmp.Indent)
}

func Test_WithTabWidth(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	opt := WithTabWidth(10)

	// --- Then ---
	opt(dmp)
	affirm.Equal(t, 10, dmp.TabWidth)
}

func Test_WithDumper(t *testing.T) {
	// --- Given ---
	dmp := Dump{Dumpers: make(map[reflect.Type]Dumper)}

	// --- When ---
	WithDumper(time.Time{}, GetTimeDumper(time.Kitchen))(&dmp)

	// --- Then ---
	have := dmp.Any(time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC))
	affirm.Equal(t, `"3:04AM"`, have)
}

func Test_New(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have := New()

		// --- Then ---
		affirm.Equal(t, false, have.Flat)
		affirm.Equal(t, 200, have.FlatStrings)
		affirm.Equal(t, false, have.Compact)
		affirm.Equal(t, TimeFormat, have.TimeFormat)
		affirm.Equal(t, "", have.DurationFormat)
		affirm.Equal(t, false, have.PtrAddr)
		affirm.Equal(t, true, have.UseAny)
		affirm.Equal(t, true, len(have.Dumpers) == 3)
		affirm.Equal(t, DefaultDepth, have.MaxDepth)
		affirm.Equal(t, DefaultIndent, have.Indent)
		affirm.Equal(t, DefaultTabWith, have.TabWidth)

		val, ok := have.Dumpers[typDur]
		affirm.Equal(t, true, ok)
		affirm.NotNil(t, val)

		val, ok = have.Dumpers[typLocation]
		affirm.Equal(t, true, ok)
		affirm.NotNil(t, val)

		val, ok = have.Dumpers[typTime]
		affirm.Equal(t, true, ok)
		affirm.NotNil(t, val)
	})
}

func Test_Dump_Any_Value_smoke_tabular(t *testing.T) {
	var itfNil types.TItf
	var itfVal, itfPtr types.TItf
	var sNil *types.TA
	itfVal = types.TVal{}
	itfPtr = &types.TPtr{}
	sPtr := &types.TPtr{Val: "a"}
	var aAnyNil any

	tt := []struct {
		testN string

		dmp  Dump
		v    any
		want string
	}{
		// Simple.
		{"bool true", New(WithFlat, WithCompact), true, "true"},
		{"int", New(WithFlat, WithCompact), 123, "123"},
		{"int8", New(WithFlat, WithCompact), int8(123), "123"},
		{"int16", New(WithFlat, WithCompact), int16(123), "123"},
		{"int32", New(WithFlat, WithCompact), int32(123), "123"},
		{"int64", New(WithFlat, WithCompact), int64(123), "123"},
		{"uint", New(WithFlat, WithCompact), uint(123), "123"},
		{"uint8", New(WithFlat, WithCompact), uint8(123), "0x7b"},
		{"byte", New(WithFlat, WithCompact), byte(123), "0x7b"},
		{"uint16", New(WithFlat, WithCompact), uint16(123), "123"},
		{"uint32", New(WithFlat, WithCompact), uint32(123), "123"},
		{"uint64", New(WithFlat, WithCompact), uint64(123), "123"},
		{"uintptr", New(WithFlat, WithCompact, WithPtrAddr), uintptr(123), "<0x7b>"},
		{"float32", New(WithFlat, WithCompact), float32(12.3), "12.3"},
		{"float64", New(WithFlat, WithCompact), 12.3, "12.3"},
		{"complex64", New(WithFlat, WithCompact), complex(float32(1), float32(2)), "(1+2i)"},
		{"complex128", New(WithFlat, WithCompact), complex(3.3, 4.4), "(3.3+4.4i)"},
		{"array", New(WithFlat, WithCompact), [2]int{}, "[2]int{0,0}"},
		{"chan", New(WithFlat, WithCompact), make(chan int), "(chan int)(<addr>)"},
		{"func", New(WithFlat, WithCompact), func() {}, "<func>(<addr>)"},
		{"interface nil", New(WithFlat, WithCompact), itfNil, ValNil},
		{"any nil", New(WithFlat, WithCompact), aAnyNil, ValNil},
		{"interface val", New(WithFlat, WithCompact), itfVal, `{Val:""}`},
		{"interface ptr", New(WithFlat, WithCompact), itfPtr, `{Val:""}`},
		{
			"map",
			New(WithFlat, WithCompact),
			map[string]string{"A": "a", "B": "b"},
			`map[string]string{"A":"a","B":"b"}`,
		},
		{"struct pointer", New(WithFlat, WithCompact), sPtr, `{Val:"a"}`},
		{"slice", New(WithFlat, WithCompact), []int{1, 2}, "[]int{1,2}"},
		{"string", New(WithFlat, WithCompact), "string", `"string"`},
		{"struct", New(WithFlat, WithCompact), struct{ F0 int }{}, "{F0:0}"},
		{
			"registered",
			New(WithFlat, WithCompact),
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			`"2000-01-02T03:04:05Z"`,
		},
		{"struct nil", New(WithFlat, WithCompact), sNil, "nil"},
		{
			"registered",
			New(WithFlat, WithCompact),
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			`"2000-01-02T03:04:05Z"`,
		},
		{
			"unsafe pointer",
			New(WithFlat, WithCompact, WithPtrAddr),
			unsafe.Pointer(sPtr),
			fmt.Sprintf("<%p>", sPtr),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveAny := tc.dmp.Any(tc.v)
			haveVal := tc.dmp.Value(reflect.ValueOf(tc.v))

			// --- Then ---
			affirm.Equal(t, tc.want, haveAny)
			affirm.Equal(t, tc.want, haveVal)
		})
	}
}

func Test_Dump_Any(t *testing.T) {
	t.Run("nil interface value", func(t *testing.T) {
		// --- Given ---
		var itfNil types.TItf
		dmp := New()

		// --- When ---
		have := dmp.Any(itfNil)

		// --- Then ---
		affirm.Equal(t, ValNil, have)
	})

	t.Run("slice of slices of any", func(t *testing.T) {
		// --- Given ---
		val := [][]any{
			{"str00", 0, "str02"},
			{"str10", 1, "str12"},
			{"str10", 1, nil},
		}
		dmp := New(WithFlat, WithCompact)

		// --- When ---
		have := dmp.Any(val)

		// --- Then ---
		want := `[][]any{{"str00",0,"str02"},{"str10",1,"str12"},{"str10",1,nil}}`
		affirm.Equal(t, want, have)
	})

	t.Run("depth", func(t *testing.T) {
		// --- Given ---
		val := struct {
			S0 struct {
				S1 struct {
					S2 struct {
						S4 struct {
							S5 struct {
								S6 struct{ VAL int }
							}
						}
					}
				}
			}
		}{}
		dmp := New(WithFlat, WithCompact)

		// --- When ---
		have := dmp.Any(val)

		// --- Then ---
		affirm.Equal(t, "{S0:{S1:{S2:{S4:{S5:{S6:{VAL:<...>}}}}}}}", have)
	})

	t.Run("format nested slices", func(t *testing.T) {
		// --- Given ---
		type Node struct {
			Value    int
			Children []*Node
		}

		val := &Node{
			Value: 1,
			Children: []*Node{
				{
					Value: 2,
				},
				{
					Value: 3,
					Children: []*Node{
						{
							Value: 4,
						},
					},
				},
			},
		}

		// --- When ---
		have := New().Any(val)

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_nested.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("format nested slices indented twice", func(t *testing.T) {
		// --- Given ---
		type Node struct {
			Value    int
			Children []*Node
		}

		val := &Node{
			Value: 1,
			Children: []*Node{
				{
					Value: 2,
				},
				{
					Value: 3,
					Children: []*Node{
						{
							Value: 4,
						},
					},
				},
			},
		}
		dmp := New(WithIndent(2))

		// --- When ---
		have := dmp.Any(val)

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_nested_with_indent.gld")
		affirm.Equal(t, want.String(), have)
	})
}

func Test_Dump_Diff_tabular(t *testing.T) {
	tt := []struct {
		testN string

		opts    []Option
		wantIn  any
		haveIn  any
		wantOut string
		haveOut string
		diffOut string
	}{
		{"same strings", nil, "abc", "abc", `"abc"`, `"abc"`, ""},
		{
			"same multiline strings",
			nil,
			"a\nb\nc",
			"a\nb\nc",
			`"a\nb\nc"`,
			`"a\nb\nc"`,
			"",
		},
		{
			"different multiline strings",
			nil,
			"a\nb\nc",
			"a\nx\nc",
			`"a\nb\nc"`,
			`"a\nx\nc"`,
			"" +
				"@@ -1,3 +1,3 @@\n" +
				" a\n" +
				"-x\n" +
				"+b\n" +
				" c",
		},
		{
			"want is single line have is multi line",
			nil,
			"abc",
			"a\nb\nc",
			`"abc"`,
			`"a\nb\nc"`,
			"" +
				"@@ -1,3 +1 @@\n" +
				"-a\n" +
				"-b\n" +
				"-c\n" +
				"+abc",
		},
		{
			"want is multiline have is single line",
			nil,
			"a\nb\nc",
			"abc",
			`"a\nb\nc"`,
			`"abc"`,
			"" +
				"@@ -1 +1,3 @@\n" +
				"-abc\n" +
				"+a\n" +
				"+b\n" +
				"+c",
		},
		{
			"both multiline strings end with a new line",
			nil,
			"a\nb\nc\n",
			"abc",
			`"a\nb\nc\n"`,
			`"abc"`,
			"" +
				"@@ -1 +1,3 @@\n" +
				"-abc\n" +
				"+a\n" +
				"+b\n" +
				"+c",
		},
		{
			"want is multiline then both should be",
			[]Option{WithFlatStrings(6)},
			"a\nb\nc\nd",
			"a\nb\nc\n",
			"a\nb\nc\nd",
			"a\nb\nc\n",
			"" +
				"@@ -2,2 +2,3 @@\n" +
				" b\n" +
				" c\n" +
				"+d",
		},
		{
			"have multiline then both should be",
			[]Option{WithFlatStrings(6)},
			"a\nb\nc\n",
			"a\nb\nc\nd",
			"a\nb\nc\n",
			"a\nb\nc\nd",
			"" +
				"@@ -2,3 +2,2 @@\n" +
				" b\n" +
				" c\n" +
				"-d",
		},
		{
			"both nil",
			[]Option{WithFlatStrings(6)},
			nil,
			nil,
			"nil",
			"nil",
			"",
		},
		{
			"want is nil",
			nil,
			nil,
			"a\nb\nc",
			"nil",
			`"a\nb\nc"`,
			"",
		},
		{
			"have is nil",
			nil,
			"a\nb\nc",
			nil,
			`"a\nb\nc"`,
			"nil",
			"",
		},
		{
			"both values are different and not multiline",
			nil,
			"abc",
			"xyz",
			`"abc"`,
			`"xyz"`,
			"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			dmp := New(tc.opts...)

			// --- When ---
			wantOut, haveOut, diffOut := dmp.Diff(tc.wantIn, tc.haveIn)

			// --- Then ---
			affirm.Equal(t, tc.wantOut, wantOut)
			affirm.Equal(t, tc.haveOut, haveOut)
			affirm.Equal(t, tc.diffOut, diffOut)
		})
	}
}

func Test_Dump_forDiff(t *testing.T) {
	t.Run("changes Flat and Compact configuration", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}
		dmp := Dump{
			Flat:        true,
			FlatStrings: 10,
			Compact:     true,
			MaxDepth:    Depth,
			Indent:      Indent,
			TabWidth:    TabWidth,
		}

		// --- When ---
		have, haveKnd := dmp.forDiff(val)

		// --- Then ---
		affirm.Equal(t, "{\n  1,\n  2,\n  3,\n}", have)
		affirm.Equal(t, reflect.Slice, haveKnd)
		affirm.Equal(t, true, dmp.Flat)
		affirm.Equal(t, 10, dmp.FlatStrings)
		affirm.Equal(t, true, dmp.Compact)
	})
}
