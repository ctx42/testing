// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
)

func Test_goitf_find(t *testing.T) {
	t.Run("method found", func(t *testing.T) {
		// --- Given ---
		itf := &goitf{
			methods: []*method{
				{name: "Method0"},
				{name: "Method1"},
				{name: "Method2"},
			},
		}

		// --- When ---
		have, err := itf.find("Method1")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, &method{name: "Method1"}, have)
	})

	t.Run("method not found", func(t *testing.T) {
		// --- Given ---
		itf := &goitf{
			methods: []*method{
				{name: "Method0"},
				{name: "Method1"},
				{name: "Method2"},
			},
		}

		// --- When ---
		have, err := itf.find("Method3")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkMet, err)
		assert.Nil(t, have)
	})

	t.Run("error - when interface without methods", func(t *testing.T) {
		// --- Given ---
		itf := &goitf{}

		// --- When ---
		have, err := itf.find("Method0")

		// --- Then ---
		assert.ErrorIs(t, ErrUnkMet, err)
		assert.Nil(t, have)
	})
}

func Test_goitf_generate(t *testing.T) {
	t.Run("single simple method", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_goitf/simple.gld"
		itf := goitf{
			name: "MyItf",
			methods: []*method{
				{
					name: "Method0",
				},
			},
		}

		// --- When ---
		have := itf.generate("MyMock", false)

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})

	t.Run("methods with args", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_goitf/methods_with_args.gld"
		itf := goitf{
			name: "MyItf",
			methods: []*method{
				{
					name: "Method0",
					args: []argument{
						{name: "a", typ: "int"},
					},
				},
				{
					name: "Method1",
					args: []argument{
						{name: "a", typ: "string"},
						{name: "b", typ: "...int"},
					},
				},
			},
		}

		// --- When ---
		have := itf.generate("MyMock", false)

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})

	t.Run("single simple with OnXXX helpers", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_goitf/simple_with_onh.gld"
		itf := goitf{
			name: "MyItf",
			methods: []*method{
				{
					name: "Method0",
				},
			},
		}

		// --- When ---
		have := itf.generate("MyMock", true)

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})
}

func Test_goitf_imports(t *testing.T) {
	t.Run("no methods", func(t *testing.T) {
		// --- Given ---
		itf := goitf{}

		// --- When ---
		have := itf.imports()

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("success", func(t *testing.T) {
		// --- Given ---
		itf := goitf{
			methods: []*method{
				{
					name: "Method0",
					args: []argument{
						{
							name: "a",
							pks: []*gopkg{
								{pkgName: "a0", pkgPath: "a0_path"},
							},
						},
						{
							name: "b",
							pks: []*gopkg{
								{pkgName: "a0", pkgPath: "a0_path"},
							},
						},
					},
					rets: []argument{
						{
							name: "a",
							pks: []*gopkg{
								{pkgName: "r0", pkgPath: "r0_path"},
							},
						},
						{
							name: "b",
							pks: []*gopkg{
								{pkgName: "r0", pkgPath: "r0_path"},
							},
						},
					},
				},
				{
					name: "Method1",
					args: []argument{
						{
							name: "a",
							pks: []*gopkg{
								{pkgName: "a0", pkgPath: "a0_path"},
							},
						},
						{
							name: "b",
							pks: []*gopkg{
								{pkgName: "a0", pkgPath: "a0_path"},
							},
						},
					},
					rets: []argument{
						{
							name: "a",
							pks: []*gopkg{
								{pkgName: "r0", pkgPath: "r0_path"},
							},
						},
						{
							name: "b",
							pks: []*gopkg{
								{pkgName: "r0", pkgPath: "r0_path"},
							},
						},
						{
							name: "c",
							pks: []*gopkg{
								{pkgName: "r1", pkgPath: "r1_path"},
							},
						},
					},
				},
			},
		}

		// --- When ---
		have := itf.imports()

		// --- Then ---
		want := []*gopkg{
			{pkgName: "a0", pkgPath: "a0_path"},
			{pkgName: "r0", pkgPath: "r0_path"},
			{pkgName: "r1", pkgPath: "r1_path"},
		}
		assert.Equal(t, want, have)
	})
}

func Test_goitf_genImports(t *testing.T) {
	t.Run("no imports", func(t *testing.T) {
		// --- Given ---
		itf := &goitf{}

		// --- When ---
		have := itf.genImports()

		// --- Then ---
		want := fmt.Sprintf("import (\n\t%q\n\t%q\n)", selfImp, testerImp)
		assert.Equal(t, want, have)
	})
}
