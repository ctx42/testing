// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package kit

import (
	"bytes"
	"log"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_AddGlobalCleanup(t *testing.T) {
	t.Setenv("___", "___")
	origLog := globLog
	buf := &bytes.Buffer{}
	globLog = log.New(buf, "", 0)
	t.Cleanup(func() { globLog = origLog })

	t.Run("add cleanup function", func(t *testing.T) {
		// --- Given ---
		var called bool
		fn := func() { called = true }

		// --- When ---
		AddGlobalCleanup(fn)

		// --- Then ---
		assert.Len(t, 1, cleanups)
		assert.Same(t, fn, cleanups[0].fn)
		assert.Equal(t, 27, cleanups[0].line)
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

	t.Run("add cleanup function", func(t *testing.T) {
		// --- Given ---
		var fn0, fn1 bool
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
