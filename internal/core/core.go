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
// wrapped nil means the interface holds a nil value of a concrete type (e.g.,
// a nil pointer or slice). It returns two booleans:
//
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

// Value returns the underlying value represented by the [reflect.Value].
// Panics for unknown [reflect.Kind].
//
// nolint: cyclop
func Value(val reflect.Value) any {
	knd := val.Kind()
	if knd == reflect.Invalid {
		return nil
	}

	if v, ok := IsSimpleType(val); ok {
		return v
	}

	if knd == reflect.Uintptr {
		return uintptr(val.Uint())
	}

	if knd == reflect.UnsafePointer {
		return val.Pointer()
	}

	switch knd := val.Kind(); knd {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct:

		if !val.CanInterface() {
			valuePtr := unsafe.Pointer(val.UnsafeAddr())
			unsafeValue := reflect.NewAt(val.Type(), valuePtr).Elem()
			return unsafeValue.Interface()
		}
		return val.Interface()

	default:
		panic("unsupported value kind")
	}
}

// IsSimpleType returns the underlying value and true when the provided
// [reflect.Value] is a simple type. Otherwise, returns untyped nil and false.
func IsSimpleType(val reflect.Value) (any, bool) {
	switch knd := val.Kind(); knd {
	case reflect.Bool:
		return val.Bool(), true
	case reflect.Int:
		return int(val.Int()), true
	case reflect.Int8:
		return int8(val.Int()), true // nolint: gosec
	case reflect.Int16:
		return int16(val.Int()), true // nolint: gosec
	case reflect.Int32:
		return int32(val.Int()), true // nolint: gosec
	case reflect.Int64:
		return val.Int(), true
	case reflect.Uint:
		return uint(val.Uint()), true
	case reflect.Uint8:
		return uint8(val.Uint()), true // nolint: gosec
	case reflect.Uint16:
		return uint16(val.Uint()), true // nolint: gosec
	case reflect.Uint32:
		return uint32(val.Uint()), true // nolint: gosec
	case reflect.Uint64:
		return val.Uint(), true
	case reflect.Float32:
		return float32(val.Float()), true
	case reflect.Float64:
		return val.Float(), true
	case reflect.Complex64:
		return complex64(val.Complex()), true
	case reflect.Complex128:
		return val.Complex(), true
	case reflect.String:
		return val.String(), true
	}
	return nil, false
}
