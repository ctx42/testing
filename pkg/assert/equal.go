// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Equal asserts both values are equal. Returns true if they are, otherwise
// marks the test as failed, writes an error message to the test log and
// returns false.
func Equal(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if err := check.Equal(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// NotEqual asserts both values are not equal. Returns true if they are not,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func NotEqual(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if err := check.NotEqual(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
