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
		affirm.Equal(t, true, msg.Rows == nil)
		affirm.Equal(t, 0, len(msg.Rows))
		affirm.Equal(t, true, errors.Is(msg, ErrNotice))
	})

	t.Run("with percent but no args", func(t *testing.T) {
		// --- When ---
		msg := New("header %s")

		// --- Then ---
		affirm.Equal(t, "header %s", msg.Header)
		affirm.Equal(t, true, msg.Rows == nil)
		affirm.Equal(t, 0, len(msg.Rows))
		affirm.Equal(t, true, errors.Is(msg, ErrNotice))
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
		affirm.Equal(t, true, orig == have)
		affirm.Equal(t, "[prefix] header row", have.Header)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
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
		affirm.Equal(t, true, orig == have)
		affirm.Equal(t, "header row", have.Header)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
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
		affirm.Equal(t, true, orig != have) // nolint: errorlint
		affirm.Equal(t, "[prefix] assertion error", have.Header)
		affirm.Equal(t, false, errors.Is(have, ErrNotice))
		affirm.Equal(t, true, errors.Is(have, orig))
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
	})

	t.Run("not instance of Error without prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig).Append("first", "%d", 1)

		// --- Then ---
		affirm.Equal(t, true, orig != have) // nolint: errorlint
		affirm.Equal(t, "assertion error", have.Header)
		affirm.Equal(t, false, errors.Is(have, ErrNotice))
		affirm.Equal(t, true, errors.Is(have, orig))
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Notice_Append(t *testing.T) {
	t.Run("append first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Append("first", "%dst", 1)

		// --- Then ---
		affirm.Equal(t, true, msg == have)
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
func Test_Notice_AppendRow(t *testing.T) {
	t.Run("append first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.AppendRow(NewRow("first", "%dst", 1))

		// --- Then ---
		affirm.Equal(t, true, msg == have)
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
func Test_Notice_Prepend(t *testing.T) {
	t.Run("prepend first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Prepend("first", "%d", 1)

		// --- Then ---
		affirm.Equal(t, true, msg == have)
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
		msg := New("header").SetTrail("type.field")

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
func Test_Notice_SetTrail(t *testing.T) {
	t.Run("add as first row", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.SetTrail("type.field")

		// --- Then ---
		affirm.Equal(t, true, msg == have)
		wRows := []Row{{Name: trail, Format: "%s", Args: []any{"type.field"}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("add to existing rows", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1)

		// --- When ---
		_ = msg.SetTrail("type.field")

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
		_ = msg.SetTrail("")

		// --- Then ---
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("setting trail again changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").SetTrail("type.field0")

		// --- When ---
		_ = msg.SetTrail("type.field1")

		// --- Then ---
		wRows := []Row{{Name: trail, Format: "%s", Args: []any{"type.field1"}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Notice_Want(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Want("%s", "row")

		// --- Then ---
		//goland:noinspection GoDirectComparisonOfErrors
		affirm.Equal(t, true, msg == have)
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
		affirm.Equal(t, true, msg == have)
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
func Test_Notice_Have(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Have("%s", "row")

		// --- Then ---
		affirm.Equal(t, true, msg == have)
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
		affirm.Equal(t, true, msg == have)
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
func Test_Notice_Wrap(t *testing.T) {
	// --- Given ---
	var errMy = errors.New("my-error")
	msg := New("header").Append("first", "%d", 1)

	// --- When ---
	have := msg.Wrap(errMy)

	// --- Then ---
	affirm.Equal(t, true, msg == have)
	affirm.Equal(t, false, errors.Is(msg, ErrNotice))
	affirm.Equal(t, true, errors.Is(msg, errMy))
	wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
	affirm.DeepEqual(t, wRows, msg.Rows)
}

func Test_Notice_Unwrap(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		err := msg.Unwrap()

		// --- Then ---
		affirm.Equal(t, true, errors.Is(err, ErrNotice))
	})

	t.Run("wrapped", func(t *testing.T) {
		// --- Given ---
		var errMy = errors.New("my-error")
		msg := New("header").Append("first", "%d", 1).Wrap(errMy)

		// --- When ---
		err := msg.Unwrap()

		// --- Then ---
		affirm.Equal(t, true, errors.Is(err, errMy))
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Notice_Remove(t *testing.T) {
	t.Run("remove existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		have := msg.Remove("first")

		// --- Then ---
		affirm.Equal(t, true, msg == have)
		wRows := []Row{{Name: "second", Format: "%d", Args: []any{2}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("remove not existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Remove("second")

		// --- Then ---
		affirm.Equal(t, true, msg == have)
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Notice_Error(t *testing.T) {
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

func Test_Notice_MetaSet(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.MetaSet("A", 0)

		// --- Then ---
		want := map[string]any{"A": 0}
		affirm.DeepEqual(t, want, have.Meta)
	})

	t.Run("overwrite", func(t *testing.T) {
		// --- Given ---
		msg := New("header").MetaSet("A", 0)

		// --- When ---
		have := msg.MetaSet("A", 1)

		// --- Then ---
		want := map[string]any{"A": 1}
		affirm.DeepEqual(t, want, have.Meta)
	})
}

func Test_Notice_MetaLookup(t *testing.T) {
	t.Run("get existing key", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{Meta: map[string]any{"A": 0}}

		// --- When ---
		haveVal, haveOK := msg.MetaLookup("A")

		// --- Then ---
		affirm.Equal(t, 0, haveVal)
		affirm.Equal(t, true, haveOK)
	})

	t.Run("get not existing key", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{Meta: map[string]any{"A": 0}}

		// --- When ---
		haveVal, haveOK := msg.MetaLookup("B")

		// --- Then ---
		affirm.Nil(t, haveVal)
		affirm.Equal(t, false, haveOK)
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
