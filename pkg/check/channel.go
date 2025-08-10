// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"time"

	"github.com/ctx42/testing/pkg/notice"
)

// ChannelWillClose checks the channel will be closed "within" a given time
// duration. Returns nil if it was, otherwise returns an error with a message
// indicating the expected and actual values.
//
// The "within" may represent duration in the form of a string, int, int64 or
// [time.Duration].
func ChannelWillClose[C any](within any, c <-chan C, opts ...any) error {
	if c == nil {
		return nil
	}

	dur, durStr, _, err := getDur(within, opts...)
	if err != nil {
		return notice.From(err, "within")
	}

	tim := time.NewTimer(dur)
	defer tim.Stop()

	for {
		select {
		case <-tim.C:
			ops := DefaultOptions(opts...)
			msg := notice.New("timeout waiting for channel to close").
				Append("within", "%s", durStr)
			return AddRows(ops, msg)

		case _, open := <-c:
			if !open {
				if !tim.Stop() {
					<-tim.C
				}
				return nil
			}
		}
	}
}
