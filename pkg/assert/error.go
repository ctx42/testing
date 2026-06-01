// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Error asserts that "err" is not nil.
//
// See [check.Error] for the error-returning form.
func Error(t tester.T, err error, opts ...any) bool {
	t.Helper()
	if e := check.Error(err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoError asserts that "err" is nil.
//
// See [check.NoError] for the error-returning form.
func NoError(t tester.T, err error, opts ...any) bool {
	t.Helper()
	if e := check.NoError(err, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// ErrorIs asserts whether any error in the "err" tree matches the "want"
// target (via [errors.Is]).
//
// See [check.ErrorIs] for the error-returning form.
func ErrorIs(t tester.T, want, err error, opts ...any) bool {
	t.Helper()
	if e := check.ErrorIs(want, err, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// ErrorIsNot asserts that no error in the "err" tree matches the "want" target.
//
// See [check.ErrorIsNot] for the error-returning form.
func ErrorIsNot(t tester.T, want, err error, opts ...any) bool {
	t.Helper()
	if e := check.ErrorIsNot(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorAs finds the first error in the "err" tree that matches the "want"
// target, and if one is found, sets the target to that error.
//
// See [check.ErrorAs] for the error-returning form.
func ErrorAs(t tester.T, want any, err error, opts ...any) bool {
	t.Helper()
	if e := check.ErrorAs(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorEqual asserts that "err" is not nil and its message equals "want".
//
// See [check.ErrorEqual] for the error-returning form.
func ErrorEqual(t tester.T, want string, err error, opts ...any) bool {
	t.Helper()
	if e := check.ErrorEqual(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorContain asserts that "err" is not nil and its message contains "want".
//
// See [check.ErrorContain] for the error-returning form.
func ErrorContain(t tester.T, want string, err error, opts ...any) bool {
	t.Helper()
	if e := check.ErrorContain(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorRegexp asserts that "err" is not nil and its message matches the "want"
// regular expression.
//
// "want" may be a regular expression string or a [*regexp.Regexp].
// [fmt.Sprint] is used for the string representation of the error.
//
// See [check.ErrorRegexp] for the error-returning form.
func ErrorRegexp(t tester.T, want string, err error, opts ...any) bool {
	t.Helper()
	if e := check.ErrorRegexp(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
