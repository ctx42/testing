// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package memfs_test

import (
	"bytes"
	"testing"

	"github.com/ctx42/testing/pkg/kit/memfs"
)

//goland:noinspection GoUnusedGlobalVariable
var bufferWrite int

func BenchmarkFileWrite(b *testing.B) {
	data := make([]byte, 1<<15)

	b.Run("fskit", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var n int

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &memfs.File{}
			n, _ = buf.Write(data)
		}
		bufferWrite = n
	})

	b.Run("bytes", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var n int

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			n, _ = buf.Write(data)
		}
		bufferWrite = n
	})
}

//goland:noinspection GoUnusedGlobalVariable
var bufferWriteByte error

func BenchmarkFileWriteByte(b *testing.B) {
	b.Run("fskit", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var err error

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &memfs.File{}
			err = buf.WriteByte(1)
		}
		bufferWriteByte = err
	})

	b.Run("bytes", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var err error

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			err = buf.WriteByte(1)
		}
		bufferWriteByte = err
	})
}

//goland:noinspection GoUnusedGlobalVariable
var bufferWriteString int

func BenchmarkFileWriteString(b *testing.B) {
	b.Run("fskit", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var n int

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &memfs.File{}
			n, _ = buf.WriteString("abcdefghijkl")
		}
		bufferWriteString = n
	})

	b.Run("bytes", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var n int

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			n, _ = buf.WriteString("abcdefghijkl")
		}
		bufferWriteString = n
	})
}

//goland:noinspection GoUnusedGlobalVariable
var bufferReadFrom int64

func BenchmarkFileReadFrom(b *testing.B) {
	data := make([]byte, 1<<15)

	b.Run("fskit", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var n int64
		src := bytes.NewReader(data)

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &memfs.File{}
			n, _ = buf.ReadFrom(src)
			src.Reset(data)
		}
		bufferReadFrom = n
	})

	b.Run("bytes", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		var n int64
		src := bytes.NewReader(data)

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			n, _ = buf.ReadFrom(src)
			src.Reset(data)
		}
		bufferReadFrom = n
	})
}
