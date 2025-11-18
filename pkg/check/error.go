// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
)

// Error checks "err" is not nil. Returns an error if it's nil.
func Error(err error, opts ...any) error {
	if err != nil {
		return nil // nolint: nilerr
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected non-nil error")
	return AddRows(ops, msg)
}

// NoError checks "err" is nil. Returns error it's not nil.
func NoError(err error, opts ...any) error {
	if err == nil {
		return nil
	}
	ops := DefaultOptions(opts...)
	const mHeader = "expected the error to be nil"
	if is := core.IsNil(err); is {
		msg := notice.New(mHeader).Want(dump.ValNil).Have("%T", err)
		return AddRows(ops, msg)
	}
	msg := notice.New(mHeader).Want(dump.ValNil).Have("%q", err.Error())
	return AddRows(ops, msg)
}

// ErrorIs checks whether any error in the "err" tree matches the "want" target.
// Returns nil if it's, otherwise returns an error with a message indicating
// the expected and actual values.
func ErrorIs(want, err error, opts ...any) error {
	if errors.Is(err, want) {
		return nil
	}
	ops := DefaultOptions(opts...)
	const hHeader = "expected error to have a target in its tree"
	msg := notice.New(hHeader).
		Want("(%T) %v", want, want).
		Have("(%T) %v", err, err)
	return AddRows(ops, msg)
}

// ErrorAs checks there is an error in the "err" tree that matches the "want"
// target, and if one is found, sets the target to that error. Returns nil if
// the target is found, otherwise returns an error with a message indicating
// the expected and actual values.
func ErrorAs(want any, err error, opts ...any) error {
	if e := Error(err); e != nil {
		return e
	}
	//goland:noinspection GoErrorsAs
	if errors.As(err, want) {
		return nil
	}
	ops := DefaultOptions(opts...)

	tgt := fmt.Sprintf("%T", want)
	if strings.HasPrefix(tgt, "**") {
		tgt = tgt[1:]
	}
	msg := notice.New("expected error to have a target in its tree").
		Append("target", "%s", tgt).
		Append("error", "%T", err)
	return AddRows(ops, msg)
}

// ErrorEqual checks "err" is not nil and its message equals to "want". Returns
// nil if it's, otherwise it returns an error with a message indicating the
// expected and actual values.
func ErrorEqual(want string, err error, opts ...any) error {
	if err != nil && want == err.Error() {
		return nil
	}
	var have any
	have = nil
	if err != nil {
		have = err.Error()
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected the error message to be").
		Want("%q", want).
		Have("%#v", have)
	return AddRows(ops, msg)
}

// ErrorContain checks "err" is not nil and its message contains "want".
// Returns nil if it's, otherwise it returns an error with a message indicating
// the expected and actual values.
func ErrorContain(want string, err error, opts ...any) error {
	if is := core.IsNil(err); is {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected error not to be nil").
			Want("<non-nil>").
			Have("%T", err)
		return AddRows(ops, msg)
	}
	if strings.Contains(err.Error(), want) {
		return nil
	}

	ops := DefaultOptions(opts...)
	var have any
	have = err.Error()
	msg := notice.New("expected the error message to contain").
		Want("%q", want).
		Have("%#v", have)
	return AddRows(ops, msg)
}

// ErrorRegexp checks "err" is not nil and its message matches the "want" regex.
// Returns nil if it is, otherwise it returns an error with a message
// indicating the expected and actual values.
//
// The "want" can be either a regular expression string or instance of
// [regexp.Regexp]. The [fmt.Sprint] is used to get string representation of
// have argument.
func ErrorRegexp(want any, err error, opts ...any) error {
	if is := core.IsNil(err); is {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected error not to be nil").
			Want("<non-nil>").
			Have("%T", err)
		return AddRows(ops, msg)
	}
	if e := Regexp(want, err.Error()); e != nil {
		ops := DefaultOptions(opts...)
		msg := notice.From(e).
			SetHeader("expected the error message to match regexp")
		return AddRows(ops, msg)
	}
	return nil
}
