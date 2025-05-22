// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/goldy"
)

func Test_dumpStruct(t *testing.T) {
	t.Run("simple struct", func(t *testing.T) {
		// --- Given ---
		s := types.TA{
			Int: 1,
			Str: "2",
			Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			Dur: 3,
			Loc: types.WAW,
			TAp: nil,
		}
		dmp := New()

		// --- When ---
		have := structDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_simple.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("simple flat & compact struct", func(t *testing.T) {
		// --- Given ---
		s := types.TA{
			Int: 1,
			Str: "2",
			Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			Dur: 3,
			Loc: types.WAW,
			TAp: nil,
		}
		dmp := New(WithFlat, WithCompact)

		// --- When ---
		have := structDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_simple_flat_compact.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("multi level struct", func(t *testing.T) {
		// --- Given ---
		s := types.T1{
			Int: 1,
			T1: &types.T1{
				Int: 2,
			},
		}
		dmp := New()

		// --- When ---
		have := structDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_level.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("multi-level struct with indent", func(t *testing.T) {
		// --- Given ---
		s := types.T1{
			Int: 1,
			T1: &types.T1{
				Int: 2,
			},
		}
		dmp := New(WithIndent(2))

		// --- When ---
		have := structDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_level_indent.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("multi-level flat & compact struct", func(t *testing.T) {
		// --- Given ---
		s := types.T1{
			Int: 1,
			T1: &types.T1{
				Int: 2,
			},
		}
		dmp := New(WithFlat, WithCompact)

		// --- When ---
		have := structDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_level_flat_compact.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("struct with a multi-line string field value", func(t *testing.T) {
		// --- Given ---
		s := struct {
			F0 int
			F1 string
		}{
			F0: 1,
			F1: "a\nb\nc\n",
		}
		dmp := New(WithFlatStrings(0))

		// --- When ---
		have := structDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_line_string_field.gld")
		affirm.Equal(t, want.String(), have)
	})
}
