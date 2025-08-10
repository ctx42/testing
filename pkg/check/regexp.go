// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"regexp"

	"github.com/ctx42/testing/pkg/notice"
)

// Regexp checks that "want" regexp matches "have". Returns nil if it does,
// otherwise, it returns an error with a message indicating the expected and
// actual values.
//
// The "want" can be either a regular expression string or instance of
// [regexp.Regexp]. The [fmt.Sprint] is used to get string representation of
// have argument.
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

// matchRegexp return true if a specified regexp matches a string.
func matchRegexp(rx, have any) (bool, error) {
	var r *regexp.Regexp
	if rr, ok := rx.(*regexp.Regexp); ok {
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
