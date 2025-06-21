// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package types provides example types used in tests.
package types

import (
	"errors"
	"fmt"
	"time"
)

// WAW represents Europe/Warsaw timezone.
var WAW *time.Location

func init() {
	var err error
	WAW, err = time.LoadLocation("Europe/Warsaw")
	if err != nil {
		panic(err)
	}
}

func TFuncA() {}

func TFuncB() {}

// /////////////////////////////////////////////////////////////////////////////

type TIntType int

// /////////////////////////////////////////////////////////////////////////////

type TStrType string

// /////////////////////////////////////////////////////////////////////////////

// TItf is an interface used in tests.
type TItf interface{ AAA() string }

// /////////////////////////////////////////////////////////////////////////////

type TInt struct{ V int }

func NewTInt(v int) (*TInt, error) {
	if v != 42 {
		return nil, errors.New("not cool")
	}
	return &TInt{V: v}, nil
}

// /////////////////////////////////////////////////////////////////////////////

type TTim struct{ Tim time.Time }

// /////////////////////////////////////////////////////////////////////////////

type TLoc struct{ Loc *time.Location }

// /////////////////////////////////////////////////////////////////////////////

type TPrv struct {
	Pub  int
	vInt int
	ptr  *TVal
	sInt []int
	aInt [2]int
	vMap map[int]int
	tim  time.Time
	fn   func() int
	ch   chan int
}

func NewTPrv() TPrv                          { return TPrv{} }
func (prv TPrv) SetInt(v int) TPrv           { prv.vInt = v; return prv }
func (prv TPrv) SetPtr(ptr *TVal) TPrv       { prv.ptr = ptr; return prv }
func (prv TPrv) SetSInt(s []int) TPrv        { prv.sInt = s; return prv }
func (prv TPrv) SetAInt(a [2]int) TPrv       { prv.aInt = a; return prv }
func (prv TPrv) SetMapII(m map[int]int) TPrv { prv.vMap = m; return prv }
func (prv TPrv) SetTim(tim time.Time) TPrv   { prv.tim = tim; return prv }
func (prv TPrv) SetFn(fn func() int) TPrv    { prv.fn = fn; return prv }
func (prv TPrv) SetCh(ch chan int) TPrv      { prv.ch = ch; return prv }

// /////////////////////////////////////////////////////////////////////////////

type TIntStr struct {
	Int int
	Str string
}

// /////////////////////////////////////////////////////////////////////////////

type TA struct {
	Int int
	Str string
	Tim time.Time
	Dur time.Duration
	Loc *time.Location
	TAp *TA

	private int // nolint: unused
}

// /////////////////////////////////////////////////////////////////////////////

type TB struct {
	TA     // Embedded by value.
	TAv TA // Non-pointer type.
}

// /////////////////////////////////////////////////////////////////////////////

type TC struct {
	TD
	Int int
}

// /////////////////////////////////////////////////////////////////////////////

type TD string

// /////////////////////////////////////////////////////////////////////////////

type TNested struct {
	SInt    []int
	STA     []TA
	STAp    []*TA
	MStrInt map[string]int
	MStrTyp map[string]TA
	MIntTyp map[int]TA
}

// /////////////////////////////////////////////////////////////////////////////

type TVal struct{ Val string } // nolint: errname

func (typ TVal) AAA() string { return typ.Val }

func (typ TVal) Error() string { return typ.Val }

// /////////////////////////////////////////////////////////////////////////////

type TPtr struct{ Val string } // nolint: errname

func (typ *TPtr) AAA() string { return typ.Val }

func (typ *TPtr) Variadic(str string, i ...int) string {
	return fmt.Sprintf("%s %s %v", typ.Val, str, i)
}

func (typ *TPtr) Error() string { return typ.Val }

func (typ *TPtr) Wrap(a, b string) string { return a + typ.Val + b }

func (typ *TPtr) Identity(val *string) *string { return val }

// /////////////////////////////////////////////////////////////////////////////

type TRec struct {
	Int int
	Rec *TRec // Recursive.
}

// /////////////////////////////////////////////////////////////////////////////

const (
	TCBoolA       bool       = true
	TCBoolB       bool       = true
	TCBoolC       bool       = false
	TCStringA     string     = "abc"
	TCStringB     string     = "abc"
	TCStringC     string     = "cba"
	TCIntA        int        = -123
	TCIntB        int        = -123
	TCIntC        int        = -321
	TCInt8A       int8       = -8
	TCInt8B       int8       = -8
	TCInt8C       int8       = -13
	TCInt16A      int16      = -16
	TCInt16B      int16      = -16
	TCInt16C      int16      = -61
	TCInt32A      int32      = -32
	TCInt32B      int32      = -32
	TCInt32C      int32      = -23
	TCInt64A      int64      = -64
	TCInt64B      int64      = -64
	TCInt64C      int64      = -46
	TCUintA       uint       = 123
	TCUintB       uint       = 123
	TCUintC       uint       = 321
	TCUint8A      uint8      = 8
	TCUint8B      uint8      = 8
	TCUint8C      uint8      = 13
	TCUint16A     uint16     = 16
	TCUint16B     uint16     = 16
	TCUint16C     uint16     = 61
	TCUint32A     uint32     = 32
	TCUint32B     uint32     = 32
	TCUint32C     uint32     = 23
	TCUint64A     uint64     = 64
	TCUint64B     uint64     = 64
	TCUint64C     uint64     = 46
	TCUintptrA    uintptr    = 42
	TCUintptrB    uintptr    = 42
	TCUintptrC    uintptr    = 24
	TCFloat32A    float32    = 32.0
	TCFloat32B    float32    = 32.0
	TCFloat32C    float32    = 23.0
	TCFloat64A    float64    = 64.0
	TCFloat64B    float64    = 64.0
	TCFloat64C    float64    = 46.0
	TCComplex64A  complex64  = 6i + 4
	TCComplex64B  complex64  = 6i + 4
	TCComplex64C  complex64  = 4i + 6
	TCComplex128A complex128 = 12i + 8
	TCComplex128B complex128 = 12i + 8
	TCComplex128C complex128 = 8i + 12

	CBoolA       = true
	CBoolB       = true
	CBoolC       = false
	CStringA     = "abc"
	CStringB     = "abc"
	CStringC     = "xyz"
	CIntA        = 123
	CIntB        = 123
	CIntC        = 321
	CFloatA      = 1.23
	CFloatB      = 1.23
	CFloatC      = 3.21
	CComplex64A  = 6i + 4
	CComplex64B  = 6i + 4
	CComplex64C  = 4i + 6
	CComplex128A = 12i + 8
	CComplex128B = 12i + 8
	CComplex128C = 8i + 12
)
