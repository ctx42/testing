// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
)

// HexPtrDumper is a generic hex dumper for pointers. It expects val to
// represent one of the kinds:
//
//   - [reflect.Uint8]
//   - [reflect.Uintptr]
//   - [reflect.UnsafePointer]
//
// Returns [valErrUsage] ("<dump-usage-error>") string if kind cannot be
// matched. It returns string representation in the format defined by [Dump]
// configuration.
func HexPtrDumper(dmp Dump, lvl int, val reflect.Value) string {
	var str string
	switch val.Kind() {
	case reflect.Uint8:
		str = fmt.Sprintf("0x%x", val.Uint())
	case reflect.Uintptr:
		if !dmp.PtrAddr {
			return ValAddr
		}
		str = fmt.Sprintf("<0x%x>", val.Uint())
	case reflect.UnsafePointer:
		if !dmp.PtrAddr {
			return ValAddr
		}
		str = fmt.Sprintf("<0x%x>", val.Pointer())
	default:
		str = ValErrUsage
	}
	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
