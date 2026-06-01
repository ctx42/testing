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
// String() representation. Returns [ValErrUsage] for wrong types.
func DurDumperString(dmp Dump, lvl int, val reflect.Value) string {
	tim, ok := val.Interface().(time.Duration)
	if !ok {
		prn := NewPrinter(dmp)
		return prn.Write(ValErrUsage).String()
	}
	val = reflect.ValueOf(tim.String())
	return SimpleDumper(dmp, lvl, val)
}

// DurDumperSeconds requires val to be dereferenced representation of
// [time.Duration]. Returns [valErrUsage] ("<dump-usage-error>") string if the
// type cannot be matched. It returns string representation in the format
// // defined by [Dump] configuration.
func DurDumperSeconds(dmp Dump, lvl int, val reflect.Value) string {
	tim, ok := val.Interface().(time.Duration)
	if !ok {
		prn := NewPrinter(dmp)
		return prn.Write(ValErrUsage).String()
	}
	val = reflect.ValueOf(tim.Seconds())
	return SimpleDumper(dmp, lvl, val)
}
