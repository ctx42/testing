package iokit

import (
	"errors"
	"io"
)

// ErrRead is the default read error.
var ErrRead = errors.New("read error")

// ErrorReader implements [io.Reader] that reads up to n bytes from an
// underlying reader then returns an error. If n is negative, it behaves as a
// standard reader without returning an error. See [ErrReader] constructor
// function for details.
type ErrorReader struct {
	*Options           // Reader options.
	r        io.Reader // Underlying error.
	n        int       // At most bytes to read without error.
	off      int       // Number of bytes read.
}

// ErrReader wraps the "src" [io.Reader] and controls how many bytes can be
// read from it (n) before it returns an error. If the "n" is negative, it
// behaves like a regular reader. With [WithReadErr] option, you can customize
// the returned error.
func ErrReader(src io.Reader, n int, opts ...Option) *ErrorReader {
	r := &ErrorReader{
		Options: defaultOptions(),
		r:       src,
		n:       n,
	}
	for _, opt := range opts {
		opt(r.Options)
	}
	if r.n < 0 {
		r.errRead = nil
	}
	return r
}

func (r *ErrorReader) Read(p []byte) (int, error) {
	// Read up to the limit - no more.
	if r.n >= 0 && r.off+len(p) > r.n {
		p = p[:r.n-r.off]
	}
	n, err := r.r.Read(p)
	r.off += n
	if err != nil {
		return n, err
	}
	if r.errRead != nil && r.off >= r.n {
		return n, r.errRead
	}
	return n, nil
}
