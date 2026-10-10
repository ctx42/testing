// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ctx42/testing/pkg/notice"
)

// Regexp checks that "want" regexp matches the string form of "have".
//
// "want" may be a regexp string or [*regexp.Regexp]. [fmt.Sprint] is used
// on "have".
//
// See assert.Regexp for the assertion wrapper.
func Regexp(want, have any, opts ...any) error {
	match, err := matchRegexp(want, have)
	if err != nil {
		ops := DefaultOptions(opts...)
		msg := notice.New("expected valid regexp").Append("error", "%q", err)
		return AddRows(ops, msg)
	}
	if match {
		return nil
	}
	ops := DefaultOptions(opts...)
	msg := notice.New("expected regexp to match").
		Append("regexp", "%s", want).
		Have("%q", have)
	return AddRows(ops, msg)
}

// matchRegexp returns true if the regexp matches the string form of have.
func matchRegexp(rx, have any) (bool, error) {
	var r *regexp.Regexp
	if rr, ok := rx.(*regexp.Regexp); ok {
		if rr == nil {
			return false, errors.New("nil regexp")
		}
		r = rr
	} else {
		var err error
		rxs := fmt.Sprint(rx)
		if r, err = regexp.Compile(rxs); err != nil {
			return false, err
		}
	}
	return r.FindStringIndex(fmt.Sprint(have)) != nil, nil
}
