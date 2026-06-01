// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"time"
)

// ZoneDumper is the built-in dumper for *time.Location (and time.Location).
// It renders the location name. Returns [ValErrUsage] for wrong types.
func ZoneDumper(dmp Dump, lvl int, val reflect.Value) string {
	loc, ok := val.Interface().(time.Location)
	if !ok {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	val = reflect.ValueOf((&loc).String())
	return SimpleDumper(dmp, lvl, val)
}
