// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock

import (
	"fmt"
	"reflect"
)

// Matcher is used for custom argument matching in expectations.
//
// Users normally obtain instances via [MatchBy], [MatchOfType], [MatchError],
// etc. rather than calling NewMatcher directly.
type Matcher struct {
	fn   reflect.Value // Matcher function.
	desc string        // Matcher description.
}

// NewMatcher creates a Matcher from a predicate function and a human-readable
// description (used in diff output). The function must have the shape
// required by [MatchBy]; otherwise NewMatcher panics.
func NewMatcher(fn any, desc string) *Matcher {
	return &Matcher{
		fn:   matcherFunc(fn),
		desc: desc,
	}
}

// Desc returns the human-readable description of the matcher (shown in
// diagnostic output when arguments fail to match).
func (am *Matcher) Desc() string { return am.desc }

// Match reports whether have satisfies the predicate. It panics if have has
// the wrong dynamic type for the matcher function.
func (am *Matcher) Match(have any) bool {
	expectType := am.fn.Type().In(0)

	var expectTypeNilSupported bool
	switch expectType.Kind() { // nolint:exhaustive
	case reflect.Slice, reflect.Map, reflect.Pointer:
		expectTypeNilSupported = true
	case reflect.Interface, reflect.Func, reflect.Chan:
		expectTypeNilSupported = true
	default:
	}

	typ := reflect.TypeOf(have)
	var val reflect.Value
	if typ == nil {
		val = reflect.New(expectType).Elem()
	} else {
		val = reflect.ValueOf(have)
	}

	if typ == nil && !expectTypeNilSupported {
		panic("attempting to call matcher with nil for non-nil expected type")
	}
	if typ == nil || typ.AssignableTo(expectType) {
		result := am.fn.Call([]reflect.Value{val})
		return result[0].Bool()
	}
	return false
}

// matcherFunc takes a function as in [MatchBy] documentation and returns its
// [reflect.Value]. Function panics on error.
func matcherFunc(fn any) reflect.Value {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("mock: \"%T\" is not a match function", fn))
	}
	if fnType.NumIn() != 1 {
		format := "mock: match function %#v does not take exactly one argument"
		panic(fmt.Sprintf(format, fn))
	}
	if fnType.NumOut() != 1 || fnType.Out(0).Kind() != reflect.Bool {
		format := "mock: match function %#v does not return a bool"
		panic(fmt.Sprintf(format, fn))
	}
	return reflect.ValueOf(fn)
}
