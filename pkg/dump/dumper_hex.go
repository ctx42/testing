// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
)

// hexPtrDumper is a generic hex dumper for pointers. It expects val to
// represent one of the kinds:
//
//   - [reflect.Uint8]
//   - [reflect.Uintptr]
//   - [reflect.UnsafePointer]
//
// Returns [valErrUsage] ("<dump-usage-error>") string if kind cannot be
// matched. It requires val to be a dereferenced value and returns its string
// representation in the format defined by [Dump] configuration.
func hexPtrDumper(dmp Dump, lvl int, val reflect.Value) string {
	var str string
	switch val.Kind() {
	case reflect.Uint8:
		str = fmt.Sprintf("0x%x", val.Interface())
	case reflect.Uintptr:
		str = fmt.Sprintf("<0x%x>", val.Uint())
	case reflect.UnsafePointer:
		str = fmt.Sprintf("<0x%x>", val.Pointer())
	default:
		str = ValErrUsage
	}

	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
