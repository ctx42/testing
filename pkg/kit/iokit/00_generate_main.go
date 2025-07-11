// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

//go:build ignore

package main

import (
	"github.com/ctx42/testing/pkg/mocker"
)

func main() {
	opts := []mocker.Option{
		mocker.WithTgtOnHelpers,
	}
	mocks := []func(opts ...mocker.Option) error{
		GenIoReaderMock,
		GenIoReadCloserMock,
		GenIoReadSeekerMock,
		GenIoReadSeekCloserMock,
	}
	for _, mock := range mocks {
		if err := mock(opts...); err != nil {
			panic(err)
		}
	}
}

func GenIoReaderMock(opts ...mocker.Option) error {
	opts = append(
		opts,
		mocker.WithSrc("io"),
		mocker.WithTgtFilename("io_reader_mock.go"),
	)
	if err := mocker.Generate("Reader", opts...); err != nil {
		return err
	}
	return nil
}

func GenIoReadCloserMock(opts ...mocker.Option) error {
	opts = append(
		opts,
		mocker.WithSrc("io"),
		mocker.WithTgtFilename("io_read_closer_mock.go"),
	)
	if err := mocker.Generate("ReadCloser", opts...); err != nil {
		return err
	}
	return nil
}

func GenIoReadSeekerMock(opts ...mocker.Option) error {
	opts = append(
		opts,
		mocker.WithSrc("io"),
		mocker.WithTgtFilename("io_read_seeker_mock.go"),
	)
	if err := mocker.Generate("ReadSeeker", opts...); err != nil {
		return err
	}
	return nil
}

func GenIoReadSeekCloserMock(opts ...mocker.Option) error {
	opts = append(
		opts,
		mocker.WithSrc("io"),
		mocker.WithTgtFilename("io_read_seek_closer_mock.go"),
	)
	if err := mocker.Generate("ReadSeekCloser", opts...); err != nil {
		return err
	}
	return nil
}
