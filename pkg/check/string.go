// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Contain checks that the string "have" contains the substring "want".
//
// See [assert.Contain] for the assertion wrapper.
func Contain(want, have string, opts ...any) error {
	if strings.Contains(have, want) {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected string to contain substring").
		Append("string", "%q", have).
		Append("substring", "%q", want)
	return AddRows(ops, msg)
}

// NotContain checks that the string "have" does **not** contain the
// substring "want".
//
// See [assert.NotContain] for the assertion wrapper.
func NotContain(want, have string, opts ...any) error {
	if strings.Contains(have, want) {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected string not to contain substring").
			Append("string", "%q", have).
			Append("substring", "%q", want)
		return AddRows(ops, msg)
	}
	return nil
}

// EqualFold checks that "want" and "have" are equal, ignoring case.
//
// See [assert.EqualFold] for the assertion wrapper.
func EqualFold(want, have string, opts ...any) error {
	if strings.EqualFold(want, have) {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected strings to be equal ignoring case").
		Append("want", "%q", want).
		Append("have", "%q", have)
	return AddRows(ops, msg)
}
