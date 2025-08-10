// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_ChannelWillClose(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		c := make(chan struct{})
		done := make(chan struct{})
		var have bool

		// --- When ---
		go func() { have = ChannelWillClose(tspy, "1s", c); close(done) }()

		// --- Then ---
		close(c)
		<-done
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		c := make(chan struct{})
		defer close(c)

		// --- When ---
		have := ChannelWillClose(tspy, "1s", c)

		// --- Then ---
		tspy.Finish().AssertExpectations()
		affirm.Equal(t, false, have)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("   trail: type.field")
		tspy.Close()

		c := make(chan struct{})
		defer close(c)
		opt := check.WithTrail("type.field")

		// --- When ---
		have := ChannelWillClose(tspy, "1s", c, opt)

		// --- Then ---
		tspy.Finish().AssertExpectations()
		affirm.Equal(t, false, have)
	})
}
