// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_ZoneDumper(t *testing.T) {
	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := ZoneDumper(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}

func Test_ZoneDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val    any
		indent int
		level  int
		want   string
	}{
		{
			"timezone",
			*testcases.WAW,
			0,
			0,
			`"Europe/Warsaw"`,
		},
		{
			"uses indent and level",
			*testcases.WAW,
			2,
			1,
			"      \"Europe/Warsaw\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			dmp := New(WithIndent(tc.indent))

			// --- When ---
			have := ZoneDumper(dmp, tc.level, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
