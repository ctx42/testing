// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ctx42/testing/pkg/dump"
)

// Types for some of the built-in types.
var (
	typTime    = reflect.TypeOf(time.Time{})
	typZone    = reflect.TypeOf(time.Location{})
	typZonePtr = reflect.TypeOf(&time.Location{})
	typByte    = reflect.TypeOf(byte(0))
)

// typeString returns a type of the value as a string.
func typeString(val reflect.Value) string {
	if !val.IsValid() {
		return "<invalid>"
	}
	return val.Type().String()
}

// isPrintableChar returns true if "v" is a printable ASCII character.
func isPrintableChar(v byte) bool {
	return v >= 32 && v <= 126
}

// valToString returns a string representation of the value.
//
// nolint: cyclop
func valToString(key reflect.Value) string {
	switch key.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", key.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return fmt.Sprintf("%d", key.Uint())

	case reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", key.Uint())

	case reflect.Float32:
		return strconv.FormatFloat(key.Float(), 'f', -1, 32)

	case reflect.Float64:
		return strconv.FormatFloat(key.Float(), 'f', -1, 64)

	case reflect.String:
		return fmt.Sprintf("%q", key.String())

	case reflect.Bool:
		return fmt.Sprintf("%v", key.Bool())

	case reflect.Struct:
		// For structs, we'll just print the type name
		// since the actual value might be complex.
		pkg := ""
		typ := key.Type()
		if pkgPath := typ.PkgPath(); pkgPath != "" {
			parts := strings.Split(pkgPath, "/")
			pkg = parts[len(parts)-1] + "."
		}
		return fmt.Sprintf("%s%s", pkg, typ.Name())

	case reflect.Ptr:
		if key.IsNil() {
			return dump.ValNil
		} else {
			return "*" + valToString(key.Elem())
		}

	case reflect.Complex64:
		return strconv.FormatComplex(key.Complex(), 'f', -1, 64)

	case reflect.Complex128:
		return strconv.FormatComplex(key.Complex(), 'f', -1, 64)

	case reflect.Array:
		return "<array>"

	case reflect.UnsafePointer:
		return fmt.Sprintf("<%p>", key.UnsafePointer())

	default:
		return "<invalid>"
	}
}
