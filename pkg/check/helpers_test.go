// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_typeString(t *testing.T) {
	tt := []struct {
		testN string

		val  reflect.Value
		want string
	}{
		{"string", reflect.ValueOf("abc"), "string"},
		{"int", reflect.ValueOf(123), "int"},
		{"invalid", reflect.ValueOf(nil), "<invalid>"},
		{"struct", reflect.ValueOf(types.TA{}), "types.TA"},
		{"ptr struct", reflect.ValueOf(&types.TA{}), "*types.TA"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := typeString(tc.val)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_isPrintableChar(t *testing.T) {
	for i := 0; i <= 31; i++ {
		if !affirm.Equal(t, false, isPrintableChar(byte(i))) {
			t.Logf("expected false for %d", i)
		}
	}
	for i := 32; i <= 126; i++ {
		if !affirm.Equal(t, true, isPrintableChar(byte(i))) {
			t.Logf("expected true for %d", i)
		}
	}
	for i := 127; i <= 255; i++ {
		if !affirm.Equal(t, false, isPrintableChar(byte(i))) {
			t.Logf("expected false for %d", i)
		}
	}
}

func Test_valToString_tabular(t *testing.T) {
	var itf, nilItf types.TItf
	itf = types.TVal{}
	var ptr, nilPtr *types.TPtr
	ptr = &types.TPtr{}

	tt := []struct {
		testN string

		key  any
		want string
	}{
		{"int", 1, "1"},
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int32(64), "64"},

		{"uint", 1, "1"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint32(64), "64"},

		{"uintptr", uintptr(123), "123"},

		{"float32", float32(1.1), "1.1"},
		{"float64", 1.2, "1.2"},

		{"string", "abc", `"abc"`},
		{"bool", true, "true"},

		{"struct", types.TA{}, "types.TA"},
		{"nil interface", nilItf, "<invalid>"},
		{"non-nil interface", itf, "types.TVal"},
		{"nil pointer", nilPtr, "<nil>"},
		{"non-nil pointer", ptr, "*types.TPtr"},

		{"complex64", complex(float32(1.0), float32(2.0)), "(1+2i)"},
		{"complex128", complex(3.0, 4.0), "(3+4i)"},
		{"array", [...]int{1, 2, 3}, "<array>"},
		{"chan", make(chan int), "<invalid>"},
		{"func", func() {}, "<invalid>"},
		{"map", map[string]int{"A": 1}, "<invalid>"},
		{"slice", []int{1, 2, 3}, "<invalid>"},
		{"unsafe pointer", unsafe.Pointer(ptr), fmt.Sprintf("<%p>", ptr)},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			key := reflect.ValueOf(tc.key)

			// --- When ---
			have := valToString(key)

			// --- Then ---
			if tc.want != have {
				format := "expected:\n\twant: %#v\n\thave: %#v"
				t.Errorf(format, tc.want, have)
			}
		})
	}
}
