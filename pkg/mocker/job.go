package mocker

import (
	"go/ast"
)

type Job struct {
	// TODO(rz): do we need it?
	// TODO(rz): test this.
	Action *Action
	Pkg    *Package
	File   *ast.File
	Itf    *ast.InterfaceType
}
