package mocker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/internal/tstmod"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_newPkg_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())
	mod1 := tstmod.New(t, "v1")

	tt := []struct {
		testN string

		dir  string
		pth  string
		want *gopkg
	}{
		{
			"current package as a dot path",
			".",
			"",
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
			"empty package from the v1 test module",
			filepath.Join(mod1.Dir, "pkg/empty"),
			"",
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
			"wd is the dot with matching import path",
			".",
			"github.com/ctx42/testing/pkg/mocker",
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
			"wd is the root of the v1 test module with the import path for its empty package",
			mod1.Dir,
			"github.com/ctx42/tst-project/pkg/empty",
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
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have, err := newPkg(tc.dir, tc.pth)

			// --- Then ---
			assert.NoError(t, err)
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_newPkg(t *testing.T) {
	t.Run("error not existing directory", func(t *testing.T) {
		// --- When ---
		have, err := newPkg("testdata/not-existing", "")

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvPkg)
		assert.Nil(t, have)
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
			"current package as an absolute path",
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

func Test_gopkg_getPkgInfo(t *testing.T) {
	t.Run("empty package from the v1 test module", func(t *testing.T) {
		// --- Given ---
		mod := tstmod.New(t, "v1")
		dir := filepath.Join(mod.Dir, "pkg/empty")
		pkg := &gopkg{wd: dir}

		// --- When ---
		err := pkg.getPkgInfo()

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvPkg)
	})
}

func Test_gopkg_getModInfo_tabular(t *testing.T) {
	wd := must.Value(os.Getwd())
	mod1 := tstmod.New(t, "v1")
	mod2 := tstmod.New(t, "v2")

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
			"wd is an empty directory in the v1 test module without an import path",
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
			"wd is an empty directory in the v1 test module with matching import path",
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
			"wd is the root of the v1 test module with an import path for its empty package",
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
			"wd is the root of the v1 test module with an import path for root of an external module",
			&gopkg{
				wd:      mod1.Dir,
				pkgPath: "github.com/ctx42/tst-a",
			},
			&gopkg{
				pkgName: "a",
				pkgPath: "github.com/ctx42/tst-a",
				pkgDir:  goModCache("github.com/ctx42/tst-a@v0.1.0"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"wd is the root of the v2 test module with an import path for root of an external module",
			&gopkg{
				wd:      mod2.Dir,
				pkgPath: "github.com/ctx42/tst-a",
			},
			&gopkg{
				pkgName: "a",
				pkgPath: "github.com/ctx42/tst-a",
				pkgDir:  goModCache("github.com/ctx42/tst-a@v0.2.0"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod2.Dir,
				wd:      mod2.Dir,
			},
		},
		{
			"wd is the root of the v1 test module with an import path for package of an external module",
			&gopkg{
				wd:      mod1.Dir,
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
			},
			&gopkg{
				pkgName: "first",
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
				pkgDir:  goModCache("github.com/ctx42/tst-b@v0.1.0/pkg/mocker/first"),
				modName: "project",
				modPath: "github.com/ctx42/tst-project",
				modDir:  mod1.Dir,
				wd:      mod1.Dir,
			},
		},
		{
			"wd is the root of the v2 test module with an import path for package of an external module",
			&gopkg{
				wd:      mod2.Dir,
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
			},
			&gopkg{
				pkgName: "first",
				pkgPath: "github.com/ctx42/tst-b/pkg/mocker/first",
				pkgDir:  goModCache("github.com/ctx42/tst-b@v0.2.0/pkg/mocker/first"),
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

func Test_gopkg_getModInfo(t *testing.T) {
	t.Run("error not existing directory", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		pkg := &gopkg{wd: filepath.Join(wd, "testdata/not-existing")}

		// --- When ---
		err := pkg.getModInfo()

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvPkg)
	})

	t.Run("error directory is not part of a go module", func(t *testing.T) {
		// --- When ---
		pkg := &gopkg{wd: t.TempDir()}

		// --- When ---
		err := pkg.getModInfo()

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvPkg)
	})

	t.Run("error unknown module", func(t *testing.T) {
		// --- When ---
		pkg := &gopkg{pkgPath: "example.com/test"}

		// --- When ---
		err := pkg.getModInfo()

		// --- Then ---
		assert.ErrorIs(t, err, ErrUnkPkg)
		assert.ErrorContain(t, "example.com/test", err)
	})
}

func Test_gopkg_isDot_tabular(t *testing.T) {
	tt := []struct {
		testN string

		alias string
		want  bool
	}{
		{"empty", "", false},
		{"not dot", "alias", false},
		{"dot", ".", true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := &gopkg{alias: tc.alias}

			// --- When ---
			have := pkg.isDot()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
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

func Test_gopkg_GoString_tabular(t *testing.T) {
	tt := []struct {
		testN string

		alias string
		path  string
		want  string
	}{
		{"without alias", "", "path", `"path"`},
		{"with alias", "alias", "path", `alias "path"`},
		{"with dot", ".", "path", `. "path"`},
		{"empty", "", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			pkg := &gopkg{alias: tc.alias, pkgPath: tc.path}

			// --- When ---
			have := pkg.GoString()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
