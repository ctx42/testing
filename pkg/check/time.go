// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

// Sentinel errors for time, timezone, and duration handling.
//
// These errors are intended to be used with [errors.Is]. When functions such as
// [Time], [Exact], [Zone], or [Duration] return an error due to unsupported
// input types or parse failures, the returned error will match one of these
// sentinels.
var (
	// ErrTimeType indicates that the provided value is not a supported time
	// representation (see package documentation for supported types).
	ErrTimeType = errors.New("not supported time type")

	// ErrTimeParse indicates that parsing a time value failed.
	ErrTimeParse = errors.New("time parsing")

	// ErrZoneType indicates that the provided value is not a supported timezone
	// representation.
	ErrZoneType = errors.New("not supported timezone type")

	// ErrZoneParse indicates that parsing a timezone value failed.
	ErrZoneParse = errors.New("timezone parsing")

	// ErrDurType indicates that the provided value is not a supported duration
	// representation.
	ErrDurType = errors.New("not supported duration type")

	// ErrDurParse indicates that parsing a duration value failed.
	ErrDurParse = errors.New("duration parsing")
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

// Time checks that "want" and "have" represent the same instant in time.
//
// See the package documentation and [Options] for supported representations
// and formatting behavior.
//
// On error due to unsupported input or parse failure, the returned error will
// satisfy `errors.Is(err, ErrTimeType)` or `errors.Is(err, ErrTimeParse)`.
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
// int64, or [time.Time]. For string representations the [Options.TimeFormat]
// is used during parsing and the returned date is always in UTC. The int and
// int64 types are interpreted as Unix Timestamp, and the date returned is also
// in UTC.
//
// On error due to unsupported input or parse failure, the returned error will
// satisfy `errors.Is(err, ErrTimeType)` or `errors.Is(err, ErrTimeParse)`.
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

// Before checks that "date" is before "mark".
//
// Both arguments may be strings, ints, int64s, or [time.Time]. See the
// package documentation or [Exact] for supported representations and
// formatting behavior.
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

// After checks that "date" is after "mark".
//
// Both arguments may be strings, ints, int64s, or [time.Time]. See the
// package documentation or [Exact] for supported representations and
// formatting behavior.
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

// BeforeOrEqual checks that "date" is before or equal to "mark".
//
// See the package documentation for supported time representations.
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

// AfterOrEqual checks that "date" is equal to or after "mark".
//
// See the package documentation for supported time representations.
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

// Within checks that the absolute difference between "want" and "have" is at
// most the given "within" duration.
//
// Arguments may be strings, ints, int64s or [time.Time] (for instants);
// "within" may also be a duration string/int64/[time.Duration]. See package
// docs or [Exact] for details.
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

// WithinChecker returns a [Checker] that checks whether two times are within
// the given duration (a partial application of [Within]).
//
// Useful for registering via [WithTypeChecker] or [RegisterTypeChecker] for
// custom time tolerance behavior.
func WithinChecker(within any) Checker {
	return func(want, have any, opts ...any) error {
		return Within(want, within, have, opts...)
	}
}

// Recent checks that "have" is within the [Options.Recent] duration of
// the current time (see [Options.now]).
//
// "have" may be a string, int, int64, or [time.Time]. See the package
// documentation or [Exact] for representation rules.
func Recent(have any, opts ...any) error {
	ops := DefaultOptions(opts...)
	return Within(ops.now(), ops.Recent, have, opts...)
}

// Zone checks that "want" and "have" timezones are equal.
//
// nil [time.Location] is treated as [time.UTC]. Arguments may be strings or
// [time.Location] values.
//
// On error due to unsupported input or parse failure, the returned error will
// satisfy `errors.Is(err, ErrZoneType)` or `errors.Is(err, ErrZoneParse)`.
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

// Duration checks that "want" and "have" durations are equal.
//
// Arguments may be strings, ints, int64s or [time.Duration].
//
// On error due to unsupported input or parse failure, the returned error will
// satisfy `errors.Is(err, ErrDurType)` or `errors.Is(err, ErrDurParse)`.
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

// getTime returns the date represented by "tim", its string representation,
// and type of the argument passed. The "tim" may represent date in the form of
// a string, int, int64, or [time.Time]. For string representations the
// [Options.TimeFormat] is used during parsing and the returned date is always
// in UTC. The int and int64 types are interpreted as Unix Timestamp, and the
// date returned is also in UTC.
//
// On error, the returned error will always satisfy
// `errors.Is(err, ErrTimeParse)` or `errors.Is(err, ErrTimeType)`.
func getTime(tim any, opts ...any) (time.Time, string, timeRep, error) {
	ops := DefaultOptions(opts...)
	switch val := tim.(type) {
	case time.Time:
		if ops.Zone != nil {
			val = val.In(ops.Zone)
		}
		return val, val.Format(time.RFC3339Nano), timeTypeTim, nil

	case string:
		if ops.TimeFormat == TimeFormatUnixStr {
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				msg := notice.New("failed to parse time").
					Append("format", "%s", ops.TimeFormat).
					Append("value", "%s", val).
					Wrap(ErrTimeParse)
				return time.Time{}, val, timeTypeStr, AddRows(ops, msg)
			}
			ts := time.Unix(n, 0)
			if ops.Zone != nil {
				ts = ts.In(ops.Zone)
			} else {
				ts = ts.UTC()
			}
			return ts, val, timeTypeStr, nil
		}

		have, err := time.Parse(ops.TimeFormat, val)
		if err == nil {
			if ops.Zone != nil {
				have = have.In(ops.Zone)
			} else {
				have = have.UTC()
			}
			return have, val, timeTypeStr, nil
		}

		if pe, ok := errors.AsType[*time.ParseError](err); ok {
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
// On error, the returned error will always satisfy
// `errors.Is(err, ErrZoneParse)` or `errors.Is(err, ErrZoneType)`.
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
// form of a string, int, int64, or [time.Duration].
//
// On error, the returned error will always satisfy
// `errors.Is(err, ErrDurParse)` or `errors.Is(err, ErrDurType)`.
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
