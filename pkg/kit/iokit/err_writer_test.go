package iokit

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_ErrWriter(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}

		// --- When ---
		have := ErrWriter(dst, 42)

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.ErrorIs(t, ErrWrite, have.errWrite)
	})

	t.Run("does not set writer error when n is negative", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}

		// --- When ---
		have := ErrWriter(dst, -1)

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, -1, have.n)
		assert.Equal(t, 0, have.off)
		assert.NoError(t, have.errWrite)
	})

	t.Run("write error set via option is not overridden", func(t *testing.T) {
		// --- Given ---
		custom := errors.New("my error")
		dst := &bytes.Buffer{}

		// --- When ---
		have := ErrWriter(dst, 42, WithWriteErr(custom))

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, 42, have.n)
		assert.Equal(t, 0, have.off)
		assert.Same(t, custom, have.errWrite)
	})
}

func Test_ErrWriter_Write(t *testing.T) {
	t.Run("no error when n is negative", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}
		ew := ErrWriter(dst, -1)

		// --- When ---
		n, err := ew.Write([]byte{0, 1, 2, 3})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, []byte{0, 1, 2, 3}, dst.Bytes())
	})

	t.Run("no error when writing less than the limit", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}
		ew := ErrWriter(dst, 3)

		// --- When ---
		n, err := ew.Write([]byte{0, 1})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.Equal(t, []byte{0, 1}, dst.Bytes())
	})

	t.Run("underlying writer error", func(t *testing.T) {
		// --- Given ---
		dst := sadStruct{}
		ew := ErrWriter(dst, 42)

		// --- When ---
		n, err := ew.Write([]byte{0, 1})

		// --- Then ---
		assert.ErrorIs(t, ErrSadWrite, err)
		assert.Equal(t, 0, n)
	})

	t.Run("error when writing more than the limit", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}
		ew := ErrWriter(dst, 3)

		// --- When ---
		n, err := ew.Write([]byte{0, 1, 2, 3, 4})

		// --- Then ---
		assert.ErrorIs(t, ErrWrite, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, []byte{0, 1, 2}, dst.Bytes())
	})

	t.Run("custom error", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}
		custom := errors.New("my error")

		// --- When ---
		ew := ErrWriter(dst, 3, WithWriteErr(custom))
		n, err := ew.Write([]byte{0, 1, 2})

		// --- Then ---
		assert.ErrorIs(t, custom, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, []byte{0, 1, 2}, dst.Bytes())
	})
}
