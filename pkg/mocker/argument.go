// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

// genName returns the argument's name. If the argument is unnamed, it
// generates a name based on its index in the method's argument list.
//
// Examples:
//
//	name
//	_a0
//	_a123
func (arg argument) genName(idx int) string {
	if arg.name != "" && arg.name != "_" {
		return arg.name
	}
	return fmt.Sprintf("_a%d", idx)
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
