// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package check provides equality toolkit used by assert package.
package check

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Count checks there is "count" occurrences of "what" in "where". Returns nil
// if it's, otherwise it returns an error with a message indicating the
// expected and actual values.
//
// Currently, only strings are supported.
func Count(count int, what, where any, opts ...any) error {
	if src, ok := where.(string); ok {
		var ok bool
		var subT string
		if subT, ok = what.(string); !ok {
			ops := DefaultOptions(opts...)
			const mHeader = "expected argument \"what\" to be string got %T"
			msg := notice.New(mHeader, what)
			return AddRows(ops, msg)
		}
		haveCnt := strings.Count(src, subT)
		if count == haveCnt {
			return nil
		}

		ops := DefaultOptions(opts...)
		msg := notice.New("expected string to contain substrings").
			Append("want count", "%d", count).
			Append("have count", "%d", haveCnt).
			Append("what", "%q", what).
			Append("where", "%q", where)
		return AddRows(ops, msg)
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("unsupported \"where\" type: %T", where)
	return AddRows(ops, msg)
}

// SameType checks that both arguments are of the same type. Returns nil if
// they are, otherwise it returns an error with a message indicating the
// expected and actual values.
//
// Check uses [reflect.TypeOf] equality to determine the type.
func SameType(want, have any, opts ...any) error {
	wTyp := reflect.TypeOf(want)
	hTyp := reflect.TypeOf(have)
	if wTyp == hTyp {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected same types").Want("%T", want).Have("%T", have)
	return AddRows(ops, msg)
}

// NotSameType checks that the arguments are not of the same type. Returns nil
// if they are not, otherwise it returns an error with a message indicating the
// expected and actual values.
//
// Check uses [reflect.TypeOf] equality to determine the type.
func NotSameType(want, have any, opts ...any) error {
	wTyp := reflect.TypeOf(want)
	hTyp := reflect.TypeOf(have)
	if wTyp != hTyp {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected different types").
		Want("%T", want).Have("%T", have)
	return AddRows(ops, msg)
}

// Type checks that the "src" can be type assigned to the pointer to the
// "target" (same as target, ok := src.(target)). Returns nil if it can be done,
// otherwise it returns an error. The "target" must be a pointer to a type.
func Type(target, src any, opts ...any) error {
	if target == nil {
		return notice.New("expected target to be a non-nil pointer")
	}
	tgtVal := reflect.ValueOf(target)
	tgtTyp := tgtVal.Type()
	if tgtTyp.Kind() != reflect.Ptr || tgtVal.IsNil() {
		return notice.New("expected target to be a non-nil pointer")
	}
	tgtType := tgtTyp.Elem()
	if reflect.TypeOf(src).AssignableTo(tgtType) {
		tgtVal.Elem().Set(reflect.ValueOf(src))
		return nil
	}

	tgtTypStr := fmt.Sprintf("%T", tgtVal.Interface())
	if tgtTypStr != "" && tgtTypStr[0] == '*' {
		tgtTypStr = tgtTypStr[1:]
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected type to be assignable to the target").
		Append("target", "%s", tgtTypStr).
		Append("src", "%T", src)
	return AddRows(ops, msg)
}

// Fields checks a struct or pointer to a struct "s" has "want" number of
// fields. Returns nil if it does, otherwise it returns an error with a message
// indicating the expected and actual values.
func Fields(want int, s any, opts ...any) error {
	sVal := reflect.Indirect(reflect.ValueOf(s))
	ops := DefaultOptions(opts...)
	if sVal.Kind() != reflect.Struct {
		msg := notice.New("expected struct type").Append("got type", "%T", s)
		return AddRows(ops, msg)
	}

	have := sVal.Type().NumField()
	if want == have {
		return nil
	}

	msg := notice.New("expected struct to have number of fields").
		Want("%d", want).
		Have("%d", have)
	return AddRows(ops, msg)
}
