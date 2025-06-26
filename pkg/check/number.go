// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/notice"
)

// Greater checks the "want" value is greater than the "have" value. Returns
// nil if the condition is met, otherwise it returns an error with a message
// indicating the expected and actual values.
func Greater[T constraints.Ordered](want, have T, opts ...Option) error {
	if want > have {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected value to be greater").
		SetTrail(ops.Trail).
		Append("greater than", "%v", want).
		Have("%v", have)
}

// Smaller checks the "want" value is smaller than the "have" value. Returns
// nil if the condition is met, otherwise it returns an error with a message
// indicating the expected and actual values.
func Smaller[T constraints.Ordered](want, have T, opts ...Option) error {
	if want < have {
		return nil
	}
	ops := DefaultOptions(opts...)
	return notice.New("expected value to be smaller").
		SetTrail(ops.Trail).
		Append("smaller than", "%v", want).
		Have("%v", have)
}

// Epsilon checks the difference between two numbers is within a given delta.
// Returns nil if it does, otherwise it returns an error with a message
// indicating the expected and actual values.
func Epsilon[T constraints.Number](
	want, epsilon, have T,
	opts ...Option,
) error {

	fWant := float64(want)
	fHave := float64(have)
	fDelta := float64(epsilon)
	diff := math.Abs(fWant - fHave)
	if diff <= fDelta {
		return nil
	}

	ops := DefaultOptions(opts...)

	wantFmt := strconv.FormatFloat(fWant, 'f', -1, 64)
	haveFmt := strconv.FormatFloat(fHave, 'f', -1, 64)
	deltaFmt := strconv.FormatFloat(fDelta, 'f', -1, 64)
	diffFmt := strconv.FormatFloat(diff, 'f', -1, 64)
	return notice.New("expected numbers to be within given epsilon").
		SetTrail(ops.Trail).
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("epsilon", "%s", deltaFmt).
		Append("diff", "%s", diffFmt)
}

// EpsilonSlice compares two slices of numbers, "have" and "want", and checks
// if the absolute difference between corresponding elements is within the
// specified delta. It returns nil if all differences are within the delta;
// otherwise, it returns an error indicating the first index where the "have"
// slice violates the epsilon condition.
func EpsilonSlice[T constraints.Number](
	want []T,
	delta T,
	have []T,
	opts ...Option,
) error {

	if err := Len(len(want), have, opts...); err != nil {
		return err
	}

	ops := DefaultOptions(opts...)
	knd := fmt.Sprintf("%T", want)

	for i, w := range want {
		iOps := ops.ArrTrail(knd, i)
		if e := Epsilon(w, delta, have[i], WithOptions(iOps)); e != nil {
			hdr := "expected all numbers to be " +
				"within given epsilon respectively"
			return notice.From(e).SetHeader(hdr)
		}
	}
	return nil
}

// Increasing checks if the given sequence has values in the increasing order.
// You may use the [WithIncreasingSoft] option to allow consecutive values to
// be equal. It returns an error if the sequence is not increasing.
func Increasing[T constraints.Ordered](seq []T, opts ...Option) error {
	if len(seq) <= 1 {
		return nil
	}

	ops := DefaultOptions(opts...)
	knd := fmt.Sprintf("%T", seq)
	var mode string

	var cmp func(T, T) bool
	if ops.IncreaseSoft {
		mode = "soft"
		cmp = func(c, p T) bool { return p <= c }
	} else {
		mode = "strict"
		cmp = func(c, p T) bool { return p < c }
	}

	prv := seq[0]
	for i := 1; i < len(seq); i++ {
		cur := seq[i]
		if !cmp(cur, prv) {
			iOps := ops.ArrTrail(knd, i)
			return notice.New("expected an increasing sequence").
				SetTrail(iOps.Trail).
				Append("mode", "%s", mode).
				Append("previous", "%v", prv).
				Append("current", "%v", cur)
		}
		prv = cur
	}
	return nil
}

// NotIncreasing is inverse of [Increasing].
func NotIncreasing[T constraints.Ordered](seq []T, opts ...Option) error {
	if err := Increasing(seq, opts...); err != nil {
		return nil
	}
	ops := DefaultOptions(opts...)
	var mode string
	if ops.IncreaseSoft {
		mode = "soft"
	} else {
		mode = "strict"
	}
	return notice.New("expected a not increasing sequence").
		Append("mode", "%s", mode)
}

// Decreasing checks if the given sequence has values in the decreasing order.
// You may use the [WithDecreasingSoft] option to allow consecutive values to
// be equal. It returns an error if the sequence is not decreasing.
func Decreasing[T constraints.Ordered](seq []T, opts ...Option) error {
	if len(seq) <= 1 {
		return nil
	}
	ops := DefaultOptions(opts...)
	knd := fmt.Sprintf("%T", seq)
	var mode string

	var cmp func(T, T) bool
	if ops.DecreaseSoft {
		mode = "soft"
		cmp = func(c, p T) bool { return p >= c }
	} else {
		mode = "strict"
		cmp = func(c, p T) bool { return p > c }
	}

	prv := seq[0]
	for i := 1; i < len(seq); i++ {
		cur := seq[i]
		if !cmp(cur, prv) {
			iOps := ops.ArrTrail(knd, i)
			return notice.New("expected a decreasing sequence").
				SetTrail(iOps.Trail).
				Append("mode", "%s", mode).
				Append("previous", "%v", prv).
				Append("current", "%v", cur)
		}
		prv = cur
	}
	return nil
}

// NotDecreasing is inverse of [Decreasing].
func NotDecreasing[T constraints.Ordered](seq []T, opts ...Option) error {
	if err := Decreasing(seq, opts...); err != nil {
		return nil
	}

	ops := DefaultOptions(opts...)
	var mode string
	if ops.DecreaseSoft {
		mode = "soft"
	} else {
		mode = "strict"
	}
	return notice.New("expected a not decreasing sequence").
		Append("mode", "%s", mode)
}
