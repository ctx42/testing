// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// MapDumper is a generic dumper for maps. It expects val to represent the
// [reflect.Map] kind. Returns [valErrUsage] ("<dump-usage-error>") string if
// the kind cannot be matched. It returns string representation in the format
// defined by [Dump] configuration.
//
// nolint: cyclop
func MapDumper(dmp Dump, lvl int, val reflect.Value) string {
	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)

	if val.Kind() != reflect.Map {
		return prn.Write(ValErrUsage).String()
	}

	if dmp.PrintType {
		keyTyp := val.Type().Key()
		valTyp := val.Type().Elem()
		valTypStr := strings.ReplaceAll(valTyp.String(), " ", "")
		if valTypStr == "interface{}" && dmp.UseAny {
			valTypStr = "any"
		}
		str := fmt.Sprintf("map[%s]%s", keyTyp.String(), valTypStr)
		prn.Write(str)
	}

	keys := val.MapKeys()
	slices.SortStableFunc(keys, valueCmp)

	if val.IsNil() {
		return prn.Write("(nil)").String()
	}

	num := val.Len()
	prn.Write("{").NLI(num)

	dmp.PrintType = false // Don't print types for map values.
	for i, key := range keys {
		last := i == num-1

		sub, _ := dmp.value(lvl+1, key)
		prn.Write(sub)
		prn.Write(":").Space()

		sub, _ = dmp.value(lvl+1, val.MapIndex(key))
		sub = strings.TrimLeft(sub, " \t")

		dmp.PrintType = true
		prn.Write(sub)
		prn.Comma(last).Sep(last).NL()
	}
	prn.Tab(dmp.Indent + lvl).Write("}")

	return prn.String()
}
