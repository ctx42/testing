// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
)

// SliceDumper is a generic dumper for slices. It expects val to represent the
// [reflect.Slice] kind. Returns [valErrUsage] ("<dump-usage-error>") string if
// the kind cannot be matched. It returns string representation in the format
// defined by [Dump] configuration.
func SliceDumper(dmp Dump, lvl int, val reflect.Value) string {
	if val.Kind() != reflect.Slice {
		prn := NewPrinter(dmp)
		return prn.Write(ValErrUsage).String()
	}
	if val.IsNil() {
		return ValNil
	}
	return ArrayDumper(dmp, lvl, val)
}
