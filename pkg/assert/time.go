// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"time"

	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Time asserts that "want" and "have" represent the same instant in time.
//
// See [check.Time] for supported representations and [check] package
// documentation for option handling.
func Time(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.Time(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Exact asserts that "want" and "have" represent the same instant and are in
// the same timezone.
//
// The arguments may be strings, integers, int64, or [time.Time]. String
// representations are parsed using [check.Options.TimeFormat]; the result is
// normalized to UTC. Integers are treated as Unix timestamps (UTC).
func Exact(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.Exact(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Before asserts that "date" is before "mark".
//
// The arguments may be strings, integers, int64, or [time.Time]. See [Exact]
// for representation and parsing details.
func Before(t tester.T, mark, date any, opts ...any) bool {
	t.Helper()
	if e := check.Before(mark, date, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// After asserts that "date" is after "mark".
//
// The arguments may be strings, integers, int64, or [time.Time]. See [Exact]
// for representation and parsing details.
func After(t tester.T, mark, date time.Time, opts ...any) bool {
	t.Helper()
	if e := check.After(mark, date, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// BeforeOrEqual asserts that "date" is before or equal to "mark".
//
// The arguments may be strings, integers, int64, or [time.Time]. See [Exact]
// for representation and parsing details.
func BeforeOrEqual(t tester.T, mark, date time.Time, opts ...any) bool {
	t.Helper()
	if e := check.BeforeOrEqual(mark, date, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// AfterOrEqual asserts that "date" is after or equal to "mark".
//
// The arguments may be strings, integers, int64, or [time.Time]. See [Exact]
// for representation and parsing details.
func AfterOrEqual(t tester.T, mark, date any, opts ...any) bool {
	t.Helper()
	if e := check.AfterOrEqual(mark, date, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Within asserts that "want" and "have" represent times within the given
// duration of each other.
//
// The arguments may be strings, integers, int64, or [time.Time]. See [Exact]
// for representation and parsing details.
func Within(t tester.T, want, within, have any, opts ...any) bool {
	t.Helper()
	if e := check.Within(want, within, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Recent asserts that "have" is within [check.Options.Recent] of [time.Now].
//
// The argument may be a string, integer, int64, or [time.Time]. See [Exact]
// for representation and parsing details.
func Recent(t tester.T, have any, opts ...any) bool {
	t.Helper()
	if e := check.Recent(have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Zone asserts that "want" and "have" timezones are equal.
func Zone(t tester.T, want, have *time.Location, opts ...any) bool {
	t.Helper()
	if e := check.Zone(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// Duration asserts that "want" and "have" durations are equal.
//
// The arguments may be strings, integers, int64, or [time.Duration].
func Duration(t tester.T, want, have any, opts ...any) bool {
	t.Helper()
	if e := check.Duration(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
