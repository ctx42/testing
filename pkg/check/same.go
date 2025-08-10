// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// Same checks "want" and "have" are generic pointers and that both of them
// reference the same object. Returns nil if they are, otherwise it returns an
// error with a message indicating the expected and actual values.
//
// Pointer variable sameness is determined based on the equality of both type
// and value. It works with pointers to objects, slices, maps and functions.
// For arrays, it always returns error.
func Same(want, have any, opts ...any) error {
	if core.Same(want, have) {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected same pointers").
		Want("%p %#v", want, want).
		Have("%p %#v", have, have)
	return AddRows(ops, msg)
}

// NotSame checks "want" and "have" are generic pointers and that both of them
// reference the same object. Returns nil if it is, otherwise it returns an
// error with a message indicating the expected and actual values.
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func NotSame(want, have any, opts ...any) error {
	if Same(want, have) == nil {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected different pointers").
			Want("%p %#v", want, want).
			Have("%p %#v", have, have)
		return AddRows(ops, msg)
	}
	return nil
}
