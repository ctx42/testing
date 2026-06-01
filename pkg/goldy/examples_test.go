// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package goldy_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/goldy"
)

func ExampleOpen() {
	gld := goldy.Open(&testing.T{}, "testdata/test_case1.gld")

	fmt.Println(gld.String())
	// Output:
	// Content #1.
	// Content #2.
}

func ExampleCreate() {
	dir, err := os.MkdirTemp("", "example-create-*")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()
	pth := filepath.Join(dir, "example.gld")

	gld := goldy.Create(&testing.T{}, pth)
	gld.SetComment("Multi\nline\ncontent")
	gld.SetContent("Content #1.\nContent #2.")
	gld.Save()

	have, err := os.ReadFile(pth)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(have))
	// Output:
	// Multi
	// line
	// content
	// ---
	// Content #1.
	// Content #2.
}

func ExampleOpen_withData() {
	tspy := &testing.T{}

	// Open supports Go text/template expansion via WithData.
	data := map[string]any{"first": 1}

	gld := goldy.Open(tspy, "testdata/test_tpl.gld", goldy.WithData(data))

	fmt.Println(gld.String())
	// Output:
	// Content #1.
}
