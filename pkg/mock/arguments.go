// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Arguments hold an array of method arguments or return values.
type Arguments []any

// Get returns the argument at the specified index. Panics when the index is
// out of bounds.
func (args Arguments) Get(idx int) any {
	if idx+1 > len(args) {
		format := "[mock] arguments: Get(%d) out of range %d max"
		panic(fmt.Sprintf(format, idx, len(args)-1))
	}
	return args[idx]
}

// Equal gets whether the objects match the specified arguments. Panics on
// error.
func (args Arguments) Equal(haves ...any) bool {
	al := len(args)
	hl := len(haves)
	if al != hl {
		format := "[must] arguments: argument lengths do not match %d != %d"
		panic(fmt.Sprintf(format, al, hl))
	}
	for i, arg := range args {
		if !reflect.DeepEqual(arg, haves[i]) {
			return false
		}
	}
	return true
}

// Diff gets a string describing the differences between the expected arguments
// and the specified values. Returns a diff string and the number of
// differences found.
//
// nolint: cyclop
func (args Arguments) Diff(vs []any) ([]string, int) {
	var out []string
	var diffCnt int

	// Pick the longer slice.
	cnt := len(args)
	if len(vs) > cnt {
		cnt = len(vs)
	}

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
			msg := "%d: FAIL: %s%s != %s"
			msg = fmt.Sprintf(msg, i, am.Desc(), panicMsg, haveFmt)
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
		msg := fmt.Sprintf("%d: FAIL: %s != %s", i, wantFmt, haveFmt)
		out = append(out, msg)
	}
	return out, diffCnt
}

// String gets the argument at the specified index cast to string. Panics for
// invalid index or when an argument cannot be cast to a string. If the index
// is set to -1, the method returns a complete string representation of the
// argument types.
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
	format := "[mock] arguments: String(%d) is of type \"%T\" not string"
	panic(fmt.Sprintf(format, idx, val))
}

// Int gets the argument at the specified index cast to int. Panics for invalid
// index or when an argument cannot be cast to an int.
func (args Arguments) Int(idx int) int {
	val := args.Get(idx)
	if got, ok := val.(int); ok {
		return got
	}
	format := "[mock] arguments: Int(%d) is of type \"%T\" not int"
	panic(fmt.Sprintf(format, idx, val))
}

// TODO(rz): Add Float32
// TODO(rz): Add Float64

// Error gets the argument at the specified index. Panics if there is no
// argument, or if the argument is of the wrong type.
func (args Arguments) Error(idx int) error {
	val := args.Get(idx)
	if val == nil {
		return nil
	}
	if got, ok := val.(error); ok {
		return got
	}
	format := "[mock] arguments: Error(%d) is of type \"%T\" not error"
	panic(fmt.Sprintf(format, idx, val))
}

// Bool gets the argument at the specified index. Panics if there is no
// argument, or if the argument is of the wrong type.
func (args Arguments) Bool(idx int) bool {
	val := args.Get(idx)
	if got, ok := val.(bool); ok {
		return got
	}
	format := "[mock] arguments: Bool(%d) is of type \"%T\" not bool"
	panic(fmt.Sprintf(format, idx, val))
}
