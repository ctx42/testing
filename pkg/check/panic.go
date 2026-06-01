// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"strings"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// TestFunc is a signature used by check functions dealing with panics.
// See [assert.Panic], [assert.NoPanic], etc.
type TestFunc func()

// Panic checks that calling "fn" causes a panic.
//
// On success (panic occurred) it returns nil.
// On failure it returns a structured error.
func Panic(fn TestFunc, opts ...any) error {
	if _, stack := core.WillPanic(fn); stack == "" {
		ops := DefaultOptions(opts...)
		msg := notice.New("func should panic")
		return AddRows(ops, msg)
	}
	return nil
}

// NoPanic checks that calling "fn" does **not** cause a panic.
//
// If a panic occurs, the recovered value and stack are included in the
// returned error.
func NoPanic(fn TestFunc, opts ...any) error {
	if val, stack := core.WillPanic(fn); stack != "" {
		ops := DefaultOptions(opts...)
		msg := notice.New("func should not panic").
			Append("panic value", "%v", val).
			Append("panic stack", "%s", notice.Indent(2, ' ', stack))
		return AddRows(ops, msg)
	}
	return nil
}

// PanicContain checks that "fn" panics and that the string representation
// of the recovered panic value contains the substring "want".
//
// See [assert.PanicContain] for the assertion wrapper.
func PanicContain(want string, fn TestFunc, opts ...any) error {
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
		msg := notice.New("func should panic with string containing").
			Append("substring", "%q", want).
			Append("panic value", "%v", val).
			Append("panic stack", "%s", notice.Indent(2, ' ', stack))
		return AddRows(ops, msg)
	}
	return nil
}

// PanicMsg checks that "fn" panics and returns the recovered panic value as a
// string. Returns an error if "fn" did not panic.
// See [assert.PanicMsg].
func PanicMsg(fn TestFunc, opts ...any) (*string, error) {
	val, stack := core.WillPanic(fn)
	if stack == "" {
		ops := DefaultOptions(opts...)
		msg := notice.New("func should panic")
		return nil, AddRows(ops, msg)
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
