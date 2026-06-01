// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package kit

import (
	"bytes"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_AddGlobalCleanup(t *testing.T) {
	t.Setenv("___", "___")
	origLog := globLog
	buf := &bytes.Buffer{}
	globLog = log.New(buf, "", 0)
	t.Cleanup(func() { globLog = origLog })

	resetCleanupsForTest()

	t.Run("add cleanup function", func(t *testing.T) {
		// --- Given ---
		var called bool
		fn := func() { called = true }

		// --- When ---
		_, file, line, _ := runtime.Caller(0)
		AddGlobalCleanup(fn)
		file = filepath.Base(file)
		expectedLine := line + 1

		// --- Then ---
		assert.Len(t, 1, cleanups)
		assert.Same(t, fn, cleanups[0].fn)
		assert.Equal(t, file, cleanups[0].file)
		assert.Equal(t, expectedLine, cleanups[0].line)
		assert.False(t, called)
		assert.Equal(t, "", buf.String())
	})
}

func Test_RunGlobalCleanups(t *testing.T) {
	t.Setenv("___", "___")
	origLog := globLog
	buf := &bytes.Buffer{}
	globLog = log.New(buf, "", 0)
	t.Cleanup(func() { globLog = origLog })

	resetCleanupsForTest()

	t.Run("add cleanup function", func(t *testing.T) {
		// --- Given ---
		var fn0, fn1 bool
		// We manually populate for this test to control the locations
		cleanups = []cleanup{
			{fn: func() { fn0 = true }, file: "fn1.go", line: 42},
			{fn: func() { fn1 = true }, file: "fn2.go", line: 44},
		}

		// --- When ---
		RunGlobalCleanups()

		// --- Then ---
		assert.Len(t, 0, cleanups)
		assert.True(t, fn0)
		assert.True(t, fn1)
		want := "" +
			"running global cleanup function registered in fn1.go:42\n" +
			"running global cleanup function registered in fn2.go:44\n"
		assert.Equal(t, want, buf.String())
	})
}

// resetCleanupsForTest clears the global cleanup state for test isolation.
func resetCleanupsForTest() {
	cleanupMx.Lock()
	cleanups = cleanups[:0]
	cleanupMx.Unlock()
}
