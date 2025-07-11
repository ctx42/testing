// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"errors"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_ErrReadCloser(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		src := must.Value(os.Open("testdata/file.txt"))
		t.Cleanup(func() { _ = src.Close() })

		// --- When ---
		have := ErrReadCloser(src, 42)

		// --- Then ---
		assert.Same(t, src, have.r)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, ErrRead, have.errRead)
		assert.Same(t, src, have.cls)
	})

	t.Run("does not set reader error when n is negative", func(t *testing.T) {
		// --- Given ---
		src := must.Value(os.Open("testdata/file.txt"))
		t.Cleanup(func() { _ = src.Close() })

		// --- When ---
		have := ErrReadCloser(src, -1)

		// --- Then ---
		assert.Same(t, src, have.r)
		assert.Equal(t, -1, have.n)
		assert.Equal(t, 0, have.off)
		assert.NoError(t, have.errRead)
		assert.Same(t, src, have.cls)
	})

	t.Run("read error set via option is not overridden", func(t *testing.T) {
		// --- Given ---
		custom := errors.New("my error")
		src := must.Value(os.Open("testdata/file.txt"))
		t.Cleanup(func() { _ = src.Close() })

		// --- When ---
		have := ErrReadCloser(src, 42, WithReadErr(custom))

		// --- Then ---
		assert.Same(t, src, have.r)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, custom, have.errRead)
		assert.Same(t, src, have.cls)
	})
}

func Test_ErrorReadCloser_Close(t *testing.T) {
	t.Run("no close error", func(t *testing.T) {
		// --- Given ---
		src := must.Value(os.Open("testdata/file.txt"))
		t.Cleanup(func() { _ = src.Close() })
		rcs := ErrReadCloser(src, -1)

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("underlying closer error", func(t *testing.T) {
		// --- Given ---
		rcs := ErrReadCloser(sadStruct{}, -1)

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.Same(t, ErrSadClose, err)
	})

	t.Run("custom close error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		rcs := ErrReadCloser(sadStruct{}, -1, WithCloseErr(exp))

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.Same(t, exp, err)
	})
}
