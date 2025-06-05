// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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
// Returns [valErrUsage] ("<dump-usage-error>") string if kind cannot be
// matched. It returns string representation in the format defined by [Dump]
// configuration.
func FuncDumper(dmp Dump, lvl int, val reflect.Value) string {
	var str string
	switch val.Kind() {
	case reflect.Func:
		ptrAddr := ValAddr
		if dmp.PtrAddr {
			ptr := reflect.ValueOf(val.Pointer())
			ptrAddr = HexPtrDumper(dmp, lvl, ptr)
		}
		str = fmt.Sprintf("%s(%s)", ValFunc, ptrAddr)
	default:
		str = ValErrUsage
	}

	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
