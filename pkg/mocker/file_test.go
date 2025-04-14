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

func Test_newFile(t *testing.T) {
	t.Run("imports are parsed", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &ast.File{
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
		}

		// --- When ---
		have, err := newFile("file.go", fil)

		// --- Then ---
		assert.NoError(t, err)
		wPks := []*gopkg{
			{
				pkgName: "mock",
				pkgPath: "github.com/ctx42/testing/pkg/mock",
				pkgDir:  filepath.Join(wd, "../mock"),
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
			{
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				pkgDir:  wd,
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		}
		assert.Equal(t, wPks, have.pks)
	})

	t.Run("import alias is recognized", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &ast.File{
			Imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: `"github.com/ctx42/testing/pkg/mocker"`,
					},
					Name: &ast.Ident{Name: "alias"},
				},
			},
		}

		// --- When ---
		have, err := newFile("file.go", fil)

		// --- Then ---
		assert.NoError(t, err)
		wPks := []*gopkg{
			{
				alias:   "alias",
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				pkgDir:  wd,
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		}
		assert.Equal(t, wPks, have.pks)
	})

	t.Run("import dot alias is recognized", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		fil := &ast.File{
			Imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: `"github.com/ctx42/testing/pkg/mocker"`,
					},
					Name: &ast.Ident{Name: "."},
				},
			},
		}

		// --- When ---
		have, err := newFile("file.go", fil)

		// --- Then ---
		assert.NoError(t, err)
		wPks := []*gopkg{
			{
				alias:   ".",
				pkgName: "mocker",
				pkgPath: "github.com/ctx42/testing/pkg/mocker",
				pkgDir:  wd,
				modName: "testing",
				modPath: "github.com/ctx42/testing",
				modDir:  filepath.Join(wd, "../.."),
				wd:      wd,
			},
		}
		assert.Equal(t, wPks, have.pks)
	})

	t.Run("error when unknown package", func(t *testing.T) {
		// --- Given ---
		fil := &ast.File{
			Imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{Value: `"example.com"`},
				},
			},
		}

		// --- When ---
		have, err := newFile("file.go", fil)

		// --- Then ---
		assert.ErrorIs(t, err, ErrUnkPkg)
		assert.ErrorContain(t, "example.com", err)
		assert.Nil(t, have)
	})

	t.Run("invalid syntax imports are skipped", func(t *testing.T) {
		// --- Given ---
		fil := &ast.File{
			Imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: `"invalid syntax`,
					},
				},
			},
		}

		// --- When ---
		have, err := newFile("file.go", fil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Len(t, 0, have.pks)
	})

	t.Run("declarations extracted", func(t *testing.T) {
		// --- Given ---
		fil := &ast.File{
			Decls: []ast.Decl{
				&ast.GenDecl{Tok: token.TYPE},
				&ast.GenDecl{Tok: token.VAR},
			},
		}

		// --- When ---
		have, err := newFile("file.go", fil)

		// --- Then ---
		assert.NoError(t, err)
		assert.Len(t, 1, have.typeDecls)
		assert.Same(t, fil.Decls[0], have.typeDecls[0])
	})
}

func Test_file_findPackage(t *testing.T) {
	t.Run("find by name", func(t *testing.T) {
		// --- Given ---
		fil := file{
			pks: []*gopkg{
				{alias: "a0", pkgName: "n0"},
				{alias: "a1", pkgName: "n1"},
				{alias: "a2", pkgName: "n2"},
			},
		}

		// --- When ---
		have, err := fil.findPackage("n1")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &gopkg{alias: "a1", pkgName: "n1"}, have)
	})

	t.Run("find by alias", func(t *testing.T) {
		// --- Given ---
		fil := file{
			pks: []*gopkg{
				{alias: "a0", pkgName: "n0"},
				{alias: "a1", pkgName: "n1"},
				{alias: "a2", pkgName: "n2"},
			},
		}

		// --- When ---
		have, err := fil.findPackage("a1")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &gopkg{alias: "a1", pkgName: "n1"}, have)
	})

	t.Run("not found", func(t *testing.T) {
		// --- Given ---
		fil := file{}

		// --- When ---
		have, err := fil.findPackage("name")

		// --- Then ---
		assert.ErrorIs(t, err, ErrUnkPkg)
		assert.Nil(t, have)
	})
}

func Test_file_dotImports(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		// --- Given ---
		fil := file{
			pks: []*gopkg{
				{alias: ".", pkgName: "n0"},
				{alias: "a1", pkgName: "n1"},
				{alias: ".", pkgName: "n2"},
			},
		}

		// --- When ---
		have := fil.dotImports()

		// --- Then ---
		want := []*gopkg{
			{alias: ".", pkgName: "n0"},
			{alias: ".", pkgName: "n2"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("no dot imports", func(t *testing.T) {
		// --- Given ---
		fil := file{}

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
			typeDecls: []*ast.GenDecl{
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
			typeDecls: []*ast.GenDecl{
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
