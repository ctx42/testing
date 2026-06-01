// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

import (
	"io"
)

// ErrorReadSeekCloser implements [io.ReadSeekCloser] by embedding an
// [ErrorReadSeeker] and adding controllable Close behavior.
type ErrorReadSeekCloser struct {
	*ErrorReadSeeker
	cls io.Closer
}

// ErrReadSeekCloser wraps the "src" [io.ReadSeekCloser] and controls how many
// ErrReadSeekCloser wraps "src" and allows up to n bytes to be read before
// returning an error. If n < 0 it behaves normally.
//
// Use With*Err options to customize read, seek, and close errors (original
// methods are still called).
func ErrReadSeekCloser(
	src io.ReadSeekCloser,
	n int,
	opts ...Option,
) *ErrorReadSeekCloser {
	return &ErrorReadSeekCloser{
		ErrorReadSeeker: ErrReadSeeker(src, n, opts...),
		cls:             src,
	}
}

func (rc *ErrorReadSeekCloser) Close() error {
	err := rc.cls.Close() // The underlying Close method is always called.
	if rc.errClose != nil {
		return rc.errClose
	}
	return err
}
