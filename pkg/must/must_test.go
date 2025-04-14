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
		// --- Given ---
		var have *types.TInt
		var panicked bool

		fn := func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			have = Value(types.NewTInt(40))
		}

		// --- When ---
		fn()

		// --- Then ---
		affirm.Equal(t, true, panicked)
		affirm.Equal(t, true, have == nil)
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
		// --- Given ---
		var t1 int
		var t2 float64
		var panicked bool

		fn := func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			t1, t2 = Values(fnBad())
		}

		// --- When ---
		fn()

		// --- Then ---
		affirm.Equal(t, true, panicked)
		affirm.Equal(t, 0, t1)
		affirm.Equal(t, 0.0, t2)
	})
}

func Test_Nil(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		Nil(nil)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		var panicked bool

		fn := func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			Nil(errors.New("test err"))
		}

		// --- When ---
		fn()

		// --- Then ---
		affirm.Equal(t, true, panicked)
	})
}

func Test_First(t *testing.T) {
	type T struct{ v int }

	t.Run("one element no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return []T{{v: 1}}, nil } // nolint:unparam

		// --- When ---
		have := First(fn())

		// --- Then ---
		affirm.Equal(t, T{v: 1}, have)
	})

	t.Run("zero elements no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return nil, nil } // nolint:unparam

		// --- When ---
		have := First(fn())

		// --- Then ---
		affirm.Equal(t, T{}, have)
	})

	t.Run("error not nil", func(t *testing.T) {
		// --- Given ---
		var msg string
		check := func() {
			defer func() {
				if r := recover(); r != nil {
					msg = r.(error).Error()
				}
			}()
			fn := func() ([]T, error) { return nil, errors.New("test msg") }
			First(fn())
		}

		// --- When ---
		check()

		// --- Then ---
		affirm.Equal(t, "test msg", msg)
	})

	t.Run("more than one element no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { // nolint:unparam
			return []T{{v: 1}, {v: 2}}, nil
		}

		// --- When ---
		have := First(fn())

		// --- Then ---
		affirm.Equal(t, T{v: 1}, have)
	})
}

func Test_Single(t *testing.T) {
	type T struct{ v int }

	t.Run("one element no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return []T{{v: 1}}, nil } // nolint:unparam

		// --- When ---
		have := Single(fn())

		// --- Then ---
		affirm.Equal(t, T{v: 1}, have)
	})

	t.Run("zero elements no error", func(t *testing.T) {
		// --- Given ---
		fn := func() ([]T, error) { return nil, nil } // nolint:unparam

		// --- When ---
		have := Single(fn())

		// --- Then ---
		affirm.Equal(t, T{}, have)
	})

	t.Run("error not nil", func(t *testing.T) {
		// --- Given ---
		var msg string
		check := func() {
			defer func() {
				if r := recover(); r != nil {
					msg = r.(error).Error()
				}
			}()
			fn := func() ([]T, error) { return nil, errors.New("test msg") }
			Single(fn())
		}

		// --- When ---
		check()

		// --- Then ---
		affirm.Equal(t, "test msg", msg)
	})

	t.Run("more than one element no error", func(t *testing.T) {
		// --- Given ---
		var e error
		check := func() {
			defer func() {
				if r := recover(); r != nil {
					e = r.(error)
				}
			}()
			s := []T{{v: 1}, {v: 2}}
			fn := func() ([]T, error) { return s, nil } // nolint: unparam
			Single(fn())
		}

		// --- When ---
		check()

		// --- Then ---
		affirm.Equal(t, true, errors.Is(e, errExpSingle))
	})
}
