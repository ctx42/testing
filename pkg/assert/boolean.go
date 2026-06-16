// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// True asserts that "have" is true using [check.True].
func True(t tester.T, have bool, opts ...any) bool {
	t.Helper()
	if e := check.True(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// False asserts that "have" is false using [check.False].
func False(t tester.T, have bool, opts ...any) bool {
	t.Helper()
	if err := check.False(have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
