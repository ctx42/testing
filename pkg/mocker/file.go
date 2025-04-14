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
	path  string         // Absolute path to the file.
	pks   []*gopkg       // Imports in the file.
	dots  []*gopkg       // Dot imports in the file.
	ast   *ast.File      // AST for the file.
	decls []*ast.GenDecl // Type declarations in the file.
}

// parseImports parses file imports populating pks slice with instances for
// each of them. It parses imports only once, caching the result for the later
// calls. The packages are not validated for existence.
func (fil *file) parseImports() {
	if fil.pks != nil {
		return
	}
	pks := make([]*gopkg, 0, len(fil.ast.Imports))
	for _, is := range fil.ast.Imports {
		pth, err := strconv.Unquote(is.Path.Value)
		if err != nil {
			continue
		}
		pkg := newPkg(filepath.Dir(fil.path), pth)
		if is.Name != nil {
			if is.Name.Name == "." {
				fil.dots = append(fil.dots, pkg)
			} else {
				pkg.alias = is.Name.Name
			}
		}
		pks = append(pks, pkg)
	}
	fil.pks = pks
}

// parseDecls parses Go type declarations from the file and stores them in the
// decls field. It parses declarations only once, caching the result for the
// later calls.
func (fil *file) parseDecls() {
	if fil.decls != nil {
		return
	}
	fil.decls = make([]*ast.GenDecl, 0, len(fil.ast.Decls))
	for _, decl := range fil.ast.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		fil.decls = append(fil.decls, gen)
	}
}

// findPackage finds the package imported by the file matching the given import
// path or its alias. It returns an [ErrUnkPkg] error if no matching import is
// found.
func (fil *file) findPackage(res *resolver, pathOrAlias string) (*gopkg, error) {
	var err error
	fil.parseImports()
	for _, pkg := range fil.pks {
		if err = res.resolve(pkg); err != nil {
			return nil, err
		}
		if pkg.pkgName == pathOrAlias || pkg.alias == pathOrAlias {
			return pkg, nil
		}
	}
	return nil, ErrUnkPkg
}

// dotImports returns packages imported to the file with dot import. The
// returned packages may not be validated for existence, you may need to call
// gopkg.resolve method before using them based on the context you use them for.
func (fil *file) dotImports() []*gopkg {
	fil.parseImports()
	return fil.dots
}

// findType returns the type declaration in the file with the given name. It
// searches the file's type declarations and returns the matching
// [ast.TypeSpec]. If no type with the given name is found, it returns nil.
func (fil *file) findType(name string) *ast.TypeSpec {
	fil.parseDecls()
	for _, dec := range fil.decls {
		for _, spec := range dec.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok {
				// Can never happen due to filter in the [file.parseDecls],
				// but it would look strange if we did not handle ok.
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
