package must_test

import (
	"fmt"

	"github.com/ctx42/testing/pkg/must"
)

// nolint:unparam
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

// nolint:unparam
func ExampleSingle() {
	type Row struct{ Name string }

	// Query to a database returning rows.
	query := func() ([]Row, error) {
		return []Row{{"a"}}, nil
	}

	// Will panic if a database returned more than one error.
	have := must.Single(query())

	fmt.Println(have)
	// Output:
	// {a}
}
