// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package memfs

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/kit/iokit"
	"github.com/ctx42/testing/pkg/must"
)

func Test_compare(t *testing.T) {
	t.Run("OSFile", func(t *testing.T) { TstFile(t, createOSFile) })
	t.Run("KitFile", func(t *testing.T) { TstFile(t, createFile) })
}

func TstFile(t *testing.T, create creator) {
	t.Helper()
	dir := t.TempDir()

	t.Run("Write - O_APPEND", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR|os.O_APPEND, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Write([]byte{3, 4})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 4}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Write - O_APPEND with seek", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR|os.O_APPEND, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		have, err := fil.Write([]byte{3, 4})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 4}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Write - partial override and append", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		have, err := fil.Write([]byte{3, 4, 5})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(4), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 3, 4, 5}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Write - override tail", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		have, err := fil.Write([]byte{3, 4})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 3, 4}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Write - override middle", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2, 3})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		have, err := fil.Write([]byte{8, 9})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 8, 9, 3}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Write - O_APPEND with seek and read", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR|os.O_APPEND, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		have, err := fil.Write([]byte{3, 4})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{}, must.Value(io.ReadAll(fil)))
	})

	t.Run("WriteAt - empty file at the beginning", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, nil)

		// --- When ---
		have, err := fil.WriteAt([]byte{0, 1, 2}, 0)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteAt - append", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.WriteAt([]byte{3, 4}, 3)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 4}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteAt - partial override and append", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.WriteAt([]byte{3, 4, 5}, 1)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 3, 4, 5}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteAt - override tail", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.WriteAt([]byte{3, 4}, 1)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 3, 4}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteAt - override middle", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2, 3})

		// --- When ---
		have, err := fil.WriteAt([]byte{8, 9}, 1)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 8, 9, 3}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteAt - beyond end", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.WriteAt([]byte{4, 5}, 4)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 0, 4, 5}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteAt - with O_APPEND", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR|os.O_APPEND, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.WriteAt([]byte{3, 4}, 1)

		// --- Then ---
		assert.ErrorEqual(t, errWriteAtInAppendMode.Error(), err)
		assert.Equal(t, 0, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, iokit.ReadAllFromStart(fil))
	})

	t.Run("WriteTo - after seek", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)
		dst := &bytes.Buffer{}

		// --- When ---
		have, err := fil.WriteTo(dst)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(2), have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{1, 2}, dst.Bytes())
	})

	t.Run("WriteString - after seek", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		have, err := fil.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(4), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 97, 98, 99}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Read - empty", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, nil)
		dst := make([]byte, 3)

		// --- When ---
		have, err := fil.Read(dst)

		// --- Then ---
		assert.ErrorIs(t, err, io.EOF)
		assert.Equal(t, 0, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 0, 0}, dst)
	})

	t.Run("Read - with a small buffer", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2, 3, 4})
		dst := make([]byte, 3)

		// --- When --- First read.
		have, err := fil.Read(dst)

		// --- Then --- First read.
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, dst)

		// --- When --- Second read.
		have, err = fil.Read(dst)

		// --- Then --- Second read.
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{3, 4, 2}, dst)

		// --- When --- Third read.
		have, err = fil.Read(dst)

		// --- Then --- Third read.
		assert.ErrorIs(t, io.EOF, err)
		assert.Equal(t, 0, have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{3, 4, 2}, dst)
	})

	t.Run("Read - beyond the end", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2, 3})
		iokit.Seek(fil, 5, io.SeekStart)
		dst := make([]byte, 3)

		// --- When ---
		have, err := fil.Read(dst)

		// --- Then ---
		assert.ErrorIs(t, err, io.EOF)
		assert.Equal(t, 0, have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 0, 0}, dst)
	})

	t.Run("Read - big buffer", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2, 3})
		dst := make([]byte, 6)

		// --- When ---
		have, err := fil.Read(dst)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 4, have)
		assert.Equal(t, int64(4), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 0, 0}, dst)
	})

	t.Run("Read - dst with bigger capacity", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		dst := make([]byte, 3, 6)

		// --- When ---
		have, err := fil.Read(dst)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, dst)
	})

	t.Run("Read - after seek", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)
		dst := make([]byte, 3)

		// --- When ---
		have, err := fil.Read(dst)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{1, 2, 0}, dst)
	})

	t.Run("ReadAt - empty", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, nil)
		dst := make([]byte, 3)

		// --- When ---
		have, err := fil.ReadAt(dst, 0)

		// --- Then ---
		assert.ErrorIs(t, io.EOF, err)
		assert.Equal(t, 0, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 0, 0}, dst)
	})

	t.Run("ReadAt - read all - buffer bigger than content", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		dst := make([]byte, 4)

		// --- When ---
		have, err := fil.ReadAt(dst, 0)

		// --- Then ---
		assert.ErrorIs(t, io.EOF, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 0}, dst)
	})

	t.Run("ReadAt - read all - buffer length equal content", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		dst := make([]byte, 3)

		// --- When ---
		have, err := fil.ReadAt(dst, 0)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, dst)
	})

	t.Run("ReadAt - beyond length", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		dst := make([]byte, 3)

		// --- When ---
		have, err := fil.ReadAt(dst, 4)

		// --- Then ---
		assert.ErrorIs(t, io.EOF, err)
		assert.Equal(t, 0, have)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 0, 0}, dst)
	})

	t.Run("ReadAt - big buffer after seek", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)
		dst := make([]byte, 4)

		// --- When ---
		have, err := fil.ReadAt(dst, 0)

		// --- Then ---
		assert.ErrorIs(t, io.EOF, err)
		assert.Equal(t, 3, have)
		assert.Equal(t, int64(1), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 0}, dst)
	})

	t.Run("ReadFrom - to empty file", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, nil)
		rdr := bytes.NewBuffer([]byte{0, 1, 2})

		// --- When ---
		have, err := fil.ReadFrom(rdr)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(3), have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, iokit.ReadAllFromStart(fil))
	})

	t.Run("ReadFrom - overwrite", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte("abc"))
		rdr := bytes.NewBuffer([]byte{0, 1, 2})

		// --- When ---
		have, err := fil.ReadFrom(rdr)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(3), have)
		assert.Equal(t, int64(3), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, iokit.ReadAllFromStart(fil))
	})

	t.Run("ReadFrom - append", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR|os.O_APPEND, []byte{0, 1, 2})
		rdr := bytes.NewBuffer([]byte{3, 4})

		// --- When ---
		have, err := fil.ReadFrom(rdr)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(2), have)
		assert.Equal(t, int64(5), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 4}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Seek - SeekCurrent zero offset after creation", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(0, io.SeekCurrent)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), have)
		assert.Equal(t, []byte{0, 1, 2}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - zero offset SeekEnd", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(0, io.SeekEnd)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(3), have)
		assert.Equal(t, []byte{}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - SeekEnd minus one offset", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(-1, io.SeekEnd)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(2), have)
		assert.Equal(t, []byte{2}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - SeekEnd negative offset to the beginning", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(-3, io.SeekEnd)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), have)
		assert.Equal(t, []byte{0, 1, 2}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - SeekEnd negative error", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(-4, io.SeekEnd)

		// --- Then ---
		var e *fs.PathError
		assert.ErrorAs(t, &e, err)
		assert.Equal(t, int64(0), have)
		assert.Equal(t, []byte{0, 1, 2}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - SeekStart zero offset", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(0, io.SeekStart)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), have)
		assert.Equal(t, []byte{0, 1, 2}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - SeekStart middle offset", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(2, io.SeekStart)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(2), have)
		assert.Equal(t, []byte{2}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Seek - SeekStart beyond length", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		have, err := fil.Seek(4, io.SeekStart)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(4), have)
		assert.Equal(t, []byte{}, must.Value(io.ReadAll(fil)))
	})

	t.Run("Truncate - to zero", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		err := fil.Truncate(0)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - to one", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		err := fil.Truncate(1)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - to zero and write", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		err := fil.Truncate(0)

		// --- Then ---
		assert.NoError(t, err)

		assert.Equal(t, 2, must.Value(fil.Write([]byte{8, 9})))
		assert.Equal(t, int64(2), iokit.Offset(fil))
		assert.Equal(t, []byte{8, 9}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - beyond length", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2, 3})

		// --- When ---
		err := fil.Truncate(6)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 0, 0}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - beyond length with O_APPEND", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR|os.O_APPEND, []byte{0, 1, 2, 3})

		// --- When ---
		err := fil.Truncate(6)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, int64(0), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2, 3, 0, 0}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - beyond capacity and write", func(t *testing.T) {
		// --- Given ---
		content := make([]byte, 3, 4)
		copy(content, []byte{0, 1, 2})
		fil := create(t, dir, os.O_RDWR, content)

		// --- When ---
		err := fil.Truncate(6)

		// --- Then ---
		assert.NoError(t, err)

		assert.Equal(t, 2, must.Value(fil.Write([]byte{8, 9})))
		assert.Equal(t, int64(2), iokit.Offset(fil))
		assert.Equal(t, []byte{8, 9, 2, 0, 0, 0}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - beyond capacity and reset and write", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		assert.NoError(t, fil.Truncate(6))
		assert.NoError(t, fil.Truncate(0))

		// --- Then ---
		assert.Equal(t, 2, must.Value(fil.Write([]byte{8, 9})))
		assert.Equal(t, int64(2), iokit.Offset(fil))
		assert.Equal(t, []byte{8, 9}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - to capacity and write", func(t *testing.T) {
		// --- Given ---
		content := make([]byte, 3, 4)
		copy(content, []byte{0, 1, 2})
		fil := create(t, dir, os.O_RDWR, content)

		// --- When ---
		err := fil.Truncate(4)

		// --- Then ---
		assert.NoError(t, err)

		assert.Equal(t, 2, must.Value(fil.Write([]byte{8, 9})))
		assert.Equal(t, int64(2), iokit.Offset(fil))
		assert.Equal(t, []byte{8, 9, 2, 0}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - above length but less than capacity", func(t *testing.T) {
		// --- Given ---
		content := make([]byte, 3, 5)
		copy(content, []byte{0, 1, 2})
		fil := create(t, dir, os.O_RDWR, content)

		// --- When ---
		err := fil.Truncate(4)

		// --- Then ---
		assert.NoError(t, err)

		assert.Equal(t, 2, must.Value(fil.Write([]byte{8, 9})))
		assert.Equal(t, int64(2), iokit.Offset(fil))
		assert.Equal(t, []byte{8, 9, 2, 0}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - to length and write", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})

		// --- When ---
		err := fil.Truncate(3)

		// --- Then ---
		assert.NoError(t, err)

		assert.Equal(t, 2, must.Value(fil.Write([]byte{8, 9})))
		assert.Equal(t, int64(2), iokit.Offset(fil))
		assert.Equal(t, []byte{8, 9, 2}, iokit.ReadAllFromStart(fil))
	})

	t.Run("Truncate - seek than negative offset", func(t *testing.T) {
		// --- Given ---
		fil := create(t, dir, os.O_RDWR, []byte{0, 1, 2})
		iokit.Seek(fil, 1, io.SeekStart)

		// --- When ---
		err := fil.Truncate(-1)

		// --- Then ---
		var e *fs.PathError
		assert.ErrorAs(t, &e, err)
		assert.Equal(t, int64(1), iokit.Offset(fil))
		assert.Equal(t, []byte{0, 1, 2}, iokit.ReadAllFromStart(fil))
	})
}
