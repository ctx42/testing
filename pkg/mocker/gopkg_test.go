// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/tstmod"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_newPkg(t *testing.T) {
	t.Run("empty import path", func(t *testing.T) {
		// --- Given ---

		// --- When ---
		have := newPkg("/dir", "")

		// --- Then ---
		assert.Equal(t, &gopkg{wd: "/dir"}, have)
	})

	t.Run("with import path", func(t *testing.T) {
		// --- Given ---

		// --- When ---
		have := newPkg("/dir", "github.com/ctx42/testing/pkg/mocker")

		// --- Then ---
		want := &gopkg{
			pkgName: "mocker",
			pkgPath: "github.com/ctx42/testing/pkg/mocker",
			wd:      "/dir",
		}
		assert.Equal(t, want, have)
	})
}

func Test_gopkg_resolve(t *testing.T) {
	t.Run("caches calls", func(t *testing.T) {
		// --- Given ---
		pkg := &gopkg{resolved: true}

		// --- When ---
		err := pkg.resolve()

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &gopkg{resolved: true}, pkg)
	})

	t.Run("error getModInfo not existing package", func(t *testing.T) {
		// --- Given ---
		pkg := newPkg("testdata/not-existing", "")

		// --- When ---
		err := pkg.resolve()

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.False(t, pkg.resolved)
	})
}

func Test_gopkg_resolve_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())
	mod1 := tstmod.New(t, "v1")

	tt := []struct {
		testN string

		dir  string
		pth  string
		want *gopkg
	}{
		{
			"current package",
			wd,
			"",
			&gopkg{
				pkgName:  "mocker",
				pkgPath:  "github.com/ctx42/testing/pkg/mocker",
				pkgDir:   wd,
				modName:  "testing",
				modPath:  "github.com/ctx42/testing",
				modDir:   filepath.Join(wd, "../.."),
				wd:       wd,
				resolved: true,
			},
		},
		{
			"empty package from the v1 test module",
			filepath.Join(mod1.Dir, "pkg/empty"),
			"",
			&gopkg{
				pkgName:  "empty",
				pkgPath:  "github.com/ctx42/tst-project/pkg/empty",
				pkgDir:   filepath.Join(mod1.Dir, "pkg/empty"),
				modName:  "project",
				modPath:  "github.com/ctx42/tst-project",
				modDir:   mod1.Dir,
				wd:       filepath.Join(mod1.Dir, "pkg/empty"),
				resolved: true,
			},
		},
		{
			"wd is the dot with a matching import path",
			wd,
			"github.com/ctx42/testing/pkg/mocker",
			&gopkg{
				pkgName:  "mocker",
				pkgPath:  "github.com/ctx42/testing/pkg/mocker",
				pkgDir:   wd,
				modName:  "testing",
				modPath:  "github.com/ctx42/testing",
				modDir:   filepath.Join(wd, "../.."),
				wd:       wd,
				resolved: true,
			},
		},
		{
			"wd is the root of the v1 test module with " +
				"the import path for its empty package",
			mod1.Dir,
			"github.com/ctx42/tst-project/pkg/empty",
			&gopkg{
				pkgName:  "empty",
				pkgPath:  "github.com/ctx42/tst-project/pkg/empty",
				pkgDir:   filepath.Join(mod1.Dir, "pkg/empty"),
				modName:  "project",
				modPath:  "github.com/ctx42/tst-project",
				modDir:   mod1.Dir,
				wd:       mod1.Dir,
				resolved: true,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := newPkg(tc.dir, tc.pth)

			// --- When ---
			err := pkg.resolve()

			// --- Then ---
			assert.NoError(t, err)
			assert.Equal(t, tc.want, pkg)
			assert.True(t, pkg.resolved)
		})
	}
}

func Test_gopkg_getPkgInfo(t *testing.T) {
	t.Run("empty package from the v1 test module", func(t *testing.T) {
		// --- Given ---
		mod := tstmod.New(t, "v1")
		dir := filepath.Join(mod.Dir, "pkg/empty")
		pkg := &gopkg{wd: dir}

		// --- When ---
		err := pkg.getPkgInfo()

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
	})
}

func Test_gopkg_getPkgInfo_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())
	mod1 := tstmod.New(t, "v1")

	tt := []struct {
		testN string

		pkg  *gopkg
		want *gopkg
	}{
		{
			"only working directory set",
			&gopkg{wd: wd},
			&gopkg{
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				pkgDir:  wd,
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		},
		{
			"working directory set from the package directory",
			&gopkg{pkgDir: wd},
			&gopkg{
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				pkgDir:  wd,
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		},
		{
			"the working directory may be set to anywhere in a module",
			&gopkg{wd: wd, pkgDir: filepath.Join(wd, "testdata/cases")},
			&gopkg{
				pkgName: "cases",
				pkgPath: "github.com/ctx42/testing/pkg/mocker/testdata/cases",
				pkgDir:  filepath.Join(wd, "testdata/cases"),
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		},
		{
			"working directory and import path",
			&gopkg{
				wd:      wd,
				pkgPath: "github.com/ctx42/testing/pkg/mocker/testdata/cases",
			},
			&gopkg{
				pkgName: "cases",
				pkgPath: "github.com/ctx42/testing/pkg/mocker/testdata/cases",
				pkgDir:  filepath.Join(wd, "testdata/cases"),
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		},
		{
			"the v1 test module root directory",
			&gopkg{wd: mod1.Dir},
			&gopkg{
				pkgName: "project",
				pkgPath: "github.com/ctx42/tst-project",
				pkgDir:  mod1.Dir,
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"a package from the v1 test module",
			&gopkg{wd: filepath.Join(mod1.Dir, "pkg/mercury")},
			&gopkg{
				pkgName: "mercury",
				pkgPath: "github.com/ctx42/tst-project/pkg/mercury",
				pkgDir:  filepath.Join(mod1.Dir, "pkg/mercury"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      filepath.Join(mod1.Dir, "pkg/mercury"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := tc.pkg.getPkgInfo()

			// --- Then ---
			assert.NoError(t, err)
			assert.Equal(t, tc.want, tc.pkg)
		})
	}
}

func Test_gopkg_getModInfo(t *testing.T) {
	t.Run("error not existing directory", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		pkg := &gopkg{wd: filepath.Join(wd, "testdata/not-existing")}

		// --- When ---
		err := pkg.getModInfo()

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
	})

	t.Run("error directory is not part of a go module", func(t *testing.T) {
		// --- When ---
		pkg := &gopkg{wd: t.TempDir()}

		// --- When ---
		err := pkg.getModInfo()

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
	})

	t.Run("error unknown module", func(t *testing.T) {
		// --- When ---
		pkg := &gopkg{pkgPath: "example.com/test"}

		// --- When ---
		err := pkg.getModInfo()

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.ErrorContain(t, "example.com/test", err)
	})
}

func Test_gopkg_getModInfo_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())
	mod1 := tstmod.New(t, "v1")
	mod2 := tstmod.New(t, "v2")

	// TODO(rz):
	t.Logf("--> Environ:\n %s", strings.Join(os.Environ(), "\n"))

	tt := []struct {
		testN string

		pkg  *gopkg
		want *gopkg
	}{
		{
			"wd is the current working directory without an import path",
			&gopkg{wd: wd},
			&gopkg{
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				pkgDir:  wd,
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		},
		{
			"wd is the root of the v1 test module without an import path",
			&gopkg{wd: mod1.Dir},
			&gopkg{
				pkgName: "project",
				pkgPath: "github.com/ctx42/tst-project",
				pkgDir:  mod1.Dir,
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"wd is a package in the v1 test module without an import path",
			&gopkg{wd: filepath.Join(mod1.Dir, "pkg/mercury")},
			&gopkg{
				pkgName: "mercury",
				pkgPath: "github.com/ctx42/tst-project/pkg/mercury",
				pkgDir:  filepath.Join(mod1.Dir, "pkg/mercury"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      filepath.Join(mod1.Dir, "pkg/mercury"),
			},
		},
		{
			"wd is an empty directory in the v1 test module without " +
				"an import path",
			&gopkg{wd: filepath.Join(mod1.Dir, "pkg/empty")},
			&gopkg{
				pkgName: "empty",
				pkgPath: "github.com/ctx42/tst-project/pkg/empty",
				pkgDir:  filepath.Join(mod1.Dir, "pkg/empty"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      filepath.Join(mod1.Dir, "pkg/empty"),
			},
		},
		{
			"wd is an empty directory in the v1 test module with " +
				"a matching import path",
			&gopkg{
				wd:      filepath.Join(mod1.Dir, "pkg/empty"),
				pkgPath: "github.com/ctx42/tst-project/pkg/empty",
			},
			&gopkg{
				pkgName: "empty",
				pkgPath: "github.com/ctx42/tst-project/pkg/empty",
				pkgDir:  filepath.Join(mod1.Dir, "pkg/empty"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      filepath.Join(mod1.Dir, "pkg/empty"),
			},
		},
		{
			"wd is the root of the v1 test module with " +
				"an import path for its empty package",
			&gopkg{
				wd:      mod1.Dir,
				pkgPath: "github.com/ctx42/tst-project/pkg/empty",
			},
			&gopkg{
				pkgName: "empty",
				pkgPath: "github.com/ctx42/tst-project/pkg/empty",
				pkgDir:  filepath.Join(mod1.Dir, "pkg/empty"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"wd is the root of the v1 test module with " +
				"an import path for the root of an external module",
			&gopkg{
				wd:      mod1.Dir,
				pkgPath: "github.com/ctx42/tst-a",
			},
			&gopkg{
				pkgName: "a",
				pkgPath: "github.com/ctx42/tst-a",
				pkgDir:  modCache("github.com/ctx42/tst-a@v0.1.0"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"wd is the root of the v2 test module with " +
				"an import path for the root of an external module",
			&gopkg{
				wd:      mod2.Dir,
				pkgPath: "github.com/ctx42/tst-a",
			},
			&gopkg{
				pkgName: "a",
				pkgPath: "github.com/ctx42/tst-a",
				pkgDir:  modCache("github.com/ctx42/tst-a@v0.2.0"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod2.Dir,
				wd:      mod2.Dir,
			},
		},
		{
			"wd is the root of the v1 test module with " +
				"an import path for a package of an external module",
			&gopkg{
				wd:      mod1.Dir,
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
			},
			&gopkg{
				pkgName: "first",
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
				pkgDir:  modCache("github.com/ctx42/tst-b@v0.1.0/pkg/mocker/first"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"wd is the root of the v2 test module with " +
				"an import path for a package of an external module",
			&gopkg{
				wd:      mod2.Dir,
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
			},
			&gopkg{
				pkgName: "first",
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
				pkgDir:  modCache("github.com/ctx42/tst-b@v0.2.0/pkg/mocker/first"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod2.Dir,
				wd:      mod2.Dir,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := tc.pkg.getModInfo()

			// --- Then ---
			assert.NoError(t, err)
			assert.Equal(t, tc.want, tc.pkg)
		})
	}
}

func Test_gopkg_parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/ignored")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		err := pkg.parse()

		// --- Then ---
		assert.NoError(t, err)
		assert.Len(t, 2, pkg.files)
		assert.Equal(t, filepath.Join(dir, "pkg.go"), pkg.files[0].path)
		assert.Equal(t, filepath.Join(dir, "with_main.go"), pkg.files[1].path)
	})

	t.Run("parses files only once", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/ignored")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		must.Nil(pkg.parse())
		must.Nil(pkg.parse())

		// --- Then ---
		assert.Len(t, 2, pkg.files)
	})

	t.Run("error unknown directory", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/not-existing")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		err := pkg.parse()

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.ErrorContain(t, dir, err)
	})

	t.Run("error parsing ast", func(t *testing.T) {
		// --- Given ---
		mod := tstmod.New(t, "v1")
		pth := mod.WriteFile("invalid.go", "")
		pkg := &gopkg{pkgDir: mod.Dir}

		// --- When ---
		err := pkg.parse()

		// --- Then ---
		var e scanner.ErrorList
		assert.ErrorAs(t, &e, err)
		assert.ErrorContain(t, pth, err)
	})
}

func Test_gopkg_findType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hTyp, err := pkg.findType("Concrete")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "other.go"), hFil.path)
		assert.Equal(t, "Concrete", hTyp.Name.String())
	})

	t.Run("error parsing package", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/not-existing")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hTyp, err := pkg.findType("Unknown")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.ErrorContain(t, dir, err)
		assert.Nil(t, hFil)
		assert.Nil(t, hTyp)
	})

	t.Run("error unknown type", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hTyp, err := pkg.findType("Unknown")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkType, err)
		assert.ErrorContain(t, "Unknown", err)
		assert.Nil(t, hFil)
		assert.Nil(t, hTyp)
	})
}

func Test_gopkg_findItf(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hItf, err := pkg.findItf("Case00")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "cases.go"), hFil.path)
		assert.NotNil(t, hItf)
	})

	t.Run("error type alias from the same package", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hItf, err := pkg.findItf("ItfAliasLocal")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkItf, err)
		assert.Nil(t, hFil)
		assert.Nil(t, hItf)
	})

	t.Run("error type alias from the other package", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hItf, err := pkg.findItf("ItfAlias")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkItf, err)
		assert.Nil(t, hFil)
		assert.Nil(t, hItf)
	})

	t.Run("error parsing package", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/not-existing")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hItf, err := pkg.findItf("Unknown")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.ErrorContain(t, dir, err)
		assert.Nil(t, hFil)
		assert.Nil(t, hItf)
	})

	t.Run("error when type is not an interface", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hItf, err := pkg.findItf("Concrete")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkItf, err)
		assert.ErrorContain(t, "Concrete", err)
		assert.Nil(t, hFil)
		assert.Nil(t, hItf)
	})

	t.Run("error unknown type", func(t *testing.T) {
		// --- Given ---
		dir := filepath.Join(must.Value(os.Getwd()), "testdata/cases")
		pkg := &gopkg{pkgDir: dir}

		// --- When ---
		hFil, hItf, err := pkg.findItf("Unknown")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkType, err)
		assert.ErrorContain(t, "Unknown", err)
		assert.Nil(t, hFil)
		assert.Nil(t, hItf)
	})
}

func Test_gopkg_isValid(t *testing.T) {
	valild := func() *gopkg {
		return &gopkg{
			alias:   "",
			pkgName: "pkgName",
			pkgPath: "pkgPath",
			pkgDir:  "pkgDir",
			modName: "modName",
			modPath: "modPath",
			modDir:  "modDir",
			wd:      "wd",
		}
	}

	tt := []struct {
		testN string

		action func(pkg *gopkg)
		want   bool
	}{
		{"valid", func(pkg *gopkg) {}, true},
		{"no pkgName", func(pkg *gopkg) { pkg.pkgName = "" }, false},
		{"no pkgPath", func(pkg *gopkg) { pkg.pkgPath = "" }, false},
		{"no pkgDir", func(pkg *gopkg) { pkg.pkgDir = "" }, false},
		{"no modName", func(pkg *gopkg) { pkg.modName = "" }, false},
		{"no modPath", func(pkg *gopkg) { pkg.modPath = "" }, false},
		{"no modDir", func(pkg *gopkg) { pkg.modDir = "" }, false},
		{"no wd", func(pkg *gopkg) { pkg.wd = "" }, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := valild()
			tc.action(pkg)

			// --- When ---
			have := pkg.isValid()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_gopkg_id_tabular(t *testing.T) {
	tt := []struct {
		testN string

		pkgPath string
		pkgDir  string
		want    string
	}{
		{"pkgPath set", "pkgPath", "", "pkgPath"},
		{"pkgDir set", "", "pkgDir", "pkgDir"},
		{"invalid", "", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := &gopkg{pkgPath: tc.pkgPath, pkgDir: tc.pkgDir}

			// --- When ---
			have := pkg.id()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_gopkg_equal_tabular(t *testing.T) {
	tt := []struct {
		testN string

		path string
		dir  string
		want bool
	}{
		{"both empty", "", "", false},
		{"import path equal", "example.com/mod", "example.com/mod", true},
		{"dir path equal", "/dir", "/dir", true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := &gopkg{pkgPath: tc.path, pkgDir: tc.dir}

			// --- When ---
			have := pkg.equal(pkg)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_gopkg_from(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		files := make([]*file, 0)
		fset := &token.FileSet{}

		from := &gopkg{
			alias:    "alias",
			pkgName:  "pkgName",
			pkgPath:  "pkgPath",
			pkgDir:   "pkgDir",
			modName:  "modName",
			modPath:  "modPath",
			modDir:   "modDir",
			wd:       "wd",
			resolved: true,
			files:    files,
			fset:     fset,
		}
		to := &gopkg{alias: "my"}

		// --- When ---
		to.from(from)

		// --- Then ---
		want := &gopkg{
			alias:    "my",
			pkgName:  "pkgName",
			pkgPath:  "pkgPath",
			pkgDir:   "pkgDir",
			modName:  "modName",
			modPath:  "modPath",
			modDir:   "modDir",
			wd:       "wd",
			resolved: true,
			files:    files,
			fset:     fset,
		}
		assert.Equal(t, want, to)
		assert.Same(t, want.files, to.files)
		assert.Same(t, want.fset, to.fset)
		assert.Fields(t, 11, gopkg{})
	})
}

func Test_gopkg_genImport_tabular(t *testing.T) {
	tt := []struct {
		testN string

		alias string
		path  string
		want  string
	}{
		{"without alias", "", "path", `"path"`},
		{"with alias", "alias", "path", `alias "path"`},
		{"with dot alias", ".", "path", `"path"`},
		{"empty", "", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := &gopkg{alias: tc.alias, pkgPath: tc.path}

			// --- When ---
			have := pkg.genImport()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
