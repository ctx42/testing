// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"regexp"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Regexp(t *testing.T) {
	t.Run("match string", func(t *testing.T) {
		// --- When ---
		err := Regexp("^abc123.*$", "abc1234")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("match instance", func(t *testing.T) {
		// --- Given ---
		rx := regexp.MustCompile("^abc123.*$")

		// --- When ---
		err := Regexp(rx, "abc1234")

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("no match", func(t *testing.T) {
		// --- When ---
		err := Regexp("^abc42.*$", "abc1234")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected regexp to match:\n" +
			"  regexp: ^abc42.*$\n" +
			"    have: \"abc1234\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Regexp("^abc42.*$", "abc1234", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected regexp to match:\n" +
			"   trail: type.field\n" +
			"  regexp: ^abc42.*$\n" +
			"    have: \"abc1234\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid regexp", func(t *testing.T) {
		// --- When ---
		err := Regexp("[p", "abc1234")

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected valid regexp:\n" +
			"  error: \"error parsing regexp: missing closing ]: `[p`\""
		affirm.Equal(t, wMsg, err.Error())
	})
}
