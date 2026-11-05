// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"cmp"
	"reflect"
	"slices"
)

// mapEntry is one key and value of a map dumped by [MapDumper].
type mapEntry struct {
	key    reflect.Value
	val    reflect.Value
	keyDmp string // Cached result of keyStr.
	hasKey bool   // Whether keyDmp holds the cached result.
	str    string // Cached result of valStr.
	has    bool   // Whether str holds the cached result.
}

// keyStr returns the dump of the key. It dumps the key at most once, because
// dumps of cyclic values are costly.
func (ent *mapEntry) keyStr(dmp Dump, lvl int) string {
	if !ent.hasKey {
		ent.keyDmp, _ = dmp.value(lvl, ent.key)
		ent.hasKey = true
	}
	return ent.keyDmp
}

// valStr returns the dump of the value. It dumps the value at most once,
// because dumps of cyclic values are costly.
func (ent *mapEntry) valStr(dmp Dump, lvl int) string {
	if !ent.has {
		ent.str, _ = dmp.value(lvl, ent.val)
		ent.has = true
	}
	return ent.str
}

// takeValStr returns the result of valStr and releases the cached dump, so
// dumps of a large map are not all held at once.
func (ent *mapEntry) takeValStr(dmp Dump, lvl int) string {
	str := ent.valStr(dmp, lvl)
	ent.str = ""
	return str
}

// MapDumper is a generic dumper for maps. It expects val to represent the
// [reflect.Map] kind. Returns [ValErrUsage] ("<dump-usage-error>") string if
// the kind cannot be matched. It returns the string representation in the
// format defined by [Dump] configuration.
//
//nolint:cyclop
func MapDumper(dmp Dump, lvl int, val reflect.Value) string {
	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)

	if val.Kind() != reflect.Map {
		return prn.Write(ValErrUsage).String()
	}

	if dmp.PrintType {
		prn.Write(typeName(dmp, val.Type()))
	}

	if val.IsNil() {
		return prn.Write("(nil)").String()
	}

	dmp.PrintType = false // Don't print types for map values.

	// Values come from the iterator: a NaN key is not found by lookup.
	entries := make([]*mapEntry, 0, val.Len())
	for itr := val.MapRange(); itr.Next(); {
		entries = append(entries, &mapEntry{key: itr.Key(), val: itr.Value()})
	}
	slices.SortStableFunc(entries, func(a, b *mapEntry) int {
		if c := valueCmp(a.key, b.key); c != 0 {
			return c
		}
		// Keys without a natural order are ordered by their depth-limited
		// dump, and keys that dump the same, such as NaN, by value.
		c := cmp.Compare(a.keyStr(dmp, lvl+1), b.keyStr(dmp, lvl+1))
		if c != 0 {
			return c
		}
		return cmp.Compare(a.valStr(dmp, lvl+1), b.valStr(dmp, lvl+1))
	})

	num := val.Len()
	prn.Write("{").NLI(num)

	for i, ent := range entries {
		last := i == num-1

		prn.Write(ent.keyStr(dmp, lvl+1))
		prn.Write(":").Space()

		prn.Write(trimIndent(dmp, lvl+1, ent.takeValStr(dmp, lvl+1)))
		prn.Comma(last).Sep(last).NL()
	}
	if num > 0 {
		prn.Tab(dmp.Indent + lvl)
	}
	prn.Write("}")

	return prn.String()
}
