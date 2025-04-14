// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// expression represents a type expression and associated imports needed by it.
//
// Examples:
//
//	int
//	map[pkga.A1]string
//	map[*pkga.A1]func(b pkgb.B1) func(a pkga.A1) error
type expression struct {
	value string   // String representation of the expression.
	pks   []*gopkg // Package pks needed by the expression.
}
