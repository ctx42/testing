package tstkit

import (
	"io"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_ReadAllFromStart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		fil := must.Value(os.Open("testdata/file.txt"))
		must.Value(fil.Seek(3, io.SeekStart))

		// --- When ---
		got := ReadAllFromStart(fil)

		// --- Then ---
		assert.Equal(t, []byte("content"), got)
		assert.Equal(t, int64(3), must.Value(fil.Seek(0, io.SeekCurrent)))
	})
}
