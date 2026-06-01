// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package examples contains demonstrations of advanced usage patterns.
package examples

import (
	"github.com/ctx42/testing/pkg/tester"
)

// IsOdd is an example test helper that demonstrates using [tester.T].
// It asserts that "have" is an odd number.
//
// See [tester_test.go] for how to test this helper using [tester.Spy].
func IsOdd(t tester.T, have int) bool {
	t.Helper()
	if have%2 == 0 {
		t.Errorf("expected %d to be odd", have)
		return false
	}
	return true
}
