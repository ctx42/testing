// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

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

func Test_WithData(t *testing.T) {
	// --- Given ---
	data := map[string]any{"A": 1}
	gld := &Goldy{}

	// --- When ---
	WithData(data)(gld)

	// --- Then ---
	affirm.Equal(t, true, core.Same(data, gld.data))
}

func Test_Open(t *testing.T) {
	t.Run("test runner set", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case1.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, true, core.Same(tspy, have.t))
	})

	t.Run("case 1 - content ends without a new line", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case1.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case1.gld", have.pth)
		affirm.Equal(t, "No new line at the end.\n", have.comment)
		affirm.Equal(t, "Content #1.\nContent #2.", string(have.content))
		affirm.Nil(t, have.tpl)
	})

	t.Run("case 2 - content ends with a new line", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case2.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case2.gld", have.pth)
		affirm.Equal(t, "New line at the end of file.\n", have.comment)
		affirm.Equal(t, "Content #1.\nContent #2.\n", string(have.content))
		affirm.Nil(t, have.tpl)
	})

	t.Run("case 3 - content ends with multiple new lines", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case3.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case3.gld", have.pth)
		affirm.Equal(t, "Multiple new lines at the file end.\n", have.comment)
		affirm.Equal(t, "Content #1.\nContent #2.\n\n", string(have.content))
		affirm.Nil(t, have.tpl)
	})

	t.Run("case 4 - multiple comment lines", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Open(tspy, "testdata/test_case4.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_case4.gld", have.pth)
		affirm.Equal(t, "Multiple\ncomment\nlines.\n", have.comment)
		affirm.Equal(t, "Content #1.\nContent #2.\n", string(have.content))
		affirm.Nil(t, have.tpl)
	})

	t.Run("open a golden file template", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		data := WithData(map[string]any{"first": 1})

		// --- When ---
		have := Open(tspy, "testdata/test_tpl.gld", data)

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.Equal(t, "testdata/test_tpl.gld", have.pth)
		affirm.Equal(t, "Golden file template.\n", have.comment)
		affirm.Equal(t, "Content #1.", string(have.content))
		affirm.Equal(t, "Content #{{ .first }}.", string(have.tpl))
	})

	t.Run("error - no marker", func(t *testing.T) {
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

	t.Run("error - not existing file", func(t *testing.T) {
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

	t.Run("error - golden file invalid template", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		data := WithData(map[string]any{"second": 2})

		// --- When ---
		have := Open(tspy, "testdata/test_tpl_invalid.gld", data)

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "unexpected \"}\" in operand"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})

	t.Run("error - golden file template missing data", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		data := WithData(map[string]any{"other": 1})

		// --- When ---
		have := Open(tspy, "testdata/test_tpl.gld", data)

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "map has no entry for key \"first\""
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})
}

func Test_Goldy_String(t *testing.T) {
	// --- Given ---
	gld := &Goldy{content: []byte("content")}

	// --- When ---
	have := gld.String()

	// --- Then ---
	affirm.Equal(t, "content", have)
}

func Test_Goldy_Bytes(t *testing.T) {
	// --- Given ---
	gld := &Goldy{content: []byte("content")}

	// --- When ---
	have := gld.Bytes()

	// --- Then ---
	affirm.DeepEqual(t, []byte("content"), have)
	affirm.Equal(t, false, core.Same(gld.content, have))
}

func Test_Goldy_SetContent(t *testing.T) {
	t.Run("regular golden file", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		gld := Open(tspy, "testdata/test_case1.gld")

		// --- When ---
		have := gld.SetContent("abc")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.DeepEqual(t, []byte("abc"), have.content)
		affirm.Equal(t, true, core.Same(gld, have))
	})

	t.Run("template golden file", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		data := WithData(map[string]any{"first": 1})
		gld := Open(tspy, "testdata/test_tpl.gld", data)

		// --- When ---
		have := gld.SetContent("32{{.first}}")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.DeepEqual(t, []byte("321"), have.content)
		affirm.DeepEqual(t, []byte("32{{.first}}"), have.tpl)
		affirm.Equal(t, true, core.Same(gld, have))
	})

	t.Run("error - rendering template", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		data := WithData(map[string]any{"first": 1})
		gld := Open(tspy, "testdata/test_tpl.gld", data)

		// --- When ---
		have := gld.SetContent("32{{.other}}")

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "map has no entry for key \"other\""
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})
}

func Test_Goldy_Save(t *testing.T) {
	t.Run("when not a template", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		pth := filepath.Join(t.TempDir(), "/golden.gld")
		gld := &Goldy{
			pth:     pth,
			comment: "Comment 1.\nComment 2.\n",
			content: []byte("content 1\ncontent 2\n"),
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

	t.Run("when a template", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		pth := filepath.Join(t.TempDir(), "/golden.gld")
		gld := &Goldy{
			pth:     pth,
			comment: "comment\n",
			content: []byte("content 1"),
			data:    map[string]any{"first": 1},
			tpl:     []byte("content {{ .first }}"),
			t:       tspy,
		}

		// --- When ---
		gld.Save()

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		have := must.Value(os.ReadFile(pth))
		want := "comment\n---\ncontent {{ .first }}"
		affirm.Equal(t, want, string(have))
	})

	t.Run("comment lines do not end with a new line", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()
		pth := filepath.Join(t.TempDir(), "/golden.gld")
		gld := &Goldy{
			pth:     pth,
			comment: "Comment 1.\nComment 2.",
			content: []byte("content 1\ncontent 2\n"),
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

	t.Run("error - when the file cannot be written", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()
		pth := filepath.Join(t.TempDir(), "sub-dir", "/golden.gld")
		gld := &Goldy{
			pth:     pth,
			comment: "Comment 1.\nComment 2.",
			content: []byte("content 1\ncontent 2\n"),
			t:       tspy,
		}

		// --- When ---
		gld.Save()

		// --- Then ---
		affirm.Equal(t, true, tspy.Failed())
		affirm.Equal(t, true, strings.Contains(tspy.Log(), pth))
	})
}
