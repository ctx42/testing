// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

// IsNil checks whether the provided interface is actual nil or wrapped nil.
// Actual nil means the interface itself has no type or value (have == nil).
func IsNil(have any) bool {
	if have == nil {
		return true
	}
	val := reflect.ValueOf(have)
	kind := val.Kind()
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map:
		return val.IsNil()
	case reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
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
//nolint:cyclop
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

	if wVal.Kind() != reflect.Pointer || hVal.Kind() != reflect.Pointer {
		return false
	}

	wTyp, hTyp := reflect.TypeOf(want), reflect.TypeOf(have)
	if wTyp != hTyp {
		return false
	}

	// Compare pointer addresses.
	return want == have
}

// sameFunc returns true when arguments represent values for functions or
// methods of the same type.
func sameFunc(want, have reflect.Value) bool {
	if !want.IsValid() || !have.IsValid() {
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
//
//nolint:gosec
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
//
//nolint:gosec
func Pointer(val reflect.Value) unsafe.Pointer {
	if !val.IsValid() || val.Kind() != reflect.Pointer {
		return nil
	}
	if val.IsNil() {
		return nil
	}
	return val.UnsafePointer()
}

// Value returns the underlying value represented by the [reflect.Value],
// including values obtained through unexported struct fields. The zero
// [reflect.Value] yields untyped nil. Returns nil, false if the value cannot
// be returned.
func Value(val reflect.Value) (any, bool) {
	if !val.IsValid() {
		return nil, true
	}
	if !val.CanInterface() {
		val = exported(val)
		if !val.CanInterface() {
			return nil, false
		}
	}
	return val.Interface(), true
}

// flagRO mirrors the read-only bits of the [reflect.Value] flag word
// (flagStickyRO and flagEmbedRO) set on values obtained through unexported
// struct fields.
const flagRO = 1<<5 | 1<<6

// exported returns a copy of val with its read-only flags cleared, so the
// value can be passed to [reflect.Value.Interface]. It relies on the layout
// of [reflect.Value]: type pointer, data pointer, flag word.
//
//nolint:gosec
func exported(val reflect.Value) reflect.Value {
	type valueHeader struct {
		_    unsafe.Pointer
		_    unsafe.Pointer
		flag uintptr
	}
	hdr := (*valueHeader)(unsafe.Pointer(&val))
	hdr.flag &^= flagRO
	return val
}

// ValueSimple returns the underlying value and true when the provided
// [reflect.Value] is a simple type. Otherwise, returns untyped nil and false.
//
//nolint:cyclop
func ValueSimple(val reflect.Value) (any, bool) {
	switch knd := val.Kind(); knd {
	case reflect.Bool:
		return val.Bool(), true
	case reflect.Int:
		return int(val.Int()), true
	case reflect.Int8:
		return int8(val.Int()), true //nolint:gosec
	case reflect.Int16:
		return int16(val.Int()), true //nolint:gosec
	case reflect.Int32:
		return int32(val.Int()), true //nolint:gosec
	case reflect.Int64:
		return val.Int(), true
	case reflect.Uint:
		return uint(val.Uint()), true
	case reflect.Uint8:
		return uint8(val.Uint()), true //nolint:gosec
	case reflect.Uint16:
		return uint16(val.Uint()), true //nolint:gosec
	case reflect.Uint32:
		return uint32(val.Uint()), true //nolint:gosec
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
