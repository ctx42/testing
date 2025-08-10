// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// FileExist asserts "pth" points to an existing file. It fails if the path
// points to a filesystem entry, which is not a file, or there is an error when
// trying to check the path. Returns true on success, otherwise marks the test
// as failed, writes an error message to the test log and returns false.
func FileExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.FileExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoFileExist asserts "pth" points to a not existing file. It fails if the
// path points to an existing filesystem entry. Returns true on success,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func NoFileExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.NoFileExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// FileContain asserts a file at "pth" can be read and its string content
// contains "want". It fails if the path points to a filesystem entry, which is
// not a file, or there is an error reading the file. The file is read in full,
// then [Contain] assertion is used to check it contains the "want" string.
// Returns true on success, otherwise marks the test as failed, writes an error
// message to the test log and returns false.
func FileContain[T check.Content](t tester.T, want T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.FileContain(want, pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// DirExist asserts "pth" points to an existing directory. It fails if the path
// points to a filesystem entry, which is not a directory, or there is an error
// when trying to check the path. Returns true on success, otherwise marks the
// test as failed, writes an error message to the test log and returns false.
func DirExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.DirExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoDirExist asserts "pth" points to not existing directory. It fails if the
// path points to an existing filesystem entry. Returns true on success,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
func NoDirExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.NoDirExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
