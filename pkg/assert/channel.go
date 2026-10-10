// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// ChannelWillClose asserts that the channel will be closed within the given
// duration. Values received before the channel closes are discarded, and a
// nil channel is reported as closed.
//
// See [check.ChannelWillClose] for details on the timeout format and
// the package documentation for options.
func ChannelWillClose[C any](
	t tester.T,
	within any,
	c <-chan C,
	opts ...any,
) bool {

	t.Helper()
	if e := check.ChannelWillClose(within, c, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
