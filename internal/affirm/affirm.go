// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package affirm is an internal package that provides simple affirmation
// functions designed to improve readability and minimize boilerplate code in
// test cases by offering concise, semantically meaningful functions.
package affirm

import (
	"fmt"
	"reflect"

	"github.com/ctx42/testing/internal/core"
)

// expected defines the log message for failed affirmations.
const expected = "" +
	"expected values to be equal:\n" +
	"  want: (%[1]T) %#[1]v\n" +
	"  have: (%[2]T) %#[2]v"

// Equal affirms two comparable values are equal using ==. Returns true if they
// are, otherwise marks the test as failed, writes an error message to the test
// log, and returns false. For interface type parameters, values whose dynamic
// type is not comparable fail instead of panicking. As with ==, 0.0 equals
// -0.0 and NaN never equals itself.
func Equal[T comparable](t core.T, want, have T) bool {
	t.Helper()
	if !isComparable(want) || !isComparable(have) {
		const format = "" +
			"expected comparable values:\n" +
			"  want: (%[1]T) %#[1]v\n" +
			"  have: (%[2]T) %#[2]v"
		t.Errorf(format, want, have)
		return false
	}
	if want != have {
		t.Errorf(expected, want, have)
		return false
	}
	return true
}

// isComparable reports whether the dynamic value of v can be compared with ==.
func isComparable(v any) bool {
	val := reflect.ValueOf(v)
	return !val.IsValid() || val.Comparable()
}

// DeepEqual affirms "want" and "have" are equal using [reflect.DeepEqual].
// Returns true if it is, otherwise marks the test as failed, writes an error
// message to the test log, and returns false. Pointers are compared by the
// values they point to, not by identity; use [Equal] or [core.Same] for that.
func DeepEqual(t core.T, want, have any) bool {
	t.Helper()
	if !reflect.DeepEqual(want, have) {
		t.Errorf(expected, want, have)
		return false
	}
	return true
}

// Nil affirms "have" is nil, including a nil pointer, slice, map, channel or
// function stored in the interface (see [core.IsNil]). Returns true if it is,
// otherwise marks the test as failed, writes an error message to the test log,
// and returns false.
func Nil(t core.T, have any) bool {
	t.Helper()
	if is := core.IsNil(have); is {
		return true
	}
	t.Errorf("expected nil value:\n  have: (%[1]T) %#[1]v", have)
	return false
}

// NotNil affirms "have" is not nil (see [core.IsNil]). Returns true if it is
// not, otherwise fails the test with Fatalf.
func NotNil(t core.T, have any) bool {
	t.Helper()
	if is := core.IsNil(have); !is {
		return true
	}
	var typ string
	if have != nil {
		typ = fmt.Sprintf("  type: %T\n", have)
	}
	const format = "expected non-nil value:\n" +
		"%s" +
		"  have: nil"
	t.Fatalf(format, typ)
	return false
}

// Panic affirms "fn" panics and returns a pointer to a string representation
// of the value used in panic(). When "fn" doesn't panic, or is nil, it fails
// the test with Fatal, so callers can dereference the result.
func Panic(t core.T, fn func()) *string {
	t.Helper()
	if fn == nil {
		t.Fatal("expected non-nil fn")
		return nil
	}
	var val any
	var stack string
	if val, stack = core.WillPanic(fn); stack == "" {
		t.Fatal("expected fn to panic")
		return nil
	}

	var str string
	switch v := val.(type) {
	case string:
		str = v
	case error:
		str = v.Error()
	default:
		str = fmt.Sprint(v)
	}
	return &str
}
