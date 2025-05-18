// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package tstkit

import (
	"sync"
	"time"
)

// ClockStartingAt returns function with the same signature as [time.Now] and
// returning time as if the current time was set to given value.
func ClockStartingAt(tim time.Time) func() time.Time {
	now := time.Now()
	guard := sync.Mutex{}
	return func() time.Time {
		guard.Lock()
		defer guard.Unlock()
		return tim.Add(time.Now().Sub(now))
	}
}

// ClockFixed returns function with the same signature as [time.Now] which
// always returns the given time.
func ClockFixed(tim time.Time) func() time.Time {
	return func() time.Time {
		return tim
	}
}

// ClockDeterministic returns function with the same signature as [time.Now]
// and returning time advancing by given tick with every call no matter now
// fast or slow you call it.
func ClockDeterministic(start time.Time, tick time.Duration) func() time.Time {
	now := start.Add(-tick)
	guard := sync.Mutex{}
	return func() time.Time {
		guard.Lock()
		defer guard.Unlock()
		now = now.Add(tick)
		return now
	}
}

// TikTak returns a deterministic clock advancing one second for each call.
func TikTak(start time.Time) func() time.Time {
	return ClockDeterministic(start, time.Second)
}
