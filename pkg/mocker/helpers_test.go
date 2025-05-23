// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
	"github.com/ctx42/testing/pkg/must"
)

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
		{"empty", "", ""},
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

func Test_appendLoopCode(t *testing.T) {
	// --- When ---
	have := genAppendFromTo("_args", "name")

	// --- Then ---
	want := "" +
		"\tfor _, _elem := range name {\n" +
		"\t\t_args = append(_args, _elem)\n" +
		"\t}\n"
	assert.Equal(t, want, have)
}

func Test_addUniquePackage(t *testing.T) {
	t.Run("add to nil", func(t *testing.T) {
		// --- Given ---
		imp := &gopkg{pkgPath: "a0_path"}

		// --- When ---
		have := addUniquePackage(nil, imp)

		// --- Then ---
		want := []*gopkg{
			{pkgPath: "a0_path"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("add not exiting", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{pkgPath: "a0_path"},
		}
		imp := &gopkg{pkgPath: "a1_path"}

		// --- When ---
		have := addUniquePackage(imps, imp)

		// --- Then ---
		want := []*gopkg{
			{pkgPath: "a0_path"},
			{pkgPath: "a1_path"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("add exiting simple", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{pkgPath: "a0_path"},
		}
		imp := &gopkg{pkgPath: "a0_path"}

		// --- When ---
		have := addUniquePackage(imps, imp)

		// --- Then ---
		want := []*gopkg{
			{pkgPath: "a0_path"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("add exiting many", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{pkgPath: "a0_path"},
			{pkgPath: "a3_path"},
		}
		imp0 := &gopkg{pkgPath: "a0_path"}
		imp1 := &gopkg{pkgPath: "a1_path"}
		imp2 := &gopkg{pkgPath: "a2_path"}
		imp3 := &gopkg{pkgPath: "a3_path"}
		imp4 := &gopkg{pkgPath: "a4_path"}

		// --- When ---
		have := addUniquePackage(imps, imp0, imp1, imp2, imp3, imp4)

		// --- Then ---
		want := []*gopkg{
			{pkgPath: "a0_path"},
			{pkgPath: "a3_path"},
			{pkgPath: "a1_path"},
			{pkgPath: "a2_path"},
			{pkgPath: "a4_path"},
		}
		assert.Equal(t, want, have)
	})
}

func Test_sortImports(t *testing.T) {
	t.Run("nil imports", func(t *testing.T) {
		// --- Given ---
		var imps []*gopkg

		// --- When ---
		have := sortImports(imps)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("no imports", func(t *testing.T) {
		// --- Given ---
		imps := make([]*gopkg, 0)

		// --- When ---
		have := sortImports(imps)

		// --- Then ---
		assert.Len(t, 0, have)
		assert.Nil(t, have)
	})

	t.Run("simple names", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{alias: "", pkgName: "pkgd", pkgPath: "bitbucket.org/pkgd"},
			{alias: "", pkgName: "pkga", pkgPath: "bitbucket.org/pkga"},
			{alias: "", pkgName: "pkgc", pkgPath: "bitbucket.org/pkgc"},
			{alias: "", pkgName: "pkgb", pkgPath: "bitbucket.org/pkgb"},
			{alias: "", pkgName: "pkge", pkgPath: "bitbucket.org/pkge"},
		}

		// --- When ---
		have := sortImports(imps)

		// --- Then ---
		want := []*gopkg{
			{alias: "", pkgName: "pkga", pkgPath: "bitbucket.org/pkga"},
			{alias: "", pkgName: "pkgb", pkgPath: "bitbucket.org/pkgb"},
			{alias: "", pkgName: "pkgc", pkgPath: "bitbucket.org/pkgc"},
			{alias: "", pkgName: "pkgd", pkgPath: "bitbucket.org/pkgd"},
			{alias: "", pkgName: "pkge", pkgPath: "bitbucket.org/pkge"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("full", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{alias: "", pkgName: "fs", pkgPath: "io/fs"},
			{alias: "mt", pkgName: "time", pkgPath: "time"},
			{alias: "", pkgName: "fmt", pkgPath: "fmt"},
			{alias: "", pkgName: "pkgd", pkgPath: "bitbucket.org/pkgd"},
			{alias: "", pkgName: "pkga", pkgPath: "bitbucket.org/pkga"},
			{alias: "", pkgName: "pkgc", pkgPath: "bitbucket.org/pkgc"},
			{alias: ".", pkgName: "pkge", pkgPath: "bitbucket.org/pkge"},
			{alias: ".", pkgName: "pkgf", pkgPath: "bitbucket.org/pkgf"},
		}

		// --- When ---
		have := sortImports(imps)

		// --- Then ---
		want := []*gopkg{
			{alias: "", pkgName: "fmt", pkgPath: "fmt"},
			{alias: "", pkgName: "fs", pkgPath: "io/fs"},
			{alias: "mt", pkgName: "time", pkgPath: "time"},
			{},
			{alias: "", pkgName: "pkga", pkgPath: "bitbucket.org/pkga"},
			{alias: "", pkgName: "pkgc", pkgPath: "bitbucket.org/pkgc"},
			{alias: "", pkgName: "pkgd", pkgPath: "bitbucket.org/pkgd"},
			{alias: ".", pkgName: "pkge", pkgPath: "bitbucket.org/pkge"},
			{alias: ".", pkgName: "pkgf", pkgPath: "bitbucket.org/pkgf"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("no std lib imports", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{alias: "", pkgName: "pkgd", pkgPath: "bitbucket.org/pkgd"},
			{alias: ".", pkgName: "pkga", pkgPath: "bitbucket.org/pkga"},
			{alias: "", pkgName: "pkgc", pkgPath: "bitbucket.org/pkgc"},
		}

		// --- When ---
		have := sortImports(imps)

		// --- Then ---
		want := []*gopkg{
			{alias: ".", pkgName: "pkga", pkgPath: "bitbucket.org/pkga"},
			{alias: "", pkgName: "pkgc", pkgPath: "bitbucket.org/pkgc"},
			{alias: "", pkgName: "pkgd", pkgPath: "bitbucket.org/pkgd"},
		}
		assert.Equal(t, want, have)
	})
}

func Test_genImports(t *testing.T) {
	t.Run("no imports", func(t *testing.T) {
		// --- When ---
		have := genImports(nil)

		// --- Then ---
		assert.Empty(t, have)
	})

	t.Run("imports", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{
			{
				alias:   "",
				pkgName: "fmt",
				pkgPath: "fmt",
			},
			{
				alias:   "",
				pkgName: "os",
				pkgPath: "os",
			},
			{
				alias:   "tm",
				pkgName: "time",
				pkgPath: "time",
			},
			{
				alias:   ".",
				pkgName: "pkgc",
				pkgPath: "github.com/ctx42/testing/pkg/mocker/testdata/pkgc",
			},
			{
				alias:   "pkga",
				pkgName: "pkga",
				pkgPath: "github.com/ctx42/testing/pkg/mocker/testdata/pkga",
			},
			{
				alias:   "",
				pkgName: "pkgb",
				pkgPath: "github.com/ctx42/testing/pkg/mocker/testdata/pkgb",
			},
		}

		// --- When ---
		have := genImports(imps)

		// --- Then ---
		want := goldy.Open(t, "testdata/imports.gld").String()
		assert.Equal(t, want, have)
	})
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

func Test_findSources(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/ignored")
		have, err := findSources(dir)

		// --- Then ---
		assert.NoError(t, err)
		want := []string{
			filepath.Join(dir, "pkg.go"),
			filepath.Join(dir, "with_main.go"),
		}
		assert.Equal(t, want, have)
	})

	t.Run("does not recurse", func(t *testing.T) {
		// --- When ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/pkga")
		have, err := findSources(dir)

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(dir, "helper.go")}
		assert.Equal(t, want, have)
	})

	t.Run("does not return test files", func(t *testing.T) {
		// --- When ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/pkgb")
		have, err := findSources(dir)

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(dir, "helper.go")}
		assert.Equal(t, want, have)
	})

	t.Run("error - when the path is not a directory", func(t *testing.T) {
		// --- When ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/pkga/helper.go")
		have, err := findSources(dir)

		// --- Then ---
		assert.ErrorContain(t, "not a directory", err)
		assert.ErrorContain(t, dir, err)
		assert.Nil(t, have)
	})

	t.Run("error - when directory doesnt exist", func(t *testing.T) {
		// --- When ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/not-existing")
		have, err := findSources(dir)

		// --- Then ---
		assert.ErrorIs(t, fs.ErrNotExist, err)
		assert.ErrorContain(t, dir, err)
		assert.Nil(t, have)
	})
}

func Test_findMethod(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		// --- Given ---
		var mts []*method

		// --- When ---
		have := findMethod(mts, "name")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("found", func(t *testing.T) {
		// --- Given ---
		mts := []*method{
			{name: "name0"},
			{name: "name1"},
			{name: "name2"},
		}

		// --- When ---
		have := findMethod(mts, "name1")

		// --- Then ---
		assert.NotNil(t, have)
		assert.Equal(t, "name1", have.name)
	})

	t.Run("not found", func(t *testing.T) {
		// --- Given ---
		mts := []*method{
			{name: "name0"},
			{name: "name1"},
			{name: "name2"},
		}

		// --- When ---
		have := findMethod(mts, "name3")

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_isBuiltin_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want bool
	}{
		{"int", true},
		{"uint64", true},
		{"any", true},
		{"SomeType", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := isBuiltin(tc.testN)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_detectDitOrImp_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())

	tt := []struct {
		testN string

		wd       string
		dirOrImp string
		wantDir  string
		wantImp  string
	}{
		{
			"absolute path to a package",
			wd,
			filepath.Join(wd, "testdata/pkga"),
			filepath.Join(wd, "testdata/pkga"),
			"",
		},
		{
			"relative path to a package",
			wd,
			"testdata/pkga",
			filepath.Join(wd, "testdata/pkga"),
			"",
		},
		{
			"not existing directory absolute path",
			wd,
			filepath.Join(wd, "testdata/not-existing"),
			filepath.Join(wd, "testdata/not-existing"),
			"",
		},
		{
			"not existing directory relative path",
			wd,
			"testdata/not-existing",
			wd,
			"testdata/not-existing",
		},
		{
			"import",
			wd,
			"github.com/ctx42/testing/pkg/mocker/testdata/cases",
			wd,
			"github.com/ctx42/testing/pkg/mocker/testdata/cases",
		},
		{
			"empty import",
			wd,
			"",
			wd,
			"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---

			// --- When ---
			haveDir, haveImp := detectDirOrImp(tc.wd, tc.dirOrImp)

			// --- Then ---
			assert.Equal(t, tc.wantDir, haveDir)
			assert.Equal(t, tc.wantImp, haveImp)
		})
	}
}
