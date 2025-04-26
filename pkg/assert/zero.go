// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Zero asserts "have" is the zero value for its type. Returns true if it is,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func Zero(t tester.T, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Zero(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotZero asserts "have" is not the zero value for its type. Returns true if
// it is not, otherwise marks the test as failed, writes an error message to
// the test log and returns false.
func NotZero(t tester.T, have any, opts ...check.Option) bool {
	t.Helper()
	if err := check.NotZero(have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
