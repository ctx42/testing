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
func Contain(want, have string, opts ...Option) error {
	if strings.Contains(have, want) {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected string to contain substring").
		SetTrail(ops.Trail).
		Append("string", "%q", have).
		Append("substring", "%q", want)
}

// NotContain checks "want" is not a substring of "have". Returns nil if it's,
// otherwise returns an error with a message indicating the expected and actual
// values.
func NotContain(want, have string, opts ...Option) error {
	if strings.Contains(have, want) {
		ops := DefaultOptions(opts...)
		return notice.New("expected string not to contain substring").
			SetTrail(ops.Trail).
			Append("string", "%q", have).
			Append("substring", "%q", want)
	}
	return nil
}
