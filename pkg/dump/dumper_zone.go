// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"time"

	"github.com/ctx42/testing/internal/core"
)

// ZoneDumper is the built-in dumper for *time.Location (and time.Location).
// It renders the location name. Returns [ValErrUsage] for wrong types.
func ZoneDumper(dmp Dump, lvl int, val reflect.Value) string {
	if val.Type() != typLocation {
		prn := NewPrinter(dmp).Tab(dmp.Indent + lvl)
		return prn.Write(ValErrUsage).String()
	}
	var loc *time.Location
	if val.CanAddr() {
		// Use the original: a copy of time.Local may not be initialized.
		itf, _ := core.Value(val.Addr())
		loc, _ = itf.(*time.Location)
	} else {
		itf, _ := core.Value(val)
		cp, _ := itf.(time.Location)
		loc = &cp
	}
	val = reflect.ValueOf(loc.String())
	return SimpleDumper(dmp, lvl, val)
}
