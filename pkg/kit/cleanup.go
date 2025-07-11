// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

// cleanup represents cleanup function and where it was created.
type cleanup struct {
	fn   func() // Cleanup function.
	file string // Filepath the [AddGlobalCleanup] was called.
	line int    // Line number the [AddGlobalCleanup] was called.
}

// cleanups represent a slice of cleanup functions to call.
var cleanups []cleanup

// cleanupMx guards cleanup slice.
var cleanupMx sync.Mutex

func init() { cleanups = make([]cleanup, 0, 10) }

// AddGlobalCleanup adds a global cleanup function.
//
// Example usage:
//
//	// TestMain is the entry point for running tests in this package.
//	func TestMain(m *testing.M) {
//	   // Run all tests and capture the exit code.
//	   exitCode := m.Run()
//
//	   // Cleanup code (runs after all tests).
//	   kit.RunGlobalCleanups()
//
//	   // Exit with the test result code.
//	   os.Exit(exitCode)
//	}
func AddGlobalCleanup(fn func()) {
	cleanupMx.Lock()
	defer cleanupMx.Unlock()
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)
	cleanups = append(cleanups, cleanup{fn: fn, file: file, line: line})
}

// RunGlobalCleanups runs global cleanup functions. Guarantees cleanup
// functions are run only once.
func RunGlobalCleanups() {
	cleanupMx.Lock()
	defer cleanupMx.Unlock()
	for _, cln := range cleanups {
		format := "running global cleanup function registered in %s:%d"
		globLog.Printf(format, cln.file, cln.line)
		cln.fn()
	}
	cleanups = cleanups[:0]
}
