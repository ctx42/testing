// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package examples

import (
	"testing"

	"github.com/ctx42/testing/pkg/tester"
)

func Test_IsOdd(t *testing.T) {
	t.Run("error - is not an odd number", func(t *testing.T) {
		// --- Given ---

		// Set up the spy with expectations
		tspy := tester.New(t)
		tspy.ExpectError()                              // Expect an error.
		tspy.ExpectLogEqual("expected %d to be odd", 2) // Expect log.
		tspy.Close()                                    // No more expectations.

		// --- When ---
		success := IsOdd(tspy, 2) // Run the helper.

		// --- Then ---
		if success { // Verify the outcome.
			t.Error("expected success to be false")
		}
		tspy.AssertExpectations() // Ensure all expectations were met.
	})

	t.Run("is an odd number", func(t *testing.T) {
		// Given
		tspy := tester.New(t)
		tspy.Close()

		// When
		success := IsOdd(tspy, 3)

		// Then
		if !success {
			t.Error("expected success to be true")
		}

		// The `tspy.AssertExpectations()` is called automatically.
	})
}
