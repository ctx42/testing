// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_ChannelWillClose(t *testing.T) {
	t.Run("closed", func(t *testing.T) {
		// --- Given ---
		c := make(chan struct{})
		done := make(chan struct{})

		// --- When ---
		var err error
		go func() { err = ChannelWillClose("1s", c); close(done) }()

		// --- Then ---
		close(c)
		<-done
		affirm.Nil(t, err)
	})

	t.Run("nil channel", func(t *testing.T) {
		// --- Given ---
		var c chan struct{}

		// --- When ---
		err := ChannelWillClose("1s", c)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("timeout", func(t *testing.T) {
		// --- Given ---
		c := make(chan struct{})
		defer close(c)
		opt := WithTrail("type.field")

		// --- When ---
		err := ChannelWillClose("5ms", c, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"timeout waiting for channel to close:\n" +
			"   trail: type.field\n" +
			"  within: 5ms"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid duration", func(t *testing.T) {
		// --- Given ---
		c := make(chan struct{})
		defer close(c)

		// --- When ---
		err := ChannelWillClose("abc", c)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[within] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})
}
