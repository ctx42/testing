// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Error asserts "err" is not nil. Returns true if it's, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
func Error(t tester.T, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.Error(err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoError asserts "err" is nil. Returns true if it is not, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
func NoError(t tester.T, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.NoError(err, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// ErrorIs asserts whether any error in "err" tree matches the "want" target.
// Returns true if it does, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
func ErrorIs(t tester.T, want, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.ErrorIs(want, err, opts...); e != nil {
		t.Fatal(e)
		return false
	}
	return true
}

// ErrorAs finds the first error in "err" tree that matches the "want" target,
// and if one is found, sets a target to that error. Returns true if it does,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func ErrorAs(t tester.T, want any, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.ErrorAs(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorEqual asserts "err" is not nil and its message equals to "want".
// Returns true if it is, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
func ErrorEqual(t tester.T, want string, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.ErrorEqual(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorContain asserts "err" is not nil and its message contains "want".
// Returns true if it does, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
func ErrorContain(t tester.T, want string, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.ErrorContain(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// ErrorRegexp asserts "err" is not nil and its message matches the "want"
// regex. Returns true if it does, otherwise marks the test as failed, writes
// an error message to the test log and returns false.
//
// The "want" can be either a regular expression string or instance of
// [regexp.Regexp]. The [fmt.Sprint] is used to get string representation of
// have argument.
func ErrorRegexp(t tester.T, want string, err error, opts ...check.Option) bool {
	t.Helper()
	if e := check.ErrorRegexp(want, err, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
