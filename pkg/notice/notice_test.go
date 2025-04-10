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
		affirm.True(t, msg.Rows != nil)
		affirm.Equal(t, 0, len(msg.Rows))
		affirm.True(t, msg.Order == nil)
		affirm.True(t, errors.Is(msg, ErrNotice))
	})

	t.Run("with percent but no args", func(t *testing.T) {
		// --- When ---
		msg := New("header %s")

		// --- Then ---
		affirm.Equal(t, "header %s", msg.Header)
		affirm.True(t, msg.Rows != nil)
		affirm.Equal(t, 0, len(msg.Rows))
		affirm.True(t, msg.Order == nil)
		affirm.True(t, errors.Is(msg, ErrNotice))
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_From(t *testing.T) {
	t.Run("with prefix", func(t *testing.T) {
		// --- Given ---
		orig := New("header %s", "row").Append("first", "0")

		// --- When ---
		have := From(orig, "prefix").Append("second", "1")

		// --- Then ---
		affirm.True(t, orig == have)
		affirm.Equal(t, "[prefix] header row", have.Header)
		affirm.True(t, errors.Is(have, ErrNotice))
		wRows := map[string]string{"first": "0", "second": "1"}
		affirm.DeepEqual(t, wRows, have.Rows)
		affirm.DeepEqual(t, []string{"first", "second"}, have.Order)
	})

	t.Run("without prefix", func(t *testing.T) {
		// --- Given ---
		orig := New("header %s", "row").Append("first", "0")

		// --- When ---
		have := From(orig).Append("second", "1")

		// --- Then ---
		affirm.True(t, orig == have)
		affirm.Equal(t, "header row", have.Header)
		affirm.True(t, errors.Is(have, ErrNotice))
		wRows := map[string]string{"first": "0", "second": "1"}
		affirm.DeepEqual(t, wRows, have.Rows)
		affirm.DeepEqual(t, []string{"first", "second"}, have.Order)
	})

	t.Run("not instance of Error with prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig, "prefix").Append("first", "0")

		// --- Then ---
		affirm.True(t, orig != have) // nolint: errorlint
		affirm.Equal(t, "[prefix] assertion error", have.Header)
		affirm.False(t, errors.Is(have, ErrNotice))
		affirm.True(t, errors.Is(have, orig))
		affirm.DeepEqual(t, map[string]string{"first": "0"}, have.Rows)
		affirm.DeepEqual(t, []string{"first"}, have.Order)
	})

	t.Run("not instance of Error without prefix", func(t *testing.T) {
		// --- Given ---
		orig := errors.New("test")

		// --- When ---
		have := From(orig).Append("first", "0")

		// --- Then ---
		affirm.True(t, orig != have) // nolint: errorlint
		affirm.Equal(t, "assertion error", have.Header)
		affirm.False(t, errors.Is(have, ErrNotice))
		affirm.True(t, errors.Is(have, orig))
		affirm.DeepEqual(t, map[string]string{"first": "0"}, have.Rows)
		affirm.DeepEqual(t, []string{"first"}, have.Order)
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
		affirm.DeepEqual(t, map[string]string{"first": "1st"}, msg.Rows)
		affirm.DeepEqual(t, []string{"first"}, msg.Order)
	})

	t.Run("append second", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "%dst", 1)

		// --- When ---
		_ = msg.Append("second", "%dnd", 2)

		// --- Then ---
		wRows := map[string]string{"first": "1st", "second": "2nd"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first", "second"}, msg.Order)
	})

	t.Run("append existing name changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0").Append("second", "1")

		// --- When ---
		_ = msg.Append("first", "2")

		// --- Then ---
		wRows := map[string]string{"first": "2", "second": "1"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"second", "first"}, msg.Order)
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
		affirm.DeepEqual(t, map[string]string{"first": "1st"}, msg.Rows)
		affirm.DeepEqual(t, []string{"first"}, msg.Order)
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
		wRows := map[string]string{"first": "1st", "second": "2nd"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first", "second"}, msg.Order)
	})

	t.Run("append existing name changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0").Append("second", "1")

		// --- When ---
		_ = msg.AppendRow(NewRow("first", "2"))

		// --- Then ---
		wRows := map[string]string{"first": "2", "second": "1"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"second", "first"}, msg.Order)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Prepend(t *testing.T) {
	t.Run("prepend first", func(t *testing.T) {
		// --- Given ---
		msg := New("header")

		// --- When ---
		have := msg.Prepend("first", "%dst", 1)

		// --- Then ---
		affirm.True(t, msg == have)
		affirm.DeepEqual(t, map[string]string{"first": "1st"}, msg.Rows)
		affirm.DeepEqual(t, []string{"first"}, msg.Order)
	})

	t.Run("prepend second", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%dst", 1)

		// --- When ---
		_ = msg.Prepend("second", "%dnd", 2)

		// --- Then ---
		wRows := map[string]string{"first": "1st", "second": "2nd"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"second", "first"}, msg.Order)
	})

	t.Run("prepend when trail row exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Trail("type.field")

		// --- When ---
		_ = msg.Prepend("second", "%dnd", 2)

		// --- Then ---
		wRows := map[string]string{trail: "type.field", "second": "2nd"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{trail, "second"}, msg.Order)
	})

	t.Run("prepend existing name changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "0").Prepend("second", "1")

		// --- When ---
		_ = msg.Prepend("second", "2")

		// --- Then ---
		wRows := map[string]string{"first": "0", "second": "2"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"second", "first"}, msg.Order)
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
		affirm.DeepEqual(t, map[string]string{trail: "type.field"}, msg.Rows)
		affirm.DeepEqual(t, []string{trail}, msg.Order)
	})

	t.Run("add to existing rows", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%dst", 1)

		// --- When ---
		_ = msg.Trail("type.field")

		// --- Then ---
		wRows := map[string]string{trail: "type.field", "first": "1st"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{trail, "first"}, msg.Order)
	})

	t.Run("is not adding empty trails", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Prepend("first", "%dst", 1)

		// --- When ---
		_ = msg.Trail("")

		// --- Then ---
		wRows := map[string]string{"first": "1st"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first"}, msg.Order)
	})

	t.Run("setting trail again changes it", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Trail("type.field[0]")

		// --- When ---
		_ = msg.Trail("type.field[1]")

		// --- Then ---
		wRows := map[string]string{trail: "type.field[1]"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{trail}, msg.Order)
	})
}

func Test_Message_Want(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0")

		// --- When ---
		have := msg.Want("want %s", "row")

		// --- Then ---
		//goland:noinspection GoDirectComparisonOfErrors
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := map[string]string{
			"first": "0",
			"want":  "want row",
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first", "want"}, msg.Order)
	})

	t.Run("want row already exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0").
			Want("orig").
			Append("second", "1")

		// --- When ---
		have := msg.Want("want %s", "row")

		// --- Then ---
		//goland:noinspection GoDirectComparisonOfErrors
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := map[string]string{
			"first":  "0",
			"want":   "want row",
			"second": "1",
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first", "want", "second"}, msg.Order)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Have(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0")

		// --- When ---
		have := msg.Have("have %s", "row")

		// --- Then ---
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := map[string]string{
			"first": "0",
			"have":  "have row",
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first", "have"}, msg.Order)
	})

	t.Run("have row already exists", func(t *testing.T) {
		// --- Given ---
		msg := New("header").
			Append("first", "0").
			Have("orig").
			Append("second", "1")

		// --- When ---
		have := msg.Have("have %s", "row")

		// --- Then ---
		affirm.True(t, msg == have)
		affirm.Equal(t, "header", msg.Header)
		wRows := map[string]string{
			"first":  "0",
			"have":   "have row",
			"second": "1",
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first", "have", "second"}, msg.Order)
	})
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Wrap(t *testing.T) {
	// --- Given ---
	var errMy = errors.New("my-error")
	msg := New("header").Append("first", "0")

	// --- When ---
	have := msg.Wrap(errMy)

	// --- Then ---
	affirm.True(t, msg == have)
	affirm.False(t, errors.Is(msg, ErrNotice))
	affirm.True(t, errors.Is(msg, errMy))
	affirm.DeepEqual(t, map[string]string{"first": "0"}, msg.Rows)
	affirm.DeepEqual(t, []string{"first"}, msg.Order)
}

//goland:noinspection GoDirectComparisonOfErrors
func Test_Message_Remove(t *testing.T) {
	t.Run("remove existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0").Append("second", "1")

		// --- When ---
		have := msg.Remove("first")

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := map[string]string{
			"second": "1",
		}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"second"}, msg.Order)
	})

	t.Run("remove not existing", func(t *testing.T) {
		// --- Given ---
		msg := New("header").Append("first", "0")

		// --- When ---
		have := msg.Remove("second")

		// --- Then ---
		affirm.True(t, msg == have)
		wRows := map[string]string{"first": "0"}
		affirm.DeepEqual(t, wRows, msg.Rows)
		affirm.DeepEqual(t, []string{"first"}, msg.Order)
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
		want := "expected values to be equal:\n" +
			"  want: 42\n" +
			"  have: 44"
		affirm.Equal(t, want, have)
	})

	t.Run("equalize name lengths", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal").
			Want("42").
			Append("longer", "44")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "expected values to be equal:\n" +
			"    want: 42\n" +
			"  longer: 44"
		affirm.Equal(t, want, have)
	})

	t.Run("multi line row value", func(t *testing.T) {
		// --- Given ---
		msg := New("expected values to be equal").
			Want("42").
			Append("longer", "\n[]int{\n  0,\n  1,\n  2,\n}")

		// --- When ---
		have := msg.Error()

		// --- Then ---
		want := "expected values to be equal:\n" +
			"    want: 42\n" +
			"  longer:\n" +
			"          []int{\n" +
			"            0,\n" +
			"            1,\n" +
			"            2,\n" +
			"          }"
		affirm.Equal(t, want, have)
	})

	t.Run("continuation header", func(t *testing.T) {
		// --- Given ---
		msg := New(ContinuationHeader).Want("42").Have("44")

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

func Test_equalizeNames(t *testing.T) {
	// --- Given ---
	msg := New("expected values to be equal").
		Want("42").
		Have("44").
		Append("other A", "A").
		Append("other B", "B")

	// --- When ---
	have := msg.equalizeNames()

	// --- Then ---
	want := []string{
		"   want",
		"   have",
		"other A",
		"other B",
	}
	affirm.DeepEqual(t, want, have)
}
