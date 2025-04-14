// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package core

import (
	"bytes"
	"fmt"
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
type Spy struct {
	HelperCalled     bool          // Tracks if Helper was called.
	ReportedError    bool          // Tracks if Error or Errorf was called.
	TriggeredFailure bool          // Tracks if Fatal or Fatalf was called.
	Messages         *bytes.Buffer // Log messages if set.
}

// NewSpy returns new instance of [Spy].
func NewSpy() *Spy { return &Spy{} }

// Capture turns on collection of messages when Error, Errorf, Fatal, and
// Fatalf are called.
func (spy *Spy) Capture() *Spy {
	spy.Messages = &bytes.Buffer{}
	return spy
}

// Helper is a no-op method that satisfies the [testing.TB.Helper] interface,
// allowing [Spy] to be used in contexts expecting a testing helper.
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
// arguments, appending the formatted message to Messages with a newline. It
// sets [Spy.ReportedError] to true.
func (spy *Spy) Errorf(format string, args ...any) {
	spy.ReportedError = true
	if spy.Messages != nil {
		_, _ = fmt.Fprintf(spy.Messages, format, args...)
		_ = spy.Messages.WriteByte('\n')
	}
}

// Fatal records a fatal error with the provided arguments, appending them as a
// space-separated message to Messages, followed by a newline. It sets
// [Spy.TriggeredFailure] to true.
func (spy *Spy) Fatal(args ...any) {
	spy.TriggeredFailure = true
	if spy.Messages != nil {
		_, _ = fmt.Fprintln(spy.Messages, args...)
	}
}

// Fatalf records a fatal error using the provided format string and arguments,
// appending the formatted message to Messages with a newline. It sets
// [Spy.TriggeredFailure] to true.
func (spy *Spy) Fatalf(format string, args ...any) {
	spy.TriggeredFailure = true
	if spy.Messages != nil {
		_, _ = fmt.Fprintf(spy.Messages, format, args...)
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
