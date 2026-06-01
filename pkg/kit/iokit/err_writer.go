// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

import (
	"errors"
	"io"
)

// ErrWrite is the default write error.
var ErrWrite = errors.New("write error")

// ErrorWriter implements [io.Writer] that writes up to a configured number of
// bytes to an underlying writer, then returns the configured error.
//
// See [ErrWriter] for the constructor and the With*Err options for customization.
type ErrorWriter struct {
	*Options
	w   io.Writer
	n   int
	off int
}

// ErrWriter wraps "dst" and allows up to n bytes to be written before returning
// an error. If n < 0 it behaves like a normal writer.
//
// Use the With*Err options to customize the error.
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

// Write implements [io.Writer]. It writes to the underlying writer up to the
// configured limit, then returns the configured error (or the underlying
// writer's error).
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
