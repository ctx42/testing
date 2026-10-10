// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package must_test

import (
	"fmt"

	"github.com/ctx42/testing/pkg/must"
)

//nolint:unparam
func ExampleFirst() {
	type Row struct{ Name string }

	// Query to a database returning rows.
	query := func() ([]Row, error) {
		return []Row{{"a"}, {"b"}}, nil
	}

	have := must.First(query())

	fmt.Println(have)
	// Output:
	// {a}
}

//nolint:unparam
func ExampleSingle() {
	type Row struct{ Name string }

	// Query to a database returning rows.
	query := func() ([]Row, error) {
		return []Row{{"a"}}, nil
	}

	// Panics if the query returned more than one row.
	have := must.Single(query())

	fmt.Println(have)
	// Output:
	// {a}
}
