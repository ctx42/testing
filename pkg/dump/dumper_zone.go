// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"time"
)

// ZoneDumper requires val to be dereferenced representation of [time.Location].
// Returns [valErrUsage] ("<dump-usage-error>") string if the kind cannot be
// matched. It returns string representation in the format defined by [Dump]
// configuration.
func ZoneDumper(dmp Dump, lvl int, val reflect.Value) string {
	loc, ok := val.Interface().(time.Location)
	if !ok {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	val = reflect.ValueOf((&loc).String())
	return SimpleDumper(dmp, lvl, val)
}
