// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Panic asserts that "fn" panics.
//
// See [check.Panic] for the error-returning form.
func Panic(t tester.T, fn check.TestFunc, opts ...any) bool {
	t.Helper()
	if e := check.Panic(fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoPanic asserts that "fn" does not panic.
//
// See [check.NoPanic] for the error-returning form.
func NoPanic(t tester.T, fn check.TestFunc, opts ...any) bool {
	t.Helper()
	if e := check.NoPanic(fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// PanicContain asserts that "fn" panics and the recovered panic value
// (as string) contains "want" (see [check.PanicContain]).
func PanicContain(
	t tester.T,
	want string,
	fn check.TestFunc,
	opts ...any,
) bool {

	t.Helper()
	if e := check.PanicContain(want, fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// PanicMsg asserts that "fn" panics and returns the recovered panic value
// as a string (see [check.PanicMsg]).
func PanicMsg(t tester.T, fn check.TestFunc, opts ...any) *string {
	t.Helper()
	msg, e := check.PanicMsg(fn, opts...)
	if e != nil {
		t.Error(e)
		return nil
	}
	return msg
}
