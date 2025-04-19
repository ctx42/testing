package cases

// THIS FILE HAS BEEN GENERATED - DO NOT EDIT!

import (
	"fmt"
	"io/fs"
	mt "time"

	"github.com/ctx42/testing/pkg/mock"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkga"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkgb"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkgc"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkgd"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkge"
	"github.com/ctx42/testing/pkg/tester"
)

type CasesMock struct {
	*mock.Mock
	t tester.T
}

func NewCasesMock(t tester.T) *CasesMock {
	t.Helper()
	return &CasesMock{Mock: mock.NewMock(t), t: t}
}

func (_mck *CasesMock) Method0() {
	_mck.t.Helper()
	var _args []any
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method1(a int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method2(a int, b int) {
	_mck.t.Helper()
	_args := []any{a, b}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method3(a int, b int) {
	_mck.t.Helper()
	_args := []any{a, b}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method4(a int, b int, c bool) {
	_mck.t.Helper()
	_args := []any{a, b, c}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method5(_a0 int) {
	_mck.t.Helper()
	_args := []any{_a0}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method6(a int, _a1 int, b bool) {
	_mck.t.Helper()
	_args := []any{a, _a1, b}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method7() error {
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

func (_mck *CasesMock) Method8() error {
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

func (_mck *CasesMock) Method9() (error, error) {
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

func (_mck *CasesMock) Method10() (int, error) {
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

func (_mck *CasesMock) Method11(_a0 int, _a1 float64) {
	_mck.t.Helper()
	_args := []any{_a0, _a1}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method12(a ...int) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method13(tim mt.Time) error {
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

func (_mck *CasesMock) Method14(_a0 func()) {
	_mck.t.Helper()
	_args := []any{_a0}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method15(_a0 func(int)) {
	_mck.t.Helper()
	_args := []any{_a0}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method16(a func(...int)) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method17() Concrete {
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

func (_mck *CasesMock) Method18() *Concrete {
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

func (_mck *CasesMock) Method19() pkga.A1 {
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

func (_mck *CasesMock) Method20() *pkga.A1 {
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

func (_mck *CasesMock) Method21(a fmt.Stringer) fs.File {
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

func (_mck *CasesMock) Method22(a Concrete) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method23(a *Concrete) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method24(a ...Concrete) int {
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

func (_mck *CasesMock) Method25(a ...*Concrete) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method26(a ...pkga.A1) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method27(a ...*pkga.A1) {
	_mck.t.Helper()
	var _args []any
	for _, _arg := range a {
		_args = append(_args, _arg)
	}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method28(a *int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method29(a pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method30(a *pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method31(a [2]int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method32(a [2]*int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method33(a [2]pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method34(a [2]*pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method35(a []int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method36(a []*int) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method37(a []pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method38(a []*pkga.A1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method39(a map[int]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method40(a map[int]*string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method41(a map[*int]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method42(a map[pkga.A1]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method43(a map[*pkga.A1]string) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method44(a map[*pkga.A1]*pkgb.B1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method45(a chan map[*pkga.A1]*pkgb.B1) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method46(a map[*pkga.A1]*pkgb.B1, b pkgc.C1) *pkgd.D1 {
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

func (_mck *CasesMock) Method47(a func(func(mt.Time, *pkga.A1))) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method48(a map[*pkga.A1]func(pkgb.B1) func(pkga.A1) error) {
	_mck.t.Helper()
	_args := []any{a}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method49(a <-chan *pkga.A1) chan<- int {
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

func (_mck *CasesMock) Method50(a int) (int, int, error) {
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

func (_mck *CasesMock) Method51(e pkge.E1) error {
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

func (_mck *CasesMock) Method52(a Concrete, b mt.Time, c int) (*Concrete, error) {
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

func (_mck *CasesMock) Method53(a int, b bool) (int, bool, string, error) {
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

func (_mck *CasesMock) Method54(a int, b Concrete, c pkga.A1, d pkge.E1) {
	_mck.t.Helper()
	_args := []any{a, b, c, d}
	_mck.Called(_args...)
}

func (_mck *CasesMock) Method55() *Other {
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

func (_mck *CasesMock) Method56(a string, b float64, c ...int) error {
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
