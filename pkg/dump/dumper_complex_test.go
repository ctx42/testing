// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_ComplexDumper(t *testing.T) {
	t.Run("error - usage", func(t *testing.T) {
		// --- Given ---
		dmp := New()

		// --- When ---
		have := ComplexDumper(dmp, 0, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, ValErrUsage, have)
	})

	t.Run("invalid usage uses level", func(t *testing.T) {
		// --- Given ---
		dmp := New()

		// --- When ---
		have := ComplexDumper(dmp, 1, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "  "+ValErrUsage, have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))

		// --- When ---
		have := ComplexDumper(dmp, 1, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})

	t.Run("unexported field", func(t *testing.T) {
		// --- Given ---
		type hold struct{ z complex128 }
		val := reflect.ValueOf(hold{z: 1 + 2i}).Field(0)

		// --- When ---
		have := ComplexDumper(New(), 0, val)

		// --- Then ---
		affirm.Equal(t, "(1+2i)", have)
	})
}

func Test_ComplexDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val  any
		want string
	}{
		{"complex64", complex(float32(1.1), float32(2.2)), "(1.1+2.2i)"},
		{"complex128", complex(3.3, 4.4), "(3.3+4.4i)"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := ComplexDumper(Dump{}, 0, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
