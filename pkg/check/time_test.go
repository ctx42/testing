// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Time(t *testing.T) {
	t.Run("equal both time Time", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, testcases.WAW)

		// --- When ---
		err := Time(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal both time Time", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 6, 0, testcases.WAW)

		// --- When ---
		err := Time(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  want: 2000-01-02T03:04:05Z\n" +
			"  have: 2000-01-02T03:04:06Z ( 2000-01-02T04:04:06+01:00 )\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal want is string", func(t *testing.T) {
		// --- Given ---
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Time("2000-01-02T04:04:05+01:00", have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal want is string", func(t *testing.T) {
		// --- Given ---
		have := time.Date(2000, 1, 2, 2, 4, 4, 0, time.UTC)

		// --- When ---
		err := Time("2000-01-02T03:04:05+01:00", have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  want: 2000-01-02T02:04:05Z ( 2000-01-02T03:04:05+01:00 )\n" +
			"  have: 2000-01-02T02:04:04Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not equal have is string", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 2, 4, 4, 0, time.UTC)

		// --- When ---
		err := Time(want, "2000-01-02T03:04:05+01:00")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  want: 2000-01-02T02:04:04Z\n" +
			"  have: 2000-01-02T02:04:05Z ( 2000-01-02T03:04:05+01:00 )\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid want date format", func(t *testing.T) {
		// --- When ---
		err := Time("2022-02-18", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[want] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Time(time.Now(), "2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid option date format", func(t *testing.T) {
		// --- When ---
		err := Time("2022-02-18", time.Now(), WithTimeFormat("abc"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[want] failed to parse time:\n" +
			"  format: abc\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 6, 0, testcases.WAW)

		// --- When ---
		err := Time(want, have, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  trail: type.field\n" +
			"   want: 2000-01-02T03:04:05Z\n" +
			"   have: 2000-01-02T03:04:06Z ( 2000-01-02T04:04:06+01:00 )\n" +
			"   diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Exact(t *testing.T) {
	t.Run("error - same instant numeric offsets", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.FixedZone("", 3600))
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, time.FixedZone("", 7200))

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  want: +01:00\n" +
			"  have: +02:00"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("exactly", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not exact date", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  want: 2000-01-02T03:04:05Z\n" +
			"  have: 2000-01-02T03:04:06Z\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, false, want.Equal(have))
	})

	t.Run("error - not exact date want is string", func(t *testing.T) {
		// --- Given ---
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Exact("2000-01-02T03:04:05Z", have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  want: 2000-01-02T03:04:05Z\n" +
			"  have: 2000-01-02T03:04:06Z\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not exact date have is string", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Exact(want, "2000-01-02T03:04:05Z")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  want: 2000-01-02T03:04:06Z\n" +
			"  have: 2000-01-02T03:04:05Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not exact timezone", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, testcases.WAW)

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  want: UTC\n" +
			"  have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid want date format", func(t *testing.T) {
		// --- When ---
		err := Exact("2022-02-18", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[want] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Exact(time.Now(), "2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Exact(want, have, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal dates:\n" +
			"  trail: type.field\n" +
			"   want: 2000-01-02T03:04:05Z\n" +
			"   have: 2000-01-02T03:04:06Z\n" +
			"   diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, false, want.Equal(have))
	})
}

func Test_Before(t *testing.T) {
	t.Run("before", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)

		// --- When ---
		err := Before(mark, date)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - equal", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Before(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be before mark:\n" +
			"  date: 2000-01-02T03:04:05Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - after", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Before(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be before mark:\n" +
			"  date: 2000-01-02T03:04:06Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid date", func(t *testing.T) {
		// --- When ---
		err := Before(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid mark", func(t *testing.T) {
		// --- When ---
		err := Before("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Before(mark, date, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be before mark:\n" +
			"  trail: type.field\n" +
			"   date: 2000-01-02T03:04:06Z\n" +
			"   mark: 2000-01-02T03:04:05Z\n" +
			"   diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_After(t *testing.T) {
	t.Run("after", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(mark, date)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - equal", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be after mark:\n" +
			"  date: 2000-01-02T03:04:05Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - before", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be after mark:\n" +
			"  date: 2000-01-02T03:04:05Z\n" +
			"  mark: 2001-01-02T03:04:05Z\n" +
			"  diff: -8784h0m0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid date", func(t *testing.T) {
		// --- When ---
		err := After(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid mark", func(t *testing.T) {
		// --- When ---
		err := After("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(mark, date, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be after mark:\n" +
			"  trail: type.field\n" +
			"   date: 2000-01-02T03:04:05Z\n" +
			"   mark: 2001-01-02T03:04:05Z\n" +
			"   diff: -8784h0m0s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_BeforeOrEqual(t *testing.T) {
	t.Run("before", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(mark, date)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(mark, date)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - after", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be equal or before mark:\n" +
			"  date: 2000-01-02T03:04:06Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid date", func(t *testing.T) {
		// --- When ---
		err := BeforeOrEqual(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid mark", func(t *testing.T) {
		// --- When ---
		err := BeforeOrEqual("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(mark, date, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be equal or before mark:\n" +
			"  trail: type.field\n" +
			"   date: 2000-01-02T03:04:06Z\n" +
			"   mark: 2000-01-02T03:04:05Z\n" +
			"   diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_AfterOrEqual(t *testing.T) {
	t.Run("after", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(mark, date)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(mark, date)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - before", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be equal or after mark:\n" +
			"  date: 2000-01-02T03:04:04Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid date", func(t *testing.T) {
		// --- When ---
		err := AfterOrEqual(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid mark", func(t *testing.T) {
		// --- When ---
		err := AfterOrEqual("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(mark, date, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected date to be equal or after mark:\n" +
			"  trail: type.field\n" +
			"   date: 2000-01-02T03:04:04Z\n" +
			"   mark: 2000-01-02T03:04:05Z\n" +
			"   diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Within(t *testing.T) {
	t.Run("within ahead", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Within(want, "1s", have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("within behind", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)

		// --- When ---
		err := Within(want, "1s", have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not within", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Within(want, "1000s", have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2001-01-02T03:04:05Z\n" +
			"  max diff +/-: 1000s\n" +
			"     have diff: -8784h0m0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid want date format", func(t *testing.T) {
		// --- When ---
		err := Within("2022-02-18", "1s", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[want] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid within format", func(t *testing.T) {
		// --- When ---
		err := Within(time.Now(), "abc", "2000-01-02T03:04:05Z")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[within] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Within(time.Now(), "1s", "2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 5e8, time.UTC)

		// --- When ---
		err := Within(want, "1s", have, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected dates to be within:\n" +
			"         trail: type.field\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:04:06.5Z\n" +
			"  max diff +/-: 1s\n" +
			"     have diff: -1.5s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_WithinChecker(t *testing.T) {
	t.Run("it is within", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := WithinChecker(time.Second)(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - it is not within", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 7, 0, time.UTC)

		// --- When ---
		err := WithinChecker(time.Second)(want, have)

		// --- Then ---
		wMsg := "" +
			"expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:04:07Z\n" +
			"  max diff +/-: 1s\n" +
			"     have diff: -2s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Recent(t *testing.T) {
	t.Run("recent with time now", func(t *testing.T) {
		// --- Given ---
		have := time.Now()

		// --- When ---
		err := Recent(have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("recent in the future", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := now.Add(9 * time.Second)

		// --- When ---
		err := Recent(have, WithNow(func() time.Time { return now }))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("recent in the past", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := now.Add(-9 * time.Second)

		// --- When ---
		err := Recent(have, WithNow(func() time.Time { return now }))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not recent in the future", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := now.Add(11 * time.Second)

		// --- When ---
		err := Recent(have, WithNow(func() time.Time { return now }))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:04:16Z\n" +
			"  max diff +/-: 10s\n" +
			"     have diff: -11s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not recent in the past", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := now.Add(-11 * time.Second)

		// --- When ---
		err := Recent(have, WithNow(func() time.Time { return now }))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:03:54Z\n" +
			"  max diff +/-: 10s\n" +
			"     have diff: 11s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Recent("2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := now.Add(-11 * time.Second)
		opts := []any{
			WithNow(func() time.Time { return now }),
			WithTrail("type.field"),
		}

		// --- When ---
		err := Recent(have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected dates to be within:\n" +
			"         trail: type.field\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:03:54Z\n" +
			"  max diff +/-: 10s\n" +
			"     have diff: 11s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Zone(t *testing.T) {
	t.Run("equal instance", func(t *testing.T) {
		// --- When ---
		err := Zone(time.UTC, time.UTC)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal string", func(t *testing.T) {
		// --- When ---
		err := Zone("UTC", "UTC")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("nil want is UTC", func(t *testing.T) {
		// --- When ---
		err := Zone(nil, time.UTC)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("nil have is UTC", func(t *testing.T) {
		// --- When ---
		err := Zone(time.UTC, nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal instance", func(t *testing.T) {
		// --- When ---
		err := Zone(time.UTC, testcases.WAW)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  want: UTC\n" +
			"  have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not equal string", func(t *testing.T) {
		// --- When ---
		err := Zone("UTC", "Europe/Warsaw")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  want: UTC\n" +
			"  have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - want is unsupported timezone type", func(t *testing.T) {
		// --- When ---
		err := Zone("123", "Europe/Warsaw")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[want] failed to parse timezone:\n" +
			"  value: 123"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - have is unsupported timezone type", func(t *testing.T) {
		// --- When ---
		err := Zone("Europe/Warsaw", "123")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"[have] failed to parse timezone:\n" +
			"  value: 123"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- When ---
		err := Zone(time.UTC, testcases.WAW, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  trail: type.field\n" +
			"   want: UTC\n" +
			"   have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Duration(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- When ---
		err := Duration(time.Second, time.Second)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- When ---
		err := Duration("1000s", "2000s")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal time durations:\n" +
			"  want: 1000s\n" +
			"  have: 2000s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid want", func(t *testing.T) {
		// --- When ---
		err := Duration("abc", "2000s")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[want] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid have", func(t *testing.T) {
		// --- When ---
		err := Duration("2000s", "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[have] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- When ---
		err := Duration(time.Second, 2*time.Second, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected equal time durations:\n" +
			"  trail: type.field\n" +
			"   want: 1s\n" +
			"   have: 2s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_formatDates(t *testing.T) {
	t.Run("same format", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2000, 1, 2, 3, 4, 5, 0, testcases.WAW)
		wTimStr := "2000-01-02T03:04:05+02:00"
		hTim := time.Date(2001, 1, 2, 3, 4, 5, 0, testcases.WAW)
		hTimStr := "2001-01-02T03:04:05+02:00"

		// --- When ---
		hW, hH := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		want := "2000-01-02T02:04:05Z ( 2000-01-02T03:04:05+02:00 )"
		affirm.Equal(t, want, hW)
		want = "2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+02:00 )"
		affirm.Equal(t, want, hH)
	})

	t.Run("shorted format when both dates in UTC", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		wTimStr := "2000-01-02T03:04:05Z"
		hTim := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		hTimStr := "2001-01-02T03:04:05Z"

		// --- When ---
		hW, hH := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		affirm.Equal(t, "2000-01-02T03:04:05Z", hW)
		affirm.Equal(t, "2001-01-02T03:04:05Z", hH)
	})

	t.Run("want date string is longer", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2001, 1, 2, 3, 4, 5, 0, testcases.WAW)
		wTimStr := "2001-01-02T03:04:05+01:00"
		hTim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		hTimStr := "2000-01-02T03:04:05Z"

		// --- When ---
		hW, hH := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		want := "2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+01:00 )"
		affirm.Equal(t, want, hW)
		affirm.Equal(t, "2000-01-02T03:04:05Z", hH)
	})

	t.Run("have date string is longer", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		wTimStr := "2000-01-02T03:04:05Z"
		hTim := time.Date(2001, 1, 2, 3, 4, 5, 0, testcases.WAW)
		hTimStr := "2001-01-02T03:04:05+01:00"

		// --- When ---
		hW, hH := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		affirm.Equal(t, "2000-01-02T03:04:05Z", hW)
		want := "2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+01:00 )"
		affirm.Equal(t, want, hH)
	})
}

func Test_getTime(t *testing.T) {
	t.Run("error - wrong time format", func(t *testing.T) {
		// --- Given ---
		opt := WithTimeFormat(time.RFC3339)

		// --- When ---
		hTim, hStr, hRep, err := getTime("2000-01-02", opt)

		// --- Then ---
		affirm.Equal(t, true, hTim.IsZero())
		affirm.Equal(t, "2000-01-02", hStr)
		affirm.Equal(t, timeTypeStr, hRep)
		wMsg := "" +
			"failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05Z07:00\n" +
			"   value: 2000-01-02"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeParse))
	})

	t.Run("error - empty option time format", func(t *testing.T) {
		// --- When ---
		have, hStr, hRep, err := getTime("2000-01-02", WithTimeFormat(""))

		// --- Then ---
		affirm.Equal(t, true, have.IsZero())
		affirm.Equal(t, "2000-01-02", hStr)
		affirm.Equal(t, timeTypeStr, hRep)
		wMsg := "" +
			"failed to parse time:\n" +
			"  format: \n" +
			"   value: 2000-01-02\n" +
			"   error: extra text: \"2000-01-02\""
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeParse))
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- Given ---
		opts := []any{
			WithTimeFormat(time.RFC3339),
			WithTrail("type.field"),
		}

		// --- When ---
		have, hStr, hRep, err := getTime("2000-01-02", opts...)

		// --- Then ---
		affirm.Equal(t, true, have.IsZero())
		affirm.Equal(t, "2000-01-02", hStr)
		affirm.Equal(t, timeTypeStr, hRep)
		wMsg := "" +
			"failed to parse time:\n" +
			"   trail: type.field\n" +
			"  format: 2006-01-02T15:04:05Z07:00\n" +
			"   value: 2000-01-02"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - str-ts format with non-numeric value", func(t *testing.T) {
		// --- Given ---
		opt := WithTimeFormat(TimeFormatUnixStr)

		// --- When ---
		hTim, hStr, hRep, err := getTime("not-a-number", opt)

		// --- Then ---
		affirm.Equal(t, true, hTim.IsZero())
		affirm.Equal(t, "not-a-number", hStr)
		affirm.Equal(t, timeTypeStr, hRep)
		wMsg := "" +
			"failed to parse time:\n" +
			"  format: str-ts\n" +
			"   value: not-a-number"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeParse))
	})

	t.Run("error - unsupported type", func(t *testing.T) {
		// --- When ---
		have, hStr, hRep, err := getTime(true)

		// --- Then ---
		affirm.Equal(t, true, have.IsZero())
		affirm.Equal(t, "true", hStr)
		affirm.Equal(t, "", hRep)
		wMsg := "" +
			"failed to parse time:\n" +
			"  cause: not supported time type\n" +
			"   type: bool"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeType))
	})
}

func Test_getTime_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		opts []any
		in   any
		wStr string
		wRep timeRep
		wTim time.Time
		wTZ  *time.Location
	}{
		{
			"time Time in UTC",
			nil,
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"2000-01-02T03:04:05Z",
			timeTypeTim,
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.UTC,
		},
		{
			"time Time in UTC - apply different timezone",
			[]any{WithZone(testcases.WAW)},
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"2000-01-02T04:04:05+01:00",
			timeTypeTim,
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			testcases.WAW,
		},
		{
			"time Time in WAW",
			nil,
			time.Date(2000, 1, 2, 3, 4, 5, 0, testcases.WAW),
			"2000-01-02T03:04:05+01:00",
			timeTypeTim,
			time.Date(2000, 1, 2, 3, 4, 5, 0, testcases.WAW),
			testcases.WAW,
		},
		{
			"RFC3339",
			nil,
			"2000-01-02T03:04:05+01:00",
			"2000-01-02T03:04:05+01:00",
			timeTypeStr,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			time.UTC,
		},
		{
			"RFC3339 - apply timezone",
			[]any{WithZone(testcases.WAW)},
			"2000-01-02T03:04:05+01:00",
			"2000-01-02T03:04:05+01:00",
			timeTypeStr,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			testcases.WAW,
		},
		{
			"Unix timestamp int",
			nil,
			946778645,
			"946778645",
			timeTypeInt,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			time.UTC,
		},
		{
			"Unix timestamp int - apply timezone",
			[]any{WithZone(testcases.WAW)},
			946778645,
			"946778645",
			timeTypeInt,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			testcases.WAW,
		},
		{
			"Unix timestamp int64",
			nil,
			int64(946778645),
			"946778645",
			timeTypeInt64,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			time.UTC,
		},
		{
			"Unix timestamp int64 - apply timezone",
			[]any{WithZone(testcases.WAW)},
			int64(946778645),
			"946778645",
			timeTypeInt64,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			testcases.WAW,
		},
		{
			"Unix timestamp string",
			[]any{WithTimeFormat(TimeFormatUnixStr)},
			"946778645",
			"946778645",
			timeTypeStr,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			time.UTC,
		},
		{
			"Unix timestamp string - apply timezone",
			[]any{WithTimeFormat(TimeFormatUnixStr), WithZone(testcases.WAW)},
			"946778645",
			"946778645",
			timeTypeStr,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			testcases.WAW,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			hTim, hStr, hRep, err := getTime(tc.in, tc.opts...)

			// --- Then ---
			affirm.Nil(t, err)
			affirm.Equal(t, tc.wStr, hStr)
			affirm.Equal(t, tc.wRep, hRep)
			affirm.Equal(t, true, tc.wTim.Equal(hTim))
			affirm.Equal(t, tc.wTZ.String(), hTim.Location().String())
		})
	}
}

func Test_getZone(t *testing.T) {
	t.Run("timezone as string", func(t *testing.T) {
		// --- When ---
		hZone, hStr, hRep, err := getZone("Europe/Warsaw")

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, zoneString, hRep)
		affirm.Equal(t, "Europe/Warsaw", hStr)
		affirm.Equal(t, testcases.WAW.String(), hZone.String())
	})

	t.Run("timezone pointer instance", func(t *testing.T) {
		// --- When ---
		hZone, hStr, hRep, err := getZone(testcases.WAW)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, zoneZone, hRep)
		affirm.Equal(t, "Europe/Warsaw", hStr)
		affirm.Equal(t, testcases.WAW.String(), hZone.String())
	})

	t.Run("timezone instance", func(t *testing.T) {
		// --- When ---
		hZone, hStr, hRep, err := getZone(*testcases.WAW)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, zoneZone, hRep)
		affirm.Equal(t, "Europe/Warsaw", hStr)
		affirm.Equal(t, testcases.WAW.String(), hZone.String())
	})

	t.Run("error - parsing error", func(t *testing.T) {
		// --- When ---
		have, hStr, hRep, err := getZone("abc///")

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, "abc///", hStr)
		affirm.Equal(t, zoneString, hRep)
		wMsg := "" +
			"failed to parse timezone:\n" +
			"  value: abc///"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrZoneParse))
	})

	t.Run("error - unsupported type", func(t *testing.T) {
		// --- When ---
		have, hStr, hRep, err := getZone(123)

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, "123", hStr)
		affirm.Equal(t, "", hRep)
		wMsg := "" +
			"failed to parse timezone:\n" +
			"  cause: not supported timezone type\n" +
			"   type: int"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrZoneType))
	})
}

func Test_getDur(t *testing.T) {
	t.Run("error - invalid string duration", func(t *testing.T) {
		// --- When ---
		hDur, hStr, hRep, err := getDur("abc")

		// --- Then ---
		affirm.Equal(t, time.Duration(0), hDur)
		affirm.Equal(t, "abc", hStr)
		affirm.Equal(t, durTypeStr, hRep)
		affirm.Equal(t, "failed to parse duration:\n  value: abc", err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrDurParse))
	})

	t.Run("error - unsupported type", func(t *testing.T) {
		// --- When ---
		hDur, hStr, hRep, err := getDur(true)

		// --- Then ---
		affirm.Equal(t, time.Duration(0), hDur)
		affirm.Equal(t, "true", hStr)
		affirm.Equal(t, "", hRep)
		wMsg := "" +
			"failed to parse duration:\n" +
			"  cause: not supported duration type\n" +
			"   type: bool"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrDurType))
	})

	t.Run("error - log message with trail", func(t *testing.T) {
		// --- When ---
		hDur, hStr, hRep, err := getDur(true, WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, time.Duration(0), hDur)
		affirm.Equal(t, "true", hStr)
		affirm.Equal(t, "", hRep)
		wMsg := "" +
			"failed to parse duration:\n" +
			"  trail: type.field\n" +
			"  cause: not supported duration type\n" +
			"   type: bool"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrDurType))
	})
}

func Test_getDur_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		in   any
		wStr string
		wRep durRep
		wDur time.Duration
	}{
		{"time Duration", time.Second, "1s", durTypeDur, time.Second},
		{"time Duration as string", "1s", "1s", durTypeStr, time.Second},
		{
			"time Duration as int",
			12345678,
			"12345678",
			durTypeInt,
			time.Duration(12345678),
		},
		{
			"time Duration as int64",
			int64(12345678),
			"12345678",
			durTypeInt64,
			time.Duration(12345678),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			hDur, hStr, hRep, err := getDur(tc.in)

			// --- Then ---
			affirm.Nil(t, err)
			affirm.Equal(t, tc.wDur, hDur)
			affirm.Equal(t, tc.wStr, hStr)
			affirm.Equal(t, tc.wRep, hRep)
		})
	}
}

// Benchmarks for time-based checks (Recent, Within, Before/After hot paths).

func Benchmark_Recent(b *testing.B) {
	now := time.Now()
	recent := 5 * time.Second

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Recent(now.Add(-time.Second), WithRecent(recent))
	}
}

func Benchmark_Within_Time(b *testing.B) {
	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Within(t1, "3s", t2)
	}
}
