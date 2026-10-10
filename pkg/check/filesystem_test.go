// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"os"
	"path/filepath"
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
		wMsg := "" +
			"expected path to an existing file:\n" +
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
		wMsg := "" +
			"expected path to be an existing file:\n" +
			"  trail: type.field\n" +
			"   path: testdata/dir"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("symlink to a file", func(t *testing.T) {
		// --- Given ---
		lnk := symlink(t, "testdata/file.txt")

		// --- When ---
		err := FileExist(lnk)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - symlink to a directory", func(t *testing.T) {
		// --- Given ---
		lnk := symlink(t, "testdata/dir")

		// --- When ---
		err := FileExist(lnk)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected path to be an existing file:\n" +
			"  path: " + lnk
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NoFileExist(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// --- When ---
		err := NoFileExist("testdata/file.txt")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected path not to exist:\n" +
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
		wMsg := "" +
			"expected path not to exist:\n" +
			"  trail: type.field\n" +
			"   path: testdata/dir"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("dangling symlink", func(t *testing.T) {
		// --- Given ---
		lnk := symlink(t, "testdata/not_existing.txt")

		// --- When ---
		err := NoFileExist(lnk)

		// --- Then ---
		affirm.Nil(t, err)
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
		wMsg := "" +
			"expected the file to contain the string:\n" +
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
		wMsg := "" +
			"expected the file to contain the string:\n" +
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
		wMsg := "" +
			"expected no error reading the file:\n" +
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
		wMsg := "" +
			"expected no error reading the file:\n" +
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
		wMsg := "" +
			"expected path to an existing directory:\n" +
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
		wMsg := "" +
			"expected the path to be an existing directory:\n" +
			"  trail: type.field\n" +
			"   path: testdata/file.txt"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("symlink to a directory", func(t *testing.T) {
		// --- Given ---
		lnk := symlink(t, "testdata/dir")

		// --- When ---
		err := DirExist(lnk)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - symlink to a file", func(t *testing.T) {
		// --- Given ---
		lnk := symlink(t, "testdata/file.txt")

		// --- When ---
		err := DirExist(lnk)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected the path to be an existing directory:\n" +
			"  path: " + lnk
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
		wMsg := "" +
			"expected path not to exist:\n" +
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
		wMsg := "" +
			"expected path not to exist:\n" +
			"  trail: type.field\n" +
			"   path: testdata/file.txt"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("dangling symlink", func(t *testing.T) {
		// --- Given ---
		lnk := symlink(t, "testdata/not_existing_dir")

		// --- When ---
		err := NoDirExist(lnk)

		// --- Then ---
		affirm.Nil(t, err)
	})
}

// symlink creates a symbolic link to "target" in a temporary directory and
// returns the link path.
func symlink(t *testing.T, target string) string {
	t.Helper()
	abs, err := filepath.Abs(target)
	if err != nil {
		t.Fatal(err)
	}
	lnk := filepath.Join(t.TempDir(), "link")
	if err = os.Symlink(abs, lnk); err != nil {
		t.Fatal(err)
	}
	return lnk
}
