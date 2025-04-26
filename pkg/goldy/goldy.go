// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package goldy

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ctx42/testing/internal/core"
)

// Text is a helper returning contents of a golden file at the given path. The
// contents start after the mandatory marker "---" line, anything before it is
// ignored. It's customary to have short documentation about golden file
// contents before the "marker".
func Text(t core.T, pth string) string {
	t.Helper()

	// Open the file
	fil, err := os.Open(pth)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer func() { _ = fil.Close() }()

	var started bool
	var lines []string

	rdr := bufio.NewReader(fil)
	for {
		line, err := rdr.ReadString('\n')
		eof := errors.Is(err, io.EOF)
		if err != nil && !eof {
			t.Fatalf("error reading file: %v", err)
			return ""
		}
		if !started {
			started = line == "---\n"
			if !started && eof {
				t.Fatal("golden file is missing \"---\" marker")
				return ""
			}
			continue
		}
		lines = append(lines, line)
		if eof {
			return strings.Join(lines, "")
		}
	}
}
