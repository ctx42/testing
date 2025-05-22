// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
)

// sliceDumper requires val to be dereferenced representation of [reflect.Slice]
// and returns its string representation in the format defined by [Dump]
// configuration.
func sliceDumper(dmp Dump, lvl int, val reflect.Value) string {
	if val.IsNil() {
		return ValNil
	}
	return arrayDumper(dmp, lvl, val)
}
