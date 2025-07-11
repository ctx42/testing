// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package timekit

import (
	"fmt"
	"time"
)

func ExampleClockFixed() {
	clk := ClockFixed(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))

	fmt.Println(clk())
	fmt.Println(clk())
	fmt.Println(clk())
	// Output:
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 00:00:00 +0000 UTC
}

func ExampleClockDeterministic() {
	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	clk := ClockDeterministic(start, time.Hour)

	fmt.Println(clk())
	fmt.Println(clk())
	fmt.Println(clk())
	// Output:
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 01:00:00 +0000 UTC
	// 2022-01-01 02:00:00 +0000 UTC
}

func ExampleTikTak() {
	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	clk := TikTak(start)

	fmt.Println(clk())
	fmt.Println(clk())
	fmt.Println(clk())
	// Output:
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 00:00:01 +0000 UTC
	// 2022-01-01 00:00:02 +0000 UTC
}
