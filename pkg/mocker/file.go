// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
)

// file represents a Go source file.
type file struct {
	// TODO(rz): is it absolute?
	path      string         // Path to the file.
	pks       []*gopkg       // Imports in the file.
	typeDecls []*ast.GenDecl // Type declarations in the file.
}

// newFile creates a file instance from an AST file node and its absolute path.
func newFile(path string, f *ast.File) (*file, error) {
	fil := &file{
		path: path,
		pks:  make([]*gopkg, 0, len(f.Imports)),
	}
	for _, is := range f.Imports {
		pth, err := strconv.Unquote(is.Path.Value)
		if err != nil {
			continue
		}
		pkg, err := newPkg(filepath.Dir(path), pth)
		if err != nil {
			return nil, err
		}
		if is.Name != nil {
			pkg.alias = is.Name.Name
		}
		fil.pks = append(fil.pks, pkg)
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		fil.typeDecls = append(fil.typeDecls, gen)
	}
	return fil, nil
}

// findPackage returns the package matching the given import path or alias.
// It returns an ErrUnkPkg error if no matching import is found.
func (fil *file) findPackage(pathOrAlias string) (*gopkg, error) {
	for _, pkg := range fil.pks {
		if pkg.pkgName == pathOrAlias || pkg.alias == pathOrAlias {
			return pkg, nil
		}
	}
	return nil, ErrUnkPkg
}

// dotImports returns packages imported to the file with dot import.
func (fil *file) dotImports() []*gopkg {
	var dots []*gopkg
	for _, pkg := range fil.pks {
		if pkg.isDot() {
			dots = append(dots, pkg)
		}
	}
	return dots
}

// findType returns the type declaration in the file with the given name. It
// searches the file's type declarations and returns the matching
// [ast.TypeSpec]. If no type with the given name is found, it returns nil.
func (fil *file) findType(name string) *ast.TypeSpec {
	for _, dec := range fil.typeDecls {
		for _, spec := range dec.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok {
				// Can never happen due to filter
				// in newFile constructor, but just in case.
				continue
			}
			if typ.Name.Name != name {
				continue
			}
			return typ
		}
	}
	return nil
}
