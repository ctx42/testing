// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"reflect"
	"sort"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Len checks "have" has "want" elements. Returns nil if it has, otherwise it
// returns an error with a message indicating the expected and actual values.
func Len(want int, have any, opts ...Option) (err error) {
	vv := reflect.ValueOf(have)
	defer func() {
		if e := recover(); e != nil {
			err = notice.New("cannot execute len(%T)", have).MetaSet("len", 0)
		}
	}()
	cnt := vv.Len()
	if want != cnt {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected %T length", have).
			SetTrail(ops.Trail).
			Want("%d", want).
			Have("%d", cnt).
			MetaSet("len", cnt)
		return msg
	}
	return nil
}

// Has checks slice has "want" value. Returns nil if it does, otherwise it
// returns an error with a message indicating the expected and actual values.
func Has[T comparable](want T, bag []T, opts ...Option) error {
	for _, got := range bag {
		if want == got {
			return nil
		}
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected slice to have a value").
		SetTrail(ops.Trail).
		Want("%#v", want).
		Append("slice", "%s", ops.Dumper.Any(bag))
}

// HasNo checks slice does not have the "want" value. Returns nil if it doesn't,
// otherwise it returns an error with a message indicating the expected and
// actual values.
func HasNo[T comparable](want T, set []T, opts ...Option) error {
	for i, got := range set {
		if want == got {
			ops := DefaultOptions(opts...)
			return notice.New("expected slice not to have value").
				SetTrail(ops.Trail).
				Want("%#v", want).
				Append("index", "%d", i).
				Append("slice", "%s", ops.Dumper.Any(set))
		}
	}
	return nil
}

// HasKey checks the map has a key. If the key exists, it returns its value and
// nil, otherwise it returns zero-value and an error with a message indicating
// the expected and actual values.
func HasKey[K comparable, V any](key K, set map[K]V, opts ...Option) (V, error) {
	val, ok := set[key]
	if ok {
		return val, nil
	}
	ops := DefaultOptions(opts...)
	return val, notice.New("expected map to have a key").
		SetTrail(ops.Trail).
		Append("key", "%#v", key).
		Append("map", "%s", ops.Dumper.Any(set))
}

// HasNoKey checks map has no key. Returns nil if it doesn't, otherwise it
// returns an error with a message indicating the expected and actual values.
func HasNoKey[K comparable, V any](key K, set map[K]V, opts ...Option) error {
	val, ok := set[key]
	if !ok {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected map not to have a key").
		SetTrail(ops.Trail).
		Append("key", "%#v", key).
		Append("value", "%#v", val).
		Append("map", "%s", ops.Dumper.Any(set))
}

// HasKeyValue checks the map has a key with a given value. Returns nil if it
// doesn't, otherwise it returns an error with a message indicating the
// expected and actual values.
func HasKeyValue[K, V comparable](key K, want V, set map[K]V, opts ...Option) error {
	have, err := HasKey(key, set, opts...)
	if err != nil {
		return err
	}
	if want == have {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected map to have a key with a value").
		SetTrail(ops.Trail).
		Append("key", "%#v", key).
		Want("%#v", want).
		Have("%#v", have)
}

// SliceSubset checks the "have" is a subset "want". In other words, all values
// in the "want" slice must be in the "have" slice. Returns nil if it does,
// otherwise returns an error with a message indicating the expected and actual
// values.
func SliceSubset[V comparable](want, have []V, opts ...Option) error {
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
	return notice.New(hHeader).
		SetTrail(ops.Trail).
		Append("missing values", "%s", ops.Dumper.Any(missing))
}

// MapSubset checks the "want" is a subset "have". In other words, all keys and
// their corresponding values in the "want" map must be in the "have" map. It
// is not an error when the "have" map has some other keys. Returns nil if
// "want" is a subset of "have", otherwise it returns an error with a message
// indicating the expected and actual values.
func MapSubset[K comparable, V any](want, have map[K]V, opts ...Option) error {
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
			SetTrail(ops.Trail).
			Append("keys", "%s", strings.Join(missing, ", "))
		err = notice.Join(err, msg)
	}
	return err
}

// MapsSubset checks all the "want" maps are subsets of corresponding "have"
// maps using [MapSubset]. Returns nil if all "want" maps are subset of
// corresponding "have" maps, otherwise it returns an error with a message
// indicating the expected and actual values.
func MapsSubset[K comparable, V any](want, have []map[K]V, opts ...Option) error {
	ops := DefaultOptions(opts...)
	if len(want) != len(have) {
		msg := "expected slices of the same length"
		return notice.New(msg).
			SetTrail(ops.Trail).
			Want("%d", len(want)).
			Have("%d", len(have))
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
