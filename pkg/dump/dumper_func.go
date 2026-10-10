// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
)

// FuncDumper is a generic dumper for functions. It expects val to represent
// one of the kinds:
//
//   - [reflect.Func]
//
// Returns [ValErrUsage] ("<dump-usage-error>") string if the kind cannot
// be matched. It returns the string representation in the format defined by
// [Dump] configuration.
func FuncDumper(dmp Dump, lvl int, val reflect.Value) string {
	var str string
	switch val.Kind() {
	case reflect.Func:
		ptrAddr := ValAddr
		if dmp.PtrAddr {
			ptrAddr = fmt.Sprintf("<0x%x>", val.Pointer())
		}
		str = fmt.Sprintf("%s(%s)", ValFunc, ptrAddr)

	default:
		str = ValErrUsage
	}

	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
