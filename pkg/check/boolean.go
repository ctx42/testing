// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"github.com/ctx42/testing/pkg/notice"
)

// True checks that "have" is true.
//
// See [assert.True] for the assertion wrapper and the package documentation
// for option handling and customization.
func True(have bool, opts ...any) error {
	if !have {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected value to be true")
		return AddRows(ops, msg)
	}
	return nil
}

// False checks that "have" is false.
//
// See [assert.False] for the assertion wrapper and the package documentation
// for option handling and customization.
func False(have bool, opts ...any) error {
	if have {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected value to be false")
		return AddRows(ops, msg)
	}
	return nil
}
