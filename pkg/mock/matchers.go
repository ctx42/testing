// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

// AnyString matches any argument whose dynamic type is string.
var AnyString = MatchOfType("string")

// AnyInt matches any argument whose dynamic type is int.
var AnyInt = MatchOfType("int")

// AnyBool matches any argument whose dynamic type is bool.
var AnyBool = MatchOfType("bool")

// AnyCtx matches any non-nil [context.Context].
var AnyCtx = MatchBy(func(ctx context.Context) bool {
	return ctx != nil
})

// MatchSame returns a matcher that reports whether the argument is the same
// pointer as want (using the same rules as check.Same).
func MatchSame(want any) *Matcher {
	return MatchBy(func(have any) bool { return core.Same(want, have) })
}

// MatchBy creates a [Matcher] backed by the supplied predicate. The function
// must have the signature func(T) bool for some type T; any other shape
// causes a panic at construction time.
//
// MatchBy is the foundation for all custom matchers and for the predefined ones
// ([MatchOfType], [MatchError], Any*, ...). The function runs while the mock is
// locked, so it must not call the mock.
//
// Example:
//
//	MatchBy(func(req *http.Request) bool { return req.Host == "localhost" })
func MatchBy(fn any) *Matcher {
	am := NewMatcher(fn, "")
	format := "[mock.MatchBy=func(%s) bool]"
	am.desc = fmt.Sprintf(format, am.fn.Type().In(0))
	return am
}

// MatchOfType returns a matcher that accepts only arguments whose
// reflect.Type.String() exactly equals typ.
//
// Examples:
//
//	MatchOfType("int")
//	MatchOfType("*http.Request")
//	MatchOfType("map[string]interface {}")
//
// It does not match interface types by name; use a custom MatchBy for that.
func MatchOfType(typ string) *Matcher {
	fn := func(have any) bool {
		haveTyp := reflect.TypeOf(have)
		if haveTyp == nil {
			return false
		}
		return typ == haveTyp.String()
	}
	desc := fmt.Sprintf("[mock.MatchOfType=%s]", typ)
	return NewMatcher(fn, desc)
}

// MatchType returns a matcher that accepts arguments whose dynamic type
// equals the type of typ (using reflect.TypeOf).
//
// Examples:
//
//	MatchType(42)
//	MatchType(true)
//	MatchType("string")
//	MatchType(mock.ExampleType{})
func MatchType(typ any) *Matcher {
	typTyp := reflect.TypeOf(typ)
	if typTyp == nil {
		panic("mock: MatchType: typ must not be nil")
	}
	fn := func(have any) bool {
		return reflect.TypeOf(have) == typTyp
	}
	desc := fmt.Sprintf("[mock.MatchType=%s]", typTyp)
	return NewMatcher(fn, desc)
}

// MatchErrorContain returns a matcher that accepts a non-nil error whose
// Error() string contains the substring want.
func MatchErrorContain(want string) *Matcher {
	return MatchBy(func(err error) bool {
		return err != nil && strings.Contains(err.Error(), want)
	})
}

// MatchError returns a matcher for non-nil errors. When want is a string it
// matches the error's Error() text exactly; when want is an error it uses
// errors.Is. Any other type for want causes a panic at construction.
func MatchError(want any) *Matcher {
	var am *Matcher
	switch w := want.(type) {
	case string:
		am = MatchBy(func(have error) bool {
			return have != nil && w == have.Error()
		})

	case error:
		am = MatchBy(func(have error) bool {
			return have != nil && errors.Is(have, w)
		})

	default:
		panic("mock: MatchError: invalid type")
	}
	return am
}

// AnySlice returns a slice of length cnt filled with [Any] sentinels.
// Useful when building expectations for variadic methods or slices.
func AnySlice(cnt int) []any {
	var str []any
	for range cnt {
		str = append(str, Any)
	}
	return str
}
