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
)

func Test_RegisterTypeChecker(t *testing.T) {
	t.Setenv("___", "___")
	affirm.Nil(t, typeCheckers)
	origLog := globLog
	buf := &bytes.Buffer{}
	globLog = log.New(buf, "", 0)
	chk := func(_, _ any, _ ...Option) error { return errors.New("123456") }
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
		typeCheckers = map[reflect.Type]Check{reflect.TypeOf(custom{}): chk}

		// --- When ---
		fn := func() { RegisterTypeChecker(custom{}, chk) }
		msg := affirm.Panic(t, fn)

		// --- Then ---
		wMsg := "cannot overwrite an existing checker"
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
		affirm.Equal(t, "cannot register a nil checker", *msg)
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
	chk := func(_, _ any, _ ...Option) error { return errors.New("123456") }
	t.Cleanup(func() { globLog = origLog; typeCheckers = nil })

	t.Run("setting", func(t *testing.T) {
		// --- Given ---
		ops := Options{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })

		// --- When ---
		have := WithTypeChecker(123, chk)(ops)

		// --- Then ---
		wChk := core.Same(chk, have.TypeCheckers[reflect.TypeOf(123)])
		affirm.Equal(t, true, wChk)
		affirm.Equal(t, "", buf.String())
	})

	t.Run("overwriting global checker", func(t *testing.T) {
		// --- Given ---
		type custom struct{}
		t.Cleanup(func() { typeCheckers = nil; buf.Reset() })

		RegisterTypeChecker(custom{}, chk) // The first call.
		buf.Reset()                        // Test later the log is empty.

		ops := Options{}

		// --- When ---
		have := WithTypeChecker(custom{}, chk)(ops)

		// --- Then ---
		wChk := core.Same(chk, have.TypeCheckers[reflect.TypeOf(custom{})])
		affirm.Equal(t, true, wChk)
		wMsg := "Overwriting the global type checker for: check.custom\n"
		affirm.Equal(t, wMsg, buf.String())
	})
}

func Test_WithTrailChecker(t *testing.T) {
	// --- Given ---
	ops := Options{}
	chk := func(want, have any, opts ...Option) error { return nil }

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

func Test_WithOptions(t *testing.T) {
	// --- Given ---
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
			UseAny:         true,
			Dumpers: map[reflect.Type]dump.Dumper{
				reflect.TypeOf(123): dump.Dumper(nil),
			},
			MaxDepth: 6,
			Indent:   2,
			TabWidth: 4,
		},
		TimeFormat:     time.RFC3339,
		Recent:         123,
		Trail:          "trail",
		TrailLog:       &trailLog,
		TypeCheckers:   make(map[reflect.Type]Check),
		TrailCheckers:  make(map[string]Check),
		SkipTrails:     make([]string, 0),
		SkipUnexported: true,
		now:            time.Now,
	}

	// --- When ---
	have := WithOptions(ops)(Options{})

	// --- Then ---
	affirm.Equal(t, true, core.Same(ops.Dumper.Dumpers, have.Dumper.Dumpers))
	affirm.Equal(t, true, core.Same(ops.TrailLog, have.TrailLog))
	affirm.Equal(t, true, core.Same(ops.TypeCheckers, have.TypeCheckers))
	affirm.Equal(t, true, core.Same(ops.TrailCheckers, have.TrailCheckers))
	affirm.Equal(t, true, core.Same(ops.SkipTrails, have.SkipTrails))
	affirm.Equal(t, true, core.Same(ops.now, have.now))

	ops.now = nil
	have.now = nil
	affirm.Equal(t, true, reflect.DeepEqual(ops, have))

	// When those fail, add fields above.
	affirm.Equal(t, 12, reflect.ValueOf(have.Dumper).NumField())
	affirm.Equal(t, 10, reflect.ValueOf(have).NumField())
}

func Test_DefaultOptions(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have := DefaultOptions()

		// --- Then ---
		affirm.Equal(t, false, have.Dumper.PtrAddr)
		affirm.Equal(t, DefaultDumpTimeFormat, have.Dumper.TimeFormat)

		affirm.Equal(t, DefaultParseTimeFormat, have.TimeFormat)
		affirm.Equal(t, DefaultRecentDuration, have.Recent)
		affirm.Equal(t, "", have.Trail)
		affirm.Equal(t, true, have.TrailLog == nil)
		affirm.Equal(t, true, have.TypeCheckers == nil)
		affirm.Equal(t, true, have.TrailCheckers == nil)
		affirm.Equal(t, true, have.SkipTrails == nil)
		affirm.Equal(t, false, have.SkipUnexported)
		affirm.Equal(t, true, core.Same(time.Now, have.now))
		affirm.Equal(t, 10, reflect.ValueOf(have).NumField())
	})

	t.Run("with options", func(t *testing.T) {
		// --- When ---
		have := DefaultOptions(WithTrail("type.field"))

		// --- Then ---
		affirm.Equal(t, false, have.Dumper.PtrAddr)
		affirm.Equal(t, DefaultDumpTimeFormat, have.Dumper.TimeFormat)

		affirm.Equal(t, DefaultParseTimeFormat, have.TimeFormat)
		affirm.Equal(t, DefaultRecentDuration, have.Recent)
		affirm.Equal(t, "type.field", have.Trail)
		affirm.Equal(t, true, have.TrailLog == nil)
		affirm.Equal(t, true, have.TypeCheckers == nil)
		affirm.Equal(t, true, have.TrailCheckers == nil)
		affirm.Equal(t, true, have.SkipTrails == nil)
		affirm.Equal(t, false, have.SkipUnexported)
		affirm.Equal(t, true, core.Same(time.Now, have.now))
		affirm.Equal(t, 10, reflect.ValueOf(have).NumField())
	})

	t.Run("TypeCheckers field is a clone of a global map", func(t *testing.T) {
		// --- Given ---
		t.Setenv("___", "___")
		affirm.Nil(t, typeCheckers)
		chk := func(_, _ any, _ ...Option) error { return errors.New("123456") }
		t.Cleanup(func() { typeCheckers = nil })
		typeCheckers = map[reflect.Type]Check{reflect.TypeOf(123): chk}

		// --- When ---
		have := DefaultOptions()

		// --- Then ---
		affirm.Equal(t, false, core.Same(typeCheckers, have.TypeCheckers))
	})
}

func Test_Options_logTrail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		list := make([]string, 0)
		ops := Options{Trail: "abc", TrailLog: &list}

		// --- When ---
		have := ops.logTrail()

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
		have := ops.logTrail()

		// --- Then ---
		affirm.DeepEqual(t, []string{}, list)
		affirm.DeepEqual(t, []string{}, *ops.TrailLog)
		affirm.DeepEqual(t, have, ops)
	})

	t.Run("does not panic when nil", func(t *testing.T) {
		// --- Given ---
		ops := Options{Trail: "abc"}

		// --- When ---
		have := ops.logTrail()

		// --- Then ---
		affirm.DeepEqual(t, have, ops)
	})
}

func Test_Options_structTrail_tabular(t *testing.T) {
	tt := []struct {
		testN string

		trail   string
		typName string
		fldName string
		want    string
	}{
		{"no trail and type", "", "type", "", "type"},                         // 1
		{"no trail and field", "", "", "field", "field"},                      // 2
		{"no trail and type and field", "", "type", "field", "type.field"},    // 3
		{"trail and type", "trail", "type", "", "trail"},                      // 4
		{"trail and field", "trail", "", "field", "trail.field"},              // 5
		{"trail and type and field", "trail", "type", "field", "trail.field"}, // 6
		{"trail[] and type", "[]", "type", "", "[]"},                          // 7
		{"trail[] and field", "[]", "", "field", "[].field"},                  // 8
		{"trail[] and type and field", "[]", "type", "field", "[].field"},     // 9
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ops := Options{Trail: tc.trail}

			// --- When ---
			have := ops.structTrail(tc.typName, tc.fldName)

			// --- Then ---
			affirm.Equal(t, ops.Trail, tc.trail)
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_Options_mapTrail_tabular(t *testing.T) {
	tt := []struct {
		testN string

		trail string
		key   string
		want  string
	}{
		{"empty trail with key", "", "key", "map[key]"},
		{"trail ends with index", "[1]", "key", "[1]map[key]"},
		{"trail ends with index", "[1]", "key", "[1]map[key]"},
		{"not empty trail", "field", "key", "field[key]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ops := Options{Trail: tc.trail}

			// --- When ---
			have := ops.mapTrail(tc.key)

			// --- Then ---
			affirm.Equal(t, ops.Trail, tc.trail)
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_Options_arrTrail_tabular(t *testing.T) {
	tt := []struct {
		testN string

		trail string
		kind  string
		key   int
		want  string
	}{
		{"empty trail with key", "", "", 1, "[1]"},
		{"empty trail with key", "", "kind", 1, "<kind>[1]"},
		{"trail ends with index", "[1]", "", 2, "[1][2]"},
		{"not empty trail", "field", "", 1, "field[1]"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ops := Options{Trail: tc.trail}

			// --- When ---
			have := ops.arrTrail(tc.kind, tc.key)

			// --- Then ---
			affirm.Equal(t, ops.Trail, tc.trail)
			affirm.Equal(t, tc.want, have)
		})
	}
}
