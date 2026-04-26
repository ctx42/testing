// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package dump.
//
// nolint: forbidigo
package dump

import (
	"bytes"
	"strings"
)

// Printer represents code printer.
type Printer struct {
	dmp Dump
	buf *strings.Builder
}

// NewPrinter returns new [Printer] configured by [Dump].
func NewPrinter(dmp Dump) Printer {
	return Printer{dmp: dmp, buf: &strings.Builder{}}
}

// NLI prints new line when not flat and at least one entry.
func (prn Printer) NLI(cnt int) Printer {
	if !prn.dmp.Flat && cnt > 0 {
		prn.buf.WriteByte('\n')
	}
	return prn
}

// NL prints new line when not flat.
func (prn Printer) NL() Printer {
	if !prn.dmp.Flat {
		prn.buf.WriteByte('\n')
	}
	return prn
}

// Comma prints Comma when not flat and not last entry.
func (prn Printer) Comma(last bool) Printer {
	if !(prn.dmp.Flat && last) {
		prn.buf.WriteByte(',')
	}
	return prn
}

// Tab prints indentation with n spaces when not flat.
func (prn Printer) Tab(n int) Printer {
	if prn.dmp.Flat {
		return prn
	}
	if n < 0 {
		n = 0
	}
	prn.buf.Write(bytes.Repeat([]byte{' '}, n*prn.dmp.TabWidth))
	return prn
}

// Space writes a space when not compact.
func (prn Printer) Space() Printer {
	if !prn.dmp.Compact {
		prn.buf.WriteByte(' ')
	}
	return prn
}

// Sep writes separator space.
func (prn Printer) Sep(last bool) Printer {
	if !prn.dmp.Compact && !last && prn.dmp.Flat {
		prn.buf.WriteByte(' ')
	}
	return prn
}

// Write writes string to the builder.
func (prn Printer) Write(s string) Printer {
	prn.buf.WriteString(s)
	return prn
}

// String returns built string.
func (prn Printer) String() string { return prn.buf.String() }
