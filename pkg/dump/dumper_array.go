// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
)

// ArrayDumper is a generic dumper for arrays. It expects val to represent one
// of the kinds:
//
//   - [reflect.Array]
//   - [reflect.Slice]
//
// Returns [valErrUsage] ("<dump-usage-error>") string if the kind cannot be
// matched. It returns string representation in the format defined by [Dump]
// configuration.
func ArrayDumper(dmp Dump, lvl int, val reflect.Value) string {
	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)

	if !(val.Kind() == reflect.Slice || val.Kind() == reflect.Array) {
		return prn.Write(ValErrUsage).String()
	}

	if dmp.PrintType {
		valTypStr := val.Type().String()
		if dmp.UseAny {
			valTypStr = strings.Replace(valTypStr, "interface {}", "any", 1)
		}
		prn.Write(valTypStr)
	}

	num := val.Len()
	prn.Write("{").NLI(num)

	dmp.PrintType = false // Don't print types for array elements.
	for i := 0; i < num; i++ {
		last := i == num-1

		sub, _ := dmp.value(lvl+1, val.Index(i))
		prn.Write(sub)
		prn.Comma(last).Sep(last).NL()
	}
	prn.Tab(dmp.Indent + lvl).Write("}")

	return prn.String()
}
