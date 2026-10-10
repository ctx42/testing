// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ctx42/testing/pkg/notice"
)

// Text is a type constraint that allows the generic parameter to be either
// string or []byte representing text.
type Text interface {
	string | []byte
}

// JSON checks that two JSON texts are equivalent (after unmarshalling).
// See assert.JSON.
//
// Example:
//
//	check.JSON(`{"hello": "world"}`, `{"foo": "bar"}`)
func JSON[W, H Text](want W, have H, opts ...any) error {
	ops := DefaultOptions(opts...)
	wantItf, err := unmarshalJSON(toBytes(want))
	if err != nil {
		msg := notice.New("did not expect the unmarshalling error").
			Append("argument", "want").
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}
	haveItf, err := unmarshalJSON(toBytes(have))
	if err != nil {
		msg := notice.New("did not expect the unmarshalling error").
			Append("argument", "have").
			Append("error", "%s", err)
		return AddRows(ops, msg)
	}

	if err = Equal(wantItf, haveItf, WithOptions(ops)); err != nil {
		w, _ := json.Marshal(wantItf) //nolint:errchkjson
		h, _ := json.Marshal(haveItf) //nolint:errchkjson
		msg := notice.New("expected JSON strings to be equal").
			Want("%v", string(w)).
			Have("%v", string(h))
		return AddRows(ops, msg)
	}
	return nil
}

// unmarshalJSON decodes data keeping numbers as [json.Number] in canonical
// form (see canonicalJSON), so integers too large for float64 are compared
// exactly.
func unmarshalJSON(data []byte) (any, error) {
	var val any
	if err := json.Unmarshal(data, &val); err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&val); err != nil {
		return nil, err
	}
	return canonicalJSON(val), nil
}

// canonicalJSON rewrites every [json.Number] in val so equal numbers have
// equal text: integral values as exact integers ("1.0" and "1e2" become "1"
// and "100"), others as the shortest float64 representation.
func canonicalJSON(val any) any {
	switch v := val.(type) {
	case map[string]any:
		for key, elem := range v {
			v[key] = canonicalJSON(elem)
		}
		return v

	case []any:
		for i, elem := range v {
			v[i] = canonicalJSON(elem)
		}
		return v

	case json.Number:
		flt, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return v
		}
		// The float64 is finite here, so the exponent is bounded.
		rat, ok := new(big.Rat).SetString(string(v))
		if ok && rat.IsInt() {
			return json.Number(rat.Num().String())
		}
		return json.Number(strconv.FormatFloat(flt, 'g', -1, 64))

	default:
		return val
	}
}

// toBytes converts a string or []byte value to []byte.
func toBytes[T Text](v T) []byte {
	if s, ok := any(v).(string); ok {
		return []byte(s)
	}
	return any(v).([]byte) //nolint:forcetypeassert
}
