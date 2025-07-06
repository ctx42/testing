// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package goldy is designed to simplify reading content from golden files.
package goldy

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/ctx42/testing/internal/core"
)

// Marker is a separator between a golden file comment and the content.
const Marker = "---\n"

// WithData is the [Open] option setting [Goldy] data for golden files which
// are text templates.
func WithData(data map[string]any) func(*Goldy) {
	return func(gld *Goldy) { gld.data = data }
}

// Goldy represents a golden file.
//
// Example golden file:
//
//	Multi line
//	golden file documentation.
//	---
//	Content line #1.
//	Content line #2.
type Goldy struct {
	pth     string         // Path to the golden file.
	comment string         // Golden file comment.
	content []byte         // Golden file content.
	data    map[string]any // Template data.
	tpl     []byte         // Raw text template.
	t       core.T         // Test manager.
}

// Open instantiates [Goldy] based on the provided path to the golden file and
// options. The golden file content starts after the mandatory [Marker] line,
// anything before it is ignored. It's customary to have short golden file
// documentation before the marker.
func Open(t core.T, pth string, opts ...func(*Goldy)) *Goldy {
	t.Helper()

	fil, err := os.Open(pth)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer func() { _ = fil.Close() }()

	mark := []byte(Marker)
	gld := &Goldy{
		pth:     pth,
		content: make([]byte, 0, 4*1024),
		t:       t,
	}
	for _, opt := range opts {
		opt(gld)
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
				gld.comment += string(line)
			}
			continue
		}
		gld.content = append(gld.content, line...)
		if eof {
			break
		}
	}

	if gld.data != nil {
		gld.tpl = gld.content
		return gld.renderTemplate()
	}
	return gld
}

// String implements [fmt.Stringer] and returns golden file content.
func (gld *Goldy) String() string { return string(gld.content) }

// Bytes return clone of the golden file content.
func (gld *Goldy) Bytes() []byte { return slices.Clone(gld.content) }

// SetComment sets a comment for the golden file. Implements fluent interface.
func (gld *Goldy) SetComment(comment string) *Goldy {
	gld.comment = comment
	return gld
}

// SetContent sets golden file content. If the golden file was a template, it
// expects a template string. Implements fluent interface.
func (gld *Goldy) SetContent(content string) *Goldy {
	gld.content = []byte(content)
	if gld.data != nil {
		gld.tpl = gld.content
		return gld.renderTemplate()
	}
	return gld
}

// Save saves the golden file to the original path.
func (gld *Goldy) Save() {
	gld.t.Helper()

	buf := &bytes.Buffer{}
	comment := gld.comment
	if !strings.HasSuffix(comment, "\n") {
		comment += "\n"
	}
	buf.WriteString(comment)
	buf.WriteString(Marker)
	if gld.data != nil {
		buf.Write(gld.tpl)
	} else {
		buf.Write(gld.content)
	}
	if err := os.WriteFile(gld.pth, buf.Bytes(), 0600); err != nil {
		gld.t.Fatalf("error writing golden file (%s): %v", gld.pth, err)
	}
}

// renderTemplate renders golden file content as a text template using data
// from [Goldy.data].
func (gld *Goldy) renderTemplate() *Goldy {
	gld.t.Helper()

	var err error
	tpl := template.New("goldy")
	tpl.Option("missingkey=error")
	if tpl, err = tpl.Parse(string(gld.content)); err != nil {
		gld.t.Fatal(err)
		return nil
	}

	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, gld.data); err != nil {
		gld.t.Fatal(err)
		return nil
	}
	gld.content = buf.Bytes()
	return gld
}
