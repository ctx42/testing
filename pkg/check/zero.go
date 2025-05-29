// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"reflect"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/notice"
)

// Zero checks "have" is the zero value for its type. Returns nil if it is,
// otherwise, it returns an error with a message indicating the expected
// and actual values.
func Zero(have any, opts ...Option) error {
	if core.IsNil(have) {
		return zeroError(have, opts...)
	}
	type z interface{ IsZero() bool }
	if zero, ok := have.(z); ok {
		if zero.IsZero() {
			return nil
		}
		return zeroError(have, opts...)
	}
	val := reflect.ValueOf(have)
	if val.IsValid() && val.IsZero() {
		return nil
	}
	return zeroError(have, opts...)
}

// zeroError returns error for non-zero value of have.
func zeroError(have any, opts ...Option) error {
	ops := DefaultOptions(opts...)
	return notice.New("expected argument to be zero value").
		SetTrail(ops.Trail).
		Want("<zero>").
		Have("%#v", have)
}

// NotZero checks "have" is not the zero value for its type. Returns nil if it
// is, otherwise it returns an error with a message indicating the expected and
// actual values.
func NotZero(have any, opts ...Option) error {
	if Zero(have) != nil {
		return nil // nolint: nilerr
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected argument not to be zero value").
		SetTrail(ops.Trail).
		Want("<non-zero>").
		Have("%#v", have)
}
