// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_GetTimeDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		format string
		val    time.Time
		want   string
	}{
		{
			"empty format",
			"",
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			`"2000-01-02T03:04:05+01:00"`,
		},
		{
			"TimeAsRFC3339",
			TimeAsRFC3339,
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			`"2000-01-02T03:04:05+01:00"`,
		},
		{
			"TimeAsUnix",
			TimeAsUnix,
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			"946778645",
		},
		{
			"custom",
			time.TimeOnly,
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			`"03:04:05"`,
		},
		{
			"TimeAsGoString",
			TimeAsGoString,
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			`time.Date(2000, time.January, 2, 3, 4, 5, 0, time.Location("Europe/Warsaw"))`,
		},
		{
			"unsupported",
			"abc",
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			`"abc"`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			dmp := GetTimeDumper(tc.format)

			// --- Then ---
			have := dmp(Dump{}, 0, reflect.ValueOf(tc.val))
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_TimeDumperFmt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW)
		val := reflect.ValueOf(tim)
		dumper := TimeDumperFmt(time.DateOnly)

		// --- When ---
		have := dumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, `"2000-01-02"`, have)
	})

	t.Run("zero value", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		val := reflect.ValueOf(time.Time{})
		dumper := TimeDumperFmt(time.DateOnly)

		// --- When ---
		have := dumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, `"0001-01-01"`, have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW)
		val := reflect.ValueOf(tim)
		dumper := TimeDumperFmt(time.DateOnly)

		// --- When ---
		have := dumper(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "      \"2000-01-02\"", have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))
		dumper := TimeDumperFmt(time.DateOnly)

		// --- When ---
		have := dumper(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}

func Test_TimeDumperUnix(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperUnix(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "946778645", have)
	})

	t.Run("start of Unix epoch", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		tim := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperUnix(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "0", have)
	})

	t.Run("zero value", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		val := reflect.ValueOf(time.Time{})

		// --- When ---
		have := TimeDumperUnix(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "-62135596800", have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperUnix(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "      946778645", have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := TimeDumperUnix(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}

func Test_TimeDumperDate(t *testing.T) {
	t.Run("UTC", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		tim := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperDate(dmp, 0, val)

		// --- Then ---
		want := "time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC)"
		affirm.Equal(t, want, have)
	})

	t.Run("non UTC timezone", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperDate(dmp, 0, val)

		// --- Then ---
		want := "time.Date(2000, time.January, 2, 3, 4, 5, 0, " +
			"time.Location(\"Europe/Warsaw\"))"
		affirm.Equal(t, want, have)
	})

	t.Run("compact", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithCompact)
		tim := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperDate(dmp, 0, val)

		// --- Then ---
		want := "time.Date(2000,time.January,2,3,4,5,6,time.UTC)"
		affirm.Equal(t, want, have)
	})

	t.Run("start of Unix epoch", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		tim := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperDate(dmp, 0, val)

		// --- Then ---
		want := "time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)"
		affirm.Equal(t, want, have)
	})

	t.Run("zero value", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		val := reflect.ValueOf(time.Time{})

		// --- When ---
		have := TimeDumperDate(dmp, 0, val)

		// --- Then ---
		want := "time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)"
		affirm.Equal(t, want, have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		tim := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		val := reflect.ValueOf(tim)

		// --- When ---
		have := TimeDumperDate(dmp, 1, val)

		// --- Then ---
		want := "      time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)"
		affirm.Equal(t, want, have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := TimeDumperDate(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}
