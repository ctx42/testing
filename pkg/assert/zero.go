// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Zero asserts that "have" is the zero value for its type using [check.Zero].
func Zero(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if e := check.Zero(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotZero asserts that "have" is not the zero value for its type using
// [check.NotZero].
func NotZero(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if err := check.NotZero(have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
