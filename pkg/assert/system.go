// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// ExitCode asserts that "err" is an [*exec.ExitError] with the given exit
// code.
//
// See [check.ExitCode] for the error-returning form.
func ExitCode(t tester.T, want int, err error, opts ...any) bool {
	t.Helper()
	if e := check.ExitCode(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
