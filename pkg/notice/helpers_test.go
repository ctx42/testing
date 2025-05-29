// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
)

func Test_Indent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "")

		// --- Then ---
		affirm.Equal(t, "", have)
	})

	t.Run("one line", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc")

		// --- Then ---
		affirm.Equal(t, "abc", have)
	})

	t.Run("multiple lines", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			" abc\n" +
			" def\n" +
			" ghi"
		affirm.Equal(t, want, have)
	})

	t.Run("empty lines are not indented", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc\ndef\n\nghi")

		// --- Then ---
		want := "" +
			" abc\n" +
			" def\n" +
			"\n" +
			" ghi"
		affirm.Equal(t, want, have)
	})

	t.Run("use tabs", func(t *testing.T) {
		// --- When ---
		have := Indent(1, '\t', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			"\tabc\n" +
			"\tdef\n" +
			"\tghi"
		affirm.Equal(t, want, have)
	})

	t.Run("no ident", func(t *testing.T) {
		// --- When ---
		have := Indent(0, ' ', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			"abc\n" +
			"def\n" +
			"ghi"
		affirm.Equal(t, want, have)
	})
}

func Test_TrialCmp_tabular(t *testing.T) {
	tt := []struct {
		testN string

		a    string
		b    string
		want int
	}{
		{"both empty", "", "", 0},
		{"2", "", "A", -1},
		{"3", "A", "", 1},
		{"equal", "ABC", "ABC", 0},
		{"5", "ABC", "XYZ", -1},
		{"6", "XYZ", "ABC", 1},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			a := New("header").SetTrail(tc.a)
			b := New("header").SetTrail(tc.b)

			// --- When ---
			have := TrialCmp(a, b)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_Pad_tabular(t *testing.T) {
	tt := []struct {
		testN string

		str    string
		length int
		want   string
	}{
		{"empty string zero pad", "", 0, ""},
		{"empty string with pad", "", 3, "   "},
		{"string with pad longer than string length", "abc", 5, "  abc"},
		{"string with pad equal to string length", "abc", 3, "abc"},
		{"string with pad shorter than string length", "abc", 2, "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := Pad(tc.str, tc.length)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_SortNotices(t *testing.T) {
	t.Run("nil chain", func(t *testing.T) {
		// --- Given ---
		var head *Notice

		// --- When ---
		have := SortNotices(head, TrialCmp)

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("single node", func(t *testing.T) {
		// --- Given ---
		msg := New("header").SetTrail("a")

		// --- When ---
		have := SortNotices(msg, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "a", have.Trail)
		affirm.Nil(t, have.prev)
		affirm.Nil(t, have.next)
	})

	t.Run("two sorted nodes", func(t *testing.T) {
		// --- Given ---
		msgA := New("hdrA").SetTrail("A")
		msgB := New("hdrB").SetTrail("B")
		_ = Join(msgA, msgB)

		// --- When ---
		have := SortNotices(msgA, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgB, have))
		affirm.Equal(t, "| hdrA (A) -> hdrB (B)", FWD(have.Head()))
		affirm.Equal(t, "| hdrB (B) -> hdrA (A)", REV(have))
	})

	t.Run("two unsorted nodes", func(t *testing.T) {
		// --- Given ---
		msgA := New("hdrA").SetTrail("A")
		msgB := New("hdrB").SetTrail("B")
		_ = Join(msgB, msgA)

		// --- When ---
		have := SortNotices(msgB, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgB, have))
		affirm.Equal(t, "| hdrA (A) -> hdrB (B)", FWD(have.Head()))
		affirm.Equal(t, "| hdrB (B) -> hdrA (A)", REV(have))
	})

	t.Run("three nodes with equal sort values", func(t *testing.T) {
		// --- Given ---
		msgA0 := New("hdr0").SetTrail("A")
		msgA1 := New("hdr1").SetTrail("A")
		msgA2 := New("hdr2").SetTrail("A")
		_ = Join(msgA0, msgA1, msgA2)

		// --- When ---
		have := SortNotices(msgA0, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgA2, have))
		affirm.Equal(t, "| hdr0 (A) -> hdr1 (A) -> hdr2 (A)", FWD(have.Head()))
		affirm.Equal(t, "| hdr2 (A) -> hdr1 (A) -> hdr0 (A)", REV(have))
	})

	t.Run("chain starts with a couple of equal sort values", func(t *testing.T) {
		// --- Given ---
		msgA0 := New("hdr0").SetTrail("A")
		msgA1 := New("hdr1").SetTrail("A")
		msgB := New("hdrB").SetTrail("B")
		_ = Join(msgA0, msgA1, msgB)

		// --- When ---
		have := SortNotices(msgA0, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgB, have))
		affirm.Equal(t, "| hdr0 (A) -> hdr1 (A) -> hdrB (B)", FWD(have.Head()))
		affirm.Equal(t, "| hdrB (B) -> hdr1 (A) -> hdr0 (A)", REV(have))
	})

	t.Run("three nodes unsorted", func(t *testing.T) {
		// --- Given ---
		msgA := New("hdrA").SetTrail("A")
		msgB := New("hdrB").SetTrail("B")
		msgC := New("hdrC").SetTrail("C")
		_ = Join(msgC, msgA, msgB)

		// --- When ---
		have := SortNotices(msgC, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		affirm.Equal(t, "| hdrA (A) -> hdrB (B) -> hdrC (C)", FWD(have.Head()))
		affirm.Equal(t, "| hdrC (C) -> hdrB (B) -> hdrA (A)", REV(have))
	})

	t.Run("nodes unsorted with repetitions", func(t *testing.T) {
		// --- Given ---
		msgA := New("hdrA").SetTrail("A")
		msgB0 := New("hdrB0").SetTrail("B")
		msgB1 := New("hdrB1").SetTrail("B")
		msgC := New("hdrC").SetTrail("C")
		_ = Join(msgC, msgA, msgB0, msgB1)

		// --- When ---
		have := SortNotices(msgC, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		affirm.Equal(t, "| hdrA (A) -> hdrB0 (B) -> hdrB1 (B) -> hdrC (C)", FWD(have.Head()))
		affirm.Equal(t, "| hdrC (C) -> hdrB1 (B) -> hdrB0 (B) -> hdrA (A)", REV(have))
	})

	t.Run("nodes unsorted with repetitions 2", func(t *testing.T) {
		// --- Given ---
		msgA := New("hdrA").SetTrail("A")
		msgB0 := New("hdrB0").SetTrail("B")
		msgB1 := New("hdrB1").SetTrail("B")
		msgC := New("hdrC").SetTrail("C")
		_ = Join(msgB0, msgC, msgA, msgB1)

		// --- When ---
		have := SortNotices(msgB0, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		affirm.Equal(t, "| hdrA (A) -> hdrB0 (B) -> hdrB1 (B) -> hdrC (C)", FWD(have.Head()))
		affirm.Equal(t, "| hdrC (C) -> hdrB1 (B) -> hdrB0 (B) -> hdrA (A)", REV(have))
	})

	t.Run("nodes with empty sort values", func(t *testing.T) {
		// --- Given ---
		msg0 := New("hdr0").SetTrail("")
		msgA := New("hdrA").SetTrail("A")
		msg1 := New("hdr2").SetTrail("")
		_ = Join(msg0, msgA, msg1)

		// --- When ---
		have := SortNotices(msg0, TrialCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgA, have))
		affirm.Equal(t, "| hdr0 () -> hdr2 () -> hdrA (A)", FWD(have.Head()))
		affirm.Equal(t, "| hdrA (A) -> hdr2 () -> hdr0 ()", REV(have))
	})
}
