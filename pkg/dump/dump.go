// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package dump can render string representation of any value.
package dump

import (
	"reflect"
	"time"
)

// Strings used by dump package to indicate special values.
const (
	valNil        = "nil"       // The [reflect.Value] is nil.
	valAddr       = "<addr>"    // The [reflect.Value] is an address.
	valFunc       = "<func>"    // The [reflect.Value] is a function.
	valChan       = "<chan>"    // The [reflect.Value] is a channel.
	valInvalid    = "<invalid>" // The [reflect.Value] is invalid.
	valMaxNesting = "<...>"     // The maximum nesting reached.

	// The [reflect.Value] is unexpected in given context.
	valErrUsage = "<dump-usage-error>"
)

// Package wide default configuration.
const (
	// DefaultTimeFormat is default format for parsing time strings.
	DefaultTimeFormat = time.RFC3339Nano

	// DefaultDepth is default depth when dumping values recursively.
	DefaultDepth = 6

	// DefaultIndent is default additional indent when dumping values.
	DefaultIndent = 0

	// DefaultTabWith is default tab width in spaces.
	DefaultTabWith = 2
)

// Package wide configuration.
var (
	// TimeFormat is configurable format for dumping [time.Time] values.
	TimeFormat = DefaultTimeFormat

	// Depth is configurable depth when dumping values recursively.
	Depth = DefaultDepth

	// Indent is configurable additional indent when dumping values.
	Indent = DefaultIndent

	// TabWidth is configurable tab width in spaces.
	TabWidth = DefaultTabWith
)

// Types for built-in dumpers.
var (
	typDur      = reflect.TypeOf(time.Duration(0))
	typLocation = reflect.TypeOf(time.Location{})
	typTime     = reflect.TypeOf(time.Time{})
)

var nilVal = reflect.ValueOf(nil)

// Dumper represents function signature for value dumpers.
type Dumper func(dmp Dump, level int, val reflect.Value) string

// Option represents [NewConfig] option.
type Option func(*Dump)

// WithFlat is option for [New] which makes [Dump] display values in one line.
func WithFlat(dmp *Dump) { dmp.Flat = true }

// WithCompact is option for [New] which makes [Dump] display values without
// unnecessary whitespaces.
func WithCompact(dmp *Dump) { dmp.Compact = true }

// WithPtrAddr is option for [New] which makes [Dump] display pointer addresses.
func WithPtrAddr(dmp *Dump) { dmp.PtrAddr = true }

// WithTimeFormat is option for [New] which makes [Dump] display [time.Time]
// using given format. The format might be standard Go time formating layout or
// one of the custom values - see [Dump.TimeFormat] for more details.
func WithTimeFormat(format string) Option {
	return func(dmp *Dump) { dmp.TimeFormat = format }
}

// WithDumper adds custom [Dumper] to the config.
func WithDumper(typ any, dumper Dumper) Option {
	return func(dmp *Dump) { dmp.Dumpers[reflect.TypeOf(typ)] = dumper }
}

// WithMaxDepth is option for [New] which controls maximum nesting when
// bumping recursive types.
func WithMaxDepth(maximum int) Option {
	return func(dmp *Dump) { dmp.MaxDepth = maximum }
}

// WithIndent is option for [New] which sets additional indentation to apply to
// dumped values.
func WithIndent(n int) Option {
	return func(dmp *Dump) { dmp.Indent = n }
}

// WithTabWidth is option for [New] setting tab width in spaces.
func WithTabWidth(n int) Option {
	return func(dmp *Dump) { dmp.TabWidth = n }
}

// Dump implements logic for dumping values and types.
type Dump struct {
	// Display values on one line.
	Flat bool

	// Do not use any indents or whitespace separators.
	Compact bool

	// Controls how [time.Time] is formated.
	//
	// Aside form Go time formating layouts following custom formats are
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
}

// New returns new instance of [Dump].
func New(opts ...Option) Dump {
	dmp := Dump{
		TimeFormat: TimeFormat,
		PrintType:  true,
		UseAny:     true,
		Dumpers:    make(map[reflect.Type]Dumper),
		MaxDepth:   Depth,
		Indent:     Indent,
		TabWidth:   TabWidth,
	}
	for _, opt := range opts {
		opt(&dmp)
	}

	if _, ok := dmp.Dumpers[typTime]; !ok {
		dmp.Dumpers[typTime] = GetTimeDumper(dmp.TimeFormat)
	}

	if _, ok := dmp.Dumpers[typLocation]; !ok {
		dmp.Dumpers[typLocation] = zoneDumper
	}

	if _, ok := dmp.Dumpers[typDur]; !ok {
		dmp.Dumpers[typDur] = GetDurDumper(dmp.DurationFormat)
	}
	return dmp
}

// Any dumps any value to its string representation.
func (dmp Dump) Any(val any) string {
	return dmp.value(0, reflect.ValueOf(val))
}

// Value dumps a [reflect.Value] representation of a value as a string.
func (dmp Dump) Value(val reflect.Value) string {
	return dmp.value(0, val)
}

// value dumps given a value as a string.
//
// nolint: cyclop
func (dmp Dump) value(lvl int, val reflect.Value) string {
	if lvl > dmp.MaxDepth {
		return valMaxNesting
	}

	var str string // One or more lines representing passed value.

	valKnd := val.Kind()
	if valKnd != reflect.Invalid {
		if fn, ok := dmp.Dumpers[val.Type()]; ok {
			return fn(dmp, lvl, val)
		}
	}

	switch valKnd {
	case reflect.Invalid:
		str = valInvalid
		if nilVal == val { // nolint: govet
			str = valNil
		}

	case reflect.Bool, reflect.Int:
		str = simpleDumper(dmp, lvl, val)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = simpleDumper(dmp, lvl, val)

	case reflect.Uint:
		str = simpleDumper(dmp, lvl, val)

	case reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = simpleDumper(dmp, lvl, val)

	case reflect.Uint8:
		str = hexPtrDumper(dmp, lvl, val)

	case reflect.Uintptr:
		str = hexPtrDumper(dmp, lvl, val)

	case reflect.Float32, reflect.Float64:
		str = simpleDumper(dmp, lvl, val)

	case reflect.Complex64, reflect.Complex128:
		str = complexDumper(dmp, lvl, val)

	case reflect.Array:
		str = arrayDumper(dmp, lvl, val)

	case reflect.Chan:
		str = chanDumper(dmp, lvl, val)

	case reflect.Func:
		str = funcDumper(dmp, lvl, val)

	case reflect.Interface:
		str = dmp.value(lvl, val.Elem())

	case reflect.Map:
		str = mapDumper(dmp, lvl, val)

	case reflect.Pointer:
		if val.IsNil() {
			str = valNil
		} else {
			str = dmp.value(lvl, val.Elem())
		}

	case reflect.Slice:
		str = sliceDumper(dmp, lvl, val)

	case reflect.String:
		str = simpleDumper(dmp, lvl, val)

	case reflect.Struct:
		str = structDumper(dmp, lvl, val)

	case reflect.UnsafePointer:
		str = hexPtrDumper(dmp, lvl, val)
	}

	return str
}
