package examples

import (
	"testing"

	"github.com/ctx42/testing/pkg/mock"
)

// Adder is the simple interface we would like to mock.
type Adder interface {
	Add(a, b float64) float64
}

// AdderMock implements Adder interface. By convention the name of the type
// should be the name of the interface it is mocking with "Mock" suffix.
type AdderMock struct {
	*mock.Mock // Embedded instance.

	// If you need additional fields you can declare them here.
}

// NewAdderMock is the constructor function we will use to instantiate the mock.
// By convention, it should be the name of the type prefixed with "New".
func NewAdderMock(t *testing.T) *AdderMock {
	return &AdderMock{mock.NewMock(t)}
}

// Add mocks the method from Adder interface.
func (_mck *AdderMock) Add(a, b float64) float64 {
	// Inform the mock the method was called with given arguments.
	// The call returns [mock.Arguments] representing return values.
	args := _mck.Called(a, b)

	// Here you can do additional logic if needed.

	// Get the first return value, cast it to expected type, and return it.
	return args.Get(0).(float64)
}

func Test_Adder_Add(t *testing.T) {
	// --- Given ---
	mck := NewAdderMock(t) // Instantiate the mock.
	mck.
		On("Add", 1.0, 2.0). // Define method and argument expectations.
		Return(3.0)          // Define return values expectations.

	// --- When ---

	// In real example the mock created above would be used in code requiring
	// Adder interface, which in turn would call Add method like below.
	have := mck.Add(1.0, 2.0)

	// --- Then ---
	// Below line will print: Result: 3.000000
	t.Logf("Result: %f", have)
}
