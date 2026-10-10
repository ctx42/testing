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
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := Equal(tspy, 42, 44)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.ReportedError {
			t.Error("expected reported error")
		}
		if !tspy.HelperCalled {
			t.Error("expected Helper to be called")
		}
		want := "" +
			"expected values to be equal:\n" +
			"  want: (int) 42\n" +
			"  have: (int) 44\n"
		if want != tspy.Log() {
			t.Errorf(expMsg, want, tspy.Log())
		}
	})

	t.Run("error - uncomparable", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		want := []int{1}
		val := []int{1}

		// --- When ---
		have := Equal[any](tspy, want, val)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		wMsg := "" +
			"expected comparable values:\n" +
			"  want: ([]int) []int{1}\n" +
			"  have: ([]int) []int{1}\n"
		if wMsg != tspy.Log() {
			t.Errorf(expMsg, wMsg, tspy.Log())
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

	t.Run("error - different types", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := DeepEqual(tspy, int64(1), 1)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.HelperCalled {
			t.Error("expected Helper to be called")
		}
		want := "" +
			"expected values to be equal:\n" +
			"  want: (int64) 1\n" +
			"  have: (int) 1\n"
		if want != tspy.Log() {
			t.Errorf(expMsg, want, tspy.Log())
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
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := Nil(tspy, 1)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.ReportedError {
			t.Error("expected reported error")
		}
		if !tspy.HelperCalled {
			t.Error("expected Helper to be called")
		}
		want := "expected nil value:\n  have: (int) 1\n"
		if want != tspy.Log() {
			t.Errorf(expMsg, want, tspy.Log())
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

	t.Run("error - typed nil", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		var ptr *int

		// --- When ---
		pan, _ := core.WillPanic(func() { _ = NotNil(tspy, ptr) })

		// --- Then ---
		if pan != core.FailNowMsg {
			t.Errorf("want panic %q, have %#v", core.FailNowMsg, pan)
		}
		want := "expected non-nil value:\n  type: *int\n  have: nil\n"
		if tspy.Log() != want {
			t.Errorf(expMsg, want, tspy.Log())
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
		val, _ := core.WillPanic(func() { _ = Panic(tspy, func() {}) })

		// --- Then ---
		if val != core.FailNowMsg {
			t.Errorf(expMsg, core.FailNowMsg, val)
		}
		if !tspy.TriggeredFailure {
			t.Error("expected fatal failure")
		}
		if tspy.Log() != "expected fn to panic\n" {
			t.Errorf(expMsg, "expected fn to panic\n", tspy.Log())
		}
	})

	t.Run("error - nil function", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		val, _ := core.WillPanic(func() { _ = Panic(tspy, nil) })

		// --- Then ---
		if val != core.FailNowMsg {
			t.Errorf(expMsg, core.FailNowMsg, val)
		}
		if !tspy.TriggeredFailure {
			t.Error("expected fatal failure")
		}
		if tspy.Log() != "expected non-nil fn\n" {
			t.Errorf(expMsg, "expected non-nil fn\n", tspy.Log())
		}
	})
}
