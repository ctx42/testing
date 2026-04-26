// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package assert provides assertion functions.
package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Count asserts there are "count" occurrences of "what" in "where". Returns
// true if the count matches, otherwise marks the test as failed, writes an
// error message to the test log, and returns false.
//
// Currently, only strings are supported.
func Count(t tester.T, count int, what, where any, opts ...any) bool {
	t.Helper()
	if e := check.Count(count, what, where, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// SameType asserts that both arguments are of the same type. Returns true if
// they are, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Assertion uses [reflect.TypeOf] equality to determine the type.
func SameType[T any](t tester.T, want T, have any, opts ...any) (T, bool) {
	t.Helper()
	h, e := check.SameType(want, have, opts...)
	if e != nil {
		t.Fatal(e)
		return h, false
	}
	return h, true
}

// NotSameType asserts that the arguments are not of the same type. Returns
// true if they are not, otherwise marks the test as failed, writes an error
// message to the test log, and returns false.
//
// Assertion uses [reflect.TypeOf] equality to determine the type.
func NotSameType(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.NotSameType(want, have, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// Type asserts that "have" can be type assigned to the pointer "want"
// (same as target, ok := src.(target)). Returns true if it can be done,
// otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Example:
//
//	var target int
//	var src any = 42
//	assert.Type(t, &target, src)
func Type(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.Type(want, have, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// Fields asserts that a struct or pointer to a struct "s" has "want" number
// of fields. Returns true if it does, otherwise marks the test as failed,
// writes an error message to the test log, and returns false.
func Fields(t tester.T, want int, s any, opts ...any) bool {
	t.Helper()
	if e := check.Fields(want, s, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
