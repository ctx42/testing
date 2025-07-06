package iokit

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_ErrReader(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		rdr := &bytes.Buffer{}

		// --- When ---
		have := ErrReader(rdr, 42)

		// --- Then ---
		assert.Same(t, rdr, have.r)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, ErrRead, have.errRead)
	})

	t.Run("does not set reader error when n is negative", func(t *testing.T) {
		// --- Given ---
		rdr := &bytes.Buffer{}

		// --- When ---
		have := ErrReader(rdr, -1)

		// --- Then ---
		assert.Same(t, rdr, have.r)
		assert.Equal(t, -1, have.n)
		assert.Equal(t, 0, have.off)
		assert.NoError(t, have.errRead)
	})

	t.Run("read error set via option is not overridden", func(t *testing.T) {
		// --- Given ---
		custom := errors.New("my error")
		rdr := &bytes.Buffer{}

		// --- When ---
		have := ErrReader(rdr, 42, WithReadErr(custom))

		// --- Then ---
		assert.Same(t, rdr, have.r)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, custom, have.errRead)
	})
}

func Test_ErrorReader_Read(t *testing.T) {
	t.Run("no read error when n is negative", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		dst := make([]byte, 5)
		rcs := ErrReader(src, -1)

		// --- When ---
		n, err := rcs.Read(dst)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, []byte{0, 1, 2, 3, 0}, dst)

		n, err = rcs.Read(dst)
		assert.Equal(t, 0, n)
		assert.ErrorIs(t, io.EOF, err)
	})

	t.Run("custom error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		dst := make([]byte, 3)
		rcs := ErrReader(src, 3, WithReadErr(exp))

		// --- When ---
		n, err := rcs.Read(dst)

		// --- Then ---
		assert.Same(t, exp, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, []byte{0, 1, 2}, dst)
	})

	t.Run("read error on last read", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		dst := make([]byte, 2)
		rcs := ErrReader(src, 3)

		// --- When ---
		n, err := rcs.Read(dst)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, []byte{0, 1}, dst)
		assert.Equal(t, 2, n)

		n, err = rcs.Read(dst)
		assert.Same(t, ErrRead, err)
		assert.Equal(t, 1, n)
		assert.Equal(t, []byte{2, 1}, dst)
	})

	t.Run("read up to n", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		dst := make([]byte, 3)
		rcs := ErrReader(src, 3)

		// --- When ---
		n, err := rcs.Read(dst)

		// --- Then ---
		assert.ErrorIs(t, ErrRead, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, []byte{0, 1, 2}, dst)
	})
}
