// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package tstkit

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ctx42/testing/pkg/tester"
)

// Golden is a helper returning contents of a golden file at given path. The
// contents start after mandatory marker "---" line, anything before it is
// ignored. It's customary to have a short documentation about golden file
// contents before the "marker".
func Golden(t tester.T, pth string) string {
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
