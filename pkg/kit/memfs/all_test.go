// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package memfs

import (
	"io"
	"os"

	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/tester"
)

// file is an interface common to [os.File] and [File].
type file interface {
	io.Seeker
	io.Reader
	io.ReaderAt
	io.Closer
	io.ReaderFrom
	io.Writer
	io.WriterAt
	io.StringWriter
	io.WriterTo

	Truncate(size int64) error
}

// creator creates a new instance of [file] with given content (if it is not
// nil) and the seek offset at its beginning. The returned file will
// automatically be closed at the test end.
type creator func(t tester.T, dir string, flag int, content []byte) file

// createOSFile creates a new [os.File] instance with given content (if it is
// not nil) and the seek offset at its beginning. The returned file will
// automatically be closed at the test end.
func createOSFile(t tester.T, dir string, flag int, content []byte) file {
	t.Helper()
	fil := must.Value(os.CreateTemp(dir, ""))
	pth := fil.Name()
	if content != nil {
		must.Value(fil.Write(content))
	}
	must.Nil(fil.Close())
	fil = must.Value(os.OpenFile(pth, flag, 0666))
	t.Cleanup(func() { t.Helper(); must.Nil(fil.Close()) })
	return fil
}

// createFile creates a [File] instance with given content and the seek offset
// at its beginning. The returned file will automatically be closed at the test
// end.
func createFile(t tester.T, _ string, flag int, content []byte) file {
	t.Helper()
	fil := FileWith(content, WithFileFlag(flag))
	t.Cleanup(func() { t.Helper(); must.Nil(fil.Close()) })
	return fil
}

// TempFile creates and opens a temporary file with contents from the data
// slice. It registers the cleanup function to close the file and remove it.
// The flag parameter is the same as in [os.Open].
//
// On error function calls t.Fatal.
func TempFile(t tester.T, flag int, content []byte) *os.File {
	t.Helper()

	fil, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
		return nil
	}
	pth := fil.Name()
	if _, err = fil.Write(content); err != nil {
		t.Fatal(err)
		return nil
	}
	if err = fil.Close(); err != nil {
		t.Fatal(err)
		return nil
	}

	fil, err = os.OpenFile(pth, flag, 0666)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return fil
}
