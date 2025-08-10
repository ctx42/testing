// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"fmt"
)

// Row represents [Notice] row.
type Row struct {
	Name   string
	Format string
	Args   []any
}

// NewRow is constructor function for [Row].
func NewRow(name, format string, args ...any) Row {
	return Row{Name: name, Format: format, Args: args}
}

// String returns formated string value.
func (r Row) String() string { return fmt.Sprintf(r.Format, r.Args...) }

// PadName left pads row name with spaces to be requested length.
func (r Row) PadName(length int) string {
	return Pad(r.Name, length)
}

type Position string

const (
	PositionBefore Position = "before"
	PositionAfter  Position = "after"
)

type PositionedRow struct {
	Row
	Position string
}
