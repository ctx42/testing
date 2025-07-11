// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package goldy_test

import (
	"fmt"

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
