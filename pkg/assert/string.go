// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Contain asserts "want" is a substring of "have". Returns true if it's,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func Contain(t tester.T, want, have string, opts ...any) bool {
	t.Helper()
	if e := check.Contain(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotContain asserts "want" is not a substring of "have". Returns true if it's
// not, otherwise marks the test as failed, writes an error message to the test
// log and returns false.
func NotContain(t tester.T, want, have string, opts ...any) bool {
	t.Helper()
	if e := check.NotContain(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
