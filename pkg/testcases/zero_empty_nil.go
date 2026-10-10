// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package testcases

import (
	"time"
)

// ZENValue represents one test case for zero / empty / nil behavior.
//
// The expectations follow the semantics of this module's checks, which differ
// from [reflect.Value.IsZero]: a nil pointer, slice, map, channel, function
// or interface is nil but not zero.
type ZENValue struct {
	Desc    string // Human-readable description.
	Val     any    // The value under test.
	IsZero  bool   // Whether check.Zero accepts Val.
	IsEmpty bool   // Whether check.Empty accepts Val.
	IsNil   bool   // Whether check.Nil accepts Val (untyped or typed nil).
}

// ZENValues returns cases covering zero, empty, and nil behavior for many
// Go kinds. Use it when testing custom zero/empty/nil logic or custom
// assertions that rely on these distinctions.
func ZENValues() []ZENValue {
	var nilPtr *TPtr
	var nilItf TItf
	var nilChan chan int
	var nilMap map[int]string
	var nilSlice []int
	nonNilChan := make(chan int)
	nonEmptyChan := make(chan int, 1)
	nonEmptyChan <- 1
	tim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

	return []ZENValue{
		{
			Desc:    "nil",
			Val:     nil,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "empty string",
			Val:     "",
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "not empty string",
			Val:     "abc",
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "zero array",
			Val:     [2]int{},
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "not zero array",
			Val:     [2]int{1, 0},
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "zero struct",
			Val:     TA{},
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "not zero struct",
			Val:     TA{Int: 1},
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "nil type pointer",
			Val:     nilPtr,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "non nil pointer but empty struct",
			Val:     &TPtr{},
			IsZero:  false,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "non nil pointer not empty struct",
			Val:     &TPtr{Val: "abc"},
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "nil interface 1",
			Val:     TItf(nil),
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "nil interface 2",
			Val:     nilItf,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "zero int",
			Val:     0,
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "non zero int",
			Val:     1,
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "zero float64",
			Val:     0.0,
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "non zero float64",
			Val:     1.0,
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "false boolean",
			Val:     false,
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "true boolean",
			Val:     true,
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "nil chan",
			Val:     nilChan,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "non nil chan",
			Val:     nonNilChan,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "not empty chan",
			Val:     nonEmptyChan,
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "time",
			Val:     tim,
			IsZero:  false,
			IsEmpty: false,
			IsNil:   false,
		},
		{
			Desc:    "zero time",
			Val:     time.Time{},
			IsZero:  true,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "nil map",
			Val:     nilMap,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "empty map",
			Val:     map[int]int{},
			IsZero:  false,
			IsEmpty: true,
			IsNil:   false,
		},
		{
			Desc:    "nil slice",
			Val:     nilSlice,
			IsZero:  false,
			IsEmpty: true,
			IsNil:   true,
		},
		{
			Desc:    "empty slice",
			Val:     []int{},
			IsZero:  false,
			IsEmpty: true,
			IsNil:   false,
		},
	}
}
