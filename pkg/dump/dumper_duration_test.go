// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_GetDurDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		format string
		val    time.Duration
		want   string
	}{
		{"empty format", "", 1500 * time.Millisecond, `"1.5s"`},
		{"DurAsString", DurAsString, 1500 * time.Millisecond, `"1.5s"`},
		{"DurAsSeconds", DurAsSeconds, 1500 * time.Millisecond, "1.5"},
		{"unsupported", "abc", 1500 * time.Millisecond, `"1.5s"`},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			dmp := New()

			// --- When ---
			dumper := GetDurDumper(tc.format)
			have := dumper(dmp, 0, reflect.ValueOf(tc.val))

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_DurDumperString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		dur := 1500 * time.Millisecond
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperString(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, `"1.5s"`, have)
	})

	t.Run("zero value", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		dur := 0 * time.Second
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperString(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, `"0s"`, have)
	})

	t.Run("uses level", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		dur := time.Second
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperString(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "  \"1s\"", have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		dur := time.Second
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperString(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "      \"1s\"", have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := DurDumperString(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, ValErrUsage, have)
	})
}

func Test_DurDumperSeconds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		dur := 1500 * time.Millisecond
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperSeconds(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "1.5", have)
	})

	t.Run("zero value", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		dur := 0 * time.Second
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperSeconds(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "0", have)
	})

	t.Run("uses level", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		dur := time.Second
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperSeconds(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "  1", have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		dur := time.Second
		val := reflect.ValueOf(dur)

		// --- When ---
		have := DurDumperSeconds(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "      1", have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := DurDumperSeconds(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, ValErrUsage, have)
	})
}
