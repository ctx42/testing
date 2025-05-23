// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_Time(t *testing.T) {
	t.Run("equal both time.Time", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, types.WAW)

		// --- When ---
		err := Time(want, have)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, true, want.Equal(have))
	})

	t.Run("not equal both time.Time", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 6, 0, types.WAW)

		// --- When ---
		err := Time(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal dates:\n" +
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

	t.Run("not equal want is string", func(t *testing.T) {
		// --- Given ---
		have := time.Date(2000, 1, 2, 2, 4, 4, 0, time.UTC)

		// --- When ---
		err := Time("2000-01-02T03:04:05+01:00", have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal dates:\n" +
			"  want: 2000-01-02T02:04:05Z ( 2000-01-02T03:04:05+01:00 )\n" +
			"  have: 2000-01-02T02:04:04Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("not equal have is string", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 2, 4, 4, 0, time.UTC)

		// --- When ---
		err := Time(want, "2000-01-02T03:04:05+01:00")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal dates:\n" +
			"  want: 2000-01-02T02:04:04Z\n" +
			"  have: 2000-01-02T02:04:05Z ( 2000-01-02T03:04:05+01:00 )\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid want date format", func(t *testing.T) {
		// --- When ---
		err := Time("2022-02-18", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[want] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Time(time.Now(), "2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid option date format", func(t *testing.T) {
		// --- Given ---
		opt := WithTimeFormat("abc")

		// --- When ---
		err := Time("2022-02-18", time.Now(), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[want] failed to parse time:\n" +
			"  format: abc\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 6, 0, types.WAW)
		opt := WithTrail("type.field")

		// --- When ---
		err := Time(want, have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal dates:\n" +
			"  trail: type.field\n" +
			"   want: 2000-01-02T03:04:05Z\n" +
			"   have: 2000-01-02T03:04:06Z ( 2000-01-02T04:04:06+01:00 )\n" +
			"   diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Exact(t *testing.T) {
	t.Run("exactly", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, true, want.Equal(have))
	})

	t.Run("error - not exact date", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal dates:\n" +
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
		wMsg := "expected equal dates:\n" +
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
		wMsg := "expected equal dates:\n" +
			"  want: 2000-01-02T03:04:06Z\n" +
			"  have: 2000-01-02T03:04:05Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not exact timezone", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, types.WAW)

		// --- When ---
		err := Exact(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same timezone:\n" +
			"  want: UTC\n" +
			"  have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, want.Equal(have))
	})

	t.Run("invalid want date format", func(t *testing.T) {
		// --- When ---
		err := Exact("2022-02-18", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[want] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Exact(time.Now(), "2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		opt := WithTrail("type.field")

		// --- When ---
		err := Exact(want, have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal dates:\n" +
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
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Before(date, mark)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Before(date, mark)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be before mark:\n" +
			"  date: 2000-01-02T03:04:05Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("after", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Before(date, mark)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be before mark:\n" +
			"  date: 2000-01-02T03:04:06Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid date", func(t *testing.T) {
		// --- When ---
		err := Before("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid mark", func(t *testing.T) {
		// --- When ---
		err := Before(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithTrail("type.field")

		// --- When ---
		err := Before(date, mark, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be before mark:\n" +
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
		date := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(date, mark)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(mark, date)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be after mark:\n" +
			"  date: 2000-01-02T03:04:05Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("before", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := After(date, mark)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be after mark:\n" +
			"  date: 2000-01-02T03:04:05Z\n" +
			"  mark: 2001-01-02T03:04:05Z\n" +
			"  diff: -8784h0m0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid date", func(t *testing.T) {
		// --- When ---
		err := After("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid mark", func(t *testing.T) {
		// --- When ---
		err := After(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithTrail("type.field")

		// --- When ---
		err := After(date, mark, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be after mark:\n" +
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
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(date, mark)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(date, mark)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("after", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := BeforeOrEqual(date, mark)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be equal or before mark:\n" +
			"  date: 2000-01-02T03:04:06Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: 1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid date", func(t *testing.T) {
		// --- When ---
		err := BeforeOrEqual("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid mark", func(t *testing.T) {
		// --- When ---
		err := BeforeOrEqual(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithTrail("type.field")

		// --- When ---
		err := BeforeOrEqual(date, mark, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be equal or before mark:\n" +
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
		date := time.Date(2000, 1, 2, 3, 4, 6, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(date, mark)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(date, mark)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("before", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := AfterOrEqual(date, mark)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be equal or after mark:\n" +
			"  date: 2000-01-02T03:04:04Z\n" +
			"  mark: 2000-01-02T03:04:05Z\n" +
			"  diff: -1s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid date", func(t *testing.T) {
		// --- When ---
		err := AfterOrEqual("abc", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[date] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid mark", func(t *testing.T) {
		// --- When ---
		err := AfterOrEqual(time.Now(), "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[mark] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		date := time.Date(2000, 1, 2, 3, 4, 4, 0, time.UTC)
		mark := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithTrail("type.field")

		// --- When ---
		err := AfterOrEqual(date, mark, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected date to be equal or after mark:\n" +
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

	t.Run("not within", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Within(want, "1000s", have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2001-01-02T03:04:05Z\n" +
			"  max diff +/-: 1000s\n" +
			"     have diff: 8784h0m0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid want date format", func(t *testing.T) {
		// --- When ---
		err := Within("2022-02-18", "1s", time.Now())

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[want] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid within format", func(t *testing.T) {
		// --- When ---
		err := Within(time.Now(), "abc", "2000-01-02T03:04:05Z")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[within] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Within(time.Now(), "1s", "2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 6, int(500*time.Millisecond), time.UTC)
		opt := WithTrail("type.field")

		// --- When ---
		err := Within(want, "1s", have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected dates to be within:\n" +
			"         trail: type.field\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:04:06.5Z\n" +
			"  max diff +/-: 1s\n" +
			"     have diff: 1.5s"
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
		opt := WithNow(func() time.Time { return now })
		have := now.Add(9 * time.Second)

		// --- When ---
		err := Recent(have, opt)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("recent in the past", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithNow(func() time.Time { return now })
		have := now.Add(-9 * time.Second)

		// --- When ---
		err := Recent(have, opt)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not recent in the future", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithNow(func() time.Time { return now })
		have := now.Add(11 * time.Second)

		// --- When ---
		err := Recent(have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:04:16Z\n" +
			"  max diff +/-: 10s\n" +
			"     have diff: 11s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not recent in the past", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		opt := WithNow(func() time.Time { return now })
		have := now.Add(-11 * time.Second)

		// --- When ---
		err := Recent(have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected dates to be within:\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:03:54Z\n" +
			"  max diff +/-: 10s\n" +
			"     have diff: -11s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid have date format", func(t *testing.T) {
		// --- When ---
		err := Recent("2022-02-18")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[have] failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05.999999999Z07:00\n" +
			"   value: 2022-02-18"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := now.Add(-11 * time.Second)
		opts := []Option{
			WithNow(func() time.Time { return now }),
			WithTrail("type.field"),
		}

		// --- When ---
		err := Recent(have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected dates to be within:\n" +
			"         trail: type.field\n" +
			"          want: 2000-01-02T03:04:05Z\n" +
			"          have: 2000-01-02T03:03:54Z\n" +
			"  max diff +/-: 10s\n" +
			"     have diff: -11s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Zone(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- When ---
		err := Zone(time.UTC, time.UTC)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("nil is equivalent to UTC", func(t *testing.T) {
		affirm.Nil(t, Zone(nil, time.UTC))
		affirm.Nil(t, Zone(time.UTC, nil))
	})

	t.Run("not equal", func(t *testing.T) {
		// --- When ---
		err := Zone(time.UTC, types.WAW)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same timezone:\n" +
			"  want: UTC\n" +
			"  have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Zone(time.UTC, types.WAW, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same timezone:\n" +
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
		wMsg := "expected equal time durations:\n" +
			"  want: 1000s\n" +
			"  have: 2000s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid want", func(t *testing.T) {
		// --- When ---
		err := Duration("abc", "2000s")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[want] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid have", func(t *testing.T) {
		// --- When ---
		err := Duration("2000s", "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[have] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Duration(time.Second, 2*time.Second, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected equal time durations:\n" +
			"  trail: type.field\n" +
			"   want: 1s\n" +
			"   have: 2s"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_formatDates(t *testing.T) {
	t.Run("same format", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW)
		wTimStr := "2000-01-02T03:04:05+02:00"
		hTim := time.Date(2001, 1, 2, 3, 4, 5, 0, types.WAW)
		hTimStr := "2001-01-02T03:04:05+02:00"

		// --- When ---
		wHave, hHave := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		affirm.Equal(t, "2000-01-02T02:04:05Z ( 2000-01-02T03:04:05+02:00 )", wHave)
		affirm.Equal(t, "2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+02:00 )", hHave)
	})

	t.Run("shorted format when both dates in UTC", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		wTimStr := "2000-01-02T03:04:05Z"
		hTim := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		hTimStr := "2001-01-02T03:04:05Z"

		// --- When ---
		wHave, hHave := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		affirm.Equal(t, "2000-01-02T03:04:05Z", wHave)
		affirm.Equal(t, "2001-01-02T03:04:05Z", hHave)
	})

	t.Run("want date string is longer", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2001, 1, 2, 3, 4, 5, 0, types.WAW)
		wTimStr := "2001-01-02T03:04:05+01:00"
		hTim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		hTimStr := "2000-01-02T03:04:05Z"

		// --- When ---
		wHave, hHave := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		affirm.Equal(t, "2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+01:00 )", wHave)
		affirm.Equal(t, "2000-01-02T03:04:05Z", hHave)
	})

	t.Run("have date string is longer", func(t *testing.T) {
		// --- Given ---
		wTim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		wTimStr := "2000-01-02T03:04:05Z"
		hTim := time.Date(2001, 1, 2, 3, 4, 5, 0, types.WAW)
		hTimStr := "2001-01-02T03:04:05+01:00"

		// --- When ---
		wHave, hHave := formatDates(wTim, wTimStr, hTim, hTimStr)

		// --- Then ---
		affirm.Equal(t, "2000-01-02T03:04:05Z", wHave)
		affirm.Equal(t, "2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+01:00 )", hHave)
	})
}

func Test_getTime(t *testing.T) {
	t.Run("wrong time format", func(t *testing.T) {
		// --- Given ---
		opt := WithTimeFormat(time.RFC3339)

		// --- When ---
		haveTim, haveStr, haveRep, err := getTime("2000-01-02", opt)

		// --- Then ---
		affirm.Equal(t, true, haveTim.IsZero())
		affirm.Equal(t, "2000-01-02", haveStr)
		affirm.Equal(t, timeTypeStr, haveRep)
		wMsg := "failed to parse time:\n" +
			"  format: 2006-01-02T15:04:05Z07:00\n" +
			"   value: 2000-01-02"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeParse))
	})

	t.Run("empty option time format", func(t *testing.T) {
		// --- Given ---
		opt := WithTimeFormat("")

		// --- When ---
		have, haveStr, haveRep, err := getTime("2000-01-02", opt)

		// --- Then ---
		affirm.Equal(t, true, have.IsZero())
		affirm.Equal(t, "2000-01-02", haveStr)
		affirm.Equal(t, timeTypeStr, haveRep)
		wMsg := "failed to parse time:\n" +
			"  format: \n" +
			"   value: 2000-01-02\n" +
			"   error: extra text: \"2000-01-02\""
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeParse))
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opts := []Option{
			WithTimeFormat(time.RFC3339),
			WithTrail("type.field"),
		}

		// --- When ---
		have, haveStr, haveRep, err := getTime("2000-01-02", opts...)

		// --- Then ---
		affirm.Equal(t, true, have.IsZero())
		affirm.Equal(t, "2000-01-02", haveStr)
		affirm.Equal(t, timeTypeStr, haveRep)
		wMsg := "failed to parse time:\n" +
			"   trail: type.field\n" +
			"  format: 2006-01-02T15:04:05Z07:00\n" +
			"   value: 2000-01-02"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("unsupported type", func(t *testing.T) {
		// --- When ---
		have, haveStr, haveRep, err := getTime(true)

		// --- Then ---
		affirm.Equal(t, true, have.IsZero())
		affirm.Equal(t, "true", haveStr)
		affirm.Equal(t, "", haveRep)
		wMsg := "failed to parse time:\n" +
			"  cause: not supported time type"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrTimeType))
	})
}

func Test_getTime_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		opts    []Option
		have    any
		haveStr string
		haveRep timeRep
		want    time.Time
		wantTZ  *time.Location
	}{
		{
			"time.Time in UTC",
			nil,
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"2000-01-02T03:04:05Z",
			timeTypeTim,
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.UTC,
		},
		{
			"time.Time in WAW",
			nil,
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			"2000-01-02T03:04:05+01:00",
			timeTypeTim,
			time.Date(2000, 1, 2, 3, 4, 5, 0, types.WAW),
			types.WAW,
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
			"Unix timestamp int",
			nil,
			946778645,
			"946778645",
			timeTypeInt,
			time.Date(2000, 1, 2, 2, 4, 5, 0, time.UTC),
			time.UTC,
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
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveTim, haveStr, haveRep, err := getTime(tc.have, tc.opts...)

			// --- Then ---
			affirm.Nil(t, err)
			affirm.Equal(t, tc.haveStr, haveStr)
			affirm.Equal(t, tc.haveRep, haveRep)
			affirm.Equal(t, true, tc.want.Equal(haveTim))
			affirm.Equal(t, tc.wantTZ.String(), haveTim.Location().String())
		})
	}
}

func Test_getDur(t *testing.T) {
	t.Run("invalid string duration", func(t *testing.T) {
		// --- When ---
		haveDur, haveStr, haveRep, err := getDur("abc")

		// --- Then ---
		affirm.Equal(t, time.Duration(0), haveDur)
		affirm.Equal(t, "abc", haveStr)
		affirm.Equal(t, durTypeStr, haveRep)
		wMsg := "failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrDurParse))
	})

	t.Run("error - unsupported type", func(t *testing.T) {
		// --- When ---
		haveDur, haveStr, haveRep, err := getDur(true)

		// --- Then ---
		affirm.Equal(t, time.Duration(0), haveDur)
		affirm.Equal(t, "true", haveStr)
		affirm.Equal(t, "", haveRep)
		wMsg := "failed to parse duration:\n" +
			"  cause: not supported duration type"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrDurType))
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		haveDur, haveStr, haveRep, err := getDur(true, opt)

		// --- Then ---
		affirm.Equal(t, time.Duration(0), haveDur)
		affirm.Equal(t, "true", haveStr)
		affirm.Equal(t, "", haveRep)
		wMsg := "failed to parse duration:\n" +
			"  trail: type.field\n" +
			"  cause: not supported duration type"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, errors.Is(err, ErrDurType))
	})
}

func Test_getDur_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		have    any
		haveStr string
		haveRep durRep
		want    time.Duration
	}{
		{"time.Duration", time.Second, "1s", durTypeDur, time.Second},
		{"time.Duration as string", "1s", "1s", durTypeStr, time.Second},
		{
			"time.Duration as int",
			12345678,
			"12345678",
			durTypeInt,
			time.Duration(12345678),
		},
		{
			"time.Duration as int",
			int64(12345678),
			"12345678",
			durTypeInt64,
			time.Duration(12345678),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveDur, haveStr, haveRep, err := getDur(tc.have)

			// --- Then ---
			affirm.Nil(t, err)
			affirm.Equal(t, tc.want, haveDur)
			affirm.Equal(t, tc.haveStr, haveStr)
			affirm.Equal(t, tc.haveRep, haveRep)
		})
	}
}
