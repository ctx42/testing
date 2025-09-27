// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_HexPtrDumper(t *testing.T) {
	t.Run("Uintptr without addresses", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		val := uintptr(123)

		// --- When ---
		have := HexPtrDumper(dmp, 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.Equal(t, "<addr>", have)
	})

	t.Run("Uintptr with addresses", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithPtrAddr)
		val := uintptr(123)

		// --- When ---
		have := HexPtrDumper(dmp, 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.Equal(t, "<0x7b>", have)
	})

	t.Run("UnsafePointer without addresses", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		v := 42
		val := unsafe.Pointer(&v)

		// --- When ---
		have := HexPtrDumper(dmp, 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.Equal(t, "<addr>", have)
	})

	t.Run("UnsafePointer with addresses", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithPtrAddr)
		v := 42
		val := unsafe.Pointer(&v)

		// --- When ---
		have := HexPtrDumper(dmp, 0, reflect.ValueOf(val))

		// --- Then ---
		affirm.Equal(t, fmt.Sprintf("<%p>", &v), have)
	})
}

func Test_HexPtrDumper_tabular(t *testing.T) {
	sPtr := &testcases.TPtr{Val: "a"}

	tt := []struct {
		testN string

		val    any
		indent int
		level  int
		want   string
	}{
		{"uintptr", uintptr(1234), 0, 0, "<0x4d2>"},
		{"byte", byte(123), 0, 0, "0x7b"},
		{"usage error", 1234, 0, 0, ValErrUsage},
		{"unsafe pointer", unsafe.Pointer(sPtr), 0, 0, fmt.Sprintf("<%p>", sPtr)},

		{"uses indent", 1234, 2, 0, "    " + ValErrUsage},
		{"uses level", 1234, 0, 1, "  " + ValErrUsage},
		{"uses indent and level", 1234, 2, 1, "      " + ValErrUsage},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			dmp := New(WithIndent(tc.indent), WithPtrAddr)

			// --- When ---
			have := HexPtrDumper(dmp, tc.level, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
