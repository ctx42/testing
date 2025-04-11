// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_New(t *testing.T) {
	t.Run("with args", func(t *testing.T) {
		// --- When ---
		msg := New("header %s", "row")

		// --- Then ---
		affirm.Equal(t, "header row", msg.Header)
		affirm.True(t, msg.Rows == nil)
		affirm.Equal(t, 0, len(msg.Rows))
		affirm.True(t, errors.Is(msg, ErrNotice))
	})

	t.Run("with percent but no args", func(t *testing.T) {
		// --- When ---
		msg := New("header %s")

		// --- Then ---
		affirm.Equal(t, "header %s", msg.Header)
		affirm.True(t, msg.Rows == nil)
		affirm.Equal(t, 0, len(msg.Rows))
		affirm.True(t, errors.Is(msg, ErrNotice))
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_From(t *testing.T) {
	t.Run("with prefix", func(t *testing.T) {
		// --- Given ---
		orig := New("header %s", "row").Append("first", "%d", 1)

		// --- When ---
		have := From(orig, "prefix").Append("second", "%d", 2)

		// --- Then ---
		affirm.True(t, orig == have)
		affirm.Equal(t, "[prefix] header row", have.Header)
		affirm.True(t, errors.Is(have, ErrNotice))
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, have.Rows)
	})

	t.Run("without prefix", func(t *testing.T) {
		// --- Given ---
		orig := New("header %s", "row").Append("first", "%d", 1)

		// --- When ---
		have := From(orig).Append("second", "%d", 2)

		// --- Then ---
		affirm.True(t, orig == have)
		affirm.Equal(t, "header row", have.Header)
		affirm.True(t, errors.Is(have, ErrNotice))
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, have.Rows)
	})

	t.Run("not instance of Error with prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig, "prefix").Append("first", "%d", 1)

		// --- Then ---
		affirm.True(t, orig != have) // nolint: errorlint
		affirm.Equal(t, "[prefix] assertion error", have.Header)
		affirm.False(t, errors.Is(have, ErrNotice))
		affirm.True(t, errors.Is(have, orig))
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
	})

	t.Run("not instance of Error without prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig).Append("first", "%d", 1)

		// --- Then ---
		affirm.True(t, orig != have) // nolint: errorlint
		affirm.Equal(t, "assertion error", have.Header)
		affirm.False(t, errors.Is(have, ErrNotice))
		affirm.True(t, errors.Is(have, orig))
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Append(t *testing.T) {
	t.Run("append first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Append("first", "%dst", 1)

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := []Row{{Name: "first", Format: "%dst", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
	})

	t.Run("append second", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%dst", 1)

		// --- When ---
		_ = msg.Append("second", "%dnd", 2)

		// --- Then ---
		wRows := []Row{
			{Name: "first", Format: "%dst", Args: []any{1}},
			{Name: "second", Format: "%dnd", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("append existing name changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		_ = msg.Append("first", "%s", "abc")

		// --- Then ---
		wRows := []Row{
			{Name: "first", Format: "%s", Args: []any{"abc"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_AppendRow(t *testing.T) {
	t.Run("append first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.AppendRow(NewRow("first", "%dst", 1))

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := []Row{{Name: "first", Format: "%dst", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("append multiple", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		_ = msg.AppendRow(
			NewRow("first", "%dst", 1),
			NewRow("second", "%dnd", 2),
		)

		// --- Then ---
		wRows := []Row{
			{Name: "first", Format: "%dst", Args: []any{1}},
			{Name: "second", Format: "%dnd", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("append existing name changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		_ = msg.AppendRow(NewRow("first", "%d", 3))

		// --- Then ---
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{3}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Prepend(t *testing.T) {
	t.Run("prepend first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Prepend("first", "%d", 1)

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("prepend second", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1)

		// --- When ---
		_ = msg.Prepend("second", "%d", 2)

		// --- Then ---
		wRows := []Row{
			{Name: "second", Format: "%d", Args: []any{2}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("prepend when trail row exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Trail("type.field")

		// --- When ---
		_ = msg.Prepend("second", "%d", 2)

		// --- Then ---
		wRows := []Row{
			{Name: "trail", Format: "%s", Args: []any{"type.field"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("prepend existing name changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1).Prepend("second", "%d", 2)

		// --- When ---
		_ = msg.Prepend("second", "%d", 3)

		// --- Then ---
		wRows := []Row{
			{Name: "second", Format: "%d", Args: []any{3}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Trail(t *testing.T) {
	t.Run("add as first row", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Trail("type.field")

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := []Row{{Name: trail, Format: "%s", Args: []any{"type.field"}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("add to existing rows", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1)

		// --- When ---
		_ = msg.Trail("type.field")

		// --- Then ---
		wRows := []Row{
			{Name: trail, Format: "%s", Args: []any{"type.field"}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("is not adding empty trails", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1)

		// --- When ---
		_ = msg.Trail("")

		// --- Then ---
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("setting trail again changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Trail("type.field0")

		// --- When ---
		_ = msg.Trail("type.field1")

		// --- Then ---
		wRows := []Row{{Name: trail, Format: "%s", Args: []any{"type.field1"}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Message_Want(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Want("%s", "row")

		// --- Then ---
		//goland:noinspection GoDirectComparisonOfErrors
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "want", Format: "%s", Args: []any{"row"}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("want row already exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Append("first", "%d", 1).
			Want("orig").
			Append("second", "%d", 2)

		// --- When ---
		have := msg.Want("%s", "row")

		// --- Then ---
		//goland:noinspection GoDirectComparisonOfErrors
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "want", Format: "%s", Args: []any{"row"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Have(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Have("%s", "row")

		// --- Then ---
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "have", Format: "%s", Args: []any{"row"}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("have row already exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Append("first", "%d", 1).
			Have("orig").
			Append("second", "%d", 2)

		// --- When ---
		have := msg.Have("%s", "row")

		// --- Then ---
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "have", Format: "%s", Args: []any{"row"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Wrap(t *testing.T) {
	// --- Given ---
	var errMy = errors.New("my-error")
	msg := New("header").Append("first", "%d", 1)

	// --- When ---
	have := msg.Wrap(errMy)

	// --- Then ---
	affirm.True(t, msg == have)
	affirm.False(t, errors.Is(msg, ErrNotice))
	affirm.True(t, errors.Is(msg, errMy))
	wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
	affirm.DeepEqual(t, wRows, msg.Rows)
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Remove(t *testing.T) {
	t.Run("remove existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		have := msg.Remove("first")

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := []Row{{Name: "second", Format: "%d", Args: []any{2}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("remove not existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Remove("second")

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Message_Error(t *testing.T) {
	t.Run("header only", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		affirm.Equal(t, "expected values to be equal", have)
	})

	t.Run("simple message", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal").
			Want("42").
			Have("44")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"expected values to be equal:\n" +
			"  want: 42\n" +
			"  have: 44"
		affirm.Equal(t, want, have)
	})

	t.Run("equalize name lengths", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal").
			Want("%d", 42).
			Append("longer", "%d", 44)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"expected values to be equal:\n" +
			"    want: 42\n" +
			"  longer: 44"
		affirm.Equal(t, want, have)
	})

	t.Run("multi line row value", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal").
			Want("%d", 42).
			Append("longer", "%s", "[]int{\n  0,\n  1,\n  2,\n}")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"expected values to be equal:\n" +
			"    want: 42\n" +
			"  longer:\n" +
			"          []int{\n" +
			"            0,\n" +
			"            1,\n" +
			"            2,\n" +
			"          }"
		affirm.Equal(t, want, have)
	})

	t.Run("force row message to start on the next line", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal").
			Append("first", "%d", 1).
			Append("second", "\n%d", 2).
			Append("third", "%d", 3)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"expected values to be equal:\n" +
			"   first: 1\n" +
			"  second:\n" +
			"          2\n" +
			"   third: 3"
		affirm.Equal(t, want, have)
	})

	t.Run("continuation header", func(t *testing.T) {
		// --- Given ---
		msg := New(ContinuationHeader).Want("%d", 42).Have("%d", 44)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := " ---\n" +
			"  want: 42\n" +
			"  have: 44"
		affirm.Equal(t, want, have)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Notice_SetData_GetData(t *testing.T) {
	t.Run("set data and get data", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.SetData("key", "value")

		// --- Then ---
		affirm.True(t, msg == have)

		val, ok := have.GetData("key")
		affirm.True(t, ok)
		affirm.Equal(t, "value", val)
	})

	t.Run("get not existing key", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		haveVal, haveOK := msg.GetData("key")

		// --- Then ---
		affirm.False(t, haveOK)
		affirm.Nil(t, haveVal)
	})
}

func Test_longestName(t *testing.T) {
	// --- Given ---
	msg := &Notice{
		Rows: []Row{
			{Name: "a"},
			{Name: "aaa"},
			{Name: "aa"},
		},
	}

	// --- When ---
	have := msg.longestName()

	// --- Then ---
	affirm.DeepEqual(t, 3, have)
}
