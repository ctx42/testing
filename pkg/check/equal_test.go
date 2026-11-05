// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Equal(t *testing.T) {
	t.Run("typed constant", func(t *testing.T) {
		// --- Given ---
		type MyInt int
		const MyIntValue MyInt = 42

		m0 := map[string]any{"A": MyIntValue}
		m1 := map[string]any{"A": MyIntValue}

		// --- When ---
		err := Equal(m0, m1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - the same underlying type", func(t *testing.T) {
		// --- Given ---
		type MyInt int
		const MyIntValue MyInt = 42

		m0 := map[string]any{"A": MyIntValue}
		m1 := map[string]any{"A": 42}

		// --- When ---
		err := Equal(m0, m1)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: map[\"A\"]\n" +
			"  want type: check.MyInt\n" +
			"  have type: int"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("WithCmpBaseTypes - the same underlying type", func(t *testing.T) {
		// --- Given ---
		type MyInt int
		const MyIntValue MyInt = 42

		m0 := map[string]any{"A": MyIntValue}
		m1 := map[string]any{"A": 42}

		// --- When ---
		err := Equal(m0, m1, WithCmpBaseTypes())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - WithCmpBaseTypes different base types", func(t *testing.T) {
		// --- When ---
		err := Equal(1, int64(1), WithCmpBaseTypes())

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: int\n" +
			"  have type: int64"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("WithOptions not modified", func(t *testing.T) {
		// --- Given ---
		ops := DefaultOptions(WithTrailChecker("a", Exact))
		src := DefaultOptions(WithOptions(ops))

		// --- When ---
		err := Equal(1, 1, WithOptions(src), WithTypeChecker(1, nil))

		// --- Then ---
		affirm.Nil(t, err)
		_, ok := src.TypeCheckers[reflect.TypeOf(1)]
		affirm.Equal(t, false, ok)
		_, ok = src.Dumper.Dumpers[typByte]
		affirm.Equal(t, false, ok)
	})

	t.Run("WithOptions with zero value options", func(t *testing.T) {
		// --- When ---
		err := Equal(1, 1, WithOptions(Options{}))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - map with interface keys", func(t *testing.T) {
		// --- Given ---
		want := map[any]int{"a": 1, "b": 2}
		have := map[any]int{"a": 1, "b": 3}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: map[\"b\"]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Equal_invalid_arguments(t *testing.T) {
	t.Run("equal both are untyped nil", func(t *testing.T) {
		// --- When ---
		err := Equal(nil, nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal untyped nil and nil interface", func(t *testing.T) {
		// --- Given ---
		var itf testcases.TItf

		// --- When ---
		err := Equal(nil, itf)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal nil interface and untyped nil ", func(t *testing.T) {
		// --- Given ---
		var itf testcases.TItf

		// --- When ---
		err := Equal(itf, nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("logs trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrail("type.field"), WithTrailLog(&trail)}

		// --- When ---
		err := Equal(nil, nil, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_one_argument_invalid(t *testing.T) {
	t.Run("error - want is invalid", func(t *testing.T) {
		// --- When ---
		err := Equal(nil, 123)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: nil\n" +
			"  have: 123"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - have is invalid", func(t *testing.T) {
		// --- When ---
		err := Equal(123, nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: 123\n" +
			"  have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - logs trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrail("type.field"), WithTrailLog(&trail)}

		// --- When ---
		err := Equal(123, nil, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: 123\n" +
			"   have: nil"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_not_matching_types(t *testing.T) {
	t.Run("error - not matching", func(t *testing.T) {
		// --- When ---
		err := Equal(123, "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: int\n" +
			"  have type: string"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - logs trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrail("type.field"), WithTrailLog(&trail)}

		// --- When ---
		err := Equal(123, "abc", opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: type.field\n" +
			"  want type: int\n" +
			"  have type: string"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_custom_trail_checkers(t *testing.T) {
	t.Run("custom checker is not used", func(t *testing.T) {
		// --- Given ---
		// Both are defining the same time in different timezone.
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, testcases.WAW)

		trail := make([]string, 0)
		opts := []any{
			WithTrail("type.other"),
			WithTrailLog(&trail),
			WithTrailChecker("type.field", Exact),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.other"}, trail)
	})

	t.Run("error - custom checker used", func(t *testing.T) {
		// --- Given ---
		// Both are defining the same time in different timezone.
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, testcases.WAW)

		trail := make([]string, 0)
		opts := []any{
			WithTrail("type.field"),
			WithTrailLog(&trail),
			WithTrailChecker("type.field", Exact),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  trail: type.field\n" +
			"   want: UTC\n" +
			"   have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_custom_checkers_aliased_pointers(t *testing.T) {
	t.Run("error - trail checker does not hide aliases", func(t *testing.T) {
		// --- Given ---
		w, h := 1, 2
		want := map[string]*int{"a": &w, "b": &w}
		have := map[string]*int{"a": &h, "b": &h}
		chk := func(_, _ any, _ ...any) error { return nil }

		// --- When ---
		err := Equal(want, have, WithTrailChecker(`map["a"]`, chk))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: map[\"b\"]\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nested checker options isolated", func(t *testing.T) {
		// --- Given ---
		type T struct{ V int }
		var calls int
		chk := func(want, have any, opts ...any) error {
			calls++
			opts = append(opts, WithTypeChecker(T{}, nil))
			return Equal(want, have, opts...)
		}

		// --- When ---
		err := Equal([]T{{1}, {2}}, []T{{1}, {2}}, WithTypeChecker(T{}, chk))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, 2, calls)
	})
}

func Test_Equal_custom_type_checkers(t *testing.T) {
	t.Run("error - use the custom type checker", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, testcases.WAW)

		trail := make([]string, 0)
		opts := []any{
			WithTrail("type.field"),
			WithTrailLog(&trail),
			WithTrailChecker("type.field", Exact),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  trail: type.field\n" +
			"   want: UTC\n" +
			"   have: Europe/Warsaw"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - use the custom checker with nils", func(t *testing.T) {
		// --- Given ---
		want, have := 1, 2

		trail := make([]string, 0)
		opts := []any{
			WithTrail("type.field"),
			WithTrailLog(&trail),
			WithTypeChecker(want, func(_, _ any, _ ...any) error {
				return errors.New("custom checker")
			}),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Equal(t, "custom checker", err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("unexported nil pointer with default checker", func(t *testing.T) {
		// --- Given ---
		type T struct{ loc *time.Location }
		var loc *time.Location

		// --- When ---
		err := Equal(T{loc: loc}, T{loc: loc})

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("custom checker gets unexported map", func(t *testing.T) {
		// --- Given ---
		type T struct{ m map[string]int }
		var lens []int
		chk := func(want, have any, _ ...any) error {
			wm, hm := want.(map[string]int), have.(map[string]int)
			lens = append(lens, len(wm), len(hm))
			return nil
		}
		w := T{m: map[string]int{"a": 1}}
		h := T{m: map[string]int{"a": 1, "b": 2}}

		// --- When ---
		err := Equal(w, h, WithTypeChecker(map[string]int{}, chk))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []int{1, 2}, lens)
	})

	t.Run("custom checker gets named simple type", func(t *testing.T) {
		// --- Given ---
		var types []string
		chk := func(want, have any, _ ...any) error {
			types = append(types, fmt.Sprintf("%T %T", want, have))
			return nil
		}
		opt := WithTypeChecker(time.Duration(0), chk)

		// --- When ---
		err := Equal(time.Second, time.Minute, opt)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"time.Duration time.Duration"}, types)
	})
}

func Test_Equal_kind_Ptr(t *testing.T) {
	t.Run("pointers to empty arrays log every trail", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B *[0]int }
		want := T{A: new([0]int), B: new([0]int)}
		have := T{A: new([0]int), B: new([0]int)}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("pointers to empty structs log every trail", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B *struct{} }
		want := T{A: new(struct{}), B: new(struct{})}
		have := T{A: new(struct{}), B: new(struct{})}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("equal structs by pointer", func(t *testing.T) {
		// --- Given ---
		want := &struct{ Int int }{Int: 123}
		have := &struct{ Int int }{Int: 123}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"Int"}, trail)
	})

	t.Run("equal time Location", func(t *testing.T) {
		// --- Given ---
		want := must.Value(time.LoadLocation("Europe/Warsaw"))
		have := must.Value(time.LoadLocation("Europe/Warsaw"))

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("error - not equal time Location", func(t *testing.T) {
		// --- Given ---
		want := must.Value(time.LoadLocation("Europe/Warsaw"))
		have := must.Value(time.LoadLocation("Europe/Paris"))

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  trail: type\n" +
			"   want: Europe/Warsaw\n" +
			"   have: Europe/Paris"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("equal both nil values", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		var want *int
		var have *int

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("error - not equal want is not nil have is nil", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		i := 123
		want := &i
		var have *int

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type\n" +
			"   want: 123\n" +
			"   have: nil"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("error - not equal want is nil have is not nil", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		i := 123
		var want *int
		have := &i

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type\n" +
			"   want: nil\n" +
			"   have: 123"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type"}, trail)
	})
}

func Test_Equal_kind_Struct(t *testing.T) {
	t.Run("empty struct fields log every trail", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B struct{} }

		trail := make([]string, 0)

		// --- When ---
		err := Equal(T{}, T{}, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("equal structs by value", func(t *testing.T) {
		// --- Given ---
		want := struct{ Int int }{Int: 123}
		have := struct{ Int int }{Int: 123}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"Int"}, trail)
	})

	t.Run("equal time struct", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("equal time struct field", func(t *testing.T) {
		// --- Given ---
		want := testcases.TTim{Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)}
		have := testcases.TTim{Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"TTim.Tim"}, trail)
	})

	t.Run("error - not equal time struct", func(t *testing.T) {
		// --- Given ---
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected equal dates:\n" +
			"  trail: type\n" +
			"   want: 2000-01-02T03:04:05Z\n" +
			"   have: 2001-01-02T03:04:05Z\n" +
			"   diff: -8784h0m0s"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal time Location struct value", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type")}

		want := must.Value(time.LoadLocation("Europe/Warsaw"))
		have := must.Value(time.LoadLocation("Europe/Warsaw"))

		// --- When ---
		err := Equal(*want, *have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("equal time Location struct field", func(t *testing.T) {
		// --- Given ---
		waw := must.Value(time.LoadLocation("Europe/Warsaw"))
		want := testcases.TLoc{Loc: waw}
		have := testcases.TLoc{Loc: waw}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"TLoc.Loc"}, trail)
	})

	t.Run("error - not equal time Location nil have", func(t *testing.T) {
		// --- Given ---
		want := testcases.TLoc{Loc: testcases.WAW}
		have := testcases.TLoc{Loc: nil}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected timezones to be equal:\n" +
			"  trail: TLoc.Loc\n" +
			"   want: Europe/Warsaw\n" +
			"   have: UTC"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"TLoc.Loc"}, trail)
	})

	t.Run("equal structs with embedded not struct field", func(t *testing.T) {
		// --- Given ---
		want := testcases.TC{TD: testcases.TD("abc"), Int: 123}
		have := testcases.TC{TD: testcases.TD("abc"), Int: 123}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"TC.TD", "TC.Int"}, trail)
	})

	t.Run("skipped fields are marked", func(t *testing.T) {
		// --- Given ---
		want := testcases.NewTPrv().SetInt(1)
		have := testcases.NewTPrv().SetInt(2)

		trail := make([]string, 0)
		opts := []any{
			WithTrailLog(&trail),
			WithSkipTrail("TPrv.vInt", "TPrv.tim"),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		wTrail := []string{
			"TPrv.Pub",
			"TPrv.vInt <skipped>",
			"TPrv.ptr",
			"TPrv.sInt",
			"TPrv.aInt[0]",
			"TPrv.aInt[1]",
			"TPrv.vMap",
			"TPrv.tim <skipped>",
			"TPrv.fn",
			"TPrv.ch",
		}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("error - not equal structs with multiple errors", func(t *testing.T) {
		// --- Given ---
		want := testcases.TIntStr{Int: 42, Str: "abc"}
		have := testcases.TIntStr{Int: 44, Str: "xyz"}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"multiple expectations violated:\n" +
			"  error: expected values to be equal\n" +
			"  trail: TIntStr.Int\n" +
			"   want: 42\n" +
			"   have: 44\n" +
			"      ---\n" +
			"  error: expected values to be equal\n" +
			"  trail: TIntStr.Str\n" +
			"   want: \"abc\"\n" +
			"   have: \"xyz\""
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"TIntStr.Int", "TIntStr.Str"}, trail)
	})

	t.Run("error - want nil struct pointer", func(t *testing.T) {
		// --- Given ---
		var want *testcases.TA
		have := &testcases.TA{Str: "abc"}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: nil\n" +
			"  have:\n" +
			"        {\n" +
			"          Int: 0,\n" +
			"          Str: \"abc\",\n" +
			"          Tim: \"0001-01-01T00:00:00Z\",\n" +
			"          Dur: \"0s\",\n" +
			"          Loc: nil,\n" +
			"          TAp: nil,\n" +
			"          private: 0,\n" +
			"        }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - have nil struct pointer", func(t *testing.T) {
		// --- Given ---
		want := &testcases.TA{Str: "abc"}
		var have *testcases.TA

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want:\n" +
			"        {\n" +
			"          Int: 0,\n" +
			"          Str: \"abc\",\n" +
			"          Tim: \"0001-01-01T00:00:00Z\",\n" +
			"          Dur: \"0s\",\n" +
			"          Loc: nil,\n" +
			"          TAp: nil,\n" +
			"          private: 0,\n" +
			"        }\n" +
			"  have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not equal deeply nested", func(t *testing.T) {
		// --- Given ---
		want := testcases.TNested{
			STAp: []*testcases.TA{{TAp: &testcases.TA{Int: 42}}},
		}
		have := testcases.TNested{
			STAp: []*testcases.TA{{TAp: &testcases.TA{Int: 44}}},
		}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithSkipUnexported()}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TNested.STAp[0].TAp.Int\n" +
			"   want: 42\n" +
			"   have: 44"
		affirm.Equal(t, wMsg, err.Error())
		wTrail := []string{
			"TNested.SInt",
			"TNested.STA",
			"TNested.STAp[0].Int",
			"TNested.STAp[0].Str",
			"TNested.STAp[0].Tim",
			"TNested.STAp[0].Dur",
			"TNested.STAp[0].Loc",
			"TNested.STAp[0].TAp.Int",
			"TNested.STAp[0].TAp.Str",
			"TNested.STAp[0].TAp.Tim",
			"TNested.STAp[0].TAp.Dur",
			"TNested.STAp[0].TAp.Loc",
			"TNested.STAp[0].TAp.TAp",
			"TNested.STAp[0].TAp.private <skipped>",
			"TNested.STAp[0].private <skipped>",
			"TNested.MStrInt",
			"TNested.MStrTyp",
			"TNested.MIntTyp",
		}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("skip trail checks", func(t *testing.T) {
		// --- Given ---
		want := testcases.TNested{
			STAp: []*testcases.TA{{TAp: &testcases.TA{Int: 42}}},
		}
		have := testcases.TNested{
			STAp: []*testcases.TA{{TAp: &testcases.TA{Int: 44}}},
		}

		trail := make([]string, 0)
		opts := []any{
			WithTrailLog(&trail),
			WithSkipTrail("TNested.STAp[0].TAp.Int"),
			WithSkipUnexported(),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		wTrail := []string{
			"TNested.SInt",
			"TNested.STA",
			"TNested.STAp[0].Int",
			"TNested.STAp[0].Str",
			"TNested.STAp[0].Tim",
			"TNested.STAp[0].Dur",
			"TNested.STAp[0].Loc",
			"TNested.STAp[0].TAp.Int <skipped>",
			"TNested.STAp[0].TAp.Str",
			"TNested.STAp[0].TAp.Tim",
			"TNested.STAp[0].TAp.Dur",
			"TNested.STAp[0].TAp.Loc",
			"TNested.STAp[0].TAp.TAp",
			"TNested.STAp[0].TAp.private <skipped>",
			"TNested.STAp[0].private <skipped>",
			"TNested.MStrInt",
			"TNested.MStrTyp",
			"TNested.MIntTyp",
		}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("error - private int fields not equal", func(t *testing.T) {
		// --- Given ---
		want := testcases.NewTPrv().SetInt(1)
		have := testcases.NewTPrv().SetInt(2)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TPrv.vInt\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
		wTrail := []string{
			"TPrv.Pub",
			"TPrv.vInt",
			"TPrv.ptr",
			"TPrv.sInt",
			"TPrv.aInt[0]",
			"TPrv.aInt[1]",
			"TPrv.vMap",
			"TPrv.tim",
			"TPrv.fn",
			"TPrv.ch",
		}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("equal private function field", func(t *testing.T) {
		// --- Given ---
		fn := func() int { return 0 }
		typ0 := testcases.NewTPrv().SetFn(fn)
		typ1 := testcases.NewTPrv().SetFn(fn)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal private function field", func(t *testing.T) {
		// --- Given ---
		typ0 := testcases.NewTPrv().SetFn(func() int { return 0 })
		typ1 := testcases.NewTPrv().SetFn(func() int { return 1 })

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TPrv.fn\n" +
			"   want: <func>(<addr>)\n" +
			"   have: <func>(<addr>)"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal private pointer field", func(t *testing.T) {
		// --- Given ---
		ptr := &testcases.TVal{Val: "abc"}
		typ0 := testcases.NewTPrv().SetPtr(ptr)
		typ1 := testcases.NewTPrv().SetPtr(ptr)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal private pointer field", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TVal{Val: "abc"}
		ptr1 := &testcases.TVal{Val: "xyz"}
		typ0 := testcases.NewTPrv().SetPtr(ptr0)
		typ1 := testcases.NewTPrv().SetPtr(ptr1)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TPrv.ptr.Val\n" +
			"   want: \"abc\"\n" +
			"   have: \"xyz\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not equal private field nil pointer", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TVal{Val: "abc"}
		typ0 := testcases.NewTPrv().SetPtr(ptr0)
		typ1 := testcases.NewTPrv().SetPtr(nil)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TPrv.ptr\n" +
			"   want:\n" +
			"         {\n" +
			"           Val: \"abc\",\n" +
			"         }\n" +
			"   have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal private slice field", func(t *testing.T) {
		// --- Given ---
		s := []int{1, 2}
		typ0 := testcases.NewTPrv().SetSInt(s)
		typ1 := testcases.NewTPrv().SetSInt(s)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal private slice field", func(t *testing.T) {
		// --- Given ---
		s0 := []int{1, 2}
		s1 := []int{1, 3}
		typ0 := testcases.NewTPrv().SetSInt(s0)
		typ1 := testcases.NewTPrv().SetSInt(s1)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TPrv.sInt[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal private array field", func(t *testing.T) {
		// --- Given ---
		a := [2]int{1, 2}
		typ0 := testcases.NewTPrv().SetAInt(a)
		typ1 := testcases.NewTPrv().SetAInt(a)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal private array field", func(t *testing.T) {
		// --- Given ---
		a0 := [2]int{1, 2}
		a1 := [2]int{1, 3}
		typ0 := testcases.NewTPrv().SetAInt(a0)
		typ1 := testcases.NewTPrv().SetAInt(a1)

		// --- When ---
		err := Equal(typ0, typ1)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TPrv.aInt[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("recursive", func(t *testing.T) {
		// --- Given ---
		s0 := testcases.TRec{Int: 1, Rec: &testcases.TRec{Int: 2}}
		s0.Rec.Rec = &s0

		s1 := testcases.TRec{Int: 1, Rec: &testcases.TRec{Int: 2}}
		s1.Rec.Rec = &s1

		// --- When ---
		err := Equal(s0, s1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("recursive maps", func(t *testing.T) {
		// --- Given ---
		m0 := map[string]any{}
		m0["self"] = m0

		m1 := map[string]any{}
		m1["self"] = m1

		// --- When ---
		err := Equal(m0, m1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("recursive slices", func(t *testing.T) {
		// --- Given ---
		s0 := []any{nil}
		s0[0] = s0

		s1 := []any{nil}
		s1[0] = s1

		// --- When ---
		err := Equal(s0, s1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - first field pointer differs", func(t *testing.T) {
		// --- Given ---
		type node struct {
			Next *node
			Val  int
		}
		n0 := &node{Next: &node{Val: 1}}
		n1 := &node{Next: &node{Val: 2}}

		// --- When ---
		err := Equal(n0, n1)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: node.Next.Val\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Equal_kind_Slice_and_Array(t *testing.T) {
	t.Run("equal slice", func(t *testing.T) {
		// --- Given ---
		want := []int{1, 2}
		have := []int{1, 2}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"<slice>[0]", "<slice>[1]"}, trail)
	})

	t.Run("equal same slice instance", func(t *testing.T) {
		// --- Given ---
		want := []int{1, 2}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal slice value", func(t *testing.T) {
		// --- Given ---
		want := []int{1, 2}
		have := []int{1, 7}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: <slice>[1]\n" +
			"   want: 2\n" +
			"   have: 7"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"<slice>[0]", "<slice>[1]"}, trail)
	})

	t.Run("error - not equal array value", func(t *testing.T) {
		// --- Given ---
		want := [...]int{1, 2}
		have := [...]int{1, 7}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: <array>[1]\n" +
			"   want: 2\n" +
			"   have: 7"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"<array>[0]", "<array>[1]"}, trail)
	})

	t.Run("error - not equal slice lengths", func(t *testing.T) {
		// --- Given ---
		want := []int{1, 2}
		have := []int{1}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"     trail: type.field\n" +
			"  want len: 2\n" +
			"  have len: 1\n" +
			"      want:\n" +
			"            []int{\n" +
			"              1,\n" +
			"              2,\n" +
			"            }\n" +
			"      have:\n" +
			"            []int{\n" +
			"              1,\n" +
			"            }\n" +
			"      diff:\n" +
			"            @@ -1,3 +1,4 @@\n" +
			"             []int{\n" +
			"            -  1,\n" +
			"            +  1,\n" +
			"            +  2,\n" +
			"             }"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - nil want and empty have slice", func(t *testing.T) {
		// --- Given ---
		want := []int(nil)
		have := make([]int, 0)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: nil\n" +
			"   have: []int{}"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - empty want and nil have slice", func(t *testing.T) {
		// --- Given ---
		want := make([]int, 0)
		have := []int(nil)

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: []int{}\n" +
			"  have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - nil and empty nested slice", func(t *testing.T) {
		// --- Given ---
		type T struct{ Items [][]int }
		want := T{Items: [][]int{nil}}
		have := T{Items: [][]int{make([]int, 0)}}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.Items[0]\n" +
			"   want: nil\n" +
			"   have: []int{}"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - nil and empty unexported field", func(t *testing.T) {
		// --- Given ---
		type T struct{ items []int }
		want := T{}
		have := T{items: make([]int, 0)}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.items\n" +
			"   want: nil\n" +
			"   have: []int{}"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("empty slices with different backing arrays", func(t *testing.T) {
		// --- Given ---
		want := make([]int, 0, 1)
		have := make([]int, 0, 2)

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("nil and empty slice with WithEquateEmpty", func(t *testing.T) {
		// --- Given ---
		want := []int(nil)
		have := make([]int, 0)

		// --- When ---
		err := Equal(want, have, WithEquateEmpty())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("nested nil and empty with WithEquateEmpty", func(t *testing.T) {
		// --- Given ---
		type T struct{ Items [][]int }
		want := T{Items: [][]int{nil}}
		have := T{Items: [][]int{make([]int, 0)}}

		// --- When ---
		err := Equal(want, have, WithEquateEmpty())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - nil and non-empty with WithEquateEmpty", func(t *testing.T) {
		// --- Given ---
		want := []int(nil)
		have := []int{1}

		// --- When ---
		err := Equal(want, have, WithEquateEmpty())

		// --- Then ---
		affirm.NotNil(t, err)
	})

	t.Run("empty slices log the trail", func(t *testing.T) {
		// --- Given ---
		want := make([]int, 0, 1)
		have := make([]int, 0, 2)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("empty arrays log the trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal([0]int{}, [0]int{}, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("WithEquateEmpty logs the trail", func(t *testing.T) {
		// --- Given ---
		want := []int(nil)
		have := make([]int, 0)

		trail := make([]string, 0)
		opts := []any{
			WithEquateEmpty(),
			WithTrailLog(&trail),
			WithTrail("type.field"),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("empty slice fields log every trail", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		want := T{A: make([]int, 0), B: make([]int, 0)}
		have := T{A: make([]int, 0), B: make([]int, 0)}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("empty slice elements log every trail", func(t *testing.T) {
		// --- Given ---
		want := [][]int{make([]int, 0), make([]int, 0)}
		have := [][]int{make([]int, 0), make([]int, 0)}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"<slice>[0]", "<slice>[1]"}, trail)
	})

	t.Run("error - empty slice hides no later difference", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		x := []int{1}
		y := []int{2}
		want := T{A: x[:0], B: x}
		have := T{A: y[:0], B: y}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.B[0]\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - each field with empty have is reported", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		x := []int{1}
		y := make([]int, 0)

		trail := make([]string, 0)

		// --- When ---
		err := Equal(T{A: x, B: x}, T{A: y, B: y}, WithTrailLog(&trail))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("error - each field with empty want is reported", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		x := make([]int, 0)
		y := []int{1}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(T{A: x, B: x}, T{A: y, B: y}, WithTrailLog(&trail))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("slices of empty structs log every trail", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []struct{} }
		want := T{A: make([]struct{}, 2), B: make([]struct{}, 2)}
		have := T{A: make([]struct{}, 2), B: make([]struct{}, 2)}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		// The runtime gives these slices one address, so elements are not
		// logged. The Go specification does not promise that.
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("error - slices of empty structs differ in length", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []struct{} }
		want := T{A: make([]struct{}, 2), B: make([]struct{}, 2)}
		have := T{A: make([]struct{}, 2), B: make([]struct{}, 3)}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		// The runtime gives these slices one address, so elements are not
		// logged. The Go specification does not promise that.
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("slices of empty structs at real addresses", func(t *testing.T) {
		// --- Given ---
		type Z struct {
			E [1]struct{}
			X int
		}
		type T struct{ A, B []struct{} }
		wZ, hZ := new(Z), new(Z)
		want := T{A: wZ.E[:], B: wZ.E[:]}
		have := T{A: hZ.E[:], B: hZ.E[:]}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A[0]", "T.B[0]"}, trail)
	})

	t.Run("WithEquateEmpty on unexported field", func(t *testing.T) {
		// --- Given ---
		type T struct{ items []int }
		want := T{}
		have := T{items: make([]int, 0)}

		// --- When ---
		err := Equal(want, have, WithEquateEmpty())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - shorter alias hides no difference", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		x := []int{1, 2}
		y := []int{1, 3}
		want := T{A: x[:1], B: x[:2]}
		have := T{A: y[:1], B: y[:2]}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.B[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - shorter nested alias hides no nil", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B [][]int }
		x := [][]int{{1}, nil}
		y := [][]int{{1}, make([]int, 0)}
		want := T{A: x[:1], B: x[:2]}
		have := T{A: y[:1], B: y[:2]}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.B[1]\n" +
			"   want: nil\n" +
			"   have: []int{}"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - shorter have alias at same address", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		x := []int{1, 2}
		y := []int{1, 2}
		want := T{A: x[:2], B: x[:2]}
		have := T{A: y[:2], B: y[:1]}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A[0]", "T.A[1]", "T.B"}, trail)
	})

	t.Run("error - shorter want alias at same address", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B []int }
		x := []int{1, 2}
		y := []int{1, 2}
		want := T{A: x[:2], B: x[:1]}
		have := T{A: y[:2], B: y[:2]}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A[0]", "T.A[1]", "T.B"}, trail)
	})

	t.Run("cycle through a shorter alias", func(t *testing.T) {
		// --- Given ---
		want := make([]any, 2)
		want[0] = want[:1]
		have := make([]any, 2)
		have[0] = have[:1]

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("nil slices", func(t *testing.T) {
		// --- When ---
		err := Equal([]int(nil), []int(nil))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal slices with multiple errors", func(t *testing.T) {
		// --- Given ---
		want := []int{1, 2}
		have := []int{2, 3}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"multiple expectations violated:\n" +
			"  error: expected values to be equal\n" +
			"  trail: <slice>[0]\n" +
			"   want: 1\n" +
			"   have: 2\n" +
			"      ---\n" +
			"  error: expected values to be equal\n" +
			"  trail: <slice>[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"<slice>[0]", "<slice>[1]"}, trail)
	})
}

func Test_Equal_kind_Map(t *testing.T) {
	t.Run("error - nil want and empty have map", func(t *testing.T) {
		// --- Given ---
		want := map[string]int(nil)
		have := make(map[string]int)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: map[string]int(nil)\n" +
			"   have: map[string]int{}"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - empty want and nil have map", func(t *testing.T) {
		// --- Given ---
		want := make(map[string]int)
		have := map[string]int(nil)

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: map[string]int{}\n" +
			"  have: map[string]int(nil)"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil and empty map with WithEquateEmpty", func(t *testing.T) {
		// --- Given ---
		want := map[string]int(nil)
		have := make(map[string]int)

		// --- When ---
		err := Equal(want, have, WithEquateEmpty())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("empty maps log the trail", func(t *testing.T) {
		// --- Given ---
		want := make(map[string]int)
		have := make(map[string]int)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("WithEquateEmpty logs the map trail", func(t *testing.T) {
		// --- Given ---
		want := map[string]int(nil)
		have := make(map[string]int)

		trail := make([]string, 0)
		opts := []any{
			WithEquateEmpty(),
			WithTrailLog(&trail),
			WithTrail("type.field"),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("empty slice map values log every trail", func(t *testing.T) {
		// --- Given ---
		want := map[string][]int{"a": make([]int, 0), "b": make([]int, 0)}
		have := map[string][]int{"a": make([]int, 0), "b": make([]int, 0)}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{`map["a"]`, `map["b"]`}, trail)
	})

	t.Run("shared empty map logs every trail", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B map[string]int }
		wMap := make(map[string]int)
		hMap := make(map[string]int)
		want := T{A: wMap, B: wMap}
		have := T{A: hMap, B: hMap}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(want, have, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("error - nil and empty unexported map", func(t *testing.T) {
		// --- Given ---
		type T struct{ m map[string]int }
		want := T{}
		have := T{m: make(map[string]int)}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.m\n" +
			"   want: map[string]int(nil)\n" +
			"   have: map[string]int{}"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - each field with empty have is reported", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B map[string]int }
		x := map[string]int{"a": 1}
		y := make(map[string]int)

		trail := make([]string, 0)

		// --- When ---
		err := Equal(T{A: x, B: x}, T{A: y, B: y}, WithTrailLog(&trail))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("error - each field with empty want is reported", func(t *testing.T) {
		// --- Given ---
		type T struct{ A, B map[string]int }
		x := make(map[string]int)
		y := map[string]int{"a": 1}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(T{A: x, B: x}, T{A: y, B: y}, WithTrailLog(&trail))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.DeepEqual(t, []string{"T.A", "T.B"}, trail)
	})

	t.Run("error - NaN key is never found", func(t *testing.T) {
		// --- Given ---
		want := map[float64]int{math.NaN(): 1}
		have := map[float64]int{math.NaN(): 1}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: map[NaN]\n" +
			"  want type: int\n" +
			"  have type: <nil>\n" +
			"       want: 1\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - NaN keys are reported by value", func(t *testing.T) {
		// --- Given ---
		want := make(map[float64]int)
		have := make(map[float64]int)
		for i := range 6 {
			want[math.NaN()] = i
			have[math.NaN()] = i
		}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		found := regexp.MustCompile(`want: (\d)`).
			FindAllStringSubmatch(err.Error(), -1)
		order := make([]string, 0, len(found))
		for _, match := range found {
			order = append(order, match[1])
		}
		affirm.DeepEqual(t, []string{"0", "1", "2", "3", "4", "5"}, order)
	})

	t.Run("error - values with tied keys are ordered once", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }
		want := make(map[float64]T)
		have := make(map[float64]T)
		for i := range 6 {
			want[math.NaN()] = T{N: i}
			have[math.NaN()] = T{N: i}
		}

		var count int
		dpr := func(dump.Dump, int, reflect.Value) string {
			count++
			return "T"
		}

		// --- When ---
		err := Equal(want, have, WithDumper(dump.WithDumper(T{}, dpr)))

		// --- Then ---
		// Each of the six values is dumped to order it and to report it.
		affirm.NotNil(t, err)
		affirm.Equal(t, 12, count)
	})

	t.Run("error - length mismatch prints NaN key value", func(t *testing.T) {
		// --- Given ---
		want := map[float64]int{math.NaN(): 1}
		have := make(map[float64]int)

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		affirm.Equal(t, true, strings.Contains(err.Error(), "NaN: 1,"))
	})

	t.Run("error - cycle under pointer keys", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		defer debug.SetMaxStack(debug.SetMaxStack(64 << 20))

		type T struct{ M map[string]any }
		cyc := make(map[string]any)
		cyc["self"] = cyc
		key0, key1 := &T{M: cyc}, &T{M: cyc}
		want := map[*T]int{key0: 1, key1: 2}
		have := map[*T]int{key0: 1}

		// --- When ---
		err := Equal(want, have, WithDumper(dump.WithMaxDepth(3)))

		// --- Then ---
		affirm.NotNil(t, err)
	})

	t.Run("error - self-referencing pointer key", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		defer debug.SetMaxStack(debug.SetMaxStack(64 << 20))

		var val any
		val = &val
		key := val.(*any)
		want := map[*any]int{key: 1}
		have := map[*any]int{key: 2}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		affirm.Equal(t, true, strings.Contains(err.Error(), dump.ValMaxNesting))
	})

	t.Run("error - missing key in unexported map", func(t *testing.T) {
		// --- Given ---
		type T struct{ m map[string]int }
		want := T{m: map[string]int{"a": 1}}
		have := T{m: map[string]int{"b": 1}}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: T.m[\"a\"]\n" +
			"  want type: int\n" +
			"  have type: <nil>\n" +
			"       want: 1\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil maps", func(t *testing.T) {
		// --- When ---
		err := Equal(map[string]int(nil), map[string]int(nil))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - struct keys are reported in order", func(t *testing.T) {
		// --- Given ---
		type K struct{ A int }
		want := map[K]int{{1}: 1, {2}: 2, {3}: 3}
		have := map[K]int{{1}: 10, {2}: 20, {3}: 30}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		wMsg := "" +
			"multiple expectations violated:\n" +
			"  error: expected values to be equal\n" +
			"  trail: map[check.K]\n" +
			"   want: 1\n" +
			"   have: 10\n" +
			"      ---\n" +
			"  error: expected values to be equal\n" +
			"  trail: map[check.K]\n" +
			"   want: 2\n" +
			"   have: 20\n" +
			"      ---\n" +
			"  error: expected values to be equal\n" +
			"  trail: map[check.K]\n" +
			"   want: 3\n" +
			"   have: 30"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - unexported byte slice is dumped", func(t *testing.T) {
		// --- Given ---
		type T struct{ b []byte }
		want := T{b: []byte{1}}
		have := T{b: []byte{1, 2}}

		// --- When ---
		err := Equal(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		hasCannotPrint := strings.Contains(err.Error(), dump.ValCannotPrint)
		affirm.Equal(t, false, hasCannotPrint)
	})

	t.Run("equal map", func(t *testing.T) {
		// --- Given ---
		want := map[int]int{1: 42}
		have := map[int]int{1: 42}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"map[1]"}, trail)
	})

	t.Run("equal same map", func(t *testing.T) {
		// --- Given ---
		want := map[int]int{1: 42}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal map", func(t *testing.T) {
		// --- Given ---
		want := map[int]int{1: 42}
		have := map[int]int{1: 44}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: map[1]\n" +
			"   want: 42\n" +
			"   have: 44"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"map[1]"}, trail)
	})

	t.Run("error - not equal have map missing keys", func(t *testing.T) {
		// --- Given ---
		want := map[int]int{1: 42, 2: 43}
		have := map[int]int{1: 42, 3: 44}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: map[2]\n" +
			"  want type: int\n" +
			"  have type: <nil>\n" +
			"       want: 43\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"map[1]"}, trail)
	})

	t.Run("error - not equal map length", func(t *testing.T) {
		// --- Given ---
		want := map[int]int{1: 42, 2: 44}
		have := map[int]int{1: 42, 2: 43, 3: 44}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"     trail: type.field\n" +
			"  want len: 2\n" +
			"  have len: 3\n" +
			"      want:\n" +
			"            map[int]int{\n" +
			"              1: 42,\n" +
			"              2: 44,\n" +
			"            }\n" +
			"      have:\n" +
			"            map[int]int{\n" +
			"              1: 42,\n" +
			"              2: 43,\n" +
			"              3: 44,\n" +
			"            }\n" +
			"      diff:\n" +
			"            @@ -1,5 +1,4 @@\n" +
			"             map[int]int{\n" +
			"               1: 42,\n" +
			"            -  2: 43,\n" +
			"            -  3: 44,\n" +
			"            +  2: 44,\n" +
			"             }"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal maps with multiple errors", func(t *testing.T) {
		// --- Given ---
		want := map[int]int{1: 42, 2: 44}
		have := map[int]int{1: 44, 2: 42}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"multiple expectations violated:\n" +
			"  error: expected values to be equal\n" +
			"  trail: type.field[1]\n" +
			"   want: 42\n" +
			"   have: 44\n" +
			"      ---\n" +
			"  error: expected values to be equal\n" +
			"  trail: type.field[2]\n" +
			"   want: 44\n" +
			"   have: 42"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field[1]", "type.field[2]"}, trail)
	})
}

func Test_Equal_kind_Interface(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		want := []any{42}
		have := []any{42}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"<slice>[0]"}, trail)
	})

	t.Run("equal both nil", func(t *testing.T) {
		// --- Given ---
		want := []any{nil}
		have := []any{nil}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"<slice>[0]"}, trail)
	})
}

func Test_Equal_kind_Bool(t *testing.T) {
	t.Run("equal true", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(true, true, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("equal false", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(false, false, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(true, false, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: true\n" +
			"   have: false"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_kind_success_tabular(t *testing.T) {
	type MyInt int
	const MyIntValue MyInt = 42

	tt := []struct {
		testN string

		want any
		have any
	}{
		{"int", 42, 42},
		{"int8", int8(42), int8(42)},
		{"int16", int16(42), int16(42)},
		{"int32", int32(42), int32(42)},
		{"int64", int64(42), int64(42)},

		{"uint", uint(42), uint(42)},
		{"uint8", uint8(42), uint8(42)},
		{"uint16", uint16(42), uint16(42)},
		{"uint32", uint32(42), uint32(42)},
		{"uint64", uint64(42), uint64(42)},

		{"float32", float32(42), float32(42)},
		{"float64", float64(42), float64(42)},

		{"complex64", complex64(42), complex64(42)},
		{"complex128", complex128(42), complex128(42)},

		{"string", "abc", "abc"},

		{"int base type", MyInt(42), MyInt(42)},
		{"const int base type", MyIntValue, MyIntValue},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			trail := make([]string, 0)
			opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

			// --- When ---
			err := Equal(tc.want, tc.have, opts...)

			// --- Then ---
			affirm.Nil(t, err)
			affirm.DeepEqual(t, []string{"type.field"}, trail)
		})
	}
}

func Test_Equal_kind_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		want    any
		have    any
		wantStr string
		haveStr string
	}{
		{"int", 42, 44, "42", "44"},
		{"int8", int8(42), int8(44), "42", "44"},
		{"int16", int16(42), int16(44), "42", "44"},
		{"int32", int32(42), int32(44), "42", "44"},
		{"int64", int64(42), int64(44), "42", "44"},

		{"uint", uint(42), uint(44), "42", "44"},
		{"uint8", uint8(42), uint8(44), "0x2a ('*')", "0x2c (',')"},
		{"uint16", uint16(42), uint16(44), "42", "44"},
		{"uint32", uint32(42), uint32(44), "42", "44"},
		{"uint64", uint64(42), uint64(44), "42", "44"},

		{"float32", float32(42), float32(44), "42", "44"},
		{"float64", float64(42), float64(44), "42", "44"},

		{"complex64", complex64(42), complex64(44), "(42+0i)", "(44+0i)"},
		{"complex128", complex128(42), complex128(44), "(42+0i)", "(44+0i)"},

		{"string", "abc", "xyz", `"abc"`, `"xyz"`},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			trail := make([]string, 0)
			opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

			// --- When ---
			err := Equal(tc.want, tc.have, opts...)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "" +
				"expected values to be equal:\n" +
				"  trail: type.field\n" +
				"   want: %s\n" +
				"   have: %s"
			wMsg = fmt.Sprintf(wMsg, tc.wantStr, tc.haveStr)
			affirm.Equal(t, wMsg, err.Error())
			affirm.DeepEqual(t, []string{"type.field"}, trail)
		})
	}
}

func Test_Equal_kind_Chan(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		want := make(chan bool)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		want := make(chan bool)
		have := make(chan bool)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: (chan bool)(<addr>)\n" +
			"   have: (chan bool)(<addr>)"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal unexported field", func(t *testing.T) {
		// --- Given ---
		want := struct{ want chan bool }{make(chan bool)}
		have := struct{ want chan bool }{make(chan bool)}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field.want\n" +
			"   want: (chan bool)(<addr>)\n" +
			"   have: (chan bool)(<addr>)"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field.want"}, trail)
	})
}

func Test_Equal_kind_Func(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		want := func() {}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		want := func() {}
		have := func() {}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: <func>(<addr>)\n" +
			"   have: <func>(<addr>)"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_kind_Uintptr(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		want := uintptr(42)
		have := uintptr(42)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal without addresses", func(t *testing.T) {
		// --- Given ---
		want := uintptr(42)
		have := uintptr(44)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: <addr>\n" +
			"   have: <addr>"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal with addresses", func(t *testing.T) {
		// --- Given ---
		want := uintptr(42)
		have := uintptr(44)

		trail := make([]string, 0)
		opts := []any{
			WithTrailLog(&trail),
			WithTrail("type.field"),
			WithDumper(dump.WithPtrAddr),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: <0x2a>\n" +
			"   have: <0x2c>"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_kind_UnsafePointer(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		w := 42
		want := unsafe.Pointer(&w)
		have := unsafe.Pointer(&w)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		w := 42
		h := 42
		want := unsafe.Pointer(&w)
		have := unsafe.Pointer(&h)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: <addr>\n" +
			"   have: <addr>"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_shared_reference_tabular(t *testing.T) {
	type TM struct{ A, B map[string]int }
	type TS struct{ A, B []int }
	type TP struct{ A, B *int }

	dict := map[string]int{"a": 1}
	list := []int{1}
	one, oneB, two := 1, 1, 2

	tt := []struct {
		testN  string
		shared any
		split  any
	}{
		{
			"map",
			TM{dict, dict},
			TM{map[string]int{"a": 1}, map[string]int{"a": 2}},
		},
		{"slice", TS{list, list}, TS{[]int{1}, []int{2}}},
		{"pointer", TP{&one, &one}, TP{&oneB, &two}},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Equal(tc.shared, tc.split)
			rev := Equal(tc.split, tc.shared)

			// --- Then ---
			affirm.NotNil(t, err)
			affirm.NotNil(t, rev)
		})
	}
}

func Test_Equal_shared_reference_trail_options(t *testing.T) {
	type T struct{ A, B []int }

	t.Run("error - skip keeps other trails", func(t *testing.T) {
		// --- Given ---
		x := []int{1, 2}
		y := []int{1, 3}

		// --- When ---
		err := Equal(T{x, x}, T{y, y}, WithSkipTrail("T.A[1]"))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.B[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.NotNil(t, err)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - trail checker keeps other trails", func(t *testing.T) {
		// --- Given ---
		x := []int{1, 2}
		y := []int{1, 3}
		chk := func(_, _ any, _ ...any) error { return nil }

		// --- When ---
		err := Equal(T{x, x}, T{y, y}, WithTrailChecker("T.A[1]", chk))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.B[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.NotNil(t, err)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - option on a later trail still runs", func(t *testing.T) {
		// --- Given ---
		x := []int{1, 2}
		y := []int{1, 2}
		chk := func(_, _ any, _ ...any) error { return errors.New("checked") }

		// --- When ---
		err := Equal(T{x, x}, T{y, y}, WithTrailChecker("T.B[1]", chk))

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "assertion error:\n  error: checked", err.Error())
	})

	t.Run("error - skip on a later trail", func(t *testing.T) {
		// --- Given ---
		x := []int{1, 2}
		y := []int{1, 3}

		// --- When ---
		err := Equal(T{x, x}, T{y, y}, WithSkipTrail("T.B[1]"))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: T.A[1]\n" +
			"   want: 2\n" +
			"   have: 3"
		affirm.NotNil(t, err)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("no option logs a shared pair once", func(t *testing.T) {
		// --- Given ---
		x := []int{1}
		y := []int{1}

		trail := make([]string, 0)

		// --- When ---
		err := Equal(T{x, x}, T{y, y}, WithTrailLog(&trail))

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A[0]"}, trail)
	})

	t.Run("unrelated option logs a shared pair once", func(t *testing.T) {
		// --- Given ---
		x := []int{1}
		y := []int{1}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithSkipTrail("other")}

		// --- When ---
		err := Equal(T{x, x}, T{y, y}, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T.A[0]"}, trail)
	})

	t.Run("walk under an option keeps an earlier mark", func(t *testing.T) {
		// --- Given ---
		type T3 struct{ A, B, C []int }
		x := []int{1}
		y := []int{1}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithSkipTrail("T3.B[0]")}

		// --- When ---
		err := Equal(T3{x, x, x}, T3{y, y, y}, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"T3.A[0]", "T3.B[0] <skipped>"}, trail)
	})

	t.Run("shared chain stays linear", func(t *testing.T) {
		// --- Given ---
		type N struct {
			L, R *N
			V    int
		}
		build := func() *N {
			nod := new(N)
			for range 12 {
				nod = &N{L: nod, R: nod}
			}
			return nod
		}

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithSkipTrail("N.L.L.V")}

		// --- When ---
		err := Equal(build(), build(), opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, true, len(trail) < 200)
	})

	t.Run("shared maps with tied keys stay linear", func(t *testing.T) {
		// --- Given ---
		type K struct{ V int }
		type M map[K]M
		build := func() M {
			nod := make(M)
			for range 10 {
				nod = M{K{1}: nod, K{2}: nod}
			}
			return nod
		}
		skip := strings.Repeat("map[check.K]", 10)

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithSkipTrail(skip)}

		// --- When ---
		err := Equal(build(), build(), opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, true, len(trail) < 50)
	})

	t.Run("skip inside a cycle stops", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		defer debug.SetMaxStack(debug.SetMaxStack(64 << 20))

		type N struct {
			V    int
			Next *N
		}
		want := &N{V: 1}
		want.Next = want
		have := &N{V: 2}
		have.Next = have

		// --- When ---
		err := Equal(want, have, WithSkipTrail("N.V"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("skip hides a value linked by a cycle", func(t *testing.T) {
		// --- Given ---
		type P struct {
			S int
			Q *struct{ P any }
		}
		type C struct {
			A *P
			B *struct{ P any }
		}
		build := func(val int) C {
			par := &P{S: val, Q: new(struct{ P any })}
			par.Q.P = par
			return C{A: par, B: par.Q}
		}

		// --- When ---
		err := Equal(build(1), build(2), WithSkipTrail("C.A.S"))

		// --- Then ---
		// The documented limit: B.P.S differs and is not reported.
		affirm.Nil(t, err)
	})

	t.Run("error - field order decides the cycle limit", func(t *testing.T) {
		// --- Given ---
		type P struct {
			S int
			Q *struct{ P any }
		}
		type C struct {
			B *struct{ P any }
			A *P
		}
		build := func(val int) C {
			par := &P{S: val, Q: new(struct{ P any })}
			par.Q.P = par
			return C{A: par, B: par.Q}
		}

		// --- When ---
		err := Equal(build(1), build(2), WithSkipTrail("C.A.S"))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: C.B.P.S\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.NotNil(t, err)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("cycles stop", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		defer debug.SetMaxStack(debug.SetMaxStack(64 << 20))

		type N struct {
			Next *N
			Kids []any
			Dict map[string]any
		}
		build := func() *N {
			nod := &N{Kids: make([]any, 1), Dict: make(map[string]any)}
			nod.Next = nod
			nod.Kids[0] = nod.Kids
			nod.Dict["self"] = nod.Dict
			return nod
		}

		// --- When ---
		err := Equal(build(), build(), WithSkipTrail("other"))

		// --- Then ---
		affirm.Nil(t, err)
	})
}

func Test_Equal_EqualCases_tabular(t *testing.T) {
	for _, tc := range testcases.EqualCases() {
		t.Run("Equal "+tc.Desc, func(t *testing.T) {
			// --- When ---
			have := Equal(tc.Val0, tc.Val1)

			// --- Then ---
			affirm.Equal(t, tc.AreEqual, have == nil)
		})
	}
}

func Test_NotEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := NotEqual(42, 44)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - fails", func(t *testing.T) {
		// --- When ---
		err := NotEqual(42, 42)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values not to be equal:\n" +
			"  want: 42\n" +
			"  have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - with bytes", func(t *testing.T) {
		// --- When ---
		err := NotEqual(byte(42), byte(42))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values not to be equal:\n" +
			"  want: 0x2a ('*')\n" +
			"  have: 0x2a ('*')"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := NotEqual(42, 42, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values not to be equal:\n" +
			"  trail: type.field\n" +
			"   want: 42\n" +
			"   have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_mapEntry_valOrder(t *testing.T) {
	t.Run("value is dumped once", func(t *testing.T) {
		// --- Given ---
		type T struct{ N int }

		var count int
		dpr := func(dump.Dump, int, reflect.Value) string {
			count++
			return "T"
		}
		dmp := dump.New(dump.WithDumper(T{}, dpr))
		ent := &mapEntry{val: reflect.ValueOf(T{})}

		// --- When ---
		have := ent.valOrder(dmp)

		// --- Then ---
		affirm.Equal(t, "T check.T", have)
		affirm.Equal(t, "T check.T", ent.valOrder(dmp))
		affirm.Equal(t, 1, count)
	})

	t.Run("dynamic type orders equal dumps", func(t *testing.T) {
		// --- Given ---
		vals := reflect.ValueOf(map[string]any{"a": 1, "b": int64(1)})
		ent0 := &mapEntry{val: vals.MapIndex(reflect.ValueOf("a"))}
		ent1 := &mapEntry{val: vals.MapIndex(reflect.ValueOf("b"))}

		// --- When ---
		hOrd0 := ent0.valOrder(dump.New())
		hOrd1 := ent1.valOrder(dump.New())

		// --- Then ---
		affirm.Equal(t, "1 int", hOrd0)
		affirm.Equal(t, "1 int64", hOrd1)
	})

	t.Run("nil interface value", func(t *testing.T) {
		// --- Given ---
		vals := reflect.ValueOf(map[string]any{"a": nil})
		ent := &mapEntry{val: vals.MapIndex(reflect.ValueOf("a"))}

		// --- When ---
		have := ent.valOrder(dump.New())

		// --- Then ---
		affirm.Equal(t, "nil interface {}", have)
	})
}

func Test_equalError(t *testing.T) {
	t.Run("without trail", func(t *testing.T) {
		// --- Given ---
		ops := DefaultOptions()

		// --- When ---
		err := equalError(42, 44, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: 42\n" +
			"  have: 44"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("with trail", func(t *testing.T) {
		// --- Given ---
		ops := DefaultOptions(WithTrail("type.field"))

		// --- When ---
		err := equalError(42, 44, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: 42\n" +
			"   have: 44"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("printable byte", func(t *testing.T) {
		// --- Given ---
		w := byte('A')
		h := byte('B')
		ops := DefaultOptions()

		// --- When ---
		err := equalError(w, h, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: 0x41 ('A')\n" +
			"  have: 0x42 ('B')"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("keeps the byte dumper already set", func(t *testing.T) {
		// --- Given ---
		w := byte('A')
		h := byte('B')
		fn := func(_ dump.Dump, _ int, _ reflect.Value) string { return "abc" }
		ops := DefaultOptions(WithDumper(dump.WithDumper(byte(0), fn)))

		// --- When ---
		err := equalError(w, h, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want: abc\n" +
			"  have: abc"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("different types", func(t *testing.T) {
		// --- Given ---
		w := byte('A')
		h := 42
		ops := DefaultOptions()

		// --- When ---
		err := equalError(w, h, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: uint8\n" +
			"  have type: int\n" +
			"       want: 0x41 ('A')\n" +
			"       have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("different types and w is nil", func(t *testing.T) {
		// --- Given ---
		h := 42
		ops := DefaultOptions()

		// --- When ---
		err := equalError(nil, h, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: <nil>\n" +
			"  have type: int\n" +
			"       want: nil\n" +
			"       have: 42"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("different types and h is nil", func(t *testing.T) {
		// --- Given ---
		w := 42
		ops := DefaultOptions()

		// --- When ---
		err := equalError(w, nil, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: int\n" +
			"  have type: <nil>\n" +
			"       want: 42\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("shows diff", func(t *testing.T) {
		// --- Given ---
		w := []int{1, 2, 3}
		h := []int{1, 2, 4}
		ops := DefaultOptions()

		// --- When ---
		err := equalError(w, h, WithOptions(ops))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want:\n" +
			"        []int{\n" +
			"          1,\n" +
			"          2,\n" +
			"          3,\n" +
			"        }\n" +
			"  have:\n" +
			"        []int{\n" +
			"          1,\n" +
			"          2,\n" +
			"          4,\n" +
			"        }\n" +
			"  diff:\n" +
			"        @@ -2,4 +2,4 @@\n" +
			"           1,\n" +
			"           2,\n" +
			"        -  4,\n" +
			"        +  3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_nilnessError(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		// --- Given ---
		wVal := reflect.ValueOf([]int(nil))
		hVal := reflect.ValueOf([]int(nil))

		// --- When ---
		err := nilnessError(wVal, hVal, DefaultOptions())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("both not nil", func(t *testing.T) {
		// --- Given ---
		wVal := reflect.ValueOf(make(map[string]int))
		hVal := reflect.ValueOf(make(map[string]int))

		// --- When ---
		err := nilnessError(wVal, hVal, DefaultOptions())

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("one nil with EquateEmpty", func(t *testing.T) {
		// --- Given ---
		wVal := reflect.ValueOf([]int(nil))
		hVal := reflect.ValueOf(make([]int, 0))
		ops := DefaultOptions(WithEquateEmpty())

		// --- When ---
		err := nilnessError(wVal, hVal, ops)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - one nil", func(t *testing.T) {
		// --- Given ---
		wVal := reflect.ValueOf([]int(nil))
		hVal := reflect.ValueOf(make([]int, 0))

		trail := make([]string, 0)
		opts := []any{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := nilnessError(wVal, hVal, DefaultOptions(opts...))

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: type.field\n" +
			"   want: nil\n" +
			"   have: []int{}"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_dumpByte(t *testing.T) {
	t.Run("printable", func(t *testing.T) {
		// --- Given ---
		dmp := dump.New()
		val := reflect.ValueOf(byte(42))

		// --- When ---
		have := dumpByte(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "0x2a ('*')", have)
	})

	t.Run("not printable", func(t *testing.T) {
		// --- Given ---
		dmp := dump.New()
		val := reflect.ValueOf(byte(1))

		// --- When ---
		have := dumpByte(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "0x01", have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := dump.New(dump.WithIndent(2))
		val := reflect.ValueOf(byte(1))

		// --- When ---
		have := dumpByte(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "      0x01", have)
	})

	t.Run("unexported field", func(t *testing.T) {
		// --- Given ---
		dmp := dump.New()
		val := reflect.ValueOf(struct{ b byte }{b: 'a'}).Field(0)

		// --- When ---
		have := dumpByte(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "0x61 ('a')", have)
	})
}

// -----------------------------------------------------------------------------
// Benchmarks
// -----------------------------------------------------------------------------

type benchUser struct {
	ID        int
	Name      string
	Email     string
	Active    bool
	CreatedAt time.Time
	Tags      []string
	Metadata  map[string]any
}

type benchOrder struct {
	ID       int
	User     benchUser
	Items    []benchOrderItem
	Total    float64
	Status   string
	Shipped  bool
	Metadata map[string]any
}

type benchOrderItem struct {
	SKU      string
	Quantity int
	Price    float64
}

func makeBenchUser() benchUser {
	return benchUser{
		ID:        42,
		Name:      "Alice Example",
		Email:     "alice@example.com",
		Active:    true,
		CreatedAt: time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC),
		Tags:      []string{"premium", "beta", "europe"},
		Metadata: map[string]any{
			"plan":   "pro",
			"region": "eu-west-1",
			"score":  87.5,
		},
	}
}

func makeBenchOrder() benchOrder {
	return benchOrder{
		ID:   1001,
		User: makeBenchUser(),
		Items: []benchOrderItem{
			{SKU: "WIDGET-42", Quantity: 3, Price: 19.99},
			{SKU: "GADGET-X", Quantity: 1, Price: 49.5},
		},
		Total:   109.47,
		Status:  "shipped",
		Shipped: true,
		Metadata: map[string]any{
			"channel": "web",
			"promo":   "SPRING25",
		},
	}
}

func Benchmark_Equal_SimpleStruct(b *testing.B) {
	want := makeBenchUser()
	have := makeBenchUser()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have)
	}
}

func Benchmark_Equal_NestedStruct(b *testing.B) {
	want := makeBenchOrder()
	have := makeBenchOrder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have)
	}
}

func Benchmark_Equal_LargeSlice(b *testing.B) {
	orders := make([]benchOrder, 0, 100)
	for i := 0; i < 100; i++ {
		o := makeBenchOrder()
		o.ID = i
		orders = append(orders, o)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(orders, orders)
	}
}

func Benchmark_Equal_WithTrail(b *testing.B) {
	want := makeBenchOrder()
	have := makeBenchOrder()
	have.User.Name = "Different Name"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have, WithTrail("order"))
	}
}

func Benchmark_Equal_WithSkipUnexported(b *testing.B) {
	type s struct {
		Exported   int
		unexported string
	}
	want := s{Exported: 1, unexported: "a"}
	have := s{Exported: 1, unexported: "b"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have, WithSkipUnexported())
	}
}

func Benchmark_Equal_DeepWithCustomChecker(b *testing.B) {
	// Custom checker that does a bit of extra work then delegates.
	custom := func(want, have any, opts ...any) error {
		// Simulate some non-trivial work without risking deep recursion
		// or expensive formatting of large structures.
		_ = len(want.(benchOrder).Items) + len(have.(benchOrder).Items)

		// Delegate to normal comparison while explicitly disabling this
		// custom checker for the recursive call (prevents infinite recursion).
		callOpts := append([]any{}, opts...)
		callOpts = append(callOpts, WithTypeChecker(benchOrder{}, nil))
		return Equal(want, have, callOpts...)
	}

	want := makeBenchOrder()
	have := makeBenchOrder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have, WithTypeChecker(benchOrder{}, custom))
	}
}

func Benchmark_Equal_WithTrail_Deep(b *testing.B) {
	type inner struct{ Val int }
	type outer struct{ Items []inner }

	want := outer{Items: []inner{{Val: 1}, {Val: 2}}}
	have := outer{Items: []inner{{Val: 1}, {Val: 99}}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have, WithTrail("Items[1].Val"))
	}
}

func Benchmark_Equal_WithOptions_Reused(b *testing.B) {
	ops := DefaultOptions(
		WithSkipUnexported(),
		WithTimeFormat(time.RFC3339),
	)

	type s struct {
		Exported   int
		unexported string
		t          time.Time
	}
	want := s{Exported: 42, unexported: "x", t: time.Now()}
	have := s{Exported: 42, unexported: "x", t: time.Now()}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Equal(want, have, WithOptions(ops))
	}
}
