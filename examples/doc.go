// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package examples contains demonstrations of advanced usage patterns for the
// testing module.
//
// The examples in this package show how to:
//
//   - Build custom assertion helpers on top of [assert] and [check] (see
//     [custom_assertions_test.go]).
//   - Write testable test helpers using [tester.T] and [tester.Spy] (see
//     [tester.go] and [tester_test.go]).
//   - Create and use mocks with the [mock] package (see [mock_test.go]).
//
// These files are intended as reference implementations rather than
// executable godoc examples.
package examples
