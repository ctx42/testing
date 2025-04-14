// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_specDir(t *testing.T) {
	t.Run("success for the current module", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		dir := filepath.Clean(filepath.Join(wd, "..", "assert"))

		// --- When ---
		have, err := specDir("", "github.com/ctx42/testing/pkg/assert")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, dir, have)
	})

	t.Run("success for an external module", func(t *testing.T) {
		// --- Given ---
		modDir := createTestModule(t)

		// --- When ---
		have, err := specDir(modDir, tstModImp)

		// --- Then ---
		assert.NoError(t, err)
		mod := os.Getenv("GOMODCACHE")
		if mod == "" {
			mod = filepath.Join(os.Getenv("GOPATH"), "pkg/mod")
		}
		assert.Equal(t, filepath.Join(mod, tstModImpCached), have)
	})

	t.Run("path without source code", func(t *testing.T) {
		// --- When ---
		have, err := specDir("", "github.com/ctx42/testing")

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.Empty(t, have)
	})

	t.Run("invalid import path", func(t *testing.T) {
		// --- When ---
		have, err := specDir("", "!!!")

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.Empty(t, have)
	})
}

func Test_dirToSpec(t *testing.T) {
	t.Run("success for the current module", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		dir := filepath.Clean(filepath.Join(wd, "..", "assert"))

		// --- When ---
		have, err := dirToSpec(dir)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "github.com/ctx42/testing/pkg/assert", have)
	})

	t.Run("success for an external module", func(t *testing.T) {
		// --- Given ---
		modDir := createTestModule(t)

		// --- When ---
		have, err := dirToSpec(modDir)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "github.com/ctx42/tst-project", have)
	})

	t.Run("path without source code", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		dir := filepath.Clean(filepath.Join(wd, "testdata", "empty"))

		// --- When ---
		have, err := dirToSpec(dir)

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.Empty(t, have)
	})

	t.Run("invalid import path", func(t *testing.T) {
		// --- When ---
		have, err := dirToSpec("!!!")

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.Empty(t, have)
	})
}

func Test_settleImport_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())

	tt := []struct {
		testN string

		have Import
		want Import
	}{
		{
			"spec for the current package",
			Import{
				Spec: "github.com/ctx42/testing/pkg/mocker",
			},
			Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
		},
		{
			"spec for another package in the same module",
			Import{
				Spec: "github.com/ctx42/testing/pkg/assert",
			},
			Import{
				Name: "assert",
				Spec: "github.com/ctx42/testing/pkg/assert",
				Dir:  filepath.Join(wd, "..", "assert"),
			},
		},
		{
			"spec for package with ignored main",
			Import{
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/ignored",
			},
			Import{
				Name: "ignored",
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/ignored",
				Dir:  filepath.Join(wd, "testdata/ignored"),
			},
		},
		{
			"spec for multi package",
			Import{
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/multi",
			},
			Import{
				Name: "multi",
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/multi",
				Dir:  filepath.Join(wd, "testdata/multi"),
			},
		},
		{
			"spec for external package",
			Import{
				Spec: "net/http",
			},
			Import{
				Name: "http",
				Spec: "net/http",
				Dir:  filepath.Join(os.Getenv("GOROOT"), "src/net/http"),
			},
		},
		{
			"dir for the current package",
			Import{
				Dir: wd,
			},
			Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
		},
		{
			"dir for another package in the same module",
			Import{
				Dir: filepath.Join(wd, "..", "assert"),
			},
			Import{
				Name: "assert",
				Spec: "github.com/ctx42/testing/pkg/assert",
				Dir:  filepath.Join(wd, "..", "assert"),
			},
		},
		{
			"dir for external package",
			Import{
				Dir: filepath.Join(os.Getenv("GOROOT"), "src/net/http"),
			},
			Import{
				Name: "http",
				Spec: "net/http",
				Dir:  filepath.Join(os.Getenv("GOROOT"), "src/net/http"),
			},
		},
		{
			"name is always set from spec",
			Import{
				Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Spec: "github.com/ctx42/testing/pkg/mocker",
			},
			Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
		},
		{
			"dir is always changed to an absolute path",
			Import{
				Dir: ".",
			},
			Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have, err := settleImport(tc.have)

			// --- Then ---
			assert.NoError(t, err)
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_settleImport(t *testing.T) {
	t.Run("error import is zero value", func(t *testing.T) {
		// --- Given ---
		imp := Import{}

		// --- When ---
		_, err := settleImport(imp)

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvImport)
	})

	t.Run("error not existing directory", func(t *testing.T) {
		// --- Given ---
		imp := Import{Dir: "not-existing"}

		// --- When ---
		_, err := settleImport(imp)

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
	})

	t.Run("error unknown or invalid spec", func(t *testing.T) {
		// --- Given ---
		imp := Import{Spec: "github.com/ctx42/testing/pkg/aaa"}

		// --- When ---
		_, err := settleImport(imp)

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
	})

	t.Run("test module", func(t *testing.T) {
		// --- Given ---
		modDir := createTestModule(t)
		imp := Import{Dir: modDir}

		// --- When ---
		have, err := settleImport(imp)

		// --- Then ---
		assert.NoError(t, err)
		want := Import{
			Name: "project",
			Spec: "github.com/ctx42/tst-project",
			Dir:  modDir,
		}
		assert.Equal(t, want, have)
	})
}

func Test_findSources(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		wd := must.Value(os.Getwd())
		have, err := findSources("testdata/ignored")

		// --- Then ---
		assert.NoError(t, err)
		want := []string{
			filepath.Join(wd, "testdata/ignored/pkg.go"),
			filepath.Join(wd, "testdata/ignored/with_main.go"),
		}
		assert.Equal(t, want, have)
	})

	t.Run("does not recurse", func(t *testing.T) {
		// --- When ---
		wd := must.Value(os.Getwd())
		have, err := findSources("testdata/pkga")

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(wd, "testdata/pkga/helper.go")}
		assert.Equal(t, want, have)
	})

	t.Run("does not return test files", func(t *testing.T) {
		// --- When ---
		wd := must.Value(os.Getwd())
		have, err := findSources("testdata/pkgb")

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(wd, "testdata/pkgb/helper.go")}
		assert.Equal(t, want, have)
	})

	t.Run("error when path is not a directory", func(t *testing.T) {
		// --- When ---
		wd := must.Value(os.Getwd())
		have, err := findSources("testdata/pkga/helper.go")

		// --- Then ---
		assert.ErrorContain(t, "not a directory", err)
		want := filepath.Join(wd, "testdata/pkga/helper.go")
		assert.ErrorContain(t, want, err)
		assert.Nil(t, have)
	})

	t.Run("error when directory doesnt exist", func(t *testing.T) {
		// --- When ---
		wd := must.Value(os.Getwd())
		have, err := findSources("testdata/not-existing")

		// --- Then ---
		assert.ErrorIs(t, err, fs.ErrNotExist)
		want := filepath.Join(wd, "testdata/not-existing")
		assert.ErrorContain(t, want, err)
		assert.Nil(t, have)
	})
}

func Test_assumedPackageName_tabular(t *testing.T) {
	tt := []struct {
		testN string

		path string
		want string
	}{
		{"stdlib", "fmt", "fmt"},
		{"stdlib sub", "go/ast", "ast"},
		{"project", "github.com/user/project", "project"},
		{"with hyphen", "github.com/user/go-project", "project"},
		{"with underscore", "github.com/user/go_project", "go_project"},
		{"with a major version", "github.com/user/project/v2", "project"},
		{"with multiple hyphens", "github.com/user/go-project-abc", "abc"},
		{"with tilde", "github.com/user/tst~project", "tst"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := assumedPackageName(tc.path)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_toLowerSnakeCase_tabular(t *testing.T) {
	tt := []struct {
		in   string
		want string
	}{
		{"toLowerSnakeCase", "to_lower_snake_case"},
		{"TestItf", "test_itf"},
		{"test", "test"},
		{"TEST", "test"},
	}

	for _, tc := range tt {
		t.Run(tc.in, func(t *testing.T) {
			// --- When ---
			have := toLowerSnakeCase(tc.in)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_fmtCmdError(t *testing.T) {
	// --- Given ---
	out := " \nabc\t\ndef "

	// --- When ---
	have := fmtCmdError(out)

	// --- Then ---
	assert.Equal(t, "abc def", have)
}

func Test_astMethods(t *testing.T) {
	t.Run("methods found", func(t *testing.T) {
		// --- Given ---
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/embedded"
		_, astItf := findItf(t, "EmbedLocal", impSpec)

		// --- When ---
		have, err := astMethods(astItf)

		// --- Then ---
		assert.NoError(t, err)
		assert.Len(t, 2, have)
	})

	t.Run("no methods found", func(t *testing.T) {
		// --- Given ---
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/cases"
		_, astItf := findItf(t, "Empty", impSpec)

		// --- When ---
		have, err := astMethods(astItf)

		// --- Then ---
		assert.ErrorIs(t, err, ErrNoMethods)
		assert.Nil(t, have)
	})
}
