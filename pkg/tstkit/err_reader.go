package tstkit

import (
	"errors"
	"io"
)

// WithReadErr is an [ErrReader] option setting custom read error.
func WithReadErr(err error) func(*ErrorReader) {
	return func(rcs *ErrorReader) { rcs.errRead = err }
}

// WithSeekErr is an [ErrReader] option setting custom seek error.
func WithSeekErr(err error) func(*ErrorReader) {
	return func(rcs *ErrorReader) { rcs.errSeek = err }
}

// WithReaderCloseErr is an [ErrReader] option setting close error.
func WithReaderCloseErr(err error) func(*ErrorReader) {
	return func(rcs *ErrorReader) { rcs.errClose = err }
}

// ErrRead is the default error used by [ErrorReader].
var ErrRead = errors.New("read error")

// ErrorReader is an [io.Reader] that reads up to n bytes from an underlying
// source returning an error. If n is negative, it behaves as a standard
// [io.Reader] with no byte limit.
//
// For readers implementing [io.Closer], use [WithReaderCloseErr] to specify an error
// to return from the Close method when called.
//
// For readers implementing [io.Seeker], use [WithSeekErr] to specify an error
// to return from the Seek method when called.
type ErrorReader struct {
	r        io.Reader // Underlying error.
	n        int       // At most bytes to read without error.
	off      int       // Number of bytes read.
	errRead  error     // Error to return after reading n bytes.
	errClose error     // Error to return when closing.
	errSeek  error     // Error to return when seeking.
}

// ErrReader wraps reader r and allows control when to return read, seek, and
// close errors. It acts like a regular reader if n < 0 and no options were
// provided.
//
// With [WithReadErr] option, you can set what error is returned.
//
// With [WithSeekErr] option, you can set an error that should be returned when
// the Seek method is called. The underlying Seek method will be called if the
// seek error is not provided and provided [io.Reader] can be cast to
// [io.Seeker].
//
// With [WithReaderCloseErr] you can set an error to be returned when the Close
// method is called. If the provided [io.Reader] can be cast to [io.Closer],
// the original Close method will always be called, but the original error
// returned from that call will only be returned if the [WithReaderCloseErr] was not
// defined.
func ErrReader(r io.Reader, n int, opts ...func(*ErrorReader)) *ErrorReader {
	rcs := &ErrorReader{
		r: r,
		n: n,
	}
	for _, opt := range opts {
		opt(rcs)
	}
	if rcs.n >= 0 && rcs.errRead == nil {
		rcs.errRead = ErrRead
	}
	return rcs
}

// Read implements io.Reader which returns an error after reading n bytes.
func (rcs *ErrorReader) Read(p []byte) (int, error) {
	// Read up to the limit - no more.
	if rcs.n >= 0 && rcs.off+len(p) > rcs.n {
		p = p[:rcs.n-rcs.off]
	}
	n, err := rcs.r.Read(p)
	rcs.off += n
	if err != nil {
		return n, err
	}
	if rcs.errRead != nil && rcs.off >= rcs.n {
		return n, rcs.errRead
	}
	return n, nil
}

func (rcs *ErrorReader) Seek(offset int64, whence int) (int64, error) {
	if s, ok := rcs.r.(io.Seeker); ok {
		if rcs.errSeek != nil {
			return 0, rcs.errSeek
		}
		return s.Seek(offset, whence)
	}
	return 0, errors.New("method Seek is not implemented")
}

func (rcs *ErrorReader) Close() error {
	if c, ok := rcs.r.(io.Closer); ok {
		err := c.Close() // The underlying Close method is always called.
		if rcs.errClose != nil {
			return rcs.errClose
		}
		return err
	}
	return errors.New("method Close is not implemented")
}
