// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_file_parseImports(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &file{
			path: filepath.Join(wd, "file.go"),
			ast: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Value: `"github.com/ctx42/testing/pkg/mock"`,
						},
					},
					{
						Path: &ast.BasicLit{
							Value: `"github.com/ctx42/testing/pkg/mocker"`,
						},
					},
				},
			},
		}

		// --- When ---
		fil.parseImports()

		// --- Then ---
		wPks := []*gopkg{
			{
				pkgName: "mock",
				pkgPath: "github.com/ctx42/testing/pkg/mock",
				wd:      wd,
			},
			{
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				wd:      wd,
			},
		}
		assert.Equal(t, wPks, fil.pks)
	})

	t.Run("import alias is recognized", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &file{
			path: filepath.Join(wd, "file.go"),
			ast: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Value: `"github.com/ctx42/testing/pkg/mocker"`,
						},
						Name: &ast.Ident{Name: "alias"},
					},
				},
			},
		}

		// --- When ---
		fil.parseImports()

		// --- Then ---
		wPks := []*gopkg{
			{
				alias:   "alias",
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				wd:      wd,
			},
		}
		assert.Equal(t, wPks, fil.pks)
	})

	t.Run("import dot alias is recognized", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &file{
			path: filepath.Join(wd, "file.go"),
			ast: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Value: `"github.com/ctx42/testing/pkg/mocker"`,
						},
						Name: &ast.Ident{Name: "."},
					},
				},
			},
		}

		// --- When ---
		fil.parseImports()

		// --- Then ---
		wPks := []*gopkg{
			{
				alias:   "",
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				wd:      wd,
			},
		}
		assert.Equal(t, wPks, fil.pks)
		assert.Equal(t, wPks, fil.dots)
	})

	t.Run("invalid syntax imports are skipped", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &file{
			path: filepath.Join(wd, "file.go"),
			ast: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Value: `"invalid syntax`,
						},
					},
				},
			},
		}

		// --- When ---
		fil.parseImports()

		// --- Then ---
		assert.Len(t, 0, fil.pks)
		assert.NotNil(t, fil.pks)
	})
}

func Test_file_parseDecls(t *testing.T) {
	t.Run("declarations extracted", func(t *testing.T) {
		// --- Given ---
		dec0 := &ast.GenDecl{Tok: token.TYPE}
		dec1 := &ast.GenDecl{Tok: token.VAR}
		fil := &file{
			ast: &ast.File{
				Decls: []ast.Decl{dec0, dec1},
			},
		}

		// --- When ---
		fil.parseDecls()

		// --- Then ---
		assert.Len(t, 1, fil.decls)
		assert.Same(t, fil.decls[0], dec0)
	})

	t.Run("parsing results are cached", func(t *testing.T) {
		// --- Given ---
		fil := &file{
			ast: &ast.File{
				Decls: []ast.Decl{
					&ast.GenDecl{Tok: token.VAR},
					&ast.GenDecl{Tok: token.TYPE},
				},
			},
		}

		fil.parseDecls()
		fil.ast.Decls = nil

		// --- When ---
		fil.parseDecls()

		// --- Then ---
		assert.Len(t, 1, fil.decls)
	})
}

func Test_file_findPackage(t *testing.T) {
	t.Run("find by name", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		fil := &file{
			pks: []*gopkg{
				{alias: "a0", pkgName: "n0", resolved: true},
				{pkgName: "n1", resolved: true},
				{alias: "a2", pkgName: "n2", resolved: true},
			},
		}

		// --- When ---
		have, err := fil.findPackage(res, "n1")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &gopkg{pkgName: "n1", resolved: true}, have)
	})

	t.Run("error - aliased import by package name", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		fil := &file{
			pks: []*gopkg{
				{alias: "a1", pkgName: "n1", resolved: true},
				{pkgName: "n2", resolved: true},
			},
		}

		// --- When ---
		have, err := fil.findPackage(res, "n1")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.Nil(t, have)
	})

	t.Run("find by alias", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		fil := &file{
			pks: []*gopkg{
				{alias: "a0", pkgName: "n0", resolved: true},
				{alias: "a1", pkgName: "n1", resolved: true},
				{alias: "a2", pkgName: "n2", resolved: true},
			},
		}

		// --- When ---
		have, err := fil.findPackage(res, "a1")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &gopkg{alias: "a1", pkgName: "n1", resolved: true}, have)
	})

	t.Run("skips unrelated import failing to resolve", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		wd := must.Value(os.Getwd())
		fil := &file{
			pks: []*gopkg{
				newPkg(wd, "example.com/does/not/exist"),
				{pkgName: "n1", resolved: true},
			},
		}

		// --- When ---
		have, err := fil.findPackage(res, "n1")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &gopkg{pkgName: "n1", resolved: true}, have)
	})

	t.Run("error - wanted import fails to resolve", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		wd := must.Value(os.Getwd())
		fil := &file{
			pks: []*gopkg{
				{pkgName: "n1", resolved: true},
				newPkg(wd, "example.com/does/not/exist"),
			},
		}

		// --- When ---
		have, err := fil.findPackage(res, "exist")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.ErrorContain(t, "example.com/does/not/exist", err)
		assert.Nil(t, have)
	})

	t.Run("error - empty packages", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		fil := &file{ast: &ast.File{}}

		// --- When ---
		have, err := fil.findPackage(res, "name")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
		assert.Nil(t, have)
	})

	t.Run("error - not found", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		fil := &file{
			pks: []*gopkg{
				{pkgName: "xyz", resolved: true},
			},
		}

		// --- When ---
		have, err := fil.findPackage(res, "abc")

		// --- Then ---
		assert.Equal(t, true, err == ErrUnkPkg) //nolint:errorlint
		assert.Nil(t, have)
	})
}

func Test_file_dotImports(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		// --- Given ---
		fil := &file{
			pks: []*gopkg{
				{alias: "a1", pkgName: "n1"},
			},
			dots: []*gopkg{
				{alias: "", pkgName: "n0"},
				{alias: "", pkgName: "n2"},
			},
		}

		// --- When ---
		have := fil.dotImports()

		// --- Then ---
		want := []*gopkg{
			{alias: "", pkgName: "n0"},
			{alias: "", pkgName: "n2"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("no dot imports", func(t *testing.T) {
		// --- Given ---
		fil := &file{ast: &ast.File{}}

		// --- When ---
		have := fil.dotImports()

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_file_findType(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		// --- Given ---
		typ0 := &ast.TypeSpec{Name: ast.NewIdent("MyType0")}
		typ1 := &ast.TypeSpec{Name: ast.NewIdent("MyType1")}

		fil := &file{
			decls: []*ast.GenDecl{
				{
					Tok:   token.TYPE,
					Specs: []ast.Spec{typ0},
				},
				{
					Tok:   token.TYPE,
					Specs: []ast.Spec{&ast.ValueSpec{}},
				},
				{
					Tok:   token.TYPE,
					Specs: []ast.Spec{typ1},
				},
			},
		}

		// --- When ---
		have := fil.findType("MyType1")

		// --- Then ---
		assert.Same(t, typ1, have)
	})

	t.Run("not found", func(t *testing.T) {
		// --- Given ---
		typ0 := &ast.TypeSpec{Name: ast.NewIdent("MyType0")}
		typ1 := &ast.TypeSpec{Name: ast.NewIdent("MyType1")}

		fil := &file{
			decls: []*ast.GenDecl{
				{
					Tok:   token.TYPE,
					Specs: []ast.Spec{typ0},
				},
				{
					Tok:   token.TYPE,
					Specs: []ast.Spec{&ast.ValueSpec{}},
				},
				{
					Tok:   token.TYPE,
					Specs: []ast.Spec{typ1},
				},
			},
		}

		// --- When ---
		have := fil.findType("MyType2")

		// --- Then ---
		assert.Nil(t, have)
	})
}
