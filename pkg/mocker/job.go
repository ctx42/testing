package mocker

import (
	"go/ast"
)

type Job struct {
	// TODO(rz): do we need it?
	// TODO(rz): test this.
	Action  *Action
	SrcPkg  *Package
	SrcFile *ast.File
	Itf     *ast.InterfaceType
}
