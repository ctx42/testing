// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package check provides equality toolkit used by assert package.
package check

import (
	"reflect"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Count checks there is "count" occurrences of "what" in "where". Returns nil
// if it's, otherwise it returns an error with a message indicating the
// expected and actual values.
//
// Currently, only strings are supported.
func Count(count int, what, where any, opts ...Option) error {
	if src, ok := where.(string); ok {
		var ok bool
		var subT string
		if subT, ok = what.(string); !ok {
			ops := DefaultOptions(opts...)
			const mHeader = "expected argument \"what\" to be string got %T"
			return notice.New(mHeader, what).SetTrail(ops.Trail)
		}
		haveCnt := strings.Count(src, subT)
		if count == haveCnt {
			return nil
		}

		ops := DefaultOptions(opts...)
		return notice.New("expected string to contain substrings").
			SetTrail(ops.Trail).
			Append("want count", "%d", count).
			Append("have count", "%d", haveCnt).
			Append("what", "%q", what).
			Append("where", "%q", where)
	}

	ops := DefaultOptions(opts...)
	return notice.New("unsupported \"where\" type: %T", where).
		SetTrail(ops.Trail)
}

// Type checks that both arguments are of the same type. Returns nil if they
// are, otherwise it returns an error with a message indicating the expected
// and actual values.
func Type(want, have any, opts ...Option) error {
	wTyp := reflect.TypeOf(want)
	hTyp := reflect.TypeOf(have)
	if wTyp == hTyp {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected same types").
		SetTrail(ops.Trail).
		Want("%T", want).
		Have("%T", have)
}

// Fields checks a struct or pointer to a struct "s" has "want" number of
// fields. Returns nil if it does, otherwise it returns an error with a message
// indicating the expected and actual values.
func Fields(want int, s any, opts ...Option) error {
	sVal := reflect.Indirect(reflect.ValueOf(s))
	ops := DefaultOptions(opts...)
	if sVal.Kind() != reflect.Struct {
		return notice.New("expected struct type").
			SetTrail(ops.Trail).
			Append("got type", "%T", s)
	}

	have := sVal.Type().NumField()
	if want == have {
		return nil
	}

	return notice.New("expected struct to have number of fields").
		SetTrail(ops.Trail).
		Want("%d", want).
		Have("%d", have)
}
