// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_ErrReadSeeker(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})

		// --- When ---
		have := ErrReadSeeker(src, 42)

		// --- Then ---
		assert.Same(t, src, have.r)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, ErrRead, have.errRead)
		assert.Same(t, src, have.seek)
	})

	t.Run("does not set reader error when n is negative", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})

		// --- When ---
		have := ErrReadSeeker(src, -1)

		// --- Then ---
		assert.Same(t, src, have.r)
		assert.Equal(t, -1, have.n)
		assert.Equal(t, 0, have.off)
		assert.NoError(t, have.errRead)
		assert.Same(t, src, have.seek)
	})

	t.Run("read error set via option is not overridden", func(t *testing.T) {
		// --- Given ---
		custom := errors.New("my error")
		src := bytes.NewReader([]byte{0, 1, 2, 3})

		// --- When ---
		have := ErrReadSeeker(src, 42, WithReadErr(custom))

		// --- Then ---
		assert.Same(t, src, have.r)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, custom, have.errRead)
		assert.Same(t, src, have.seek)
	})
}

func Test_ErrorReadSeeker_Seek(t *testing.T) {
	t.Run("seek error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		rcs := ErrReadSeeker(src, -1, WithSeekErr(exp))

		// --- When ---
		n, err := rcs.Seek(10, io.SeekStart)

		// --- Then ---
		assert.Same(t, exp, err)
		assert.Equal(t, int64(0), n)
	})

	t.Run("underlying seeker error", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		rcs := ErrReadSeeker(src, -1)

		// --- When ---
		n, err := rcs.Seek(-1, io.SeekStart)

		// --- Then ---
		assert.ErrorEqual(t, "bytes.Reader.Seek: negative position", err)
		assert.Equal(t, int64(0), n)
	})
}
