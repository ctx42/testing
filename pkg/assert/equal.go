// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Equal asserts that want and have are equal using [check.Equal].
//
// See the package documentation for customization options
// ([check.WithTrail], [check.WithTypeChecker], [dump] options, etc.).
func Equal(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if err := check.Equal(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// NotEqual asserts that want and have are not equal using [check.NotEqual].
//
// See the package documentation for customization options.
func NotEqual(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if err := check.NotEqual(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
