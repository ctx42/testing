// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
	"time"
)

// Special format values for [time.Time] accepted by [WithTimeFormat],
// [GetTimeDumper], and [Dump.TimeFormat].
const (
	// TimeAsRFC3339 is the default (empty) format, producing RFC3339Nano.
	TimeAsRFC3339 = ""

	// TimeAsUnix produces a Unix timestamp (seconds since epoch) as a
	// decimal number.
	TimeAsUnix = "<unix>"

	// TimeAsGoString produces output identical to time.Time.GoString().
	TimeAsGoString = "<go-str>"
)

// GetTimeDumper returns a [Dumper] for [time.Time] values using the given
// format (one of the TimeAs* constants or a standard Go time layout).
func GetTimeDumper(format string) Dumper {
	switch format {
	case "":
		return TimeDumperFmt(time.RFC3339Nano)
	case TimeAsUnix:
		return TimeDumperUnix
	case TimeAsGoString:
		return TimeDumperDate
	default:
		return TimeDumperFmt(format)
	}
}

// TimeDumperFmt returns a [Dumper] that formats [time.Time] using the
// supplied Go time layout. The returned function returns [ValErrUsage] if
// the value is not a time.Time.
func TimeDumperFmt(format string) Dumper {
	return func(dmp Dump, lvl int, val reflect.Value) string {
		tim, ok := val.Interface().(time.Time)
		if !ok {
			prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
			return prn.Write(ValErrUsage).String()
		}
		val = reflect.ValueOf(tim.Format(format))
		return SimpleDumper(dmp, lvl, val)
	}
}

// TimeDumperUnix requires val to be a value representing [time.Time] and
// returns its string representation as a Unix timestamp. Returns [valErrUsage]
// ("<dump-usage-error>") string if the type cannot be matched.
func TimeDumperUnix(dmp Dump, lvl int, val reflect.Value) string {
	ts, ok := val.Interface().(time.Time)
	if !ok {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	val = reflect.ValueOf(ts.Unix())
	return SimpleDumper(dmp, lvl, val)
}

// TimeDumperDate requires val to be a value representing [time.Time] and
// returns its representation using [time.Time.GoString] method. Returns
// [valErrUsage] ("<dump-usage-error>") string if the type cannot be matched.
func TimeDumperDate(dmp Dump, lvl int, val reflect.Value) string {
	ts, ok := val.Interface().(time.Time)
	if !ok {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	str := ts.GoString()
	if dmp.Compact {
		str = strings.ReplaceAll(str, " ", "")
	}
	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
