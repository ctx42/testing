// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
)

func Test_New(t *testing.T) {
	t.Run("with args", func(t *testing.T) {
		t.Run("with args", func(t *testing.T) {
			// --- When ---
			msg := New("header %s", "row")

			// --- Then ---
			affirm.Equal(t, "header row", msg.Header)
			affirm.Equal(t, "", msg.Trail)
			affirm.Nil(t, msg.Rows)
			affirm.Nil(t, msg.Meta)
			affirm.Equal(t, true, errors.Is(msg.err, ErrNotice))
			affirm.Nil(t, msg.prev)
			affirm.Nil(t, msg.next)
		})

	})

	t.Run("with percent but no args", func(t *testing.T) {
		// --- When ---
		msg := New("header %s")

		// --- Then ---
		affirm.Equal(t, "header %s", msg.Header)
		affirm.Equal(t, "", msg.Trail)
		affirm.Nil(t, msg.Rows)
		affirm.Nil(t, msg.Meta)
		affirm.Equal(t, true, errors.Is(msg.err, ErrNotice))
		affirm.Nil(t, msg.prev)
		affirm.Nil(t, msg.next)
	})
}

func Test_From(t *testing.T) {
	t.Run("without prefix", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row").Append("first", "%d", 1)

		// --- When ---
		have := From(msg).Append("second", "%d", 2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header row", have.Header)
		affirm.Equal(t, "", have.Trail)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, have.Rows)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
	})

	t.Run("with prefix", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row").Append("first", "%d", 1)

		// --- When ---
		have := From(msg, "prefix").Append("second", "%d", 2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "[prefix] header row", have.Header)
		affirm.Equal(t, "", have.Trail)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, have.Rows)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
	})

	t.Run("not an instance of Notice without a prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig).Append("first", "%d", 1)

		// --- Then ---
		affirm.Equal(t, false, core.Same(orig, have))
		affirm.Equal(t, "assertion error", have.Header)
		affirm.Equal(t, "", have.Trail)
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
		affirm.Equal(t, true, errors.Is(have, orig))
		affirm.Equal(t, false, errors.Is(have, ErrNotice))
	})

	t.Run("not an instance of Notice with a prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig, "prefix").Append("first", "%d", 1)

		// --- Then ---
		affirm.Equal(t, false, core.Same(orig, have))
		affirm.Equal(t, "[prefix] assertion error", have.Header)
		affirm.Equal(t, "", have.Trail)
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, have.Rows)
		affirm.Equal(t, true, errors.Is(have, orig))
		affirm.Equal(t, false, errors.Is(have, ErrNotice))
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		have := From(nil)

		// --- Then ---
		affirm.Nil(t, have)
	})
}

func Test_Notice_SetHeader(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row")

		// --- When ---
		have := msg.SetHeader("new header")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "new header", have.Header)
	})

	t.Run("with args", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row")

		// --- When ---
		have := msg.SetHeader("new header %s", "row")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "new header row", have.Header)
	})
}

func Test_Notice_Append(t *testing.T) {
	t.Run("append first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Append("first", "%dst", 1)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
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

	t.Run("append an existing name overwrites", func(t *testing.T) {
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

func Test_Notice_AppendRow(t *testing.T) {
	t.Run("append first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.AppendRow(NewRow("first", "%dst", 1))

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
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

	t.Run("append an existing name overwrites", func(t *testing.T) {
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

func Test_Notice_Prepend(t *testing.T) {
	t.Run("prepend first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Prepend("first", "%d", 1)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
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

func Test_Notice_SetTrail(t *testing.T) {
	t.Run("add as first row", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.SetTrail("type.field")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		want := &Notice{
			Header: "header",
			Trail:  "type.field",
			err:    ErrNotice,
		}
		affirm.DeepEqual(t, want, have)
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
}

func Test_Notice_Want(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Want("%s", "row")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "want", Format: "%s", Args: []any{"row"}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("row already exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Append("first", "%d", 1).
			Want("orig").
			Append("second", "%d", 2)

		// --- When ---
		have := msg.Want("%s", "row")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "want", Format: "%s", Args: []any{"row"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Notice_Have(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Have("%s", "row")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "have", Format: "%s", Args: []any{"row"}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("row already exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Append("first", "%d", 1).
			Have("orig").
			Append("second", "%d", 2)

		// --- When ---
		have := msg.Have("%s", "row")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header", msg.Header)
		wRows := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "have", Format: "%s", Args: []any{"row"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Notice_Wrap(t *testing.T) {
	// --- Given ---
	var errMy = errors.New("my-error")
	msg := New("header").Append("first", "%d", 1)

	// --- When ---
	have := msg.Wrap(errMy)

	// --- Then ---
	affirm.Equal(t, true, core.Same(msg, have))
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
		affirm.Equal(t, true, core.Same(ErrNotice, err))
	})

	t.Run("wrapped", func(t *testing.T) {
		// --- Given ---
		var errMy = errors.New("my-error")
		msg := New("header").Append("first", "%d", 1).Wrap(errMy)

		// --- When ---
		err := msg.Unwrap()

		// --- Then ---
		affirm.Equal(t, true, core.Same(errMy, err))
	})
}

func Test_Notice_Remove(t *testing.T) {
	t.Run("remove existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		have := msg.Remove("first")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		wRows := []Row{{Name: "second", Format: "%d", Args: []any{2}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})

	t.Run("remove not existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Remove("second")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		wRows := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, wRows, msg.Rows)
	})
}

func Test_Notice_Is(t *testing.T) {
	t.Run("is", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Is(ErrNotice)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("is not", func(t *testing.T) {
		// --- Given ---
		err := errors.New("my-error")
		msg := New("header")

		// --- When ---
		have := msg.Is(err)

		// --- Then ---
		affirm.Equal(t, false, have)
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

	t.Run("force a row message to start on the next line", func(t *testing.T) {
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

	t.Run("with only trail row", func(t *testing.T) {
		// --- Given ---
		msg := New("header").SetTrail("type.field")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"header:\n" +
			"  trail: type.field"
		affirm.Equal(t, want, have)
	})

	t.Run("with rows and trial", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Want("%d", 42).
			Have("%d", 44).
			SetTrail("type.field")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"header:\n" +
			"  trail: type.field\n" +
			"   want: 42\n" +
			"   have: 44"
		affirm.Equal(t, want, have)
	})

	t.Run("joined", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0").Want("%d", 42).Have("%d", 44)
		msg1 := New("header1").Want("%d", 11).Have("%d", 7)
		msg := Join(msg0, msg1)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"multiple expectations violated:\n" +
			"  error: header0\n" +
			"   want: 42\n" +
			"   have: 44\n" +
			"      ---\n" +
			"  error: header1\n" +
			"   want: 11\n" +
			"   have: 7"
		affirm.Equal(t, want, have)
	})

	t.Run("joined the last has a simple header", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0").Want("%d", 42).Have("%d", 44)
		msg1 := New("header1")
		msg := Join(msg0, msg1)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"multiple expectations violated:\n" +
			"  error: header0\n" +
			"   want: 42\n" +
			"   have: 44\n" +
			"      ---\n" +
			"  error: header1"
		affirm.Equal(t, want, have)
	})

	t.Run("joined the first has a simple header", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1").Want("%d", 11).Have("%d", 7)
		msg := Join(msg0, msg1)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"multiple expectations violated:\n" +
			"  error: header0\n" +
			"      ---\n" +
			"  error: header1\n" +
			"   want: 11\n" +
			"   have: 7"
		affirm.Equal(t, want, have)
	})

	t.Run("joined with empty error row", func(t *testing.T) {
		// --- Given ---
		msg0 := New("").Append("name", "%s", "value")
		msg1 := New("header0").Want("%d", 42).Have("%d", 44)
		msg2 := New("header1").Want("%d", 11).Have("%d", 7)
		msg := Join(msg0, msg1, msg2)

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "" +
			"multiple expectations violated:\n" +
			"   name: value\n" +
			"      ---\n" +
			"  error: header0\n" +
			"   want: 42\n" +
			"   have: 44\n" +
			"      ---\n" +
			"  error: header1\n" +
			"   want: 11\n" +
			"   have: 7"
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

func Test_Notice_longest(t *testing.T) {
	t.Run("empty trail", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{
			Rows: []Row{
				{Name: "a"},
				{Name: "aaa"},
				{Name: "aa"},
			},
		}

		// --- When ---
		have := msg.longest()

		// --- Then ---
		affirm.DeepEqual(t, 3, have)
	})

	t.Run("trail set and all shorter than the trail", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{
			Trail: "type.field",
			Rows: []Row{
				{Name: "a"},
				{Name: "aaa"},
				{Name: "aa"},
			},
		}

		// --- When ---
		have := msg.longest()

		// --- Then ---
		affirm.DeepEqual(t, 5, have)
	})

	t.Run("longer than the trail", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{
			Rows: []Row{
				{Name: "a"},
				{Name: "long-name"},
				{Name: "aa"},
			},
		}

		// --- When ---
		have := msg.longest()

		// --- Then ---
		affirm.DeepEqual(t, 9, have)
	})
}

func Test_Notice_Chain(t *testing.T) {
	// --- Given ---
	msg0 := New("header0").Want("%d", 42).Have("%d", 44)
	msg1 := New("header1").Want("%d", 11).Have("%d", 7)

	// --- When ---
	have := msg1.Chain(msg0)

	// --- Then ---
	affirm.Equal(t, true, core.Same(msg1, have))
	affirm.Nil(t, msg0.prev)
	affirm.Equal(t, true, core.Same(msg0.next, msg1))
	affirm.Equal(t, true, core.Same(msg1.prev, msg0))
	affirm.Nil(t, msg1.next)
}

func Test_Notice_Head(t *testing.T) {
	t.Run("without parent", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		have := msg.Head()

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
	})

	t.Run("finds parent", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		msg2 := New("header2")
		_ = Join(msg0, msg1, msg2)

		// --- When ---
		have := msg2.Head()

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg0, have))
	})
}

func Test_Notice_Next(t *testing.T) {
	t.Run("no next in the chain", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		have := msg.Next()

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("next", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		_ = Join(msg0, msg1)

		// --- When ---
		have := msg0.Next()

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
	})
}

func Test_Notice_Prev(t *testing.T) {
	t.Run("no previous in the chain", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		have := msg.Prev()

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("previous", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		_ = Join(msg0, msg1)

		// --- When ---
		have := msg1.Prev()

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg0, have))
	})
}

func Test_Notice_collect(t *testing.T) {
	t.Run("without parent", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		have := msg.collect()

		// --- Then ---
		affirm.DeepEqual(t, []*Notice{msg}, have)
	})

	t.Run("multiple notices in the chain", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		msg2 := New("header2")
		_ = Join(msg0, msg1, msg2)

		// --- When ---
		have := msg2.collect()

		// --- Then ---
		affirm.DeepEqual(t, []*Notice{msg0, msg1, msg2}, have)
	})
}

func Test_Join(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		// --- When ---
		have := Join()

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("single nil error", func(t *testing.T) {
		// --- When ---
		have := Join(nil)

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("multiple nil errors", func(t *testing.T) {
		// --- When ---
		have := Join(nil, nil, nil)

		// --- Then ---
		isNil, isWrapped := core.IsNil(have)
		affirm.Equal(t, true, isNil)
		affirm.Equal(t, false, isWrapped)
	})

	t.Run("join two errors", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")

		// --- When ---
		have := Join(msg0, msg1)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
		affirm.Equal(t, true, core.Same(msg0.next, msg1))
		affirm.Equal(t, true, core.Same(msg1.prev, msg0))
	})

	t.Run("skip nil errors errors", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")

		// --- When ---
		have := Join(msg0, nil, msg1, nil)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
		affirm.Equal(t, true, core.Same(msg0.next, msg1))
		affirm.Equal(t, true, core.Same(msg1.prev, msg0))
	})
}
