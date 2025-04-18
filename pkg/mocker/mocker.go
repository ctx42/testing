// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"errors"
	"fmt"
	"go/ast"
)

// Sentinel errors.
var (
	// ErrInvImport is returned when [Import] is invalid.
	ErrInvImport = errors.New("invalid import instance")

	// ErrInvSpec is returned when an import spec is invalid.
	//
	// This error occurs in the following cases:
	//   - The specified directory is empty or contains no source files.
	//   - The provided path cannot be resolved to a valid import spec.
	//   - The import spec cannot be resolved to a valid directory.
	ErrInvSpec = errors.New("invalid import spec")

	// ErrAstParse is returned when AST parsing encounters an error.
	ErrAstParse = errors.New("error parsing sources")

	// ErrItfNotFound is returned when an interface cannot be found.
	ErrItfNotFound = errors.New("interface not found")

	// ErrNoMethods is returned when interface we want to mock has no methods.
	ErrNoMethods = errors.New("interface has no methods")
)

// Mocker represents the main structure for creating interface mocks.
type Mocker struct {
	packages []*Package
}

func NewMocker() *Mocker {
	// TODO(rz): test this.
	// TODO(rz): document this.
	return &Mocker{}
}

// Run runs given mocking action.
func (mck *Mocker) Run(act *Action) error {
	pkg, err := mck.findPackage(act.Src)
	if err != nil {
		return err
	}

	astFil, astItf := pkg.find(act.SrcName)
	if astFil == nil || astItf == nil {
		return ErrItfNotFound
	}
	job := Job{Action: act, Pkg: pkg, File: astFil, Itf: astItf}
	if err = mck.parse(job); err != nil {
		return err
	}
	return nil
}

func (mck *Mocker) findPackage(imp Import) (*Package, error) {
	for _, pkg := range mck.packages {
		if pkg.imp.Spec == imp.Spec {
			return pkg, nil
		}
	}
	pkg, err := NewPackage(imp)
	if err != nil {
		return nil, err
	}
	mck.packages = append(mck.packages, pkg)
	return pkg, nil
}

func (mck *Mocker) parse(job Job) error {
	if job.Itf.Methods == nil || len(job.Itf.Methods.List) == 0 {
		return ErrNoMethods
	}
	mts, err := mck.methods(job)
	if err != nil {
		return err
	}
	_ = mts

	return nil
}

func (mck *Mocker) methods(job Job) ([]*Method, error) {
	fls := job.Itf.Methods.List
	mts := make([]*Method, 0, len(fls))
	for _, fld := range fls {
		metFld, err := mck.parseMethodField(job, fld)
		if err != nil {
			return nil, err
		}
		for _, met := range metFld {
			// Do not add duplicates. Note that we do not check if
			// the signatures are the same, we relay on the fact that
			// the code is correct.
			if found := findMethod(mts, met.name); found != nil {
				continue
			}
			mts = append(mts, met)
		}
	}
	return mts, nil
}

func (mck *Mocker) parseMethodField(job Job, fld *ast.Field) ([]*Method, error) {
	switch v := fld.Type.(type) {
	case *ast.FuncType:
		met, err := mck.parseMethod(job, 0, v)
		if err != nil {
			return nil, err
		}
		met.name = fld.Names[0].Name
		return []*Method{met}, nil

	case *ast.Ident:
		if v.Obj != nil {
			fil, itf := job.Pkg.find(v.Obj.Name)
			if fil == nil || itf == nil {
				return nil, ErrItfNotFound
			}
			// TODO(rz):
			// emb := mck.copy()
			// emb.srcItfName = v.Obj.Name
			// emb.srcItf = itf
			// emb.srcFil = fil
			// emb.srcFilImps = parseImports(mck.srcFil.Imports)
			// return emb.methods()
			return nil, nil
		}

	case *ast.SelectorExpr:
		// TODO(rz):
		// pkg := v.X.(*ast.Ident).Name // nolint: forcetypeassert
		// name := v.Sel.Name
		// imp, err := mck.findImport(pkg)
		// if err != nil {
		// 	return nil, err
		// }
		// opts := []Option{
		// 	WithSrcSpec(imp.spec),
		// 	WithDstSpec(mck.dstImport.spec),
		// }
		// _ = opts
		// subMck, err := newMocker(name, opts...)
		// if err != nil {
		//     return nil, err
		// }
		// return subMck.methods()
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected method field type: %T", fld.Type)
}

func (mck *Mocker) parseMethod(job Job, lvl int, fn *ast.FuncType) (*Method, error) {
	var err error
	met := &Method{}
	if fn.Params != nil {
		met.args, err = mck.parseArgs(job, lvl, fn.Params.List)
		if err != nil {
			return nil, err
		}
	}
	if fn.Results != nil {
		met.rets, err = mck.parseArgs(job, lvl, fn.Results.List)
		if err != nil {
			return nil, err
		}
	}
	return met, nil
}

// parseArgs parses slice of [ast.Field] instances as method arguments.
func (mck *Mocker) parseArgs(job Job, lvl int, fields []*ast.Field) ([]Argument, error) {
	var args []Argument
	for _, fld := range fields {
		idents := fld.Names
		if len(idents) == 0 {
			idents = []*ast.Ident{{}}
		}
		for _, ident := range idents {
			exp, err := mck.parseExpr(job, lvl, fld.Type)
			if err != nil {
				return nil, err
			}
			arg := Argument{
				// TODO(rz):
				// level:   lvl,
				Name:    ident.Name,
				Type:    exp.value,
				Imports: exp.imports,
			}
			args = append(args, arg)
		}
	}
	return args, nil
}

// parseExpr parses code expression.
//
// nolint: gocognit, cyclop
func (mck *Mocker) parseExpr(job Job, lvl int, e ast.Expr) (Expression, error) {
	switch v := e.(type) {

	case *ast.Ident:
		// In our context an identifier has a non-nil [ast.Object] associated
		// with it when the method parameter is a local type.
		if v.Obj != nil {
			// When the generated mock is put in a different package than the
			// mocked interface we need to provide qualified expression and
			// add the source package to the imports needed by the parameter.

			// TODO(rz):
			if v.Obj.Kind == ast.Typ &&
				job.Pkg.imp.Spec != job.Action.Dst.Spec {

				// TODO(rz): why do we set empty alias?
				imp := job.Pkg.imp.SetAlias("")
				exp := Expression{}
				exp.value = imp.Name + "." + v.Name
				exp.imports = append(exp.imports, imp)
				return exp, nil
			}

			exp := Expression{value: v.Name}
			return exp, nil
		}

		// TODO(rz):
		// Local type from different file.
		// if imp, err := mck.findLocalType(v.Name); err == nil {
		// 	exp := Expression{}
		// 	if mck.srcPkg.spec != mck.dstImport.spec {
		// 		exp.value = imp.name + "." + v.Name
		// 		exp.imports = append(exp.imports, imp)
		// 	} else {
		// 		exp.value = v.Name
		// 	}
		// 	return exp, nil
		// }

		// If parameter type is not builtin or local type
		// it must have been imported by the dot-import.
		if !isBuiltinType(v.Name) {
			// TODO(rz):
			// imp, err := mck.findDotType(v.Name)
			// if err != nil {
			// 	return Expression{}, err
			// }
			// exp := Expression{}
			// exp.value = imp.name + "." + v.Name
			// exp.imports = append(exp.imports, imp)
			// return exp, nil
			return Expression{}, nil
		}

		exp := Expression{}
		exp.value = v.Name
		return exp, nil

	case *ast.StarExpr:
		got, err := mck.parseExpr(job, lvl, v.X)
		if err != nil {
			return Expression{}, err
		}

		exp := Expression{}
		exp.value += "*"
		exp.value += got.value
		exp.imports = append(exp.imports, got.imports...)
		return exp, nil

	case *ast.SelectorExpr:
		// TODO(rz):
		// pkg := v.X.(*ast.Ident).Name // nolint: forcetypeassert
		// imp, err := mck.findImport(pkg)
		// if err != nil {
		// 	return Expression{}, err
		// }
		//
		// exp := Expression{}
		// exp.value = pkg + "." + v.Sel.Name
		// exp.imports = append(exp.imports, imp)
		// return exp, nil
		return Expression{}, nil

	case *ast.MapType:
		got, err := mck.parseExpr(job, lvl, v.Key)
		if err != nil {
			return Expression{}, err
		}

		exp := Expression{}
		exp.value += "map[" + got.value + "]"
		exp.imports = append(exp.imports, got.imports...)

		got, err = mck.parseExpr(job, lvl, v.Value)
		if err != nil {
			return Expression{}, err
		}
		exp.value += got.value
		exp.imports = append(exp.imports, got.imports...)
		return exp, nil

	case *ast.ArrayType:
		got, err := mck.parseExpr(job, lvl, v.Elt)
		if err != nil {
			return Expression{}, err
		}
		sb := "[]"
		if v.Len != nil {
			sb = "[" + v.Len.(*ast.BasicLit).Value + "]" // nolint: forcetypeassert
		}
		exp := Expression{}
		exp.value = sb + got.value
		exp.imports = append(exp.imports, got.imports...)
		return exp, nil

	case *ast.Ellipsis:
		got, err := mck.parseExpr(job, lvl, v.Elt)
		if err != nil {
			return Expression{}, err
		}
		exp := Expression{}
		exp.value = "..." + got.value
		exp.imports = append(exp.imports, got.imports...)
		return exp, nil

	case *ast.FuncType:
		// met, err := mck.parseMethod(job, lvl+1, v)
		// if err != nil {
		// 	return Expression{}, err
		// }
		// exp := Expression{}
		// exp.value = met.genSigCode("")
		// exp.imports = append(exp.imports, met.imports()...)
		// return exp, nil
		return Expression{}, nil

	case *ast.ChanType:
		got, err := mck.parseExpr(job, lvl, v.Value)
		if err != nil {
			return Expression{}, err
		}
		ch := "chan"
		switch v.Dir {
		case ast.SEND:
			ch += "<-"
		case ast.RECV:
			ch = "<-" + ch
		}
		exp := Expression{}
		exp.value = ch + " " + got.value
		exp.imports = append(exp.imports, got.imports...)
		return exp, nil
	}
	return Expression{}, ErrAstParse
}
