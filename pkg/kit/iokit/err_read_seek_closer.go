// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"io"
)

// ErrorReadSeekCloser implements [io.ReadSeekCloser] by embedding
// [ErrorReadSeeker] instance and adding Close method which behavior you may
// control with options. See [ErrReadSeekCloser] constructor function for
// details.
type ErrorReadSeekCloser struct {
	*ErrorReadSeeker
	cls io.Closer
}

// ErrReadSeekCloser wraps the "src" [io.ReadSeekCloser] and controls how many
// bytes can be read from it (n) before it returns an error. If the "n" is
// negative, it behaves like a regular reader. With [WithReadErr] option, you
// can customize the returned error.
//
// By default, the [ErrorReadSeekCloser.Seek] method calls the original Seek
// method and returns whatever it returned. You may customize what error is
// returned from [ErrorReadSeekCloser.Seek] with a [WithSeekErr] option. When a
// [WithSeekErr] option is used, the original Seek method is also called, but
// its return value is ignored and the one provided with the [WithSeekErr]
// option is used.
//
// By default, the [ErrorReadSeekCloser.Close] method calls the original Close
// method and returns whatever it returned. You may customize what error is
// returned from [ErrorReadSeekCloser.Close] with a [WithCloseErr] option. When
// a [WithCloseErr] option is used, the original Close method is also called,
// but its return value is ignored.
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
