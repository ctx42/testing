// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package core provides foundational utilities for the Go Testing Module,
// offering low-level functions and types to support readable and robust test
// cases.
package core

import (
	"reflect"
	"runtime"
	"runtime/debug"
	"unsafe"
)

var nilVal = reflect.ValueOf(nil)

// IsNil checks whether the provided interface is actual nil or wrapped nil.
// Actual nil means the interface itself has no type or value (have == nil). A
// wrapped nil means the interface holds a nil value of a concrete type
// (e.g., a nil pointer or slice). It returns two booleans:
//   - isNil: true if the interface is actual nil.
//   - isWrapped: true if the interface holds a nil value of a type.
func IsNil(have any) (isNil, isWrapped bool) {
	if have == nil {
		return true, false
	}
	val := reflect.ValueOf(have)
	kind := val.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice {
		return val.IsNil(), true
	}
	return false, false
}

// WillPanic returns not empty stack trace if the passed function panicked when
// executed and the value that was passed to panic. When a function does not
// panic, it returns nil and empty stack.
func WillPanic(fn func()) (val any, stack string) {
	defer func() {
		if val = recover(); val != nil {
			stack = string(debug.Stack())
		}
	}()

	fn() // Call the target function
	return
}

// Same returns true when two generic pointers point to the same memory.
//
// It works with pointers to objects, slices, maps and functions. For arrays,
// it always returns false.
//
// nolint: cyclop
func Same(want, have any) bool {
	wVal, hVal := reflect.ValueOf(want), reflect.ValueOf(have)
	wKnd, hKnd := wVal.Kind(), hVal.Kind()

	if wKnd == reflect.Func || hKnd == reflect.Func {
		return sameFunc(wVal, hVal)
	}

	if wKnd == reflect.Slice && hKnd == reflect.Slice {
		return same(wVal, hVal)
	}

	if wKnd == reflect.Array && hKnd == reflect.Array {
		return false
	}

	if wKnd == reflect.Map && hKnd == reflect.Map {
		return same(wVal, hVal)
	}

	if wKnd == reflect.Chan && hKnd == reflect.Chan {
		return same(wVal, hVal)
	}

	if wVal.Kind() != reflect.Ptr || hVal.Kind() != reflect.Ptr {
		return false
	}

	wTyp, hTyp := reflect.TypeOf(want), reflect.TypeOf(have)
	if wTyp != hTyp {
		return false
	}

	// Compare pointer addresses.
	return want == have
}

// The sameFunc returns true when arguments represent values for functions or
// methods of the same type.
func sameFunc(want, have reflect.Value) bool {
	if want.Equal(nilVal) || have.Equal(nilVal) {
		return false
	}

	if !want.Type().AssignableTo(have.Type()) {
		return false
	}

	fn0pc := runtime.FuncForPC(want.Pointer())
	fn1pc := runtime.FuncForPC(have.Pointer())

	wName := fn0pc.Name()
	hName := fn1pc.Name()
	if wName == "" && hName == "" {
		return want.Type().AssignableTo(have.Type())
	}
	return wName == hName
}

// same returns true if values represent pointers to the same memory.
func same(want, have reflect.Value) bool {
	if !want.Type().AssignableTo(have.Type()) {
		return false
	}

	wPtr := unsafe.Pointer(want.Pointer())
	hPtr := unsafe.Pointer(have.Pointer())
	return wPtr == hPtr
}
