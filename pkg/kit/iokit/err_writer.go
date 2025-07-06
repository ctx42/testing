package iokit

import (
	"errors"
	"io"
)

// ErrWrite is the default write error.
var ErrWrite = errors.New("write error")

// ErrorWriter implements [io.Writer] that writes up to n bytes from an
// underlying writer then returns an error. If n is negative, it behaves as a
// standard writer without returning an error. See [ErrWriter] constructor
// function for details.
type ErrorWriter struct {
	*Options           // Writer options.
	w        io.Writer // Underlying writer.
	n        int       // At most bytes to write without error.
	off      int       // Number of bytes written.
}

// ErrWriter wraps the "dst" [io.Writer] and controls how many bytes can be
// written to it (n) before it returns an error. If the "n" is negative, it
// behaves like a regular writer. With [WithWriteErr] option, you can customize
// the returned error.
func ErrWriter(dst io.Writer, n int, opts ...Option) *ErrorWriter {
	ew := &ErrorWriter{
		Options: defaultOptions(),
		w:       dst,
		n:       n,
	}
	for _, opt := range opts {
		opt(ew.Options)
	}
	if ew.n < 0 {
		ew.errWrite = nil
	}
	return ew
}

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
