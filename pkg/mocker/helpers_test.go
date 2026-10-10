// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

func Test_genAppendFromTo(t *testing.T) {
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

	t.Run("add not existing", func(t *testing.T) {
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

	t.Run("add existing simple", func(t *testing.T) {
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

	t.Run("add existing many", func(t *testing.T) {
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
			{alias: "", pkgName: "fmt", pkgPath: "fmt"},
			{alias: "", pkgName: "os", pkgPath: "os"},
			{alias: "tm", pkgName: "time", pkgPath: "time"},
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

func Test_dedupeImports(t *testing.T) {
	t.Run("keeps first import of a path", func(t *testing.T) {
		// --- Given ---
		pkg0 := &gopkg{pkgName: "mock", pkgPath: "example.com/mock"}
		pkg1 := &gopkg{
			alias:   "mk",
			pkgName: "mk",
			pkgPath: "example.com/mock",
		}
		met := &method{args: []argument{{typ: "mk.Mock", pks: []*gopkg{pkg1}}}}
		imps := []*gopkg{pkg0, pkg1}
		mts := []*method{met}

		// --- When ---
		have := dedupeImports(imps, mts)

		// --- Then ---
		assert.Equal(t, []*gopkg{pkg0}, have)
		assert.Equal(t, "mock.Mock", met.args[0].typ)
	})

	t.Run("keeps imports without path", func(t *testing.T) {
		// --- Given ---
		imps := []*gopkg{{}, nil}

		// --- When ---
		have := dedupeImports(imps, nil)

		// --- Then ---
		assert.Len(t, 2, have)
	})
}

func Test_renameShadowingArgs(t *testing.T) {
	// --- Given ---
	pkg := &gopkg{pkgName: "url", pkgPath: "net/url"}
	met := &method{args: []argument{{name: "url"}, {name: "any"}, {name: "a"}}}
	imps := []*gopkg{pkg, nil}
	mts := []*method{met}

	// --- When ---
	renameShadowingArgs(imps, mts)

	// --- Then ---
	want := []argument{{name: ""}, {name: ""}, {name: "a"}}
	assert.Equal(t, want, met.args)
}

func Test_resolveImportNameClashes(t *testing.T) {
	t.Run("renames second package with same name", func(t *testing.T) {
		// --- Given ---
		a := &gopkg{pkgName: "foo", pkgPath: "example.com/a/foo"}
		b := &gopkg{pkgName: "foo", pkgPath: "example.com/b/foo"}
		imps := []*gopkg{a, b}
		met := &method{
			args: []argument{{
				typ: "foo.T",
				pks: []*gopkg{b},
			}},
			rets: []argument{{
				typ: "map[afoo.U]foo.V",
				pks: []*gopkg{b},
			}},
		}
		mts := []*method{met}

		// --- When ---
		resolveImportNameClashes(imps, mts)

		// --- Then ---
		assert.Equal(t, "foo", a.pkgName)
		assert.Equal(t, "foo2", b.pkgName)
		assert.Equal(t, "foo2", b.alias)
		assert.Equal(t, "foo2.T", met.args[0].typ)
		assert.Equal(t, "map[afoo.U]foo2.V", met.rets[0].typ)
	})

	t.Run("clash uses alias not package name", func(t *testing.T) {
		// --- Given ---
		a := &gopkg{alias: "mt", pkgName: "time", pkgPath: "time"}
		b := &gopkg{pkgName: "mt", pkgPath: "example.com/mt"}
		imps := []*gopkg{a, b}
		met := &method{
			args: []argument{{
				typ: "mt.Thing",
				pks: []*gopkg{b},
			}},
		}
		mts := []*method{met}

		// --- When ---
		resolveImportNameClashes(imps, mts)

		// --- Then ---
		assert.Equal(t, "mt", a.ident())
		assert.Equal(t, "mt2", b.pkgName)
		assert.Equal(t, "mt2", b.alias)
		assert.Equal(t, "mt2.Thing", met.args[0].typ)
	})

	t.Run("reserved mock first keeps mock name", func(t *testing.T) {
		// --- Given ---
		mck := &gopkg{
			pkgName: "mock",
			pkgPath: "github.com/ctx42/testing/pkg/mock",
		}
		user := &gopkg{pkgName: "mock", pkgPath: "example.com/other/mock"}
		imps := []*gopkg{mck, user}
		met := &method{
			args: []argument{{
				typ: "mock.Type",
				pks: []*gopkg{user},
			}},
		}
		mts := []*method{met}

		// --- When ---
		resolveImportNameClashes(imps, mts)

		// --- Then ---
		assert.Equal(t, "mock", mck.pkgName)
		assert.Equal(t, "", mck.alias)
		assert.Equal(t, "mock2", user.pkgName)
		assert.Equal(t, "mock2.Type", met.args[0].typ)
	})

	t.Run("no clash leaves names", func(t *testing.T) {
		// --- Given ---
		a := &gopkg{pkgName: "a", pkgPath: "example.com/a"}
		b := &gopkg{pkgName: "b", pkgPath: "example.com/b"}
		imps := []*gopkg{a, b}

		// --- When ---
		resolveImportNameClashes(imps, nil)

		// --- Then ---
		assert.Equal(t, "a", a.pkgName)
		assert.Equal(t, "b", b.pkgName)
		assert.Equal(t, "", a.alias)
		assert.Equal(t, "", b.alias)
	})
}

func Test_nextImportAlias(t *testing.T) {
	// --- Given ---
	used := map[string]*gopkg{"pkg": {}, "pkg2": {}}

	// --- When ---
	have := nextImportAlias("pkg", used)

	// --- Then ---
	assert.Equal(t, "pkg3", have)
}

func Test_rewriteTypePkgPrefix(t *testing.T) {
	// --- Given ---
	pkg := &gopkg{pkgName: "foo", pkgPath: "example.com/foo"}
	other := &gopkg{pkgName: "bar", pkgPath: "example.com/bar"}
	met := &method{
		args: []argument{
			{typ: "foo.T", pks: []*gopkg{pkg}},
			{typ: "foo.T", pks: []*gopkg{other}},
		},
		rets: []argument{{typ: "[]foo.T", pks: []*gopkg{pkg}}},
	}
	mts := []*method{met, nil}

	// --- When ---
	rewriteTypePkgPrefix(mts, pkg, "foo", "foo2")

	// --- Then ---
	assert.Equal(t, "foo2.T", met.args[0].typ)
	assert.Equal(t, "foo.T", met.args[1].typ)
	assert.Equal(t, "[]foo2.T", met.rets[0].typ)
}

func Test_replacePkgPrefix_tabular(t *testing.T) {
	tt := []struct {
		testN string
		typ   string
		old   string
		neu   string
		want  string
	}{
		{"simple", "foo.T", "foo", "foo2", "foo2.T"},
		{
			"suffix name preserved",
			"map[afoo.T]foo.U",
			"foo",
			"foo2",
			"map[afoo.T]foo2.U",
		},
		{"no match", "bar.T", "foo", "foo2", "bar.T"},
		{"start of string", "foo.T", "foo", "x", "x.T"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := replacePkgPrefix(tc.typ, tc.old, tc.neu)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_isIdentRune_tabular(t *testing.T) {
	tt := []struct {
		testN string

		r    rune
		want bool
	}{
		{"underscore", '_', true},
		{"letter", 'a', true},
		{"digit", '7', true},
		{"non-ASCII letter", 'ż', true},
		{"dot", '.', false},
		{"star", '*', false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := isIdentRune(tc.r)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_argUsesPkg(t *testing.T) {
	t.Run("uses", func(t *testing.T) {
		// --- Given ---
		pkg := &gopkg{pkgPath: "example.com/foo"}
		arg := &argument{pks: []*gopkg{{}, pkg}}

		// --- When ---
		have := argUsesPkg(arg, pkg)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("same import path", func(t *testing.T) {
		// --- Given ---
		arg := &argument{pks: []*gopkg{{pkgPath: "example.com/foo"}}}

		// --- When ---
		have := argUsesPkg(arg, &gopkg{pkgPath: "example.com/foo"})

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("does not use", func(t *testing.T) {
		// --- Given ---
		arg := &argument{pks: []*gopkg{nil, {pkgPath: "example.com/bar"}}}

		// --- When ---
		have := argUsesPkg(arg, &gopkg{pkgPath: "example.com/foo"})

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_toLowerSnakeCase_tabular(t *testing.T) {
	tt := []struct {
		testN string

		in   string
		want string
	}{
		{"camel case", "toLowerSnakeCase", "to_lower_snake_case"},
		{"pascal case", "TestItf", "test_itf"},
		{"lower case", "test", "test"},
		{"upper case", "TEST", "test"},
		{"non-ASCII", "ŻółwMock", "żółw_mock"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := toLowerSnakeCase(tc.in)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_findSources(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/ignored")

		// --- When ---
		have, err := findSources(dir)

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(dir, "pkg.go")}
		assert.Equal(t, want, have)
	})

	t.Run("does not recurse", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/pkga")

		// --- When ---
		have, err := findSources(dir)

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(dir, "helper.go")}
		assert.Equal(t, want, have)
	})

	t.Run("does not return test files", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/pkgb")

		// --- When ---
		have, err := findSources(dir)

		// --- Then ---
		assert.NoError(t, err)
		want := []string{filepath.Join(dir, "helper.go")}
		assert.Equal(t, want, have)
	})

	t.Run("error - when the path is not a directory", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/pkga/helper.go")

		// --- When ---
		have, err := findSources(dir)

		// --- Then ---
		assert.ErrorContain(t, "not a directory", err)
		assert.ErrorContain(t, dir, err)
		assert.Nil(t, have)
	})

	t.Run("error - when directory doesnt exist", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/not-existing")

		// --- When ---
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

func Test_detectDirOrImp_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())

	tt := []struct {
		testN string

		wd       string
		dirOrImp string
		wDir     string
		wImp     string
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
		{"empty import", wd, "", wd, ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			hDir, hImp := detectDirOrImp(tc.wd, tc.dirOrImp)

			// --- Then ---
			assert.Equal(t, tc.wDir, hDir)
			assert.Equal(t, tc.wImp, hImp)
		})
	}
}

func Test_isDir(t *testing.T) {
	t.Run("directory", func(t *testing.T) {
		// --- When ---
		have := isDir(t.TempDir())

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("file", func(t *testing.T) {
		// --- When ---
		have := isDir("helpers.go")

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("missing", func(t *testing.T) {
		// --- When ---
		have := isDir(filepath.Join(t.TempDir(), "missing"))

		// --- Then ---
		assert.False(t, have)
	})
}
