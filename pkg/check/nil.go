// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// Nil checks that "have" is nil (including typed nils such as (*T)(nil)).
//
// See [assert.Nil] for the assertion wrapper.
func Nil(have any, opts ...any) error {
	if is := core.IsNil(have); is {
		return nil
	}
	ops := DefaultOptions(opts...)
	const mHeader = "expected value to be nil"
	msg := notice.New(mHeader).Want("nil").Have("%s", ops.Dumper.Any(have))
	return AddRows(ops, msg)
}

// NotNil checks that "have" is not nil.
//
// See [assert.NotNil] for the assertion wrapper.
func NotNil(have any, opts ...any) error {
	if is := core.IsNil(have); !is {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected non-nil value")
	return AddRows(ops, msg)
}
