// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package kit provides a curated collection of sub-packages and helpers
// designed to streamline common testing tasks in Go.
//
// In addition to the sub-packages, a few top-level helpers are provided:
//   - SHA1 file/reader utilities
//   - [AddGlobalCleanup] / [RunGlobalCleanups] — global post-test cleanup
//     coordination (see important warnings in the godoc)
//
// Current sub-packages:
//   - [iokit] — I/O and buffer helpers with error injection and test
//     cleanup support
//   - [timekit] — controllable and deterministic clocks for
//     time-dependent testing
//   - [reflectkit] — lightweight reflection utilities (primarily for
//     struct field inspection in tests)
//
// See the package [README] for an overview of the collection and links to
// the sub-package documentation.
package kit
