// Package iokit provides I/O and buffer related helpers.
package iokit

import (
	"bytes"
	"io"
)

// ReadAllFromStart seeks to the beginning of "rs" and reads it till gets
// [io.EOF] or any other error. Then seek back to the position where "rs" was
// before the call. Panics on error.
func ReadAllFromStart(rs io.ReadSeeker) []byte {
	cur, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	if _, err = rs.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	defer func() { _, _ = rs.Seek(cur, io.SeekStart) }()

	ret := &bytes.Buffer{}
	if _, err = ret.ReadFrom(rs); err != nil {
		panic(err)
	}

	return ret.Bytes()
}

// Offset returns the current offset of the seeker. Panics on error.
func Offset(s io.Seeker) int64 { return Seek(s, 0, io.SeekCurrent) }

// Seek sets the offset for the next Read or Write operation to offset,
// interpreted according to whence. Seek returns the new offset relative to the
// start of the s. Panics on error.
func Seek(s io.Seeker, offset int64, whence int) int64 {
	off, err := s.Seek(offset, whence)
	if err != nil {
		panic(err)
	}
	return off
}
