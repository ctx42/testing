// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package goldy

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/ctx42/testing/internal/core"
)

// Marker denotes a separator between a comment and the content in a golden
// test file.
const Marker = "---\n"

// byteseq represents byte sequences.
type byteseq interface {
	~string | ~[]byte
}

// Goldy represents golden file.
type Goldy struct {
	Path    string // Path to the golden file.
	Comment string // Golden file comment.
	Content []byte // Golden file content after the marker.
	t       core.T // Test manager.
}

// Open instantiates [Golden] based on the provided path to the golden file.
// The contents start after the mandatory [Marker] line, anything before it is
// ignored. It's customary to have short documentation about golden file
// contents before the marker.
func Open(t core.T, pth string) *Goldy {
	t.Helper()

	// Open the file
	fil, err := os.Open(pth)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer func() { _ = fil.Close() }()

	mark := []byte(Marker)
	gld := &Goldy{
		Path:    pth,
		Content: make([]byte, 0, 4*1024),
		t:       t,
	}

	var started bool
	rdr := bufio.NewReader(fil)
	for {
		line, err := rdr.ReadBytes('\n')
		eof := errors.Is(err, io.EOF)
		if err != nil && !eof {
			t.Fatalf("error reading file: %v", err)
			return nil
		}
		if !started {
			started = bytes.Equal(line, mark)
			if !started {
				if eof {
					m := strings.TrimSpace(Marker)
					t.Fatalf("the golden file is missing the %q marker", m)
					return nil
				}
				gld.Comment += string(line)
			}
			continue
		}
		gld.Content = append(gld.Content, line...)
		if eof {
			break
		}
	}
	return gld
}

// New returns new instance of [Goldy].
func New[C byteseq](t core.T, pth, comment string, content C) *Goldy {
	t.Helper()
	return &Goldy{
		Path:    pth,
		Comment: comment,
		Content: append([]byte{}, content...),
		t:       t,
	}
}

// String implements [fmt.Stringer] interface and returns golden file content as
// string.
func (gld *Goldy) String() string { return string(gld.Content) }

// Bytes return clone of the [Goldy.Content].
func (gld *Goldy) Bytes() []byte { return slices.Clone(gld.Content) }

// Save saves the golden file to the [Goldy.Path].
func (gld *Goldy) Save() {
	gld.t.Helper()

	buf := &bytes.Buffer{}
	comment := gld.Comment
	if !strings.HasSuffix(comment, "\n") {
		comment += "\n"
	}
	buf.WriteString(comment)
	buf.WriteString(Marker)
	buf.Write(gld.Content)
	if err := os.WriteFile(gld.Path, buf.Bytes(), 0600); err != nil {
		gld.t.Fatalf("error writing golden file (%s): %v", gld.Path, err)
	}
}
