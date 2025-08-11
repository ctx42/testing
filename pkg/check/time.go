// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ctx42/testing/pkg/notice"
)

// Sentinel errors.
var (
	// ErrTimeType is returned when time representation is not supported.
	ErrTimeType = fmt.Errorf("not supported time type")

	// ErrTimeParse is used when date parsing fails for whatever reason.
	ErrTimeParse = fmt.Errorf("time parsing")

	// ErrZoneType is returned when timezone representation is not supported.
	ErrZoneType = fmt.Errorf("not supported timezone type")

	// ErrZoneParse is used when timezone parsing fails for whatever reason.
	ErrZoneParse = fmt.Errorf("timezone parsing")

	// ErrDurType is returned when duration representation is not supported.
	ErrDurType = fmt.Errorf("not supported duration type")

	// ErrDurParse is used when duration parsing fails for whatever reason.
	ErrDurParse = fmt.Errorf("duration parsing")
)

// timeRep is time representation.
type timeRep string

// The time representations the [Time] supports.
const (
	timeTypeTim   timeRep = "tim-tim"
	timeTypeStr   timeRep = "tim-string"
	timeTypeInt   timeRep = "tim-int"
	timeTypeInt64 timeRep = "tim-int64"
)

// zoneRep is timezone representation.
type zoneRep string

// The timezone representations the [Zone] supports.
const (
	zoneString = "zone-string"
	zoneZone   = "zone-zone"
)

// durRep is duration representation.
type durRep string

// The duration representations the [Duration] supports.
const (
	durTypeDur   durRep = "dur-dur"
	durTypeStr   durRep = "dur-str"
	durTypeInt   durRep = "dur-int"
	durTypeInt64 durRep = "dur-int64"
)

// Time checks "want" and "have" dates are equal. Returns nil if they are,
// otherwise returns an error with a message indicating the expected and actual
// values.
//
// The "want" and "have" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
func Time(want, have any, opts ...any) error {
	ops := DefaultOptions(opts...)

	wTim, wStr, _, err := getTime(want, opts...)
	if err != nil {
		return notice.From(err, "want")
	}
	hTim, hStr, _, err := getTime(have, opts...)
	if err != nil {
		return notice.From(err, "have")
	}
	if wTim.Equal(hTim) {
		return nil
	}

	diff := wTim.Sub(hTim)
	wantFmt, haveFmt := formatDates(wTim, wStr, hTim, hStr)
	msg := notice.New("expected equal dates").
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("diff", "%s", diff.String())
	return AddRows(ops, msg)
}

// Exact checks "want" and "have" dates are equal and are in the same timezone.
// Returns nil they are, otherwise returns an error with a message indicating
// the expected and actual values.
//
// The "want" and "have" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
func Exact(want, have any, opts ...any) error {
	wTim, wStr, _, err := getTime(want, opts...)
	if err != nil {
		return notice.From(err, "want")
	}
	hTim, hStr, _, err := getTime(have, opts...)
	if err != nil {
		return notice.From(err, "have")
	}

	if !wTim.Equal(hTim) {
		diff := wTim.Sub(hTim)
		ops := DefaultOptions(opts...)
		wantFmt, haveFmt := formatDates(wTim, wStr, hTim, hStr)
		msg := notice.New("expected equal dates").
			Want("%s", wantFmt).
			Have("%s", haveFmt).
			Append("diff", "%s", diff.String())
		return AddRows(ops, msg)
	}

	return Zone(wTim.Location(), hTim.Location(), opts...)
}

// Before checks "date" is before "mark". Returns nil if it is, otherwise it
// returns an error with a message indicating the expected and actual values.
//
// The "want" and "have" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
func Before(mark, date any, opts ...any) error {
	dTim, dStr, _, err := getTime(date, opts...)
	if err != nil {
		return notice.From(err, "date")
	}
	mTim, mStr, _, err := getTime(mark, opts...)
	if err != nil {
		return notice.From(err, "mark")
	}
	if dTim.Before(mTim) {
		return nil
	}

	diff := dTim.Sub(mTim)
	markFmt, dateFmt := formatDates(mTim, mStr, dTim, dStr)
	ops := DefaultOptions(opts...)
	msg := notice.New("expected date to be before mark").
		Append("date", "%s", dateFmt).
		Append("mark", "%s", markFmt).
		Append("diff", "%s", diff.String())
	return AddRows(ops, msg)
}

// After checks "date" is after "mark". Returns nil if it is, otherwise it
// returns an error with a message indicating the expected and actual values.
//
// The "date" and "mark" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
func After(mark, date any, opts ...any) error {
	dTim, dStr, _, err := getTime(date, opts...)
	if err != nil {
		return notice.From(err, "date")
	}
	mTim, mStr, _, err := getTime(mark, opts...)
	if err != nil {
		return notice.From(err, "mark")
	}
	if dTim.After(mTim) {
		return nil
	}

	diff := dTim.Sub(mTim)
	markFmt, dateFmt := formatDates(mTim, mStr, dTim, dStr)
	ops := DefaultOptions(opts...)
	msg := notice.New("expected date to be after mark").
		Append("date", "%s", dateFmt).
		Append("mark", "%s", markFmt).
		Append("diff", "%s", diff.String())
	return AddRows(ops, msg)
}

// BeforeOrEqual checks "date" is equal or before "mark". Returns nil if it is,
// otherwise it returns an error with a message indicating the expected and
// actual values.
//
// The "date" and "mark" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
func BeforeOrEqual(mark, date any, opts ...any) error {
	dTim, dStr, _, err := getTime(date, opts...)
	if err != nil {
		return notice.From(err, "date")
	}
	mTim, mStr, _, err := getTime(mark, opts...)
	if err != nil {
		return notice.From(err, "mark")
	}
	if dTim.Equal(mTim) || dTim.Before(mTim) {
		return nil
	}

	diff := dTim.Sub(mTim)
	markFmt, dateFmt := formatDates(mTim, mStr, dTim, dStr)
	ops := DefaultOptions(opts...)
	msg := notice.New("expected date to be equal or before mark").
		Append("date", "%s", dateFmt).
		Append("mark", "%s", markFmt).
		Append("diff", "%s", diff.String())
	return AddRows(ops, msg)
}

// AfterOrEqual checks "date" is equal or after "mark". Returns nil if it is,
// otherwise it returns an error with a message indicating the expected and
// actual values.
//
// The "date" and "mark" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
func AfterOrEqual(mark, date any, opts ...any) error {
	dTim, dStr, _, err := getTime(date, opts...)
	if err != nil {
		return notice.From(err, "date")
	}
	mTim, mStr, _, err := getTime(mark, opts...)
	if err != nil {
		return notice.From(err, "mark")
	}
	if dTim.Equal(mTim) || dTim.After(mTim) {
		return nil
	}

	diff := dTim.Sub(mTim)
	markFmt, dateFmt := formatDates(mTim, mStr, dTim, dStr)
	ops := DefaultOptions(opts...)
	msg := notice.New("expected date to be equal or after mark").
		Append("date", "%s", dateFmt).
		Append("mark", "%s", markFmt).
		Append("diff", "%s", diff.String())
	return AddRows(ops, msg)
}

// Within checks "want" and "have" dates are equal "within" given duration.
// Returns nil if they are, otherwise returns an error with a message
// indicating the expected and actual values.
//
// The "want" and "have" may represent dates in the form of a string, int,
// int64 or [time.Time]. For string representations the [Options.TimeFormat] is
// used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
//
// The "within" may represent duration in the form of a string, int, int64 or
// [time.Duration].
func Within(want, within, have any, opts ...any) error {
	wTim, wStr, _, err := getTime(want, opts...)
	if err != nil {
		return notice.From(err, "want")
	}
	hTim, hStr, _, err := getTime(have, opts...)
	if err != nil {
		return notice.From(err, "have")
	}
	dur, durStr, _, err := getDur(within, opts...)
	if err != nil {
		return notice.From(err, "within")
	}

	diff := wTim.Sub(hTim.In(wTim.Location()))
	if math.Abs(float64(diff)) <= math.Abs(float64(dur)) {
		return nil
	}

	wantFmt, haveFmt := formatDates(wTim, wStr, hTim, hStr)
	ops := DefaultOptions(opts...)
	msg := notice.New("expected dates to be within").
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("max diff +/-", "%s", durStr).
		Append("have diff", "%s", diff.String())
	return AddRows(ops, msg)
}

// WithinChecker is a partial application function for [Within].
func WithinChecker(within any) Checker {
	return func(want, have any, opts ...any) error {
		return Within(want, within, have, opts...)
	}
}

// Recent checks "have" is within [Options.Recent] from [time.Now]. Returns nil
// if it is, otherwise returns an error with a message indicating the expected
// and actual values.
//
// The "have" may represent date in the form of a string, int, int64 or
// [time.Time]. For string representations the [Options.TimeFormat] is used
// during parsing and the returned date is always in UTC. The int and int64
// types are interpreted as Unix Timestamp, and the date returned is also in
// UTC.
func Recent(have any, opts ...any) error {
	ops := DefaultOptions(opts...)
	return Within(ops.now(), ops.Recent, have, opts...)
}

// Zone checks "want" and "have" timezones are equal. Returns nil if they are,
// otherwise returns an error with a message indicating the expected and actual
// values.
//
// Note nil [time.Location] is the same as [time.UTC].
//
// The "want" and "have" may represent timezones in the form of a string,
// nil (UTC), or [time.Location].
func Zone(want, have any, opts ...any) error {
	wZone, wStr, _, err := getZone(want, opts...)
	if err != nil {
		return notice.From(err, "want")
	}
	hZone, hStr, _, err := getZone(have, opts...)
	if err != nil {
		return notice.From(err, "have")
	}
	if wZone.String() == hZone.String() {
		return nil
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected timezones to be equal").
		Want("%s", wStr).
		Have("%s", hStr)
	return AddRows(ops, msg)
}

// Duration checks "want" and "have" durations are equal. Returns nil if they
// are, otherwise returns an error with a message indicating the expected and
// actual values.
//
// The "want" and "have" may represent duration in the form of a string, int,
// int64 or [time.Duration].
func Duration(want, have any, opts ...any) error {
	wDur, wStr, _, err := getDur(want, opts...)
	if err != nil {
		return notice.From(err, "want")
	}
	hDur, hStr, _, err := getDur(have, opts...)
	if err != nil {
		return notice.From(err, "have")
	}

	if wDur == hDur {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected equal time durations").
		Want("%s", wStr).
		Have("%s", hStr)
	return AddRows(ops, msg)
}

// formatDates formats two dates for comparison in an error message.
//
// Example:
//
//	2000-01-02T03:04:05Z ( 2000-01-02T03:04:05Z      )
//	2001-01-02T02:04:05Z ( 2001-01-02T03:04:05+01:00 )
func formatDates(
	wTim time.Time, wTimStr string,
	hTim time.Time, hTimStr string,
) (string, string) {

	wTimUTC := wTim.In(time.UTC).Format(time.RFC3339Nano)
	hTimUTC := hTim.In(time.UTC).Format(time.RFC3339Nano)

	wTimStrLen := len(wTimStr)
	hTimStrLen := len(hTimStr)

	var wTimPad, hTimPad string
	if wTimStrLen < hTimStrLen {
		wTimPad = strings.Repeat(" ", hTimStrLen-wTimStrLen)
	}
	if hTimStrLen < wTimStrLen {
		hTimPad = strings.Repeat(" ", wTimStrLen-hTimStrLen)
	}

	var want, have string
	if wTimUTC == wTimStr {
		want = wTimUTC
	} else {
		want = fmt.Sprintf("%s ( %s %s)", wTimUTC, wTimStr, wTimPad)

	}

	if hTimUTC == hTimStr {
		have = hTimUTC
	} else {
		have = fmt.Sprintf("%s ( %s %s)", hTimUTC, hTimStr, hTimPad)
	}

	return want, have
}

// getTime returns the date represented by "tim", its string representation and
// type of the argument passed. The "tim" may represent date in the form of a
// string, int, int64 or [time.Time]. For string representations the
// [Options.TimeFormat] is used during parsing and the returned date is always
// in UTC. The int and int64 types are interpreted as Unix Timestamp, and the
// date returned is also in UTC.
//
// When an error is returned, it will always have [ErrTimeParse], [ErrTimeType]
// in its chain.
func getTime(tim any, opts ...any) (time.Time, string, timeRep, error) {
	ops := DefaultOptions(opts...)
	switch val := tim.(type) {
	case time.Time:
		if ops.Zone != nil {
			val = val.In(ops.Zone)
		}
		return val, val.Format(time.RFC3339Nano), timeTypeTim, nil

	case string:
		have, err := time.Parse(ops.TimeFormat, val)
		if err == nil {
			if ops.Zone != nil {
				have = have.In(ops.Zone)
			} else {
				have = have.UTC()
			}
			return have, val, timeTypeStr, nil
		}

		var pe *time.ParseError
		if errors.As(err, &pe) {
			msg := notice.New("failed to parse time").
				Append("format", "%s", ops.TimeFormat).
				Append("value", "%s", pe.Value).
				Wrap(ErrTimeParse)
			if pe.Message != "" {
				msg = msg.Append("error", "%s", strings.Trim(pe.Message, " :"))
			}
			err = AddRows(ops, msg)
		}
		return time.Time{}, val, timeTypeStr, err

	case int:
		str := strconv.Itoa(val)
		ts := time.Unix(int64(val), 0)
		if ops.Zone != nil {
			ts = ts.In(ops.Zone)
		} else {
			ts = ts.UTC()
		}
		return ts, str, timeTypeInt, nil

	case int64:
		str := strconv.FormatInt(val, 10)
		ts := time.Unix(val, 0)
		if ops.Zone != nil {
			ts = ts.In(ops.Zone)
		} else {
			ts = ts.UTC()
		}
		return ts, str, timeTypeInt64, nil

	default:
		str := fmt.Sprintf("%v", val)
		msg := notice.New("failed to parse time").
			Append("cause", "%s", ErrTimeType).
			Wrap(ErrTimeType)
		return time.Time{}, str, "", AddRows(ops, msg)
	}
}

// getZone returns timezone represented by "zone", its string representation,
// and the type of the argument passed. The "zone" may represent a timezone in
// the form of a string, nil (UTC) or instance of [time.Location].
//
// When an error is returned, it will always have [ErrDurParse], [ErrZoneType]
// in its chain.
func getZone(zone any, opts ...any) (*time.Location, string, zoneRep, error) {
	switch val := zone.(type) {
	case nil:
		return time.UTC, "UTC", zoneZone, nil

	case time.Location:
		v := &val
		vStr := v.String()
		return v, vStr, zoneZone, nil

	case string:
		z, err := time.LoadLocation(val)
		if err == nil {
			return z, val, zoneString, nil
		}
		ops := DefaultOptions(opts...)
		msg := notice.New("failed to parse timezone").
			Append("value", "%s", zone).
			Wrap(ErrZoneParse)
		return nil, val, zoneString, AddRows(ops, msg)

	case *time.Location:
		valStr := val.String()
		return val, valStr, zoneZone, nil

	default:
		str := fmt.Sprintf("%v", val)
		ops := DefaultOptions(opts...)
		msg := notice.New("failed to parse timezone").
			Append("cause", "%s", ErrZoneType).
			Wrap(ErrZoneType)
		return nil, str, "", AddRows(ops, msg)
	}
}

// getDur returns duration represented by "dur", its string representation, and
// the type of the argument passed. The "dur" may represent duration in the
// form of a string, int, int64 or [time.Duration].
//
// When an error is returned, it will always have [ErrDurParse], [ErrDurType]
// in its chain.
func getDur(dur any, opts ...any) (time.Duration, string, durRep, error) {
	switch val := dur.(type) {
	case time.Duration:
		return val, val.String(), durTypeDur, nil

	case string:
		have, err := time.ParseDuration(val)
		if err == nil {
			return have, val, durTypeStr, nil
		}

		ops := DefaultOptions(opts...)
		msg := notice.New("failed to parse duration").
			Append("value", "%s", dur).
			Wrap(ErrDurParse)
		return 0, val, durTypeStr, AddRows(ops, msg)

	case int:
		str := strconv.Itoa(val)
		return time.Duration(val), str, durTypeInt, nil

	case int64:
		str := fmt.Sprintf("%v", val)
		return time.Duration(val), str, durTypeInt64, nil

	default:
		str := fmt.Sprintf("%v", val)
		ops := DefaultOptions(opts...)
		msg := notice.New("failed to parse duration").
			Append("cause", "%s", ErrDurType).
			Wrap(ErrDurType)
		return 0, str, "", AddRows(ops, msg)
	}
}
