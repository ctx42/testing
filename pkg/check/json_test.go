// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_JSON(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		want := ` {"hello": "world"} `
		have := `{"hello": "world"}`

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("not equal", func(t *testing.T) {
		// --- Given ---
		want := ` {"hello": "world"} `
		have := `{"hello": "ms"}`
		opt := WithTrail("type.field")

		// --- When ---
		err := JSON(want, have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := `expected JSON strings to be equal:
  trail: type.field
   want: {"hello":"world"}
   have: {"hello":"ms"}`
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid want JSON", func(t *testing.T) {
		// --- Given ---
		want := `{!!!}`
		have := `{"hello": "world"}`
		opt := WithTrail("type.field")

		// --- When ---
		err := JSON(want, have, opt)

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

	t.Run("invalid have JSON", func(t *testing.T) {
		// --- Given ---
		want := `{"hello": "world"}`
		have := `{!!!}`
		opt := WithTrail("type.field")

		// --- When ---
		err := JSON(want, have, opt)

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
