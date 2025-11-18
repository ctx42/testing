// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// Nil checks "have" is nil. Returns nil if it's, otherwise returns an error
// with a message indicating the expected and actual values.
func Nil(have any, opts ...any) error {
	if is := core.IsNil(have); is {
		return nil
	}
	ops := DefaultOptions(opts...)
	const mHeader = "expected value to be nil"
	msg := notice.New(mHeader).Want("nil").Have("%s", ops.Dumper.Any(have))
	return AddRows(ops, msg)
}

// NotNil checks if "have" is not nil. Returns nil if it is not nil, otherwise
// returns an error with a message indicating the expected and actual values.
//
// The returned error might be one or more errors joined with [errors.Join].
func NotNil(have any, opts ...any) error {
	if is := core.IsNil(have); !is {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected non-nil value")
	return AddRows(ops, msg)
}
