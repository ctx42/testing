// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/notice"
)

func Test_RegisterTypeChecker(t *testing.T) {
	t.Setenv("___", "___")
	affirm.Nil(t, typeCheckers)
	origLog := globLog
	buf := &bytes.Buffer{}
	globLog = log.New(buf, "", 0)
	chk := func(_, _ any, _ ...any) error { return errors.New("123456") }
	t.Cleanup(func() { globLog = origLog; typeCheckers = nil })

	t.Run("is registered", func(t *testing.T) {
		// --- Given ---
		type custom struct{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })

		// --- When ---
		RegisterTypeChecker(custom{}, chk)

		// --- Then ---
		affirm.Equal(t, 1, len(typeCheckers))
		wChk := core.Same(chk, typeCheckers[reflect.TypeOf(custom{})])
		affirm.Equal(t, true, wChk)
		wMsg := "Registering type checker for: check.custom\n"
		affirm.Equal(t, wMsg, buf.String())
	})

	t.Run("panics if already registered", func(t *testing.T) {
		// --- Given ---
		type custom struct{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })
		typeCheckers = map[reflect.Type]Checker{reflect.TypeOf(custom{}): chk}

		// --- When ---
		fn := func() { RegisterTypeChecker(custom{}, chk) }
		msg := affirm.Panic(t, fn)

		// --- Then ---
		wMsg := "cannot overwrite an existing type checker"
		affirm.Equal(t, true, strings.Contains(*msg, wMsg))
		affirm.Equal(t, "", buf.String())
	})

	t.Run("panics if the checker is nil", func(t *testing.T) {
		// --- Given ---
		type custom struct{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })

		// --- When ---
		fn := func() { RegisterTypeChecker(custom{}, nil) }
		msg := affirm.Panic(t, fn)

		// --- Then ---
		affirm.Equal(t, "cannot register a nil type checker", *msg)
		affirm.Equal(t, "", buf.String())
	})
}

func Test_WithTrail(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithTrail("type.field")(ops)

	// --- Then ---
	affirm.Equal(t, "", ops.Trail)
	affirm.Equal(t, "type.field", have.Trail)
}

func Test_WithTrailLog(t *testing.T) {
	// --- Given ---
	buf := make([]string, 0)
	ops := Options{}

	// --- When ---
	have := WithTrailLog(&buf)(ops)

	// --- Then ---
	affirm.Equal(t, true, core.Same(&buf, have.TrailLog))
}

func Test_WithTimeFormat(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithTimeFormat(time.RFC3339)(ops)

	// --- Then ---
	affirm.Equal(t, time.RFC3339, have.TimeFormat)
}

func Test_WithZone(t *testing.T) {
	// --- Given ---
	waw := must.Value(time.LoadLocation("Europe/Warsaw"))
	ops := Options{}

	// --- When ---
	have := WithZone(waw)(ops)

	// --- Then ---
	affirm.Equal(t, true, core.Same(waw, have.Zone))
}

func Test_WithRecent(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithRecent(time.Second)(ops)

	// --- Then ---
	affirm.Equal(t, time.Second, have.Recent)
}

func Test_WithDumper(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithDumper(dump.WithMaxDepth(100))(ops)

	// --- Then ---
	affirm.Equal(t, 0, ops.Dumper.MaxDepth)
	affirm.Equal(t, 100, have.Dumper.MaxDepth)
}

func Test_WithTypeChecker(t *testing.T) {
	t.Setenv("___", "___")
	affirm.Nil(t, typeCheckers)
	origLog := globLog
	buf := &bytes.Buffer{}
	globLog = log.New(buf, "", 0)
	cChk := func(_, _ any, _ ...any) error { return errors.New("123456") }
	t.Cleanup(func() { globLog = origLog; typeCheckers = nil })

	t.Run("setting", func(t *testing.T) {
		// --- Given ---
		ops := Options{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })

		// --- When ---
		have := WithTypeChecker(123, cChk)(ops)

		// --- Then ---
		wChk := core.Same(cChk, have.TypeCheckers[reflect.TypeOf(123)])
		affirm.Equal(t, true, wChk)
		affirm.Equal(t, "", buf.String())
	})

	t.Run("overwriting global checker", func(t *testing.T) {
		// --- Given ---
		type custom struct{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })

		RegisterTypeChecker(custom{}, cChk) // The first call.
		buf.Reset()                         // Test later the log is empty.

		ops := Options{}

		// --- When ---
		have := WithTypeChecker(custom{}, cChk)(ops)

		// --- Then ---
		wChk := core.Same(cChk, have.TypeCheckers[reflect.TypeOf(custom{})])
		affirm.Equal(t, true, wChk)
		wMsg := "Overwriting the global type checker for: check.custom\n"
		affirm.Equal(t, wMsg, buf.String())
	})
}

func Test_WithTrailChecker(t *testing.T) {
	// --- Given ---
	ops := Options{}
	chk := func(want, have any, opts ...any) error { return nil }

	// --- When ---
	have := WithTrailChecker("type.field", chk)(ops)

	// --- Then ---
	haveChk, _ := have.TrailCheckers["type.field"]
	affirm.Equal(t, true, core.Same(chk, haveChk))
}

func Test_WithSkipTrail(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithSkipTrail("type.field1", "type.field2")(ops)

	// --- Then ---
	affirm.Equal(t, true, ops.SkipTrails == nil)
	affirm.DeepEqual(t, []string{"type.field1", "type.field2"}, have.SkipTrails)
}

func Test_WithSkipUnexported(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithSkipUnexported(ops)

	// --- Then ---
	affirm.Equal(t, false, ops.SkipUnexported)
	affirm.Equal(t, true, have.SkipUnexported)
}

func Test_WithIncreasingSoft(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithIncreasingSoft(ops)

	// --- Then ---
	affirm.Equal(t, true, have.IncreaseSoft)
}

func Test_WithDecreasingSoft(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithDecreasingSoft(ops)

	// --- Then ---
	affirm.Equal(t, true, have.DecreaseSoft)
}

func Test_WithCmpBaseTypes(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithCmpBaseTypes(ops)

	// --- Then ---
	affirm.Equal(t, true, have.CmpSimpleType)
}

func Test_WithWaitThrottle(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithWaitThrottle(time.Second)(ops)

	// --- Then ---
	affirm.Equal(t, time.Second, have.WaitThrottle)
}

func Test_WithComment(t *testing.T) {
	// --- Given ---
	ops := Options{}

	// --- When ---
	have := WithComment("A%d", 42)(ops)

	// --- Then ---
	affirm.Equal(t, "A42", have.Comment)
}

func Test_WithOptions(t *testing.T) {
	// --- Given ---
	waw := must.Value(time.LoadLocation("Europe/Warsaw"))
	trailLog := make([]string, 0)
	ops := Options{
		Dumper: dump.Dump{
			Flat:           true,
			FlatStrings:    100,
			Compact:        true,
			TimeFormat:     time.Kitchen,
			DurationFormat: "DurAsString",
			PtrAddr:        true,
			PrintType:      true,
			PrintPrivate:   true,
			UseAny:         true,
			Dumpers: map[reflect.Type]dump.Dumper{
				reflect.TypeOf(123): dump.Dumper(nil),
			},
			MaxDepth: 6,
			Indent:   2,
			TabWidth: 4,
		},
		TimeFormat:     time.RFC3339,
		Zone:           waw,
		Recent:         123,
		Trail:          "trail",
		TrailLog:       &trailLog,
		TypeCheckers:   make(map[reflect.Type]Checker),
		TrailCheckers:  make(map[string]Checker),
		SkipTrails:     make([]string, 0),
		SkipUnexported: true,
		CmpSimpleType:  true,
		IncreaseSoft:   true,
		DecreaseSoft:   true,
		WaitThrottle:   10 * time.Millisecond,
		Comment:        "comment",
		now:            time.Now,
	}

	// --- When ---
	have := WithOptions(ops)(Options{})

	// --- Then ---
	affirm.Equal(t, true, core.Same(ops.Dumper.Dumpers, have.Dumper.Dumpers))
	affirm.Equal(t, true, core.Same(ops.Zone, have.Zone))
	affirm.Equal(t, true, core.Same(ops.TrailLog, have.TrailLog))
	affirm.Equal(t, true, core.Same(ops.TypeCheckers, have.TypeCheckers))
	affirm.Equal(t, true, core.Same(ops.TrailCheckers, have.TrailCheckers))
	affirm.Equal(t, true, core.Same(ops.SkipTrails, have.SkipTrails))
	affirm.Equal(t, true, core.Same(ops.now, have.now))

	ops.now = nil
	have.now = nil
	affirm.Equal(t, true, reflect.DeepEqual(ops, have))

	// When those fail, add fields above.
	affirm.Equal(t, 14, reflect.ValueOf(have.Dumper).NumField())
	affirm.Equal(t, 16, reflect.ValueOf(have).NumField())
}

func Test_DefaultOptions(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have := DefaultOptions()

		// --- Then ---
		affirm.Equal(t, false, have.Dumper.PtrAddr)
		affirm.Equal(t, DefaultDumpTimeFormat, have.Dumper.TimeFormat)

		affirm.Equal(t, DefaultParseTimeFormat, have.TimeFormat)
		affirm.Nil(t, have.Zone)
		affirm.Equal(t, DefaultRecentDuration, have.Recent)
		affirm.Equal(t, "", have.Trail)
		affirm.Equal(t, true, have.TrailLog == nil)
		affirm.Equal(t, false, have.TypeCheckers == nil)
		affirm.Equal(t, true, have.TrailCheckers == nil)
		affirm.Equal(t, true, core.Same(Time, have.TypeCheckers[typTime]))
		affirm.Equal(t, true, core.Same(Zone, have.TypeCheckers[typZone]))
		affirm.Equal(t, true, core.Same(Zone, have.TypeCheckers[typZonePtr]))
		affirm.Equal(t, true, have.SkipTrails == nil)
		affirm.Equal(t, false, have.SkipUnexported)
		affirm.Equal(t, false, have.CmpSimpleType)
		affirm.Equal(t, false, have.IncreaseSoft)
		affirm.Equal(t, false, have.DecreaseSoft)
		affirm.Equal(t, 10*time.Millisecond, have.WaitThrottle)
		affirm.Equal(t, "", have.Comment)
		affirm.Equal(t, true, core.Same(time.Now, have.now))
		affirm.Equal(t, 16, reflect.ValueOf(have).NumField())
	})

	t.Run("with options", func(t *testing.T) {
		// --- When ---
		have := DefaultOptions(WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have.Dumper.PtrAddr)
		affirm.Equal(t, DefaultDumpTimeFormat, have.Dumper.TimeFormat)

		affirm.Equal(t, DefaultParseTimeFormat, have.TimeFormat)
		affirm.Nil(t, have.Zone)
		affirm.Equal(t, DefaultRecentDuration, have.Recent)
		affirm.Equal(t, "type.field", have.Trail)
		affirm.Equal(t, true, have.TrailLog == nil)
		affirm.Equal(t, true, have.TrailCheckers == nil)
		affirm.Equal(t, true, core.Same(Time, have.TypeCheckers[typTime]))
		affirm.Equal(t, true, core.Same(Zone, have.TypeCheckers[typZone]))
		affirm.Equal(t, true, core.Same(Zone, have.TypeCheckers[typZonePtr]))
		affirm.Equal(t, true, have.TrailCheckers == nil)
		affirm.Equal(t, true, have.SkipTrails == nil)
		affirm.Equal(t, false, have.SkipUnexported)
		affirm.Equal(t, false, have.CmpSimpleType)
		affirm.Equal(t, false, have.IncreaseSoft)
		affirm.Equal(t, false, have.DecreaseSoft)
		affirm.Equal(t, 10*time.Millisecond, have.WaitThrottle)
		affirm.Equal(t, "", have.Comment)
		affirm.Equal(t, true, core.Same(time.Now, have.now))
		affirm.Equal(t, 16, reflect.ValueOf(have).NumField())
	})

	t.Run("TypeCheckers field is a clone of a global map", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		affirm.Nil(t, typeCheckers)
		chk := func(_, _ any, _ ...any) error { return errors.New("123456") }
		t.Cleanup(func() { typeCheckers = nil })
		typeCheckers = map[reflect.Type]Checker{reflect.TypeOf(123): chk}

		// --- When ---
		have := DefaultOptions()

		// --- Then ---
		affirm.Equal(t, false, core.Same(typeCheckers, have.TypeCheckers))
	})

	t.Run("the time check is not overwritten when set", func(t *testing.T) {
		// --- Given ---
		chk := func(_, _ any, _ ...any) error { return nil }
		opt := WithTypeChecker(time.Time{}, chk)

		// --- When ---
		have := DefaultOptions(opt)

		// --- Then ---
		affirm.Equal(t, true, core.Same(chk, have.TypeCheckers[typTime]))
	})

	t.Run("the timezone check is not overwritten when set", func(t *testing.T) {
		// --- Given ---
		chk := func(_, _ any, _ ...any) error { return nil }
		opt := WithTypeChecker(time.Location{}, chk)

		// --- When ---
		have := DefaultOptions(opt)

		// --- Then ---
		affirm.Equal(t, true, core.Same(chk, have.TypeCheckers[typZone]))
	})

	t.Run("timezone ptr check is not overwritten when set", func(t *testing.T) {
		// --- Given ---
		chk := func(_, _ any, _ ...any) error { return nil }
		opt := WithTypeChecker(&time.Location{}, chk)

		// --- When ---
		have := DefaultOptions(opt)

		// --- Then ---
		affirm.Equal(t, true, core.Same(chk, have.TypeCheckers[typZonePtr]))
	})

	t.Run("with formated comment", func(t *testing.T) {
		// --- When ---
		ops := DefaultOptions("A%d", 42)

		// --- Then ---
		affirm.Equal(t, "A42", ops.Comment)
	})

	t.Run("with option", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		have := DefaultOptions(opt)

		// --- Then ---
		affirm.Equal(t, "type.field", have.Trail)
	})

	t.Run("with Option type", func(t *testing.T) {
		// --- Given ---
		opt := WithCmpBaseTypes

		// --- When ---
		have := DefaultOptions(opt)

		// --- Then ---
		affirm.Equal(t, true, have.CmpSimpleType)
	})

	t.Run("with option and formated comment", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		have := DefaultOptions("A%s", opt, "BC")

		// --- Then ---
		affirm.Equal(t, "ABC", have.Comment)
		affirm.Equal(t, "type.field", have.Trail)
	})

	t.Run("error - cannot use a non-string comment format", func(t *testing.T) {
		// --- When ---
		msg := affirm.Panic(t, func() { DefaultOptions(42, "A%d") })

		// --- Then ---
		affirm.Equal(t, "cannot use a non-string comment format", *msg)
	})
}

func Test_Options_LogTrail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		list := make([]string, 0)
		ops := Options{Trail: "abc", TrailLog: &list}

		// --- When ---
		have := ops.LogTrail()

		// --- Then ---
		affirm.DeepEqual(t, []string{"abc"}, list)
		affirm.DeepEqual(t, []string{"abc"}, *ops.TrailLog)
		affirm.DeepEqual(t, have, ops)
	})

	t.Run("does not log empty trails", func(t *testing.T) {
		// --- Given ---
		list := make([]string, 0)
		ops := Options{Trail: "", TrailLog: &list}

		// --- When ---
		have := ops.LogTrail()

		// --- Then ---
		affirm.DeepEqual(t, []string{}, list)
		affirm.DeepEqual(t, []string{}, *ops.TrailLog)
		affirm.DeepEqual(t, have, ops)
	})

	t.Run("does not panic when nil", func(t *testing.T) {
		// --- Given ---
		ops := Options{Trail: "abc"}

		// --- When ---
		have := ops.LogTrail()

		// --- Then ---
		affirm.DeepEqual(t, have, ops)
	})
}

func Test_Options_StructTrail_tabular(t *testing.T) {
	tt := []struct {
		testN string

		trail   string
		typName string
		fldName string
		want    string
	}{
		{"no trail and type", "", "type", "", "type"},
		{"no trail and field", "", "", "field", "field"},
		{"no trail and type and field", "", "type", "field", "type.field"},
		{"trail and type", "trail", "type", "", "trail"},
		{"trail and field", "trail", "", "field", "trail.field"},
		{"trail and type and field", "trail", "type", "field", "trail.field"},
		{"trail[] and type", "[]", "type", "", "[]"},
		{"trail[] and field", "[]", "", "field", "[].field"},
		{"trail[] and type and field", "[]", "type", "field", "[].field"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ops := Options{Trail: tc.trail}

			// --- When ---
			have := ops.StructTrail(tc.typName, tc.fldName)

			// --- Then ---
			affirm.Equal(t, ops.Trail, tc.trail)
			affirm.Equal(t, tc.want, have.Trail)
		})
	}
}

func Test_Options_MapTrail_tabular(t *testing.T) {
	tt := []struct {
		testN string

		trail string
		key   string
		want  string
	}{
		{"empty trail with a key", "", "key", "map[key]"},
		{"trail ends with index", "[1]", "key", "[1]map[key]"},
		{"trail ends with index", "[1]", "key", "[1]map[key]"},
		{"not empty trail", "field", "key", "field[key]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ops := Options{Trail: tc.trail}

			// --- When ---
			have := ops.MapTrail(tc.key)

			// --- Then ---
			affirm.Equal(t, ops.Trail, tc.trail)
			affirm.Equal(t, tc.want, have.Trail)
		})
	}
}

func Test_Options_ArrTrail_tabular(t *testing.T) {
	tt := []struct {
		testN string

		trail string
		kind  string
		key   int
		want  string
	}{
		{"empty trail with a key", "", "", 1, "[1]"},
		{"empty trail with a key", "", "kind", 1, "<kind>[1]"},
		{"trail ends with index", "[1]", "", 2, "[1][2]"},
		{"not empty trail", "field", "", 1, "field[1]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ops := Options{Trail: tc.trail}

			// --- When ---
			have := ops.ArrTrail(tc.kind, tc.key)

			// --- Then ---
			affirm.Equal(t, ops.Trail, tc.trail)
			affirm.Equal(t, tc.want, have.Trail)
		})
	}
}

func Test_FieldName(t *testing.T) {
	t.Run("empty trail", func(t *testing.T) {
		// --- Given ---
		ops := Options{SkipUnexported: true}

		// --- When ---
		have := FieldName(ops, "myType")("myField")

		// --- Then ---
		affirm.Equal(t, "myType.myField", have(Options{}).Trail)
		affirm.Equal(t, true, have(Options{}).SkipUnexported)
	})

	t.Run("existing trail", func(t *testing.T) {
		// --- Given ---
		ops := Options{Trail: "ABC[1]", SkipUnexported: true}

		// --- When ---
		have := FieldName(ops, "myType")("myField")

		// --- Then ---
		affirm.Equal(t, "ABC[1].myField", have(Options{}).Trail)
		affirm.Equal(t, true, have(Options{}).SkipUnexported)
	})
}

func Test_AddRows(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		// --- Given ---
		ops := Options{}
		msg := notice.New("header").Want("%s", "want")

		// --- When ---
		have := AddRows(ops, msg)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		affirm.Equal(t, "header:\n  want: want", have.Error())
	})

	t.Run("with trail", func(t *testing.T) {
		// --- Given ---
		ops := Options{Trail: "type.field"}
		msg := notice.New("header").Want("%s", "want")

		// --- When ---
		have := AddRows(ops, msg)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		wMsg := "" +
			"header:\n" +
			"  trail: type.field\n" +
			"   want: want"
		affirm.Equal(t, wMsg, have.Error())
	})

	t.Run("with comment", func(t *testing.T) {
		// --- Given ---
		ops := Options{Comment: "abc"}
		msg := notice.New("header").Want("%s", "want")

		// --- When ---
		have := AddRows(ops, msg)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		wMsg := "" +
			"header:\n" +
			"  comment: abc\n" +
			"     want: want"
		affirm.Equal(t, wMsg, have.Error())
	})

	t.Run("with trail and comment", func(t *testing.T) {
		// --- Given ---
		ops := Options{Trail: "type.field", Comment: "abc"}
		msg := notice.New("header").Want("%s", "want")

		// --- When ---
		have := AddRows(ops, msg)

		// --- Then ---
		affirm.Equal(t, true, core.Same(msg, have))
		wMsg := "" +
			"header:\n" +
			"    trail: type.field\n" +
			"  comment: abc\n" +
			"     want: want"
		affirm.Equal(t, wMsg, have.Error())
	})
}
