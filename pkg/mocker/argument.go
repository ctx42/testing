// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import (
	"fmt"
	"strings"
)

// argument represents method's argument or return value.
//
// Examples:
//
//	a int
//	a map[pkg.A1]string
//	a map[*pkg.A1]func(b pkg.B1) func(a pkg.A1) error
type argument struct {
	name string   // Name which may be empty.
	typ  string   // String representation of the argument's type.
	pks  []*gopkg // Packages needed by the types used by the argument type.
}

// genName returns the argument's name. If the argument is unnamed, blank, or
// collides with a local the generator emits, it returns a synthetic `_aN`
// name based on the argument index.
//
// Examples:
//
//	name
//	_a0
//	_a123
func (arg argument) genName(idx int) string {
	if arg.name != "" && arg.name != "_" && !isGenLocal(arg.name) {
		return arg.name
	}
	return fmt.Sprintf("_a%d", idx)
}

// isGenLocal reports whether name is reserved for mock method body locals
// (_mck, _args, _rets, _rN, …). User parameters with those names would
// produce invalid Go if kept as-is.
func isGenLocal(name string) bool {
	switch name {
	case "_mck", "_args", "_arg", "_rets", "_r", "_rFn", "_elem":
		return true
	}
	if len(name) < 3 || name[0] != '_' {
		return false
	}
	// Synthetic argument names (_a0) and return temps (_r0).
	if name[1] != 'a' && name[1] != 'r' {
		return false
	}
	for _, c := range name[2:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// genArg generates code for the argument. If the argument is unnamed, it
// generates its name based on its index in the method's argument list.
//
// Examples:
//
//	name int
//	_a0 int
//	_a123 int
func (arg argument) genArg(idx int) string {
	return arg.genName(idx) + " " + arg.typ
}

// getType returns argument's type.
//
// Examples:
//
//	error
//	int
//	...int
//	map[pkg.A1]string
//	map[*pkg.A1]func(b pkg.B1) func(a pkg.A1) error
func (arg argument) getType() string {
	return arg.typ
}

// isVariadic returns true if the argument is variadic.
func (arg argument) isVariadic() bool {
	return strings.HasPrefix(arg.typ, "...")
}
