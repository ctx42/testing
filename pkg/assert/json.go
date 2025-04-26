// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// JSON asserts that two JSON strings are equivalent. Returns true if they are,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
//
// Example:
//
//	assert.JSON(t, `{"hello": "world"}`, `{"foo": "bar"}`)
func JSON(t tester.T, want, have string, opts ...check.Option) bool {
	t.Helper()
	if e := check.JSON(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
