// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package must

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_Value(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// --- When ---
		have := Value(types.NewTInt(42))

		// --- Then ---
		affirm.Equal(t, 42, have.V)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		msg := affirm.Panic(t, func() { Value(types.NewTInt(40)) })

		// --- Then ---
		affirm.Equal(t, "not cool", *msg)
	})
}

func Test_Values(t *testing.T) {
	fnGood := func() (int, float64, error) { return 1, 2, nil }
	fnBad := func() (int, float64, error) { return 0, 0, errors.New("test") }

	t.Run("no error", func(t *testing.T) {
		// --- When ---
		have1, have2 := Values(fnGood())

		// --- Then ---
		affirm.Equal(t, 1, have1)
		affirm.Equal(t, 2.0, have2)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		msg := affirm.Panic(t, func() { Values(fnBad()) })

		// --- Then ---
		affirm.Equal(t, "test", *msg)
	})
}

func Test_Nil(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		Nil(nil)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		msg := affirm.Panic(t, func() { Nil(errors.New("test err")) })

		// --- Then ---
		affirm.Equal(t, "test err", *msg)
	})
}

func Test_First(t *testing.T) {
	type T struct{ V int }

	t.Run("one element no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return []T{{V: 1}}, nil } // nolint:unparam

		// --- When ---
		have := First(fn())

		// --- Then ---
		affirm.Equal(t, T{V: 1}, have)
	})

	t.Run("zero elements no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return nil, nil } // nolint:unparam

		// --- When ---
		have := First(fn())

		// --- Then ---
		affirm.Equal(t, T{}, have)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return nil, errors.New("test msg") }

		// --- When ---
		msg := affirm.Panic(t, func() { First(fn()) })

		// --- Then ---
		affirm.Equal(t, "test msg", *msg)
	})

	t.Run("more than one element no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { // nolint:unparam
			return []T{{V: 1}, {V: 2}}, nil
		}

		// --- When ---
		have := First(fn())

		// --- Then ---
		affirm.Equal(t, T{V: 1}, have)
	})
}

func Test_Single(t *testing.T) {
	type T struct{ V int }

	t.Run("one element no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return []T{{V: 1}}, nil } // nolint:unparam

		// --- When ---
		have := Single(fn())

		// --- Then ---
		affirm.Equal(t, T{V: 1}, have)
	})

	t.Run("zero elements no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return nil, nil } // nolint:unparam

		// --- When ---
		have := Single(fn())

		// --- Then ---
		affirm.Equal(t, T{}, have)
	})

	t.Run("error - not nil", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return nil, errors.New("test msg") }

		// --- When ---
		msg := affirm.Panic(t, func() { Single(fn()) })

		// --- Then ---
		affirm.Equal(t, "test msg", *msg)
	})

	t.Run("more than one element no error", func(t *testing.T) {
		// --- Given ---
		s := []T{{V: 1}, {V: 2}}
		fn := func() ([]T, error) { return s, nil } // nolint: unparam

		// --- When ---
		msg := affirm.Panic(t, func() { Single(fn()) })

		// --- Then ---
		affirm.Equal(t, "expected a single result", *msg)
	})
}
