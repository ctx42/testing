// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Same asserts that "want" and "have" are pointers to the same object using
// [check.Same].
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
//
// See the Design section in the root README for the layered assert/check/notice
// architecture. Errors are built with [notice] and can be customized with
// options from [check] and [dump].
func Same(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.Same(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotSame asserts that "want" and "have" are not pointers to the same object
// using [check.NotSame].
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
//
// See the Design section in the root README for the layered assert/check/notice
// architecture. Errors are built with [notice] and can be customized with
// options from [check] and [dump].
func NotSame(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.NotSame(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
