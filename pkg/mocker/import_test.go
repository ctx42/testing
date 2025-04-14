// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewImport(t *testing.T) {
	// --- When ---
	have := NewImport("github.com/ctx42/testing")

	// --- Then ---
	want := Import{
		Alias: "",
		Name:  "testing",
		Spec:  "github.com/ctx42/testing",
		Dir:   "",
	}
	assert.Equal(t, want, have)
}

func Test_Import_SetSpec(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		imp := Import{}

		// --- When ---
		have := imp.SetSpec("github.com/ctx42/testing/pkg/mocker")

		// --- Then ---
		want := Import{
			Alias: "",
			Name:  "mocker",
			Spec:  "github.com/ctx42/testing/pkg/mocker",
			Dir:   "",
		}
		assert.Equal(t, want, have)
	})

	t.Run("does not change the directory", func(t *testing.T) {
		// --- Given ---
		imp := Import{Dir: "dir"}

		// --- When ---
		have := imp.SetSpec("github.com/ctx42/testing/pkg/mocker")

		// --- Then ---
		assert.Equal(t, "dir", have.Dir)
	})
}

func Test_Import_SetAlias(t *testing.T) {
	// --- Given ---
	imp := Import{}

	// --- When ---
	have := imp.SetAlias("alias")

	// --- Then ---
	assert.Equal(t, "alias", have.Alias)
}

func Test_Import_SetDir(t *testing.T) {
	// --- Given ---
	imp := Import{}

	// --- When ---
	have := imp.SetDir("dir")

	// --- Then ---
	assert.Equal(t, "dir", have.Dir)
}

func Test_Import_IsDot_tabular(t *testing.T) {
	tt := []struct {
		testN string

		alias string
		want  bool
	}{
		{"empty", "", false},
		{"not dot", "alias", false},
		{"dot", ".", true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			imp := Import{Alias: tc.alias}

			// --- When ---
			have := imp.IsDot()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_Import_IsZero_tabular(t *testing.T) {
	tt := []struct {
		testN string

		imp  Import
		want bool
	}{
		{"only spec field set", Import{Spec: "abc"}, false},
		{"only dir field set", Import{Dir: "abc"}, false},
		{"alias and name set", Import{Alias: "a", Name: "n"}, true},
		{"no fields set", Import{}, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := tc.imp.IsZero()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_Import_GoString_tabular(t *testing.T) {
	tt := []struct {
		testN string

		alias string
		path  string
		want  string
	}{
		{"without alias", "", "path", `"path"`},
		{"with alias", "alias", "path", `alias "path"`},
		{"empty", "", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			imp := Import{Alias: tc.alias, Spec: tc.path}

			// --- When ---
			have := imp.GoString()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
