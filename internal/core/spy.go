// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package core

import (
	"bytes"
	"fmt"
	"strings"
)

// T defines an interface capturing a subset of [testing.TB] methods, used to
// test helper functions that accept `*testing.T` and provide reusable assertions
// for test cases. Implemented by [Spy], it allows verification of interactions
// between test helpers and `*testing.T` behavior without causing actual test
// errors or failures.
type T interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
}

// Spy is a testing utility that captures and tracks calls to error and failure
// reporting methods, used to mock `*testing.T` when testing helper functions.
// It logs calls to Error, Errorf, Fatal, and Fatalf, allowing verification of
// test behavior without causing actual test failures.
//
// Spy is not safe for concurrent use.
type Spy struct {
	HelperCalled     bool          // Tracks if Helper was called.
	ReportedError    bool          // Tracks if Error or Errorf was called.
	TriggeredFailure bool          // Tracks if Fatal or Fatalf was called.
	Messages         *bytes.Buffer // Log messages if set.
}

// FailNowMsg is the panic message used by [Spy.Fatal] and [Spy.Fatalf] so
// helpers stop after a fatal path, matching [tester.Spy]. Real [testing.TB]
// uses runtime.Goexit rather than this string.
const FailNowMsg = "FailNow was called directly"

// NewSpy returns new instance of [Spy].
func NewSpy() *Spy { return &Spy{} }

// Capture turns on the collection of messages when Error, Errorf, Fatal, and
// Fatalf are called.
func (spy *Spy) Capture() *Spy {
	spy.Messages = &bytes.Buffer{}
	return spy
}

// Helper records the call in [Spy.HelperCalled].
func (spy *Spy) Helper() { spy.HelperCalled = true }

// Error records a non-fatal error with the provided arguments, appending them
// as a space-separated message to Messages, followed by a newline. It sets
// [Spy.ReportedError] to true.
func (spy *Spy) Error(args ...any) {
	spy.ReportedError = true
	if spy.Messages != nil {
		_, _ = fmt.Fprintln(spy.Messages, args...)
	}
}

// Errorf records a non-fatal error using the provided format string and
// arguments, appending the formatted message to Messages, followed by a
// newline unless it already ends with one. It sets [Spy.ReportedError] to
// true.
func (spy *Spy) Errorf(format string, args ...any) {
	spy.ReportedError = true
	spy.logf(format, args...)
}

// Fatal records a fatal error with the provided arguments, appending them as a
// space-separated message to Messages, followed by a newline. It sets
// [Spy.TriggeredFailure] to true and panics with [FailNowMsg]. Unlike
// [testing.TB], which calls runtime.Goexit, the panic can be recovered by
// the helper under test.
func (spy *Spy) Fatal(args ...any) {
	spy.TriggeredFailure = true
	if spy.Messages != nil {
		_, _ = fmt.Fprintln(spy.Messages, args...)
	}
	panic(FailNowMsg)
}

// Fatalf records a fatal error using the provided format string and arguments,
// appending the formatted message to Messages, followed by a newline unless
// it already ends with one. It sets [Spy.TriggeredFailure] to true and panics
// with [FailNowMsg].
func (spy *Spy) Fatalf(format string, args ...any) {
	spy.TriggeredFailure = true
	spy.logf(format, args...)
	panic(FailNowMsg)
}

// logf appends the formatted message to Messages, when capture is on, and
// terminates it with a newline unless it already ends with one.
func (spy *Spy) logf(format string, args ...any) {
	if spy.Messages == nil {
		return
	}
	msg := fmt.Sprintf(format, args...)
	_, _ = spy.Messages.WriteString(msg)
	if !strings.HasSuffix(msg, "\n") {
		_ = spy.Messages.WriteByte('\n')
	}
}

// Failed reports whether an error or failure has been recorded, returning true
// if either ReportedError or TriggeredFailure is set. This mimics the behavior
// of [testing.T.Failed], allowing tests to check the Spy's state without
// modifying it.
func (spy *Spy) Failed() bool {
	return spy.ReportedError || spy.TriggeredFailure
}

// Log returns the logged messages. If logging has not been enabled via
// [Spy.Capture], it will panic to indicate misuse.
func (spy *Spy) Log() string {
	if spy.Messages == nil {
		panic("Spy.Log called without log capture being on")
	}
	return spy.Messages.String()
}
