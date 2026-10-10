// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package testcases

import (
	"time"
)

// EqualCase represents one test case for equality behavior.
//
// AreEqual follows check.Equal, which differs from [reflect.DeepEqual]: it
// compares time.Time values by instant (ignoring the location), functions by
// identity and pointers by the values they point to.
type EqualCase struct {
	Desc     string // Human-readable description of the case.
	Val0     any    // First value.
	Val1     any    // Second value.
	AreEqual bool   // Whether check.Equal considers the values equal.
}

// EqualCases returns a large, battle-tested set of cases for verifying
// custom equality logic.
//
// The collection covers primitives, structs (exported/unexported fields),
// pointers, interfaces, slices, maps, channels, time values, and more.
// It is the same data used inside this module to test check.Equal and
// registered type checkers.
//
// Recommended usage when testing your own equality helper:
//
//	for _, tc := range testcases.EqualCases() {
//	    have := myEqual(tc.Val0, tc.Val1)
//	    if have != tc.AreEqual { ... }
//	}
func EqualCases() []EqualCase {
	var itfVal0, itfVal1, itfPtr0, itfPtr1, itfNil TItf
	itfVal0, itfVal1 = TVal{}, TVal{}
	itfPtr0, itfPtr1 = &TPtr{}, &TPtr{}
	mPtr := new(map[string]int{"A": 1})
	tim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	ch0, ch1 := make(chan int), make(chan int)

	cases := []EqualCase{
		{"empty string slice", []string{}, []string{}, true},
		{"nil string slice", []string(nil), []string(nil), true},
		{"equal int slice", []int{42, 44}, []int{42, 44}, true},
		{"not equal int slice", []int{42, 44}, []int{42, 45}, false},
		{"empty int slice", []int{}, []int{}, true},
		{"nil int slice", []int(nil), []int(nil), true},
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
			new(TStrType("abc")),
			new(TStrType("abc")),
			true,
		},
		{
			"not equal type value pointer",
			new(TStrType("ab")),
			new(TStrType("abc")),
			false,
		},
		{"func", TFuncA, TFuncA, true},
		{"not equal func", TFuncA, TFuncB, false},
		{"func ptr", new(TFuncA), new(TFuncA), true},
		{"not equal func ptr", new(TFuncA), new(TFuncB), false},
		{
			"equal any slice",
			[]any{1, "b", 3.4, tim},
			[]any{1, "b", 3.4, tim},
			true,
		},
		{
			"not equal any slice",
			[]any{1, "b", 3.4, tim},
			[]any{1, "b", 3.4, tim.Add(time.Second)},
			false,
		},
		{
			"not equal any slice length",
			[]any{1, "b", 3.4, tim},
			[]any{1, "b", 3.4},
			false,
		},

		{
			"equal nested any slice",
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
			"not equal nested any slice",
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
			"equal string to int map",
			map[string]int{"A": TCIntA, "B": 2},
			map[string]int{"A": TCIntB, "B": 2},
			true,
		},
		{
			"not equal string to int map",
			map[string]int{"A": 1, "B": 2},
			map[string]int{"A": 1, "B": 3},
			false,
		},
		{
			"not equal string to int map length",
			map[string]int{"A": 1, "B": 2},
			map[string]int{"A": 1, "B": 2, "C": 3},
			false,
		},
		{
			"not equal string to int map same length different keys",
			map[string]int{"A": 1, "B": 2},
			map[string]int{"A": 1, "C": 3},
			false,
		},
		{
			"equal time Time same timezone",
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			true,
		},
		{
			"equal time Time different timezone",
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2000, 1, 2, 4, 4, 5, 0, WAW),
			true,
		},
		{
			"not equal time Time",
			time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2001, 1, 2, 4, 4, 5, 0, time.UTC),
			false,
		},
		{
			"equal time Location",
			time.UTC,
			time.UTC,
			true,
		},
		{
			"not equal time Location",
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
		{"string and named string", "abc", TStrType("abc"), false},

		{"equal int", 42, 42, true},
		{"not equal int", 42, 44, false},
		{"int and int64", 42, int64(42), false},
		{"int and named int", 1, TIntType(1), false},

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

// EqualConstants returns cases to test equality for typed and untyped
// constants.
func EqualConstants() []EqualCase {
	return []EqualCase{
		{"TCBool A eq B", TCBoolA, TCBoolB, true},
		{"TCString A eq B", TCStringA, TCStringB, true},
		{"TCInt A eq B", TCIntA, TCIntB, true},
		{"TCInt8 A eq B", TCInt8A, TCInt8B, true},
		{"TCInt16 A eq B", TCInt16A, TCInt16B, true},
		{"TCInt32 A eq B", TCInt32A, TCInt32B, true},
		{"TCInt64 A eq B", TCInt64A, TCInt64B, true},
		{"TCUint A eq B", TCUintA, TCUintB, true},
		{"TCUint8 A eq B", TCUint8A, TCUint8B, true},
		{"TCUint16 A eq B", TCUint16A, TCUint16B, true},
		{"TCUint32 A eq B", TCUint32A, TCUint32B, true},
		{"TCUint64 A eq B", TCUint64A, TCUint64B, true},
		{"TCUintptr A eq B", TCUintptrA, TCUintptrB, true},
		{"TCFloat32 A eq B", TCFloat32A, TCFloat32B, true},
		{"TCFloat64 A eq B", TCFloat64A, TCFloat64B, true},
		{"TCComplex64 A eq B", TCComplex64A, TCComplex64B, true},
		{"TCComplex128 A eq B", TCComplex128A, TCComplex128B, true},

		{"TCBool A ne C", TCBoolA, TCBoolC, false},
		{"TCString A ne C", TCStringA, TCStringC, false},
		{"TCInt A ne C", TCIntA, TCIntC, false},
		{"TCInt8 A ne C", TCInt8A, TCInt8C, false},
		{"TCInt16 A ne C", TCInt16A, TCInt16C, false},
		{"TCInt32 A ne C", TCInt32A, TCInt32C, false},
		{"TCInt64 A ne C", TCInt64A, TCInt64C, false},
		{"TCUint A ne C", TCUintA, TCUintC, false},
		{"TCUint8 A ne C", TCUint8A, TCUint8C, false},
		{"TCUint16 A ne C", TCUint16A, TCUint16C, false},
		{"TCUint32 A ne C", TCUint32A, TCUint32C, false},
		{"TCUint64 A ne C", TCUint64A, TCUint64C, false},
		{"TCUintptr A ne C", TCUintptrA, TCUintptrC, false},
		{"TCFloat32 A ne C", TCFloat32A, TCFloat32C, false},
		{"TCFloat64 A ne C", TCFloat64A, TCFloat64C, false},
		{"TCComplex64 A ne C", TCComplex64A, TCComplex64C, false},
		{"TCComplex128 A ne C", TCComplex128A, TCComplex128C, false},

		{"CBool A eq B", CBoolA, CBoolB, true},
		{"CBool A ne C", CBoolA, CBoolC, false},
		{"CString A eq B", CStringA, CStringB, true},
		{"CString A ne C", CStringA, CStringC, false},
		{"CInt A eq B", CIntA, CIntB, true},
		{"CInt A ne C", CIntA, CIntC, false},
		{"CFloat64 A eq B", CFloatA, CFloatB, true},
		{"CFloat64 A ne C", CFloatA, CFloatC, false},
		{"CComplex64 A eq B", CComplex64A, CComplex64B, true},
		{"CComplex64 A ne C", CComplex64A, CComplex64C, false},
		{"CComplex128 A eq B", CComplex128A, CComplex128B, true},
		{"CComplex128 A ne C", CComplex128A, CComplex128C, false},

		{"TCBool A eq bool", TCBoolA, true, true},
		{"TCString A eq string", TCStringA, "abc", true},
		{"TCInt A eq int", TCIntA, -123, true},
		{"TCInt8 A eq int8", TCInt8A, int8(-8), true},
		{"TCInt16 A eq int16", TCInt16A, int16(-16), true},
		{"TCInt32 A eq int32", TCInt32A, int32(-32), true},
		{"TCInt64 A eq int64", TCInt64A, int64(-64), true},
		{"TCUint A eq uint", TCUintA, uint(123), true},
		{"TCUint8 A eq uint8", TCUint8A, uint8(8), true},
		{"TCUint16 A eq uint16", TCUint16A, uint16(16), true},
		{"TCUint32 A eq uint32", TCUint32A, uint32(32), true},
		{"TCUint64 A eq uint64", TCUint64A, uint64(64), true},
		{"TCUintptr A eq uintptr", TCUintptrA, uintptr(42), true},
		{"TCFloat32 A eq float32", TCFloat32A, float32(32.0), true},
		{"TCFloat64 A eq float64", TCFloat64A, 64.0, true},
		{"TCComplex64 A eq complex64", TCComplex64A, complex64(6i + 4), true},
		{"TCComplex128 A eq complex128", TCComplex128A, 12i + 8, true},

		{"TCBool A ne bool", TCBoolA, false, false},
		{"TCString A ne string", TCStringA, "cba", false},
		{"TCInt A ne int", TCIntA, -321, false},
		{"TCInt8 A ne int8", TCInt8A, int8(-13), false},
		{"TCInt16 A ne int16", TCInt16A, int16(-61), false},
		{"TCInt32 A ne int32", TCInt32A, int32(-23), false},
		{"TCInt64 A ne int64", TCInt64A, int64(-46), false},
		{"TCUint A ne uint", TCUintA, uint(321), false},
		{"TCUint8 A ne uint8", TCUint8A, uint8(13), false},
		{"TCUint16 A ne uint16", TCUint16A, uint16(61), false},
		{"TCUint32 A ne uint32", TCUint32A, uint32(23), false},
		{"TCUint64 A ne uint64", TCUint64A, uint64(46), false},
		{"TCUintptr A ne uintptr", TCUintptrA, uintptr(24), false},
		{"TCFloat32 A ne float32", TCFloat32A, float32(23.0), false},
		{"TCFloat64 A ne float64", TCFloat64A, 46.0, false},
		{"TCComplex64 A ne complex64", TCComplex64A, complex64(4i + 6), false},
		{"TCComplex128 A ne complex128", TCComplex128A, 8i + 12, false},

		{"CBool A eq bool", CBoolA, true, true},
		{"CString A eq string", CStringA, "abc", true},
		{"CInt A eq int", CIntA, 123, true},
		{"CFloat64 A eq float", CFloatA, 1.23, true},
		{"CComplex128 A eq complex128", CComplex128A, 12i + 8, true},

		{"CBool A ne bool", CBoolA, false, false},
		{"CString A ne string", CStringA, "cba", false},
		{"CInt A ne int", CIntA, 321, false},
		{"CFloat64 A ne float", CFloatA, 3.21, false},
		{"CComplex128 A ne complex128", CComplex128A, 8i + 12, false},
	}
}
