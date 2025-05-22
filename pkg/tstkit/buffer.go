// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package tstkit

import (
	"bytes"
	"sync"

	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"
)

// Buffer kinds define the behavior of a [Buffer] during test cleanup.
const (
	BufferDry   = "dry"     // BufferDry enforces no data is written.
	BufferWet   = "wet"     // BufferWet enforces data is written and examined.
	BuffDefault = "default" // BuffDefault applies no cleanup checks.
)

// Buffer is a thread-safe buffer that wraps [bytes.Buffer]. It tracks write
// and read operations and supports test cleanup checks based on its kind.
type Buffer struct {
	name string        // Buffer name for identification in the test log.
	kind string        // Buffer kind: [BufferDry], [BufferWet], [BuffDefault].
	buf  *bytes.Buffer // Underlying buffer for data storage.
	mx   sync.Mutex    // Ensures thread-safe access to the buffer.
	wc   int           // Count of write operations.
	rc   int           // Count of read operations via [Buffer.String].

	// The test must examine the written content by calling [Buffer.String]
	// method. By default, it is set to true. You may change this behavior by
	// calling [Buffer.SkipExamine].
	examine bool
}

// NewBuffer creates a new thread-safe [Buffer] with the [BuffDefault] kind. An
// optional name can be provided. If no name is provided, [Buffer.Name] returns
// an empty string.
//
// Example:
//
//	buf := NewBuffer("my-buffer") // Named buffer.
//	buf := NewBuffer()            // Unnamed buffer.
func NewBuffer(names ...string) *Buffer {
	tsb := &Buffer{
		kind:    BuffDefault,
		buf:     &bytes.Buffer{},
		mx:      sync.Mutex{},
		examine: true,
	}
	if len(names) > 0 {
		tsb.name = names[0]
	}
	return tsb
}

// Name returns the buffer's name or an empty string if no name was provided
// during creation.
func (buf *Buffer) Name() string { return buf.name }

// Kind returns the buffer's kind.
func (buf *Buffer) Kind() string { return buf.kind }

// SkipExamine disables the cleanup requirement for the test to examine the
// buffer. This is useful when the buffer is [WetBuffer] but we do not want to
// examine what was written to it. Implements fluent interface.
func (buf *Buffer) SkipExamine() *Buffer {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.examine = false
	return buf
}

// Write writes the byte slice p to the buffer, incrementing the write counter.
// It is thread-safe and implements the [io.Writer] interface.
func (buf *Buffer) Write(p []byte) (n int, err error) {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.wc++
	return buf.buf.Write(p)
}

// WriteString writes the string s to the buffer, incrementing the write
// counter. It is thread-safe and implements the [io.StringWriter] interface.
func (buf *Buffer) WriteString(s string) (n int, err error) {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.wc++
	return buf.buf.WriteString(s)
}

// MustWriteString writes the string s to the buffer, incrementing the write
// counter. It panics if the write operation fails. This is useful in tests
// where write failures are unexpected and should halt execution.
func (buf *Buffer) MustWriteString(s string) int {
	n, _ := buf.WriteString(s) // Panics when out-of-memory.
	return n
}

// String returns the current contents of the buffer as a string, incrementing
// the read counter. It is thread-safe and implements the [fmt.Stringer]
// interface and is intended for inspecting buffer data during tests.
func (buf *Buffer) String() string {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	return buf.string(true)
}

// string returns the buffer's contents as a string. If "inc" is true, it
// increments the read counter. This method assumes the caller holds the lock.
func (buf *Buffer) string(inc bool) string {
	if inc {
		buf.rc++
	}
	return buf.buf.String()
}

// Reset clears the buffer's contents and resets the write and read counters.
// It is thread-safe and prepares the buffer for reuse in tests.
func (buf *Buffer) Reset() {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.wc = 0
	buf.rc = 0
	buf.buf.Reset()
}

// DryBuffer creates a thread-safe [Buffer] with the dry kind [BufferDry],
// which fails the test during cleanup if any data was written to it. An
// optional name can be provided for use in test log output. The provided
// [tester.T] is used to register cleanup checks and report failures.
//
// Example:
//
//	buf := DryBuffer(t, "dry-buffer")
//	buf.WriteString("data") // Will fail test during cleanup.
func DryBuffer(t tester.T, names ...string) *Buffer {
	t.Helper()
	buf := NewBuffer(names...)
	buf.kind = BufferDry
	t.Cleanup(func() {
		t.Helper()
		buf.mx.Lock()
		defer buf.mx.Unlock()
		if out := buf.string(false); out != "" {
			msg := notice.New("expected buffer to be empty").
				Want("%s", dump.ValEmpty).
				Have("%s", out)
			if buf.name != "" {
				_ = msg.Prepend("name", "%s", buf.name)
			}
			t.Error(msg)
		}
	})
	return buf
}

// WetBuffer creates a thread-safe [Buffer] with the wet kind [BufferWet],
// which fails the test during cleanup if no data was written or if the
// contents were not read via [Buffer.String]. An optional name can be provided
// for use in test log output. The provided [tester.T] is used to register
// cleanup checks and report failures.
//
// Example:
//
//	buf := WetBuffer(t, "wet-buffer")
//	buf.WriteString("data")
//	// Must call buf.SkipExamine() or buf.String() to avoid test failure.
func WetBuffer(t tester.T, names ...string) *Buffer {
	t.Helper()
	buf := NewBuffer(names...)
	buf.kind = BufferWet
	t.Cleanup(func() {
		t.Helper()
		buf.mx.Lock()
		defer buf.mx.Unlock()
		out := buf.string(false)
		if out == "" {
			msg := notice.New("expected buffer not to be empty")
			if buf.name != "" {
				_ = msg.Append("name", "%s", buf.name)
			}
			t.Error(msg)
			return
		}
		if !buf.examine {
			return
		}
		if buf.rc == 0 {
			msg := notice.New("expected buffer contents to be examined")
			if buf.name != "" {
				_ = msg.Append("name", "%s", buf.name)
			}
			t.Error(msg)
		}
	})
	return buf
}
