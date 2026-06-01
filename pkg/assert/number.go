// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Greater asserts that "have" is greater than "want".
//
// See [check.Greater] for the error-returning form.
func Greater[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.Greater(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// GreaterOrEqual asserts that "have" is greater than or equal to "want".
//
// See [check.GreaterOrEqual] for the error-returning form.
func GreaterOrEqual[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.GreaterOrEqual(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Smaller asserts that "have" is smaller than "want".
//
// See [check.Smaller] for the error-returning form.
func Smaller[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.Smaller(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// SmallerOrEqual asserts that "have" is smaller than or equal to "want".
//
// See [check.SmallerOrEqual] for the error-returning form.
func SmallerOrEqual[T constraints.Ordered](
	t tester.T,
	want, have T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.SmallerOrEqual(want, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Delta asserts that the absolute relative difference between "want" and
// "have" is at most "delta":
//
//	|w-h|/|w| <= delta
//
// See [check.Delta] for the error-returning form and options.
func Delta[T, E constraints.Number](
	t tester.T,
	want T,
	delta E,
	have T,
	opts ...any,
) bool {

	t.Helper()
	if e := check.Delta(want, delta, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// DeltaSlice asserts that the absolute difference between corresponding
// elements of "want" and "have" is at most "delta":
//
//	|w[i]-h[i]| <= delta
//
// See [check.DeltaSlice].
func DeltaSlice[T, E constraints.Number](
	t tester.T,
	want []T,
	delta E,
	have []T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.DeltaSlice(want, delta, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Epsilon asserts that the relative error between "want" and "have" is less
// than "epsilon":
//
//	|w-h|/|w| <= epsilon
//
// See [check.Epsilon] for the error-returning form and options.
func Epsilon[T, E constraints.Number](
	t tester.T,
	want T,
	epsilon E,
	have T,
	opts ...any,
) bool {

	t.Helper()
	if e := check.Epsilon(want, epsilon, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// EpsilonSlice asserts that the relative error between corresponding elements
// of "want" and "have" is less than "epsilon":
//
//	|w[i]-h[i]|/|w[i]| <= epsilon
//
// See [check.EpsilonSlice].
func EpsilonSlice[T, E constraints.Number](
	t tester.T,
	want []T,
	epsilon E,
	have []T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.EpsilonSlice(want, epsilon, have, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Increasing asserts that the given sequence is in strictly increasing order.
// Use [check.WithIncreasingSoft] to allow equal consecutive values.
// See [check.Increasing].
func Increasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.Increasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// NotIncreasing is the inverse of [Increasing].
// See [check.NotIncreasing].
func NotIncreasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.NotIncreasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// Decreasing asserts that the given sequence is in strictly decreasing order.
// Use [check.WithDecreasingSoft] to allow equal consecutive values.
// See [check.Decreasing].
func Decreasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.Decreasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}

// NotDecreasing is the inverse of [Decreasing].
// See [check.NotDecreasing].
func NotDecreasing[T constraints.Ordered](
	t tester.T,
	seq []T,
	opts ...any,
) bool {

	t.Helper()
	if err := check.NotDecreasing(seq, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
