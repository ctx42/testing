// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ctx42/testing/internal/core"
)

// AnyString is a helper matching any argument value of string type.
var AnyString = MatchOfType("string")

// AnyInt is a helper matching any argument value of integer type.
var AnyInt = MatchOfType("int")

// AnyBool is a helper matching any argument value of boolean type.
var AnyBool = MatchOfType("bool")

// AnyCtx matches any non-nil context.
var AnyCtx = MatchBy(func(ctx context.Context) bool {
	return ctx != nil
})

// MatchSame matches two generic pointers point to the same object using
// [is.SamePointers].
func MatchSame(want any) *Matcher {
	return MatchBy(func(have any) bool { return core.Same(want, have) })
}

// MatchBy constructs an [Matcher] instance which validates arguments using
// a given function. The function must be accepting a single argument (of the
// expected type) and return a true if argument matches expectations or false
// when it doesn't. If function doesn't match the required signature, [MatchBy]
// panics.
//
// Examples:
//
//	fn := func(have int) bool { ... }
//	fn := func(have float64) bool { ... }
//	fn := func(have ExampleItf) bool { ... }
//	fn := func(have ExampleType) bool { ... }
//	fn := func(have *ExampleType) bool { ... }
//
// MatchBy can be used to match complex mocked method argument like function,
// structure, channel, map, ...
//
// Example:
//
//	MatchBy(func(req *http.Request) bool { return req.Host == "localhost" })
func MatchBy(fn any) *Matcher {
	val := reflect.ValueOf(fn)
	desc := fmt.Sprintf(
		"[mock.MatchBy=func(%s) bool]",
		val.Type().In(0).String(),
	)
	return NewMatcher(fn, desc)
}

// MatchOfType constructs an argument matcher (Matcher) instance which
// ensures argument is of given type.
//
// Examples:
//
//	MatchOfType("int")
//	MatchOfType("string")
//	MatchOfType("mock.ExampleType")
//	MatchOfType("*mock.ExampleType")
//	MatchOfType("map[string]interface {}")
//
// The MatchOfType will not match if the string is an interface name.
func MatchOfType(typ string) *Matcher {
	fn := func(have any) bool {
		haveTyp := reflect.TypeOf(have)
		haveStr := haveTyp.String()
		return typ == haveStr
	}
	desc := fmt.Sprintf("[mock.MatchOfType=%s]", typ)
	return NewMatcher(fn, desc)
}

// MatchType constructs an argument matcher ([Matcher]) instance which
// ensures argument is of the same type as the [MatchType] argument.
//
// Examples:
//
//	MatchType(42)
//	MatchType(true)
//	MatchType("string")
//	MatchType(mock.ExampleType{})
//	MatchType(*mock.ExampleType{})
func MatchType(typ any) *Matcher {
	typTyp := reflect.TypeOf(typ)
	typStr := typTyp.String()

	fn := func(have any) bool {
		haveTyp := reflect.TypeOf(have)
		haveStr := haveTyp.String()
		return typStr == haveStr
	}
	desc := fmt.Sprintf("[mock.MatchType=%s]", typStr)
	return NewMatcher(fn, desc)
}

// MatchErrorContain constructs an argument matcher ([Matcher]) instance
// which ensures argument is a non nil error with given message.
func MatchErrorContain(want string) *Matcher {
	return MatchBy(func(err error) bool {
		return strings.Contains(err.Error(), want)
	})
}

// MatchError creates an [Matcher] that verifies the argument is a non-nil
// error. The "want" parameter specifies the expected error behavior:
//   - If want is a string, the error's message (via Error()) is matched.
//   - If want is an error, [errors.Is] is used to check.
//
// It will panic if "want" is neither.
func MatchError(want any) *Matcher {
	var mby *Matcher
	switch w := want.(type) {
	case string:
		mby = MatchBy(func(have error) bool { return w == have.Error() })
	case error:
		mby = MatchBy(func(have error) bool { return errors.Is(have, w) })
	default:
		panic("mock: MatchError: invalid type")
	}
	return mby
}

// AnySlice is a helper to create slice of length cnt of [mock.Any] values.
func AnySlice(cnt int) []any {
	var str []any
	for i := 0; i < cnt; i++ {
		str = append(str, Any)
	}
	return str
}
