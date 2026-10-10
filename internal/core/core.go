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
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.Slice, reflect.UnsafePointer:

		return val.IsNil()

	default:
		return false
	}
}

// WillPanic returns the value passed to panic and a non-empty stack trace if
// fn panicked when executed. When fn does not panic, it returns nil and an
// empty stack.
func WillPanic(fn func()) (val any, stack string) {
	defer func() {
		if val = recover(); val != nil {
			stack = string(debug.Stack())
		}
	}()

	fn()
	return
}

// Same returns true when two generic pointers point to the same memory.
//
// It works with pointers to objects, slices, maps, channels and functions.
// Functions are the same when they have the same code. For arrays, and for
// untyped nil arguments, it always returns false; typed nil values of the
// same type are the same.
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

	if wKnd != reflect.Pointer || hKnd != reflect.Pointer {
		return false
	}

	wTyp, hTyp := reflect.TypeOf(want), reflect.TypeOf(have)
	if wTyp != hTyp {
		return false
	}

	// Compare pointer addresses.
	return want == have
}

// sameFunc returns true when arguments are functions of the same type with
// the same code (or both without a name, e.g. nil).
func sameFunc(want, have reflect.Value) bool {
	if !want.IsValid() || !have.IsValid() {
		return false
	}

	if !want.Type().AssignableTo(have.Type()) {
		return false
	}

	fn0pc := runtime.FuncForPC(want.Pointer())
	fn1pc := runtime.FuncForPC(have.Pointer())

	return fn0pc.Name() == fn1pc.Name()
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
