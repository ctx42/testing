// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Panic asserts "fn" panics. Returns true if it panicked, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
func Panic(t tester.T, fn check.TestFunc, opts ...any) bool {
	t.Helper()
	if e := check.Panic(fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoPanic asserts "fn" does not panic. Returns true if it did not panic,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func NoPanic(t tester.T, fn check.TestFunc, opts ...any) bool {
	t.Helper()
	if e := check.NoPanic(fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// PanicContain asserts "fn" panics, and the recovered panic value represented
// as a string contains "want". Returns true if it panics and does contain the
// wanted string, otherwise marks the test as failed, writes an error message
// to the test log and returns false.
func PanicContain(t tester.T, want string, fn check.TestFunc, opts ...any) bool {
	t.Helper()
	if e := check.PanicContain(want, fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// PanicMsg asserts the "fn" panics and returns the recovered panic value
// represented as a string. If the function did not panic, it marks the test as
// failed and writes an error message to the test log.
func PanicMsg(t tester.T, fn check.TestFunc, opts ...any) *string {
	t.Helper()
	msg, e := check.PanicMsg(fn, opts...)
	if e != nil {
		t.Error(e)
		return nil
	}
	return msg
}
