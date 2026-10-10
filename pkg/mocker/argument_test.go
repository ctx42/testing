// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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
		{"argument with name", 0, "name", "name"},
		{"argument without name", 0, "", "_a0"},
		{"argument without name and non-zero index", 123, "", "_a123"},
		{"argument with underscore as a name", 3, "_", "_a3"},
		{"reserved local _args", 0, "_args", "_a0"},
		{"reserved local _mck", 1, "_mck", "_a1"},
		{"reserved return temp _r0", 2, "_r0", "_a2"},
		{"reserved synthetic arg name _a0", 3, "_a0", "_a3"},
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

func Test_isGenLocal_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name string
		want bool
	}{
		{"mock receiver", "_mck", true},
		{"ok local", "_ok", true},
		{"synthetic argument", "_a12", true},
		{"return temp", "_r0", true},
		{"short", "_a", false},
		{"other letter", "_x1", false},
		{"letters after prefix", "_ab", false},
		{"user name", "name", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := isGenLocal(tc.name)

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
