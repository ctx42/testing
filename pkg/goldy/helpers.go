// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package goldy

import (
	"bytes"
)

// crlfToLF replaces a trailing CRLF line ending with LF.
func crlfToLF(line []byte) []byte {
	if bytes.HasSuffix(line, []byte("\r\n")) {
		line = append(line[:len(line)-2], '\n')
	}
	return line
}
