// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Wait waits for "fn" to return true, but no longer than the given timeout.
// The "within" may be a string, int, int64, or [time.Duration].
//
// Calls to "fn" are throttled (see [check.Options.WaitThrottle] and
// [check.WithWaitThrottle]).
//
// See [check.Wait] for the error-returning form.
func Wait(t tester.T, within any, fn func() bool, opts ...any) bool {
	t.Helper()
	if e := check.Wait(within, fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
