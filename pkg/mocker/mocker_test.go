// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/internal/tstmod"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
	"github.com/ctx42/testing/pkg/must"
)

func Test_New(t *testing.T) {
	// --- When ---
	have := New()

	// --- Then ---
	assert.NotNil(t, have.res)
}

func Test_Mocker_Generate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		mod := tstmod.New(t, "v2")

		opts := []Option{
			WithSrc("testdata/cases"),
			WithTgt(mod.Dir),
		}
		mck := New()

		// --- When ---
		err := mck.Generate("Case54", opts...)

		// --- Then ---
		assert.NoError(t, err)

		gld := goldy.Open(t, "testdata/generate_success.gld")
		have := must.Value(os.ReadFile(mod.Path("case54_mock.go")))
		assert.Equal(t, gld.String(), string(have))
	})

	t.Run("error - configuration", func(t *testing.T) {
		// --- Given ---
		mck := New()

		// --- When ---
		err := mck.Generate("")

		// --- Then ---
		assert.ErrorEqual(t, "interface name is required for mocking", err)
	})

	t.Run("error - unknown interface", func(t *testing.T) {
		// --- Given ---
		opts := []Option{
			WithTgtOutput(&bytes.Buffer{}), // Do not create the output file.
		}
		mck := New()

		// --- When ---
		err := mck.Generate("Unknown", opts...)

		// --- Then ---
		assert.ErrorIs(t, ErrUnkType, err)
	})

	t.Run("error - not an interface", func(t *testing.T) {
		// --- Given ---
		opts := []Option{
			WithTgtOutput(&bytes.Buffer{}), // Do not create the output file.
			WithSrc("testdata/cases"),
		}
		mck := New()

		// --- When ---
		err := mck.Generate("Concrete", opts...)

		// --- Then ---
		assert.ErrorIs(t, ErrUnkItf, err)
		assert.ErrorContain(t, "Concrete is not an interface", err)
	})

	t.Run("error - not an interface alias type", func(t *testing.T) {
		// --- Given ---
		opts := []Option{
			WithTgtOutput(&bytes.Buffer{}), // Do not create the output file.
			WithSrc("testdata/cases"),
		}
		mck := New()

		// --- When ---
		err := mck.Generate("Alias", opts...)

		// --- Then ---
		assert.ErrorIs(t, ErrUnkItf, err)
		assert.ErrorContain(t, "Alias is not an interface", err)
	})

	t.Run("error - empty interface", func(t *testing.T) {
		// --- Given ---
		opts := []Option{
			WithTgtOutput(&bytes.Buffer{}), // Do not create the output file.
			WithSrc("testdata/cases"),
		}
		mck := New()

		// --- When ---
		err := mck.Generate("Empty", opts...)

		// --- Then ---
		assert.ErrorIs(t, ErrNoMethods, err)
	})

	t.Run("error - an any interface", func(t *testing.T) {
		// --- Given ---
		opts := []Option{
			WithTgtOutput(&bytes.Buffer{}), // Do not create the output file.
			WithSrc("testdata/cases"),
		}
		mck := New()

		// --- When ---
		err := mck.Generate("EmptyAny", opts...)

		// --- Then ---
		assert.ErrorIs(t, ErrUnkItf, err)
		assert.ErrorContain(t, "EmptyAny is not an interface", err)
	})
}

func Test_Mocker_Generate_tabular(t *testing.T) {
	tt := []struct {
		testN string

		itfName string
		srcPath string
		tgtPath string
	}{
		{"Case00", "Case00", "cases", "golden"},
		{"Case01", "Case01", "cases", "golden"},
		{"Case02", "Case02", "cases", "golden"},
		{"Case03", "Case03", "cases", "golden"},
		{"Case04", "Case04", "cases", "golden"},
		{"Case05", "Case05", "cases", "golden"},
		{"Case06", "Case06", "cases", "golden"},
		{"Case07", "Case07", "cases", "golden"},
		{"Case08", "Case08", "cases", "golden"},
		{"Case09", "Case09", "cases", "golden"},
		{"Case10", "Case10", "cases", "golden"},
		{"Case11", "Case11", "cases", "golden"},
		{"Case12", "Case12", "cases", "golden"},
		{"Case13", "Case13", "cases", "golden"},
		{"Case14", "Case14", "cases", "golden"},
		{"Case15", "Case15", "cases", "golden"},
		{"Case16", "Case16", "cases", "golden"},
		{"Case17", "Case17", "cases", "golden"},
		{"Case17_dst_cases", "Case17", "cases", "cases"},
		{"Case18", "Case18", "cases", "golden"},
		{"Case19", "Case19", "cases", "golden"},
		{"Case20", "Case20", "cases", "golden"},
		{"Case21", "Case21", "cases", "golden"},
		{"Case22", "Case22", "cases", "golden"},
		{"Case23", "Case23", "cases", "golden"},
		{"Case24", "Case24", "cases", "golden"},
		{"Case25", "Case25", "cases", "golden"},
		{"Case26", "Case26", "cases", "golden"},
		{"Case27", "Case27", "cases", "golden"},
		{"Case28", "Case28", "cases", "golden"},
		{"Case29", "Case29", "cases", "golden"},
		{"Case30", "Case30", "cases", "golden"},
		{"Case31", "Case31", "cases", "golden"},
		{"Case32", "Case32", "cases", "golden"},
		{"Case33", "Case33", "cases", "golden"},
		{"Case34", "Case34", "cases", "golden"},
		{"Case35", "Case35", "cases", "golden"},
		{"Case36", "Case36", "cases", "golden"},
		{"Case37", "Case37", "cases", "golden"},
		{"Case38", "Case38", "cases", "golden"},
		{"Case39", "Case39", "cases", "golden"},
		{"Case40", "Case40", "cases", "golden"},
		{"Case41", "Case41", "cases", "golden"},
		{"Case42", "Case42", "cases", "golden"},
		{"Case43", "Case43", "cases", "golden"},
		{"Case44", "Case44", "cases", "golden"},
		{"Case45", "Case45", "cases", "golden"},
		{"Case46", "Case46", "cases", "golden"},
		{"Case47", "Case47", "cases", "golden"},
		{"Case48", "Case48", "cases", "golden"},
		{"Case48_dst_cases", "Case48", "cases", "cases"},
		{"Case48_dst_pkgb", "Case48", "cases", "pkgb"},
		{"Case49", "Case49", "cases", "golden"},
		{"Case50", "Case50", "cases", "golden"},
		{"Case51", "Case51", "cases", "golden"},
		{"Case52", "Case52", "cases", "golden"},
		{"Case53", "Case53", "cases", "golden"},
		{"Case54", "Case54", "cases", "golden"},
		{"Case54_dst_cases", "Case54", "cases", "cases"},
		{"Case54_dst_pkge", "Case54", "cases", "pkge"},
		{"Case55", "Case55", "cases", "golden"},
		{"Case55_dst_cases", "Case55", "cases", "cases"},
		{"Case56", "Case56", "cases", "golden"},

		{"ItfA", "ItfA", "cases", "golden"},
		{"ItfB", "ItfB", "cases", "golden"},
		{"EmbedLocal", "EmbedLocal", "cases", "golden"},
		{"Embedder", "Embedder", "cases", "golden"},
		{"EmptyEmbed", "EmptyEmbed", "cases", "golden"},
		{"Massive", "Massive", "cases", "golden"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			buf := &bytes.Buffer{}
			impPath := "github.com/ctx42/testing/pkg/mocker/testdata/"
			opts := []Option{
				WithSrc(impPath + tc.srcPath),
				WithTgt(impPath + tc.tgtPath),
				WithTgtName(tc.itfName),
				WithTgtOutput(buf),
			}

			// --- When ---
			err := New().Generate(tc.itfName, opts...)

			// --- Then ---
			assert.NoError(t, err)

			gfp := filepath.Join("testdata/golden", tc.testN+".gld")
			// nolint: gocritic
			// goldy.New(t, gfp, "", buf.String()).Save()

			want := goldy.Open(t, gfp)
			assert.Equal(t, want.String(), buf.String())
		})
	}
}
