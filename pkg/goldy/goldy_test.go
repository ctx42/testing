// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package goldy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/tester"
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
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := Open(tspy, "testdata/test_case1.gld")

		// --- Then ---
		affirm.Equal(t, false, tspy.Failed())
		affirm.NotNil(t, have.comment)
		affirm.Equal(t, true, core.Same(tspy, have.t))
	})

	t.Run("case 1 - content ends without a new line", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogEqual("the golden file is missing the \"---\" marker")
		tspy.Close()

		// --- When ---
		have := Open(tspy, "testdata/test_no_marker.gld")

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("error - not existing file", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "error opening file: open testdata/not-existing.txt: " +
			"no such file or directory"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := Open(tspy, "testdata/not-existing.txt")

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("error - golden file invalid template", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogEqual("template: goldy:2: unexpected \"}\" in operand")
		tspy.Close()

		data := WithData(map[string]any{"second": 2})

		// --- When ---
		have := Open(tspy, "testdata/test_tpl_invalid.gld", data)

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("error - golden file template missing data", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "template: goldy:1:12: executing \"goldy\" at <.first>: " +
			"map has no entry for key \"first\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		data := WithData(map[string]any{"other": 1})

		// --- When ---
		have := Open(tspy, "testdata/test_tpl.gld", data)

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("error - golden file not found", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "error opening file: open testdata/does_not_exist.gld: " +
			"no such file or directory"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := Open(tspy, "testdata/does_not_exist.gld")

		// --- Then ---
		affirm.Nil(t, have)
	})
}

func Test_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		pth := filepath.Join(t.TempDir(), "empty.gld")

		// --- When ---
		have := Create(tspy, pth)

		// --- Then ---
		affirm.NotNil(t, have)
		affirm.Equal(t, false, tspy.Failed())
		affirm.NotNil(t, have.comment)
		affirm.Equal(t, true, core.Same(tspy, have.t))

		content, err := os.ReadFile(pth)
		affirm.Nil(t, err)
		affirm.Equal(t, 0, len(content))
	})

	t.Run("error - cannot create the file", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		pth := filepath.Join(t.TempDir(), "not-existing-dir", "empty.gld")
		wMsg := "error creating file: open " + pth +
			": no such file or directory"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := Create(tspy, pth)

		// --- Then ---
		affirm.Nil(t, have)
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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "template: goldy:1:4: executing \"goldy\" at <.other>: " +
			"map has no entry for key \"other\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		data := WithData(map[string]any{"first": 1})
		gld := Open(tspy, "testdata/test_tpl.gld", data)

		// --- When ---
		have := gld.SetContent("32{{.other}}")

		// --- Then ---
		affirm.Nil(t, have)
	})
}

func Test_Goldy_Save(t *testing.T) {
	t.Run("when not a template", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.Close()

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
		tspy := tester.New(t)
		tspy.ExpectError()
		pth := filepath.Join(t.TempDir(), "sub-dir", "/golden.gld")
		wMsg := "error writing golden file (" + pth + "): open " + pth +
			": no such file or directory"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		gld := &Goldy{
			pth:     pth,
			comment: "Comment 1.\nComment 2.",
			content: []byte("content 1\ncontent 2\n"),
			t:       tspy,
		}

		// --- When ---
		gld.Save()
	})
}

// Benchmarks for goldy I/O and templating hot paths.

func Benchmark_Goldy_Open_Small(b *testing.B) {
	dir := b.TempDir()
	pth := filepath.Join(dir, "small.gld")
	content := "header comment\n---\nline1\nline2\n"
	if err := os.WriteFile(pth, []byte(content), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g := Open(&testing.T{}, pth)
		_ = g.Bytes()
	}
}

func Benchmark_Goldy_Open_WithTemplate(b *testing.B) {
	dir := b.TempDir()
	pth := filepath.Join(dir, "tpl.gld")
	content := "data for {{.Name}}\n---\nHello {{.Name}}!\n"
	if err := os.WriteFile(pth, []byte(content), 0o644); err != nil {
		b.Fatal(err)
	}

	data := map[string]any{"Name": "World"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g := Open(&testing.T{}, pth, WithData(data))
		_ = g.Bytes()
	}
}

func Benchmark_Goldy_Save(b *testing.B) {
	dir := b.TempDir()
	pth := filepath.Join(dir, "out.gld")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g := &Goldy{
			pth:     pth,
			comment: "Benchmark comment\n",
			content: []byte("benchmark content line\n"),
			t:       &testing.T{},
		}
		g.Save()
	}
}
