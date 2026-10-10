// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// FileExist asserts that "pth" points to an existing file.
//
// See [check.FileExist] for the error-returning form.
func FileExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.FileExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoFileExist asserts that nothing exists at "pth". A directory at "pth"
// fails the assertion too.
//
// See [check.NoFileExist] for the error-returning form.
func NoFileExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.NoFileExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// FileContain asserts that the file at "pth" contains "want" (after reading
// it in full). Uses [check.FileContain] internally.
// See [check.FileContain].
func FileContain[T check.Content](
	t tester.T,
	want T,
	pth string,
	opts ...any,
) bool {

	t.Helper()
	if e := check.FileContain(want, pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// DirExist asserts that "pth" points to an existing directory.
// See [check.DirExist].
func DirExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.DirExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// NoDirExist asserts that nothing exists at "pth". A file at "pth" fails the
// assertion too.
// See [check.NoDirExist].
func NoDirExist(t tester.T, pth string, opts ...any) bool {
	t.Helper()
	if e := check.NoDirExist(pth, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
