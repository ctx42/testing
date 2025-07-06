package iokit

import (
	"io"
	"os"
	"path/filepath"
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
		have := ReadAllFromStart(fil)

		// --- Then ---
		assert.Equal(t, []byte("content"), have)
		assert.Equal(t, int64(3), must.Value(fil.Seek(0, io.SeekCurrent)))
	})
}

func Test_Offset(t *testing.T) {
	// --- Given ---
	pth := filepath.Join(t.TempDir(), "test_file.txt")
	content := "line1\nline2\nend"
	must.Nil(os.WriteFile(pth, []byte(content), 0600))

	fil := must.Value(os.Open(pth))
	must.Value(fil.Read(make([]byte, 3)))

	// --- When ---
	have := Offset(fil)

	// --- Then ---
	assert.Equal(t, int64(3), have)
}

func Test_Seek(t *testing.T) {
	// --- Given ---
	pth := filepath.Join(t.TempDir(), "test_file.txt")
	content := "line1\nline2\nend"
	must.Nil(os.WriteFile(pth, []byte(content), 0600))

	fil := must.Value(os.Open(pth))

	// --- When ---
	have := Seek(fil, 4, io.SeekStart)

	// --- Then ---
	assert.Equal(t, int64(4), have)
	assert.Equal(t, "1\nline2\nend", string(must.Value(io.ReadAll(fil))))
}
