// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package affirm

import (
	"errors"
	"runtime"
	"testing"

	"github.com/ctx42/testing/internal/core"
)

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

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Equal(tspy, 42, 44)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected test error")
		}
		if !tspy.ReportedError {
			t.Error("expected test error")
		}
	})
}

func Test_DeepEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := DeepEqual(tspy, []int{42}, []int{42})

		// --- Then ---
		if !have || tspy.Failed() {
			t.Error("expected test error")
		}
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := DeepEqual(tspy, []int{42}, []int{44})

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected test error")
		}
		if !tspy.ReportedError {
			t.Error("expected test error")
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

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		err := errors.New("m0")

		// --- When ---
		have := Nil(tspy, err)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected test error")
		}
		if !tspy.ReportedError {
			t.Error("expected test error")
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

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		var err error

		// --- When ---
		have := NotNil(tspy, err)

		// --- Then ---
		if have || !tspy.Failed() {
			t.Error("expected failed test")
		}
		if !tspy.TriggeredFailure {
			t.Error("expected test failure")
		}
	})
}

func Test_Panic(t *testing.T) {
	const expMsg = "expected values to be equal:\n  want: %q\n  have: %q"

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

	t.Run("other type", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() { panic(123) })

		// --- Then ---
		if tspy.ReportedError {
			t.Error("expected error")
		}
		if *have != "123" {
			t.Errorf(expMsg, "123", *have)
		}
	})

	t.Run("function panics with nil", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() { panic(nil) }) // nolint: govet

		// --- Then ---
		if tspy.ReportedError {
			t.Error("expected passed test")
		}
		want := (&runtime.PanicNilError{}).Error()
		if *have != want {
			t.Errorf(expMsg, want, *have)
		}
	})

	t.Run("error - function does not panic", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Panic(tspy, func() {})

		// --- Then ---
		if !tspy.Failed() {
			t.Error("expected failed test")
		}
		if have != nil {
			t.Errorf(expMsg, "", *have)
		}
	})
}
