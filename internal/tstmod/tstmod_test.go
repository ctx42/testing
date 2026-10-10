// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tstmod

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_New(t *testing.T) {
	t.Run("module builds", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectTempDir(1)
		tspy.Close()

		// --- When ---
		have := New(tspy, "v1")

		// --- Then ---
		affirm.Equal(t, 2, len(have.ExternalDirs))

		cmd := exec.Command("go", "list", "./...")
		cmd.Dir = have.Dir
		out := string(must.Value(cmd.Output()))
		want := "" +
			"github.com/ctx42/tst-project\n" +
			"github.com/ctx42/tst-project/pkg/mercury\n"
		affirm.Equal(t, want, out)
	})

	t.Run("error - unknown version", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.ExpectLogEqual("unknown test project version: v3")
		tspy.Close()

		// --- When ---
		have := affirm.Panic(t, func() { New(tspy, "v3") })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *have)
	})
}

func Test_Module_WriteFile(t *testing.T) {
	t.Run("creates parent directories", func(t *testing.T) {
		// --- Given ---
		mod := New(t, "v2")

		// --- When ---
		have := mod.WriteFile("pkg/new/file.go", "package new")

		// --- Then ---
		affirm.Equal(t, filepath.Join(mod.Dir, "pkg/new/file.go"), have)
		affirm.Equal(t, "package new", string(must.Value(os.ReadFile(have))))
	})

	t.Run("absolute path", func(t *testing.T) {
		// --- Given ---
		mod := New(t, "v2")
		pth := filepath.Join(t.TempDir(), "file.txt")

		// --- When ---
		have := mod.WriteFile(pth, "content")

		// --- Then ---
		affirm.Equal(t, pth, have)
		affirm.Equal(t, "content", string(must.Value(os.ReadFile(pth))))
	})

	t.Run("error - parent is a file", func(t *testing.T) {
		// --- Given ---
		mod := New(t, "v2")

		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()
		mod.t = tspy

		// --- When ---
		have := affirm.Panic(t, func() { mod.WriteFile("go.mod/x", "") })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *have)
	})
}

func Test_Module_CreateDir(t *testing.T) {
	t.Run("relative path", func(t *testing.T) {
		// --- Given ---
		mod := New(t, "v2")

		// --- When ---
		have := mod.CreateDir("pkg/dir")

		// --- Then ---
		affirm.Equal(t, filepath.Join(mod.Dir, "pkg/dir"), have)
		affirm.Equal(t, true, must.Value(os.Stat(have)).IsDir())
	})

	t.Run("absolute path", func(t *testing.T) {
		// --- Given ---
		mod := New(t, "v2")
		pth := filepath.Join(t.TempDir(), "dir")

		// --- When ---
		have := mod.CreateDir(pth)

		// --- Then ---
		affirm.Equal(t, pth, have)
		affirm.Equal(t, true, must.Value(os.Stat(pth)).IsDir())
	})

	t.Run("error - parent is a file", func(t *testing.T) {
		// --- Given ---
		mod := New(t, "v2")

		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()
		mod.t = tspy

		// --- When ---
		have := affirm.Panic(t, func() { mod.CreateDir("go.mod/x") })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *have)
	})
}

func Test_Module_Path(t *testing.T) {
	// --- Given ---
	mod := &Module{Dir: "/mod"}

	// --- When ---
	have := mod.Path("pkg", "file.go")

	// --- Then ---
	affirm.Equal(t, filepath.Join("/mod", "pkg", "file.go"), have)
}

func Test_goMod(t *testing.T) {
	// --- When ---
	have := goMod("v0.1.0")

	// --- Then ---
	want := "" +
		"module github.com/ctx42/tst-project\n\n" +
		"go 1.24.0\n\n" +
		"require github.com/ctx42/tst-b v0.1.0\n\n" +
		"require github.com/ctx42/tst-a v0.1.0 // indirect\n\n" +
		"replace github.com/ctx42/tst-a v0.1.0 => ../tst-a\n" +
		"replace github.com/ctx42/tst-b v0.1.0 => ../tst-b\n"
	affirm.Equal(t, want, have)
}

func Test_variant_goSum(t *testing.T) {
	// --- Given ---
	vrt := variant{
		version: "v1.0.0",
		sumA:    "h1:a",
		sumAMod: "h1:am",
		sumB:    "h1:b",
		sumBMod: "h1:bm",
	}

	// --- When ---
	have := vrt.goSum()

	// --- Then ---
	want := "" +
		"github.com/ctx42/tst-a v1.0.0 h1:a\n" +
		"github.com/ctx42/tst-a v1.0.0/go.mod h1:am\n" +
		"github.com/ctx42/tst-b v1.0.0 h1:b\n" +
		"github.com/ctx42/tst-b v1.0.0/go.mod h1:bm\n"
	affirm.Equal(t, want, have)
}
