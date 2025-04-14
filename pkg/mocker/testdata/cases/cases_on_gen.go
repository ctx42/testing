package cases

// THIS FILE HAS BEEN GENERATED - DO NOT EDIT!

import (
	"fmt"
	"io/fs"
	mt "time"

	"github.com/ctx42/testing/internal/mocker/testdata/pkga"
	"github.com/ctx42/testing/internal/mocker/testdata/pkgb"
	"github.com/ctx42/testing/internal/mocker/testdata/pkgc"
	"github.com/ctx42/testing/internal/mocker/testdata/pkgd"
	"github.com/ctx42/testing/internal/mocker/testdata/pkge"
	"github.com/ctx42/testing/pkg/mock"
	"github.com/ctx42/testing/pkg/tester"
)

type CasesMockOn struct {
	*mock.Mock
	t tester.T
}

func NewCasesMockOn(t tester.T) *CasesMockOn {
	t.Helper()
	return &CasesMockOn{Mock: mock.NewMock(t), t: t}
}

func (_mck *CasesMockOn) Method0() {
	_mck.t.Helper()
	var _args []any
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod0() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method0", _args...)
}

func (_mck *CasesMockOn) Method1(a int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod1(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method1", _args...)
}

func (_mck *CasesMockOn) Method2(a int, b int) {
	_mck.t.Helper()
	_args := []any{a, b}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod2(a any, b any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b}
	return _mck.On("Method2", _args...)
}

func (_mck *CasesMockOn) Method3(a int, b int) {
	_mck.t.Helper()
	_args := []any{a, b}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod3(a any, b any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b}
	return _mck.On("Method3", _args...)
}

func (_mck *CasesMockOn) Method4(a int, b int, c bool) {
	_mck.t.Helper()
	_args := []any{a, b, c}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod4(a any, b any, c any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b, c}
	return _mck.On("Method4", _args...)
}

func (_mck *CasesMockOn) Method5(_a0 int) {
	_mck.t.Helper()
	_args := []any{_a0}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod5(_a0 any) *mock.Call {
	_mck.t.Helper()
	_args := []any{_a0}
	return _mck.On("Method5", _args...)
}

func (_mck *CasesMockOn) Method6(a int, _a1 int, b bool) {
	_mck.t.Helper()
	_args := []any{a, _a1, b}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod6(a any, _a1 any, b any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, _a1, b}
	return _mck.On("Method6", _args...)
}

func (_mck *CasesMockOn) Method7() error {
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

func (_mck *CasesMockOn) OnMethod7() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method7", _args...)
}

func (_mck *CasesMockOn) Method8() error {
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

func (_mck *CasesMockOn) OnMethod8() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method8", _args...)
}

func (_mck *CasesMockOn) Method9() (error, error) {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 2 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 error
	if _rFn, ok := _rets.Get(0).(func() error); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(error)
	}
	var _r1 error
	if _rFn, ok := _rets.Get(1).(func() error); ok {
		_r1 = _rFn()
	} else if _r := _rets.Get(1); _r != nil {
		_r1 = _r.(error)
	}
	return _r0, _r1
}

func (_mck *CasesMockOn) OnMethod9() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method9", _args...)
}

func (_mck *CasesMockOn) Method10() (int, error) {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 2 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 int
	if _rFn, ok := _rets.Get(0).(func() int); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(int)
	}
	var _r1 error
	if _rFn, ok := _rets.Get(1).(func() error); ok {
		_r1 = _rFn()
	} else if _r := _rets.Get(1); _r != nil {
		_r1 = _r.(error)
	}
	return _r0, _r1
}

func (_mck *CasesMockOn) OnMethod10() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method10", _args...)
}

func (_mck *CasesMockOn) Method11(_a0 int, _a1 float64) {
	_mck.t.Helper()
	_args := []any{_a0, _a1}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod11(_a0 any, _a1 any) *mock.Call {
	_mck.t.Helper()
	_args := []any{_a0, _a1}
	return _mck.On("Method11", _args...)
}

func (_mck *CasesMockOn) Method12(a ...int) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod12(a ...any) *mock.Call {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	return _mck.On("Method12", _args...)
}

func (_mck *CasesMockOn) Method13(tim mt.Time) error {
	_mck.t.Helper()
	_args := []any{tim}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 error
	if _rFn, ok := _rets.Get(0).(func(mt.Time) error); ok {
		_r0 = _rFn(tim)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(error)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod13(tim any) *mock.Call {
	_mck.t.Helper()
	_args := []any{tim}
	return _mck.On("Method13", _args...)
}

func (_mck *CasesMockOn) Method14(_a0 func()) {
	_mck.t.Helper()
	_args := []any{_a0}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod14(_a0 any) *mock.Call {
	_mck.t.Helper()
	_args := []any{_a0}
	return _mck.On("Method14", _args...)
}

func (_mck *CasesMockOn) Method15(_a0 func(int)) {
	_mck.t.Helper()
	_args := []any{_a0}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod15(_a0 any) *mock.Call {
	_mck.t.Helper()
	_args := []any{_a0}
	return _mck.On("Method15", _args...)
}

func (_mck *CasesMockOn) Method16(a func(...int)) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod16(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method16", _args...)
}

func (_mck *CasesMockOn) Method17() Concrete {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 Concrete
	if _rFn, ok := _rets.Get(0).(func() Concrete); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(Concrete)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod17() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method17", _args...)
}

func (_mck *CasesMockOn) Method18() *Concrete {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 *Concrete
	if _rFn, ok := _rets.Get(0).(func() *Concrete); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(*Concrete)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod18() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method18", _args...)
}

func (_mck *CasesMockOn) Method19() pkga.A1 {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 pkga.A1
	if _rFn, ok := _rets.Get(0).(func() pkga.A1); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(pkga.A1)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod19() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method19", _args...)
}

func (_mck *CasesMockOn) Method20() *pkga.A1 {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 *pkga.A1
	if _rFn, ok := _rets.Get(0).(func() *pkga.A1); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(*pkga.A1)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod20() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method20", _args...)
}

func (_mck *CasesMockOn) Method21(a fmt.Stringer) fs.File {
	_mck.t.Helper()
	_args := []any{a}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 fs.File
	if _rFn, ok := _rets.Get(0).(func(fmt.Stringer) fs.File); ok {
		_r0 = _rFn(a)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(fs.File)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod21(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method21", _args...)
}

func (_mck *CasesMockOn) Method22(a Concrete) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod22(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method22", _args...)
}

func (_mck *CasesMockOn) Method23(a *Concrete) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod23(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method23", _args...)
}

func (_mck *CasesMockOn) Method24(a ...Concrete) int {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 int
	if _rFn, ok := _rets.Get(0).(func(...Concrete) int); ok {
		_r0 = _rFn(a...)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(int)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod24(a ...any) *mock.Call {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	return _mck.On("Method24", _args...)
}

func (_mck *CasesMockOn) Method25(a ...*Concrete) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod25(a ...any) *mock.Call {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	return _mck.On("Method25", _args...)
}

func (_mck *CasesMockOn) Method26(a ...pkga.A1) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod26(a ...any) *mock.Call {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	return _mck.On("Method26", _args...)
}

func (_mck *CasesMockOn) Method27(a ...*pkga.A1) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod27(a ...any) *mock.Call {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	return _mck.On("Method27", _args...)
}

func (_mck *CasesMockOn) Method28(a *int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod28(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method28", _args...)
}

func (_mck *CasesMockOn) Method29(a pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod29(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method29", _args...)
}

func (_mck *CasesMockOn) Method30(a *pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod30(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method30", _args...)
}

func (_mck *CasesMockOn) Method31(a [2]int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod31(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method31", _args...)
}

func (_mck *CasesMockOn) Method32(a [2]*int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod32(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method32", _args...)
}

func (_mck *CasesMockOn) Method33(a [2]pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod33(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method33", _args...)
}

func (_mck *CasesMockOn) Method34(a [2]*pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod34(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method34", _args...)
}

func (_mck *CasesMockOn) Method35(a []int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod35(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method35", _args...)
}

func (_mck *CasesMockOn) Method36(a []*int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod36(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method36", _args...)
}

func (_mck *CasesMockOn) Method37(a []pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod37(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method37", _args...)
}

func (_mck *CasesMockOn) Method38(a []*pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod38(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method38", _args...)
}

func (_mck *CasesMockOn) Method39(a map[int]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod39(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method39", _args...)
}

func (_mck *CasesMockOn) Method40(a map[int]*string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod40(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method40", _args...)
}

func (_mck *CasesMockOn) Method41(a map[*int]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod41(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method41", _args...)
}

func (_mck *CasesMockOn) Method42(a map[pkga.A1]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod42(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method42", _args...)
}

func (_mck *CasesMockOn) Method43(a map[*pkga.A1]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod43(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method43", _args...)
}

func (_mck *CasesMockOn) Method44(a map[*pkga.A1]*pkgb.B1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod44(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method44", _args...)
}

func (_mck *CasesMockOn) Method45(a chan map[*pkga.A1]*pkgb.B1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod45(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method45", _args...)
}

func (_mck *CasesMockOn) Method46(a map[*pkga.A1]*pkgb.B1, b pkgc.C1) *pkgd.D1 {
	_mck.t.Helper()
	_args := []any{a, b}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 *pkgd.D1
	if _rFn, ok := _rets.Get(0).(func(map[*pkga.A1]*pkgb.B1, pkgc.C1) *pkgd.D1); ok {
		_r0 = _rFn(a, b)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(*pkgd.D1)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod46(a any, b any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b}
	return _mck.On("Method46", _args...)
}

func (_mck *CasesMockOn) Method47(a func(func(mt.Time, *pkga.A1))) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod47(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method47", _args...)
}

func (_mck *CasesMockOn) Method48(a map[*pkga.A1]func(pkgb.B1) func(pkga.A1) error) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod48(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method48", _args...)
}

func (_mck *CasesMockOn) Method49(a <-chan *pkga.A1) chan<- int {
	_mck.t.Helper()
	_args := []any{a}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 chan<- int
	if _rFn, ok := _rets.Get(0).(func(<-chan *pkga.A1) chan<- int); ok {
		_r0 = _rFn(a)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(chan<- int)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod49(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method49", _args...)
}

func (_mck *CasesMockOn) Method50(a int) (int, int, error) {
	_mck.t.Helper()
	_args := []any{a}
	_rets := _mck.Called(_args...)
	if len(_rets) != 3 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 int
	if _rFn, ok := _rets.Get(0).(func(int) int); ok {
		_r0 = _rFn(a)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(int)
	}
	var _r1 int
	if _rFn, ok := _rets.Get(1).(func(int) int); ok {
		_r1 = _rFn(a)
	} else if _r := _rets.Get(1); _r != nil {
		_r1 = _r.(int)
	}
	var _r2 error
	if _rFn, ok := _rets.Get(2).(func(int) error); ok {
		_r2 = _rFn(a)
	} else if _r := _rets.Get(2); _r != nil {
		_r2 = _r.(error)
	}
	return _r0, _r1, _r2
}

func (_mck *CasesMockOn) OnMethod50(a any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a}
	return _mck.On("Method50", _args...)
}

func (_mck *CasesMockOn) Method51(e pkge.E1) error {
	_mck.t.Helper()
	_args := []any{e}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 error
	if _rFn, ok := _rets.Get(0).(func(pkge.E1) error); ok {
		_r0 = _rFn(e)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(error)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod51(e any) *mock.Call {
	_mck.t.Helper()
	_args := []any{e}
	return _mck.On("Method51", _args...)
}

func (_mck *CasesMockOn) Method52(a Concrete, b mt.Time, c int) (*Concrete, error) {
	_mck.t.Helper()
	_args := []any{a, b, c}
	_rets := _mck.Called(_args...)
	if len(_rets) != 2 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 *Concrete
	if _rFn, ok := _rets.Get(0).(func(Concrete, mt.Time, int) *Concrete); ok {
		_r0 = _rFn(a, b, c)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(*Concrete)
	}
	var _r1 error
	if _rFn, ok := _rets.Get(1).(func(Concrete, mt.Time, int) error); ok {
		_r1 = _rFn(a, b, c)
	} else if _r := _rets.Get(1); _r != nil {
		_r1 = _r.(error)
	}
	return _r0, _r1
}

func (_mck *CasesMockOn) OnMethod52(a any, b any, c any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b, c}
	return _mck.On("Method52", _args...)
}

func (_mck *CasesMockOn) Method53(a int, b bool) (int, bool, string, error) {
	_mck.t.Helper()
	_args := []any{a, b}
	_rets := _mck.Called(_args...)
	if len(_rets) != 4 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 int
	if _rFn, ok := _rets.Get(0).(func(int, bool) int); ok {
		_r0 = _rFn(a, b)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(int)
	}
	var _r1 bool
	if _rFn, ok := _rets.Get(1).(func(int, bool) bool); ok {
		_r1 = _rFn(a, b)
	} else if _r := _rets.Get(1); _r != nil {
		_r1 = _r.(bool)
	}
	var _r2 string
	if _rFn, ok := _rets.Get(2).(func(int, bool) string); ok {
		_r2 = _rFn(a, b)
	} else if _r := _rets.Get(2); _r != nil {
		_r2 = _r.(string)
	}
	var _r3 error
	if _rFn, ok := _rets.Get(3).(func(int, bool) error); ok {
		_r3 = _rFn(a, b)
	} else if _r := _rets.Get(3); _r != nil {
		_r3 = _r.(error)
	}
	return _r0, _r1, _r2, _r3
}

func (_mck *CasesMockOn) OnMethod53(a any, b any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b}
	return _mck.On("Method53", _args...)
}

func (_mck *CasesMockOn) Method54(a int, b Concrete, c pkga.A1, d pkge.E1) {
	_mck.t.Helper()
	_args := []any{a, b, c, d}
	_mck.Called(_args...)
}

func (_mck *CasesMockOn) OnMethod54(a any, b any, c any, d any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b, c, d}
	return _mck.On("Method54", _args...)
}

func (_mck *CasesMockOn) Method55() *Other {
	_mck.t.Helper()
	var _args []any
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 *Other
	if _rFn, ok := _rets.Get(0).(func() *Other); ok {
		_r0 = _rFn()
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(*Other)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod55() *mock.Call {
	_mck.t.Helper()
	var _args []any
	return _mck.On("Method55", _args...)
}

func (_mck *CasesMockOn) Method56(a string, b float64, c ...int) error {
	_mck.t.Helper()
	_args := []any{a, b}
	for _, _arg := range c {
		_args = append(_args, _arg)
	}
	_rets := _mck.Called(_args...)
	if len(_rets) != 1 {
		_mck.t.Fatal("number of mocked method returns does not match")
	}

	var _r0 error
	if _rFn, ok := _rets.Get(0).(func(string, float64, ...int) error); ok {
		_r0 = _rFn(a, b, c...)
	} else if _r := _rets.Get(0); _r != nil {
		_r0 = _r.(error)
	}
	return _r0
}

func (_mck *CasesMockOn) OnMethod56(a any, b any, c ...any) *mock.Call {
	_mck.t.Helper()
	_args := []any{a, b}
	for _, _arg := range c {
		_args = append(_args, _arg)
	}
	return _mck.On("Method56", _args...)
}
