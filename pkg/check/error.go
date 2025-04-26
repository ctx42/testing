// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"strings"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// Error checks "err" is not nil. Returns an error if it's nil.
func Error(err error, opts ...Option) error {
	if err != nil {
		return nil // nolint: nilerr
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected non-nil error").Trail(ops.Trail)
}

// NoError checks "err" is nil. Returns error it's not nil.
func NoError(err error, opts ...Option) error {
	if err == nil {
		return nil
	}
	ops := DefaultOptions(opts...)
	const mHeader = "expected the error to be nil"
	if core.IsNil(err) {
		return notice.New(mHeader).
			Trail(ops.Trail).
			Want("<nil>").
			Have("%T", err)
	}
	return notice.New(mHeader).
		Trail(ops.Trail).
		Want("<nil>").
		Have("%q", err.Error())
}

// ErrorIs checks whether any error in "err" tree matches target. Returns nil
// if it's, otherwise returns an error with a message indicating the expected
// and actual values.
func ErrorIs(err, target error, opts ...Option) error {
	if errors.Is(err, target) {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected err to have a target in its tree").
		Trail(ops.Trail).
		Want("(%T) %v", target, target).
		Have("(%T) %v", err, err)
}

// ErrorAs checks there is an error in the "err" tree that matches a target,
// and if one is found, sets the target to that error. Returns nil if the
// target is found, otherwise returns an error with a message indicating the
// expected and actual values.
func ErrorAs(err error, target any, opts ...Option) error {
	if e := Error(err); e != nil {
		return e
	}
	//goland:noinspection GoErrorsAs
	if errors.As(err, target) {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected err to have a target in its tree").
		Trail(ops.Trail).
		Want("(%T) %#v", err, err).
		Have("(%T) %#v", target, target)
}

// ErrorEqual checks "err" is not nil and its message equals to "want". Returns
// nil if it's, otherwise it returns an error with a message indicating the
// expected and actual values.
func ErrorEqual(want string, err error, opts ...Option) error {
	if err != nil && want == err.Error() {
		return nil
	}
	var have any
	have = nil
	if err != nil {
		have = err.Error()
	}

	ops := DefaultOptions(opts...)
	return notice.New("expected error message to be").
		Trail(ops.Trail).
		Want("%q", want).
		Have("%#v", have)
}

// ErrorContain checks "err" is not nil and its message contains "want".
// Returns nil if it's, otherwise it returns an error with a message indicating
// the expected and actual values.
func ErrorContain(want string, err error, opts ...Option) error {
	if core.IsNil(err) {
		ops := DefaultOptions(opts...)
		return notice.New("expected error not to be nil").
			Trail(ops.Trail).
			Want("<non-nil>").
			Have("%T", err)
	}
	if strings.Contains(err.Error(), want) {
		return nil
	}

	ops := DefaultOptions(opts...)
	var have any
	have = err.Error()
	return notice.New("expected error message to contain").
		Trail(ops.Trail).
		Want("%q", want).
		Have("%#v", have)
}

// ErrorRegexp checks "err" is not nil and its message matches the "want" regex.
// Returns nil if it is, otherwise it returns an error with a message
// indicating the expected and actual values.
//
// The "want" can be either regular expression string or instance of
// [regexp.Regexp]. The [fmt.Sprint] is used to get string representation of
// have argument.
func ErrorRegexp(want any, err error, opts ...Option) error {
	if core.IsNil(err) {
		ops := DefaultOptions(opts...)
		return notice.New("expected error not to be nil").
			Trail(ops.Trail).
			Want("<non-nil>").
			Have("%T", err)
	}
	if e := Regexp(want, err.Error()); e != nil {
		ops := DefaultOptions(opts...)
		return notice.From(e).
			Trail(ops.Trail).
			SetHeader("expected error message to match regexp")
	}
	return nil
}
