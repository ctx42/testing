// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"time"
)

// Formats used by [GetDurDumper].
const (
	DurAsString  = ""          // Same format as [time.Duration.String].
	DurAsSeconds = "<seconds>" // Duration as seconds float.
)

// GetDurDumper returns [time.Duration] dumper based on format.
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

// DurDumperString requires val to be dereferenced representation of
// [time.Duration] and returns its string representation in format defined
// by [Dump] configuration.
func DurDumperString(dmp Dump, lvl int, val reflect.Value) string {
	tim := val.Interface().(time.Duration) // nolint: forcetypeassert
	val = reflect.ValueOf(tim.String())
	return simpleDumper(dmp, lvl, val)
}

// DurDumperSeconds requires val to be dereferenced representation of
// [time.Duration] and returns its string representation in format defined
// by [Dump] configuration.
func DurDumperSeconds(dmp Dump, lvl int, val reflect.Value) string {
	tim := val.Interface().(time.Duration) // nolint: forcetypeassert
	val = reflect.ValueOf(tim.Seconds())
	return simpleDumper(dmp, lvl, val)
}
