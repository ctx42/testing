// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// Same checks that "want" and "have" are pointers to the same object.
//
// Sameness is based on type + value equality (via [internal/core.Same]).
// Works for pointers, slices, maps, funcs. Arrays always fail.
//
// See [assert.Same] for the assertion wrapper.
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

// NotSame checks that "want" and "have" are not pointers to the same object.
//
// NotSame checks that "want" and "have" are not pointers to the same object.
//
// Both must be pointer variables. Sameness check is type + value equality.
//
// See [assert.NotSame] for the assertion wrapper.
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
