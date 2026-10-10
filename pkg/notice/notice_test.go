// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package notice

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
)

func Test_New(t *testing.T) {
	t.Run("with args", func(t *testing.T) {
		// --- When ---
		have := New("header %s", "row")

		// --- Then ---
		affirm.Equal(t, "header row", have.HeaderText)
		affirm.Equal(t, "", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		affirm.Nil(t, have.Rows)
		affirm.Nil(t, have.Meta)
		affirm.Equal(t, true, errors.Is(have.err, ErrNotice))
		affirm.Nil(t, have.prev)
		affirm.Nil(t, have.next)
	})

	t.Run("with percent but no args", func(t *testing.T) {
		// --- When ---
		have := New("header %s")

		// --- Then ---
		affirm.Equal(t, "header %s", have.HeaderText)
		affirm.Equal(t, "", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		affirm.Nil(t, have.Rows)
		affirm.Nil(t, have.Meta)
		affirm.Equal(t, true, errors.Is(have.err, ErrNotice))
		affirm.Nil(t, have.prev)
		affirm.Nil(t, have.next)
	})
}

func Test_From(t *testing.T) {
	t.Run("typed nil notice", func(t *testing.T) {
		// --- Given ---
		var msg *Notice

		// --- When ---
		have := From(msg, "prefix")

		// --- Then ---
		affirm.Equal(t, true, have == nil)
	})

	t.Run("prefix on whole chain", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		_ = Join(msg0, msg1)

		// --- When ---
		have := From(msg0, "prefix")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg0, have))
		affirm.Equal(t, "prefix", msg0.HeaderPrefix)
		affirm.Equal(t, "prefix", msg1.HeaderPrefix)
	})

	t.Run("not a notice", func(t *testing.T) {
		// --- When ---
		have := From(errors.New("boom"))

		// --- Then ---
		affirm.Equal(t, "assertion error:\n  error: boom", have.Error())
	})

	t.Run("without prefix", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row").Append("first", "%d", 1)

		// --- When ---
		have := From(msg).Append("second", "%d", 2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header row", have.HeaderText)
		affirm.Equal(t, "", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, have.Rows)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
	})

	t.Run("with prefix", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row").Append("first", "%d", 1)

		// --- When ---
		have := From(msg, "prefix").Append("second", "%d", 2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header row", have.HeaderText)
		affirm.Equal(t, "prefix", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, have.Rows)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
	})

	t.Run("override prefix", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row").Append("first", "%d", 1)
		msg = From(msg, "prefix 0")

		// --- When ---
		have := From(msg, "prefix 1").Append("second", "%d", 2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header row", have.HeaderText)
		affirm.Equal(t, "prefix 1", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, have.Rows)
		affirm.Equal(t, true, errors.Is(have, ErrNotice))
	})

	t.Run("prefix on chain of two", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		_ = Join(msg0, msg1)

		// --- When ---
		have := From(msg1, "prefix")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
		affirm.Equal(t, "prefix", msg1.HeaderPrefix)
		affirm.Equal(t, "prefix", msg0.HeaderPrefix)
	})

	t.Run("prefix on chain of three", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		msg2 := New("header2")
		_ = Join(msg0, msg1, msg2)

		// --- When ---
		have := From(msg2, "prefix")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg2, have))
		affirm.Equal(t, "prefix", msg2.HeaderPrefix)
		affirm.Equal(t, "prefix", msg1.HeaderPrefix)
		affirm.Equal(t, "prefix", msg0.HeaderPrefix)
	})

	t.Run("not an instance of Notice without a prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig).Append("first", "%d", 1)

		// --- Then ---
		affirm.Equal(t, false, core.Same(orig, have))
		affirm.Equal(t, "assertion error", have.HeaderText)
		affirm.Equal(t, "", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		want := []Row{
			{Name: "error", Format: "%s", Args: []any{orig}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, want, have.Rows)
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
		affirm.Equal(t, "assertion error", have.HeaderText)
		affirm.Equal(t, "prefix", have.HeaderPrefix)
		affirm.Equal(t, "", have.Trail)
		want := []Row{
			{Name: "error", Format: "%s", Args: []any{orig}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, want, have.Rows)
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
		affirm.Equal(t, "new header", have.HeaderText)
		affirm.Equal(t, "", have.HeaderPrefix)
	})

	t.Run("with args", func(t *testing.T) {
		// --- Given ---
		msg := New("header %s", "row")

		// --- When ---
		have := msg.SetHeader("new header %s", "row")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "new header row", have.HeaderText)
		affirm.Equal(t, "", have.HeaderPrefix)
	})
}

func Test_Notice_SetPrefix(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.SetPrefix("my prefix")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "my prefix", have.HeaderPrefix)
		affirm.Equal(t, "header", have.HeaderText)
	})

	t.Run("with args", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.SetPrefix("my %s", "prefix")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "my prefix", have.HeaderPrefix)
		affirm.Equal(t, "header", have.HeaderText)
	})
}

func Test_Notice_Header(t *testing.T) {
	t.Run("no prefix", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{HeaderText: "expected values to be equal"}

		// --- When ---
		have := msg.Header()

		// --- Then ---
		affirm.Equal(t, "expected values to be equal", have)
	})

	t.Run("with prefix", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{
			HeaderText:   "expected values to be equal",
			HeaderPrefix: "my prefix",
		}

		// --- When ---
		have := msg.Header()

		// --- Then ---
		affirm.Equal(t, "[my prefix] expected values to be equal", have)
	})

	t.Run("empty header with prefix", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{HeaderPrefix: "my prefix"}

		// --- When ---
		have := msg.Header()

		// --- Then ---
		affirm.Equal(t, "[my prefix] ", have)
	})

	t.Run("empty header no prefix", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		have := msg.Header()

		// --- Then ---
		affirm.Equal(t, "", have)
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
		want := []Row{{Name: "first", Format: "%dst", Args: []any{1}}}
		affirm.DeepEqual(t, want, have.Rows)
	})

	t.Run("append second", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%dst", 1)

		// --- When ---
		_ = msg.Append("second", "%dnd", 2)

		// --- Then ---
		want := []Row{
			{Name: "first", Format: "%dst", Args: []any{1}},
			{Name: "second", Format: "%dnd", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
	})

	t.Run("append existing name", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		_ = msg.Append("first", "%s", "abc")

		// --- Then ---
		want := []Row{
			{Name: "first", Format: "%s", Args: []any{"abc"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
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
		want := []Row{{Name: "first", Format: "%dst", Args: []any{1}}}
		affirm.DeepEqual(t, want, msg.Rows)
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
		want := []Row{
			{Name: "first", Format: "%dst", Args: []any{1}},
			{Name: "second", Format: "%dnd", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
	})

	t.Run("append existing name", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1).Append("second", "%d", 2)

		// --- When ---
		_ = msg.AppendRow(NewRow("first", "%d", 3))

		// --- Then ---
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{3}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
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
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
	})

	t.Run("prepend second", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1)

		// --- When ---
		_ = msg.Prepend("second", "%d", 2)

		// --- Then ---
		want := []Row{
			{Name: "second", Format: "%d", Args: []any{2}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
	})

	t.Run("prepend existing name", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Prepend("first", "%d", 1).
			Prepend("second", "%d", 2)

		// --- When ---
		_ = msg.Prepend("second", "%d", 3)

		// --- Then ---
		want := []Row{
			{Name: "second", Format: "%d", Args: []any{3}},
			{Name: "first", Format: "%d", Args: []any{1}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
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
			HeaderText: "header",
			Trail:      "type.field",
			err:        ErrNotice,
		}
		affirm.DeepEqual(t, want, have)
	})

	t.Run("empty trail", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%d", 1)

		// --- When ---
		_ = msg.SetTrail("")

		// --- Then ---
		affirm.Equal(t, "", msg.Trail)
		want := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, want, msg.Rows)
		affirm.Equal(t, "header:\n  first: 1", msg.Error())
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
		affirm.Equal(t, "header", msg.HeaderText)
		affirm.Equal(t, "", msg.HeaderPrefix)
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
		affirm.Equal(t, "header", msg.HeaderText)
		affirm.Equal(t, "", msg.HeaderPrefix)
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
		affirm.Equal(t, "header", msg.HeaderText)
		affirm.Equal(t, "", msg.HeaderPrefix)
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "have", Format: "%s", Args: []any{"row"}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
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
		affirm.Equal(t, "header", msg.HeaderText)
		affirm.Equal(t, "", msg.HeaderPrefix)
		want := []Row{
			{Name: "first", Format: "%d", Args: []any{1}},
			{Name: "have", Format: "%s", Args: []any{"row"}},
			{Name: "second", Format: "%d", Args: []any{2}},
		}
		affirm.DeepEqual(t, want, msg.Rows)
	})
}

func Test_Notice_Wrap(t *testing.T) {
	// --- Given ---
	msg := New("header").Append("first", "%d", 1)
	errMy := errors.New("my-error")

	// --- When ---
	have := msg.Wrap(errMy)

	// --- Then ---
	affirm.Equal(t, true, core.Same(msg, have))
	affirm.Equal(t, false, errors.Is(msg, ErrNotice))
	affirm.Equal(t, true, errors.Is(msg, errMy))
	want := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
	affirm.DeepEqual(t, want, msg.Rows)
}

func Test_Notice_Unwrap(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Unwrap()

		// --- Then ---
		affirm.Equal(t, true, core.Same(ErrNotice, have))
	})

	t.Run("wrapped", func(t *testing.T) {
		// --- Given ---
		var errMy = errors.New("my-error")
		msg := New("header").Append("first", "%d", 1).Wrap(errMy)

		// --- When ---
		have := msg.Unwrap()

		// --- Then ---
		affirm.Equal(t, true, core.Same(errMy, have))
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
		want := []Row{{Name: "second", Format: "%d", Args: []any{2}}}
		affirm.DeepEqual(t, want, msg.Rows)
	})

	t.Run("remove not existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%d", 1)

		// --- When ---
		have := msg.Remove("second")

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		want := []Row{{Name: "first", Format: "%d", Args: []any{1}}}
		affirm.DeepEqual(t, want, msg.Rows)
	})
}

func Test_Notice_Is(t *testing.T) {
	t.Run("base error wrapping a chain notice", func(t *testing.T) {
		// --- Given ---
		msg := New("header")
		joined := Join(msg, fmt.Errorf("context: %w", msg))

		// --- When ---
		have := errors.Is(joined, io.EOF)

		// --- Then ---
		affirm.Equal(t, false, have)

		affirm.Equal(t, true, errors.Is(joined, ErrNotice))
	})

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
		msg := New("header")
		err := errors.New("my-error")

		// --- When ---
		have := msg.Is(err)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("matches wrapped error on earlier chain notice", func(t *testing.T) {
		// --- Given ---
		errA := errors.New("err-a")
		errB := errors.New("err-b")
		joined := Join(New("one").Wrap(errA), New("two").Wrap(errB))

		// --- When ---
		have := errors.Is(joined, errA)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("matches wrapped error on the tail notice", func(t *testing.T) {
		// --- Given ---
		errA := errors.New("err-a")
		errB := errors.New("err-b")
		joined := Join(New("one").Wrap(errA), New("two").Wrap(errB))

		// --- When ---
		have := errors.Is(joined, errB)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("matches a notice on the chain", func(t *testing.T) {
		// --- Given ---
		one := New("one")
		joined := Join(one, New("two"))

		// --- When ---
		have := errors.Is(joined, one)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("does not match a notice off the chain", func(t *testing.T) {
		// --- Given ---
		joined := Join(New("one"), New("two"))

		// --- When ---
		have := errors.Is(joined, New("one"))

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

	t.Run("with prefix", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal")
		msg.HeaderPrefix = "prefix"

		// --- When ---
		have := msg.Error()

		// --- Then ---
		affirm.Equal(t, "[prefix] expected values to be equal", have)
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

	t.Run("leading newline in row", func(t *testing.T) {
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

	t.Run("with rows and trail", func(t *testing.T) {
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
		affirm.DeepEqual(t, map[string]any{"A": 0}, have.Meta)
	})

	t.Run("overwrite", func(t *testing.T) {
		// --- Given ---
		msg := New("header").MetaSet("A", 0)

		// --- When ---
		have := msg.MetaSet("A", 1)

		// --- Then ---
		affirm.DeepEqual(t, map[string]any{"A": 1}, have.Meta)
	})
}

func Test_Notice_MetaLookup(t *testing.T) {
	t.Run("get existing key", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{Meta: map[string]any{"A": 0}}

		// --- When ---
		hVal, hOK := msg.MetaLookup("A")

		// --- Then ---
		affirm.Equal(t, 0, hVal)
		affirm.Equal(t, true, hOK)
	})

	t.Run("get not existing key", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{Meta: map[string]any{"A": 0}}

		// --- When ---
		hVal, hOK := msg.MetaLookup("B")

		// --- Then ---
		affirm.Nil(t, hVal)
		affirm.Equal(t, false, hOK)
	})

	t.Run("nil meta", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		hVal, hOK := msg.MetaLookup("A")

		// --- Then ---
		affirm.Nil(t, hVal)
		affirm.Equal(t, false, hOK)
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
		affirm.Equal(t, 3, have)
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
		affirm.Equal(t, 5, have)
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
		affirm.Equal(t, 9, have)
	})
}

func Test_Notice_Chain(t *testing.T) {
	t.Run("prev among successors", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		msg2 := New("header2")
		_ = Join(msg0, msg1, msg2)

		// --- When ---
		have := msg1.Chain(msg2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
		affirm.DeepEqual(t, []*Notice{msg0, msg1, msg2}, msg0.All())
		affirm.Nil(t, msg2.next)
	})

	t.Run("links two notices", func(t *testing.T) {
		// --- Given ---
		msg1 := New("header1").Want("%d", 11).Have("%d", 7)
		msg0 := New("header0").Want("%d", 42).Have("%d", 44)

		// --- When ---
		have := msg1.Chain(msg0)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
		affirm.Nil(t, msg0.prev)
		affirm.Equal(t, true, core.Same(msg0.next, msg1))
		affirm.Equal(t, true, core.Same(msg1.prev, msg0))
		affirm.Nil(t, msg1.next)
	})

	t.Run("rejoin", func(t *testing.T) {
		// --- Given ---
		msgA := New("a")
		msgB := New("b")
		_ = msgB.Chain(msgA)

		msgC := New("c")

		// --- When ---
		// Chain onto msgA replaces its successor; Join would append after
		// the tail so multi-notice operands stay intact.
		_ = msgC.Chain(msgA)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgA.next, msgC))
		affirm.Equal(t, true, core.Same(msgC.prev, msgA))
		affirm.Nil(t, msgB.prev)
		affirm.Nil(t, msgB.next)
		affirm.Equal(t, true, core.Same(msgB.Head(), msgB))
		affirm.Equal(t, 1, len(msgB.All()))
	})

	t.Run("self chain", func(t *testing.T) {
		// --- Given ---
		msg := New("a")

		// --- When ---
		have := msg.Chain(msg)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Nil(t, msg.prev)
		affirm.Nil(t, msg.next)
	})

	t.Run("msg with successors", func(t *testing.T) {
		// --- Given ---
		msgA := New("a")
		msgB := New("b")
		msgC := New("c")
		_ = Join(msgA, msgB, msgC)

		msgX := New("x")

		// --- When ---
		have := msgB.Chain(msgX)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgB, have))
		affirm.Nil(t, msgA.next)
		affirm.Equal(t, true, core.Same(msgX.next, msgB))
		affirm.Equal(t, true, core.Same(msgB.prev, msgX))
		affirm.Equal(t, true, core.Same(msgB.next, msgC))
		affirm.Equal(t, true, core.Same(msgC.prev, msgB))
		affirm.DeepEqual(t, []*Notice{msgX, msgB, msgC}, msgX.All())
	})
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

func Test_Notice_All(t *testing.T) {
	t.Run("without parent", func(t *testing.T) {
		// --- Given ---
		msg := &Notice{}

		// --- When ---
		have := msg.All()

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
		have := msg2.All()

		// --- Then ---
		affirm.DeepEqual(t, []*Notice{msg0, msg1, msg2}, have)
	})
}

func Test_Join(t *testing.T) {
	t.Run("same chain twice", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		joined := Join(msg0, msg1)

		// --- When ---
		have := Join(joined, joined)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
		affirm.DeepEqual(t, []*Notice{msg0, msg1}, msg0.All())
	})

	t.Run("rejoin already joined notices", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header0")
		msg1 := New("header1")
		msg2 := New("header2")
		_ = Join(msg0, msg1, msg2)

		// --- When ---
		have := Join(msg0, msg1, msg2)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg2, have))
		affirm.DeepEqual(t, []*Notice{msg0, msg1, msg2}, msg0.All())
	})

	t.Run("typed nil notice", func(t *testing.T) {
		// --- Given ---
		var msg0 *Notice
		msg1 := New("header1")

		// --- When ---
		have := Join(msg0, msg1)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg1, have))
	})

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
		affirm.Equal(t, true, have == nil)
	})

	t.Run("multiple nil errors", func(t *testing.T) {
		// --- When ---
		have := Join(nil, nil, nil)

		// --- Then ---
		affirm.Equal(t, true, have == nil)
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

	t.Run("joins multi-notice operands", func(t *testing.T) {
		// --- Given ---
		a0 := New("a0")
		a1 := New("a1")
		b0 := New("b0")
		b1 := New("b1")
		left := Join(a0, a1)
		right := Join(b0, b1)

		// --- When ---
		have := Join(left, right)

		// --- Then ---
		affirm.Equal(t, true, core.Same(b1, have))
		affirm.DeepEqual(t, []*Notice{a0, a1, b0, b1}, have.(*Notice).All())
	})

	t.Run("skip nil errors", func(t *testing.T) {
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

func Test_TrailCmp_tabular(t *testing.T) {
	tt := []struct {
		testN string

		a    string
		b    string
		want int
	}{
		{"both empty", "", "", 0},
		{"empty before not empty", "", "A", -1},
		{"not empty after empty", "A", "", 1},
		{"equal", "ABC", "ABC", 0},
		{"less", "ABC", "XYZ", -1},
		{"greater", "XYZ", "ABC", 1},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			msgA := New("header").SetTrail(tc.a)
			msgB := New("header").SetTrail(tc.b)

			// --- When ---
			have := TrailCmp(msgA, msgB)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_SortNotices(t *testing.T) {
	t.Run("cmp returning any sign", func(t *testing.T) {
		// --- Given ---
		msgC := New("hdrC").SetTrail("ccc")
		msgA := New("hdrA").SetTrail("a")
		msgB := New("hdrB").SetTrail("bb")
		_ = Join(msgA, msgC, msgB)

		cmp := func(x, y *Notice) int { return len(x.Trail) - len(y.Trail) }

		// --- When ---
		have := SortNotices(msgA, cmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		wFwd := "| hdrA (a) -> hdrB (bb) -> hdrC (ccc)"
		affirm.Equal(t, wFwd, FWD(have.Head()))
	})

	t.Run("not a head node", func(t *testing.T) {
		// --- Given ---
		msgB := New("hdrB").SetTrail("B")
		msgC := New("hdrC").SetTrail("C")
		msgA := New("hdrA").SetTrail("A")
		_ = Join(msgB, msgC, msgA)

		// --- When ---
		have := SortNotices(msgC, TrailCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		affirm.Equal(t, "| hdrA (A) -> hdrB (B) -> hdrC (C)", FWD(have.Head()))
		affirm.Equal(t, "| hdrC (C) -> hdrB (B) -> hdrA (A)", REV(have))
	})

	t.Run("nil chain", func(t *testing.T) {
		// --- Given ---
		var head *Notice

		// --- When ---
		have := SortNotices(head, TrailCmp)

		// --- Then ---
		affirm.Nil(t, have)
	})

	t.Run("single node", func(t *testing.T) {
		// --- Given ---
		msg := New("header").SetTrail("a")

		// --- When ---
		have := SortNotices(msg, TrailCmp)

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
		have := SortNotices(msgA, TrailCmp)

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
		have := SortNotices(msgB, TrailCmp)

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
		have := SortNotices(msgA0, TrailCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgA2, have))
		affirm.Equal(t, "| hdr0 (A) -> hdr1 (A) -> hdr2 (A)", FWD(have.Head()))
		affirm.Equal(t, "| hdr2 (A) -> hdr1 (A) -> hdr0 (A)", REV(have))
	})

	t.Run("chain starts with equal sort values", func(t *testing.T) {
		// --- Given ---
		msgA0 := New("hdr0").SetTrail("A")
		msgA1 := New("hdr1").SetTrail("A")
		msgB := New("hdrB").SetTrail("B")
		_ = Join(msgA0, msgA1, msgB)

		// --- When ---
		have := SortNotices(msgA0, TrailCmp)

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
		have := SortNotices(msgC, TrailCmp)

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
		have := SortNotices(msgC, TrailCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		wFwd := "| hdrA (A) -> hdrB0 (B) -> hdrB1 (B) -> hdrC (C)"
		affirm.Equal(t, wFwd, FWD(have.Head()))
		wRev := "| hdrC (C) -> hdrB1 (B) -> hdrB0 (B) -> hdrA (A)"
		affirm.Equal(t, wRev, REV(have))
	})

	t.Run("nodes unsorted with repetitions 2", func(t *testing.T) {
		// --- Given ---
		msgA := New("hdrA").SetTrail("A")
		msgB0 := New("hdrB0").SetTrail("B")
		msgB1 := New("hdrB1").SetTrail("B")
		msgC := New("hdrC").SetTrail("C")
		_ = Join(msgB0, msgC, msgA, msgB1)

		// --- When ---
		have := SortNotices(msgB0, TrailCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgC, have))
		wFwd := "| hdrA (A) -> hdrB0 (B) -> hdrB1 (B) -> hdrC (C)"
		affirm.Equal(t, wFwd, FWD(have.Head()))
		wRev := "| hdrC (C) -> hdrB1 (B) -> hdrB0 (B) -> hdrA (A)"
		affirm.Equal(t, wRev, REV(have))
	})

	t.Run("nodes with empty sort values", func(t *testing.T) {
		// --- Given ---
		msg0 := New("hdr0").SetTrail("")
		msgA := New("hdrA").SetTrail("A")
		msg1 := New("hdr2").SetTrail("")
		_ = Join(msg0, msgA, msg1)

		// --- When ---
		have := SortNotices(msg0, TrailCmp)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msgA, have))
		affirm.Equal(t, "| hdr0 () -> hdr2 () -> hdrA (A)", FWD(have.Head()))
		affirm.Equal(t, "| hdrA (A) -> hdr2 () -> hdr0 ()", REV(have))
	})
}

// Benchmarks for the notice error model and hot formatting paths.

func Benchmark_Notice_Error_Simple(b *testing.B) {
	msg := New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz")

	b.ReportAllocs()
	for b.Loop() {
		_ = msg.Error()
	}
}

func Benchmark_Notice_Join_And_All(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		joined := Join(New("first"), New("second"), New("third"))
		if n, ok := errors.AsType[*Notice](joined); ok {
			_ = n.All()
		}
	}
}

func Benchmark_Notice_Is_And_Unwrap(b *testing.B) {
	myErr := errors.New("root")
	msg := New("failed").Wrap(myErr)

	b.ReportAllocs()
	for b.Loop() {
		_ = errors.Is(msg, myErr)
		_ = msg.Unwrap()
	}
}
