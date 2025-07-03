package tstkit

import (
	"errors"
	"io"
)

// WithWriteErr is an [ErrWriter] option setting custom write error.
func WithWriteErr(err error) func(writer *ErrorWriter) {
	return func(ew *ErrorWriter) { ew.errWrite = err }
}

// WithWriteCloseErr is an [ErrWriter] option setting custom close error.
func WithWriteCloseErr(err error) func(writer *ErrorWriter) {
	return func(ew *ErrorWriter) { ew.errClose = err }
}

// ErrWrite is the default write error used by [ErrorWriter].
var ErrWrite = errors.New("write error")

// ErrorWriter is an [io.Writer] that writes up to n bytes to an underlying
// source returning an error. If n is negative, it behaves as a standard
// [io.Writer] with no byte limit.
//
// For writers implementing [io.Closer], use [WithWriteCloseErr] to specify an
// error to return from the Close method when called.
type ErrorWriter struct {
	w        io.Writer // Underlying writer.
	off      int       // Number of written bytes.
	n        int       // At most bytes to write without error.
	errWrite error     // Error to return after writing n bytes.
	errClose error     // Error to return when closing.
}

// ErrWriter wraps writer w and allows control when to return read and close
// errors. It acts like a regular writer if n < 0 and no options were provided.
//
// With [WithWriteErr] option, you can set what error is returned.
//
// With [WithWriteCloseErr] you can set an error to be returned when the Close
// method is called. If the provided [io.Writer] can be cast to [io.Closer],
// the original Close method will always be called, but the original error
// returned from that call will only be returned if the [WithReaderCloseErr]
// was not defined.
func ErrWriter(w io.Writer, n int, opts ...func(*ErrorWriter)) *ErrorWriter {
	ew := &ErrorWriter{
		w: w,
		n: n,
	}
	for _, opt := range opts {
		opt(ew)
	}
	if ew.n >= 0 && ew.errWrite == nil {
		ew.errWrite = ErrWrite
	}
	return ew
}

// Write writes to the underlying buffer and returns error if number of written
// bytes is equal to the predefined limit.
func (ew *ErrorWriter) Write(p []byte) (int, error) {
	// Write no more than n bytes.
	if ew.n >= 0 && ew.off+len(p) > ew.n {
		p = p[:ew.n-ew.off]
	}
	n, err := ew.w.Write(p)
	ew.off += n
	if err != nil {
		return n, err
	}
	if ew.errWrite != nil && ew.off >= ew.n {
		return n, ew.errWrite
	}
	return n, nil
}

func (ew *ErrorWriter) Close() error {
	if c, ok := ew.w.(io.Closer); ok {
		err := c.Close() // The underlying Close method is always called.
		if ew.errClose != nil {
			return ew.errClose
		}
		return err
	}
	return errors.New("method Close is not implemented")
}
