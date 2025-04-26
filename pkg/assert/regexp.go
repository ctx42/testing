// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Regexp asserts that "want" regexp matches a string representation of "have.
// Returns true if it is, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
//
// The "want" can be either a regular expression string or instance of
// [regexp.Regexp]. The [fmt.Sprint] s used to get string representation of
// have argument.
func Regexp(t tester.T, want, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Regexp(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
