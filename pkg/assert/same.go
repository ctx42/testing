// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Same asserts "want" and "have" are generic pointers and that both reference
// the same object. Returns true if they are, otherwise marks the test as
// failed, writes an error message to the test log and returns false.
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func Same(t tester.T, want, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Same(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotSame asserts "want" and "have" are generic pointers and that both do not
// reference the same object. Returns true if they are not, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func NotSame(t tester.T, want, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.NotSame(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
