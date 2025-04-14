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
		assert.NotNil(t, have)
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
			{},
			{alias: ".", pkgName: "pkge", pkgPath: "bitbucket.org/pkge"},
			{alias: ".", pkgName: "pkgf", pkgPath: "bitbucket.org/pkgf"},
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
		want := goldy.Text(t, "testdata/imports.gld")
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
