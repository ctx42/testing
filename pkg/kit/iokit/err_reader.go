// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

import (
	"errors"
	"io"
)

// ErrRead is the default error returned by error readers when no custom error
// is provided via [WithReadErr].
var ErrRead = errors.New("read error")

// ErrorReader implements [io.Reader] that reads up to a configured number of
// bytes from an underlying reader, then returns the configured error.
//
// See [ErrReader] for the constructor and the With*Err options for customization.
type ErrorReader struct {
	*Options
	r   io.Reader
	n   int
	off int
}

// ErrReader wraps "src" and allows up to n bytes to be read before returning
// an error. If n < 0 it behaves like a normal reader (no error injection).
//
// Use the With*Err options to customize the error returned.
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

// Read implements [io.Reader]. It reads from the underlying reader up to the
// configured limit, then returns the configured error (or the underlying
// reader's error).
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
