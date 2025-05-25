// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"strings"
)

// Indent indents lines with n number of runes. Lines are indented only if
// there are more than one line.
func Indent(n int, r rune, lns string) string {
	if lns == "" {
		return ""
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

// Pad left pads the string with spaces.
func Pad(str string, length int) string {
	l := len(str)
	if length > l {
		return strings.Repeat(" ", length-l) + str
	}
	return str
}
