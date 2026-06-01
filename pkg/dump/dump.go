// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package dump provides a configurable renderer that turns any Go value
// into a human-readable string representation.
//
// It is the foundation for high-quality diagnostic output across the
// module: used by [github.com/ctx42/testing/pkg/mock] for expectation
// diffs, [github.com/ctx42/testing/pkg/notice] for structured messages,
// [github.com/ctx42/testing/pkg/assert] failures, and golden-file
// comparisons in [github.com/ctx42/testing/pkg/goldy].
//
// The package supports deep recursion with cycle detection, custom type
// dumpers, fine-grained formatting options (flat, compact, indentation,
// time/duration formats, etc.), and both simple one-shot use via [Any]
// and reusable configured instances via [New].
//
// See the package [README] for usage patterns, configuration examples,
// and extensibility. See [examples_test.go] for executable demonstrations
// (many of which are also wired into the README via gmdoceg markers).
//
// Key types and entry points:
//   - [New] + [Dump] — the core configurable dumper
//   - [Any] — one-shot convenience using defaults
//   - [Option] functions (WithFlat, WithTimeFormat, WithDumper, ...)
//   - [RegisterTypeDumper] — optional global custom type handlers (see
//     package documentation for the full customization model)
//   - Time/duration sentinels: [TimeAsUnix], [DurAsString], ...
//   - [Diff] / [DiffValue] — for rich want/have comparisons
//
// The customization model supports both global defaults (via
// [RegisterTypeDumper] and package-level variables) and per-instance
// configuration via options passed to [New].
package dump

import (
	"fmt"
	"log"
	"maps"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ctx42/testing/internal/diff"
)

// globLog is a global logger used package-wide.
var globLog = log.New(os.Stderr, "*** DUMP ", log.Llongfile)

// Special sentinel strings that appear in rendered output for values that
// cannot be printed in the normal way or that carry semantic meaning.
const (
	// ValNotNil is shown for any non-nil value when only nil-ness matters.
	ValNotNil = "<not-nil>"

	// ValNil is the representation of a nil value.
	ValNil = "nil"

	// ValAddr appears when pointer addresses are requested ([WithPtrAddr]).
	ValAddr = "<addr>"

	// ValFunc is shown for function values (which have no printable content).
	ValFunc = "<func>"

	// ValChan is shown for channel values.
	ValChan = "<chan>"

	// ValInvalid is the representation of an invalid reflect.Value.
	ValInvalid = "<invalid>"

	// ValMaxNesting is emitted when recursion exceeds [Dump.MaxDepth].
	ValMaxNesting = "<...>"

	// ValEmpty is used for empty collections or zero-length strings in some
	// contexts.
	ValEmpty = "<empty>"

	// ValErrUsage is emitted when a dumper is invoked on a value of the
	// wrong type (internal contract violation).
	ValErrUsage = "<dump-usage-error>"

	// ValCannotPrint is shown for values that cannot be printed, typically
	// because they are unexported and [Dump.PrintPrivate] is false.
	ValCannotPrint = "<dump-cannot-print>"
)

// Default values used by [New] when no options override them.
const (
	// DefaultTimeFormat is the default [time.Time] layout.
	DefaultTimeFormat = time.RFC3339Nano

	// DefaultDepth is the default maximum recursion depth.
	DefaultDepth = 6

	// DefaultIndent is the default extra indentation per level.
	DefaultIndent = 0

	// DefaultTabWith is the default number of spaces per indentation level.
	DefaultTabWith = 2
)

// Package-level defaults. These are read by [New] at construction time and
// can be changed to affect all subsequent default-configured dumpers.
// For per-instance configuration prefer the [Option] functions.
var (
	// TimeFormat is the package default for [time.Time] formatting.
	TimeFormat = DefaultTimeFormat

	// Depth is the package default maximum recursion depth.
	Depth = DefaultDepth

	// Indent is the package default extra indentation.
	Indent = DefaultIndent

	// TabWidth is the package default spaces per indentation level.
	TabWidth = DefaultTabWith
)

// Types for built-in dumpers.
var (
	typDur      = reflect.TypeFor[time.Duration]()
	typLocation = reflect.TypeFor[time.Location]()
	typTime     = reflect.TypeFor[time.Time]()
	typError    = reflect.TypeFor[error]()
)

var nilVal = reflect.ValueOf(nil)

// Dumper is the signature for custom value renderers.
//
// Receives the active [Dump] (for recursive calls), current nesting
// level, and the [reflect.Value] to render. Must return the string
// representation. Used with [RegisterTypeDumper] and [WithDumper].
type Dumper func(dmp Dump, level int, val reflect.Value) string

// typeDumpers holds globally registered custom dumpers (see
// [RegisterTypeDumper]).
var typeDumpers map[reflect.Type]Dumper

// RegisterTypeDumper registers a custom [Dumper] that will be used for
// all values of the given type by default (across all [New] calls that
// do not override it with [WithDumper]).
//
// It panics if a dumper for the same type is already registered.
// The registration is logged at debug level for visibility into
// overrides.
func RegisterTypeDumper(typ any, dmp Dumper) {
	if dmp == nil {
		panic("cannot register a nil type dumper")
	}
	if typeDumpers == nil {
		typeDumpers = make(map[reflect.Type]Dumper)
	}
	rt := reflect.TypeOf(typ)
	msg := fmt.Sprintf("Registering type dumper for: %s", rt)
	if _, ok := typeDumpers[rt]; ok {
		panic("cannot overwrite an existing type dumper: " + msg)
	}
	_ = globLog.Output(2, msg)
	typeDumpers[rt] = dmp
}

// Option configures a [Dump] created by [New].
type Option func(*Dump)

// WithFlat makes the dumper render values on a single line (no newlines).
func WithFlat(dmp *Dump) { dmp.Flat = true }

// WithFlatStrings sets the threshold (number of newlines) above which
// strings are forced to multiline rendering even under [WithFlat].
// Zero disables the feature.
func WithFlatStrings(n int) Option {
	return func(dmp *Dump) { dmp.FlatStrings = n }
}

// WithCompact removes most whitespace separators for the most compact output.
func WithCompact(dmp *Dump) { dmp.Compact = true }

// WithAlwaysMultiline forces all strings (even single-line ones) to be
// rendered as multiline blocks.
func WithAlwaysMultiline(dmp *Dump) { dmp.AlwaysMultiline = true }

// WithPtrAddr includes pointer addresses in the output (normally hidden).
func WithPtrAddr(dmp *Dump) { dmp.PtrAddr = true }

// WithNoPrivate suppresses printing of unexported struct fields.
func WithNoPrivate(dmp *Dump) { dmp.PrintPrivate = false }

// WithTimeFormat sets the format used for [time.Time] values. In addition
// to standard Go layouts, the special values [TimeAsUnix] and [TimeAsGoString]
// are supported. See [Dump.TimeFormat].
func WithTimeFormat(format string) Option {
	return func(dmp *Dump) { dmp.TimeFormat = format }
}

// WithDumper registers a per-instance custom [Dumper] for the given type.
// It overrides any globally registered dumper for that type (see
// [RegisterTypeDumper]).
func WithDumper(typ any, dumper Dumper) Option {
	return func(dmp *Dump) {
		rt := reflect.TypeOf(typ)
		if _, ok := typeDumpers[rt]; ok {
			format := "Overwriting the global type dumper for: %s"
			_ = globLog.Output(2, fmt.Sprintf(format, rt))
		}
		dmp.Dumpers[rt] = dumper
	}
}

// WithMaxDepth limits recursion depth when dumping nested or recursive
// structures (default is [DefaultDepth]).
func WithMaxDepth(maximum int) Option {
	return func(dmp *Dump) { dmp.MaxDepth = maximum }
}

// WithIndent sets additional indentation (in spaces) applied to each level.
func WithIndent(n int) Option {
	return func(dmp *Dump) { dmp.Indent = n }
}

// WithTabWidth sets the number of spaces used for one level of indentation.
func WithTabWidth(n int) Option {
	return func(dmp *Dump) { dmp.TabWidth = n }
}

// Dump holds configuration and state for rendering values. Obtain an
// instance via [New] (or use the package-level [Any] for defaults).
//
// Most fields are exported for advanced use or inspection; the preferred
// way to configure is through the [Option] functions.
type Dump struct {
	// Flat renders everything on a single line (no newlines or indentation).
	Flat bool

	// FlatStrings treats strings longer than this as multiline (0 disables).
	FlatStrings int

	// Compact removes almost all whitespace for the densest output.
	Compact bool

	// AlwaysMultiline forces strings to be rendered as multiline blocks.
	AlwaysMultiline bool

	// TimeFormat controls formatting of [time.Time]. In addition to standard
	// Go layouts, the special values [TimeAsUnix], [TimeAsGoString], and
	// [TimeAsRFC3339] are recognized. Default: [time.RFC3339Nano].
	TimeFormat string

	// DurationFormat controls formatting of [time.Duration]. Special values:
	// [DurAsString] (default) and [DurAsSeconds].
	DurationFormat string

	// PtrAddr shows pointer addresses (normally hidden for readability).
	PtrAddr bool

	// PrintType includes type information in the output.
	PrintType bool

	// PrintPrivate includes unexported struct field values.
	PrintPrivate bool

	// UseAny uses the "any" alias instead of "interface{}" in type names.
	UseAny bool

	// Dumpers holds per-instance custom type renderers. These take
	// precedence over the global registrations from [RegisterTypeDumper].
	// Default built-in dumpers for time.Time, time.Duration and
	// time.Location are added automatically by [New] unless overridden.
	Dumpers map[reflect.Type]Dumper

	// MaxDepth limits recursion depth for nested/recursive data (default
	// [DefaultDepth]). Deeper values render as [ValMaxNesting].
	MaxDepth int

	// Indent is additional spaces of indentation per nesting level.
	Indent int

	// TabWidth is the number of spaces that represent one indentation level.
	TabWidth int

	// internal: mirrors FlatStrings for nested string fields inside structs
	// when we want them forced flat.
	flatStrings bool
}

// New creates a new [Dump] with sensible defaults, applies the provided
// options, and ensures built-in dumpers for time.Time, time.Duration and
// time.Location are present (unless overridden via options).
func New(opts ...Option) Dump {
	dmp := Dump{
		FlatStrings:  200,
		TimeFormat:   TimeFormat,
		PrintType:    true,
		PrintPrivate: true,
		UseAny:       true,
		Dumpers:      maps.Clone(typeDumpers),
		MaxDepth:     Depth,
		Indent:       Indent,
		TabWidth:     TabWidth,
	}
	if dmp.Dumpers == nil {
		dmp.Dumpers = make(map[reflect.Type]Dumper)
	}
	for _, opt := range opts {
		opt(&dmp)
	}

	if _, ok := dmp.Dumpers[typTime]; !ok {
		dmp.Dumpers[typTime] = GetTimeDumper(dmp.TimeFormat)
	}

	if _, ok := dmp.Dumpers[typLocation]; !ok {
		dmp.Dumpers[typLocation] = ZoneDumper
	}

	if _, ok := dmp.Dumpers[typDur]; !ok {
		dmp.Dumpers[typDur] = GetDurDumper(dmp.DurationFormat)
	}
	return dmp
}

// Any returns a string representation of val using default configuration.
// It is the simplest entry point for one-off diagnostics.
func Any(val any) string { return New().Any(val) }

// Any renders val using this Dump's configuration.
func (dmp Dump) Any(val any) string {
	str, _ := dmp.value(0, reflect.ValueOf(val))
	return str
}

// Diff renders want and have, then returns their string forms plus a unified
// diff (empty when identical). Useful for rich assertion failure messages.
func (dmp Dump) Diff(want, have any) (string, string, string) {
	wVal := reflect.ValueOf(want)
	hVal := reflect.ValueOf(have)
	return dmp.DiffValue(wVal, hVal)
}

// DiffValue is like [Diff] but accepts reflect.Value directly.
func (dmp Dump) DiffValue(wVal, hVal reflect.Value) (string, string, string) {
	// Format values for display.
	wStr, _ := dmp.value(0, wVal)
	hStr, _ := dmp.value(0, hVal)
	if wStr == hStr {
		return wStr, hStr, ""
	}

	if wStr == ValNil || hStr == ValNil {
		return wStr, hStr, ""
	}

	// If one of the values is multiline, force the other to be as well.
	wMlStr := strings.Contains(wStr, "\n")
	hMlStr := strings.Contains(hStr, "\n")
	if wMlStr != hMlStr {
		dmp2 := dmp
		dmp2.Flat = false
		dmp2.FlatStrings = 0
		if wMlStr {
			hStr, _ = dmp2.value(0, hVal)
		} else {
			wStr, _ = dmp2.value(0, wVal)
		}
	}

	// Format values for diff.
	wDiffStr, _ := dmp.forDiff(wVal)
	hDiffStr, _ := dmp.forDiff(hVal)
	wDiffMlStr := strings.Contains(wDiffStr, "\n")
	hDiffMlStr := strings.Contains(hDiffStr, "\n")

	// If both values are not multiline, don't show diff.
	if !wDiffMlStr && !hDiffMlStr {
		return wStr, hStr, ""
	}

	edits := diff.Strings(hDiffStr, wDiffStr)
	// Error can't happen: edits are consistent.
	unified, _ := diff.CtxToUnified("want", "have", hDiffStr, edits, 2)
	return wStr, hStr, strings.TrimRight(unified, "\n")
}

// forDiff prepares a value for diffing by formatting it into a string. Returns
// the string representation of the value and kind.
func (dmp Dump) forDiff(val reflect.Value) (string, reflect.Kind) {
	dmp.Flat = false
	dmp.FlatStrings = 0
	dmp.Compact = false

	str, knd := dmp.value(0, val)
	if s, err := strconv.Unquote(str); err == nil {
		str = s
	}
	return str, knd
}

// Value renders a reflect.Value using this configuration.
func (dmp Dump) Value(val reflect.Value) string {
	str, _ := dmp.value(0, val)
	return str
}

// value dumps given a value as a string.
//
// nolint: cyclop
func (dmp Dump) value(lvl int, val reflect.Value) (string, reflect.Kind) {
	if lvl > dmp.MaxDepth {
		return ValMaxNesting, reflect.Invalid
	}

	var str string // One or more lines representing passed value.

	knd := val.Kind()
	if knd != reflect.Invalid {
		if fn, ok := dmp.Dumpers[val.Type()]; ok {
			return fn(dmp, lvl, val), knd
		}
	}

	if val.IsValid() {
		typ := val.Type()
		// Special case for type: error.
		if typ == typError || typ.String() == "*errors.errorString" {
			str = ValNil
			if !val.IsNil() {
				if val.CanInterface() {
					err := val.Interface().(error) // nolint: forcetypeassert
					str = fmt.Sprintf("%q", err.Error())
				} else {
					str = ValCannotPrint
				}
			}
			prn := NewPrinter(dmp)
			return prn.Tab(dmp.Indent + lvl).Write(str).String(), knd
		}
	}

	switch knd {
	case reflect.Invalid:
		str = ValInvalid
		if nilVal == val { // nolint: govet
			str = ValNil
		}

	case reflect.Bool, reflect.Int:
		str = SimpleDumper(dmp, lvl, val)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = SimpleDumper(dmp, lvl, val)

	case reflect.Uint:
		str = SimpleDumper(dmp, lvl, val)

	case reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = SimpleDumper(dmp, lvl, val)

	case reflect.Uint8:
		str = HexPtrDumper(dmp, lvl, val)

	case reflect.Uintptr:
		str = HexPtrDumper(dmp, lvl, val)

	case reflect.Float32, reflect.Float64:
		str = SimpleDumper(dmp, lvl, val)

	case reflect.Complex64, reflect.Complex128:
		str = ComplexDumper(dmp, lvl, val)

	case reflect.Array:
		str = ArrayDumper(dmp, lvl, val)

	case reflect.Chan:
		str = ChanDumper(dmp, lvl, val)

	case reflect.Func:
		str = FuncDumper(dmp, lvl, val)

	case reflect.Interface:
		str, knd = dmp.value(lvl, val.Elem())

	case reflect.Map:
		str = MapDumper(dmp, lvl, val)

	case reflect.Pointer:
		if val.IsNil() {
			str = ValNil
		} else {
			str, knd = dmp.value(lvl, val.Elem())
		}

	case reflect.Slice:
		str = SliceDumper(dmp, lvl, val)

	case reflect.String:
		str = SimpleDumper(dmp, lvl, val)

	case reflect.Struct:
		str = StructDumper(dmp, lvl, val)

	case reflect.UnsafePointer:
		str = HexPtrDumper(dmp, lvl, val)
	}

	return str, knd
}
