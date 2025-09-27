// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package testcases

import (
	"time"
)

// EqualCase represents two values and if they are considered equal.
type EqualCase struct {
	Desc     string // The case description.
	Val0     any    // The first value.
	Val1     any    // The second value.
	AreEqual bool   // Are the values equal?
}

// EqualCases returns cases to test equality.
func EqualCases() []EqualCase {
	var itfVal0, itfVal1, itfPtr0, itfPtr1, itfNil TItf
	itfVal0, itfVal1 = TVal{}, TVal{}
	itfPtr0, itfPtr1 = &TPtr{}, &TPtr{}
	mPtr := ptr(map[string]int{"A": 1})
	tim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	ch0, ch1 := make(chan int), make(chan int)

	cases := []EqualCase{
		{"empty string slice", []string{}, []string{}, true},
		{"nil string slice", []string(nil), []string(nil), true},
		{"equal []int", []int{42, 44}, []int{42, 44}, true},
		{"not equal []int", []int{42, 44}, []int{42, 45}, false},
		{"empty []int", []int{}, []int{}, true},
		{"nil []int", []int(nil), []int(nil), true},
		{
			"equal type value",
			TStrType("abc"),
			TStrType("abc"),
			true,
		},
		{
			"not equal type value",
			TStrType("ab"),
			TStrType("abc"),
			false,
		},
		{
			"equal type pointer",
			ptr(TStrType("abc")),
			ptr(TStrType("abc")),
			true,
		},
		{
			"not equal type value pointer",
			ptr(TStrType("ab")),
			ptr(TStrType("abc")),
			false,
		},
		{"func", TFuncA, TFuncA, true},
		{"not equal func", TFuncA, TFuncB, false},
		{"func ptr", ptr(TFuncA), ptr(TFuncA), true},
		{"not equal func ptr", ptr(TFuncA), ptr(TFuncB), false},
		{"equal []any", []any{1, "b", 3.4, tim}, []any{1, "b", 3.4, tim}, true},
		{
			"not equal []any",
			[]any{1, "b", 3.4, tim},
			[]any{1, "b", 3.4, tim.Add(time.Second)},
			false,
		},
		{
			"not equal []any length",
			[]any{1, "b", 3.4, tim},
			[]any{1, "b", 3.4},
			false,
		},

		{
			"equal [][]any",
			[][]any{
				{1, "b", 3.4, tim},
				{2, "c", 5.6, tim.Add(time.Second)},
			},
			[][]any{
				{1, "b", 3.4, tim},
				{2, "c", 5.6, tim.Add(time.Second)},
			},
			true,
		},
		{
			"not equal [][]any",
			[][]any{
				{1, "b", 3.4, tim},
				{2, "c", 5.6, tim.Add(time.Second)},
			},
			[][]any{
				{1, "b", 3.4, tim},
				{1000, "c", 5.6, tim.Add(time.Second)},
			},
			false,
		},
		{
			"equal map[string]int",
			map[string]int{"A": TCIntA, "B": 2},
			map[string]int{"A": TCIntB, "B": 2},
			true,
		},
		{
			"not equal map[string]int",
			map[string]int{"A": 1, "B": 2},
			map[string]int{"A": 1, "B": 3},
			false,
		},
		{
			"not equal map[string]int length",
			map[string]int{"A": 1, "B": 2},
			map[string]int{"A": 1, "B": 2, "C": 3},
			false,
		},
		{
			"not equal map[string]int same length different keys",
			map[string]int{"A": 1, "B": 2},
			map[string]int{"A": 1, "C": 3},
			false,
		},
		{
			"equal time.Time same timezone",
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			true,
		},
		{
			"equal time.Time different timezone",
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2000, 1, 2, 4, 4, 5, 0, WAW),
			true,
		},
		{
			"not equal time.Time",
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2001, 1, 2, 4, 4, 5, 0, time.UTC),
			false,
		},
		{
			"equal time.Location ",
			time.UTC,
			time.UTC,
			true,
		},
		{
			"not equal time.Location",
			time.UTC,
			WAW,
			false,
		},

		{"itf val 00", itfVal0, itfVal0, true},
		{"itf val 01", itfVal0, itfVal1, true},
		{"itf ptr 00", itfPtr0, itfPtr0, true},
		{"itf ptr 01", itfPtr0, itfPtr1, true},
		{"itf ptr nil 00", itfPtr0, itfNil, false},
		{"itf ptr nil 01", itfNil, itfPtr0, false},
		{"val", TPtr{}, TPtr{}, true},
		{"val with val", TPtr{Val: "A"}, TPtr{Val: "A"}, true},
		{"ptr", &TPtr{}, &TPtr{}, true},
		{"ptr with val", &TPtr{Val: "A"}, &TPtr{Val: "A"}, true},

		{"int slice", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"map", map[string]int{"A": 1}, map[string]int{"A": 1}, true},
		{"map ptr", mPtr, mPtr, true},
		{"nil and nil slice", nil, []string(nil), false},
		{"nil slice and nil", []string(nil), nil, false},
		{"nil and nil", nil, nil, true},

		{"chan ptr", &ch0, &ch0, true},
		{"not equal chan ptr", &ch0, &ch1, false},
	}

	cases = append(cases, EqualPrimitives()...)
	cases = append(cases, EqualConstants()...)
	return cases
}

// EqualPrimitives returns cases to test equality for primitive types.
func EqualPrimitives() []EqualCase {
	return []EqualCase{
		{"equal bool true", true, true, true},
		{"equal bool false", false, false, true},
		{"not equal bool", true, false, false},

		{"equal string", "abc", "abc", true},
		{"not equal string", "ab", "abc", false},

		{"equal int", 42, 42, true},
		{"not equal int", 42, 44, false},

		{"equal int8", int8(42), int8(42), true},
		{"not equal int8", int8(42), int8(44), false},
		{"equal int16", int16(42), int16(42), true},
		{"not equal int16", int16(42), int16(44), false},
		{"equal int32", int32(42), int32(42), true},
		{"not equal int32", int32(42), int32(44), false},
		{"equal int64", int64(42), int64(42), true},
		{"not equal int64", int64(42), int64(44), false},

		{"equal uint8", uint8(42), uint8(42), true},
		{"not equal uint8", uint8(42), uint8(44), false},
		{"equal uint16", uint16(42), uint16(42), true},
		{"not equal uint16", uint16(42), uint16(44), false},
		{"equal uint32", uint32(42), uint32(42), true},
		{"not equal uint32", uint32(42), uint32(44), false},
		{"equal uint64", uint64(42), uint64(42), true},
		{"not equal uint64", uint64(42), uint64(44), false},

		{"equal uintptr", uintptr(42), uintptr(42), true},
		{"not equal uintptr", uintptr(42), uintptr(44), false},

		{"equal float64", 42.0, 42.0, true},
		{"not equal float64", 42.0, 44.0, false},
		{"equal float32", float32(42.0), float32(42.0), true},
		{"not equal float32", float32(42.0), float32(44.0), false},
		{
			"equal complex64",
			complex(float32(1.0), float32(2.0)),
			complex(float32(1.0), float32(2.0)),
			true,
		},
		{
			"not equal complex64",
			complex(float32(1.0), float32(2.0)),
			complex(float32(1.0), float32(3.0)),
			false,
		},
		{"equal complex128", complex(1.0, 2.0), complex(1.0, 2.0), true},
		{"not equal complex128", complex(1.0, 2.0), complex(1.0, 3.0), false},
	}
}

// EqualConstants returns cases to test equality for typed constants.
func EqualConstants() []EqualCase {
	return []EqualCase{
		{"TCBool A==B", TCBoolA, TCBoolB, true},
		{"TCString A==B", TCStringA, TCStringB, true},
		{"TCInt A==B", TCIntA, TCIntB, true},
		{"TCInt8 A==B", TCInt8A, TCInt8B, true},
		{"TCInt16 A==B", TCInt16A, TCInt16B, true},
		{"TCInt32 A==B", TCInt32A, TCInt32B, true},
		{"TCInt64 A==B", TCInt64A, TCInt64B, true},
		{"TCUint A==B", TCUintA, TCUintB, true},
		{"TCUint8 A==B", TCUint8A, TCUint8B, true},
		{"TCUint16 A==B", TCUint16A, TCUint16B, true},
		{"TCUint32 A==B", TCUint32A, TCUint32B, true},
		{"TCUint64 A==B", TCUint64A, TCUint64B, true},
		{"TCUintptr A==B", TCUintptrA, TCUintptrB, true},
		{"TCFloat32 A==B", TCFloat32A, TCFloat32B, true},
		{"TCFloat64 A==B", TCFloat64A, TCFloat64B, true},
		{"TCComplex64 A==B", TCComplex64A, TCComplex64B, true},
		{"TCComplex128 A==B", TCComplex128A, TCComplex128B, true},

		{"TCBool A!=C", TCBoolA, TCBoolC, false},
		{"TCString A!=C", TCStringA, TCStringC, false},
		{"TCInt A!=C", TCIntA, TCIntC, false},
		{"TCInt8 A!=C", TCInt8A, TCInt8C, false},
		{"TCInt16 A!=C", TCInt16A, TCInt16C, false},
		{"TCInt32 A!=C", TCInt32A, TCInt32C, false},
		{"TCInt64 A!=C", TCInt64A, TCInt64C, false},
		{"TCUint A!=C", TCUintA, TCUintC, false},
		{"TCUint8 A!=C", TCUint8A, TCUint8C, false},
		{"TCUint16 A!=C", TCUint16A, TCUint16C, false},
		{"TCUint32 A!=C", TCUint32A, TCUint32C, false},
		{"TCUint64 A!=C", TCUint64A, TCUint64C, false},
		{"TCUintptr A!=C", TCUintptrA, TCUintptrC, false},
		{"TCFloat32 A!=C", TCFloat32A, TCFloat32C, false},
		{"TCFloat64 A!=C", TCFloat64A, TCFloat64C, false},
		{"TCComplex64 A!=C", TCComplex64A, TCComplex64C, false},
		{"TCComplex128 A!=C", TCComplex128A, TCComplex128C, false},

		{"CBool A==B", CBoolA, CBoolB, true},
		{"CBool A!=C", CBoolA, CBoolC, false},
		{"CString A==B", CStringA, CStringB, true},
		{"CString A!=C", CStringA, CStringC, false},
		{"CInt A==B", CIntA, CIntB, true},
		{"CInt A!=C", CIntA, CIntC, false},
		{"CFloat64 A==B", CFloatA, CFloatB, true},
		{"CFloat64 A!=C", CFloatA, CFloatC, false},
		{"CComplex64 A==B", CComplex64A, CComplex64B, true},
		{"CComplex64 A!=C", CComplex64A, CComplex64C, false},
		{"CComplex128 A==B", CComplex128A, CComplex128B, true},
		{"CComplex128 A!=C", CComplex128A, CComplex128C, false},

		{"TCBool A==bool", TCBoolA, true, true},
		{"TCString A==string", TCStringA, "abc", true},
		{"TCInt A==int", TCIntA, -123, true},
		{"TCInt8 A==int8", TCInt8A, int8(-8), true},
		{"TCInt16 A==int16", TCInt16A, int16(-16), true},
		{"TCInt32 A==int32", TCInt32A, int32(-32), true},
		{"TCInt64 A==int64", TCInt64A, int64(-64), true},
		{"TCUint A==uint", TCUintA, uint(123), true},
		{"TCUint8 A==uint8", TCUint8A, uint8(8), true},
		{"TCUint16 A==uint16", TCUint16A, uint16(16), true},
		{"TCUint32 A==uint32", TCUint32A, uint32(32), true},
		{"TCUint64 A==uint64", TCUint64A, uint64(64), true},
		{"TCUintptr A==uintptr", TCUintptrA, uintptr(42), true},
		{"TCFloat32 A==float32", TCFloat32A, float32(32.0), true},
		{"TCFloat64 A==float64", TCFloat64A, 64.0, true},
		{"TCComplex64 A==complex64", TCComplex64A, complex64(6i + 4), true},
		{"TCComplex128 A==complex128", TCComplex128A, 12i + 8, true},

		{"TCBool A!=bool", TCBoolA, false, false},
		{"TCString A!=string", TCStringA, "cba", false},
		{"TCInt A!=int", TCIntA, -321, false},
		{"TCInt8 A!=int8", TCInt8A, int8(-13), false},
		{"TCInt16 A!=int16", TCInt16A, int16(-61), false},
		{"TCInt32 A!=int32", TCInt32A, int32(-23), false},
		{"TCInt64 A!=int64", TCInt64A, int64(-46), false},
		{"TCUint A!=uint", TCUintA, uint(321), false},
		{"TCUint8 A!=uint8", TCUint8A, uint8(13), false},
		{"TCUint16 A!=uint16", TCUint16A, uint16(61), false},
		{"TCUint32 A!=uint32", TCUint32A, uint32(23), false},
		{"TCUint64 A!=uint64", TCUint64A, uint64(46), false},
		{"TCUintptr A!=uintptr", TCUintptrA, uintptr(24), false},
		{"TCFloat32 A!=float32", TCFloat32A, float32(23.0), false},
		{"TCFloat64 A!=float64", TCFloat64A, 46.0, false},
		{"TCComplex64 A!=complex64", TCComplex64A, complex64(4i + 6), false},
		{"TCComplex128 A!=complex128", TCComplex128A, 8i + 12, false},

		{"CBool A==bool", CBoolA, true, true},
		{"CString A==string", CStringA, "abc", true},
		{"CInt A==int", CIntA, 123, true},
		{"CFloat64 A==float", CFloatA, 1.23, true},
		{"CComplex128 A==complex128", CComplex128A, 12i + 8, true},

		{"CBool A!=bool", CBoolA, false, false},
		{"CString A!=string", CStringA, "cba", false},
		{"CInt A!=int", CIntA, 321, false},
		{"CFloat64 A!=float", CFloatA, 3.21, false},
		{"CComplex128 A!=complex128", CComplex128A, 8i + 12, false},
	}
}

// ptr returns the pointer to any type.
func ptr[M any](v M) *M { return &v }
