// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Greater checks the "have" value is greater than the "want" value. Returns
// true if it is, otherwise marks the test as failed, writes an error message
// to the test log and returns false.
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

// GreaterOrEqual checks the "have" value is greater or equal than the "want"
// value. Returns true if it is, otherwise marks the test as failed, writes an
// error message to the test log and returns false.
func GreaterOrEqual[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.GreaterOrEqual(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Smaller checks the "have" value is smaller than the "want" value. Returns
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

// SmallerOrEqual checks the "have" value is smaller or equal than the "want"
// value. Returns true if it is, otherwise marks the test as failed, writes an
// error message to the test log and returns false.
func SmallerOrEqual[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.SmallerOrEqual(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Delta asserts both values are within the given delta. Returns true if they
// are, otherwise marks the test as failed, writes an error message to the test
// log and returns false.
//
//	|w-h|/|w| <= delta
func Delta[T, E constraints.Number](
	t tester.T,
	want T, delta E, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if e := check.Delta(want, delta, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// DeltaSlice asserts values are within the given delta for all respective
// slice indexes. It returns true if all differences are within the delta;
// otherwise, marks the test as failed, writes an error message to the test log
// and returns false.
//
//	|w[i]-h[i]| <= delta
func DeltaSlice[T, E constraints.Number](
	t tester.T,
	want []T, delta E, have []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.DeltaSlice(want, delta, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Epsilon asserts the relative error is less than epsilon. Returns true if it
// is, otherwise marks the test as failed, writes an error message to the test
// log and returns false.
//
//	|w-h|/|w| <= epsilon
func Epsilon[T, E constraints.Number](
	t tester.T,
	want T, epsilon E, have T,
	opts ...check.Option,
) bool {

	t.Helper()
	if e := check.Epsilon(want, epsilon, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// EpsilonSlice asserts the relative error is less than epsilon for all
// respective values in the provided slices. It returns true if all differences
// are within the delta; otherwise, marks the test as failed, writes an error
// message to the test log and returns false.
//
//	|w[i]-h[i]|/|w[i]| <= epsilon
func EpsilonSlice[T, E constraints.Number](
	t tester.T,
	want []T, epsilon E, have []T,
	opts ...check.Option,
) bool {

	t.Helper()
	if err := check.EpsilonSlice(want, epsilon, have, opts...); err != nil {
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
