// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package kit

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// globLog is a global logger used package-wide.
var globLog = log.New(os.Stderr, "*** KIT ", 0)

// cleanup records a function registered for global post-test cleanup,
// along with the source location where it was registered.
type cleanup struct {
	fn   func() // The function to call during cleanup.
	file string // Basename of the file where AddGlobalCleanup was called.
	line int    // Line number in that file.
}

// cleanupLogFormat is the format string used by RunGlobalCleanups when
// logging which registered cleanup functions are being executed.
const cleanupLogFormat = "running global cleanup function registered in %s:%d"

// cleanups holds the list of functions registered via AddGlobalCleanup.
// It is protected by cleanupMx.
var cleanups []cleanup

// cleanupMx protects access to the cleanups slice.
var cleanupMx sync.Mutex

func init() { cleanups = make([]cleanup, 0, 10) }

// AddGlobalCleanup registers a function that will be invoked by
// [RunGlobalCleanups].
//
// This is a global mechanism primarily intended for performing cleanup
// after all tests have completed (typically from a TestMain function).
//
// WARNING:
//   - This API uses package-level mutable state (protected by a mutex).
//   - All registered cleanups are visible across the entire program.
//   - The package uses an init() function and a global logger for diagnostics.
//   - The registration records the call site (file + line) for logging.
//
// Prefer [testing.TB.Cleanup], [testing.TB.TempDir], or per-test cleanup
// mechanisms whenever possible. Use this global facility only when you
// genuinely need post-test actions that must run after every test in the
// package or module has finished.
//
// Example usage:
//
//	// TestMain is the entry point for running tests in this package.
//	func TestMain(m *testing.M) {
//	   exitCode := m.Run()
//	   kit.RunGlobalCleanups()
//	   os.Exit(exitCode)
//	}
func AddGlobalCleanup(fn func()) {
	cleanupMx.Lock()
	defer cleanupMx.Unlock()
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)
	cleanups = append(cleanups, cleanup{fn: fn, file: file, line: line})
}

// RunGlobalCleanups executes every function previously registered with
// [AddGlobalCleanup], in registration order, and then clears the list.
//
// Each cleanup is logged (file:line) via the package's internal logger
// before it is invoked. This logging cannot currently be disabled.
//
// This function is safe to call multiple times; after the first call the
// list will be empty until new cleanups are registered.
//
// See [AddGlobalCleanup] for important warnings about global state,
// the package init(), and the global logger.
func RunGlobalCleanups() {
	cleanupMx.Lock()
	defer cleanupMx.Unlock()
	for _, cln := range cleanups {
		globLog.Printf(cleanupLogFormat, cln.file, cln.line)
		cln.fn()
	}
	cleanups = cleanups[:0]
}
