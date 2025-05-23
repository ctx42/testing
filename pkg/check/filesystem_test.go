// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_FileExist(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// --- When ---
		err := FileExist("testdata/file.txt")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - does not exist", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := FileExist("testdata/not_existing.txt", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to an existing file:\n" +
			"  trail: type.field\n" +
			"   path: testdata/not_existing.txt"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - is a directory", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := FileExist("testdata/dir", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to be existing file:\n" +
			"  trail: type.field\n" +
			"   path: testdata/dir"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NoFileExist(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// --- When ---
		err := NoFileExist("testdata/file.txt")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to not existing file:\n" +
			"  path: testdata/file.txt"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("does not exist", func(t *testing.T) {
		// --- When ---
		err := NoFileExist("testdata/not_existing.txt")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("is a directory", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NoFileExist("testdata/dir", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to be not existing file:\n" +
			"  trail: type.field\n" +
			"   path: testdata/dir"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_FileContain(t *testing.T) {
	t.Run("contains string", func(t *testing.T) {
		// --- When ---
		err := FileContain("ghi\njkl", "testdata/file.txt")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("contains byte slice", func(t *testing.T) {
		// --- When ---
		err := FileContain([]byte("ghi\njkl"), "testdata/file.txt")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("does not contain string", func(t *testing.T) {
		// --- When ---
		err := FileContain("not there", "testdata/file.txt")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected file to contain string:\n" +
			"  path: testdata/file.txt\n" +
			"  want: \"not there\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("does not contain byte slice", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := FileContain([]byte("not there"), "testdata/file.txt", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected file to contain string:\n" +
			"  trail: type.field\n" +
			"   path: testdata/file.txt\n" +
			"   want: \"not there\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("file does not exist", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := FileContain("content", "testdata/not_existing.txt", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected no error reading file:\n" +
			"  trail: type.field\n" +
			"   path: testdata/not_existing.txt\n" +
			"  error: open testdata/not_existing.txt: no such file or directory"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("is a directory", func(t *testing.T) {
		// --- When ---
		err := FileContain("content", "testdata/dir")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected no error reading file:\n" +
			"   path: testdata/dir\n" +
			"  error: read testdata/dir: is a directory"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_DirExist(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// --- When ---
		err := DirExist("testdata/dir")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("does not exist", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := DirExist("testdata/not_existing_dir", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to an existing directory:\n" +
			"  trail: type.field\n" +
			"   path: testdata/not_existing_dir"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("is a file", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := DirExist("testdata/file.txt", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to be existing directory:\n" +
			"  trail: type.field\n" +
			"   path: testdata/file.txt"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NoDirExist(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NoDirExist("testdata/dir", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to not existing directory:\n" +
			"  trail: type.field\n" +
			"   path: testdata/dir"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("does not exist", func(t *testing.T) {
		// --- When ---
		err := NoDirExist("testdata/not_existing_dir")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("is a file", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NoDirExist("testdata/file.txt", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected path to be not existing directory:\n" +
			"  trail: type.field\n" +
			"   path: testdata/file.txt"
		affirm.Equal(t, wMsg, err.Error())
	})
}
