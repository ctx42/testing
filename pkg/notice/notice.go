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

// ErrAssert represents an error that occurs during an assertion. It is
// typically used when a condition fails to meet the expected value.
var ErrAssert = errors.New("assert error")

// Notice represents structured notice message consisting of a header and
// multiple named rows giving context to it.
//
// nolint: errname
type Notice struct {
	Header string            // Header message.
	Rows   map[string]string // Context rows.
	Data   map[string]any    // Any useful data (ignored by Notice).
	Order  []string          // Order to display rows in.
	err    error             // Base error.
}

// New creates a new [Notice] with the specified header which is constructed
// using [fmt.Sprintf] from format and args. By default, the base error is
// set to [ErrAssert].
func New(header string, args ...any) *Notice {
	msg := &Notice{
		Rows: make(map[string]string, 2),
		err:  ErrAssert,
	}
	return msg.SetHeader(header, args...)
}

// From returns instance of [Notice] if it is in err's tree. If prefix is not
// empty header will be prefixed with the first element in the slice. If err is
// not instance of [Notice] it will create a new one and wrap err.
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
// [fmt.Sprintf] from format and args. If a row with the same name already
// exists, it is moved to the end of the [Notice.Order] slice. Implements
// fluent interface.
func (msg *Notice) Append(name, format string, args ...any) *Notice {
	if _, exists := msg.Rows[name]; exists {
		msg.Order = slices.DeleteFunc(msg.Order, func(s string) bool {
			return name == s
		})
	}
	msg.Rows[name] = fmt.Sprintf(format, args...)
	msg.Order = append(msg.Order, name)
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
// [fmt.Sprintf] from format and args. If a row with the same name already
// exists, it is moved to the beginning of the [Notice.Order] slice.
// Implements fluent interface.
func (msg *Notice) Prepend(name, format string, args ...any) *Notice {
	hasTrail := slices.Contains(msg.Order, trail)
	msg.Order = slices.DeleteFunc(msg.Order, func(s string) bool {
		if hasTrail && s == trail {
			return true
		}
		return name == s
	})
	msg.Rows[name] = fmt.Sprintf(format, args...)
	var prepend []string
	if hasTrail && name != trail {
		prepend = []string{trail}
	}
	prepend = append(prepend, name)
	msg.Order = append(prepend, msg.Order...)
	return msg
}

// Trail adds trail row if "tr" is not empty string. If the trail row already
// exists it overwrites it. Implements fluent interface.
//
// Trail examples:
//
//   - Type
//   - Type[1].Field
//   - Type["key"].Field
func (msg *Notice) Trail(tr string) *Notice {
	if tr == "" {
		return msg
	}
	return msg.Prepend(trail, "%s", tr)
}

// Want uses Append method to append a row with "want" name. If the "want" row
// already exists it will just replace its value.
func (msg *Notice) Want(format string, args ...any) *Notice {
	if _, ok := msg.Rows["want"]; ok {
		msg.Rows["want"] = fmt.Sprintf(format, args...)
		return msg
	}
	return msg.Append("want", format, args...)
}

// Have uses Append method to append a row with "have" name. If the "have" row
// already exists it will just replace its value.
func (msg *Notice) Have(format string, args ...any) *Notice {
	if _, ok := msg.Rows["have"]; ok {
		msg.Rows["have"] = fmt.Sprintf(format, args...)
		return msg
	}
	return msg.Append("have", format, args...)
}

// Wrap wraps base error with provided one.
func (msg *Notice) Wrap(err error) *Notice {
	msg.err = fmt.Errorf("%w: %w", msg.err, err)
	return msg
}

// Remove removes named row.
func (msg *Notice) Remove(name string) *Notice {
	msg.Order = slices.DeleteFunc(msg.Order, func(s string) bool {
		return name == s
	})
	delete(msg.Rows, name)
	return msg
}

func (msg *Notice) Is(target error) bool { return errors.Is(msg.err, target) }

// Notice returns a formatted string representation of the Notice.
func (msg *Notice) Error() string {
	m := msg.Header
	if len(msg.Order) > 0 {
		if msg.Header != ContinuationHeader {
			m += ":"
		}
		m += "\n"
	}
	names := msg.equalizeNames()
	for i := range msg.Order {
		display := names[i]
		given := msg.Order[i]
		value := msg.Rows[given]
		format := "  %s: %s"
		if strings.IndexByte(value, '\n') >= 0 {
			format = "  %s:\n%s"
			value = Indent(len(display)+4, ' ', value)
		}
		m += fmt.Sprintf(format, display, value)
		if i < len(msg.Order)-1 {
			m += "\n"
		}
	}
	return m
}

// SetData sets data. To get it back use [Notice.GetData] method.
func (msg *Notice) SetData(key string, val any) *Notice {
	if msg.Data == nil {
		msg.Data = make(map[string]any)
	}
	msg.Data[key] = val
	return msg
}

// GetData returns data set by [Notice.SetData]. Returns nil and false if the
// key was never set.
func (msg *Notice) GetData(key string) (any, bool) {
	if msg.Data == nil {
		return nil, false
	}
	val, ok := msg.Data[key]
	return val, ok
}

// equalizeNames returns a slice of row names where each name is the same
// length as the longest row name. The shorter names have spaces prepended to
// their names. The order of the returned slice is the same as Order field.
func (msg *Notice) equalizeNames() []string {
	var maxLen int
	for _, name := range msg.Order {
		if maxLen < len(name) {
			maxLen = len(name)
		}
	}
	var names []string
	for _, name := range msg.Order {
		name = strings.Repeat(" ", maxLen-len(name)) + name
		names = append(names, name)
	}
	return names
}
