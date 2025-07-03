package tstkit

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithReadErr(t *testing.T) {
	// --- Given ---
	custom := errors.New("my error")
	rcs := &ErrorReader{}

	// --- When --
	WithReadErr(custom)(rcs)

	// --- Then ---
	assert.Same(t, custom, rcs.errRead)
}

func Test_WithSeekErr(t *testing.T) {
	// --- Given ---
	me := errors.New("my error")
	rcs := &ErrorReader{}

	// --- When --
	WithSeekErr(me)(rcs)

	// --- Then ---
	assert.Same(t, me, rcs.errSeek)
}

func Test_WithReaderCloseErr(t *testing.T) {
	// --- Given ---
	me := errors.New("my error")
	rcs := &ErrorReader{}

	// --- When --
	WithReaderCloseErr(me)(rcs)

	// --- Then ---
	assert.Same(t, me, rcs.errClose)
}

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
		assert.NoError(t, have.errClose)
		assert.NoError(t, have.errSeek)
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
		assert.NoError(t, have.errClose)
		assert.NoError(t, have.errSeek)
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
		assert.NoError(t, have.errClose)
		assert.NoError(t, have.errSeek)
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

func Test_ErrorReader_Seek(t *testing.T) {
	t.Run("seek error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		rcs := ErrReader(src, -1, WithSeekErr(exp))

		// --- When ---
		n, err := rcs.Seek(10, io.SeekStart)

		// --- Then ---
		assert.Same(t, exp, err)
		assert.Equal(t, int64(0), n)
	})

	t.Run("underlying seeker error", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		rcs := ErrReader(src, -1)

		// --- When ---
		n, err := rcs.Seek(-1, io.SeekStart)

		// --- Then ---
		assert.ErrorEqual(t, "bytes.Reader.Seek: negative position", err)
		assert.Equal(t, int64(0), n)
	})

	t.Run("seek not defined on reader", func(t *testing.T) {
		// --- Given ---
		src := sadReader{}
		rcs := ErrReader(src, -1)

		// --- When ---
		n, err := rcs.Seek(-1, io.SeekStart)

		// --- Then ---
		assert.ErrorEqual(t, "method Seek is not implemented", err)
		assert.Equal(t, int64(0), n)
	})
}

func Test_ErrorReader_Close(t *testing.T) {
	t.Run("no close error", func(t *testing.T) {
		// --- Given ---
		src := io.NopCloser(bytes.NewReader([]byte{0, 1, 2, 3}))
		rcs := ErrReader(src, -1)

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("underlying closer error", func(t *testing.T) {
		// --- Given ---
		rcs := ErrReader(sadReader{}, -1)

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.Same(t, ErrSadClose, err)
	})

	t.Run("close not defined on reader", func(t *testing.T) {
		// --- Given ---
		src := bytes.NewReader([]byte{0, 1, 2, 3})
		rcs := ErrReader(src, -1)

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.ErrorEqual(t, "method Close is not implemented", err)
	})

	t.Run("custom close error", func(t *testing.T) {
		// --- Given ---
		exp := errors.New("test message")
		rcs := ErrReader(sadReader{}, -1, WithReaderCloseErr(exp))

		// --- When ---
		err := rcs.Close()

		// --- Then ---
		assert.Same(t, exp, err)
	})
}

var (
	ErrSadRead  = errors.New("sad read error")
	ErrSadClose = errors.New("sad close error")
)

type sadReader struct{}

func (r sadReader) Read(_ []byte) (int, error) { return 0, ErrSadRead }
func (r sadReader) Close() error               { return ErrSadClose }
