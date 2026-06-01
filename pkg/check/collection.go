// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Len checks that "have" has "want" length.
//
// See the package documentation for option handling.
func Len(want int, have any, opts ...any) (err error) {
	vv := reflect.ValueOf(have)
	defer func() {
		if e := recover(); e != nil {
			msg := notice.New("cannot execute len(%T)", have).MetaSet("len", 0)
			ops := DefaultOptions(opts...)
			err = AddRows(ops, msg)
		}
	}()
	cnt := vv.Len()
	if want != cnt {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected %T length", have).
			Want("%d", want).
			Have("%d", cnt).
			MetaSet("len", cnt)
		return AddRows(ops, msg)
	}
	return nil
}

// Cap checks that "have" has "want" capacity.
//
// See the package documentation for option handling.
func Cap(want int, have any, opts ...any) (err error) {
	vv := reflect.ValueOf(have)
	defer func() {
		if e := recover(); e != nil {
			msg := notice.New("cannot execute cap(%T)", have).MetaSet("cap", 0)
			ops := DefaultOptions(opts...)
			err = AddRows(ops, msg)
		}
	}()
	cnt := vv.Cap()
	if want != cnt {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected %T capacity", have).
			Want("%d", want).
			Have("%d", cnt).
			MetaSet("cap", cnt)
		return AddRows(ops, msg)
	}
	return nil
}

// Has checks that the slice "bag" contains the value "want".
//
// See [assert.Has] for the assertion form and the package documentation for
// option handling ([DefaultOptions], [WithTrail], custom checkers, etc.).
func Has[T comparable](want T, bag []T, opts ...any) error {
	if slices.Contains(bag, want) {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected slice to have a value").
		Want("%#v", want).
		Append("slice", "%s", ops.Dumper.Any(bag))
	return AddRows(ops, msg)
}

// HasNo checks that the slice "set" does not contain the value "want".
//
// See [assert.HasNo] for the assertion form.
func HasNo[T comparable](want T, set []T, opts ...any) error {
	for i, got := range set {
		if want == got {
			ops := DefaultOptions(opts...)
			msg := notice.New("expected slice not to have value").
				Want("%#v", want).
				Append("index", "%d", i).
				Append("slice", "%s", ops.Dumper.Any(set))
			return AddRows(ops, msg)
		}
	}
	return nil
}

// HasKey checks that the map "set" contains the key "key".
//
// On success it returns the corresponding value and nil.
// On failure it returns the zero value for V and a descriptive error.
func HasKey[K comparable, V any](key K, set map[K]V, opts ...any) (V, error) {
	val, ok := set[key]
	if ok {
		return val, nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected map to have a key").
		Append("key", "%#v", key).
		Append("map", "%s", ops.Dumper.Any(set))
	return val, AddRows(ops, msg)
}

// HasNoKey checks that the map "set" does not contain the key "key".
//
// See [assert.HasNoKey] for the assertion wrapper.
func HasNoKey[K comparable, V any](key K, set map[K]V, opts ...any) error {
	val, ok := set[key]
	if !ok {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected map not to have a key").
		Append("key", "%#v", key).
		Append("value", "%#v", val).
		Append("map", "%s", ops.Dumper.Any(set))
	return AddRows(ops, msg)
}

// HasKeyValue checks that the map contains "key" mapped to exactly "want".
//
// See [assert.HasKeyValue] for the assertion wrapper.
func HasKeyValue[K, V comparable](
	key K,
	want V,
	set map[K]V,
	opts ...any,
) error {

	have, err := HasKey(key, set, opts...)
	if err != nil {
		return err
	}
	if err = Equal(want, have, opts...); err == nil {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected map to have a key with a value").
		Append("key", "%#v", key).
		Want("%#v", want).
		Have("%#v", have)
	return AddRows(ops, msg)
}

// SliceSubset checks that every element in "want" also appears in "have"
// (i.e. "want" is a subset of "have").
//
// See [assert.SliceSubset] for the corresponding assertion.
func SliceSubset[V comparable](want, have []V, opts ...any) error {
	var missing []V
	for _, wantVal := range want {
		found := slices.Contains(have, wantVal)
		if !found {
			missing = append(missing, wantVal)
		}
	}
	if len(missing) == 0 {
		return nil
	}

	ops := DefaultOptions(opts...)
	const hHeader = "expected \"want\" slice to be a subset of \"have\" slice"
	msg := notice.New(hHeader).
		Append("missing values", "%s", ops.Dumper.Any(missing))
	return AddRows(ops, msg)
}

// MapSubset checks that "want" is a subset of "have". All keys and values from
// "want" must exist in "have". Extra keys in "have" are allowed.
// See [assert.MapSubset] for the assertion form.
func MapSubset[K comparable, V any](want, have map[K]V, opts ...any) error {
	ops := DefaultOptions(opts...)

	var err error
	var missing []string
	for wKey, wVal := range want {
		wKeyStr := valToString(reflect.ValueOf(wKey))
		hVal, exist := have[wKey]
		if !exist {
			missing = append(missing, wKeyStr)
			continue
		}
		kOps := ops.MapTrail(wKeyStr)
		if e := Equal(wVal, hVal, WithOptions(kOps)); e != nil {
			err = notice.Join(err, e)
		}
	}

	if err != nil {
		err = notice.SortNotices(notice.From(err).Head(), notice.TrailCmp)
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		msg := notice.New("expected the map to have keys").
			Append("keys", "%s", strings.Join(missing, ", "))
		msg = AddRows(ops, msg)
		err = notice.Join(err, msg)
	}
	return err
}

// MapsSubset checks that every map in "want" is a subset of the corresponding
// map in "have" (using [MapSubset]).
//
// See [assert.MapsSubset] for the assertion form.
func MapsSubset[K comparable, V any](want, have []map[K]V, opts ...any) error {
	ops := DefaultOptions(opts...)
	if len(want) != len(have) {
		const hHeader = "expected slices of the same length"
		msg := notice.New(hHeader).Want("%d", len(want)).Have("%d", len(have))
		return AddRows(ops, msg)
	}

	var err error
	for i := range want {
		iOps := ops.ArrTrail("slice", i)
		if e := MapSubset(want[i], have[i], WithOptions(iOps)); e != nil {
			err = notice.Join(err, e)
		}
	}
	return err
}
