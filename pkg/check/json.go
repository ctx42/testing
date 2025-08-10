// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"encoding/json"

	"github.com/ctx42/testing/pkg/notice"
)

// JSON checks that two JSON strings are equivalent. Returns nil if they are,
// otherwise it returns an error with a message indicating the expected and
// actual values.
//
// Example:
//
//	check.JSON(`{"hello": "world"}`, `{"foo": "bar"}`)
func JSON(want, have string, opts ...any) error {
	var wantItf, haveItf any

	ops := DefaultOptions(opts...)
	if err := json.Unmarshal([]byte(want), &wantItf); err != nil {
		msg := notice.New("did not expect the unmarshalling error").
			Append("argument", "want").
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if err := json.Unmarshal([]byte(have), &haveItf); err != nil {
		msg := notice.New("did not expect the unmarshalling error").
			Append("argument", "have").
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}

	if err := Equal(wantItf, haveItf, WithOptions(ops)); err != nil {
		w, _ := json.Marshal(wantItf) // nolint:errchkjson
		h, _ := json.Marshal(haveItf) // nolint:errchkjson
		msg := notice.New("expected JSON strings to be equal").
			Want("%v", string(w)).
			Have("%v", string(h))
		return AddRows(ops, msg)
	}
	return nil
}
