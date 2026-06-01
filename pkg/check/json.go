// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"encoding/json"

	"github.com/ctx42/testing/pkg/notice"
)

// Text is a type constraint that allows the generic parameter to be either
// string or []byte representing text.
type Text interface {
	string | []byte
}

// JSON checks that two JSON texts are equivalent (after unmarshalling).
// See [assert.JSON].
//
// Example:
//
//	check.JSON(`{"hello": "world"}`, `{"foo": "bar"}`)
func JSON[W, H Text](want W, have H, opts ...any) error {
	var wantItf, haveItf any

	ops := DefaultOptions(opts...)
	if err := json.Unmarshal(toBytes(want), &wantItf); err != nil {
		msg := notice.New("did not expect the unmarshalling error").
			Append("argument", "want").
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	if err := json.Unmarshal(toBytes(have), &haveItf); err != nil {
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

// toBytes converts a string or []byte value to []byte.
func toBytes[T Text](v T) []byte {
	if s, ok := any(v).(string); ok {
		return []byte(s)
	}
	return any(v).([]byte) // nolint: forcetypeassert
}
