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

// Greater checks the "have" value is greater than the "want" value. Returns
// nil if the condition is met, otherwise it returns an error with a message
// indicating the expected and actual values.
func Greater[T constraints.Ordered](want, have T, opts ...any) error {
	if want < have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected value to be greater").
		Append("greater than", "%v", want).
		Have("%v", have)
	return AddRows(ops, msg)
}

// GreaterOrEqual checks the "have" value is greater or equal than the "want"
// value. Returns nil if the condition is met, otherwise it returns an error
// with a message indicating the expected and actual values.
func GreaterOrEqual[T constraints.Ordered](want, have T, opts ...any) error {
	if want <= have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected value to be greater or equal").
		Append("greater or equal than", "%v", want).
		Have("%v", have)
	return AddRows(ops, msg)
}

// Smaller checks the "have" value is smaller than the "want" value. Returns
// nil if the condition is met, otherwise it returns an error with a message
// indicating the expected and actual values.
func Smaller[T constraints.Ordered](want, have T, opts ...any) error {
	if want > have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected value to be smaller").
		Append("smaller than", "%v", want).
		Have("%v", have)
	return AddRows(ops, msg)
}

// SmallerOrEqual checks the "have" value is smaller or equal than the "want"
// value. Returns nil if the condition is met, otherwise it returns an error
// with a message indicating the expected and actual values.
func SmallerOrEqual[T constraints.Ordered](want, have T, opts ...any) error {
	if want >= have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected value to be smaller or equal").
		Append("smaller or equal than", "%v", want).
		Have("%v", have)
	return AddRows(ops, msg)
}

// Delta checks the both values are within the given delta. Returns nil if they
// are, otherwise it returns an error with a message indicating the expected
// and actual values.
//
//	|w-h| <= delta
func Delta[T, E constraints.Number](
	want T, delta E, have T,
	opts ...any,
) error {

	fWant := float64(want)
	fHave := float64(have)
	fwDelta := float64(delta)
	fhDelta := math.Abs(fWant - fHave)
	if fwDelta >= fhDelta {
		return nil
	}

	ops := DefaultOptions(opts...)

	wantFmt := strconv.FormatFloat(fWant, 'f', -1, 64)
	haveFmt := strconv.FormatFloat(fHave, 'f', -1, 64)
	wDeltaFmt := strconv.FormatFloat(fwDelta, 'f', -1, 64)
	hDeltaFmt := strconv.FormatFloat(fhDelta, 'f', -1, 64)
	msg := notice.New("expected numbers to be within the given delta").
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("want delta", "%s", wDeltaFmt).
		Append("have delta", "%s", hDeltaFmt)
	return AddRows(ops, msg)
}

// DeltaSlice checks values are within the given delta for all respective
// slice indexes. It returns nil if all differences are within the delta;
// otherwise, it returns an error indicating the first index where the "have"
// slice violates the epsilon condition.
//
//	|w[i]-h[i]| <= delta
func DeltaSlice[T, E constraints.Number](
	want []T,
	delta E,
	have []T,
	opts ...any,
) error {

	if err := Len(len(want), have, opts...); err != nil {
		return err
	}

	fDelta := float64(delta)
	ops := DefaultOptions(opts...)
	knd := fmt.Sprintf("%T", want)

	for i, w := range want {
		iOps := ops.ArrTrail(knd, i)
		if e := Delta(w, fDelta, have[i], WithOptions(iOps)); e != nil {
			const hHeader = "expected all numbers to be " +
				"within the given delta respectively"
			msg := notice.From(e).SetHeader(hHeader)
			return AddRows(iOps, msg)
		}
	}
	return nil
}

// Epsilon checks the relative error is less than epsilon. Returns nil if it
// does, otherwise it returns an error with a message indicating the expected
// and actual values.
//
//	|w-h|/|w| <= epsilon
func Epsilon[T, E constraints.Number](
	want T, epsilon E, have T,
	opts ...any,
) error {

	fWant := float64(want)
	fHave := float64(have)
	fwEpsilon := float64(epsilon)
	fhEpsilon := math.Abs(fWant-fHave) / math.Abs(fWant)
	if fwEpsilon >= fhEpsilon {
		return nil
	}

	ops := DefaultOptions(opts...)

	wantFmt := strconv.FormatFloat(fWant, 'f', -1, 64)
	haveFmt := strconv.FormatFloat(fHave, 'f', -1, 64)
	wEpsilonFmt := strconv.FormatFloat(fwEpsilon, 'f', -1, 64)
	hEpsilonFmt := strconv.FormatFloat(fhEpsilon, 'f', -1, 64)
	msg := notice.New("expected numbers to be within the given epsilon").
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("want epsilon", "%s", wEpsilonFmt).
		Append("have epsilon", "%s", hEpsilonFmt)
	return AddRows(ops, msg)
}

// EpsilonSlice checks the relative error is less than epsilon for all
// respective values in the provided slices. It returns nil if all differences
// are within the epsilon; otherwise, it returns an error indicating the first
// index where the "have" slice violates the epsilon condition.
//
//	|w[i]-h[i]|/|w[i]| <= epsilon
func EpsilonSlice[T, E constraints.Number](
	want []T,
	epsilon E,
	have []T,
	opts ...any,
) error {

	if err := Len(len(want), have, opts...); err != nil {
		return err
	}

	fEpsilon := float64(epsilon)
	ops := DefaultOptions(opts...)
	knd := fmt.Sprintf("%T", want)

	for i, w := range want {
		iOps := ops.ArrTrail(knd, i)
		if e := Epsilon(w, fEpsilon, have[i], WithOptions(iOps)); e != nil {
			const hHeader = "expected all numbers to be " +
				"within the given epsilon respectively"
			msg := notice.From(e).SetHeader(hHeader)
			return AddRows(iOps, msg)
		}
	}
	return nil
}

// Increasing checks if the given sequence has values in the increasing order.
// You may use the [WithIncreasingSoft] option to allow consecutive values to
// be equal. It returns an error if the sequence is not increasing.
func Increasing[T constraints.Ordered](seq []T, opts ...any) error {
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
			msg := notice.New("expected an increasing sequence").
				Append("mode", "%s", mode).
				Append("previous", "%v", prv).
				Append("current", "%v", cur)
			return AddRows(iOps, msg)
		}
		prv = cur
	}
	return nil
}

// NotIncreasing is inverse of [Increasing].
func NotIncreasing[T constraints.Ordered](seq []T, opts ...any) error {
	if err := Increasing(seq, opts...); err != nil {
		return nil // nolint: nilerr
	}
	ops := DefaultOptions(opts...)
	var mode string
	if ops.IncreaseSoft {
		mode = "soft"
	} else {
		mode = "strict"
	}
	msg := notice.New("expected a not increasing sequence").
		Append("mode", "%s", mode)
	return AddRows(ops, msg)
}

// Decreasing checks if the given sequence has values in the decreasing order.
// You may use the [WithDecreasingSoft] option to allow consecutive values to
// be equal. It returns an error if the sequence is not decreasing.
func Decreasing[T constraints.Ordered](seq []T, opts ...any) error {
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
			msg := notice.New("expected a decreasing sequence").
				Append("mode", "%s", mode).
				Append("previous", "%v", prv).
				Append("current", "%v", cur)
			return AddRows(iOps, msg)
		}
		prv = cur
	}
	return nil
}

// NotDecreasing is inverse of [Decreasing].
func NotDecreasing[T constraints.Ordered](seq []T, opts ...any) error {
	if err := Decreasing(seq, opts...); err != nil {
		return nil // nolint: nilerr
	}

	ops := DefaultOptions(opts...)
	var mode string
	if ops.DecreaseSoft {
		mode = "soft"
	} else {
		mode = "strict"
	}
	msg := notice.New("expected a not decreasing sequence").
		Append("mode", "%s", mode)
	return AddRows(ops, msg)
}
