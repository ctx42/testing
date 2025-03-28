// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
)

func Test_chanDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val  any
		want string
	}{
		{"chan0", make(chan int), "(chan int)(<addr>)"},
		{"chan1", make(chan func()), "(chan func())(<addr>)"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			val := reflect.ValueOf(tc.val)

			// --- When ---
			have := chanDumper(Dump{}, 0, val)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_chanDumper(t *testing.T) {
	t.Run("nil channel", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig(WithPtrAddr))
		var ch chan int
		val := reflect.ValueOf(ch)

		// --- When ---
		have := chanDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "(chan int)(<0x0>)", have)
	})

	t.Run("usage error", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig())
		val := reflect.ValueOf(1234)

		// --- When ---
		have := chanDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, valErrUsage, have)
	})

	t.Run("usage error uses level", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig())
		val := reflect.ValueOf(1234)

		// --- When ---
		have := chanDumper(dmp, 2, val)

		// --- Then ---
		affirm.Equal(t, "\t\t"+valErrUsage, have)
	})

	t.Run("print pointer address", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig(WithPtrAddr))
		val := reflect.ValueOf(make(chan int))

		// --- When ---
		have := chanDumper(dmp, 0, val)

		// --- Then ---
		want := fmt.Sprintf("(chan int)(<0x%x>)", val.Pointer())
		affirm.Equal(t, want, have)
	})

	t.Run("uses level", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig())
		val := reflect.ValueOf(make(chan int))

		// --- When ---
		have := chanDumper(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "\t(chan int)(<addr>)", have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig(WithIndent(2)))
		val := reflect.ValueOf(make(chan int))

		// --- When ---
		have := chanDumper(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "\t\t\t(chan int)(<addr>)", have)
	})
}
