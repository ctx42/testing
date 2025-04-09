// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"strings"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// TestFunc is signature used by check functions dealing with panics.
type TestFunc func()

// Panic checks "fn" panics. Returns nil if it does, otherwise it returns an
// error with a message with value passed to panic and stack trace.
func Panic(fn TestFunc, opts ...Option) error {
	if _, stack := core.WillPanic(fn); stack == "" {
		ops := DefaultOptions(opts...)
		return notice.New("func should panic").Trail(ops.Trail)
	}
	return nil
}

// NoPanic checks "fn" does not panic. Returns nil if it doesn't, otherwise it
// returns an error with a message with value passed to panic and stack trace.
func NoPanic(fn TestFunc, opts ...Option) error {
	if val, stack := core.WillPanic(fn); stack != "" {
		ops := DefaultOptions(opts...)
		return notice.New("func should not panic").
			Trail(ops.Trail).
			Append("panic value", "%v", val).
			Append("panic stack", "\n%s", notice.Indent(2, ' ', stack))
	}
	return nil
}

// PanicContain checks "fn" panics, and the recovered panic value represented
// as a string contains "want". Returns nil if it does, otherwise it returns an
// error with a message with value passed to panic and stack trace.
func PanicContain(want string, fn TestFunc, opts ...Option) error {
	val, stack := core.WillPanic(fn)
	if stack == "" {
		return notice.New("func should panic")
	}

	var msg string
	switch v := val.(type) {
	case string:
		msg = v
	case error:
		msg = v.Error()
	default:
		msg = fmt.Sprint(v)
	}
	if !strings.Contains(msg, want) {
		ops := DefaultOptions(opts...)
		return notice.New("func should panic with string containing").
			Trail(ops.Trail).
			Append("substring", "%q", want).
			Append("panic value", "%v", val).
			Append("panic stack", "\n%s", notice.Indent(2, ' ', stack))
	}
	return nil
}

// PanicMsg checks "fn" panics, and returns the recovered panic value as a
// string. If function didn't panic, it returns nil and an error with a
// detailed message indicating the expected behaviour.
func PanicMsg(fn TestFunc, opts ...Option) (*string, error) {
	val, stack := core.WillPanic(fn)
	if stack == "" {
		ops := DefaultOptions(opts...)
		return nil, notice.New("func should panic").Trail(ops.Trail)
	}
	var msg string
	switch v := val.(type) {
	case string:
		msg = v
	case error:
		msg = v.Error()
	default:
		msg = fmt.Sprint(v)
	}
	return &msg, nil
}
