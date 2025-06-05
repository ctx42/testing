package check

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/cases"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/must"
)

func Test_Equal_invalid_arguments(t *testing.T) {
	t.Run("equal both are untyped nil", func(t *testing.T) {
		// --- When ---
		err := Equal(nil, nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal untyped nil and nil interface", func(t *testing.T) {
		// --- Given ---
		var itf types.TItf

		// --- When ---
		err := Equal(nil, itf)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal nil interface and untyped nil ", func(t *testing.T) {
		// --- Given ---
		var itf types.TItf

		// --- When ---
		err := Equal(itf, nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("logs trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrail("type.field"), WithTrailLog(&trail)}

		// --- When ---
		err := Equal(nil, nil, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_one_argument_invalid(t *testing.T) {
	t.Run("want is invalid", func(t *testing.T) {
		// --- When ---
		err := Equal(nil, 123)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: <nil>\n" +
			"  have type: int\n" +
			"       want: nil\n" +
			"       have: 123"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("have is invalid", func(t *testing.T) {
		// --- When ---
		err := Equal(123, nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: int\n" +
			"  have type: <nil>\n" +
			"       want: 123\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("logs trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrail("type.field"), WithTrailLog(&trail)}

		// --- When ---
		err := Equal(123, nil, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: type.field\n" +
			"  want type: int\n" +
			"  have type: <nil>\n" +
			"       want: 123\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_not_matching_types(t *testing.T) {
	t.Run("not matching", func(t *testing.T) {
		// --- When ---
		err := Equal(123, "abc")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  want type: int\n" +
			"  have type: string\n" +
			"       want: 123\n" +
			"       have: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("logs trail", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrail("type.field"), WithTrailLog(&trail)}

		// --- When ---
		err := Equal(123, "abc", opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: type.field\n" +
			"  want type: int\n" +
			"  have type: string\n" +
			"       want: 123\n" +
			"       have: \"abc\""
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_custom_trail_checkers(t *testing.T) {
	t.Run("custom checker is not used", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{
			WithTrail("type.other"),
			WithTrailLog(&trail),
			WithTrailChecker("type.field", Exact),
		}

		// Both are defining the same time in different timezone.
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, types.WAW)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.other"}, trail)
	})

	t.Run("custom checker used", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{
			WithTrail("type.field"),
			WithTrailLog(&trail),
			WithTrailChecker("type.field", Exact),
		}

		// Both are defining the same time in different timezone.
		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, types.WAW)

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

func Test_Equal_custom_type_checkers(t *testing.T) {
	t.Run("use the custom type checker", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{
			WithTrail("type.field"),
			WithTrailLog(&trail),
			WithTrailChecker("type.field", Exact),
		}

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 4, 4, 5, 0, types.WAW)

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

	t.Run("use the custom checker with nils", func(t *testing.T) {
		// --- Given ---
		var want, have = 1, 2

		trail := make([]string, 0)
		opts := []Option{
			WithTrail("type.field"),
			WithTrailLog(&trail),
			WithTypeChecker(want, func(_, _ any, _ ...Option) error {
				return errors.New("custom checker")
			}),
		}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Equal(t, "custom checker", err.Error())
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})
}

func Test_Equal_kind_Ptr(t *testing.T) {
	t.Run("equal structs by pointer", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := &struct{ Int int }{Int: 123}
		have := &struct{ Int int }{Int: 123}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"Int"}, trail)
	})

	t.Run("equal time.Location struct pointer", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

		want := must.Value(time.LoadLocation("Europe/Warsaw"))
		have := must.Value(time.LoadLocation("Europe/Warsaw"))

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("equal both nil values", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

		var want *int
		var have *int

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("not equal want is not nil have is nil", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

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

	t.Run("not equal want is nil have is not nil", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

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
	t.Run("equal structs by value", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := struct{ Int int }{Int: 123}
		have := struct{ Int int }{Int: 123}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"Int"}, trail)
	})

	t.Run("equal time struct", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("equal time struct field", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := types.TTim{Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)}
		have := types.TTim{Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"TTim.Tim"}, trail)
	})

	t.Run("not equal time struct", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

		want := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		have := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)

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

	t.Run("equal time.Location struct value", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

		want := must.Value(time.LoadLocation("Europe/Warsaw"))
		have := must.Value(time.LoadLocation("Europe/Warsaw"))

		// --- When ---
		err := Equal(*want, *have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type"}, trail)
	})

	t.Run("equal time.Location struct field", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := types.TLoc{Loc: must.Value(time.LoadLocation("Europe/Warsaw"))}
		have := types.TLoc{Loc: must.Value(time.LoadLocation("Europe/Warsaw"))}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"TLoc.Loc"}, trail)
	})

	t.Run("not equal time.Location struct value", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type")}

		want := must.Value(time.LoadLocation("Europe/Warsaw"))
		have := must.Value(time.LoadLocation("Europe/Paris"))

		// --- When ---
		err := Equal(*want, *have, opts...)

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

	t.Run("not equal time.Location nil have", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := types.TLoc{Loc: types.WAW}
		have := types.TLoc{Loc: nil}

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
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := types.TC{TD: types.TD("abc"), Int: 123}
		have := types.TC{TD: types.TD("abc"), Int: 123}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"TC.TD", "TC.Int"}, trail)
	})

	t.Run("private fields must be skipped to avoid error", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithSkipTrail("TIntPrv.v")}

		want := types.NewTIntPrv(42, 1)
		have := types.NewTIntPrv(42, 2)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		wTrail := []string{"TIntPrv.Int", "TIntPrv.v <skipped>"}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("error - private fields not equal", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := types.NewTIntPrv(42, 1)
		have := types.NewTIntPrv(42, 2)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TIntPrv.v\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
		wTrail := []string{"TIntPrv.Int", "TIntPrv.v"}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("not equal with private fields", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithSkipTrail("TIntPrv.v")}

		want := types.NewTIntPrv(42, 1)
		have := types.NewTIntPrv(44, 2)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TIntPrv.Int\n" +
			"   want: 42\n" +
			"   have: 44"
		affirm.Equal(t, wMsg, err.Error())
		wTrail := []string{"TIntPrv.Int", "TIntPrv.v <skipped>"}
		affirm.DeepEqual(t, wTrail, trail)
	})

	t.Run("not equal structs with multiple errors", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := types.TIntStr{Int: 42, Str: "abc"}
		have := types.TIntStr{Int: 44, Str: "xyz"}

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

	t.Run("not equal when want is the nil struct pointer", func(t *testing.T) {
		// --- Given ---
		var want *types.TA
		have := &types.TA{Str: "abc"}

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
			"        }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("not equal when have is a nil struct pointer", func(t *testing.T) {
		// --- Given ---
		want := &types.TA{Str: "abc"}
		var have *types.TA

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
			"        }\n" +
			"  have: nil"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("not equal deeply nested", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithSkipUnexported}

		want := types.TNested{STAp: []*types.TA{{TAp: &types.TA{Int: 42}}}}
		have := types.TNested{STAp: []*types.TA{{TAp: &types.TA{Int: 44}}}}

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
		trail := make([]string, 0)
		opts := []Option{
			WithTrailLog(&trail),
			WithSkipTrail("TNested.STAp[0].TAp.Int"),
			WithSkipUnexported,
		}

		want := types.TNested{STAp: []*types.TA{{TAp: &types.TA{Int: 42}}}}
		have := types.TNested{STAp: []*types.TA{{TAp: &types.TA{Int: 44}}}}

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
}

func Test_Equal_kind_Slice_and_Array(t *testing.T) {
	t.Run("equal slice", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := []int{1, 2}
		have := []int{1, 2}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"<slice>[0]", "<slice>[1]"}, trail)
	})

	t.Run("equal same slice instance", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := []int{1, 2}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal slice value", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := []int{1, 2}
		have := []int{1, 7}

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

	t.Run("not equal array value", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := [...]int{1, 2}
		have := [...]int{1, 7}

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

	t.Run("not equal slice lengths", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := []int{1, 2}
		have := []int{1}

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

	t.Run("not equal slices with multiple errors", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := []int{1, 2}
		have := []int{2, 3}

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
	t.Run("equal map", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := map[int]int{1: 42}
		have := map[int]int{1: 42}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"map[1]"}, trail)
	})

	t.Run("equal same map", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := map[int]int{1: 42}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal map", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := map[int]int{1: 42}
		have := map[int]int{1: 44}

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

	t.Run("not equal have map missing keys", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := map[int]int{1: 42, 2: 43}
		have := map[int]int{1: 42, 3: 44}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected values to be equal:\n" +
			"      trail: map[2]\n" +
			"  want type: map[int]int\n" +
			"  have type: <nil>\n" +
			"       want:\n" +
			"             map[int]int{\n" +
			"               1: 42,\n" +
			"               3: 44,\n" +
			"             }\n" +
			"       have: nil"
		affirm.Equal(t, wMsg, err.Error())
		affirm.DeepEqual(t, []string{"map[1]"}, trail)
	})

	t.Run("not equal map length", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := map[int]int{1: 42, 2: 44}
		have := map[int]int{1: 42, 2: 43, 3: 44}

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

	t.Run("not equal maps with multiple errors", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := map[int]int{1: 42, 2: 44}
		have := map[int]int{1: 44, 2: 42}

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
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}

		want := []any{42}
		have := []any{42}

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"<slice>[0]"}, trail)
	})

	t.Run("equal both nil", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail)}
		want := []any{nil}
		have := []any{nil}

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
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(true, true, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("equal false", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		// --- When ---
		err := Equal(false, false, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

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

func Test_Equal_kind_numbers_success_tabular(t *testing.T) {
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
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			trail := make([]string, 0)
			opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

			// --- When ---
			err := Equal(tc.want, tc.have, opts...)

			// --- Then ---
			affirm.Nil(t, err)
			affirm.DeepEqual(t, []string{"type.field"}, trail)
		})
	}
}

func Test_Equal_kind_numbers_error_tabular(t *testing.T) {
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
			opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

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
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := make(chan bool)

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := make(chan bool)
		have := make(chan bool)

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

	t.Run("not equal unexported field", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := struct{ want chan bool }{make(chan bool)}
		have := struct{ want chan bool }{make(chan bool)}

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
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := func() {}

		// --- When ---
		err := Equal(want, want, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := func() {}
		have := func() {}

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
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := uintptr(42)
		have := uintptr(42)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal without addresses", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		want := uintptr(42)
		have := uintptr(44)

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

	t.Run("not equal with addresses", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{
			WithTrailLog(&trail),
			WithTrail("type.field"),
			WithDumper(dump.WithPtrAddr),
		}

		want := uintptr(42)
		have := uintptr(44)

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
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		w := 42
		want := unsafe.Pointer(&w)
		have := unsafe.Pointer(&w)

		// --- When ---
		err := Equal(want, have, opts...)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.DeepEqual(t, []string{"type.field"}, trail)
	})

	t.Run("not equal", func(t *testing.T) {
		// --- Given ---
		trail := make([]string, 0)
		opts := []Option{WithTrailLog(&trail), WithTrail("type.field")}

		w := 42
		h := 42
		want := unsafe.Pointer(&w)
		have := unsafe.Pointer(&h)

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

func Test_Equal_EqualCases_tabular(t *testing.T) {
	for _, tc := range cases.EqualCases() {
		t.Run("Equal "+tc.Desc, func(t *testing.T) {
			// --- When ---
			have := Equal(tc.Val0, tc.Val1)

			// --- Then ---
			if tc.AreEqual && have != nil {
				format := "expected nil error:\n  have: %#v"
				t.Errorf(format, have)
			}
			if !tc.AreEqual && have == nil {
				format := "expected not-nil error:\n  have: %#v"
				t.Errorf(format, have)
			}
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

	t.Run("error", func(t *testing.T) {
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

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NotEqual(42, 42, opt)

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

func Test_equalError(t *testing.T) {
	t.Run("without trail", func(t *testing.T) {
		// --- Given ---
		ops := DefaultOptions()

		// --- When ---
		err := equalError(42, 44, WithOptions(ops))

		// --- Then ---
		wMsg := "expected values to be equal:\n" +
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

	t.Run("it does not override the already set byte dumper", func(t *testing.T) {
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
}
