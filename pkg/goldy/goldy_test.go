package goldy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/must"
)

func Test_Open(t *testing.T) {
	t.Run("test runner set", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case1.gld")

		// --- Then ---
		affirm.Equal(t, true, core.Same(tspy, have.t))
	})

	t.Run("success case 1", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case1.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case1.gld", have.Path)
		affirm.Equal(t, "No new line at the end.\n", have.Comment)
		affirm.Equal(t, "Content #1.\nContent #2.", have.String())
	})

	t.Run("success case 2", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case2.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case2.gld", have.Path)
		affirm.Equal(t, "New line at the end of file.\n", have.Comment)
		affirm.Equal(t, "Content #1.\nContent #2.\n", have.String())
	})

	t.Run("success case 3", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case3.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case3.gld", have.Path)
		affirm.Equal(t, "Multiple new lines at the file end.\n", have.Comment)
		affirm.Equal(t, "Content #1.\nContent #2.\n\n", have.String())
	})

	t.Run("success case 4", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case4.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case4.gld", have.Path)
		affirm.Equal(t, "Multiple\ncomment\nlines.\n", have.Comment)
		affirm.Equal(t, "Content #1.\nContent #2.\n", have.String())
	})

	t.Run("error no marker", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := Open(tspy, "testdata/test_no_marker.gld")

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "the golden file is missing the \"---\" marker"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})

	t.Run("not existing file", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := Open(tspy, "testdata/not-existing.txt")

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "no such file or directory"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
		wMsg = "testdata/not-existing.txt"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})
}

func Test_New(t *testing.T) {
	t.Run("open with string content", func(t *testing.T) {
		// --- Given ---
		pth := filepath.Join(t.TempDir(), "test.gld")
		tspy := core.NewSpy()

		// --- When ---
		have := New(tspy, pth, "comment", "content")

		// --- Then ---
		affirm.Equal(t, pth, have.Path)
		affirm.Equal(t, "comment", have.Comment)
		affirm.DeepEqual(t, []byte("content"), have.Content)
		affirm.Equal(t, true, core.Same(tspy, have.t))
	})

	t.Run("open with byte content", func(t *testing.T) {
		// --- Given ---
		pth := filepath.Join(t.TempDir(), "test.gld")
		tspy := core.NewSpy()

		// --- When ---
		have := New(tspy, pth, "comment", []byte("content"))

		// --- Then ---
		affirm.Equal(t, pth, have.Path)
		affirm.Equal(t, "comment", have.Comment)
		affirm.DeepEqual(t, []byte("content"), have.Content)
		affirm.Equal(t, true, core.Same(tspy, have.t))
	})
}

func Test_Goldy_String(t *testing.T) {
	// --- Given ---
	gld := &Goldy{Content: []byte("content")}

	// --- When ---
	have := gld.String()

	// --- Then ---
	affirm.Equal(t, "content", have)
}

func Test_Goldy_Bytes(t *testing.T) {
	// --- Given ---
	gld := &Goldy{Content: []byte("content")}

	// --- When ---
	have := gld.Bytes()

	// --- Then ---
	affirm.DeepEqual(t, []byte("content"), have)
	affirm.Equal(t, false, core.Same(gld.Content, have))
}

func Test_Goldy_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		pth := filepath.Join(t.TempDir(), "/golden.gld")
		gld := &Goldy{
			Path:    pth,
			Comment: "Comment 1.\nComment 2.\n",
			Content: []byte("content 1\ncontent 2\n"),
			t:       tspy,
		}

		// --- When ---
		gld.Save()

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		have := must.Value(os.ReadFile(pth))
		want := "Comment 1.\nComment 2.\n---\ncontent 1\ncontent 2\n"
		affirm.Equal(t, want, string(have))
	})

	t.Run("comment lines do not end with a new line", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		pth := filepath.Join(t.TempDir(), "/golden.gld")
		gld := &Goldy{
			Path:    pth,
			Comment: "Comment 1.\nComment 2.",
			Content: []byte("content 1\ncontent 2\n"),
			t:       tspy,
		}

		// --- When ---
		gld.Save()

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		have := must.Value(os.ReadFile(pth))
		want := "Comment 1.\nComment 2.\n---\ncontent 1\ncontent 2\n"
		affirm.Equal(t, want, string(have))
	})

	t.Run("error when the file cannot be written", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		pth := filepath.Join(t.TempDir(), "sub-dir", "/golden.gld")
		gld := &Goldy{
			Path:    pth,
			Comment: "Comment 1.\nComment 2.",
			Content: []byte("content 1\ncontent 2\n"),
			t:       tspy,
		}

		// --- When ---
		gld.Save()

		// --- Then ---
		affirm.Equal(t, true, tspy.Failed())
		affirm.Equal(t, true, strings.Contains(tspy.Log(), pth))
	})
}
