package iokit

import (
	"io"
)

// ErrorReadSeeker implements [io.ReadSeeker] by embedding [ErrorReader]
// instance and adding Seek method which behavior may be controlled with
// options. See [ErrReadSeeker] constructor function for details.
type ErrorReadSeeker struct {
	*ErrorReader
	seek io.Seeker
}

// ErrReadSeeker wraps the "src" [io.ReadSeeker] and controls how many bytes
// can be read from it (n) before it returns an error. If "n" is negative,
// it behaves like a regular reader. With [WithReadErr] option, you can
// customize the returned error.
//
// By default, the [ErrorReadSeeker.Seek] method calls the original Seek method
// and returns whatever it returned. You may customize the returned error
// from [ErrorReadSeeker.Seek] with a [WithSeekErr] option. When a [WithSeekErr]
// option is used, the original Seek method is also called, but its return
// value is ignored and the one provided with the [WithSeekErr] option is used.
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
