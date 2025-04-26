// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package affirm is an internal package that provides simple affirmation
// functions designed to improve readability and minimize boilerplate code in
// test cases by offering concise, semantically meaningful functions.
package affirm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ctx42/testing/internal/core"
)

// expected defines the log message for failed affirmations.
const expected = "" +
	"expected values to be equal:\n" +
	"  type: %[1]T\n" +
	"  want: %#[1]v\n" +
	"  have: %#v"

// Equal affirms two comparable types are equal. Returns true if it is,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func Equal[T comparable](t core.T, want, have T) bool {
	t.Helper()
	if want != have {
		t.Errorf(expected, want, have)
		return false
	}
	return true
}

// DeepEqual affirms "want" and "have" are equal using [reflect.DeepEqual].
// Returns true if it is, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
func DeepEqual(t core.T, want, have any) bool {
	t.Helper()
	if !reflect.DeepEqual(want, have) {
		t.Errorf(expected, want, have)
		return false
	}
	return true
}

// Nil affirms "have" is nil. Returns true if it is, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
func Nil(t core.T, have any) bool {
	t.Helper()
	if core.IsNil(have) {
		return true
	}
	t.Errorf(expected, nil, have)
	return false
}

// NotNil affirms "have" is not nil. Returns true if it is not, otherwise
// marks the test as failed, writes an error message to the test log and
// returns false.
func NotNil(t core.T, have any) bool {
	t.Helper()
	if !core.IsNil(have) {
		return true
	}
	typ := fmt.Sprintf("  type: %T\n", have)
	if strings.Contains(typ, "<nil>") {
		typ = ""
	}
	const format = "expected values to be equal:\n" +
		"%s" +
		"  want: nil\n" +
		"  have: <not-nil>"
	t.Fatalf(format, typ)
	return false
}

// Panic affirms "fn" panics. When "fn" panicked, it returns a pointer to a
// string representation of the value used in panic(). When "fn" doesn't panic,
// it returns nil, marks the test as failed and writes an error message to the
// test.
func Panic(t core.T, fn func()) *string {
	t.Helper()
	var val any
	var stack string
	if val, stack = core.WillPanic(fn); stack == "" {
		t.Error("expected fn to panic")
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
