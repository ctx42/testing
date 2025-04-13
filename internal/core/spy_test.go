package core

import (
	"testing"
)

func Test_Spy_Error(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Error("a", 1, 1.1)

		// --- Then ---
		have := spy.Messages
		want := "a 1 1.1\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
		if !spy.ReportedError {
			t.Errorf("expected spy.ReportedError to be true")
		}
		if spy.TriggeredFailure {
			t.Errorf("expected spy.TriggeredFailure to be false")
		}
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Error("a", 1, 1.1)
		spy.Error("b", 2, 2.2)

		// --- Then ---
		have := spy.Messages
		want := "a 1 1.1\nb 2 2.2\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
	})
}

func Test_Spy_Errorf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Errorf("a%s", "bc")

		// --- Then ---
		have := spy.Messages
		want := "abc\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
		if !spy.ReportedError {
			t.Errorf("expected spy.ReportedError to be true")
		}
		if spy.TriggeredFailure {
			t.Errorf("expected spy.TriggeredFailure to be false")
		}
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Errorf("a%s", "bc")
		spy.Errorf("x%s", "yz")

		// --- Then ---
		have := spy.Messages
		want := "abc\nxyz\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
	})
}

func Test_Spy_Fatal(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Fatal("a", 1, 1.1)

		// --- Then ---
		have := spy.Messages
		want := "a 1 1.1\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
		if spy.ReportedError {
			t.Errorf("expected spy.ReportedError to be false")
		}
		if !spy.TriggeredFailure {
			t.Errorf("expected spy.TriggeredFailure to be true")
		}
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Fatal("a", 1, 1.1)
		spy.Fatal("b", 2, 2.2)

		// --- Then ---
		have := spy.Messages
		want := "a 1 1.1\nb 2 2.2\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
	})
}

func Test_Spy_Fatalf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Fatalf("a%s", "bc")

		// --- Then ---
		have := spy.Messages
		want := "abc\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
		}
		if spy.ReportedError {
			t.Errorf("expected spy.ReportedError to be false")
		}
		if !spy.TriggeredFailure {
			t.Errorf("expected spy.TriggeredFailure to be true")
		}
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		spy := &Spy{}

		// --- When ---
		spy.Fatalf("a%s", "bc")
		spy.Fatalf("x%s", "yz")

		// --- Then ---
		have := spy.Messages
		want := "abc\nxyz\n"
		if want != have {
			wMsg := "expected same:\n  want: %q\n  have: %q"
			t.Errorf(wMsg, want, have)
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
