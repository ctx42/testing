// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_FuncDumper(t *testing.T) {
	t.Run("nil function", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithPtrAddr)

		var fn func()
		val := reflect.ValueOf(fn)

		// --- When ---
		have := FuncDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "<func>(<0x0>)", have)
	})

	t.Run("usage error", func(t *testing.T) {
		// --- Given ---
		dmp := New()
		val := reflect.ValueOf(1234)

		// --- When ---
		have := FuncDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, ValErrUsage, have)
	})

	t.Run("print pointer address", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithPtrAddr)
		fn := func() {}
		val := reflect.ValueOf(fn)
		want := fmt.Sprintf("<func>(<0x%x>)", val.Pointer())

		// --- When ---
		have := FuncDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, want, have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(WithIndent(2))
		val := reflect.ValueOf(1234)

		// --- When ---
		have := FuncDumper(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "      "+ValErrUsage, have)
	})
}

func Test_FuncDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val  any
		want string
	}{
		{"func0", func() {}, "<func>(<addr>)"},
		{"func1", func(int) error { return nil }, "<func>(<addr>)"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			val := reflect.ValueOf(tc.val)

			// --- When ---
			have := FuncDumper(Dump{}, 0, val)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}
