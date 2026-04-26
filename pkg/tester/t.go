// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester

import (
	"context"
)

// TODO(rz): add missing methods.
// var _ testing.TB = (T)(nil)

// T is a subset of [testing.TB] interface it has been defined to allow mocking.
type T interface {
	// Cleanup registers a function to be called when the test and all its
	// subtests complete. Cleanup functions will be called in last added,
	// first called order.
	Cleanup(func())

	// Error is equivalent to Log followed by Fail.
	Error(args ...any)

	// Errorf is equivalent to Logf followed by Fail.
	Errorf(format string, args ...any)

	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...any)

	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...any)

	// FailNow marks the function as having failed and stops its execution
	// by calling [runtime.Goexit] (which then runs all deferred calls in the
	// current goroutine). Execution will continue at the next test or
	// benchmark. FailNow must be called from the goroutine running the test or
	// benchmark function, not from other goroutines created during the test.
	// Calling FailNow does not stop.
	FailNow()

	// Failed reports whether the function has failed.
	Failed() bool

	// Helper marks the calling function as a test helper function. When
	// printing file and line information, that function will be skipped.
	// Helper may be called simultaneously from multiple goroutines.
	Helper()

	// Log formats its arguments using default formatting, analogous to Println,
	// and records the text in the error log. For tests, the text will be
	// printed only if the test fails or the -test.v flag is set. For
	// benchmarks, the text is always printed to avoid having performance
	// depend on the value of the -test.v flag.
	Log(args ...any)

	// Logf formats its arguments according to the format, analogous to Printf,
	// and records the text in the error log. A final newline is added if not
	// provided. For tests, the text will be printed only if the test fails or
	// the -test.v flag is set. For benchmarks, the text is always printed to
	// avoid having performance depend on the value of the -test.v flag.
	Logf(format string, args ...any)

	// Name returns the name of the running (sub-) test.
	//
	// The name will include the name of the test along with the names of any
	// nested subtests. If two sibling subtests have the same name, Name will
	// append a suffix to guarantee the returned name is unique.
	Name() string

	// Setenv calls os.Setenv(key, value) and uses Cleanup to restore the
	// environment variable to its original value after the test.
	//
	// This cannot be used in parallel tests.
	Setenv(key, value string)

	// Skip is equivalent to Log followed by SkipNow.
	Skip(args ...any)

	// TempDir returns a temporary directory for the test to use. The directory
	// is automatically removed by Cleanup when the test and all its subtests
	// complete. Each subsequent call to TempDir returns a unique directory;
	// if the directory creation fails, TempDir terminates the test by calling
	// Fatal.
	TempDir() string

	// Context returns a context that is canceled just before Cleanup
	// registered functions are called. Cleanup functions can wait for any
	// resources that shut down on Context.Done before the test or benchmark
	// completes.
	Context() context.Context
}
