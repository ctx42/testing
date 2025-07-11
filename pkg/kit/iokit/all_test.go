// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"errors"
)

var (
	ErrSadRead  = errors.New("sad read error")
	ErrSadClose = errors.New("sad close error")
	ErrSadSeek  = errors.New("sad seek error")
	ErrSadWrite = errors.New("sad write error")
)

// sadStruct is sad because all its methods return errors.
type sadStruct struct{}

func (_ sadStruct) Read(_ []byte) (int, error)     { return 0, ErrSadRead }
func (_ sadStruct) Seek(int64, int) (int64, error) { return 0, ErrSadSeek }
func (_ sadStruct) Write(_ []byte) (int, error)    { return 0, ErrSadWrite }
func (_ sadStruct) Close() error                   { return ErrSadClose }

// blissfulWriter is happy because it does not do anything and never complains.
type blissfulWriter struct{}

func (w blissfulWriter) Write(_ []byte) (int, error) { return 0, nil }
func (w blissfulWriter) Close() error                { return nil }
