// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"os"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// FileExist checks that "pth" points to an existing file.
//
// See [assert.FileExist] for the assertion wrapper.
func FileExist(pth string, opts ...any) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		ops := DefaultOptions(opts...)
		if os.IsNotExist(err) {
			msg := notice.New("expected path to an existing file").
				Append("path", "%s", pth)
			return AddRows(ops, msg)

		}
		msg := notice.New("expected os.Lstat to succeed").
			Append("path", "%s", pth).
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if inf.IsDir() {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected path to be an existing file").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}
	return nil
}

// NoFileExist checks that "pth" does not point to an existing file.
//
// See [assert.NoFileExist] for the assertion wrapper.
func NoFileExist(pth string, opts ...any) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		ops := DefaultOptions(opts...)
		msg := notice.New("expected os.Lstat to succeed").
			Append("path", "%s", pth).
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if inf.IsDir() {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected path to not be an existing file").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected path to not be an existing file").
		Append("path", "%s", pth)
	return AddRows(ops, msg)
}

// Content declares type constraint for file content.
type Content interface {
	string | []byte
}

// FileContain checks that the file at "pth" can be read and contains "want"
// (full read + [strings.Contains]). See [assert.FileContain].
func FileContain[T Content](want T, pth string, opts ...any) error {
	// G304: path comes from test assertions on filesystem content.
	content, err := os.ReadFile(pth) // nolint:gosec
	if err != nil {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected no error reading the file").
			Append("path", "%s", pth).
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if strings.Contains(string(content), string(want)) {
		return nil
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected the file to contain the string").
		Append("path", "%s", pth).
		Want("%q", want)
	return AddRows(ops, msg)
}

// DirExist checks that "pth" points to an existing directory.
// See [assert.DirExist].
func DirExist(pth string, opts ...any) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			ops := DefaultOptions(opts...)
			msg := notice.New("expected path to an existing directory").
				Append("path", "%s", pth)
			return AddRows(ops, msg)
		}

		ops := DefaultOptions(opts...)
		msg := notice.New("expected os.Lstat to succeed").
			Append("path", "%s", pth).
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if !inf.IsDir() {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected the path to be an existing directory").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}
	return nil
}

// NoDirExist checks that "pth" does not point to an existing directory.
// See [assert.NoDirExist].
func NoDirExist(pth string, opts ...any) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		ops := DefaultOptions(opts...)
		msg := notice.New("expected os.Lstat to succeed").
			Append("path", "%s", pth).
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if !inf.IsDir() {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected path to not be an existing directory").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected path to not be an existing directory").
		Append("path", "%s", pth)
	return AddRows(ops, msg)
}
