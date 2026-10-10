// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// This file demonstrates how to write custom assertion helpers using
// [assert], [check], and [tester.Spy] for testing the helpers themselves.
package examples_test

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/testcases"
	"github.com/ctx42/testing/pkg/tester"

	"github.com/ctx42/testing/pkg/assert"
)

func AssertTA(t tester.T, want, have *testcases.TA, opts ...any) bool {
	t.Helper()

	ops := check.DefaultOptions(opts...)
	fName := check.FieldName(ops, "TA")

	if !assert.NotNil(t, want, fName("")) {
		return false
	}
	if !assert.NotNil(t, have, fName("")) {
		return false
	}
	if !assert.Equal(t, want.Int, have.Int, fName("Int")) {
		return false
	}
	if !assert.Equal(t, want.Str, have.Str, fName("Str")) {
		return false
	}
	if !assert.Equal(t, want.Tim, have.Tim, fName("Tim")) {
		return false
	}
	if !assert.Equal(t, want.Dur, have.Dur, fName("Dur")) {
		return false
	}
	if !assert.Equal(t, want.Loc, have.Loc, fName("Loc")) {
		return false
	}
	if !assert.Equal(t, want.TAp, have.TAp, fName("TAp")) {
		return false
	}
	return true
}

func Test_AssertTA(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		w := &testcases.TA{Int: 1, Str: "a", TAp: &testcases.TA{Int: 2, Str: "b"}}
		h := &testcases.TA{Int: 1, Str: "a", TAp: &testcases.TA{Int: 2, Str: "b"}}

		// --- When ---
		have := AssertTA(tspy, w, h)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("not equal field", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TA.TAp.Int\n" +
			"   want: 2\n" +
			"   have: 3"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		w := &testcases.TA{Int: 1, Str: "a", TAp: &testcases.TA{Int: 2, Str: "b"}}
		h := &testcases.TA{Int: 1, Str: "a", TAp: &testcases.TA{Int: 3, Str: "b"}}

		// --- When ---
		have := AssertTA(tspy, w, h)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("not equal duration", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TA.Dur\n" +
			"   want: 1\n" +
			"   have: 2"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		w := &testcases.TA{Dur: 1}
		h := &testcases.TA{Dur: 2}

		// --- When ---
		have := AssertTA(tspy, w, h)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("both want and have are nil", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		wMsg := "expected non-nil value:\n  trail: TA"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		affirm.Panic(t, func() { AssertTA(tspy, nil, nil) })
	})

	t.Run("error - have is nil", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		wMsg := "expected non-nil value:\n  trail: TA"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		w := &testcases.TA{Int: 1, Str: "a"}

		// --- When ---
		affirm.Panic(t, func() { AssertTA(tspy, w, nil) })
	})
}
