// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package iokit provides I/O and buffer-related helpers for testing.
//
// It offers two main categories of tools:
//
//   - [Buffer] and its constructors ([NewBuffer], [DryBuffer], [WetBuffer])
//     for controlled, thread-safe buffering with automatic test cleanup checks.
//   - Error-injecting readers and writers ([ErrReader], [ErrWriter], and
//     their Closer/Seeker variants) that let you simulate I/O errors after
//     a configurable number of bytes.
//
// See the package [README] for more detailed usage examples.
//
// All helpers are designed to work well with [tester.T] and the assertion
// packages.
package iokit

import (
	"bytes"
	"io"
)

// ReadAllFromStart seeks to the beginning of "rs", reads until [io.EOF] (or
// another error), then seeks back to the original position. It panics on any
// seek or read error.
//
// See also [Offset] and [Seek].
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

// Offset returns the current offset of the seeker. It panics on error.
//
// See also [Seek] and [ReadAllFromStart].
func Offset(s io.Seeker) int64 { return Seek(s, 0, io.SeekCurrent) }

// Seek sets the offset for the next Read or Write operation, interpreted
// according to "whence". It returns the new offset relative to the start of
// the seeker and panics on error.
//
// See also [Offset] and [ReadAllFromStart].
func Seek(s io.Seeker, offset int64, whence int) int64 {
	off, err := s.Seek(offset, whence)
	if err != nil {
		panic(err)
	}
	return off
}
