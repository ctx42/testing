// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
)

func Test_Indent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "")

		// --- Then ---
		affirm.Equal(t, "", have)
	})

	t.Run("one line", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc")

		// --- Then ---
		affirm.Equal(t, "abc", have)
	})

	t.Run("multiple lines", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			" abc\n" +
			" def\n" +
			" ghi"
		affirm.Equal(t, want, have)
	})

	t.Run("empty lines are not indented", func(t *testing.T) {
		// --- When ---
		have := Indent(1, ' ', "abc\ndef\n\nghi")

		// --- Then ---
		want := "" +
			" abc\n" +
			" def\n" +
			"\n" +
			" ghi"
		affirm.Equal(t, want, have)
	})

	t.Run("use tabs", func(t *testing.T) {
		// --- When ---
		have := Indent(1, '\t', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			"\tabc\n" +
			"\tdef\n" +
			"\tghi"
		affirm.Equal(t, want, have)
	})

	t.Run("no ident", func(t *testing.T) {
		// --- When ---
		have := Indent(0, ' ', "abc\ndef\nghi")

		// --- Then ---
		want := "" +
			"abc\n" +
			"def\n" +
			"ghi"
		affirm.Equal(t, want, have)
	})
}

func Test_Unwrap(t *testing.T) {
	t.Run("unwrap multiple", func(t *testing.T) {
		// --- Given ---
		err0 := errors.New("e0")
		err1 := errors.New("e1")
		ers := errors.Join(err0, err1)

		// --- When ---
		have := Unwrap(ers)

		// --- Then ---
		affirm.DeepEqual(t, []error{err0, err1}, have)
	})

	t.Run("unwrap not multi error", func(t *testing.T) {
		// --- Given ---
		err := errors.New("e0")

		// --- When ---
		have := Unwrap(err)

		// --- Then ---
		affirm.DeepEqual(t, []error{err}, have)
	})

	t.Run("unwrap nil", func(t *testing.T) {
		// --- When ---
		have := Unwrap(nil)

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, 0, len(have))
	})
}

func Test_Join(t *testing.T) {
	t.Run("single error", func(t *testing.T) {
		// --- Given ---
		e := errors.New("e")

		// --- When ---
		have := Join(e)

		// --- Then ---
		affirm.True(t, core.Same(e, have))
	})

	t.Run("joined errors", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("e0")
		e1 := errors.New("e1")
		msg := errors.Join(e0, e1)

		// --- When ---
		have := Join(msg)

		// --- Then ---
		affirm.False(t, core.Same(msg, have))
		ers := have.(multi).Unwrap() // nolint: errorlint
		affirm.Equal(t, 2, len(ers))
		affirm.True(t, core.Same(e0, ers[0]))
		affirm.True(t, core.Same(e1, ers[1]))
	})

	t.Run("joined and single errors", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("e0")
		e1 := errors.New("e1")
		ej := errors.Join(e0, e1)
		e2 := errors.New("e2")

		// --- When ---
		have := Join(ej, e2)

		// --- Then ---
		ers := have.(multi).Unwrap() // nolint: errorlint
		affirm.Equal(t, 3, len(ers))
		affirm.True(t, core.Same(e0, ers[0]))
		affirm.True(t, core.Same(e1, ers[1]))
		affirm.True(t, core.Same(e2, ers[2]))
	})

	t.Run("nil errors", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("e0")
		e1 := errors.New("e1")
		ej := errors.Join(e0, e1)
		e2 := errors.New("e2")

		// --- When ---
		have := Join(ej, nil, e2)

		// --- Then ---
		ers := have.(multi).Unwrap() // nolint: errorlint
		affirm.Equal(t, 3, len(ers))
		affirm.True(t, core.Same(e0, ers[0]))
		affirm.True(t, core.Same(e1, ers[1]))
		affirm.True(t, core.Same(e2, ers[2]))
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		have := Join(nil)

		// --- Then ---
		affirm.Nil(t, have)
	})
}

func Test_multi_Error(t *testing.T) {
	t.Run("multiple errors with consecutive headers", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header").Want("%s", "want 0").Have("%s", "have 0")
		msg1 := New("header").Want("%s", "want 1").Have("%s", "have 1")
		me := Join(errors.Join(msg0, msg1))

		// --- When ---
		have := me.Error()

		// --- Then ---
		wMsg := "" +
			"header:\n" +
			"  want: want 0\n" +
			"  have: have 0\n" +
			" ---\n" +
			"  want: want 1\n" +
			"  have: have 1"
		affirm.Equal(t, wMsg, have)
	})

	t.Run("multiple errors without consecutive headers", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header").Want("%s", "want 0").Have("%s", "have 0")
		msg1 := New("other").Want("%s", "want 1").Have("%s", "have 1")
		msg2 := New("header").Want("%s", "want 2").Have("%s", "have 2")
		me := Join(errors.Join(msg0, msg1, msg2))

		// --- When ---
		have := me.Error()

		// --- Then ---
		wMsg := "" +
			"header:\n" +
			"  want: want 0\n" +
			"  have: have 0\n" +
			"\n" +
			"other:\n" +
			"  want: want 1\n" +
			"  have: have 1\n" +
			"\n" +
			"header:\n" +
			"  want: want 2\n" +
			"  have: have 2"
		affirm.Equal(t, wMsg, have)
	})

	t.Run("not notice error", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header").Want("%s", "want 0").Have("%s", "have 0")
		msg1 := errors.New("not notice")
		msg2 := New("header").Want("%s", "want 2").Have("%s", "have 2")
		me := Join(errors.Join(msg0, msg1, msg2))

		// --- When ---
		have := me.Error()

		// --- Then ---
		wMsg := "header:\n" +
			"  want: want 0\n" +
			"  have: have 0\n" +
			"\n" +
			"not notice\n" +
			"\n" +
			"header:\n" +
			"  want: want 2\n" +
			"  have: have 2"
		affirm.Equal(t, wMsg, have)
	})

	t.Run("multiple errors serialized multiple times", func(t *testing.T) {
		// --- Given ---
		msg0 := New("header").Want("%s", "want 0").Have("%s", "have 0")
		msg1 := New("header").Want("%s", "want 1").Have("%s", "have 1")
		me := Join(errors.Join(msg0, msg1))

		// --- When ---
		have0 := me.Error()
		have1 := me.Error()

		// --- Then ---
		wMsg := "" +
			"header:\n" +
			"  want: want 0\n" +
			"  have: have 0\n" +
			" ---\n" +
			"  want: want 1\n" +
			"  have: have 1"
		affirm.Equal(t, wMsg, have0)
		affirm.Equal(t, have0, have1)
	})
}
