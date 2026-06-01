// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Regexp asserts that "want" regexp matches a string representation of "have".
//
// "want" may be a regular expression string or a [*regexp.Regexp].
// [fmt.Sprint] produces the string form of "have".
//
// See [check.Regexp] for the error-returning form.
func Regexp(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.Regexp(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
