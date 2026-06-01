// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

import (
	"io"
)

// ErrorReadSeeker implements [io.ReadSeeker] by embedding an [ErrorReader]
// and adding controllable Seek behavior.
type ErrorReadSeeker struct {
	*ErrorReader
	seek io.Seeker
}

// ErrReadSeeker wraps "src" and allows up to n bytes to be read before
// returning an error. If n < 0 it behaves normally.
//
// Use With*Err options to customize read and seek errors (original Seek
// is still called).
func ErrReadSeeker(src io.ReadSeeker, n int, opts ...Option) *ErrorReadSeeker {
	return &ErrorReadSeeker{
		ErrorReader: ErrReader(src, n, opts...),
		seek:        src,
	}
}

func (rs *ErrorReadSeeker) Seek(offset int64, whence int) (int64, error) {
	n, err := rs.seek.Seek(offset, whence)
	if rs.errSeek != nil {
		return 0, rs.errSeek
	}
	return n, err
}
