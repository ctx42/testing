// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package goldy simplifies golden file testing.
//
// It loads files that contain an optional comment section followed by
// the [Marker] separator and the expected content. The loaded content
// can be compared with the assert and check packages (often via
// [Goldy.String] / [Goldy.Bytes] or after template expansion with [WithData]).
//
// Golden files are typically stored under testdata/ and committed.
// Failures are reported through [tester.T].
//
// See the package README and examples for typical usage.
package goldy

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/ctx42/testing/pkg/tester"
)

// Marker is a separator between a golden file comment and the content.
const Marker = "---\n"

// WithData is the [Open] option setting [Goldy] data for golden files that
// are text templates. The file is treated as a template even when data is
// nil or empty.
func WithData(data map[string]any) func(*Goldy) {
	return func(gld *Goldy) { gld.data = data; gld.isTpl = true }
}

// bom is the UTF-8 byte order mark.
const bom = "\xef\xbb\xbf"

// Goldy represents a golden file. Build it with [Open] or [Create]; the zero
// value is not usable.
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
	isTpl   bool           // True when the content is a text template.
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
	fil, err := os.Open(pth) //nolint:gosec
	if err != nil {
		t.Fatalf("error opening file: %v", err)
		return nil
	}
	defer func() { _ = fil.Close() }()

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
	for first := true; ; first = false {
		line, err := rdr.ReadBytes('\n')
		eof := errors.Is(err, io.EOF)
		if err != nil && !eof {
			t.Fatalf("error reading file %s: %v", pth, err)
			return nil
		}
		if first {
			line = bytes.TrimPrefix(line, []byte(bom))
		}
		line = crlfToLF(line)
		if !started {
			started = isMarker(line)
			if !started {
				if eof {
					format := "the golden file %s is missing the %q marker"
					t.Fatalf(format, pth, strings.TrimSpace(Marker))
					return nil
				}
				gld.comment += string(line)
			}
			if eof {
				break
			}
			continue
		}
		gld.content = append(gld.content, line...)
		if eof {
			break
		}
	}

	if gld.isTpl {
		gld.tpl = gld.content
		return gld.renderTemplate()
	}
	return gld
}

// Create creates a new [Goldy] instance representing an empty golden file
// at the given path. The file is written only by [Goldy.Save]; Create checks
// that the parent directory exists and fails the test when it does not.
func Create(t tester.T, pth string) *Goldy {
	t.Helper()

	dir := filepath.Dir(pth)
	inf, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("error checking golden file directory: %v", err)
		return nil
	}
	if !inf.IsDir() {
		format := "error checking golden file directory: %s is not a directory"
		t.Fatalf(format, dir)
		return nil
	}

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
// It implements a fluent interface.
func (gld *Goldy) SetComment(comment string) *Goldy {
	gld.comment = comment
	return gld
}

// SetContent sets the golden file content. When data was provided via
// [WithData], the content is treated as a text template.
//
// It implements a fluent interface.
func (gld *Goldy) SetContent(content string) *Goldy {
	gld.content = []byte(content)
	if gld.isTpl {
		gld.tpl = gld.content
		return gld.renderTemplate()
	}
	return gld
}

// Save writes the golden file (comment + [Marker] + content, the raw template
// for templated files) to its path. It reports errors via the test's t.Error,
// including a comment containing a [Marker] line, which would move the
// content when the file is opened.
func (gld *Goldy) Save() {
	gld.t.Helper()

	comment := gld.comment
	if comment != "" && !strings.HasSuffix(comment, "\n") {
		comment += "\n"
	}
	for _, line := range strings.SplitAfter(comment, "\n") {
		if isMarker([]byte(line)) {
			format := "golden file (%s) comment must not contain the %q marker"
			gld.t.Errorf(format, gld.pth, strings.TrimSpace(Marker))
			return
		}
	}

	buf := &bytes.Buffer{}
	_, _ = buf.WriteString(comment)
	_, _ = buf.WriteString(Marker)
	if gld.isTpl {
		_, _ = buf.Write(gld.tpl)
	} else {
		_, _ = buf.Write(gld.content)
	}
	if err := os.WriteFile(gld.pth, buf.Bytes(), 0o600); err != nil {
		gld.t.Errorf("error writing golden file (%s): %v", gld.pth, err)
	}
}

// renderTemplate renders golden file content as a text template using data
// supplied via [WithData].
func (gld *Goldy) renderTemplate() *Goldy {
	gld.t.Helper()

	var err error
	tpl := template.New(gld.pth)
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

// isMarker reports whether line is the [Marker] line, ignoring trailing
// spaces, tabs and the line ending.
func isMarker(line []byte) bool {
	return string(bytes.TrimRight(line, " \t\r\n")) == strings.TrimSpace(Marker)
}
