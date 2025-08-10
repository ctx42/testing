// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Contain checks "want" is a substring of "have". Returns nil if it's,
// otherwise returns an error with a message indicating the expected and actual
// values.
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

// NotContain checks "want" is not a substring of "have". Returns nil if it's,
// otherwise returns an error with a message indicating the expected and actual
// values.
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
