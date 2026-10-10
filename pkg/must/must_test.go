// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package must

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Value(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// --- When ---
		have := Value(testcases.NewTInt(42))

		// --- Then ---
		affirm.Equal(t, 42, have.V)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- When ---
		have := affirm.Panic(t, func() { Value(testcases.NewTInt(40)) })

		// --- Then ---
		affirm.Equal(t, "not cool", *have)
	})
}

func Test_Values(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// --- When ---
		hInt, hFloat := Values(1, 2.0, nil)

		// --- Then ---
		affirm.Equal(t, 1, hInt)
		affirm.Equal(t, 2.0, hFloat)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- Given ---
		err := errors.New("test")

		// --- When ---
		have := affirm.Panic(t, func() { Values(0, 0.0, err) })

		// --- Then ---
		affirm.Equal(t, "test", *have)
	})
}

func Test_Nil(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// --- When ---
		have, _ := core.WillPanic(func() { Nil(nil) })

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- When ---
		have := affirm.Panic(t, func() { Nil(errors.New("test err")) })

		// --- Then ---
		affirm.Equal(t, "test err", *have)
	})
}

func Test_First(t *testing.T) {
	type T struct{ V int }

	t.Run("one element", func(t *testing.T) {
		// --- When ---
		have := First([]T{{V: 1}}, nil)

		// --- Then ---
		affirm.Equal(t, T{V: 1}, have)
	})

	t.Run("zero elements", func(t *testing.T) {
		// --- When ---
		have := First([]T(nil), nil)

		// --- Then ---
		affirm.Equal(t, T{}, have)
	})

	t.Run("more than one element", func(t *testing.T) {
		// --- When ---
		have := First([]T{{V: 1}, {V: 2}}, nil)

		// --- Then ---
		affirm.Equal(t, T{V: 1}, have)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- Given ---
		err := errors.New("test msg")

		// --- When ---
		have, _ := core.WillPanic(func() { First([]T{{V: 1}}, err) })

		// --- Then ---
		affirm.Equal(t, true, have == any(err))
	})

	t.Run("error - wraps the single result error", func(t *testing.T) {
		// --- Given ---
		_, single := single([]T{{V: 1}, {V: 2}}, nil)
		err := fmt.Errorf("query: %w", single)

		// --- When ---
		have, _ := core.WillPanic(func() { First([]T{{V: 1}}, err) })

		// --- Then ---
		affirm.Equal(t, true, have == any(err))
	})
}

func Test_Single(t *testing.T) {
	type T struct{ V int }

	t.Run("one element", func(t *testing.T) {
		// --- When ---
		have := Single([]T{{V: 1}}, nil)

		// --- Then ---
		affirm.Equal(t, T{V: 1}, have)
	})

	t.Run("zero elements", func(t *testing.T) {
		// --- When ---
		have := Single([]T(nil), nil)

		// --- Then ---
		affirm.Equal(t, T{}, have)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- Given ---
		err := errors.New("test msg")

		// --- When ---
		have, _ := core.WillPanic(func() { Single([]T{{V: 1}}, err) })

		// --- Then ---
		affirm.Equal(t, true, have == any(err))
	})

	t.Run("error - more than one element", func(t *testing.T) {
		// --- When ---
		have := affirm.Panic(t, func() { Single([]T{{V: 1}, {V: 2}}, nil) })

		// --- Then ---
		affirm.Equal(t, "expected a single result: got 2", *have)
	})
}
