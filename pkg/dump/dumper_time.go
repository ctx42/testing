// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
	"time"
)

// Formats used by [GetTimeDumper].
const (
	TimeAsRFC3339  = ""         // Formats time as [time.RFC3339Nano].
	TimeAsUnix     = "<unix>"   // Formats time as Unix timestamp (seconds).
	TimeAsGoString = "<go-str>" // Formats time the same way as [time.GoString].
)

// GetTimeDumper returns [time.Time] dumper based on format.
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

// TimeDumperFmt returns [Dumper] for [time.Time] using the given format. The
// returned function requires val to be a value representing [time.Time] and
// returns its string representation in the format defined by [Dump]
// configuration. Returns [valErrUsage] ("<dump-usage-error>") string if the
// type cannot be matched.
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
