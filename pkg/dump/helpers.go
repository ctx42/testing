// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
)

// valueCmp returns whether the first value should sort before the second one.
//
// nolint: cyclop
func valueCmp(a, b reflect.Value) int {
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
		av, bv := a.Float(), b.Float()
		if av == bv {
			return 0
		}
		if av < bv {
			return -1
		}
		return 1

	default:
		av, bv := a.String(), b.String()
		if av == bv {
			return 0
		}
		if av < bv {
			return -1
		}
		return 1
	}
}
