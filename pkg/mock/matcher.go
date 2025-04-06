package mock

import (
	"fmt"
	"reflect"
)

// Matcher represents argument matcher.
type Matcher struct {
	fn   reflect.Value // Matcher function.
	desc string        // Matcher description.
}

// NewMatcher returns new instance of an [Matcher]. Takes a function as
// in [MatchBy] documentation and argument matcher description. Function panics
// on error.
func NewMatcher(fn any, desc string) *Matcher {
	return &Matcher{
		fn:   matcherFunc(fn),
		desc: desc,
	}
}

// Desc returns argument matcher description.
func (am *Matcher) Desc() string { return am.desc }

// Match runs matcher function with "have" argument and returns true if it
// matches expected value, otherwise returns false. It will panic if the "have"
// doesn't match expected type for the matcher.
func (am *Matcher) Match(have any) bool {
	expectType := am.fn.Type().In(0)

	var expectTypeNilSupported bool
	switch expectType.Kind() { // nolint:exhaustive
	case reflect.Slice, reflect.Map, reflect.Ptr:
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
