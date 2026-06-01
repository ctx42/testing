// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package goldy simplifies golden file testing.
//
// It loads files that contain an optional comment section followed by
// the [Marker] separator and the expected content. The loaded content
// can be used directly with [assert.Equal] or [check.Equal] (often via
// [Goldy.Content] or after template expansion with [WithData]).
//
// Golden files are typically stored under testdata/ and committed.
// The package integrates with [tester.T] for failure reporting.
//
// See the package [README] and examples for typical usage with
// [github.com/ctx42/testing/pkg/assert] and [github.com/ctx42/testing/pkg/check].
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

	"github.com/ctx42/testing/pkg/tester"
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
	t       tester.T       // Test manager.
}

// Open creates a [Goldy] from the file at pth.
//
// Content begins after the mandatory [Marker] line; anything before it is
// treated as a comment. Use options such as [WithData] for templated golden
// files.
func Open(t tester.T, pth string, opts ...func(*Goldy)) *Goldy {
	t.Helper()

	// G304: path comes from test code calling golden file helpers.
	fil, err := os.Open(pth) // nolint:gosec
	if err != nil {
		t.Errorf("error opening file: %v", err)
		return nil
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
			t.Errorf("error reading file: %v", err)
			return nil
		}
		if !started {
			started = bytes.Equal(line, mark)
			if !started {
				if eof {
					m := strings.TrimSpace(Marker)
					t.Errorf("the golden file is missing the %q marker", m)
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

// Create creates a new [Goldy] instance representing an empty golden file.
//
// The file is created (or truncated) at the given path.
func Create(t tester.T, pth string) *Goldy {
	t.Helper()

	// G304: path comes from test code calling golden file helpers.
	fil, err := os.Create(pth) // nolint:gosec
	if err != nil {
		t.Errorf("error creating file: %v", err)
		return nil
	}
	defer func() { _ = fil.Close() }()

	return &Goldy{
		pth:     pth,
		content: make([]byte, 0, 4*1024),
		t:       t,
	}
}

// String implements [fmt.Stringer] and returns the golden file content.
func (gld *Goldy) String() string { return string(gld.content) }

// Bytes returns a clone of the golden file content.
func (gld *Goldy) Bytes() []byte { return slices.Clone(gld.content) }

// SetComment sets the comment section for the golden file.
// Implements fluent interface.
func (gld *Goldy) SetComment(comment string) *Goldy {
	gld.comment = comment
	return gld
}

// SetContent sets the golden file content. When data was provided via
// [WithData], the content is treated as a text template.
//
// Implements fluent interface.
func (gld *Goldy) SetContent(content string) *Goldy {
	gld.content = []byte(content)
	if gld.data != nil {
		gld.tpl = gld.content
		return gld.renderTemplate()
	}
	return gld
}

// Save writes the golden file (comment + [Marker] + content) back to its
// original path. It reports errors via the test's t.Error.
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
		gld.t.Errorf("error writing golden file (%s): %v", gld.pth, err)
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
		gld.t.Error(err)
		return nil
	}

	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, gld.data); err != nil {
		gld.t.Error(err)
		return nil
	}
	gld.content = buf.Bytes()
	return gld
}
