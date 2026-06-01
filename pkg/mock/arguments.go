// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// Arguments holds a slice of method arguments or return values.
//
// It provides typed accessors (Get, String, Int, ...) that panic with rich
// [notice.Notice] diagnostics on misuse, and Diff for detailed mismatch
// reporting used by the expectation engine.
type Arguments []any

// Get returns the value at the given index. It panics with a descriptive
// [notice.Notice] when the index is out of range.
func (args Arguments) Get(idx int) any {
	if idx+1 > len(args) {
		mHeader := "[mock] arguments: Get(%d) out of range %d max"
		msg := notice.New(mHeader, idx, len(args)-1)
		panic(msg)
	}
	return args[idx]
}

// Equal reports whether the provided values match the arguments using
// reflect.DeepEqual. It panics (with a [notice.Notice]) when the lengths
// differ.
func (args Arguments) Equal(haves ...any) bool {
	al := len(args)
	hl := len(haves)
	if al != hl {
		mHeader := "[must] arguments: argument lengths do not match %d != %d"
		msg := notice.New(mHeader, al, hl)
		panic(msg)
	}
	for i, arg := range args {
		if !reflect.DeepEqual(arg, haves[i]) {
			return false
		}
	}
	return true
}

// Diff produces a human-readable comparison between the expected arguments
// and the supplied values. It returns the lines of the diff and the count of
// mismatches. Used internally by the expectation matching logic.
//
// nolint: cyclop
func (args Arguments) Diff(vs []any) ([]string, int) {
	var out []string
	var diffCnt int

	// Pick the longer slice.
	cnt := max(len(vs), len(args))

	for i := 0; i < cnt; i++ {
		var want, have any = "(Missing)", "(Missing)"
		wantFmt, haveFmt := "(Missing)", "(Missing)"

		if len(vs) > i {
			have = vs[i]
			if _, ok := have.(context.Context); ok {
				haveFmt = fmt.Sprintf("(%[1]T=%[1]T)", have)
			} else {
				haveFmt = fmt.Sprintf("(%[1]T=%#[1]v)", have)
			}
		}

		if len(args) > i {
			want = args[i]
			if want == Any {
				wantFmt = "(any=mock.Any)"
			} else if _, ok := have.(context.Context); ok {
				wantFmt = fmt.Sprintf("(%[1]T=%[1]T)", want)
			} else {
				wantFmt = fmt.Sprintf("(%[1]T=%#[1]v)", want)
			}
		}

		if am, ok := want.(*Matcher); ok {
			var match bool
			var panicMsg string
			func() {
				defer func() {
					if r := recover(); r != nil {
						format := " {panic: %#v}"
						panicMsg = fmt.Sprintf(format, r)
					}
				}()
				match = am.Match(have)
			}()

			if match {
				msg := "%d: PASS: %s == %s"
				msg = fmt.Sprintf(msg, i, am.Desc(), haveFmt)
				out = append(out, msg)
				continue
			}

			diffCnt++
			msg := formatFail(i, am.Desc()+panicMsg, haveFmt)
			out = append(out, msg)
			continue
		}

		// Normal checking.
		if haveFmt != "(Missing)" &&
			(reflect.DeepEqual(want, Any) || reflect.DeepEqual(want, have)) {
			msg := fmt.Sprintf("%d: PASS: %s == %s", i, wantFmt, haveFmt)
			out = append(out, msg)
			continue
		}

		// Not match
		diffCnt++
		msg := formatFail(i, wantFmt, haveFmt)
		out = append(out, msg)
	}
	return out, diffCnt
}

// String returns the argument at idx as a string. When idx == -1 it returns
// a comma-separated list of the argument types instead. Panics with a
// [notice.Notice] on out-of-range access or wrong type.
func (args Arguments) String(idx int) string {
	if idx == -1 {
		// Return a string representation of the arg types.
		var argsStr []string
		for _, arg := range args {
			argsStr = append(argsStr, fmt.Sprintf("%T", arg))
		}
		return strings.Join(argsStr, ", ")
	}

	val := args.Get(idx)
	if got, ok := val.(string); ok {
		return got
	}
	mHeader := "[mock] arguments: String(%d) is of type \"%T\" not string"
	msg := notice.New(mHeader, idx, val)
	panic(msg)
}

// Int returns the argument at idx as int. Panics with a [notice.Notice] on
// out-of-range access or wrong type.
func (args Arguments) Int(idx int) int {
	val := args.Get(idx)
	if got, ok := val.(int); ok {
		return got
	}
	mHeader := "[mock] arguments: Int(%d) is of type \"%T\" not int"
	msg := notice.New(mHeader, idx, val)
	panic(msg)
}

// Float32 returns the argument at idx as float32. Panics with a
// [notice.Notice] on out-of-range access or wrong type.
func (args Arguments) Float32(idx int) float32 {
	val := args.Get(idx)
	if got, ok := val.(float32); ok {
		return got
	}
	mHeader := "[mock] arguments: Float32(%d) is of type \"%T\" not float32"
	msg := notice.New(mHeader, idx, val)
	panic(msg)
}

// Float64 returns the argument at idx as float64. Panics with a
// [notice.Notice] on out-of-range access or wrong type.
func (args Arguments) Float64(idx int) float64 {
	val := args.Get(idx)
	if got, ok := val.(float64); ok {
		return got
	}
	mHeader := "[mock] arguments: Float64(%d) is of type \"%T\" not float64"
	msg := notice.New(mHeader, idx, val)
	panic(msg)
}

// Error returns the argument at idx as error (nil is allowed). Panics with a
// [notice.Notice] on out-of-range access or wrong type.
func (args Arguments) Error(idx int) error {
	val := args.Get(idx)
	if val == nil {
		return nil
	}
	if got, ok := val.(error); ok {
		return got
	}
	mHeader := "[mock] arguments: Error(%d) is of type \"%T\" not error"
	msg := notice.New(mHeader, idx, val)
	panic(msg)
}

// Bool returns the argument at idx as bool. Panics with a [notice.Notice] on
// out-of-range access or wrong type.
func (args Arguments) Bool(idx int) bool {
	val := args.Get(idx)
	if got, ok := val.(bool); ok {
		return got
	}
	mHeader := "[mock] arguments: Bool(%d) is of type \"%T\" not bool"
	msg := notice.New(mHeader, idx, val)
	panic(msg)
}

// formatFail returns a single- or multi-line failure description.
// When either side is long, it uses a vertical layout for readability.
func formatFail(i int, left, right string) string {
	const maxLen = 80
	if len(left) > maxLen || len(right) > maxLen {
		format := "%d: FAIL:\n    want: %s\n    have: %s"
		return fmt.Sprintf(format, i, left, right)
	}
	return fmt.Sprintf("%d: FAIL: %s != %s", i, left, right)
}
