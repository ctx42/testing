// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package examples

// File contains example usages of tester.T and tester.Spy.

import (
	"github.com/ctx42/testing/pkg/tester"
)

// IsOdd asserts "have" is odd number. Returns true if it is, otherwise marks
// the test as failed, writes an error message to the test log and returns false.
func IsOdd(t tester.T, have int) bool {
	t.Helper()
	if have%2 == 0 {
		t.Errorf("expected %d to be odd", have)
		return false
	}
	return true
}
