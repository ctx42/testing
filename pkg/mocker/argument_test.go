// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_argument_genName_tabular(t *testing.T) {
	tt := []struct {
		testN string

		index int
		name  string
		want  string
	}{
		{
			"argument with name",
			0,
			"name",
			"name",
		},
		{
			"argument without name",
			0,
			"",
			"_a0",
		},
		{
			"argument without name and non-zero index",
			123,
			"",
			"_a123",
		},
		{
			"argument with underscore as a name ",
			3,
			"_",
			"_a3",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			arg := argument{name: tc.name}

			// --- When ---
			have := arg.genName(tc.index)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_argument_genArg_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name  string
		typ   string
		index int
		want  string
	}{
		{
			"argument with a name",
			"name",
			"int",
			0,
			"name int",
		},
		{
			"argument without a name",
			"",
			"int",
			123,
			"_a123 int",
		},
		{
			"argument with underscore as a name",
			"_",
			"int",
			2,
			"_a2 int",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			arg := argument{name: tc.name, typ: tc.typ}

			// --- When ---
			have := arg.genArg(tc.index)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_argument_getType_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name  string
		typ   string
		index int
		want  string
	}{
		{
			"argument with name",
			"name",
			"int",
			0,
			"int",
		},
		{
			"argument without name",
			"",
			"int",
			0,
			"int",
		},
		{
			"argument without name and non-zero index",
			"",
			"int",
			123,
			"int",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			arg := argument{
				name: tc.name,
				typ:  tc.typ,
			}

			// --- When ---
			have := arg.getType()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_argument_isVariadic_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  string
		want bool
	}{
		{"single", "int", false},
		{"variadic", "...int", true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			arg := argument{typ: tc.typ}

			// --- When ---
			have := arg.isVariadic()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
