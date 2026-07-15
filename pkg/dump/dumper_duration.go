// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"time"
)

// Formats used by [GetDurDumper].
// Special format values for [time.Duration] accepted by [WithTimeFormat]
// (via DurationFormat), [GetDurDumper], and [Dump.DurationFormat].
const (
	// DurAsString is the default (empty) format, using time.Duration.String.
	DurAsString = ""

	// DurAsSeconds renders the duration as a floating-point number of seconds.
	DurAsSeconds = "<seconds>"
)

// GetDurDumper returns a [Dumper] for [time.Duration] values using the given
// format (one of the DurAs* constants).
func GetDurDumper(format string) Dumper {
	switch format {
	case DurAsString:
		return DurDumperString
	case DurAsSeconds:
		return DurDumperSeconds
	default:
		return DurDumperString
	}
}

// DurDumperString is the built-in dumper for [time.Duration] using the
// String() representation. Returns [ValErrUsage] for wrong kinds. Uses
// [reflect.Value.Int] so unexported duration fields still dump.
func DurDumperString(dmp Dump, lvl int, val reflect.Value) string {
	if val.Kind() != reflect.Int64 {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	tim := time.Duration(val.Int())
	val = reflect.ValueOf(tim.String())
	return SimpleDumper(dmp, lvl, val)
}

// DurDumperSeconds requires val to be a [time.Duration] (int64 kind) and
// returns seconds as configured by [Dump]. Returns [ValErrUsage] for wrong
// kinds. Uses [reflect.Value.Int] so unexported duration fields still dump.
func DurDumperSeconds(dmp Dump, lvl int, val reflect.Value) string {
	if val.Kind() != reflect.Int64 {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	tim := time.Duration(val.Int())
	val = reflect.ValueOf(tim.Seconds())
	return SimpleDumper(dmp, lvl, val)
}
