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
