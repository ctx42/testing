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
		// --- When ---
		have := New(t, "v1")

		// --- Then ---
		cmd := exec.Command("go", "list", "./...")
		cmd.Dir = have.Dir
		out, err := cmd.CombinedOutput()
		affirm.Nil(t, err)
		affirm.Equal(t, true, len(out) > 0)
		affirm.Equal(t, 2, len(have.ExternalDirs))
	})

	t.Run("error - unknown version", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.ExpectLogEqual("unknown test project version: v3")
		tspy.Close()

		// --- When ---
		msg := affirm.Panic(t, func() { New(tspy, "v3") })

		// --- Then ---
		affirm.Equal(t, tester.FailNowMsg, *msg)
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
}
