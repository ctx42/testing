// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package cases

import (
	"time"

	"github.com/ctx42/testing/internal/types"
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
	var itfVal0, itfVal1, itfPtr0, itfPtr1, itfNil types.TItf
	itfVal0, itfVal1 = types.TVal{}, types.TVal{}
	itfPtr0, itfPtr1 = &types.TPtr{}, &types.TPtr{}
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
			types.TStrType("abc"),
			types.TStrType("abc"),
			true,
		},
		{
			"not equal type value",
			types.TStrType("ab"),
			types.TStrType("abc"),
			false,
		},
		{
			"equal type pointer",
			ptr(types.TStrType("abc")),
			ptr(types.TStrType("abc")),
			true,
		},
		{
			"not equal type value pointer",
			ptr(types.TStrType("ab")),
			ptr(types.TStrType("abc")),
			false,
		},
		{"func", types.TFuncA, types.TFuncA, true},
		{"not equal func", types.TFuncA, types.TFuncB, false},
		{"func ptr", ptr(types.TFuncA), ptr(types.TFuncA), true},
		{"not equal func ptr", ptr(types.TFuncA), ptr(types.TFuncB), false},
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
			map[string]int{"A": types.TCIntA, "B": 2},
			map[string]int{"A": types.TCIntB, "B": 2},
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
			time.Date(2000, 1, 2, 4, 4, 5, 0, types.WAW),
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
			types.WAW,
			false,
		},

		{"itf val 00", itfVal0, itfVal0, true},
		{"itf val 01", itfVal0, itfVal1, true},
		{"itf ptr 00", itfPtr0, itfPtr0, true},
		{"itf ptr 01", itfPtr0, itfPtr1, true},
		{"itf ptr nil 00", itfPtr0, itfNil, false},
		{"itf ptr nil 01", itfNil, itfPtr0, false},
		{"val", types.TPtr{}, types.TPtr{}, true},
		{"val with val", types.TPtr{Val: "A"}, types.TPtr{Val: "A"}, true},
		{"ptr", &types.TPtr{}, &types.TPtr{}, true},
		{"ptr with val", &types.TPtr{Val: "A"}, &types.TPtr{Val: "A"}, true},

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
		{"TCBool A==B", types.TCBoolA, types.TCBoolB, true},
		{"TCString A==B", types.TCStringA, types.TCStringB, true},
		{"TCInt A==B", types.TCIntA, types.TCIntB, true},
		{"TCInt8 A==B", types.TCInt8A, types.TCInt8B, true},
		{"TCInt16 A==B", types.TCInt16A, types.TCInt16B, true},
		{"TCInt32 A==B", types.TCInt32A, types.TCInt32B, true},
		{"TCInt64 A==B", types.TCInt64A, types.TCInt64B, true},
		{"TCUint A==B", types.TCUintA, types.TCUintB, true},
		{"TCUint8 A==B", types.TCUint8A, types.TCUint8B, true},
		{"TCUint16 A==B", types.TCUint16A, types.TCUint16B, true},
		{"TCUint32 A==B", types.TCUint32A, types.TCUint32B, true},
		{"TCUint64 A==B", types.TCUint64A, types.TCUint64B, true},
		{"TCUintptr A==B", types.TCUintptrA, types.TCUintptrB, true},
		{"TCFloat32 A==B", types.TCFloat32A, types.TCFloat32B, true},
		{"TCFloat64 A==B", types.TCFloat64A, types.TCFloat64B, true},
		{"TCComplex64 A==B", types.TCComplex64A, types.TCComplex64B, true},
		{"TCComplex128 A==B", types.TCComplex128A, types.TCComplex128B, true},

		{"TCBool A!=C", types.TCBoolA, types.TCBoolC, false},
		{"TCString A!=C", types.TCStringA, types.TCStringC, false},
		{"TCInt A!=C", types.TCIntA, types.TCIntC, false},
		{"TCInt8 A!=C", types.TCInt8A, types.TCInt8C, false},
		{"TCInt16 A!=C", types.TCInt16A, types.TCInt16C, false},
		{"TCInt32 A!=C", types.TCInt32A, types.TCInt32C, false},
		{"TCInt64 A!=C", types.TCInt64A, types.TCInt64C, false},
		{"TCUint A!=C", types.TCUintA, types.TCUintC, false},
		{"TCUint8 A!=C", types.TCUint8A, types.TCUint8C, false},
		{"TCUint16 A!=C", types.TCUint16A, types.TCUint16C, false},
		{"TCUint32 A!=C", types.TCUint32A, types.TCUint32C, false},
		{"TCUint64 A!=C", types.TCUint64A, types.TCUint64C, false},
		{"TCUintptr A!=C", types.TCUintptrA, types.TCUintptrC, false},
		{"TCFloat32 A!=C", types.TCFloat32A, types.TCFloat32C, false},
		{"TCFloat64 A!=C", types.TCFloat64A, types.TCFloat64C, false},
		{"TCComplex64 A!=C", types.TCComplex64A, types.TCComplex64C, false},
		{"TCComplex128 A!=C", types.TCComplex128A, types.TCComplex128C, false},

		{"CBool A==B", types.CBoolA, types.CBoolB, true},
		{"CBool A!=C", types.CBoolA, types.CBoolC, false},
		{"CString A==B", types.CStringA, types.CStringB, true},
		{"CString A!=C", types.CStringA, types.CStringC, false},
		{"CInt A==B", types.CIntA, types.CIntB, true},
		{"CInt A!=C", types.CIntA, types.CIntC, false},
		{"CFloat64 A==B", types.CFloatA, types.CFloatB, true},
		{"CFloat64 A!=C", types.CFloatA, types.CFloatC, false},
		{"CComplex64 A==B", types.CComplex64A, types.CComplex64B, true},
		{"CComplex64 A!=C", types.CComplex64A, types.CComplex64C, false},
		{"CComplex128 A==B", types.CComplex128A, types.CComplex128B, true},
		{"CComplex128 A!=C", types.CComplex128A, types.CComplex128C, false},

		{"TCBool A==bool", types.TCBoolA, true, true},
		{"TCString A==string", types.TCStringA, "abc", true},
		{"TCInt A==int", types.TCIntA, -123, true},
		{"TCInt8 A==int8", types.TCInt8A, int8(-8), true},
		{"TCInt16 A==int16", types.TCInt16A, int16(-16), true},
		{"TCInt32 A==int32", types.TCInt32A, int32(-32), true},
		{"TCInt64 A==int64", types.TCInt64A, int64(-64), true},
		{"TCUint A==uint", types.TCUintA, uint(123), true},
		{"TCUint8 A==uint8", types.TCUint8A, uint8(8), true},
		{"TCUint16 A==uint16", types.TCUint16A, uint16(16), true},
		{"TCUint32 A==uint32", types.TCUint32A, uint32(32), true},
		{"TCUint64 A==uint64", types.TCUint64A, uint64(64), true},
		{"TCUintptr A==uintptr", types.TCUintptrA, uintptr(42), true},
		{"TCFloat32 A==float32", types.TCFloat32A, float32(32.0), true},
		{"TCFloat64 A==float64", types.TCFloat64A, 64.0, true},
		{"TCComplex64 A==complex64", types.TCComplex64A, complex64(6i + 4), true},
		{"TCComplex128 A==complex128", types.TCComplex128A, 12i + 8, true},

		{"TCBool A!=bool", types.TCBoolA, false, false},
		{"TCString A!=string", types.TCStringA, "cba", false},
		{"TCInt A!=int", types.TCIntA, -321, false},
		{"TCInt8 A!=int8", types.TCInt8A, int8(-13), false},
		{"TCInt16 A!=int16", types.TCInt16A, int16(-61), false},
		{"TCInt32 A!=int32", types.TCInt32A, int32(-23), false},
		{"TCInt64 A!=int64", types.TCInt64A, int64(-46), false},
		{"TCUint A!=uint", types.TCUintA, uint(321), false},
		{"TCUint8 A!=uint8", types.TCUint8A, uint8(13), false},
		{"TCUint16 A!=uint16", types.TCUint16A, uint16(61), false},
		{"TCUint32 A!=uint32", types.TCUint32A, uint32(23), false},
		{"TCUint64 A!=uint64", types.TCUint64A, uint64(46), false},
		{"TCUintptr A!=uintptr", types.TCUintptrA, uintptr(24), false},
		{"TCFloat32 A!=float32", types.TCFloat32A, float32(23.0), false},
		{"TCFloat64 A!=float64", types.TCFloat64A, 46.0, false},
		{"TCComplex64 A!=complex64", types.TCComplex64A, complex64(4i + 6), false},
		{"TCComplex128 A!=complex128", types.TCComplex128A, 8i + 12, false},

		{"CBool A==bool", types.CBoolA, true, true},
		{"CString A==string", types.CStringA, "abc", true},
		{"CInt A==int", types.CIntA, 123, true},
		{"CFloat64 A==float", types.CFloatA, 1.23, true},
		{"CComplex128 A==complex128", types.CComplex128A, 12i + 8, true},

		{"CBool A!=bool", types.CBoolA, false, false},
		{"CString A!=string", types.CStringA, "cba", false},
		{"CInt A!=int", types.CIntA, 321, false},
		{"CFloat64 A!=float", types.CFloatA, 3.21, false},
		{"CComplex128 A!=complex128", types.CComplex128A, 8i + 12, false},
	}
}

// ptr returns pointer to any type.
func ptr[M any](v M) *M { return &v }
