// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Greater checks the "want" value is greater than the "have" value. Returns
// true if it is, otherwise marks the test as failed, writes an error message to
// the test log and returns false.
func Greater[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.Greater(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Smaller checks the "want" value is smaller than the "have" value. Returns
// true if it is, otherwise marks the test as failed, writes an error message
// to the test log and returns false.
func Smaller[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.Smaller(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Epsilon asserts the difference between two numbers is within a given delta.
// Returns true if it is, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
func Epsilon[T constraints.Number](
	t tester.T,
	want, delta, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if e := check.Epsilon(want, delta, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// EpsilonSlice compares two slices of numbers, "have" and "want", and checks
// if the absolute difference between corresponding elements is within the
// specified delta. It returns true if all differences are within the delta;
// otherwise, marks the test as failed, writes an error message to the test log
// and returns false.
func EpsilonSlice[T constraints.Number](
	t tester.T,
	want []T, delta T, have []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.EpsilonSlice(want, delta, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Increasing checks if the given sequence has values in the increasing order.
// You may use the [check.WithIncreasingSoft] option to allow consecutive
// values to be equal. It returns true if the sequence is increasing otherwise,
// marks the test as failed, writes an error message to the test log and
// returns false.
func Increasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.Increasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// NotIncreasing is inverse of [Increasing].
func NotIncreasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.NotIncreasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Decreasing checks if the given sequence has values in the decreasing order.
// You may use the [check.WithDecreasingSoft] option to allow consecutive
// values to be equal. It returns true if the sequence is decreasing otherwise,
// marks the test as failed, writes an error message to the test log and
// returns false.
func Decreasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.Decreasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// NotDecreasing is inverse of [Decreasing].
func NotDecreasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.NotDecreasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
