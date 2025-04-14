// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Package represents a package.
type Package struct {
	imp   Import               // [Import] describing the package.
	files map[string]*ast.File // Files in the package.
	fset  *token.FileSet       // Manages sources and their positions.
}

// NewPackage retrieves information about a package based on the provided
// working directory and import specification. An empty working directory is
// the same as the current working directory.
//
// The working directory must be within a Go module that uses the specified
// package. This ensures the function can resolve the correct package version,
// as different modules on the same system may depend on different versions of
// the same package.
func NewPackage(imp Import) (*Package, error) {
	pkg := &Package{
		imp:  imp,
		fset: token.NewFileSet(),
	}
	if err := pkg.parse(); err != nil {
		return nil, err
	}
	return pkg, nil
}

// parse parses package files to their AST representation.
func (pkg *Package) parse() error {
	names, err := findSources(pkg.imp.Dir)
	if err != nil {
		return err
	}
	pkg.files = make(map[string]*ast.File, len(names))
	for _, name := range names {
		fil, err := parser.ParseFile(pkg.fset, name, nil, 0)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrAstParse, err)
		}
		pkg.files[name] = fil
	}
	return nil
}

// find returns the [ast.File] where the named interface was found along with
// its [ast.InterfaceType] representation. If an interface with the given name
// does not exist in the package, it returns [ErrItfNotFound].
func (pkg *Package) find(name string) (*ast.File, *ast.InterfaceType) {
	for _, fil := range pkg.files {
		for _, dec := range fil.Decls {
			gen, ok := dec.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue // Can never happen, but just in case.
				}
				itf, ok := typ.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				if typ.Name.Name != name {
					continue
				}
				return fil, itf
			}
		}
	}
	return nil, nil
}
