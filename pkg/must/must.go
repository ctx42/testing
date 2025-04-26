// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package must provide a set of helper functions which panic on error.
//
// Functions are designed to simplify error handling in test code by panicking
// on errors, reduce boilerplate, and error checking in test cases.
package must

import (
	"errors"
)

// Value is a helper for functions returning (T, error) panicking if returned
// error is not nil. Example:
//
//	fil := Value(os.Open("file.txt"))
func Value[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Values is a helper for functions returning (T1, T2, error) panicking if
// returned error is not nil. Example:
//
//	t1, t2 := Values(DoSomething())
func Values[T, TT any](val0 T, val1 TT, err error) (T, TT) {
	if err != nil {
		panic(err)
	}
	return val0, val1
}

// Nil is a helper function which panics if an error is not nil.
func Nil(err error) {
	if err != nil {
		panic(err)
	}
}

// First returns the first element in the slice or T's zero value if the slice
// is empty. It panics with value of err if err is not nil.
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

// Single returns the first element in the slice or T's zero value if the slice
// is empty. It panics with value of err if err is not nil or with value of
// errExpSingle if slices has more than one element.
func Single[T any](s []T, err error) T {
	v, err := single(s, err)
	if err != nil {
		panic(err)
	}
	return v
}

// errExpSingle is an error returned when [single] receives a slice with more
// than one element and nil err.
var errExpSingle = errors.New("expected single result")

// single returns the first element in the slice or T's zero value if the slice
// is empty. It returns T's zero value and error if err is not nil. If slice has
// more than one element it returns the first element and errExpSingle error.
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
