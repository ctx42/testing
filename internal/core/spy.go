// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package core

import (
	"fmt"
	"strings"
)

// T defines an interface for types that mimic the behavior of [testing.T] in
// Go's testing framework, used for capturing error and failure reports in
// tests. It is implemented by [Spy] to mock test interactions, enabling
// verification of test behavior without triggering real test failures.
type T interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
}

// Spy is a testing utility that captures and tracks calls to error and failure
// reporting methods, typically used to mock `*testing.T` in unit tests. It
// records messages logged via Error, Errorf, Fatal, and Fatalf, and tracks
// whether errors or fatal failures were reported. Use Spy to verify test
// behavior without triggering real test failures.
type Spy struct {
	ReportedError    bool   // Tracks if Error or Errorf was called.
	TriggeredFailure bool   // Tracks if Fatal or Fatalf was called.
	Messages         string // Accumulated log of all error and failure messages.
}

// Helper is a no-op method that satisfies the [testing.TB.Helper] interface,
// allowing [Spy] to be used in contexts expecting a testing helper.
func (spy *Spy) Helper() {}

// Error records a non-fatal error with the provided arguments, appending them
// as a space-separated message to Messages, followed by a newline. It sets
// [Spy.ReportedError] to true.
func (spy *Spy) Error(args ...any) {
	spy.ReportedError = true
	spy.Messages += spy.log(args...)
}

// Errorf records a non-fatal error using the provided format string and
// arguments, appending the formatted message to Messages with a newline. It
// sets [Spy.ReportedError] to true.
func (spy *Spy) Errorf(format string, args ...any) {
	spy.ReportedError = true
	spy.Messages += fmt.Sprintf(format, args...) + "\n"
}

// Fatal records a fatal error with the provided arguments, appending them as a
// space-separated message to Messages, followed by a newline. It sets
// [Spy.TriggeredFailure] to true.
func (spy *Spy) Fatal(args ...any) {
	spy.TriggeredFailure = true
	spy.Messages += spy.log(args...)
}

// Fatalf records a fatal error using the provided format string and arguments,
// appending the formatted message to Messages with a newline. It sets
// [Spy.TriggeredFailure] to true.
func (spy *Spy) Fatalf(format string, args ...any) {
	spy.TriggeredFailure = true
	spy.Messages += fmt.Sprintf(format, args...) + "\n"
}

// log formats the provided arguments into a space-separated string, terminated
// with a newline, for consistent message logging.
func (spy *Spy) log(args ...any) string {
	var buf strings.Builder
	for i, arg := range args {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprintf("%v", arg))
	}
	buf.WriteByte('\n')
	return buf.String()
}

// Failed reports whether an error or failure has been recorded, returning true
// if either ReportedError or TriggeredFailure is set. This mimics the behavior
// of [testing.T.Failed], allowing tests to check the Spy's state without
// modifying it.
func (spy *Spy) Failed() bool {
	return spy.ReportedError || spy.TriggeredFailure
}
