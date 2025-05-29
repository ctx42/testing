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

// trail represents a row name with special meaning representing a trail (path)
// to the field / element or key the notice message is about.
//
// Trail examples:
//
//   - Type
//   - Type[1].Field
//   - Type["key"].Field
const trail = "trail"

// ContinuationHeader is a special [Notice.Header] separating notices with the
// same header.
//
// Example:
//
//	header:
//	  want: want 0
//	  have: have 0
//	 ---
//	  want: want 1
//	  have: have 1
const ContinuationHeader = " ---"

// ErrNotice is a sentinel error automatically wrapped by all instances of
// [Notice] unless changed with the [Notice.Wrap] method.
var ErrNotice = errors.New("notice error")

// Notice represents a structured notice message consisting of a header and
// multiple named rows giving context to it.
//
// nolint: errname
type Notice struct {
	Header string         // Header message.
	Rows   []Row          // Context rows.
	Meta   map[string]any // Useful metadata.
	err    error          // Base error (default: [ErrNotice]).
}

// New creates a new [Notice] with the specified header which is constructed
// using [fmt.Sprintf] from format and args. By default, the base error is
// set to [ErrNotice].
func New(header string, args ...any) *Notice {
	msg := &Notice{err: ErrNotice}
	return msg.SetHeader(header, args...)
}

// From returns instance of [Notice] if it is in err's tree. If the prefix is
// not empty, the header will be prefixed with the first element in the slice.
// If "err" is not an instance of [Notice], it will create a new one and wrap
// the "err".
func From(err error, prefix ...string) *Notice {
	var e *Notice
	if errors.As(err, &e) {
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
	fn := func(row Row) bool { return row.Name == name }
	if idx := slices.IndexFunc(msg.Rows, fn); idx >= 0 {
		msg.Rows[idx].Format = format
		msg.Rows[idx].Args = args
		return msg
	}
	msg.Rows = append(msg.Rows, NewRow(name, format, args...))
	return msg
}

// AppendRow appends description rows to the message.
func (msg *Notice) AppendRow(desc ...Row) *Notice {
	for _, row := range desc {
		_ = msg.Append(row.Name, row.Format, row.Args...)
	}
	return msg
}

// Prepend prepends a new row with the specified name and value built using
// [fmt.Sprintf] from format and args. Implements fluent interface.
func (msg *Notice) Prepend(name, format string, args ...any) *Notice {
	fn := func(row Row) bool { return row.Name == name }
	if idx := slices.IndexFunc(msg.Rows, fn); idx >= 0 {
		msg.Rows[idx].Format = format
		msg.Rows[idx].Args = args
		return msg
	}
	var idx int
	if len(msg.Rows) != 0 && msg.Rows[0].Name == trail {
		idx = 1
	}
	msg.Rows = slices.Insert(msg.Rows, idx, NewRow(name, format, args...))
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
func (msg *Notice) SetTrail(tr string) *Notice {
	if tr == "" {
		return msg
	}
	return msg.Prepend(trail, "%s", tr)
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
func (msg *Notice) Error() string {
	m := msg.Header
	if len(msg.Rows) > 0 {
		if msg.Header != ContinuationHeader {
			m += ":"
		}
		m += "\n"
	}
	longest := msg.longestName()
	for i := range msg.Rows {
		row := msg.Rows[i]
		name := row.PadName(longest)
		value := row.String()
		format := "  %s: %s"
		if idx := strings.IndexByte(value, '\n'); idx >= 0 {
			if idx == 0 {
				format = "  %s:%s"
			} else {
				format = "  %s:\n%s"
			}
			value = Indent(len(name)+4, ' ', value)
		}
		m += fmt.Sprintf(format, name, value)
		if i < len(msg.Rows)-1 {
			m += "\n"
		}
	}
	return m
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

// longestName returns the longest row name in [Notice.Rows].
func (msg *Notice) longestName() int {
	var maxLen int
	for _, row := range msg.Rows {
		if maxLen < len(row.Name) {
			maxLen = len(row.Name)
		}
	}
	return maxLen
}
