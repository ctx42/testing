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
		spy := core.NewSpy()

		// --- When ---
		have := Equal(spy, 42, 42)

		// --- Then ---
		if !have || spy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Equal(spy, 42, 44)

		// --- Then ---
		if have || !spy.Failed() {
			t.Error("expected test error")
		}
		if !spy.ReportedError {
			t.Error("expected test error")
		}
	})
}

func Test_DeepEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := DeepEqual(spy, []int{42}, []int{42})

		// --- Then ---
		if !have || spy.Failed() {
			t.Error("expected test error")
		}
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := DeepEqual(spy, []int{42}, []int{44})

		// --- Then ---
		if have || !spy.Failed() {
			t.Error("expected test error")
		}
		if !spy.ReportedError {
			t.Error("expected test error")
		}
	})
}

func Test_Nil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Nil(spy, nil)

		// --- Then ---
		if !have || spy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()
		err := errors.New("m0")

		// --- When ---
		have := Nil(spy, err)

		// --- Then ---
		if have || !spy.Failed() {
			t.Error("expected test error")
		}
		if !spy.ReportedError {
			t.Error("expected test error")
		}
	})
}

func Test_NotNil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := NotNil(spy, errors.New("m0"))

		// --- Then ---
		if !have || spy.Failed() {
			t.Error("expected passed test")
		}
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy().Capture()
		var err error

		// --- When ---
		have := NotNil(spy, err)

		// --- Then ---
		if have || !spy.Failed() {
			t.Error("expected failed test")
		}
		if !spy.TriggeredFailure {
			t.Error("expected test failure")
		}
	})
}

func Test_Panic(t *testing.T) {
	const expMsg = "expected values to be equal:\n  want: %q\n  have: %q"

	t.Run("success string message", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Panic(spy, func() { panic("abc") })

		// --- Then ---
		if spy.ReportedError {
			t.Error("expected passed test")
		}
		if *have != "abc" {
			t.Errorf(expMsg, "abc", *have)
		}
	})

	t.Run("success error", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Panic(spy, func() { panic(errors.New("abc")) })

		// --- Then ---
		if spy.ReportedError {
			t.Error("expected passed test")
		}
		if *have != "abc" {
			t.Errorf(expMsg, "abc", *have)
		}
	})

	t.Run("success other type", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Panic(spy, func() { panic(123) })

		// --- Then ---
		if spy.ReportedError {
			t.Error("expected error")
		}
		if *have != "123" {
			t.Errorf(expMsg, "123", *have)
		}
	})

	t.Run("success function panics with nil", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Panic(spy, func() { panic(nil) }) // nolint: govet

		// --- Then ---
		if spy.ReportedError {
			t.Error("expected passed test")
		}
		want := (&runtime.PanicNilError{}).Error()
		if *have != want {
			t.Errorf(expMsg, want, *have)
		}
	})

	t.Run("error function does not panic", func(t *testing.T) {
		// --- Given ---
		spy := core.NewSpy()

		// --- When ---
		have := Panic(spy, func() {})

		// --- Then ---
		if !spy.Failed() {
			t.Error("expected failed test")
		}
		if have != nil {
			t.Errorf(expMsg, "", *have)
		}
	})
}
