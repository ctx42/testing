// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Empty asserts "have" is empty. Returns true if it's, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
//
// See [check.Empty] for the list of values which are considered empty.
func Empty(t tester.T, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Empty(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NotEmpty asserts "have" is not empty. Returns true if it is not, otherwise
// marks the test as failed, writes an error message to the test log and
// returns false.
//
// See [check.Empty] for the list of values which are considered empty.
func NotEmpty(t tester.T, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.NotEmpty(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
