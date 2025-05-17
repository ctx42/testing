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

// TimeDumperFmt returns [Dumper] for [time.Time] using given format. The
// returned function requires val to be dereferenced representation of
// [reflect.Struct] which can be cast to [time.Time] and returns its string
// representation in format defined by [Dump] configuration.
func TimeDumperFmt(format string) Dumper {
	return func(dmp Dump, lvl int, val reflect.Value) string {
		tim := val.Interface().(time.Time) // nolint: forcetypeassert
		val = reflect.ValueOf(tim.Format(format))
		return simpleDumper(dmp, lvl, val)
	}
}

// TimeDumperUnix requires val to be dereferenced representation of
// [reflect.Struct] which can be cast to [time.Time] and returns its string
// representation as Unix timestamp.
func TimeDumperUnix(dmp Dump, lvl int, val reflect.Value) string {
	ts := val.Interface().(time.Time).Unix() // nolint: forcetypeassert
	val = reflect.ValueOf(ts)
	return simpleDumper(dmp, lvl, val)
}

// TimeDumperDate requires val to be dereferenced representation of
// [reflect.Struct] which can be cast to [time.Time] and returns its
// representation using [time.Time.GoString] method.
func TimeDumperDate(dmp Dump, lvl int, val reflect.Value) string {
	ts := val.Interface().(time.Time) // nolint: forcetypeassert
	str := ts.GoString()
	if dmp.Compact {
		str = strings.ReplaceAll(str, " ", "")
	}
	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
