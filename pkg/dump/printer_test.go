// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_NewPrinter(t *testing.T) {
	// --- Given ---
	dmp := New()

	// --- When ---
	have := NewPrinter(dmp)

	// --- Then ---
	affirm.NotNil(t, have.buf)
}

func Test_Printer_NLI(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false}

		// --- When ---
		have := NewPrinter(dmp).NLI(1)

		// --- Then ---
		affirm.Equal(t, "\n", have.String())
	})

	t.Run("default with zero count", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false}

		// --- When ---
		have := NewPrinter(dmp).NLI(0)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: true}

		// --- When ---
		have := NewPrinter(dmp).NLI(1)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})
}

func Test_Printer_NL(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false}

		// --- When ---
		have := NewPrinter(dmp).NL()

		// --- Then ---
		affirm.Equal(t, "\n", have.String())
	})

	t.Run("flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: true}

		// --- When ---
		have := NewPrinter(dmp).NL()

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})
}

func Test_Printer_Comma(t *testing.T) {
	t.Run("not last and not flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false}

		// --- When ---
		have := NewPrinter(dmp).Comma(false)

		// --- Then ---
		affirm.Equal(t, ",", have.String())
	})

	t.Run("not last and flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: true}

		// --- cfgn ---
		have := NewPrinter(dmp).Comma(false)

		// --- Then ---
		affirm.Equal(t, ",", have.String())
	})

	t.Run("last and not flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false}

		// --- When ---
		have := NewPrinter(dmp).Comma(true)

		// --- Then ---
		affirm.Equal(t, ",", have.String())
	})

	t.Run("last and flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: true}

		// --- When ---
		have := NewPrinter(dmp).Comma(true)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})
}

func Test_Printer_Tab(t *testing.T) {
	t.Run("default and positive n", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false, TabWidth: 2}

		// --- When ---
		have := NewPrinter(dmp).Tab(2)

		// --- Then ---
		affirm.Equal(t, "    ", have.String())
	})

	t.Run("default and negative n", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: false}

		// --- When ---
		have := NewPrinter(dmp).Tab(-2)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("flat", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Flat: true}

		// --- When ---
		have := NewPrinter(dmp).Tab(2)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})
}

func Test_Printer_Space(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: false}

		// --- When ---
		have := NewPrinter(dmp).Space()

		// --- Then ---
		affirm.Equal(t, " ", have.String())
	})

	t.Run("compact", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: true}

		// --- When ---
		have := NewPrinter(dmp).Space()

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})
}

func Test_Printer_Sep(t *testing.T) {
	t.Run("default and not last", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: false}

		// --- When ---
		have := NewPrinter(dmp).Sep(false)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("default and last", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: false}

		// --- When ---
		have := NewPrinter(dmp).Sep(true)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("compact and flat and last", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: true, Flat: true}

		// --- When ---
		have := NewPrinter(dmp).Sep(true)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("compact and not flat and last", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: true, Flat: false}

		// --- When ---
		have := NewPrinter(dmp).Sep(true)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("not compact and not flat and last", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: true, Flat: false}

		// --- When ---
		have := NewPrinter(dmp).Sep(true)

		// --- Then ---
		affirm.Equal(t, "", have.String())
	})

	t.Run("not compact and flat and not last", func(t *testing.T) {
		// --- Given ---
		dmp := Dump{Compact: false, Flat: true}

		// --- When ---
		have := NewPrinter(dmp).Sep(false)

		// --- Then ---
		affirm.Equal(t, " ", have.String())
	})
}

func Test_Printer_Write_String(t *testing.T) {
	// --- Given ---
	dmp := New()

	// --- When ---
	have := NewPrinter(dmp).Write("test")

	// --- Then ---
	affirm.Equal(t, "test", have.String())
}
