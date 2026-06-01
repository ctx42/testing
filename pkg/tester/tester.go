// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package tester provides the [T] interface (a deliberate subset of
// [testing.TB]) and the [Spy] implementation for testing custom test helpers
// and assertion functions.
//
// Writing reusable test helpers is a common practice for keeping tests
// readable and DRY. Those helpers themselves need to be tested. The tester
// package makes this possible without circular dependencies or fragile
// test setups.
//
// [Spy] is the primary tool: pass it to your Helper Under Test (HUT) in
// place of *testing.T, set expectations with the Expect* methods, call
// [Spy.Close], exercise the helper, then verify with [Spy.AssertExpectations]
// (or let the automatic cleanup do it).
//
// See the package [README] for the full usage guide, lifecycle patterns,
// log matching strategies, and examples. See [examples/tester.go] and
// [examples/tester_test.go] (in the module root) for complete, tested
// demonstrations.
//
// Key types and entry points:
//   - [T] — the interface your helpers should accept
//   - [New] — creates a [Spy] bound to a real *testing.T
//   - [Spy] — the spy implementation with Expect*, Close, AssertExpectations
//   - [Strategy] + ExpectLog* — for precise log message verification
package tester
