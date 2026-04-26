// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_valueCmp_tabular(t *testing.T) {
	tt := []struct {
		testN string

		a    any
		b    any
		want int
	}{
		{"bool false & false", false, false, 0},
		{"bool true & true", true, true, 0},
		{"bool true & false", true, false, 1},
		{"bool false & true", false, true, -1},

		{"int 0 & 0", 0, 0, 0},
		{"int 1 & 1", 0, 0, 0},
		{"int 2 & 1", 2, 1, 1},
		{"int 1 & 2", 1, 2, -1},

		{"int8 0 & 0", int8(0), int8(0), 0},
		{"int8 1 & 1", int8(0), int8(0), 0},
		{"int8 2 & 1", int8(2), int8(1), 1},
		{"int8 1 & 2", int8(1), int8(2), -1},

		{"int16 0 & 0", int16(0), int16(0), 0},
		{"int16 1 & 1", int16(0), int16(0), 0},
		{"int16 2 & 1", int16(2), int16(1), 1},
		{"int16 1 & 2", int16(1), int16(2), -1},

		{"int32 0 & 0", int32(0), int32(0), 0},
		{"int32 1 & 1", int32(0), int32(0), 0},
		{"int32 2 & 1", int32(2), int32(1), 1},
		{"int32 1 & 2", int32(1), int32(2), -1},

		{"int64 0 & 0", int64(0), int64(0), 0},
		{"int64 1 & 1", int64(0), int64(0), 0},
		{"int64 2 & 1", int64(2), int64(1), 1},
		{"int64 1 & 2", int64(1), int64(2), -1},

		{"uint 0 & 0", uint(0), uint(0), 0},
		{"uint 1 & 1", uint(0), uint(0), 0},
		{"uint 2 & 1", uint(2), uint(1), 1},
		{"uint 1 & 2", uint(1), uint(2), -1},

		{"uint8 0 & 0", uint8(0), uint8(0), 0},
		{"uint8 1 & 1", uint8(0), uint8(0), 0},
		{"uint8 2 & 1", uint8(2), uint8(1), 1},
		{"uint8 1 & 2", uint8(1), uint8(2), -1},

		{"uint16 0 & 0", uint16(0), uint16(0), 0},
		{"uint16 1 & 1", uint16(0), uint16(0), 0},
		{"uint16 2 & 1", uint16(2), uint16(1), 1},
		{"uint16 1 & 2", uint16(1), uint16(2), -1},

		{"uint32 0 & 0", uint32(0), uint32(0), 0},
		{"uint32 1 & 1", uint32(0), uint32(0), 0},
		{"uint32 2 & 1", uint32(2), uint32(1), 1},
		{"uint32 1 & 2", uint32(1), uint32(2), -1},

		{"uint64 0 & 0", uint64(0), uint64(0), 0},
		{"uint64 1 & 1", uint64(0), uint64(0), 0},
		{"uint64 2 & 1", uint64(2), uint64(1), 1},
		{"uint64 1 & 2", uint64(1), uint64(2), -1},

		{"uintptr 0 & 0", uintptr(0), uintptr(0), 0},
		{"uintptr 1 & 1", uintptr(0), uintptr(0), 0},
		{"uintptr 2 & 1", uintptr(2), uintptr(1), 1},
		{"uintptr 1 & 2", uintptr(1), uintptr(2), -1},

		{"float32 0 & 0", float32(0.0), float32(0.0), 0},
		{"float32 1 & 1", float32(0.0), float32(0.0), 0},
		{"float32 2 & 1", float32(2.0), float32(1.0), 1},
		{"float32 1 & 2", float32(1.0), float32(2.0), -1},

		{"float64 0 & 0", 0.0, 0.0, 0},
		{"float64 1 & 1", 0.0, 0.0, 0},
		{"float64 2 & 1", 2.0, 1.0, 1},
		{"float64 1 & 2", 1.0, 2.0, -1},

		{"string empty & empty", "", "", 0},
		{"string abc & abc", "abc", "abc", 0},
		{"string xyz & abc", "xyz", "abc", 1},
		{"string abc & xyz", "abc", "xyz", -1},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			a := reflect.ValueOf(tc.a)
			b := reflect.ValueOf(tc.b)

			// --- When ---
			have := valueCmp(a, b)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
