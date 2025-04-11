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
	"github.com/ctx42/testing/internal/tstkit"
	"github.com/ctx42/testing/internal/types"
)

func Test_WithFlat(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithFlat(dmp)

	// --- Then ---
	affirm.True(t, dmp.Flat)
}

func Test_WithFlatStrings(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithFlatStrings(123)(dmp)

	// --- Then ---
	affirm.Equal(t, 123, dmp.FlatStings)
}

func Test_WithCompact(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithCompact(dmp)

	// --- Then ---
	affirm.True(t, dmp.Compact)
}

func Test_WithPtrAddr(t *testing.T) {
	// --- Given ---
	dmp := &Dump{}

	// --- When ---
	WithPtrAddr(dmp)

	// --- Then ---
	affirm.True(t, dmp.PtrAddr)
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
		affirm.False(t, have.Flat)
		affirm.Equal(t, 200, have.FlatStings)
		affirm.False(t, have.Compact)
		affirm.Equal(t, TimeFormat, have.TimeFormat)
		affirm.Equal(t, "", have.DurationFormat)
		affirm.False(t, have.PtrAddr)
		affirm.True(t, have.UseAny)
		affirm.True(t, len(have.Dumpers) == 3)
		affirm.Equal(t, DefaultDepth, have.MaxDepth)
		affirm.Equal(t, DefaultIndent, have.Indent)
		affirm.Equal(t, DefaultTabWith, have.TabWidth)

		val, ok := have.Dumpers[typDur]
		affirm.True(t, ok)
		affirm.NotNil(t, val)

		val, ok = have.Dumpers[typLocation]
		affirm.True(t, ok)
		affirm.NotNil(t, val)

		val, ok = have.Dumpers[typTime]
		affirm.True(t, ok)
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
		{"uintptr", New(WithFlat, WithCompact), uintptr(123), "<0x7b>"},
		{"float32", New(WithFlat, WithCompact), float32(12.3), "12.3"},
		{"float64", New(WithFlat, WithCompact), 12.3, "12.3"},
		{"complex64", New(WithFlat, WithCompact), complex(float32(1), float32(2)), "(1+2i)"},
		{"complex128", New(WithFlat, WithCompact), complex(3.3, 4.4), "(3.3+4.4i)"},
		{"array", New(WithFlat, WithCompact), [2]int{}, "[2]int{0,0}"},
		{"chan", New(WithFlat, WithCompact), make(chan int), "(chan int)(<addr>)"},
		{"func", New(WithFlat, WithCompact), func() {}, "<func>(<addr>)"},
		{"interface nil", New(WithFlat, WithCompact), itfNil, valNil},
		{"any nil", New(WithFlat, WithCompact), aAnyNil, valNil},
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
			New(WithFlat, WithCompact),
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
		affirm.Equal(t, valNil, have)
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
		want := tstkit.Golden(t, "testdata/struct_nested.txt")
		affirm.Equal(t, want, have)
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
		want := tstkit.Golden(t, "testdata/struct_nested_with_indent.txt")
		affirm.Equal(t, want, have)
	})
}
