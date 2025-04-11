// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"errors"
	"strings"
	"unsafe"
)

// Indent indents lines with n number of runes. Lines are indented only if
// there are more than one line.
func Indent(n int, r rune, lns string) string {
	if lns == "" {
		return ""
	}
	rows := strings.Split(lns, "\n")
	if len(rows) == 1 {
		return lns
	}
	for i, lin := range rows {
		var ind string
		if lin != "" {
			ind = strings.Repeat(string(r), n)
		}
		rows[i] = ind + lin
	}
	return strings.Join(rows, "\n")
}

// Unwrap unwraps joined errors. Returns nil if err is nil, unwraps only
// non-nil errors.
func Unwrap(err error) []error {
	if err == nil {
		return nil
	}
	var ers []error
	if es, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range es.Unwrap() {
			ers = append(ers, e)
		}
	} else {
		ers = append(ers, err)
	}
	return ers
}

// Join wraps errors in an instance of multi decorator if it's an error joined
// with [errors.Join].
func Join(err ...error) error {
	var ers []error
	for _, e := range err {
		if e == nil {
			continue
		}
		if es, ok := e.(interface{ Unwrap() []error }); ok {
			ers = append(ers, es.Unwrap()...)
		} else {
			ers = append(ers, e)
		}
	}
	switch len(ers) {
	case 0:
		return nil
	case 1:
		return ers[0]
	default:
		return multi{ers: ers}
	}
}

// multi is a decorator that formats multiple errors for output. It includes
// specialized handling for Notice errors but supports any error type.
type multi struct{ ers []error }

func (e multi) Error() string {
	if len(e.ers) == 1 {
		return e.ers[0].Error()
	}

	var prev string
	var msg *Notice
	if errors.As(e.ers[0], &msg) {
		prev = msg.Header
	}
	buf := []byte(e.ers[0].Error())

	for _, err := range e.ers[1:] {
		if errors.As(err, &msg) {
			tmp := msg.Header
			if prev == msg.Header {
				msg.Header = ContinuationHeader
				buf = append(buf, '\n')
				buf = append(buf, msg.Error()...)
				msg.Header = tmp
				continue
			}
		}
		prev = ""
		buf = append(buf, '\n', '\n')
		buf = append(buf, err.Error()...)
	}

	return unsafe.String(&buf[0], len(buf))
}

func (e multi) Unwrap() []error { return e.ers }
