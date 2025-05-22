// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
)

// complexDumper is a generic dumper for complex values. It expects val to
// represent one of the kinds:
//
//   - [reflect.Complex64]
//   - [reflect.Complex128]
//
// Returns [valErrUsage] ("<dump-usage-error>") string if kind cannot be
// matched. It requires val to be a dereferenced value and returns its string
// representation in the format defined by [Dump] configuration.
func complexDumper(dmp Dump, lvl int, val reflect.Value) string {
	var str string
	switch val.Kind() {
	case reflect.Complex64, reflect.Complex128:
		str = fmt.Sprintf("%v", val.Interface())
	default:
		str = ValErrUsage
	}

	prn := NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
