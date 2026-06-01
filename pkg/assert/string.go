// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Contain asserts that "want" is a substring of "have" using [check.Contain].
//
// See the package documentation for the overall design and option handling.
func Contain(t tester.T, want, have string, opts ...any) bool {
	t.Helper()
	if e := check.Contain(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotContain asserts that "want" is not a substring of "have" using
// [check.NotContain].
//
// See the package documentation for the overall design and option handling.
func NotContain(t tester.T, want, have string, opts ...any) bool {
	t.Helper()
	if e := check.NotContain(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
