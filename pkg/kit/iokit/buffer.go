// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

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

// Buffer is a thread-safe wrapper around [bytes.Buffer] that tracks write and
// read operations and supports automatic test cleanup verification.
//
// See [NewBuffer], [DryBuffer], and [WetBuffer] for the recommended
// constructors. The type works well with [tester.T] for registering cleanup
// checks.
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

// NewBuffer creates a new thread-safe [Buffer] with the [BuffDefault] kind
// (no automatic cleanup checks). An optional name can be provided.
//
// Prefer [DryBuffer] or [WetBuffer] for most test scenarios.
//
// See the package [README] for usage examples.
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

// Name returns the buffer's name (or empty if none was provided at creation).
func (buf *Buffer) Name() string { return buf.name }

// Kind returns the buffer's kind (Dry, Wet, or Default).
func (buf *Buffer) Kind() string { return buf.kind }

// SkipExamine disables the automatic cleanup check that the test must call
// String() to examine the buffer. Useful for [WetBuffer] when you don't care
// about the content. Implements fluent interface.
func (buf *Buffer) SkipExamine() *Buffer {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.examine = false
	return buf
}

// Write implements [io.Writer]. Thread-safe; increments internal write count.
func (buf *Buffer) Write(p []byte) (n int, err error) {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.wc++
	return buf.buf.Write(p)
}

// WriteString implements [io.StringWriter]. Thread-safe; increments write count.
func (buf *Buffer) WriteString(s string) (n int, err error) {
	buf.mx.Lock()
	defer buf.mx.Unlock()
	buf.wc++
	return buf.buf.WriteString(s)
}

// MustWriteString writes s to the buffer and panics on failure. Useful in
// tests where write errors are not expected.
func (buf *Buffer) MustWriteString(s string) int {
	n, _ := buf.WriteString(s) // Panics when out-of-memory.
	return n
}

// String returns the current contents as a string (thread-safe) and
// increments the read counter. Implements [fmt.Stringer]. This is the
// primary way tests inspect what was written.
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

// DryBuffer creates a thread-safe [Buffer] with the [BufferDry] kind.
//
// It registers a cleanup check via the provided [tester.T] that fails the test
// if any data was written to the buffer.
//
// See the package [README] and [WetBuffer] for the complementary "wet" behavior.
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

// WetBuffer creates a thread-safe [Buffer] with the [BufferWet] kind.
//
// It registers a cleanup check via the provided [tester.T] that fails the test
// if no data was written, or (by default) if the contents were never examined
// via [Buffer.String].
//
// Use [Buffer.SkipExamine] to disable the examination requirement.
//
// See the package [README] and [DryBuffer] for the complementary "dry" behavior.
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
