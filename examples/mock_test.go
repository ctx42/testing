// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package examples

import (
	"testing"

	"github.com/ctx42/testing/pkg/mock"
)

// Adder is a simple interface to mock.
type Adder interface {
	Add(a, b float64) float64
}

// AdderMock implements the [Adder] interface. By convention, mock types are
// named after the interface with a "Mock" suffix.
type AdderMock struct {
	*mock.Mock // Embedded mock instance.

	// Add custom fields here if needed.
}

// NewAdderMock creates a new [AdderMock] instance. By convention, constructor
// functions are prefixed with "New". More complex mocks may accept additional
// parameters.
func NewAdderMock(t *testing.T) *AdderMock {
	return &AdderMock{mock.NewMock(t)}
}

// Add implements the Add method from the [Adder] interface.
func (mck *AdderMock) Add(a, b float64) float64 {
	// Record the method call with its arguments, returning [mock.Arguments]
	// containing the defined return values.
	rets := mck.Called(a, b)

	// Add custom logic here if needed.

	// Extract and return the first return value, cast to float64.
	return rets.Get(0).(float64)
}

// Test_Adder_Add demonstrates using AdderMock in a test.
func Test_Adder_Add(t *testing.T) {
	// --- Given ---
	mck := NewAdderMock(t) // Create the mock.
	mck.
		On("Add", 1.0, 2.0). // Specify expected method and arguments.
		Return(3.0)          // Define the return value.

	// --- When ---
	// In a real test, the mock would be passed to code requiring the Adder
	// interface, which would invoke Add as shown.
	have := mck.Add(1.0, 2.0)

	// --- Then ---
	// Prints: Result: 3.000000
	t.Logf("Result: %f", have)
}
