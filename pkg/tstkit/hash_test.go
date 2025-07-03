package tstkit

import (
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_SHA1Reader(t *testing.T) {
	// --- Given ---
	fil := must.Value(os.Open("testdata/file.txt"))

	// --- When ---
	have := SHA1Reader(fil)

	// --- Then ---
	assert.Equal(t, "040f06fd774092478d450774f5ba30c5da78acc8", have)
}

func Test_SHA1File(t *testing.T) {
	// --- When ---
	have := SHA1File("testdata/file.txt")

	// --- Then ---
	assert.Equal(t, "040f06fd774092478d450774f5ba30c5da78acc8", have)
}
