// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package affirm

import (
	"errors"
	"runtime"
	"testing"

	"github.com/ctx42/testing/internal/core"
)

// expMsg is the message for a failed test assertion.
const expMsg = "expected values to be equal:\n  want: %q\n  have: %q"

func Test_Equal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Equal(tspy, 42, 42)

		// --- Then ---
		if !have || tspy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Equal(tspy, 42, 44)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.ReportedError {
			t.Error("expected reported error")
		}
	})
}

func Test_DeepEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		want := []int{42}
		val := []int{42}

		// --- When ---
		have := DeepEqual(tspy, want, val)

		// --- Then ---
		if !have || tspy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		want := []int{42}
		val := []int{44}

		// --- When ---
		have := DeepEqual(tspy, want, val)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.ReportedError {
			t.Error("expected reported error")
		}
	})
}

func Test_Nil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Nil(tspy, nil)

		// --- Then ---
		if !have || tspy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Nil(tspy, errors.New("m0"))

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.ReportedError {
			t.Error("expected reported error")
		}
	})
}

func Test_NotNil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := NotNil(tspy, errors.New("m0"))

		// --- Then ---
		if !have || tspy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error - nil", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		var err error

		// --- When ---
		var pan any
		var have bool
		func() {
			defer func() { pan = recover() }()
			have = NotNil(tspy, err)
		}()

		// --- Then ---
		if pan != core.FailNowMsg {
			t.Errorf("want panic %q, have %#v", core.FailNowMsg, pan)
		}
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.TriggeredFailure {
			t.Error("expected test failure")
		}
		want := "expected non-nil value:\n  have: nil\n"
		if tspy.Messages.String() != want {
			t.Errorf("want log %q, have %q", want, tspy.Messages.String())
		}
	})
}

func Test_Panic(t *testing.T) {
	t.Run("string message", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() { panic("abc") })

		// --- Then ---
		if tspy.ReportedError {
			t.Error("expected passed test")
		}
		if *have != "abc" {
			t.Errorf(expMsg, "abc", *have)
		}
	})

	t.Run("success error", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() { panic(errors.New("abc")) })

		// --- Then ---
		if tspy.ReportedError {
			t.Error("expected passed test")
		}
		if *have != "abc" {
			t.Errorf(expMsg, "abc", *have)
		}
	})

	t.Run("integer type", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() { panic(123) })

		// --- Then ---
		if tspy.ReportedError {
			t.Error("expected passed test")
		}
		if *have != "123" {
			t.Errorf(expMsg, "123", *have)
		}
	})

	t.Run("panic nil", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() { panic(nil) }) //nolint:govet

		// --- Then ---
		if tspy.ReportedError {
			t.Error("expected passed test")
		}
		want := (&runtime.PanicNilError{}).Error()
		if *have != want {
			t.Errorf(expMsg, want, *have)
		}
	})

	t.Run("error - no panic", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		val, _ := core.WillPanic(func() { Panic(tspy, func() {}) })

		// --- Then ---
		if val != core.FailNowMsg {
			t.Errorf(expMsg, core.FailNowMsg, val)
		}
		if !tspy.TriggeredFailure {
			t.Error("expected fatal failure")
		}
		if want, have := "expected fn to panic\n", tspy.Log(); want != have {
			t.Errorf(expMsg, want, have)
		}
	})

	t.Run("error - nil function", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		val, _ := core.WillPanic(func() { Panic(tspy, nil) })

		// --- Then ---
		if val != core.FailNowMsg {
			t.Errorf(expMsg, core.FailNowMsg, val)
		}
		if !tspy.TriggeredFailure {
			t.Error("expected fatal failure")
		}
	})
}

func Test_Equal_uncomparable(t *testing.T) {
	// --- Given ---
	tspy := core.NewSpy().Capture()

	// --- When ---
	have := Equal[any](tspy, []int{1}, []int{1})

	// --- Then ---
	if have {
		t.Error("expected false")
	}
	want := "" +
		"expected comparable values:\n" +
		"  want: ([]int) []int{1}\n" +
		"  have: ([]int) []int{1}\n"
	if want != tspy.Log() {
		t.Errorf(expMsg, want, tspy.Log())
	}
}

func Test_messages_tabular(t *testing.T) {
	tt := []struct {
		testN string

		fn   func(t core.T) bool
		want string
	}{
		{
			"Equal",
			func(t core.T) bool { return Equal(t, int64(1), int64(2)) },
			"" +
				"expected values to be equal:\n" +
				"  want: (int64) 1\n" +
				"  have: (int64) 2\n",
		},
		{
			"DeepEqual",
			func(t core.T) bool { return DeepEqual(t, int64(1), 1) },
			"" +
				"expected values to be equal:\n" +
				"  want: (int64) 1\n" +
				"  have: (int) 1\n",
		},
		{
			"Nil",
			func(t core.T) bool { return Nil(t, 1) },
			"expected nil value:\n  have: (int) 1\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			tspy := core.NewSpy().Capture()

			// --- When ---
			have := tc.fn(tspy)

			// --- Then ---
			if have {
				t.Error("expected false")
			}
			if !tspy.HelperCalled {
				t.Error("expected Helper to be called")
			}
			if tc.want != tspy.Log() {
				t.Errorf(expMsg, tc.want, tspy.Log())
			}
		})
	}
}
