// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Empty asserts that "have" is empty.
//
// See [check.Empty] for the error-returning form and the package
// documentation for option handling.
func Empty(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if e := check.Empty(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotEmpty asserts that "have" is not empty.
//
// See [check.NotEmpty] and [check.Empty] for the error-returning forms.
func NotEmpty(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if e := check.NotEmpty(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
