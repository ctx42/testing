// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package must provides helpers that panic on error.
//
// These functions eliminate repetitive `if err != nil { panic(err) }`
// boilerplate in test setup and assertions, making test code more
// concise while keeping failures explicit.
//
// The helpers integrate naturally with [tester.T] and the assertion
// packages when writing test fixtures or one-off checks that must
// succeed.
//
// See the package [README] for motivation and patterns. See
// [examples_test.go] for executable demonstrations.
//
// Key entry points:
//   - [Value] / [Values] — return values or panic on error
//   - [Nil] — assert no error (panic on any error)
//   - [First] / [Single] — extract elements from slices with error handling
package must

import (
	"errors"
)

// Value returns val if err is nil; otherwise it panics with err.
//
// Generic equivalent of the common "must succeed" pattern for single-value
// + error results.
func Value[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Values returns val0 and val1 if err is nil; otherwise it panics with err.
//
// Useful for functions that return two values plus an error.
func Values[T, TT any](val0 T, val1 TT, err error) (T, TT) {
	if err != nil {
		panic(err)
	}
	return val0, val1
}

// Nil panics with err if err is not nil. The simplest form for asserting
// that an operation produced no error.
func Nil(err error) {
	if err != nil {
		panic(err)
	}
}

// First returns the first element of s (or the zero value if empty).
// Panics if err is not nil.
//
// See [Single] when you require exactly one element.
func First[T any](s []T, err error) T {
	v, err := single(s, err)
	if errors.Is(err, errExpSingle) {
		return v
	}
	if err != nil {
		panic(err)
	}
	return v
}

// Single returns the single element of s (or the zero value if empty).
// Panics if err is not nil, or with a specific error if s has more than
// one element.
//
// See [First] when you only need the first element.
func Single[T any](s []T, err error) T {
	v, err := single(s, err)
	if err != nil {
		panic(err)
	}
	return v
}

// errExpSingle is an error returned when [single] receives a slice with more
// than one element and nil err.
var errExpSingle = errors.New("expected a single result")

// single returns the first element in the slice or T's zero value if the slice
// is empty. It returns T's zero value and error if err is not nil. If a slice
// has more than one element, it returns the first element and errExpSingle
// error.
func single[T any](s []T, err error) (T, error) {
	var t T
	if err != nil {
		return t, err
	}
	switch len(s) {
	case 0:
		return t, nil
	case 1:
		return s[0], nil
	default:
		return s[0], errExpSingle
	}
}
