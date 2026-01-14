// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package goldy_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/goldy"
)

func ExampleOpen() {
	tspy := core.NewSpy()

	gld := goldy.Open(tspy, "testdata/test_case1.gld")

	fmt.Println(gld.String())
	// Output:
	// Content #1.
	// Content #2.
}

func ExampleCreate() {
	tspy := core.NewSpy()

	dir, err := os.MkdirTemp("", "example-create-*")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()
	pth := filepath.Join(dir, "example.gld")

	gld := goldy.Create(tspy, pth)
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
