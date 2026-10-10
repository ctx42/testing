// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester_test

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/tester"
)

// isOdd is a test helper reporting an error when have is even.
func isOdd(t tester.T, have int) bool {
	t.Helper()
	if have%2 == 0 {
		t.Errorf("expected %d to be odd", have)
		return false
	}
	return true
}

func ExampleSpy() {
	spy := tester.New(&testing.T{})
	spy.ExpectError()
	spy.ExpectLogEqual("expected %d to be odd", 2)
	spy.Close()

	have := isOdd(spy, 2)

	fmt.Println(have, spy.AssertExpectations())
	// Output: false true
}
