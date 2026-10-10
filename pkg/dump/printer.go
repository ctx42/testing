// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"bytes"
	"strings"
)

// Printer builds the text of a dumped value. Build it with [NewPrinter]; the
// zero value is not usable.
type Printer struct {
	dmp Dump
	buf *strings.Builder
}

// NewPrinter returns a new [Printer] configured by [Dump].
func NewPrinter(dmp Dump) Printer {
	return Printer{dmp: dmp, buf: &strings.Builder{}}
}

// NLI prints a new line when not flat and there is at least one entry.
func (prn Printer) NLI(cnt int) Printer {
	if !prn.dmp.Flat && cnt > 0 {
		_ = prn.buf.WriteByte('\n')
	}
	return prn
}

// NL prints a new line when not flat.
func (prn Printer) NL() Printer {
	if !prn.dmp.Flat {
		_ = prn.buf.WriteByte('\n')
	}
	return prn
}

// Comma prints a comma unless the output is flat and the entry is the last.
func (prn Printer) Comma(last bool) Printer {
	if !prn.dmp.Flat || !last {
		_ = prn.buf.WriteByte(',')
	}
	return prn
}

// Tab prints n levels of indentation, each [Dump.TabWidth] spaces wide, when
// not flat.
func (prn Printer) Tab(n int) Printer {
	if prn.dmp.Flat {
		return prn
	}
	if width := n * prn.dmp.TabWidth; width > 0 {
		_, _ = prn.buf.Write(bytes.Repeat([]byte{' '}, width))
	}
	return prn
}

// Space writes a space when not compact.
func (prn Printer) Space() Printer {
	if !prn.dmp.Compact {
		_ = prn.buf.WriteByte(' ')
	}
	return prn
}

// Sep writes a separator space between flat entries when not compact.
func (prn Printer) Sep(last bool) Printer {
	if !prn.dmp.Compact && !last && prn.dmp.Flat {
		_ = prn.buf.WriteByte(' ')
	}
	return prn
}

// Write writes a string to the builder.
func (prn Printer) Write(s string) Printer {
	_, _ = prn.buf.WriteString(s)
	return prn
}

// String returns the built string.
func (prn Printer) String() string { return prn.buf.String() }
