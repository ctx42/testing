// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package assert provides assertion functions.
package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Count asserts there is "count" occurrences of "what" in "where". Returns
// true if the count matches, otherwise marks the test as failed, writes an
// error message to the test log and returns false.
//
// Currently, only strings are supported.
func Count(t tester.T, count int, what, where any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Count(count, what, where, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Type asserts that both arguments are of the same type. Returns true if
// they are, otherwise marks the test as failed, writes an error message to the
// test log and returns false.
func Type(t tester.T, want, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Type(want, have, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// Fields asserts struct or pointer to a struct "s" has "want" number of
// fields. Returns true if it does, otherwise marks the test as failed, writes
// an error message to the test log and returns false.
func Fields(t tester.T, want int, s any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Fields(want, s, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
