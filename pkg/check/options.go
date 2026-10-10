// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
)

// Package-wide default configuration.
const (
	// DefaultParseTimeFormat is the default format for parsing time strings.
	DefaultParseTimeFormat = time.RFC3339Nano

	// DefaultRecentDuration is the default duration when comparing recent
	// dates.
	DefaultRecentDuration = 10 * time.Second

	// DefaultDumpTimeFormat is the default format for dumping [time.Time]
	// values.
	DefaultDumpTimeFormat = time.RFC3339Nano

	// DefaultDumpDepth is the default depth when dumping values recursively in
	// log messages.
	DefaultDumpDepth = 6
)

// TimeFormatUnixStr is a special [WithTimeFormat] sentinel indicating that a
// string time value represents a UNIX timestamp (seconds since epoch).
const TimeFormatUnixStr = "str-ts"

// Package-wide configuration.
var (
	// ParseTimeFormat is a configurable format for parsing time strings.
	ParseTimeFormat = DefaultParseTimeFormat

	// RecentDuration is a configurable duration when comparing recent dates.
	RecentDuration = DefaultRecentDuration

	// DumpTimeFormat is the format for dumping [time.Time] values.
	DumpTimeFormat = DefaultDumpTimeFormat

	// DumpDepth is a configurable depth when dumping values in log messages.
	DumpDepth = DefaultDumpDepth
)

// Checker is the signature for a generic check function. It compares two
// values and returns an error (typically a *[notice.Notice] or joined errors)
// if they do not match.
//
// Custom checkers can be registered globally with [RegisterTypeChecker] or
// per-call with [WithTypeChecker]. See the Design section in the root README.
//
// Options can be [Option] values, strings (for comments), or other supported
// forms. See [DefaultOptions] and the With* functions.
//
// Example:
//
//	Checker("want", "have", "A%d", 42, WithTrail("type.field"))
type Checker func(want, have any, opts ...any) error

// typeCheckers holds globally registered custom checkers (see
// [RegisterTypeChecker]), guarded by typeCheckersMx.
var (
	typeCheckers   map[reflect.Type]Checker
	typeCheckersMx sync.RWMutex
)

// RegisterTypeChecker globally registers a custom [Checker] for values of the
// given type. It will be used by default in all checks (via [DefaultOptions])
// unless overridden per-call with [WithTypeChecker].
//
// It panics if typ is nil or a checker for the same type is already
// registered. It is safe to call concurrently with checks.
//
// See the Design section in the root README for the overall customization
// model, and [WithTypeChecker] for per-check overrides.
func RegisterTypeChecker(typ any, chk Checker) {
	if chk == nil {
		panic("cannot register a nil type checker")
	}
	if typ == nil {
		panic("cannot register a type checker for a nil type")
	}
	typeCheckersMx.Lock()
	defer typeCheckersMx.Unlock()
	if typeCheckers == nil {
		typeCheckers = make(map[reflect.Type]Checker)
	}
	rt := reflect.TypeOf(typ)
	if _, ok := typeCheckers[rt]; ok {
		format := "cannot overwrite an existing type checker for: %s"
		panic(fmt.Sprintf(format, rt))
	}
	typeCheckers[rt] = chk
}

// Option configures a check (see [DefaultOptions]).
type Option func(Options) Options

// WithTrail sets the initial trail (path to the field/element/key being
// checked). Trails are used in error messages and for trail-specific custom
// checkers ([WithTrailChecker]).
func WithTrail(pth string) Option {
	return func(ops Options) Options {
		ops.Trail = pth
		return ops
	}
}

// WithTrailLog is a [Checker] option turning on a collection of checked
// fields/elements/keys. The trails are added to the provided slice.
func WithTrailLog(list *[]string) Option {
	return func(ops Options) Options {
		ops.TrailLog = list
		return ops
	}
}

// WithTimeFormat is a [Checker] option setting the time format when parsing
// dates.
func WithTimeFormat(format string) Option {
	return func(ops Options) Options {
		ops.TimeFormat = format
		return ops
	}
}

// WithZone is a [Checker] option that specifies the timezone to apply to the
// "want" date before comparing times. It ensures consistent timezone handling
// for string-based date inputs. For [time.Time] values, it calls
// [time.Time.In] to adjust the timezone.
//
// Example:
//
//	assert.Exact(t, "2000-01-02T02:04:05Z", have, check.WithZone(WAW))
//
// In the above example the string date will be parsed then [time.Time.In]
// method will set its timezone.
func WithZone(zone *time.Location) Option {
	return func(ops Options) Options {
		ops.Zone = zone
		return ops
	}
}

// WithRecent is a [Checker] option setting the duration used to compare recent
// dates.
func WithRecent(recent time.Duration) Option {
	return func(ops Options) Options {
		ops.Recent = recent
		return ops
	}
}

// WithDumper is a [Checker] option applying [dump.Option] options.
func WithDumper(optsD ...dump.Option) Option {
	return func(ops Options) Options {
		ops.Dumper.Dumpers = maps.Clone(ops.Dumper.Dumpers)
		for _, opt := range optsD {
			opt(&ops.Dumper)
		}
		return ops
	}
}

// WithTypeChecker registers a custom [Checker] for values of the given type
// (used only for checks with these options). It overrides any globally
// registered checker from [RegisterTypeChecker] for that type. It panics if
// typ is nil.
//
// See the Design section in the root README for the customization model.
func WithTypeChecker(typ any, chk Checker) Option {
	if typ == nil {
		panic("cannot set a type checker for a nil type")
	}
	rt := reflect.TypeOf(typ)
	return func(ops Options) Options {
		ops.TypeCheckers = maps.Clone(ops.TypeCheckers)
		if ops.TypeCheckers == nil {
			ops.TypeCheckers = make(map[reflect.Type]Checker)
		}
		ops.TypeCheckers[rt] = chk
		return ops
	}
}

// WithTrailChecker registers a custom [Checker] that will only be used when
// the current trail exactly matches the given string.
func WithTrailChecker(trail string, chk Checker) Option {
	return func(ops Options) Options {
		ops.TrailCheckers = maps.Clone(ops.TrailCheckers)
		if ops.TrailCheckers == nil {
			ops.TrailCheckers = make(map[string]Checker)
		}
		ops.TrailCheckers[trail] = chk
		return ops
	}
}

// WithSkipTrail is a [Checker] option setting trails to skip.
func WithSkipTrail(skip ...string) Option {
	return func(ops Options) Options {
		ops.SkipTrails = append(slices.Clip(ops.SkipTrails), skip...)
		return ops
	}
}

// WithSkipUnexported is a [Checker] option instructing equality checks to skip
// unexported fields.
func WithSkipUnexported() Option {
	return func(ops Options) Options {
		ops.SkipUnexported = true
		return ops
	}
}

// WithIncreasingSoft is an option used by the [Increasing] check allowing
// consecutive values to be equal to each other.
func WithIncreasingSoft() Option {
	return func(ops Options) Options {
		ops.IncreaseSoft = true
		return ops
	}
}

// WithDecreasingSoft is an option used by the [Decreasing] check allowing
// consecutive values to be equal to each other.
func WithDecreasingSoft() Option {
	return func(ops Options) Options {
		ops.DecreaseSoft = true
		return ops
	}
}

// WithCmpBaseTypes is a [Checker] option turning on simple base type
// comparisons.
//
// During a normal operation, when comparing values with different types, the
// error is returned. When this option is used, and both values have the same
// underlying simple type the values of that base type will be compared.
//
// Simple types are defined as:
//   - [reflect.Bool]
//   - [reflect.Int]
//   - [reflect.Int8]
//   - [reflect.Int16]
//   - [reflect.Int32]
//   - [reflect.Int64]
//   - [reflect.Uint]
//   - [reflect.Uint8]
//   - [reflect.Uint16]
//   - [reflect.Uint32]
//   - [reflect.Uint64]
//   - [reflect.Float32]
//   - [reflect.Float64]
//   - [reflect.Complex64]
//   - [reflect.Complex128]
//   - [reflect.String]
//
// Example:
//
//	// --- Given ---
//	type MyInt int
//	const MyIntValue MyInt = 42
//
//	m0 := map[string]any{"A": MyIntValue}
//	m1 := map[string]any{"A": 42}
//
//	// --- When ---
//	err := Equal(m0, m1, WithCmpBaseTypes())
//
//	// --- Then ---
//	assert.NoError(t, err)
func WithCmpBaseTypes() Option {
	return func(ops Options) Options {
		ops.CmpSimpleType = true
		return ops
	}
}

// WithWaitThrottle is a [Checker] option setting how calls to the test
// function should be throttled.
func WithWaitThrottle(throttle time.Duration) Option {
	return func(ops Options) Options {
		ops.WaitThrottle = throttle
		return ops
	}
}

// WithComment is a [Checker] option setting a comment to be added to all error
// messages.
func WithComment(format string, args ...any) Option {
	return func(ops Options) Options {
		ops.Comment = fmt.Sprintf(format, args...)
		return ops
	}
}

// WithOptions is a [Checker] option that passes all options.
func WithOptions(src Options) Option {
	return func(ops Options) Options {
		ops.Dumper = src.Dumper
		ops.TimeFormat = src.TimeFormat
		ops.Zone = src.Zone
		ops.Recent = src.Recent
		ops.Trail = src.Trail
		ops.TrailLog = src.TrailLog
		ops.TypeCheckers = src.TypeCheckers
		ops.TrailCheckers = src.TrailCheckers
		ops.SkipTrails = src.SkipTrails
		ops.SkipUnexported = src.SkipUnexported
		ops.CmpSimpleType = src.CmpSimpleType
		ops.IncreaseSoft = src.IncreaseSoft
		ops.DecreaseSoft = src.DecreaseSoft
		ops.WaitThrottle = src.WaitThrottle
		ops.Comment = src.Comment
		ops.now = src.now
		return ops
	}
}

// Options represents options used by [Checker] functions.
type Options struct {
	// Dump configuration.
	Dumper dump.Dump

	// Time format when parsing time strings (default: [time.RFC3339Nano]).
	TimeFormat string

	// See [WithZone].
	Zone *time.Location

	// Duration when comparing recent dates.
	Recent time.Duration

	// Trail is the field/element/key breadcrumb trail being checked.
	Trail string

	// List of visited trails.
	// The skipped trails have a "<skipped>" suffix.
	TrailLog *[]string

	// Custom checks to run for a given type.
	TypeCheckers map[reflect.Type]Checker

	// Custom checker for given trail.
	TrailCheckers map[string]Checker

	// List of trails to skip.
	SkipTrails []string

	// Skips all unexported fields during equality checks.
	SkipUnexported bool

	// See [WithCmpBaseTypes].
	CmpSimpleType bool

	// Option for [Increasing] allowing consecutive values to be equal.
	IncreaseSoft bool

	// Option for [Decreasing] allowing consecutive values to be equal.
	DecreaseSoft bool

	// Option for [Wait] throttling the calls to a test function.
	WaitThrottle time.Duration

	// Comment row for all error messages.
	Comment string

	// Function used to get the current time. Used mainly to inject a clock in
	// tests of checks and assertions using [time.Now].
	now func() time.Time
}

// DefaultOptions builds an [Options] struct from the provided arguments
// (mix of [Option], strings for comments, etc.). It is the central entry
// point for customizing checks (trail, custom checkers via global
// [RegisterTypeChecker] or [WithTypeChecker], dumper behavior, etc.).
//
// See the With* functions and the root README Design section for the
// overall architecture.
func DefaultOptions(opts ...any) Options {
	ops := Options{
		Dumper: dump.New(
			dump.WithTimeFormat(DumpTimeFormat),
			dump.WithMaxDepth(DumpDepth),
		),
		Recent:       RecentDuration,
		TimeFormat:   ParseTimeFormat,
		Zone:         nil,
		TypeCheckers: cloneTypeCheckers(),
		WaitThrottle: 10 * time.Millisecond,
		now:          time.Now,
	}

	cmtFmt, cmtArgs, chkOpts := parseCheckArgs(opts...)
	if cmtFmt != "" {
		chkOpts = append(chkOpts, WithComment(cmtFmt, cmtArgs...))
	}

	ops = ops.set(chkOpts)
	defaults := map[reflect.Type]Checker{
		typTime:    Time,
		typZonePtr: Zone,
		typZone:    Zone,
	}
	var cloned bool
	for typ, chk := range defaults {
		if _, ok := ops.TypeCheckers[typ]; ok {
			continue
		}
		if !cloned {
			// Options passed by WithOptions share the map with the caller.
			ops.TypeCheckers = maps.Clone(ops.TypeCheckers)
			if ops.TypeCheckers == nil {
				ops.TypeCheckers = make(map[reflect.Type]Checker, 3)
			}
			cloned = true
		}
		ops.TypeCheckers[typ] = chk
	}
	return ops
}

// cloneTypeCheckers returns a copy of the globally registered checkers.
func cloneTypeCheckers() map[reflect.Type]Checker {
	typeCheckersMx.RLock()
	defer typeCheckersMx.RUnlock()
	return maps.Clone(typeCheckers)
}

// parseCheckArgs separates the variadic arguments for [DefaultOptions] into
// a comment format plus its arguments, and a slice of [Option] values.
//
// The first string value is treated as the comment format. All subsequent
// values (strings or otherwise) become format arguments. Any value whose
// type is [Option] or func(Options) Options is collected as an option to
// apply instead.
//
// It panics if a non-string appears before the first string is seen.
func parseCheckArgs(opts ...any) (string, []any, []Option) {
	var cmtFmt string
	var cmtArgs []any
	var chkOpts []Option

	for _, arg := range opts {
		switch a := arg.(type) {
		case func(Options) Options:
			chkOpts = append(chkOpts, a)
		case Option:
			chkOpts = append(chkOpts, a)

		case string:
			if cmtFmt == "" {
				cmtFmt = a
			} else {
				cmtArgs = append(cmtArgs, a)
			}

		default:
			if cmtFmt == "" {
				panic("cannot use a non-string comment format")
			}
			cmtArgs = append(cmtArgs, a)
		}
	}

	return cmtFmt, cmtArgs, chkOpts
}

// set sets [Options] from a slice of [Option] functions.
func (ops Options) set(opts []Option) Options {
	dst := ops
	for _, opt := range opts {
		dst = opt(dst)
	}
	return dst
}

// LogTrail logs non-empty [Options.Trail] to [Options.TrailLog].
func (ops Options) LogTrail() Options {
	if ops.TrailLog != nil && ops.Trail != "" {
		*ops.TrailLog = append(*ops.TrailLog, ops.Trail)
	}
	return ops
}

// StructTrail updates [Options.Trail] with a struct type and/or field name
// considering an already existing trail.
//
// Example trails:
//
//	Type.Field
//	Type.Field.Field
//	Type.Field[1].Field
//	Type.Field["A"].Field
func (ops Options) StructTrail(typeName, fldName string) Options {
	left := ops.Trail
	if typeName != "" && ops.Trail == "" {
		left = typeName
	}
	if left != "" && fldName != "" {
		ops.Trail = left + "." + fldName
		return ops
	}
	if left == "" && fldName != "" {
		ops.Trail = fldName
		return ops
	}

	ops.Trail = left
	return ops
}

// MapTrail updates [Options.Trail] with the trail of the map value considering
// already existing trails.
//
// Example trails:
//
//	map[1]
//	["A"]map[1]
//	[1]map["A"]
//	field["A"]
func (ops Options) MapTrail(key string) Options {
	next := ops.Trail
	if ops.Trail == "" {
		next = "map"
	}
	if next[len(next)-1] == ']' {
		next += "map"
	}
	ops.Trail = next + "[" + key + "]"
	return ops
}

// ArrTrail updates [Options.Trail] with a slice or array index considering an
// already existing trail.
//
// Example trails:
//
//	arr[1]
//	[1]
func (ops Options) ArrTrail(kind string, idx int) Options {
	next := ops.Trail
	if next == "" && kind != "" {
		next = "<" + kind + ">"
	}
	ops.Trail = next + "[" + strconv.Itoa(idx) + "]"
	return ops
}

// FieldName returns a helper function which updates [Options.Trail].
//
// It is useful when constructing trails in custom struct checkers.
//
// Example:
//
//	func fileCheck(want, have any, opts ...any) error {
//		ops := check.DefaultOptions(opts...)
//		err := check.Type(file{}, have, check.WithOptions(ops))
//		if err != nil {
//			return err
//		}
//		w, h := want.(file), have.(file)
//
//		fName := check.FieldName(ops, "file")
//		ers := []error{
//			check.Equal(w.path, h.path, fName("path")),
//			check.Equal(w.pks, h.pks, fName("pks")),
//			// Not all fields are compared.
//			check.Fields(4, w, fName("{field count}")),
//		}
//		return notice.Join(ers...)
//	}
func FieldName(ops Options, typeName string) func(fldName string) Option {
	return func(fldName string) Option {
		return WithOptions(ops.StructTrail(typeName, fldName))
	}
}

// AddRows adds rows based on the [Options] instance to the [notice.Notice].
func AddRows(ops Options, msg *notice.Notice) *notice.Notice {
	if ops.Comment != "" {
		_ = msg.Prepend("comment", ops.Comment)
	}
	return msg.SetTrail(ops.Trail)
}
