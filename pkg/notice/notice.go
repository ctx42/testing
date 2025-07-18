// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package notice simplifies building structured assertion messages.
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

// ErrNotice is a sentinel error automatically wrapped by all instances of
// [Notice] unless changed with the [Notice.Wrap] method.
var ErrNotice = errors.New("notice error")

// Notice represents a structured expectation violation message consisting of a
// header, trail and multiple named rows with context.
//
// nolint: errname
type Notice struct {
	Header string // Header message.

	// Is a trail to the field, element or key the notice message is about.
	Trail string

	Rows []Row          // Context rows.
	Meta map[string]any // Useful metadata.
	err  error          // Base error (default: [ErrNotice]).
	prev *Notice        // Next message in the chain.
	next *Notice        // Previous message in the chain.
}

// New creates a new [Notice] with a header formatted using [fmt.Sprintf] from
// the provided format string and optional arguments. The resulting [Notice]
// has its base error set to [ErrNotice]. The header is set by calling. If no
// arguments are provided, the format string is used as-is.
//
// Example:
//
//	n := New("header: %s", "name") // Header: "header: name"
//	n := New("generic error")      // Header: "generic error"
func New(header string, args ...any) *Notice {
	msg := &Notice{err: ErrNotice}
	return msg.SetHeader(header, args...)
}

// From returns instance of [Notice] if it is in err's tree. If the prefix is
// not empty, the header will be prefixed with the first element in the slice.
// If "err" is not an instance of [Notice], it will create a new one and wrap
// the "err".
func From(err error, prefix ...string) *Notice {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Notice); ok { // nolint: errorlint
		if len(prefix) > 0 {
			e.Header = fmt.Sprintf("[%s] %s", prefix[0], e.Header)
		}
		return e
	}

	header := "assertion error"
	if len(prefix) > 0 {
		header = fmt.Sprintf("[%s] %s", prefix[0], header)
	}
	return New(header).Wrap(err)
}

// SetHeader sets the header message. Implements fluent interface.
func (msg *Notice) SetHeader(header string, args ...any) *Notice {
	if len(args) > 0 {
		header = fmt.Sprintf(header, args...)
	}
	msg.Header = header
	return msg
}

// Append appends a new row with the specified name and value build using
// [fmt.Sprintf] from format and args. Implements fluent interface.
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

// Want uses the Append method to append a row with the "want" name. If the
// "want" row already exists, it will just replace its value.
func (msg *Notice) Want(format string, args ...any) *Notice {
	return msg.Append("want", format, args...)
}

// Have uses the Append method to append a row with the "have" name. If the
// "have" row already exists, it will just replace its value.
func (msg *Notice) Have(format string, args ...any) *Notice {
	return msg.Append("have", format, args...)
}

// Wrap sets base error with provided one.
func (msg *Notice) Wrap(err error) *Notice {
	msg.err = err
	return msg
}

// Unwrap returns wrapped error. By default, it returns [ErrNotice] unless a
// different error was specified using [Notice.Wrap].
func (msg *Notice) Unwrap() error {
	return msg.err
}

// Remove removes named row.
func (msg *Notice) Remove(name string) *Notice {
	fn := func(row Row) bool { return row.Name == name }
	msg.Rows = slices.DeleteFunc(msg.Rows, fn)
	return msg
}

func (msg *Notice) Is(target error) bool { return errors.Is(msg.err, target) }

// Notice returns a formatted string representation of the Notice.
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

		if multiMsg && m.Header != "" {
			buf.WriteString("  ")
			buf.WriteString(Pad("error", longest))
			buf.WriteString(": ")
			buf.WriteString(m.Header)
		} else {
			buf.WriteString(m.Header)
		}

		if len(rows) > 0 && m.Header != "" {
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

// MetaSet sets data. To get it back, use the [Notice.MetaLookup] method.
func (msg *Notice) MetaSet(key string, val any) *Notice {
	if msg.Meta == nil {
		msg.Meta = make(map[string]any)
	}
	msg.Meta[key] = val
	return msg
}

// MetaLookup returns the data set by [Notice.MetaSet]. Returns nil and false
// if the key was never set.
func (msg *Notice) MetaLookup(key string) (any, bool) {
	if msg.Meta == nil {
		return nil, false
	}
	val, ok := msg.Meta[key]
	return val, ok
}

// The longest returns the length of the longest row name among all [Notice.Rows]
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

// Chain adds the current [Notice] as next in the chain after "prev" and
// returns the current instance.
func (msg *Notice) Chain(prev *Notice) *Notice {
	msg.prev = prev
	prev.next = msg
	return msg
}

// Head returns the head of the notice chain, if the current notice instance is
// the head, it returns self.
func (msg *Notice) Head() *Notice {
	if msg.prev == nil {
		return msg
	}
	return msg.prev.Head()
}

// Next returns the next [Notice] in the chain or nil.
func (msg *Notice) Next() *Notice { return msg.next }

// Prev returns the previous [Notice] in the chain or nil.
func (msg *Notice) Prev() *Notice { return msg.prev }

// collect collects all the notices in the chain starting with the Head notice.
func (msg *Notice) collect() []*Notice {
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

// Join joins multiple notices into one. Returns the last not nil joined notice.
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
