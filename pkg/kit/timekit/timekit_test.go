// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package timekit

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

var past = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
var future = time.Now().Add(time.Hour).Truncate(time.Second)

func Test_ClockStartingAt_past(t *testing.T) {
	t.Run("past", func(t *testing.T) {
		// --- Given ---
		clk := ClockStartingAt(past)

		// --- When ---
		tim1 := clk()
		tim2 := clk()

		// --- Then ---
		assert.True(t, tim1.Before(tim2))
		assert.True(t, tim2.Before(time.Now()))
	})

	t.Run("future", func(t *testing.T) {
		// --- Given ---
		clk := ClockStartingAt(future)

		// --- When ---
		tim1 := clk()
		tim2 := clk()

		// --- Then ---
		assert.True(t, tim1.Before(tim2))
		assert.True(t, tim2.After(time.Now()))
	})
}

func Test_ClockFixed(t *testing.T) {
	// --- Given ---
	want := time.Now()
	clk := ClockFixed(want)

	// --- When ---
	tim0 := clk()
	tim1 := clk()
	tim2 := clk()

	// --- Then ---
	assert.Equal(t, want, tim0)
	assert.Equal(t, want, tim1)
	assert.Equal(t, want, tim2)
}

func Test_ClockDeterministic(t *testing.T) {
	// --- Given ---
	clk := ClockDeterministic(past, time.Second)

	// --- When ---
	tim0 := clk()
	tim1 := clk()
	tim2 := clk()

	// --- Then ---
	assert.Equal(t, past.Add(0*time.Second), tim0)
	assert.Equal(t, past.Add(1*time.Second), tim1)
	assert.Equal(t, past.Add(2*time.Second), tim2)
}

func Test_TikTak(t *testing.T) {
	// --- Given ---
	clk := TikTak(past)

	// --- When ---
	tim0 := clk()
	tim1 := clk()
	tim2 := clk()

	// --- Then ---
	assert.Equal(t, past.Add(0*time.Second), tim0)
	assert.Equal(t, past.Add(1*time.Second), tim1)
	assert.Equal(t, past.Add(2*time.Second), tim2)
}
