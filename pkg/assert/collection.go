// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"
)

// Len asserts that "have" has "want" length using [check.Len].
//
// See the Design section in the root README for the layered assert/check/notice
// architecture. Errors are built with [notice] and can be customized with
// options from [check] and [dump].
func Len(t tester.T, want int, have any, opts ...any) bool {
	t.Helper()
	if e := check.Len(want, have, opts...); e != nil {
		var cnt int
		if val, ok := notice.From(e).MetaLookup("len"); ok {
			cnt = val.(int) // nolint: forcetypeassert
		}
		if want > cnt {
			t.Fatal(e)
		} else {
			t.Error(e)
		}
		return false
	}
	return true
}

// Cap asserts that "have" has "want" capacity. See [check.Cap].
func Cap(t tester.T, want int, have any, opts ...any) bool {
	t.Helper()
	if e := check.Cap(want, have, opts...); e != nil {
		var cnt int
		if val, ok := notice.From(e).MetaLookup("cap"); ok {
			cnt = val.(int) // nolint: forcetypeassert
		}
		if want > cnt {
			t.Fatal(e)
		} else {
			t.Error(e)
		}
		return false
	}
	return true
}

// Has asserts that the slice contains the "want" value.
// See [check.Has] for the error-returning form and option details.
func Has[T comparable](t tester.T, want T, bag []T, opts ...any) bool {
	t.Helper()
	if e := check.Has(want, bag, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// HasNo asserts that the slice does not contain the "want" value.
// See [check.HasNo] for the error-returning form.
func HasNo[T comparable](t tester.T, want T, bag []T, opts ...any) bool {
	t.Helper()
	if e := check.HasNo(want, bag, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// HasKey asserts that the map contains the key. When the key exists, the
// associated value is also returned. See [check.HasKey].
func HasKey[K comparable, V any](
	t tester.T,
	key K,
	set map[K]V,
	opts ...any,
) (V, bool) {

	t.Helper()
	val, e := check.HasKey(key, set, opts...)
	if e != nil {
		t.Error(e)
		return val, false
	}
	return val, true
}

// HasNoKey asserts that the map does not contain the key.
// See [check.HasNoKey].
func HasNoKey[K comparable, V any](
	t tester.T,
	key K,
	set map[K]V,
	opts ...any,
) bool {

	t.Helper()
	if e := check.HasNoKey(key, set, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// HasKeyValue asserts that the map contains the key with the given value.
// See [check.HasKeyValue] for the error-returning form.
func HasKeyValue[K, V comparable](
	t tester.T,
	key K,
	want V,
	set map[K]V,
	opts ...any,
) bool {

	t.Helper()
	if e := check.HasKeyValue(key, want, set, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// SliceSubset asserts that "want" is a subset of "have". All values in "want"
// must be present in "have". See [check.SliceSubset].
func SliceSubset[T comparable](t tester.T, want, have []T, opts ...any) bool {
	t.Helper()
	if e := check.SliceSubset(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// MapSubset asserts that "want" is a subset of "have". All keys and values in
// "want" must be present in "have". It is not an error if "have" contains
// additional keys. See [check.MapSubset].
func MapSubset[K comparable, V any](
	t tester.T,
	want, have map[K]V,
	opts ...any,
) bool {

	t.Helper()
	if e := check.MapSubset(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// MapsSubset asserts that every map in "want" is a subset of the corresponding
// map in "have" (using [MapSubset]). See [check.MapsSubset].
func MapsSubset[K comparable, V any](
	t tester.T,
	want, have []map[K]V,
	opts ...any,
) bool {

	t.Helper()
	if e := check.MapsSubset(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
