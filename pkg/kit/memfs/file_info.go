package memfs

import (
	"io/fs"
	"time"
)

var _ fs.FileInfo = &FileInfo{} // Compile file check.

// FileInfo implements [fs.FileInfo] interface.
type FileInfo struct {
	name string
	size int64
}

func (fi FileInfo) Name() string       { return fi.name }
func (fi FileInfo) Size() int64        { return fi.size }
func (fi FileInfo) Mode() fs.FileMode  { return 0444 }
func (fi FileInfo) ModTime() time.Time { return time.Time{} }
func (fi FileInfo) IsDir() bool        { return false }
func (fi FileInfo) Sys() any           { return nil }
