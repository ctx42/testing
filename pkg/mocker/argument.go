// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// Argument represents method's argument or return value.
//
// Examples:
//
//	a int
//	a map[pkg.A1]string
//	a map[*pkg.A1]func(b pkg.B1) func(a pkg.A1) error
type Argument struct {
	Name    string   // Argument name may be empty.
	Type    string   // String representation of the argument's type.
	Imports []Import // Imports needed by the argument's type.
}
