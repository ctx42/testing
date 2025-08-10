// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"reflect"
	"sort"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Len checks "have" has "want" length. Returns nil if it has, otherwise it
// returns an error with a message indicating the expected and actual values.
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

// Cap checks "have" has "want" capacity. Returns nil if it has, otherwise it
// returns an error with a message indicating the expected and actual values.
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

// Has checks slice has "want" value. Returns nil if it does, otherwise it
// returns an error with a message indicating the expected and actual values.
func Has[T comparable](want T, bag []T, opts ...any) error {
	for _, got := range bag {
		if want == got {
			return nil
		}
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected slice to have a value").
		Want("%#v", want).
		Append("slice", "%s", ops.Dumper.Any(bag))
	return AddRows(ops, msg)
}

// HasNo checks slice does not have the "want" value. Returns nil if it doesn't,
// otherwise it returns an error with a message indicating the expected and
// actual values.
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

// HasKey checks the map has a key. If the key exists, it returns its value and
// nil, otherwise it returns zero-value and an error with a message indicating
// the expected and actual values.
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

// HasNoKey checks map has no key. Returns nil if it doesn't, otherwise it
// returns an error with a message indicating the expected and actual values.
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

// HasKeyValue checks the map has a key with a given value. Returns nil if it
// doesn't, otherwise it returns an error with a message indicating the
// expected and actual values.
func HasKeyValue[K, V comparable](key K, want V, set map[K]V, opts ...any) error {
	have, err := HasKey(key, set, opts...)
	if err != nil {
		return err
	}
	if want == have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected map to have a key with a value").
		Append("key", "%#v", key).
		Want("%#v", want).
		Have("%#v", have)
	return AddRows(ops, msg)
}

// SliceSubset checks the "have" is a subset "want". In other words, all values
// in the "want" slice must be in the "have" slice. Returns nil if it does,
// otherwise returns an error with a message indicating the expected and actual
// values.
func SliceSubset[V comparable](want, have []V, opts ...any) error {
	var missing []V
	for _, wantVal := range want {
		found := false
		for _, haveVal := range have {
			if wantVal == haveVal {
				found = true
				break
			}
		}
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

// MapSubset checks the "want" is a subset "have". In other words, all keys and
// their corresponding values in the "want" map must be in the "have" map. It
// is not an error when the "have" map has some other keys. Returns nil if
// "want" is a subset of "have", otherwise it returns an error with a message
// indicating the expected and actual values.
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
		err = notice.SortNotices(notice.From(err).Head(), notice.TrialCmp)
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

// MapsSubset checks all the "want" maps are subsets of corresponding "have"
// maps using [MapSubset]. Returns nil if all "want" maps are subset of
// corresponding "have" maps, otherwise it returns an error with a message
// indicating the expected and actual values.
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
