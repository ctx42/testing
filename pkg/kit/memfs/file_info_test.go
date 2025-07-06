package memfs

import (
	"io/fs"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_FileInfo_Name(t *testing.T) {
	// --- Given ---
	fi := FileInfo{name: "my-name"}

	// --- When ---
	have := fi.Name()

	// --- Then ---
	assert.Equal(t, "my-name", have)
}

func Test_FileInfo_Size(t *testing.T) {
	// --- Given ---
	fi := FileInfo{size: 42}

	// --- When ---
	have := fi.Size()

	// --- Then ---
	assert.Equal(t, int64(42), have)
}

func Test_FileInfo_Mode(t *testing.T) {
	// --- Given ---
	fi := FileInfo{}

	// --- When ---
	have := fi.Mode()

	// --- Then ---
	assert.Equal(t, fs.FileMode(0444), have)
}

func Test_FileInfo_ModTime(t *testing.T) {
	// --- Given ---
	fi := FileInfo{}

	// --- When ---
	have := fi.ModTime()

	// --- Then ---
	assert.Zero(t, have)
}

func Test_FileInfo_IsDir(t *testing.T) {
	// --- Given ---
	fi := FileInfo{}

	// --- When ---
	have := fi.IsDir()

	// --- Then ---
	assert.False(t, have)
}

func Test_FileInfo_Sys(t *testing.T) {
	// --- Given ---
	fi := FileInfo{}

	// --- When ---
	have := fi.Sys()

	// --- Then ---
	assert.Nil(t, have)
}
