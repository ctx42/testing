// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// JSON asserts that two JSON strings are equivalent.
//
// See [check.JSON] for the error-returning form and supported types
// (via the [check.Text] constraint).
func JSON[W, H check.Text](t tester.T, want W, have H, opts ...any) bool {
	t.Helper()
	if e := check.JSON(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
