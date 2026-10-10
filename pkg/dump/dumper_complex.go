// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
)

// ComplexDumper is a generic dumper for complex values. It expects val to
// represent one of the kinds:
//
//   - [reflect.Complex64]
//   - [reflect.Complex128]
//
// Returns [ValErrUsage] ("<dump-usage-error>") string if the kind cannot
// be matched. It returns the string representation in the format defined by
// [Dump] configuration.
func ComplexDumper(dmp Dump, lvl int, val reflect.Value) string {
	var str string
	switch val.Kind() {
	case reflect.Complex64:
		// Use Complex(), not Interface(), so unexported fields still dump.
		str = fmt.Sprintf("%v", complex64(val.Complex()))

	case reflect.Complex128:
		str = fmt.Sprintf("%v", val.Complex())

	default:
		str = ValErrUsage
	}

	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
