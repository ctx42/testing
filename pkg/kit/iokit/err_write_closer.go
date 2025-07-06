package iokit

import (
	"io"
)

// ErrorWriteCloser implements [io.WriteCloser] by embedding [ErrorWriter]
// instance and adding a Close method which behavior you may control with options.
// See [ErrWriteCloser] constructor function for details.
type ErrorWriteCloser struct {
	*ErrorWriter
	cls io.Closer
}

// ErrWriteCloser wraps the "dst" [io.WriteCloser] and controls how many bytes
// can be written to it (n) before it returns an error. If the "n" is negative,
// it behaves like a regular writer. With [WithWriteErr] option, you can
// customize the returned error.
//
// By default, the [ErrorWriteCloser.Close] method calls the original Close
// method and returns whatever it returned. You may customize what error is
// returned from [ErrorWriteCloser.Close] with a [WithCloseErr] option. When a
// [WithCloseErr] option is used, the original Close method is also called, but
// its return value is ignored.
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

func (wc *ErrorWriteCloser) Close() error {
	err := wc.cls.Close() // The underlying Close method is always called.
	if wc.errClose != nil {
		return wc.errClose
	}
	return err
}
