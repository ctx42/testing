// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"time"

	"github.com/ctx42/testing/pkg/notice"
)

// ChannelWillClose checks channel will be closed "within" given time duration.
// Returns nil if it was, otherwise returns an error with a message indicating
// the expected and actual values.
//
// The "within" may represent duration in form of a string, int, int64 or
// [time.Duration].
func ChannelWillClose[C any](within any, c <-chan C, opts ...Option) error {
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
			return notice.New("timeout waiting for channel to close").
				SetTrail(ops.Trail).
				Append("within", "%s", durStr)

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
