// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package affirm

import (
	"reflect"
)

// isComparable reports whether the dynamic value of v can be compared with ==.
func isComparable(v any) bool {
	val := reflect.ValueOf(v)
	return !val.IsValid() || val.Comparable()
}
