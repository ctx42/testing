// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"cmp"
	"fmt"
	"reflect"
	"strings"
)

// trimIndent removes the indentation [Printer.Tab] adds at level lvl from
// the start of str, leaving any leading whitespace of the value itself.
func trimIndent(dmp Dump, lvl int, str string) string {
	if dmp.Flat {
		return str
	}
	ind := strings.Repeat(" ", (dmp.Indent+lvl)*dmp.TabWidth)
	return strings.TrimPrefix(str, ind)
}

// typeName returns the name of typ, using "any" for the empty interface when
// [Dump.UseAny] is set.
func typeName(dmp Dump, typ reflect.Type) string {
	name := typ.String()
	if dmp.UseAny {
		name = strings.ReplaceAll(name, "interface {}", "any")
	}
	return name
}

// valueCmp compares two map keys for sorting. It returns a negative number
// when a sorts before b, zero when they are equal, and a positive number
// otherwise. Keys of interface type are compared by their dynamic values,
// ordered by kind and type first.
//
// nolint: cyclop
func valueCmp(a, b reflect.Value) int {
	if a.Kind() == reflect.Interface {
		a = a.Elem()
	}
	if b.Kind() == reflect.Interface {
		b = b.Elem()
	}
	if !a.IsValid() || !b.IsValid() {
		return cmp.Compare(boolInt(a.IsValid()), boolInt(b.IsValid()))
	}
	if a.Kind() != b.Kind() {
		return cmp.Compare(a.Kind(), b.Kind())
	}
	if a.Type() != b.Type() {
		return cmp.Compare(a.Type().String(), b.Type().String())
	}

	switch a.Kind() {
	case reflect.Bool:
		av, bv := a.Bool(), b.Bool()
		if av == bv {
			return 0
		}
		if av {
			return 1
		}
		return -1

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		av, bv := a.Int(), b.Int()
		if av == bv {
			return 0
		}
		if av < bv {
			return -1
		}
		return 1

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:

		av, bv := a.Uint(), b.Uint()
		if av == bv {
			return 0
		}
		if av < bv {
			return -1
		}
		return 1

	case reflect.Float32, reflect.Float64:
		// Orders NaN before any other value.
		return cmp.Compare(a.Float(), b.Float())

	case reflect.Complex64, reflect.Complex128:
		av, bv := a.Complex(), b.Complex()
		if c := cmp.Compare(real(av), real(bv)); c != 0 {
			return c
		}
		return cmp.Compare(imag(av), imag(bv))

	case reflect.String:
		return cmp.Compare(a.String(), b.String())

	default:
		// Structs, arrays, pointers and channels: compare their formatted
		// values, which fmt renders even for unexported fields.
		return cmp.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
	}
}

// boolInt returns 1 for true and 0 for false.
func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
