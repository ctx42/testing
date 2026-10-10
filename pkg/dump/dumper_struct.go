// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
)

// StructDumper is a generic dumper for maps. It expects val to represent the
// [reflect.Struct] kind. Returns [ValErrUsage] ("<dump-usage-error>") string
// if the kind cannot be matched. It returns the string representation in the
// format defined by [Dump] configuration.
func StructDumper(dmp Dump, lvl int, val reflect.Value) string {
	dmp.flatStrings = true

	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)

	if val.Kind() != reflect.Struct {
		return prn.Write(ValErrUsage).String()
	}

	vTyp := val.Type()

	num := val.NumField() // Total number of fields.
	lastPrivate := false
	var shown int // Number of fields printed.
	prn.Write("{").NLI(num)

	for i := range num {
		last := i == num-1

		fld := vTyp.Field(i)

		if !fld.IsExported() && !dmp.PrintPrivate {
			lastPrivate = last
			continue
		}

		// Field name.
		shown++
		prn.Tab(dmp.Indent + lvl + 1)
		prn.Write(fld.Name)
		prn.Write(":").Space()

		// Field value.
		dmp.PrintType = true
		sub, _ := dmp.value(lvl+1, val.Field(i))
		sub = trimIndent(dmp, lvl+1, sub)

		prn.Write(sub)
		prn.Comma(last).Sep(last).NL()
	}

	if lastPrivate && dmp.Flat {
		return strings.TrimRight(prn.String(), ", ") + "}"
	}
	if shown > 0 {
		prn.Tab(dmp.Indent + lvl)
	}
	prn.Write("}")

	return prn.String()
}
