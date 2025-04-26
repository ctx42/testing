// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/internal/constraints"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

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
