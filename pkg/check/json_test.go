// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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
		have := " { \"hello\"\t:\n\n \"world\"\n\n\t}    "

		// --- When ---
		err := JSON(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("not equal", func(t *testing.T) {
		// --- Given ---
		want := ` {"hello": "world"} `
		have := " { \"hello\"\t:\n\n \"ms\"\n\n\t}    "
		opt := WithTrail("type.field")

		// --- When ---
		err := JSON(want, have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected JSON strings to be equal:\n" +
			"  trail: type.field\n" +
			"   want: {\"hello\":\"world\"}\n" +
			"   have: {\"hello\":\"ms\"}"
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
}
