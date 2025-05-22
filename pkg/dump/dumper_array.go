// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
)

// arrayDumper requires val to be dereferenced representation of [reflect.Array]
// or [reflect.Slice] and returns its string representation in format defined
// by [Dump] configuration.
func arrayDumper(dmp Dump, lvl int, val reflect.Value) string {
	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)

	if dmp.PrintType {
		valTypStr := val.Type().String()
		if dmp.UseAny {
			switch {
			case valTypStr == "interface{}":
				valTypStr = "any"
			case strings.HasSuffix(valTypStr, "]interface {}"):
				valTypStr = strings.Replace(valTypStr, "interface {}", "any", 1)
			}
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
