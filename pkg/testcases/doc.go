// Package testcases provides a rich, curated collection of test values and
// cases designed specifically for people writing custom assertions or using
// [tester.Spy] to test their own test helpers.
//
// See the Design section in the root README for the overall layered
// architecture and the recommended patterns for testing custom helpers with
// [tester.Spy].
//
// The package contains:
//   - Primitive values of every Go kind (including many edge cases)
//   - Structs exercising exported/unexported field combinations
//   - Types implementing common interfaces
//   - Ready-made collections via [EqualCases], [ZENValues], etc.
//
// These values are heavily used inside this module itself to verify
// [check.Equal], [assert.Equal], custom type checkers, and Spy-based helper
// tests. They are the recommended data source so you do not have to hand-craft
// dozens of interesting cases when validating your own checkers across Go's
// type system.
//
// Example usage when writing a custom assertion:
//
//	func AssertMyType(t tester.T, want, have MyType) bool {
//	    t.Helper()
//	    return assert.Equal(t, want, have, testcases.WithTrail("MyType"))
//	}
//
// See the [examples] package for complete, tested patterns of using
// testcases together with [tester.Spy].
package testcases
