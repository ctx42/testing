package tstkit

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithWriteErr(t *testing.T) {
	// --- Given ---
	custom := errors.New("my error")
	ew := &ErrorWriter{}

	// --- When ---
	WithWriteErr(custom)(ew)

	// --- Then ---
	assert.Same(t, custom, ew.errWrite)
}

func Test_WithWriteCloseErr(t *testing.T) {
	// --- Given ---
	custom := errors.New("my error")
	ew := &ErrorWriter{}

	// --- When ---
	WithWriteCloseErr(custom)(ew)

	// --- Then ---
	assert.Same(t, custom, ew.errClose)
}

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
		have := ErrWriter(dst, -1, WithWriteErr(custom))

		// --- Then ---
		assert.Same(t, dst, have.w)
		assert.Equal(t, -1, have.n)
		assert.Equal(t, 0, have.off)
		assert.Same(t, custom, have.errWrite)
	})
}

func Test_ErrWriter_Write(t *testing.T) {
	t.Run("no read error when n is negative", func(t *testing.T) {
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
		dst := sadWriter{}
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

func Test_ErrorWriter_Close(t *testing.T) {
	t.Run("no close error", func(t *testing.T) {
		// --- Given ---
		dst := &happyWriter{}
		ew := ErrWriter(dst, -1)

		// --- When ---
		err := ew.Close()

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("underlying closer error", func(t *testing.T) {
		// --- Given ---
		ew := ErrWriter(sadWriter{}, -1)

		// --- When ---
		err := ew.Close()

		// --- Then ---
		assert.Same(t, ErrSadClose, err)
	})

	t.Run("close not defined on reader", func(t *testing.T) {
		// --- Given ---
		dst := &bytes.Buffer{}
		rcs := ErrWriter(dst, -1)

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.ErrorEqual(t, "method Close is not implemented", err)
	})

	t.Run("custom close error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		rcs := ErrWriter(sadWriter{}, -1, WithWriteCloseErr(exp))

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.Same(t, exp, err)
	})
}

var ErrSadWrite = errors.New("sad write error")

type sadWriter struct{}

func (w sadWriter) Write(_ []byte) (int, error) { return 0, ErrSadWrite }
func (w sadWriter) Close() error                { return ErrSadClose }

type happyWriter struct{}

func (w happyWriter) Write(_ []byte) (int, error) { return 0, nil }
func (w happyWriter) Close() error                { return nil }
