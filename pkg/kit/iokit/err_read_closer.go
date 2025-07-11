// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"io"
)

// ErrorReadCloser implements [io.ReadCloser] by embedding [ErrorReader]
// instance and adding a Close method which behavior may be controlled with
// options. See [ErrReadCloser] constructor function for details.
type ErrorReadCloser struct {
	*ErrorReader
	cls io.Closer
}

// ErrReadCloser wraps the "src" [io.ReadCloser] and controls how many bytes
// can be read from it (n) before it returns an error. If the "n" is negative,
// it behaves like a regular reader. The retuned error may be customized with
// [WithReadErr] option.
//
// By default, the [ErrorReadCloser.Close] method calls the original Close
// method and returns whatever it returned. You may customize what error is
// returned from [ErrorReadCloser.Close] with a [WithCloseErr] option. When a
// [WithCloseErr] option is used, the original Close method is also called, but
// its return value is ignored.
func ErrReadCloser(src io.ReadCloser, n int, opts ...Option) *ErrorReadCloser {
	return &ErrorReadCloser{
		ErrorReader: ErrReader(src, n, opts...),
		cls:         src,
	}
}

func (rc *ErrorReadCloser) Close() error {
	err := rc.cls.Close() // The underlying Close method is always called.
	if rc.errClose != nil {
		return rc.errClose
	}
	return err
}
