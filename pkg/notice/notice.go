// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package notice provides the builder for rich, structured assertion messages
// used throughout the module.
//
// It is the formatting layer in the layered architecture: `assert` (user-facing)
// is built on `check` (composable, returns error) which produces `*Notice`
// values for detailed, trail-aware diagnostics.
//
// Notices implement [error]. They wrap a base error (default [ErrNotice])
// and support [errors.Is] / [errors.As] via the wrapped error. Use
// [Notice.Wrap] to change the base.
//
// Notices can be chained with [Join] (or [Notice.Chain]) into a linked list.
// Walk the chain with [Notice.Head], [Notice.Next], [Notice.Prev], or collect
// all with [Notice.All]. Chains are mutable.
//
// See the package [README] for usage and [examples_test.go] for executable
// examples.
//
// Key types and entry points:
//   - [Notice] — core structured message (header + trail + rows + metadata)
//   - [New] / [From] — create notices
//   - [Join] — chain multiple notices or errors
//   - [Notice.Head] / [Notice.Next] / [Notice.Prev] / [Notice.All] —
//     walk chains
//   - [Row] / [NewRow] — context lines
package notice

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	// trail represents the header name containing the field trail (path) name.
	trail = "trail"

	// multiHeader is a header for multiple notices.
	multiHeader = "multiple expectations violated"
)

// ErrNotice is the default base error wrapped by every [Notice].
// It can be changed per instance with [Notice.Wrap].
//
// All notices implement [error] and delegate [errors.Is] / [errors.As] to
// their base error (see [Notice.Is] and [Notice.Unwrap]).
var ErrNotice = errors.New("notice error")

// Notice represents a structured expectation violation message.
//
// It holds a header, optional trail (e.g. "User[0].Name"), context rows,
// and metadata. Notices implement [error] and form mutable chains via
// [Join] / [Notice.Chain]. Walk with [Notice.Head], [Notice.Next],
// [Notice.Prev] or collect with [Notice.All].
//
// The default base error is [ErrNotice]; change it with [Notice.Wrap]
// to control what [errors.Is] and [errors.As] see (see [Notice.Is]).
//
// nolint: errname
type Notice struct {
	HeaderText   string // HeaderText message.
	HeaderPrefix string // HeaderText prefix message.

	// Is a trail to the field, element, or key the notice message is about.
	Trail string

	Rows []Row          // Context rows.
	Meta map[string]any // Useful metadata.
	err  error          // Base error (default: [ErrNotice]).
	prev *Notice        // Next message in the chain.
	next *Notice        // Previous message in the chain.
}

// New creates a [Notice] with the given header (formatted with [fmt.Sprintf]
// if args are provided). The notice starts with base error [ErrNotice].
//
// Example:
//
//	n := New("expected %s to be equal", "values")
//	n := New("generic error")
func New(header string, args ...any) *Notice {
	msg := &Notice{err: ErrNotice}
	return msg.SetHeader(header, args...)
}

// From extracts a [Notice] from err's error tree (or creates one wrapping err).
// If a prefix is provided, it is prepended to the header in "[prefix] ..." form.
//
// If err is already a *Notice, it is returned (after optional prefixing).
// Otherwise a new notice with header "assertion error" (or prefixed) is
// created and the original err is wrapped via [Notice.Wrap].
func From(err error, prefix ...string) *Notice {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Notice); ok { // nolint: errorlint
		if len(prefix) > 0 && prefix[0] != "" {
			e.HeaderPrefix = prefix[0]
		}
		return e
	}

	header := "assertion error"
	pref := ""
	if len(prefix) > 0 && prefix[0] != "" {
		pref = prefix[0]
	}
	msg := New(header).Wrap(err)
	msg.HeaderPrefix = pref
	return msg
}

// SetHeader sets the header message. Implements fluent interface.
func (msg *Notice) SetHeader(header string, args ...any) *Notice {
	if len(args) > 0 {
		header = fmt.Sprintf(header, args...)
	}
	msg.HeaderText = header
	return msg
}

// SetPrefix sets the header prefix. Implements fluent interface.
func (msg *Notice) SetPrefix(prefix string, args ...any) *Notice {
	if len(args) > 0 {
		prefix = fmt.Sprintf(prefix, args...)
	}
	msg.HeaderPrefix = prefix
	return msg
}

// Header returns the formatted Header string. When HeaderPrefix is set, it
// returns "[Prefix] HeaderText"; otherwise it returns HeaderText unchanged.
func (msg *Notice) Header() string {
	if msg.HeaderPrefix != "" {
		return fmt.Sprintf("[%s] %s", msg.HeaderPrefix, msg.HeaderText)
	}
	return msg.HeaderText
}

// Append adds a row (name + formatted value). Replaces any existing row with
// the same name. Implements fluent interface.
func (msg *Notice) Append(name, format string, args ...any) *Notice {
	return msg.appendRow(NewRow(name, format, args...))
}

// AppendRow appends description rows to the message.
func (msg *Notice) AppendRow(desc ...Row) *Notice {
	for _, row := range desc {
		_ = msg.appendRow(row)
	}
	return msg
}

// appendRow appends [Row] to the [Notice]. Implements fluent interface.
func (msg *Notice) appendRow(row Row) *Notice {
	fn := func(have Row) bool { return have.Name == row.Name }
	if idx := slices.IndexFunc(msg.Rows, fn); idx >= 0 {
		msg.Rows[idx].Format = row.Format
		msg.Rows[idx].Args = row.Args
		return msg
	}
	msg.Rows = append(msg.Rows, row)
	return msg
}

// Prepend prepends a new row with the specified name and value built using
// [fmt.Sprintf] from format and args. Implements fluent interface.
func (msg *Notice) Prepend(name, format string, args ...any) *Notice {
	return msg.prependRow(NewRow(name, format, args...))
}

// prependRow prepends [Row] to the [Notice]. Implements fluent interface.
func (msg *Notice) prependRow(row Row) *Notice {
	fn := func(have Row) bool { return have.Name == row.Name }
	if idx := slices.IndexFunc(msg.Rows, fn); idx >= 0 {
		msg.Rows[idx].Format = row.Format
		msg.Rows[idx].Args = row.Args
		return msg
	}
	msg.Rows = slices.Insert(msg.Rows, 0, row)
	return msg
}

// SetTrail adds trail row if "tr" is not an empty string. If the trail row
// already exists, it overwrites it. Implements fluent interface.
//
// SetTrail examples:
//
//   - Type
//   - Type[1].Field
//   - Type["key"].Field
func (msg *Notice) SetTrail(trail string) *Notice {
	msg.Trail = trail
	return msg
}

// Want is a convenience for Append("want", ...). Replaces any prior "want" row.
func (msg *Notice) Want(format string, args ...any) *Notice {
	return msg.Append("want", format, args...)
}

// Have is a convenience for Append("have", ...). Replaces any prior "have" row.
func (msg *Notice) Have(format string, args ...any) *Notice {
	return msg.Append("have", format, args...)
}

// Wrap sets the base error returned by [Notice.Unwrap] and used by
// [Notice.Is] for [errors.Is] / [errors.As] delegation.
// Defaults to [ErrNotice] if never called.
func (msg *Notice) Wrap(err error) *Notice {
	msg.err = err
	return msg
}

// Unwrap returns the base error (for [errors.Is] / [errors.As] support).
// Defaults to [ErrNotice] unless changed with [Notice.Wrap].
func (msg *Notice) Unwrap() error {
	return msg.err
}

// Remove removes the row with the given name, if present.
// Implements fluent interface.
func (msg *Notice) Remove(name string) *Notice {
	fn := func(row Row) bool { return row.Name == name }
	msg.Rows = slices.DeleteFunc(msg.Rows, fn)
	return msg
}

// Is reports whether the base error matches target via [errors.Is].
// This makes every *Notice satisfy errors.Is(msg, ErrNotice) by default.
func (msg *Notice) Is(target error) bool { return errors.Is(msg.err, target) }

// Error returns a formatted string representation of the Notice.
//
// nolint: gocognit, cyclop
func (msg *Notice) Error() string {
	mgs := msg.collect()

	var longest int
	for _, m := range mgs {
		if ln := m.longest(); longest < ln {
			longest = ln
		}
	}

	buf := &strings.Builder{}
	multiMsg := len(mgs) > 1

	if multiMsg {
		if longest < len("error") {
			longest = len("error")
		}
		buf.WriteString(multiHeader)
		buf.WriteString(":\n")
	}

	for im, m := range mgs {
		lastMsg := im == len(mgs)-1

		rows := m.Rows
		if m.Trail != "" {
			rows = append([]Row{NewRow(trail, "%s", m.Trail)}, m.Rows...)
		}

		header := m.Header()
		if multiMsg && header != "" {
			buf.WriteString("  ")
			buf.WriteString(Pad("error", longest))
			buf.WriteString(": ")
			buf.WriteString(header)
		} else {
			buf.WriteString(header)
		}

		if len(rows) > 0 && header != "" {
			if !multiMsg {
				buf.WriteString(":")
			}
			buf.WriteString("\n")
		}

		for ir, r := range rows {
			lastRow := ir == len(rows)-1
			name := r.PadName(longest)
			value := r.String()

			buf.WriteString("  ")
			buf.WriteString(name)
			buf.WriteString(":")

			if idx := strings.IndexByte(value, '\n'); idx >= 0 {
				value = Indent(len(name)+4, ' ', value)
				if idx != 0 {
					buf.WriteString("\n")
				}
			} else {
				buf.WriteString(" ")
			}
			buf.WriteString(value)

			if !lastRow {
				buf.WriteString("\n")
			}
		}

		if !lastMsg {
			buf.WriteString("\n")
			buf.WriteString(Pad("---", longest+4))
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// MetaSet stores arbitrary data under key for later retrieval.
// To retrieve, use [Notice.MetaLookup]. Implements fluent interface.
func (msg *Notice) MetaSet(key string, val any) *Notice {
	if msg.Meta == nil {
		msg.Meta = make(map[string]any)
	}
	msg.Meta[key] = val
	return msg
}

// MetaLookup returns the value stored by [Notice.MetaSet] for key, or
// nil and false if the key was never set or Meta was nil.
func (msg *Notice) MetaLookup(key string) (any, bool) {
	if msg.Meta == nil {
		return nil, false
	}
	val, ok := msg.Meta[key]
	return val, ok
}

// longest returns the length of the longest row name among all [Notice.Rows]
// and the [Notice.Trail] string. If there are no rows and the trail is empty,
// it returns 0.
func (msg *Notice) longest() int {
	var maxLen int
	if msg.Trail != "" {
		maxLen = len(trail)
	}
	for _, row := range msg.Rows {
		if maxLen < len(row.Name) {
			maxLen = len(row.Name)
		}
	}
	return maxLen
}

// Chain links msg after prev in the chain (mutates both) and returns msg.
//
// After the call: prev.Next() == msg and msg.Prev() == prev.
// Use [Join] for a safer way to build chains from multiple values.
func (msg *Notice) Chain(prev *Notice) *Notice {
	msg.prev = prev
	prev.next = msg
	return msg
}

// Head returns the first notice in the chain (or self if it is the head).
//
// Follows the prev pointers to the start of the linked list.
func (msg *Notice) Head() *Notice {
	if msg.prev == nil {
		return msg
	}
	return msg.prev.Head()
}

// Next returns the following notice in the chain (nil if this is the last).
//
// The link is set by [Chain] or [Join].
func (msg *Notice) Next() *Notice { return msg.next }

// Prev returns the preceding notice in the chain (nil if this is the first).
//
// The link is set by [Chain] or [Join].
func (msg *Notice) Prev() *Notice { return msg.prev }

// All returns every notice in the chain, starting from the head.
// The slice is newly allocated on each call.
//
// See [Join] for building chains and [Notice.Head] for finding the start.
func (msg *Notice) All() []*Notice {
	head := msg.Head()
	var mgs []*Notice
	for {
		mgs = append(mgs, head)
		if head.next == nil {
			break
		}
		head = head.next
	}
	return mgs
}

// collect collects all the notices in the chain starting with the Head notice.
func (msg *Notice) collect() []*Notice { return msg.All() }

// Join chains multiple notices or errors into a linked list.
//
// Each non-nil argument is converted via [From] and linked using
// [Notice.Chain]. The result can be walked with [Notice.Head],
// [Notice.Next], [Notice.Prev] or collected with [Notice.All].
//
// Returns the last non-nil notice in the chain (or nil). The returned
// value implements [error] and supports [errors.Is]/[errors.As] through
// the notices' wrapped errors.
func Join(ers ...error) error {
	if len(ers) == 0 {
		return nil
	}

	var err *Notice
	for _, next := range ers {
		if next == nil {
			continue
		}
		ne := From(next)
		if err == nil {
			err = ne
			continue
		}
		err = ne.Chain(err)
	}
	if err == nil {
		return nil
	}
	return err
}
