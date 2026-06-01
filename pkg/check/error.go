// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

// Error checks that "err" is not nil.
//
// See [assert.Error] for the assertion wrapper.
func Error(err error, opts ...any) error {
	if err != nil {
		return nil // nolint: nilerr
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected non-nil error")
	return AddRows(ops, msg)
}

// NoError checks that "err" is nil.
//
// See [assert.NoError] for the assertion wrapper.
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

// ErrorIs checks that the error tree rooted at "err" contains an error
// that matches "want" according to [errors.Is].
//
// See [assert.ErrorIs] for the assertion wrapper.
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

// ErrorIsNot checks that no error in the tree rooted at "err" matches
// "want" according to [errors.Is].
//
// See [assert.ErrorIsNot] for the assertion wrapper.
func ErrorIsNot(want, err error, opts ...any) error {
	if !errors.Is(err, want) {
		return nil
	}
	ops := DefaultOptions(opts...)
	const hHeader = "expected error to not have a target in its tree"
	msg := notice.New(hHeader).
		Want("(%T) %v", want, want).
		Have("(%T) %v", err, err)
	return AddRows(ops, msg)
}

// ErrorAs checks that an error in the tree rooted at "err" matches the
// target "want" according to [errors.As], and if so, assigns it into the
// pointer provided in "want".
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

// ErrorEqual checks that "err" is not nil and that err.Error() exactly
// equals the string "want".
//
// See [assert.ErrorEqual] for the assertion wrapper.
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

// ErrorContain checks that "err" is not nil and its message contains "want".
//
// See [assert.ErrorContain] for the assertion wrapper.
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
	have := err.Error()
	msg := notice.New("expected the error message to contain").
		Want("%q", want).
		Have("%#v", have)
	return AddRows(ops, msg)
}

// ErrorRegexp checks that "err" is not nil and its message matches the "want"
// regexp (string or [*regexp.Regexp]). Uses [fmt.Sprint] for error text.
//
// See [assert.ErrorRegexp] for the assertion wrapper.
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
