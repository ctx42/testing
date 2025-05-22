// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_complexDumper_tabular(t *testing.T) {
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
			have := complexDumper(Dump{}, 0, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_complexDumper(t *testing.T) {
	t.Run("invalid usage", func(t *testing.T) {
		// --- Given ---
		dmp := New()

		// --- When ---
		have := complexDumper(dmp, 0, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, ValErrUsage, have)
	})

	t.Run("invalid usage uses level", func(t *testing.T) {
		// --- Given ---
		dmp := New()

		// --- When ---
		have := complexDumper(dmp, 1, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "  "+ValErrUsage, have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))

		// --- When ---
		have := complexDumper(dmp, 1, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}
