// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/goldy"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_StructDumper(t *testing.T) {
	t.Run("simple struct with private fields", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		s := testcases.TA{
			Int: 1,
			Str: "2",
			Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			Dur: 3,
			Loc: testcases.WAW,
			TAp: nil,
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_simple.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("simple struct without private fields", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithNoPrivate)
		s := testcases.TA{
			Int: 1,
			Str: "2",
			Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			Dur: 3,
			Loc: testcases.WAW,
			TAp: nil,
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_simple_no_private.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("simple flat and compact struct", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithFlat, WithCompact)
		s := testcases.TA{
			Int: 1,
			Str: "2",
			Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			Dur: 3,
			Loc: testcases.WAW,
			TAp: nil,
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_simple_flat_compact.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("multi level struct", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		s := testcases.TRec{
			Int: 1,
			Rec: &testcases.TRec{
				Int: 2,
			},
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_level.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("multi-level struct with indent", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		s := testcases.TRec{
			Int: 1,
			Rec: &testcases.TRec{
				Int: 2,
			},
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_level_indent.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("multi-level flat and compact struct", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithFlat, WithCompact)
		s := testcases.TRec{
			Int: 1,
			Rec: &testcases.TRec{
				Int: 2,
			},
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_level_flat_compact.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("struct with a multiline string field value", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithFlatStrings(0))
		s := struct {
			F0 int
			F1 string
		}{
			F0: 1,
			F1: "a\nb\nc\n",
		}

		// --- When ---
		have := StructDumper(dmp, 0, reflect.ValueOf(s))

		// --- Then ---
		want := goldy.Open(t, "testdata/struct_multi_line_string_field.gld")
		affirm.Equal(t, want.String(), have)
	})

	t.Run("error - invalid type", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(1))

		// --- When ---
		have := StructDumper(dmp, 2, reflect.ValueOf(123))

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}
