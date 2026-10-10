// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package testcases provides a rich, curated collection of test values and
// cases designed specifically for people writing custom assertions or using
// tester.Spy to test their own test helpers.
//
// See the Design section in the root README for the overall layered
// architecture and the recommended patterns for testing custom helpers with
// tester.Spy.
//
// The package contains:
//   - Primitive values of every Go kind (including many edge cases)
//   - Structs exercising exported/unexported field combinations
//   - Types implementing common interfaces
//   - Ready-made collections via [EqualCases], [ZENValues], etc.
//
// These values are heavily used inside this module itself to verify
// check.Equal, check.Zero, check.Empty, check.Nil, custom type checkers, and
// Spy-based helper tests. They are the recommended data source so you do not
// have to hand-craft dozens of interesting cases when validating your own
// checkers across Go's type system.
//
// Example usage when testing a custom zero-value check:
//
//	for _, tc := range testcases.ZENValues() {
//	    have := myIsZero(tc.Val)
//	    if have != tc.IsZero { ... }
//	}
//
// See the examples package for complete, tested patterns of using
// testcases together with tester.Spy.
package testcases

import (
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
