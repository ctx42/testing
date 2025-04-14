// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_NewPackage(t *testing.T) {
	t.Run("from the current module", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		impSpec := "github.com/ctx42/testing/pkg/mocker"
		imp := must.Value(settleImport(NewImport(impSpec)))

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, imp, pkg.imp)
		wFiles := []string{
			filepath.Join(wd, "action.go"),
			filepath.Join(wd, "argument.go"),
			filepath.Join(wd, "helpers.go"),
			filepath.Join(wd, "import.go"),
			filepath.Join(wd, "interface.go"),
			filepath.Join(wd, "method.go"),
			filepath.Join(wd, "mocker.go"),
			filepath.Join(wd, "package.go"),
		}
		assert.Equal(t, wFiles, keys(pkg.files))
		assert.NotEmpty(t, pkg.files)
	})

	t.Run("from another package in the same module", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/cases"
		imp := must.Value(settleImport(NewImport(impSpec)))

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, imp, pkg.imp)
		wFiles := []string{
			filepath.Join(wd, "testdata/cases/cases.go"),
			filepath.Join(wd, "testdata/cases/cases_gen.go"),
			filepath.Join(wd, "testdata/cases/cases_on_gen.go"),
			filepath.Join(wd, "testdata/cases/main.go"),
			filepath.Join(wd, "testdata/cases/other.go"),
		}
		assert.Equal(t, wFiles, keys(pkg.files))
	})

	t.Run("from a test module", func(t *testing.T) {
		// --- Given ---
		modDir := createTestModule(t)
		imp := must.Value(settleImport(Import{Dir: modDir}))

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, imp, pkg.imp)
		wFiles := []string{filepath.Join(modDir, "project.go")}
		assert.Equal(t, wFiles, keys(pkg.files))
	})

	t.Run("single package with ignored main", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/ignored"
		imp := must.Value(settleImport(NewImport(impSpec)))

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.NoError(t, err)
		wFiles := []string{
			filepath.Join(wd, "testdata/ignored/pkg.go"),
			filepath.Join(wd, "testdata/ignored/with_main.go"),
		}
		assert.Equal(t, wFiles, keys(pkg.files))
	})

	t.Run("multiple ast packages error", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/multi"
		imp := must.Value(settleImport(NewImport(impSpec)))

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.NoError(t, err)
		wFiles := []string{
			filepath.Join(wd, "testdata/multi/makefile.go"),
			filepath.Join(wd, "testdata/multi/project.go"),
		}
		assert.Equal(t, wFiles, keys(pkg.files))
	})

	t.Run("error parsing file", func(t *testing.T) {
		// --- Given ---
		modDir := createTestModule(t)
		path := filepath.Join(modDir, "invalid.go")
		content := []byte("package project\n\nfunc")
		must.Nil(os.WriteFile(path, content, 0644))
		imp := must.Value(settleImport(Import{Dir: modDir}))

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.ErrorIs(t, err, ErrAstParse)
		assert.ErrorContain(t, path, err)
		assert.Nil(t, pkg)
	})

	t.Run("error not existing import directory", func(t *testing.T) {
		// --- Given ---
		imp := Import{Dir: "not/existing"}

		// --- When ---
		pkg, err := NewPackage(imp)

		// --- Then ---
		assert.ErrorContain(t, "no such file or directory", err)
		assert.ErrorContain(t, "not/existing", err)
		assert.Nil(t, pkg)
	})
}

func Test_Package_Find(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/embedded"
		imp := must.Value(settleImport(NewImport(impSpec)))
		pkg := must.Value(NewPackage(imp))

		// --- When ---
		astFile, astItf := pkg.find("EmbedLocal")

		// --- Then ---
		assert.NotNil(t, astFile)
		assert.NotNil(t, astItf)
	})

	t.Run("error interface not found", func(t *testing.T) {
		// --- Given ---
		impSpec := "github.com/ctx42/testing/pkg/mocker/testdata/embedded"
		imp := must.Value(settleImport(NewImport(impSpec)))
		pkg := must.Value(NewPackage(imp))

		// --- When ---
		astFile, astItf := pkg.find("NotFound")

		// --- Then ---
		assert.Nil(t, astFile)
		assert.Nil(t, astItf)
	})
}
