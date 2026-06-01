// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package check provides composable checks that return errors rather than
// calling methods on [testing.T].
//
// See the Design section in the root README for the overall layered
// architecture (assert built on check built on notice).
//
// The customization model supports both global defaults and per-use
// overrides:
//   - Global type checkers via [RegisterTypeChecker] (affect all checks
//     unless overridden).
//   - Per-check overrides via [WithTypeChecker], [WithTrailChecker], etc.
//   - Full option objects via [WithOptions] and [DefaultOptions].
//
// These checks form the foundation of the [assert] package. They are useful
// when building custom assertion helpers or when performing multiple checks
// before deciding how to report failure.
//
// On success, most functions return nil. On failure they return a
// *[notice.Notice] containing a structured message.
//
// Checks accept variadic options (see [DefaultOptions]) that control how values
// are rendered and how errors are formatted.
package check

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Count checks that there are "count" occurrences of "what" in "where".
//
// Currently only strings are supported for "where".
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

// SameType checks that both arguments are of the same type.
//
// On match, it returns the value cast to T. Type comparison uses
// [reflect.TypeOf] equality.
func SameType[T any](want T, have any, opts ...any) (T, error) {
	wTyp := reflect.TypeOf(want)
	hTyp := reflect.TypeOf(have)
	if wTyp == hTyp {
		return have.(T), nil // nolint: forcetypeassert
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected same types").Want("%T", want).Have("%T", have)
	var zero T
	return zero, AddRows(ops, msg)
}

// NotSameType checks that the arguments are not of the same type.
//
// Type comparison uses [reflect.TypeOf] equality.
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

// Type checks that "src" is assignable to the pointer "target"
// (equivalent to `ok := src.(target)`). "target" must be a non-nil pointer.
func Type(target, src any, opts ...any) error {
	if target == nil {
		return notice.New("expected target to be a non-nil pointer")
	}
	tgtVal := reflect.ValueOf(target)
	tgtTyp := tgtVal.Type()
	if tgtTyp.Kind() != reflect.Pointer || tgtVal.IsNil() {
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

// Fields checks that a struct (or pointer to struct) "s" has exactly "want"
// number of fields.
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
