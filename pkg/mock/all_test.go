package mock

import (
	"errors"
)

// ExampleImpl is a test implementation of ExampleItf
type ExampleImpl struct {
	*Mock
}

// NewExampleImpl returns new instance of ExampleImpl which uses
// given Mock instance.
func NewExampleImpl(mck *Mock) *ExampleImpl {
	return &ExampleImpl{Mock: mck}
}

func (imp *ExampleImpl) MethodInts(a, b, c int) (int, error) {
	rets := imp.Called(a, b, c)
	return rets.Int(0), errors.New("whoops")
}

func (imp *ExampleImpl) MethodBool(val bool) {
	imp.Called(val)
}

func (imp *ExampleImpl) MethodPtr(et *ExampleType) error {
	rets := imp.Called(et)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodItf(v ExampleItf) error {
	rets := imp.Called(v)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodChan(ch chan struct{}) error {
	rets := imp.Called(ch)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodMap(m map[string]bool) error {
	rets := imp.Called(m)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodBoolS(slice []bool) error {
	rets := imp.Called(slice)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodFunc(fn func(string) error) error {
	rets := imp.Called(fn)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodIntVar(a ...int) error {
	var args []any
	for _, v := range a {
		args = append(args, v)
	}
	rets := imp.Called(args...)
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodAnyVar(a ...any) error {
	rets := imp.Called(a...) // nolint: asasalint
	return rets.Error(0)
}

func (imp *ExampleImpl) MethodIntIntVar(a int, b ...int) error {
	args := []any{a}
	for _, v := range b {
		args = append(args, v)
	}
	rets := imp.Called(args...)
	return rets.Error(0)
}

type ExampleFuncType func(string) error

func (imp *ExampleImpl) MethodFuncType(fn ExampleFuncType) error {
	rets := imp.Called(fn)
	return rets.Error(0)
}

// -----------------------------------------------------------------------------

// ExampleItf represents test interface.
type ExampleItf interface{ HasRan() bool }

// ExampleType represents test type implementing ExampleItf.
type ExampleType struct{ ran bool }

// HasRun satisfies ExampleItf interface.
func (et *ExampleType) HasRan() bool { return et.ran }
