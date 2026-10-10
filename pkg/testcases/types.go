// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package testcases

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

// TFuncA is a function used in function value test cases.
func TFuncA() {}

// TFuncB is a second function, different from [TFuncA].
func TFuncB() {}

// /////////////////////////////////////////////////////////////////////////////

// TIntType is a named int type.
type TIntType int

// /////////////////////////////////////////////////////////////////////////////

// TStrType is a named string type.
type TStrType string

// /////////////////////////////////////////////////////////////////////////////

// TItf is an interface used in tests.
type TItf interface{ AAA() string }

// /////////////////////////////////////////////////////////////////////////////

// TInt is a struct with a single exported int field.
type TInt struct{ V int }

// NewTInt returns a [TInt] for 42 and an error for any other value.
func NewTInt(v int) (*TInt, error) {
	if v != 42 {
		return nil, errors.New("not cool")
	}
	return &TInt{V: v}, nil
}

// /////////////////////////////////////////////////////////////////////////////

// TTim is a struct with a [time.Time] field.
type TTim struct{ Tim time.Time }

// /////////////////////////////////////////////////////////////////////////////

// TLoc is a struct with a [time.Location] pointer field.
type TLoc struct{ Loc *time.Location }

// /////////////////////////////////////////////////////////////////////////////

// TPrv is a struct with one exported and many unexported fields of
// different kinds. Use [NewTPrv] and the Set methods to set them.
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

// NewTPrv returns a zero [TPrv].
func NewTPrv() TPrv { return TPrv{} }

// SetInt returns a copy of prv with the unexported int field set.
func (prv TPrv) SetInt(v int) TPrv { prv.vInt = v; return prv }

// SetPtr returns a copy of prv with the unexported pointer field set.
func (prv TPrv) SetPtr(ptr *TVal) TPrv { prv.ptr = ptr; return prv }

// SetSInt returns a copy of prv with the unexported slice field set.
func (prv TPrv) SetSInt(s []int) TPrv { prv.sInt = s; return prv }

// SetAInt returns a copy of prv with the unexported array field set.
func (prv TPrv) SetAInt(a [2]int) TPrv { prv.aInt = a; return prv }

// SetMapII returns a copy of prv with the unexported map field set.
func (prv TPrv) SetMapII(m map[int]int) TPrv { prv.vMap = m; return prv }

// SetTim returns a copy of prv with the unexported time field set.
func (prv TPrv) SetTim(tim time.Time) TPrv { prv.tim = tim; return prv }

// SetFn returns a copy of prv with the unexported function field set.
func (prv TPrv) SetFn(fn func() int) TPrv { prv.fn = fn; return prv }

// SetCh returns a copy of prv with the unexported channel field set.
func (prv TPrv) SetCh(ch chan int) TPrv { prv.ch = ch; return prv }

// /////////////////////////////////////////////////////////////////////////////

// TIntStr is a struct with an int and a string field.
type TIntStr struct {
	Int int
	Str string
}

// /////////////////////////////////////////////////////////////////////////////

// TA is a struct with fields of common kinds, a pointer to itself and an
// unexported field.
type TA struct {
	Int int
	Str string
	Tim time.Time
	Dur time.Duration
	Loc *time.Location
	TAp *TA

	private int //nolint:unused
}

// /////////////////////////////////////////////////////////////////////////////

// TB is a struct embedding [TA] and holding a [TA] value.
type TB struct {
	TA     // Embedded by value.
	TAv TA // Non-pointer type.
}

// /////////////////////////////////////////////////////////////////////////////

// TC is a struct embedding the named string type [TD].
type TC struct {
	TD
	Int int
}

// /////////////////////////////////////////////////////////////////////////////

// TD is a named string type embedded in [TC].
type TD string

// /////////////////////////////////////////////////////////////////////////////

// TNested is a struct with slice and map fields of values and pointers.
type TNested struct {
	SInt    []int
	STA     []TA
	STAp    []*TA
	MStrInt map[string]int
	MStrTyp map[string]TA
	MIntTyp map[int]TA
}

// /////////////////////////////////////////////////////////////////////////////

// TVal is a value type implementing [TItf] and error.
type TVal struct{ Val string } //nolint:errname

// AAA returns the value.
func (typ TVal) AAA() string { return typ.Val }

// Error returns the value.
func (typ TVal) Error() string { return typ.Val }

// /////////////////////////////////////////////////////////////////////////////

// TPtr is a type implementing [TItf] and error with pointer receivers.
type TPtr struct{ Val string } //nolint:errname

// AAA returns the value.
func (typ *TPtr) AAA() string { return typ.Val }

// Variadic formats the value, str and i.
func (typ *TPtr) Variadic(str string, i ...int) string {
	return fmt.Sprintf("%s %s %v", typ.Val, str, i)
}

// Error returns the value.
func (typ *TPtr) Error() string { return typ.Val }

// Wrap returns the value between a and b.
func (typ *TPtr) Wrap(a, b string) string { return a + typ.Val + b }

// Identity returns val.
func (typ *TPtr) Identity(val *string) *string { return val }

// /////////////////////////////////////////////////////////////////////////////

// TRec is a struct with a pointer to its own type.
type TRec struct {
	Int int
	Rec *TRec // Recursive.
}

// /////////////////////////////////////////////////////////////////////////////

// Typed constants of every basic type. The A and B variants are equal, the C
// variant differs.
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
)

// Untyped constants. The A and B variants are equal, the C variant differs.
const (
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
