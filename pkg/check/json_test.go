// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_JSON(t *testing.T) {
	t.Run("error - integers beyond float64 precision", func(t *testing.T) {
		// --- Given ---
		want := `{"id": 9007199254740993}`
		have := `{"id": 9007199254740992}`

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		wMsg := "" +
			"expected JSON strings to be equal:\n" +
			"  want: {\"id\":9007199254740993}\n" +
			"  have: {\"id\":9007199254740992}"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal numbers in different notation", func(t *testing.T) {
		// --- Given ---
		want := `{"a": 1.0, "b": 1e2, "c": [0.5]}`
		have := `{"a": 1, "b": 100, "c": [5e-1]}`

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		want := ` {"hello": "world"} `
		have := `{"hello": "world"}`

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not equal", func(t *testing.T) {
		// --- Given ---
		want := ` {"hello": "world"} `
		have := `{"hello": "ms"}`

		// --- When ---
		err := JSON(want, have, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected JSON strings to be equal:\n" +
			"  trail: type.field\n" +
			"   want: {\"hello\":\"world\"}\n" +
			"   have: {\"hello\":\"ms\"}"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid want JSON", func(t *testing.T) {
		// --- Given ---
		want := `{!!!}`
		have := `{"hello": "world"}`

		// --- When ---
		err := JSON(want, have, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"did not expect the unmarshalling error:\n" +
			"     trail: type.field\n" +
			"  argument: want\n" +
			"     error: invalid character '!' looking for beginning of " +
			"object key string"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - invalid have JSON", func(t *testing.T) {
		// --- Given ---
		want := `{"hello": "world"}`
		have := `{!!!}`

		// --- When ---
		err := JSON(want, have, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"did not expect the unmarshalling error:\n" +
			"     trail: type.field\n" +
			"  argument: have\n" +
			"     error: invalid character '!' looking for beginning of " +
			"object key string"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("equal bytes want", func(t *testing.T) {
		// --- Given ---
		want := []byte(` {"hello": "world"} `)
		have := `{"hello": "world"}`

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal bytes have", func(t *testing.T) {
		// --- Given ---
		want := ` {"hello": "world"} `
		have := []byte(`{"hello": "world"}`)

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("equal bytes", func(t *testing.T) {
		// --- Given ---
		want := []byte(` {"hello": "world"} `)
		have := []byte(`{"hello": "world"}`)

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})
}
