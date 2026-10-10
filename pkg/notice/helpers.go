// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package notice

import (
	"strings"
	"unicode/utf8"
)

// Indent prepends n copies of rune r to each non-empty line in lns.
// Single-line input is returned unchanged.
func Indent(n int, r rune, lns string) string {
	if lns == "" || n <= 0 {
		return lns
	}
	rows := strings.Split(lns, "\n")
	if len(rows) == 1 {
		return lns
	}
	for i, lin := range rows {
		var ind string
		if lin != "" {
			ind = strings.Repeat(string(r), n)
		}
		rows[i] = ind + lin
	}
	return strings.Join(rows, "\n")
}

// Pad left-pads str with spaces until it is at least length runes long.
func Pad(str string, length int) string {
	l := utf8.RuneCountInString(str)
	if length > l {
		return strings.Repeat(" ", length-l) + str
	}
	return str
}
