// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"log"
	"maps"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/ctx42/testing/pkg/dump"
)

// globLog is a global logger used package-wide.
var globLog = log.New(os.Stderr, "*** CHECK ", log.Llongfile)

// Package wide default configuration.
const (
	// DefaultParseTimeFormat is default format for dumping [time.Time] values.
	DefaultParseTimeFormat = time.RFC3339Nano

	// DefaultRecentDuration is default duration when comparing recent dates.
	DefaultRecentDuration = 10 * time.Second

	// DefaultDumpTimeFormat is default format for parsing time strings.
	DefaultDumpTimeFormat = time.RFC3339Nano

	// DefaultDumpDepth is default depth when dumping values recursively in log
	// messages.
	DefaultDumpDepth = 6
)

// Package-wide configuration.
var (
	// ParseTimeFormat is a configurable format for parsing time strings.
	ParseTimeFormat = DefaultParseTimeFormat

	// RecentDuration is a configurable duration when comparing recent dates.
	RecentDuration = DefaultRecentDuration

	// DumpTimeFormat is configurable format for dumping [time.Time] values.
	DumpTimeFormat = DefaultDumpTimeFormat

	// DumpDepth is a configurable depth when dumping values in log messages.
	DumpDepth = DefaultDumpDepth
)

// typeCheckers is the global map of custom checkers for given types.
var typeCheckers map[reflect.Type]Check

// RegisterTypeChecker globally registers a custom checker for a given type.
// It panics if a checker for the same type is already registered.
func RegisterTypeChecker(typ any, chk Check) {
	if chk == nil {
		panic("cannot register a nil checker")
	}
	if typeCheckers == nil {
		typeCheckers = make(map[reflect.Type]Check)
	}
	rt := reflect.TypeOf(typ)
	msg := fmt.Sprintf("Registering type checker for: %s", rt)
	if _, ok := typeCheckers[rt]; ok {
		panic("cannot overwrite an existing checker: " + msg)
	}
	_ = globLog.Output(2, msg)
	typeCheckers[rt] = chk
}

// Check is signature for generic check function comparing two arguments
// returning error if they are not. The returned error might be one or more
// errors joined with [errors.Join].
type Check func(want, have any, opts ...Option) error

// Option represents [Check] option.
type Option func(Options) Options

// WithTrail is a [Check] option setting initial field/element/key breadcrumb
// trail.
func WithTrail(pth string) Option {
	return func(ops Options) Options {
		ops.Trail = pth
		return ops
	}
}

// WithTrailLog is a [Check] option turning on a collection of checked
// fields/elements/keys. The trails are added to the provided slice.
func WithTrailLog(list *[]string) Option {
	return func(ops Options) Options {
		ops.TrailLog = list
		return ops
	}
}

// WithTimeFormat is a [Check] option setting time format when parsing dates.
func WithTimeFormat(format string) Option {
	return func(ops Options) Options {
		ops.TimeFormat = format
		return ops
	}
}

// WithRecent is a [Check] option setting duration used to compare recent dates.
func WithRecent(recent time.Duration) Option {
	return func(ops Options) Options {
		ops.Recent = recent
		return ops
	}
}

// WithDumper is [Check] option setting [dump.Config] options.
func WithDumper(optsD ...dump.Option) Option {
	return func(optsC Options) Options {
		for _, opt := range optsD {
			opt(&optsC.Dumper)
		}
		return optsC
	}
}

// WithTypeChecker is a [Check] option setting custom checker for a type.
func WithTypeChecker(typ any, chk Check) Option {
	return func(ops Options) Options {
		if ops.TypeCheckers == nil {
			ops.TypeCheckers = make(map[reflect.Type]Check)
		}
		rt := reflect.TypeOf(typ)
		if _, ok := typeCheckers[rt]; ok {
			format := "Overwriting the global type checker for: %s"
			_ = globLog.Output(2, fmt.Sprintf(format, rt))
		}
		ops.TypeCheckers[rt] = chk
		return ops
	}
}

// WithTrailChecker is a [Check] option setting a custom checker for a given
// trail.
func WithTrailChecker(trail string, chk Check) Option {
	return func(ops Options) Options {
		if ops.TrailCheckers == nil {
			ops.TrailCheckers = make(map[string]Check)
		}
		ops.TrailCheckers[trail] = chk
		return ops
	}
}

// WithSkipTrail is a [Check] option setting trails to skip.
func WithSkipTrail(skip ...string) Option {
	return func(ops Options) Options {
		ops.SkipTrails = append(ops.SkipTrails, skip...)
		return ops
	}
}

// WithSkipUnexported is a [Check] option instructing equality checks to skip
// exported fields.
func WithSkipUnexported(ops Options) Options {
	ops.SkipUnexported = true
	return ops
}

// WithOptions is a [Check] option which passes all options.
func WithOptions(src Options) Option {
	return func(ops Options) Options {
		ops.Dumper = src.Dumper
		ops.TimeFormat = src.TimeFormat
		ops.Recent = src.Recent
		ops.Trail = src.Trail
		ops.TrailLog = src.TrailLog
		ops.TypeCheckers = src.TypeCheckers
		ops.TrailCheckers = src.TrailCheckers
		ops.SkipTrails = src.SkipTrails
		ops.SkipUnexported = src.SkipUnexported
		ops.now = src.now
		return ops
	}
}

// Options represent options used by [Check] functions.
type Options struct {
	// Dump configuration.
	Dumper dump.Dump

	// Time format when parsing time strings (default: [time.RFC3339]).
	TimeFormat string

	// Duration when comparing recent dates.
	Recent time.Duration

	// Field/element/key breadcrumb trail being checked.
	Trail string

	// List of visited trails.
	// The skipped trails have " <skipped>" suffix.
	TrailLog *[]string

	// Custom checks to run for a given type.
	TypeCheckers map[reflect.Type]Check

	// Custom checker for given trail.
	TrailCheckers map[string]Check

	// List of trails to skip.
	SkipTrails []string

	// Skips all unexported fields during equality checks.
	SkipUnexported bool

	// Function used to get current time. Used preliminary to inject a clock in
	// tests of checks and assertions using [time.Now].
	now func() time.Time
}

// DefaultOptions returns default [Options].
func DefaultOptions(opts ...Option) Options {
	ops := Options{
		Dumper: dump.New(
			dump.WithTimeFormat(DumpTimeFormat),
			dump.WithMaxDepth(DumpDepth),
		),
		Recent:       RecentDuration,
		TimeFormat:   ParseTimeFormat,
		TypeCheckers: maps.Clone(typeCheckers),
		now:          time.Now,
	}
	return ops.set(opts)
}

// set sets [Options] from a slice of [Option] functions.
func (ops Options) set(opts []Option) Options {
	dst := ops
	for _, opt := range opts {
		dst = opt(dst)
	}
	return dst
}

// logTrail logs non-empty [Options.Trail] to [Options.TrailLog].
func (ops Options) logTrail() Options {
	if ops.TrailLog != nil && ops.Trail != "" {
		*ops.TrailLog = append(*ops.TrailLog, ops.Trail)
	}
	return ops
}

// structTrail updates [Options.Trail] with struct type and/or field name
// considering an already existing trail.
//
// Example trails:
//
//	Type.Field
//	Type.Field.Field
//	Type.Field[1].Field
//	Type.Field["A"].Field
func (ops Options) structTrail(typeName, fldName string) string {
	left := ops.Trail
	if typeName != "" && ops.Trail == "" {
		left = typeName
	}
	if left != "" && fldName != "" {
		return left + "." + fldName
	}
	if left == "" && fldName != "" {
		return fldName
	}
	return left
}

// mapTrail updates [Options.Trail] with trail of the map value considering
// already existing trails.
//
// Example trails:
//
//	map[1]
//	["A"]map[1]
//	[1]map["A"]
//	field["A"]
func (ops Options) mapTrail(key string) string {
	next := ops.Trail
	if ops.Trail == "" {
		next = "map"
	}
	if next[len(next)-1] == ']' {
		next += "map"
	}
	next += "[" + key + "]"
	return next
}

// arrTrail updates [Options.Trail] with slice or array index considering
// already existing trail.
//
// Example trails:
//
//	arr[1]
//	[1]
func (ops Options) arrTrail(kind string, idx int) string {
	next := ops.Trail
	if next == "" && kind != "" {
		next = "<" + kind + ">"
	}
	next += "[" + strconv.Itoa(idx) + "]"
	return next
}
