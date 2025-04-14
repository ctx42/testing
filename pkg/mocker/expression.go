// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// Expression represents a type Expression.
//
// Examples:
//
//	int
//	map[pkga.A1]string
//	map[*pkga.A1]func(b pkgb.B1) func(a pkga.A1) error
type Expression struct {
	// TODO(rz): test this.
	// TODO(rz): document this.
	value   string   // String representation of the expression.
	imports []Import // Imports needed by the expression.
}
