// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package tstkit

import (
	"bytes"
	"sync"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_NewBuffer(t *testing.T) {
	t.Run("without name", func(t *testing.T) {
		// --- When ---
		buf := NewBuffer()

		// --- Then ---
		assert.Empty(t, buf.name)
		assert.Equal(t, BuffDefault, buf.kind)
		assert.Empty(t, buf.buf.String())
		assert.True(t, buf.examine)
		assert.Equal(t, 0, buf.wc)
		assert.Equal(t, 0, buf.rc)
	})

	t.Run("with name", func(t *testing.T) {
		// --- When ---
		buf := NewBuffer("name")

		// --- Then ---
		assert.Equal(t, "name", buf.name)
		assert.Equal(t, BuffDefault, buf.kind)
		assert.Empty(t, buf.buf.String())
		assert.True(t, buf.examine)
		assert.Equal(t, 0, buf.wc)
		assert.Equal(t, 0, buf.rc)
	})
}

func Test_Buffer_Name(t *testing.T) {
	// --- Given ---
	buf := &Buffer{name: "name"}

	// --- When ---
	have := buf.Name()

	// --- Then ---
	assert.Equal(t, "name", have)
}

func Test_Buffer_Kind(t *testing.T) {
	// --- Given ---
	buf := &Buffer{kind: BufferDry}

	// --- When ---
	have := buf.Kind()

	// --- Then ---
	assert.Equal(t, BufferDry, have)
}

func Test_Buffer_SkipExamine(t *testing.T) {
	// --- Given ---
	buf := &Buffer{examine: true}

	// --- When ---
	have := buf.SkipExamine()

	// --- Then ---
	assert.False(t, buf.examine)
	assert.Same(t, buf, have)
}

func Test_Buffer_Write(t *testing.T) {
	// --- Given ---
	buf := NewBuffer()

	// --- When ---
	n, err := buf.Write([]byte{97, 98, 99})

	// --- Then ---
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 1, buf.wc)
	assert.Equal(t, 0, buf.rc)
	assert.Equal(t, "abc", buf.buf.String())
}

func Test_Buffer_WriteString(t *testing.T) {
	// --- Given ---
	buf := NewBuffer()

	// --- When ---
	n, err := buf.WriteString("abc")

	// --- Then ---
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 1, buf.wc)
	assert.Equal(t, 0, buf.rc)
}

func Test_Buffer_MustWriteString(t *testing.T) {
	// --- Given ---
	buf := NewBuffer()

	// --- When ---
	have := buf.MustWriteString("abc")

	// --- Then ---
	assert.Equal(t, 3, have)
	assert.Equal(t, 1, buf.wc)
	assert.Equal(t, 0, buf.rc)
	assert.Equal(t, "abc", buf.buf.String())
}

func Test_Buffer_String(t *testing.T) {
	// --- Given ---
	buf := &Buffer{buf: bytes.NewBuffer([]byte{97, 98, 99})}

	// --- When ---
	have := buf.String()

	// --- Then ---
	assert.Equal(t, 0, buf.wc)
	assert.Equal(t, 1, buf.rc)
	assert.Equal(t, "abc", have)
}

func Test_Buffer_string(t *testing.T) {
	t.Run("do not increase read counter", func(t *testing.T) {
		// --- Given ---
		buf := &Buffer{buf: bytes.NewBuffer([]byte{97, 98, 99})}

		// --- When ---
		have := buf.string(false)

		// --- Then ---
		assert.Equal(t, 0, buf.wc)
		assert.Equal(t, 0, buf.rc)
		assert.Equal(t, "abc", have)
	})

	t.Run("do increase read counter", func(t *testing.T) {
		// --- Given ---
		buf := &Buffer{buf: bytes.NewBuffer([]byte{97, 98, 99})}

		// --- When ---
		have := buf.string(true)

		// --- Then ---
		assert.Equal(t, 0, buf.wc)
		assert.Equal(t, 1, buf.rc)
		assert.Equal(t, "abc", have)
	})
}

func Test_Buffer_Reset(t *testing.T) {
	// --- Given ---
	buf := Buffer{
		name:    "name",
		kind:    BuffDefault,
		buf:     bytes.NewBuffer([]byte{97, 98, 99}),
		mx:      sync.Mutex{},
		examine: false,
		wc:      1,
		rc:      2,
	}

	// --- When ---
	buf.Reset()

	// --- Then ---
	assert.Equal(t, "name", buf.name)
	assert.Equal(t, BuffDefault, buf.kind)
	assert.Equal(t, "", buf.buf.String())
	assert.False(t, buf.examine)
	assert.Equal(t, 0, buf.wc)
	assert.Equal(t, 0, buf.rc)
}

func Test_DryBuffer(t *testing.T) {
	t.Run("kind set", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)

		// --- Then ---
		assert.Equal(t, BufferDry, buf.Kind())
	})

	t.Run("buffer is not written to", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)

		// --- Then ---
		assert.NotNil(t, buf)
	})

	t.Run("error - buffer written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer to be empty:\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("calling SkipExamine does not change behavior", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer to be empty:\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy).SkipExamine()
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("error - buffer named and written to", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer to be empty:\n" +
			"  name: buf-name\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy, "buf-name")
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})
}

func Test_WetBuffer(t *testing.T) {
	t.Run("kind set", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)

		// --- Then ---
		assert.Equal(t, BufferWet, buf.Kind())

		buf.SkipExamine().MustWriteString("abc") // So the test passes.
	})

	t.Run("buffer written to", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "abc", buf.String())
	})

	t.Run("error - buffer written to but not examined", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer contents to be examined"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("error - named buffer written to but not examined", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer contents to be examined:\n  name: buf-name"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy, "buf-name")
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("can skip examination", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy).SkipExamine()
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("error - buffer is not written to", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer not to be empty"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)

		// --- Then ---
		assert.NotNil(t, buf)
	})

	t.Run("error - buffer named and not written to", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer not to be empty:\n  name: buf-name"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy, "buf-name")

		// --- Then ---
		assert.NotNil(t, buf)
	})
}
