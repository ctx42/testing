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
	spy := NewSpy()

	// --- When ---
	have := spy.Capture()

	// --- Then ---
	if have != spy {
		t.Errorf("expected have to be the same instance as spy")
	}
	if have.Messages == nil {
		t.Errorf("expected Spy.Messages not to be nil")
	}
}

func Test_Spy_Helper(t *testing.T) {
	// --- Given ---
	spy := NewSpy()

	// --- When ---
	spy.Helper()

	// --- Then ---
	if !spy.HelperCalled {
		t.Errorf("expected Spy.HelperCalled to be true")
	}
}

func Test_Spy_Error(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy()

		// --- When ---
		spy.Error("a", 1, 1.1)

		// --- Then ---
		if !spy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be true")
		}
		if spy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be false")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy().Capture()

		// --- When ---
		spy.Error("a", 1, 1.1)
		spy.Error("b", 2, 2.2)

		// --- Then ---
		have := spy.Messages.String()
		want := "a 1 1.1\nb 2 2.2\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}

func Test_Spy_Errorf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy()

		// --- When ---
		spy.Errorf("a%s", "bc")

		// --- Then ---
		if !spy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be true")
		}
		if spy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be false")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy().Capture()

		// --- When ---
		spy.Errorf("a%s", "bc")
		spy.Errorf("x%s", "yz")

		// --- Then ---
		have := spy.Messages.String()
		want := "abc\nxyz\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}

func Test_Spy_Fatal(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy()

		// --- When ---
		spy.Fatal("a", 1, 1.1)

		// --- Then ---
		if spy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be false")
		}
		if !spy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be true")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy().Capture()

		// --- When ---
		spy.Fatal("a", 1, 1.1)
		spy.Fatal("b", 2, 2.2)

		// --- Then ---
		have := spy.Messages.String()
		want := "a 1 1.1\nb 2 2.2\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
		}
	})
}

func Test_Spy_Fatalf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy()

		// --- When ---
		spy.Fatalf("a%s", "bc")

		// --- Then ---
		if spy.ReportedError {
			t.Errorf("expected Spy.ReportedError to be false")
		}
		if !spy.TriggeredFailure {
			t.Errorf("expected Spy.TriggeredFailure to be true")
		}
	})

	t.Run("with capture", func(t *testing.T) {
		// --- Given ---
		spy := NewSpy().Capture()

		// --- When ---
		spy.Fatalf("a%s", "bc")
		spy.Fatalf("x%s", "yz")

		// --- Then ---
		have := spy.Messages.String()
		want := "abc\nxyz\n"
		if want != have {
			t.Errorf(sameLogMsg, want, have)
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
			spy := &Spy{
				ReportedError:    tc.ReportedError,
				TriggeredFailure: tc.TriggeredFailure,
			}

			// --- When ---
			have := spy.Failed()

			// --- Then ---
			if tc.want != have {
				wMsg := "expected Spy.Failed:\n  want: %t\n  have: %t"
				t.Errorf(wMsg, tc.want, have)
			}
		})
	}
}
