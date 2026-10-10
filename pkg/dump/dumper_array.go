// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
)

// ArrayDumper is a generic dumper for arrays. It expects val to represent one
// of the kinds:
//
//   - [reflect.Array]
//   - [reflect.Slice]
//
// Returns [ValErrUsage] ("<dump-usage-error>") string if the kind cannot be
// matched. It returns the string representation in the format defined by [Dump]
// configuration.
func ArrayDumper(dmp Dump, lvl int, val reflect.Value) string {
	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return prn.Write(ValErrUsage).String()
	}

	if dmp.PrintType {
		prn.Write(typeName(dmp, val.Type()))
	}

	num := val.Len()
	prn.Write("{").NLI(num)

	dmp.PrintType = false // Don't print types for array elements.
	for i := range num {
		last := i == num-1

		sub, _ := dmp.value(lvl+1, val.Index(i))
		prn.Write(sub)
		prn.Comma(last).Sep(last).NL()
	}
	if num > 0 {
		prn.Tab(dmp.Indent + lvl)
	}
	prn.Write("}")

	return prn.String()
}
