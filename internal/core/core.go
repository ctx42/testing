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
// Actual nil means the interface itself has no type or value (have == nil).
func IsNil(have any) bool {
	if have == nil {
		return true
	}
	val := reflect.ValueOf(have)
	val.IsValid()
	kind := val.Kind()
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map:
		return val.IsNil()
	case reflect.Pointer, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
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

// Pointer checks if the argument represents a pointer type and returns its
// memory address as an [unsafe.Pointer], otherwise returns nil.
func Pointer(val reflect.Value) unsafe.Pointer {
	if !val.IsValid() || val.Kind() != reflect.Ptr {
		return nil
	}
	if val.IsNil() {
		return nil
	}
	if val.CanAddr() {
		return unsafe.Pointer(val.UnsafeAddr())
	}
	if val.CanInterface() || !val.CanAddr() {
		return val.UnsafePointer()
	}
	// TODO(rz): test this case.
	return nil
}

// Value returns the underlying value represented by the [reflect.Value].
// Returns nil, false if the underlying value cannot be returned.
func Value(val reflect.Value) (any, bool) {
	// TODO(rz): improve coverage.
	if nilVal.Equal(val) {
		return nil, true
	}

	if v, ok := ValueSimple(val); ok {
		return v, true
	}

	if val.CanInterface() {
		return val.Interface(), true
	}

	if val.CanAddr() {
		valPtr := unsafe.Pointer(val.UnsafeAddr())
		valUnsafe := reflect.NewAt(val.Type(), valPtr).Elem()
		return valUnsafe.Interface(), true
	}

	switch knd := val.Kind(); knd {
	case reflect.Pointer:
		if val.Elem().Kind() == reflect.Struct {
			return value(val.Type(), val.Elem())
		}
		return nil, false

	case reflect.Func, reflect.Chan:
		return value(val.Type(), val)

	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		return value(val.Type(), val)

	default:
		return nil, false
	}
}

// Value extracts the underlying value from a [reflect.Value] representing an
// unexported or unaddressable field, returning it as an "any" with a boolean
// indicating success.
//
// The typ parameter specifies the expected type of the value, and val is the
// [reflect.Value] to extract. For unexported fields, it uses unsafe to bypass
// reflection restrictions. If val is unaddressable (e.g., from a struct passed
// by value), it creates an addressable copy.
//
// The returned bool is true if the extraction succeeds, false otherwise (e.g.,
// nil pointer or invalid value).
//
// The function is unsafe and assumes val is a valid field of the specified
// type.
func value(typ reflect.Type, val reflect.Value) (any, bool) {
	v := reflect.New(typ).Elem()

	if val.CanAddr() {
		valPtr := unsafe.Pointer(val.UnsafeAddr())
		*(*unsafe.Pointer)(unsafe.Pointer(v.UnsafeAddr())) = valPtr
		return v.Interface(), true
	}

	// Unaddressable: Copy the value's memory byte-by-byte.
	size := int(typ.Size())
	// Access the [reflect.Value] internal representation.
	type valueHeader struct {
		_    uintptr        // Type pointer.
		data unsafe.Pointer // Data pointer.
		_    uintptr        // Flag or padding.
	}
	vHeader := (*valueHeader)(unsafe.Pointer(&val))
	if vHeader.data == nil {
		return nil, false
	}

	// Copy the bytes from vHeader.data to v.
	destPtr := unsafe.Pointer(v.UnsafeAddr())
	srcPtr := vHeader.data
	for i := 0; i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(destPtr) + uintptr(i))) =
			*(*byte)(unsafe.Pointer(uintptr(srcPtr) + uintptr(i)))
	}

	// Return the struct value as any
	return v.Interface(), true
}

// ValueSimple returns the underlying value and true when the provided
// [reflect.Value] is a simple type. Otherwise, returns untyped nil and false.
//
// nolint: cyclop
func ValueSimple(val reflect.Value) (any, bool) {
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
	default:
		return nil, false
	}
}
