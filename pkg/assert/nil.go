// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Nil asserts that "have" is nil using [check.Nil].
//
// See the Design section in the root README for the layered assert/check/notice
// architecture. Errors are built with [notice] and can be customized with
// options from [check] and [dump].
func Nil(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if e := check.Nil(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotNil asserts that "have" is not nil using [check.NotNil].
//
// See the Design section in the root README for the layered assert/check/notice
// architecture. Errors are built with [notice] and can be customized with
// options from [check] and [dump].
func NotNil(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if e := check.NotNil(have, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}
