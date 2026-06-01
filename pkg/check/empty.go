// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"reflect"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
)

// Empty checks that "have" is empty.
//
// Empty values are:
//   - nil
//   - int(0)
//   - float64(0)
//   - float32(0)
//   - false
//   - len(array) == 0
//   - len(slice) == 0
//   - len(map) == 0
//   - len(chan) == 0
//   - time.Time{}
//
// See [assert.Empty] for the assertion form.
func Empty(have any, opts ...any) error {
	if isEmpty(have) {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected argument to be empty").
		Want(dump.ValEmpty).
		Have("%#v", have)
	return AddRows(ops, msg)
}

// isEmpty returns true if "have" is empty.
func isEmpty(have any) bool {
	if is := core.IsNil(have); is {
		return true
	}

	val := reflect.ValueOf(have)
	switch val.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice:
		if val.Len() == 0 {
			return true
		}

	case reflect.Pointer:
		return isEmpty(val.Elem().Interface())

	default:
		zero := reflect.Zero(val.Type())
		if reflect.DeepEqual(have, zero.Interface()) {
			return true
		}
	}

	return false
}

// NotEmpty checks that "have" is not empty.
//
// See [Empty] and [assert.NotEmpty] for related functions.
func NotEmpty(have any, opts ...any) error {
	if !isEmpty(have) {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected non-empty value")
	return AddRows(ops, msg)
}
