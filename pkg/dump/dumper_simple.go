// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SimpleDumper is a generic dumper for simple types. It expects val to
// represent one of the kinds:
//
//   - [reflect.Bool]
//   - [reflect.Int]
//   - [reflect.Int8]
//   - [reflect.Int16]
//   - [reflect.Int32]
//   - [reflect.Int64]
//   - [reflect.Uint]
//   - [reflect.Uint8]
//   - [reflect.Uint16]
//   - [reflect.Uint32]
//   - [reflect.Uint64]
//   - [reflect.Uintptr]
//   - [reflect.Float32]
//   - [reflect.Float64]
//   - [reflect.String]
//
// It returns string representation in the format defined by [Dump]
// configuration.
func SimpleDumper(dmp Dump, lvl int, val reflect.Value) string {
	v := val.Interface()

	var format string
	switch val.Kind() {
	case reflect.String:
		length := val.Len()
		switch {
		case dmp.flatStrings:
			format = `%q`
		case dmp.Flat:
			format = `%q`
		case dmp.FlatStrings > 0 && length <= dmp.FlatStrings:
			format = `%q`
		case strings.Contains(val.String(), "\n"):
			format = "%v"
		default:
			format = "%q"
		}

	case reflect.Float32:
		format = "%s"
		f := float64(v.(float32)) // nolint: forcetypeassert
		v = strconv.FormatFloat(f, 'f', -1, 32)

	case reflect.Float64:
		format = "%s"
		f := v.(float64) // nolint: forcetypeassert
		v = strconv.FormatFloat(f, 'f', -1, 64)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		format = "%d"

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		format = "%d"

	default:
		format = "%v"
	}

	prn := NewPrinter(dmp)
	str := fmt.Sprintf(format, v)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}
