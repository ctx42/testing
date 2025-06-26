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
		h := have[i]
		if e := Epsilon(w, delta, h, WithOptions(iOps)); e != nil {
			hdr := "expected all numbers in a slice to be within given epsilon " +
				"respectively"
			return notice.From(e).SetHeader(hdr)
		}
	}
	return nil
}
