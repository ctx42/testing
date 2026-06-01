// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package assert provides a rich set of assertion functions designed for use
// in Go tests.
//
// See the Design section in the root README for the overall layered
// architecture (assert built on check built on notice).
//
// The customization model supports both global defaults and fine-grained
// control per assertion:
//   - Global type checkers via [check.RegisterTypeChecker].
//   - Per-assertion options such as [check.WithTrail], [check.WithTypeChecker],
//     and [dump] options.
//   - Full control via [check.DefaultOptions] and [check.WithOptions].
//
// All assertion functions follow a consistent pattern:
//   - Accept [tester.T] first.
//   - Return true on success.
//   - Report a structured error (via t.Error or t.Fatal) and return false
//     on failure.
//
// Errors are built with [notice] and can be customized with options from
// [check] and [dump].
//
// See individual function documentation for behavior specific to each
// assertion.
package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Count asserts that there are exactly "count" occurrences of "what" in
// "where".
//
// Currently only strings are supported for "where".
func Count(t tester.T, count int, what, where any, opts ...any) bool {
	t.Helper()
	if e := check.Count(count, what, where, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// SameType asserts that both arguments are of the same type.
//
// Type comparison uses [reflect.TypeOf] equality.
func SameType[T any](t tester.T, want T, have any, opts ...any) (T, bool) {
	t.Helper()
	h, e := check.SameType(want, have, opts...)
	if e != nil {
		t.Fatal(e)
		return h, false
	}
	return h, true
}

// NotSameType asserts that the arguments are not of the same type.
//
// Type comparison uses [reflect.TypeOf] equality.
func NotSameType(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.NotSameType(want, have, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// Type asserts that "have" is assignable to the pointer "want"
// (equivalent to `ok := src.(target)`).
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
// of fields.
func Fields(t tester.T, want int, s any, opts ...any) bool {
	t.Helper()
	if e := check.Fields(want, s, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
