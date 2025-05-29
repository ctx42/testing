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
func FileExist(pth string, opts ...Option) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		ops := DefaultOptions(opts...)
		if os.IsNotExist(err) {
			return notice.New("expected path to an existing file").
				SetTrail(ops.Trail).
				Append("path", "%s", pth)
		}
		return notice.New("expected os.Lstat to succeed").
			SetTrail(ops.Trail).
			Append("path", "%s", pth).
			Append("error", "%s", err)
	}
	if inf.IsDir() {
		ops := DefaultOptions(opts...)
		return notice.New("expected path to be existing file").
			SetTrail(ops.Trail).
			Append("path", "%s", pth)
	}
	return nil
}

// NoFileExist checks "pth" points to not existing file. Returns an error if
// the path points to an existing filesystem entry. On success, it returns nil.
func NoFileExist(pth string, opts ...Option) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		ops := DefaultOptions(opts...)
		return notice.New("expected os.Lstat to succeed").
			SetTrail(ops.Trail).
			Append("path", "%s", pth).
			Append("error", "%s", err)
	}
	if inf.IsDir() {
		ops := DefaultOptions(opts...)
		return notice.New("expected path to be not existing file").
			SetTrail(ops.Trail).
			Append("path", "%s", pth)
	}

	ops := DefaultOptions(opts...)
	return notice.New("expected path to not existing file").
		SetTrail(ops.Trail).
		Append("path", "%s", pth)
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
func FileContain[T Content](want T, pth string, opts ...Option) error {
	content, err := os.ReadFile(pth)
	if err != nil {
		ops := DefaultOptions(opts...)
		return notice.New("expected no error reading file").
			SetTrail(ops.Trail).
			Append("path", "%s", pth).
			Append("error", "%s", err)
	}
	if strings.Contains(string(content), string(want)) {
		return nil
	}

	ops := DefaultOptions(opts...)
	return notice.New("expected file to contain string").
		SetTrail(ops.Trail).
		Append("path", "%s", pth).
		Want("%q", want)
}

// DirExist checks "pth" points to an existing directory. It fails if the path
// points to a filesystem entry which is not a directory or there is an error
// when trying to check the path. When it fails it returns an error with a
// detailed message indicating the expected and actual values.
func DirExist(pth string, opts ...Option) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			ops := DefaultOptions(opts...)
			return notice.New("expected path to an existing directory").
				SetTrail(ops.Trail).
				Append("path", "%s", pth)
		}

		ops := DefaultOptions(opts...)
		return notice.New("expected os.Lstat to succeed").
			SetTrail(ops.Trail).
			Append("path", "%s", pth).
			Append("error", "%s", err)
	}
	if !inf.IsDir() {
		ops := DefaultOptions(opts...)
		return notice.New("expected path to be existing directory").
			SetTrail(ops.Trail).
			Append("path", "%s", pth)
	}
	return nil
}

// NoDirExist checks "pth" points to not existing directory. It fails if the
// path points to an existing filesystem entry. When it fails it returns an
// error with a detailed message indicating the expected and actual values.
func NoDirExist(pth string, opts ...Option) error {
	inf, err := os.Lstat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		ops := DefaultOptions(opts...)
		return notice.New("expected os.Lstat to succeed").
			SetTrail(ops.Trail).
			Append("path", "%s", pth).
			Append("error", "%s", err)
	}
	if !inf.IsDir() {
		ops := DefaultOptions(opts...)
		return notice.New("expected path to be not existing directory").
			SetTrail(ops.Trail).
			Append("path", "%s", pth)
	}

	ops := DefaultOptions(opts...)
	return notice.New("expected path to not existing directory").
		SetTrail(ops.Trail).
		Append("path", "%s", pth)
}
