// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"math"
	"strconv"

	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/notice"
)

// Epsilon checks the difference between two numbers is within a given delta.
// Returns nil if it does, otherwise it returns an error with a message
// indicating the expected and actual values.
func Epsilon[T constraints.Number](want, epsilon, have T, opts ...Option) error {
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
