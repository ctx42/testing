// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package notice

import (
	"fmt"
)

// Row is a single named context line in a [Notice].
// It holds a name (e.g. "want", "have") plus a format string and args
// that are rendered via [fmt.Sprintf] when [Row.String] is called. A row
// without args renders its format verbatim.
type Row struct {
	Name   string
	Format string
	Args   []any
}

// NewRow creates a [Row] for use with [Notice.AppendRow].
// The format and args are stored and rendered later with [fmt.Sprintf].
func NewRow(name, format string, args ...any) Row {
	return Row{Name: name, Format: format, Args: args}
}

// String returns the formatted value using [fmt.Sprintf] on Format and Args.
// Without Args, it returns Format unchanged.
func (row Row) String() string {
	if len(row.Args) == 0 {
		return row.Format
	}
	return fmt.Sprintf(row.Format, row.Args...)
}

// PadName left-pads the row Name with spaces to the given length.
func (row Row) PadName(length int) string {
	return Pad(row.Name, length)
}
