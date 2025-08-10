// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"os"
	"strings"

	"github.com/ctx42/testing/pkg/notice"
)

// FileExist checks "pth" points to an existing file. Returns an error if the
// path points to a filesystem entry which is not a file or there is an error
// when trying to check the path. On success, it returns nil.
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
		msg := notice.New("expected path to be existing file").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}
	return nil
}

// NoFileExist checks "pth" points to not existing file. Returns an error if
// the path points to an existing filesystem entry. On success, it returns nil.
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
		msg := notice.New("expected path to be not existing file").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected path to not existing file").
		Append("path", "%s", pth)
	return AddRows(ops, msg)
}

// Content declares type constraint for file content.
type Content interface {
	string | []byte
}

// FileContain checks file at "pth" can be read and its string content contains
// "want". It fails if the path points to a filesystem entry which is not a
// file or there is an error reading the file. The file is read in full then
// [strings.Contains] is used to check it contains "want" string. When it fails
// it returns an error with a message indicating the expected and actual values.
func FileContain[T Content](want T, pth string, opts ...any) error {
	content, err := os.ReadFile(pth)
	if err != nil {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected no error reading file").
			Append("path", "%s", pth).
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if strings.Contains(string(content), string(want)) {
		return nil
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected file to contain string").
		Append("path", "%s", pth).
		Want("%q", want)
	return AddRows(ops, msg)
}

// DirExist checks "pth" points to an existing directory. It fails if the path
// points to a filesystem entry which is not a directory or there is an error
// when trying to check the path. When it fails it returns an error with a
// detailed message indicating the expected and actual values.
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
		msg := notice.New("expected path to be existing directory").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}
	return nil
}

// NoDirExist checks "pth" points to not existing directory. It fails if the
// path points to an existing filesystem entry. When it fails it returns an
// error with a detailed message indicating the expected and actual values.
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
		msg := notice.New("expected path to be not existing directory").
			Append("path", "%s", pth)
		return AddRows(ops, msg)
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected path to not existing directory").
		Append("path", "%s", pth)
	return AddRows(ops, msg)
}
