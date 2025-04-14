package embedded

// THIS FILE HAS BEEN GENERATED - DO NOT EDIT!

import (
	"github.com/ctx42/testing/pkg/mock"
	"github.com/ctx42/testing/pkg/tester"
)

type EmbedLocalAndStdMock struct {
	*mock.Mock
	t tester.T
}

func NewEmbedLocalAndStdMock(t tester.T) *EmbedLocalAndStdMock {
	t.Helper()
	return &EmbedLocalAndStdMock{Mock: mock.NewMock(t), t: t}
}

func (_mck *EmbedLocalAndStdMock) Method0() {
	_mck.t.Helper()
	var _args []any
	_mck.Called(_args...)
}

func (_mck *EmbedLocalAndStdMock) Method1() error {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 error
	if _rFn, ok := _rets.Get(0).(func() error); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(error)
	}
	return _r0
}

func (_mck *EmbedLocalAndStdMock) Method2(a int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *EmbedLocalAndStdMock) Close() error {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 error
	if _rFn, ok := _rets.Get(0).(func() error); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(error)
	}
	return _r0
}

func (_mck *EmbedLocalAndStdMock) AAA() string {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 string
	if _rFn, ok := _rets.Get(0).(func() string); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(string)
	}
	return _r0
}
