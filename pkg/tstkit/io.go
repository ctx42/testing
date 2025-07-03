package tstkit

import (
	"bytes"
	"io"
)

// ReadAllFromStart seeks to the beginning of "rs" and reads it till gets
// [io.EOF] or any other error. Then seeks back to the position where "rs" was
// before the call. Panics on error.
func ReadAllFromStart(rs io.ReadSeeker) []byte {
	cur, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
		return nil
	}

	if _, err = rs.Seek(0, io.SeekStart); err != nil {
		panic(err)
		return nil
	}

	defer func() { _, _ = rs.Seek(cur, io.SeekStart) }()

	ret := &bytes.Buffer{}
	if _, err = ret.ReadFrom(rs); err != nil {
		panic(err)
		return nil
	}

	return ret.Bytes()
}
