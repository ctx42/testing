// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package dump can render a string representation of any value.
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

// Strings used by dump package to indicate special values.
const (
	ValNotNil     = "<not-nil>"          // Represents any not-nil value.
	ValNil        = "nil"                // The [reflect.Value] is nil.
	ValAddr       = "<addr>"             // The [reflect.Value] is an address.
	ValFunc       = "<func>"             // The [reflect.Value] is a function.
	ValChan       = "<chan>"             // The [reflect.Value] is a channel.
	ValInvalid    = "<invalid>"          // The [reflect.Value] is invalid.
	ValMaxNesting = "<...>"              // The maximum nesting reached.
	ValEmpty      = "<empty>"            // Empty value.
	ValErrUsage   = "<dump-usage-error>" // The [reflect.Value] is unexpected in the given context.
)

// Package wide default configuration.
const (
	// DefaultTimeFormat is default format for parsing time strings.
	DefaultTimeFormat = time.RFC3339Nano

	// DefaultDepth is the default depth when dumping values recursively.
	DefaultDepth = 6

	// DefaultIndent is default additional indent when dumping values.
	DefaultIndent = 0

	// DefaultTabWith is the default tab width in spaces.
	DefaultTabWith = 2
)

// Package-wide configuration.
var (
	// TimeFormat is configurable format for dumping [time.Time] values.
	TimeFormat = DefaultTimeFormat

	// Depth is configurable depth when dumping values recursively.
	Depth = DefaultDepth

	// Indent is a configurable additional indent when dumping values.
	Indent = DefaultIndent

	// TabWidth is a configurable tab width in spaces.
	TabWidth = DefaultTabWith
)

// Types for built-in dumpers.
var (
	typDur      = reflect.TypeOf(time.Duration(0))
	typLocation = reflect.TypeOf(time.Location{})
	typTime     = reflect.TypeOf(time.Time{})
	typError    = reflect.TypeOf((*error)(nil)).Elem()
)

var nilVal = reflect.ValueOf(nil)

// Dumper represents function signature for value dumpers.
type Dumper func(dmp Dump, level int, val reflect.Value) string

// typeDumpers is the global map of custom dumpers for given types.
var typeDumpers map[reflect.Type]Dumper

// RegisterTypeDumper globally registers a custom dumper for a given type.
// It panics if a dumper for the same type is already registered.
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

// Option represents a [NewConfig] option.
type Option func(*Dump)

// WithFlat is an option for [New] which makes [Dump] display values in one
// line.
func WithFlat(dmp *Dump) { dmp.Flat = true }

// WithFlatStrings configures the maximum length of strings to be represented
// as flat in the output. Strings longer than the specified length may be
// formatted differently, depending on the configuration. This option is
// similar to [WithFlat] but applies specifically to strings based on their
// length. Set to zero to turn this feature off.
func WithFlatStrings(n int) Option {
	return func(dmp *Dump) { dmp.FlatStrings = n }
}

// WithCompact is an option for [New] which makes [Dump] display values without
// unnecessary whitespaces.
func WithCompact(dmp *Dump) { dmp.Compact = true }

// WithAlwaysMultiline is an option for [New] which makes [Dump] treat all
// strings as multiline strings, even if they are not.
func WithAlwaysMultiline(dmp *Dump) { dmp.AlwaysMultiline = true }

// WithPtrAddr is an option for [New] which makes [Dump] display pointer
// addresses.
func WithPtrAddr(dmp *Dump) { dmp.PtrAddr = true }

// WithNoPrivate is an option for [New] which makes [Dump] skip displaying
// values for not exported fields.
func WithNoPrivate(dmp *Dump) { dmp.PrintPrivate = false }

// WithTimeFormat is an option for [New] which makes [Dump] display [time.Time]
// using a given format. The format might be a standard Go time formating
// layout or one of the custom values - see [Dump.TimeFormat] for more details.
func WithTimeFormat(format string) Option {
	return func(dmp *Dump) { dmp.TimeFormat = format }
}

// WithDumper adds custom [Dumper] to the config.
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

// WithMaxDepth is an option for [New] which controls maximum nesting when
// bumping recursive types.
func WithMaxDepth(maximum int) Option {
	return func(dmp *Dump) { dmp.MaxDepth = maximum }
}

// WithIndent is an option for [New] which sets additional indentation to apply
// to dumped values.
func WithIndent(n int) Option {
	return func(dmp *Dump) { dmp.Indent = n }
}

// WithTabWidth is an option for [New] setting tab width in spaces.
func WithTabWidth(n int) Option {
	return func(dmp *Dump) { dmp.TabWidth = n }
}

// Dump implements logic for dumping values and types.
type Dump struct {
	// Display values on one line.
	Flat bool

	// Display strings shorter that given value as with Flat.
	FlatStrings int

	// Do not use any indents or whitespace separators.
	Compact bool

	// Always treat strings as if they were multiline.
	AlwaysMultiline bool

	// Controls how [time.Time] is formated.
	//
	// Aside from Go time formating layouts, the following custom formats are
	// available:
	//
	//  - [TimeAsUnix] - Unix timestamp,
	//
	// By default (empty value) [time.RFC3339Nano] is used.
	TimeFormat string

	// Controls how [time.Duration] is formated.
	//
	// Supports formats:
	//
	//  - [DurAsString]
	//  - [DurAsSeconds]
	DurationFormat string

	// Show pointer addresses.
	PtrAddr bool

	// Print types.
	PrintType bool

	// Controls if the not exported field values should be printed.
	PrintPrivate bool

	// Use "any" instead of "interface{}".
	UseAny bool

	// Custom type dumpers.
	//
	// By default, dumpers for types:
	//   - [time.Time]
	//   - [time.Duration]
	//   - [time.Location]
	//
	// are automatically registered.
	Dumpers map[reflect.Type]Dumper

	// Controls maximum nesting when dumping recursive types.
	// The depth is also used to properly indent values being dumped.
	MaxDepth int

	// How much additional indentation to apply to values being dumped.
	Indent int

	// Default tab with in spaces.
	TabWidth int

	// In cases of nested structures like structs, we want to force string
	// fields to be dumped in flat representation. This value has the same
	// meaning as the Flat option.
	flatStrings bool
}

// New returns new instance of [Dump].
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

// Any dumps any value to its string representation.
func (dmp Dump) Any(val any) string {
	str, _ := dmp.value(0, reflect.ValueOf(val))
	return str
}

// Diff compares two values and returns their formatted representations and
// diff. The first result is the formatted "want" value, the second is the
// formatted "have" value, and the third is the unified diff if they differ. If
// the values are identical, the diff result will be an empty string.
func (dmp Dump) Diff(want, have any) (string, string, string) {
	wVal := reflect.ValueOf(want)
	hVal := reflect.ValueOf(have)
	return dmp.DiffValue(wVal, hVal)
}

// DiffValue works like [Diff] but uses [reflect.Value] instances.
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

// Value dumps a [reflect.Value] representation of a value as a string.
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
				err := val.Interface().(error) // nolint: forcetypeassert
				str = fmt.Sprintf("%q", err.Error())
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
