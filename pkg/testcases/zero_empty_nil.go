// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package testcases

import (
	"time"
)

// ZENValue represents a value and if it's considered zero, empty or nil value.
type ZENValue struct {
	Desc         string // The value description.
	Val          any    // The value.
	IsZero       bool   // Is Val considered zero value.
	IsEmpty      bool   // Is Val considered empty value.
	IsNil        bool   // Is Val considered nil value.
	IsWrappedNil bool   // Is Val a wrapped nil value.
}

// ZENValues returns cases for zero, empty and nil values.
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
			Desc:         "nil",
			Val:          nil,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: false,
		},
		{
			Desc:         "nil type pointer",
			Val:          nilPtr,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: true,
		},
		{
			Desc:         "non nil pointer but empty struct",
			Val:          &TPtr{},
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: true,
		},
		{
			Desc:         "non nil pointer not empty struct",
			Val:          &TPtr{Val: "abc"},
			IsZero:       false,
			IsEmpty:      false,
			IsNil:        false,
			IsWrappedNil: true,
		},
		{
			Desc:         "nil interface 1",
			Val:          TItf(nil),
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: false,
		},
		{
			Desc:         "nil interface 2",
			Val:          nilItf,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: false,
		},
		{
			Desc:         "zero int",
			Val:          0,
			IsZero:       true,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "non zero int",
			Val:          1,
			IsZero:       false,
			IsEmpty:      false,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "zero float64",
			Val:          0.0,
			IsZero:       true,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "non zero float64",
			Val:          1.0,
			IsZero:       false,
			IsEmpty:      false,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "false boolean",
			Val:          false,
			IsZero:       true,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "true boolean",
			Val:          true,
			IsZero:       false,
			IsEmpty:      false,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "nil chan",
			Val:          nilChan,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: true,
		},
		{
			Desc:         "non nil chan",
			Val:          nonNilChan,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: true,
		},
		{
			Desc:         "not empty chan",
			Val:          nonEmptyChan,
			IsZero:       false,
			IsEmpty:      false,
			IsNil:        false,
			IsWrappedNil: true,
		},
		{
			Desc:         "time",
			Val:          tim,
			IsZero:       false,
			IsEmpty:      false,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "zero time",
			Val:          time.Time{},
			IsZero:       true,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: false,
		},
		{
			Desc:         "nil map",
			Val:          nilMap,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: true,
		},
		{
			Desc:         "empty map",
			Val:          map[int]int{},
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: true,
		},
		{
			Desc:         "nil slice",
			Val:          nilSlice,
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        true,
			IsWrappedNil: true,
		},
		{
			Desc:         "empty slice",
			Val:          []int{},
			IsZero:       false,
			IsEmpty:      true,
			IsNil:        false,
			IsWrappedNil: true,
		},
	}
}
