// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package core

import (
	"testing"
)

// sameLogMsg test log message shown when messages arte not the same.
const sameLogMsg = "expected log to be the same:\n  want: %q\n  have: %q"

func Test_NewSpy(t *testing.T) {
	// --- When ---
	have := NewSpy()

	// --- Then ---
	if have.HelperCalled {
		t.Errorf("expected Spy.HelperCalled to be false")
	}
	if have.ReportedError {
		t.Errorf("expected Spy.ReportedError to be false")
	}
	if have.TriggeredFailure {
		t.Errorf("expected Spy.TriggeredFailure to be false")
	}
	if have.Messages != nil {
		t.Errorf("expected Spy.Messages to be nil")
	}
}

func Test_Spy_Capture(t *testing.T) {
	// --- Given ---
	tspy := NewSpy()

	// --- When ---
	have := tspy.Capture()

	// --- Then ---
	if have != tspy {
		t.Errorf("expected `have` to be the same instance as spy")
	}
	if have.Messages == nil {
		t.Errorf("expected Spy.Messages not to be nil")
	}
}

func Test_Spy_Helper(t *testing.T) {
	// --- Given ---
	tspy := NewSpy()

	// --- When ---
	tspy.Helper()

	// --- Then ---
	if !tspy.HelperCalled {
		t.Errorf("expected Spy.HelperCalled to be true")
	}
}

func Test_Spy_Error(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy()

		// --- When ---
		tspy.Error("a", 1, 1.1)

		// --- Then ---
		if !tspy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be true")
		}
		if tspy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be false")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()

		// --- When ---
		tspy.Error("a", 1, 1.1)
		tspy.Error("b", 2, 2.2)

		// --- Then ---
		have := tspy.Messages.String()
		want := "a 1 1.1\nb 2 2.2\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}

func Test_Spy_Errorf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy()

		// --- When ---
		tspy.Errorf("a%s", "bc")

		// --- Then ---
		if !tspy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be true")
		}
		if tspy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be false")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()

		// --- When ---
		tspy.Errorf("a%s", "bc")
		tspy.Errorf("x%s", "yz")

		// --- Then ---
		have := tspy.Messages.String()
		want := "abc\nxyz\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})

	t.Run("message ending with new line", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()

		// --- When ---
		tspy.Errorf("a%s\n", "bc")

		// --- Then ---
		have := tspy.Messages.String()
		want := "abc\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}

func Test_Spy_Fatal(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy()

		// --- When ---
		var pan any
		func() {
			defer func() { pan = recover() }()
			tspy.Fatal("a", 1, 1.1)
		}()

		// --- Then ---
		if pan != FailNowMsg {
			t.Errorf("want panic %q, have %#v", FailNowMsg, pan)
		}
		if tspy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be false")
		}
		if !tspy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be true")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()

		// --- When ---
		func() {
			defer func() { _ = recover() }()
			tspy.Fatal("a", 1, 1.1)
		}()
		func() {
			defer func() { _ = recover() }()
			tspy.Fatal("b", 2, 2.2)
		}()

		// --- Then ---
		have := tspy.Messages.String()
		want := "a 1 1.1\nb 2 2.2\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}

func Test_Spy_Fatalf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy()

		// --- When ---
		var pan any
		func() {
			defer func() { pan = recover() }()
			tspy.Fatalf("a%s", "bc")
		}()

		// --- Then ---
		if pan != FailNowMsg {
			t.Errorf("want panic %q, have %#v", FailNowMsg, pan)
		}
		if tspy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be false")
		}
		if !tspy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be true")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()

		// --- When ---
		func() {
			defer func() { _ = recover() }()
			tspy.Fatalf("a%s", "bc")
		}()
		func() {
			defer func() { _ = recover() }()
			tspy.Fatalf("x%s", "yz")
		}()

		// --- Then ---
		have := tspy.Messages.String()
		want := "abc\nxyz\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})

	t.Run("message ending with new line", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()

		// --- When ---
		have, _ := WillPanic(func() { tspy.Fatalf("a%s\n", "bc") })

		// --- Then ---
		if have != FailNowMsg {
			t.Errorf("want panic %q, have %#v", FailNowMsg, have)
		}
		want := "abc\n"
		if want != tspy.Messages.String() {
			t.Errorf(sameLogMsg, want, tspy.Messages.String())
		}
	})
}

func Test_Spy_Failed_tabular(t *testing.T) {
	tt := []struct {
		testN string

		ReportedError    bool
		TriggeredFailure bool
		want             bool
	}{
		{"both false", false, false, false},
		{"only error", true, false, true},
		{"only fatal", false, true, true},
		{"both true", true, true, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			tspy := &Spy{
				ReportedError:    tc.ReportedError,
				TriggeredFailure: tc.TriggeredFailure,
			}

			// --- When ---
			have := tspy.Failed()

			// --- Then ---
			if tc.want != have {
				wMsg := "expected Spy.Failed:\n  want: %t\n  have: %t"
				t.Errorf(wMsg, tc.want, have)
			}
		})
	}
}

func Test_Spy_Log(t *testing.T) {
	t.Run("log is returned", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy().Capture()
		tspy.Errorf("a%s", "bc")
		func() {
			defer func() { _ = recover() }()
			tspy.Fatalf("x%s", "yz")
		}()

		// --- When ---
		have := tspy.Log()

		// --- Then ---
		want := "abc\nxyz\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})

	t.Run("panic when log capture is not turned on", func(t *testing.T) {
		// --- Given ---
		tspy := NewSpy()
		tspy.Errorf("a%s", "bc")

		// --- When ---
		have, _ := WillPanic(func() { tspy.Log() })

		// --- Then ---
		want := "Spy.Log called without log capture being on"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}
