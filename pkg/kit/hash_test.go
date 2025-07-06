package kit

import (
	"bytes"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_SHA1Reader(t *testing.T) {
	// --- Given ---
	rdr := bytes.NewReader([]byte("content"))

	// --- When ---
	have := SHA1Reader(rdr)

	// --- Then ---
	assert.Equal(t, "040f06fd774092478d450774f5ba30c5da78acc8", have)
}

func Test_SHA1File(t *testing.T) {
	// --- When ---
	have := SHA1File("testdata/file.txt")

	// --- Then ---
	assert.Equal(t, "040f06fd774092478d450774f5ba30c5da78acc8", have)
}
