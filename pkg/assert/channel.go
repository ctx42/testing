// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// ChannelWillClose asserts the channel will be closed "within" a given time
// duration. Returns true if it was, otherwise marks the test as failed, writes
// an error message to the test log and returns false.
//
// The "within" may represent duration in the form of a string, int, int64 or
// [time.Duration].
func ChannelWillClose[C any](t tester.T, within any, c <-chan C, opts ...any) bool {
	t.Helper()
	if err := check.ChannelWillClose(within, c, opts...); err != nil {
		t.Error(err)
		return false
	}
	return true
}
