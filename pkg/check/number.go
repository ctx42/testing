// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/notice"
)

// Greater checks that "have" is strictly greater than "want".
//
// See assert.Greater for the assertion wrapper.
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

// GreaterOrEqual checks that "have" is greater than or equal to "want".
//
// See assert.GreaterOrEqual for the assertion wrapper.
func GreaterOrEqual[T constraints.Ordered](want, have T, opts ...any) error {
	if want <= have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected value to be greater or equal").
		Append("greater or equal to", "%v", want).
		Have("%v", have)
	return AddRows(ops, msg)
}

// Smaller checks that "have" is strictly smaller than "want".
//
// See assert.Smaller for the assertion wrapper.
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

// SmallerOrEqual checks that "have" is smaller than or equal to "want".
//
// See assert.SmallerOrEqual for the assertion wrapper.
func SmallerOrEqual[T constraints.Ordered](want, have T, opts ...any) error {
	if want >= have {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected value to be smaller or equal").
		Append("smaller or equal to", "%v", want).
		Have("%v", have)
	return AddRows(ops, msg)
}

// Delta checks that the absolute difference between "want" and "have"
// is at most "delta".
//
//	|w-h| <= delta
//
// See assert.Delta for the assertion form and options.
func Delta[T, E constraints.Number](
	want T,
	delta E,
	have T,
	opts ...any,
) error {

	ops := DefaultOptions(opts...)
	fwDelta := float64(delta)
	if fwDelta < 0 {
		msg := notice.New("expected non-negative delta").
			Append("delta", "%s", formatFloat(fwDelta))
		return AddRows(ops, msg)
	}

	wantFmt, haveFmt, hDeltaFmt, pass := deltaInt(want, have, fwDelta)
	if wantFmt == "" {
		fWant := float64(want)
		fHave := float64(have)
		fhDelta := math.Abs(fWant - fHave)
		pass = fwDelta >= fhDelta
		wantFmt, haveFmt = formatFloat(fWant), formatFloat(fHave)
		hDeltaFmt = formatFloat(fhDelta)
	}
	if pass {
		return nil
	}

	wDeltaFmt := formatFloat(fwDelta)
	msg := notice.New("expected numbers to be within the given delta").
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("want delta", "%s", wDeltaFmt).
		Append("have delta", "%s", hDeltaFmt)
	return AddRows(ops, msg)
}

// DeltaSlice checks that the absolute difference between corresponding
// elements of "want" and "have" is at most "delta".
//
//	|w[i]-h[i]| <= delta
//
// See assert.DeltaSlice.
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
			const hHeader = "" +
				"expected all numbers to be " +
				"within the given delta respectively"
			msg := notice.From(e).SetHeader(hHeader)
			return AddRows(iOps, msg)
		}
	}
	return nil
}

// Epsilon checks that the relative error between "want" and "have"
// is at most "epsilon":
//
//	|w-h|/|w| <= epsilon
//
// When want is zero, relative error is undefined: both want and have zero
// succeeds, otherwise the check fails.
//
// See assert.Epsilon for the assertion form and options.
func Epsilon[T, E constraints.Number](
	want T,
	epsilon E,
	have T,
	opts ...any,
) error {

	ops := DefaultOptions(opts...)
	fwEpsilon := float64(epsilon)
	if fwEpsilon < 0 {
		msg := notice.New("expected non-negative epsilon").
			Append("epsilon", "%s", formatFloat(fwEpsilon))
		return AddRows(ops, msg)
	}

	wantFmt, haveFmt, hEpsilonFmt, pass := epsilonInt(want, have, fwEpsilon)
	if wantFmt == "" {
		fWant := float64(want)
		fHave := float64(have)
		var fhEpsilon float64
		if fWant == 0 {
			// Relative error is undefined at zero; both zero is a match.
			fhEpsilon = math.Inf(1)
			if fHave == 0 {
				fhEpsilon = 0
			}
		} else {
			fhEpsilon = math.Abs(fWant-fHave) / math.Abs(fWant)
		}
		pass = fwEpsilon >= fhEpsilon
		wantFmt, haveFmt = formatFloat(fWant), formatFloat(fHave)
		hEpsilonFmt = formatFloat(fhEpsilon)
	}
	if pass {
		return nil
	}

	wEpsilonFmt := formatFloat(fwEpsilon)
	msg := notice.New("expected numbers to be within the given epsilon").
		Want("%s", wantFmt).
		Have("%s", haveFmt).
		Append("want epsilon", "%s", wEpsilonFmt).
		Append("have epsilon", "%s", hEpsilonFmt)
	return AddRows(ops, msg)
}

// EpsilonSlice checks that the relative error between corresponding elements
// of "want" and "have" is at most "epsilon".
//
//	|w[i]-h[i]|/|w[i]| <= epsilon
//
// See assert.EpsilonSlice.
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
			const hHeader = "" +
				"expected all numbers to be " +
				"within the given epsilon respectively"
			msg := notice.From(e).SetHeader(hHeader)
			return AddRows(iOps, msg)
		}
	}
	return nil
}

// Increasing checks that the sequence is strictly increasing.
// Use [WithIncreasingSoft] to allow equal consecutive values.
// See assert.Increasing.
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

// NotIncreasing is the inverse of [Increasing].
// See assert.NotIncreasing.
func NotIncreasing[T constraints.Ordered](seq []T, opts ...any) error {
	if err := Increasing(seq, opts...); err != nil {
		return nil //nolint:nilerr
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

// Decreasing checks that the sequence is strictly decreasing.
// Use [WithDecreasingSoft] to allow equal consecutive values.
// See assert.Decreasing.
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

// NotDecreasing is the inverse of [Decreasing].
// See assert.NotDecreasing.
func NotDecreasing[T constraints.Ordered](seq []T, opts ...any) error {
	if err := Decreasing(seq, opts...); err != nil {
		return nil //nolint:nilerr
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

// formatFloat formats f with the minimal number of digits.
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// intAbsDiff returns |a-b| computed exactly when T is an integer type, and
// false for floating-point types.
func intAbsDiff[T constraints.Number](a, b T) (*big.Int, bool) {
	aVal, bVal := reflect.ValueOf(a), reflect.ValueOf(b)
	var aInt, bInt big.Int
	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_ = aInt.SetInt64(aVal.Int())
		_ = bInt.SetInt64(bVal.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		_ = aInt.SetUint64(aVal.Uint())
		_ = bInt.SetUint64(bVal.Uint())

	default:
		return nil, false
	}
	diff := new(big.Int).Sub(&aInt, &bInt)
	return diff.Abs(diff), true
}

// deltaInt checks |want-have| <= delta exactly for integer types. It returns
// the formatted want, have and difference, and whether the check passed. For
// floating-point types, or a NaN delta, the formatted values are empty.
func deltaInt[T constraints.Number](
	want T,
	have T,
	delta float64,
) (string, string, string, bool) {

	diff, ok := intAbsDiff(want, have)
	if !ok || math.IsNaN(delta) {
		return "", "", "", false
	}
	pass := big.NewFloat(delta).Cmp(new(big.Float).SetInt(diff)) >= 0
	return fmt.Sprint(want), fmt.Sprint(have), diff.String(), pass
}

// epsilonInt checks |want-have|/|want| <= epsilon exactly for integer types.
// It returns the formatted want, have and relative error, and whether the
// check passed. For floating-point types, or a NaN epsilon, the formatted
// values are empty.
func epsilonInt[T constraints.Number](
	want T,
	have T,
	epsilon float64,
) (string, string, string, bool) {

	diff, ok := intAbsDiff(want, have)
	if !ok || math.IsNaN(epsilon) {
		return "", "", "", false
	}
	wAbs, _ := intAbsDiff(want, 0)
	wantFmt, haveFmt := fmt.Sprint(want), fmt.Sprint(have)
	if wAbs.Sign() == 0 {
		// Relative error is undefined at zero; both zero is a match.
		if diff.Sign() == 0 {
			return wantFmt, haveFmt, "0", true
		}
		return wantFmt, haveFmt, formatFloat(math.Inf(1)), false
	}
	rel := new(big.Float).SetPrec(256).SetInt(diff)
	_ = rel.Quo(rel, new(big.Float).SetPrec(256).SetInt(wAbs))
	pass := big.NewFloat(epsilon).Cmp(rel) >= 0
	relFlt, _ := rel.Float64()
	return wantFmt, haveFmt, formatFloat(relFlt), pass
}
