// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

import (
	"io"
)

// ErrorWriteCloser implements [io.WriteCloser] by embedding an [ErrorWriter]
// and adding controllable Close behavior.
//
// See [ErrWriteCloser] for the constructor.
type ErrorWriteCloser struct {
	*ErrorWriter
	cls io.Closer
}

// ErrWriteCloser wraps "dst" and allows up to n bytes to be written before
// returning an error. If n < 0 it behaves normally.
//
// Use With*Err options to customize write and close errors (original Close
// is still called).
func ErrWriteCloser(
	dst io.WriteCloser,
	n int,
	opts ...Option,
) *ErrorWriteCloser {
	return &ErrorWriteCloser{
		ErrorWriter: ErrWriter(dst, n, opts...),
		cls:         dst,
	}
}

// Close implements [io.Closer]. The underlying Close is always invoked.
// If a custom close error was set via [WithCloseErr], that error is returned
// instead of the underlying result.
func (wc *ErrorWriteCloser) Close() error {
	err := wc.cls.Close() // The underlying Close method is always called.
	if wc.errClose != nil {
		return wc.errClose
	}
	return err
}
