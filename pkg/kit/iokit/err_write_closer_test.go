// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_ErrWriteCloser(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		pth := filepath.Join(t.TempDir(), "file.txt")
		dst := must.Value(os.Create(pth))
		defer t.Cleanup(func() { _ = dst.Close() })

		// --- When ---
		have := ErrWriteCloser(dst, 42)

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, ErrWrite, have.errWrite)
		assert.Same(t, dst, have.cls)
	})

	t.Run("does not set writer error when n is negative", func(t *testing.T) {
		// --- Given ---
		pth := filepath.Join(t.TempDir(), "file.txt")
		dst := must.Value(os.Create(pth))
		defer t.Cleanup(func() { _ = dst.Close() })

		// --- When ---
		have := ErrWriteCloser(dst, -1)

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, -1, have.n)
		assert.Equal(t, 0, have.off)
		assert.NoError(t, have.errWrite)
		assert.Same(t, dst, have.cls)
	})

	t.Run("write error set via option is not overridden", func(t *testing.T) {
		// --- Given ---
		custom := errors.New("my error")
		pth := filepath.Join(t.TempDir(), "file.txt")
		dst := must.Value(os.Create(pth))
		defer t.Cleanup(func() { _ = dst.Close() })

		// --- When ---
		have := ErrWriteCloser(dst, 42, WithWriteErr(custom))

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.Same(t, custom, have.errWrite)
		assert.Same(t, dst, have.cls)
	})
}

func Test_ErrorWriteCloser_Close(t *testing.T) {
	t.Run("no close error", func(t *testing.T) {
		// --- Given ---
		dst := &blissfulWriter{}
		ew := ErrWriteCloser(dst, -1)

		// --- When ---
		err := ew.Close()

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("underlying closer error", func(t *testing.T) {
		// --- Given ---
		ew := ErrWriteCloser(sadStruct{}, -1)

		// --- When ---
		err := ew.Close()

		// --- Then ---
		assert.Same(t, ErrSadClose, err)
	})

	t.Run("custom close error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		rcs := ErrWriteCloser(sadStruct{}, -1, WithCloseErr(exp))

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.Same(t, exp, err)
	})
}
