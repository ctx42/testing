// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff_test

import (
	"testing"

	"github.com/ctx42/testing/internal/diff"
)

func TestMerge(t *testing.T) {
	// For convenience, we test Merge using strings, not sequences of edits,
	// though this does put us at the mercy of the diff algorithm.
	for _, test := range []struct {
		base, x, y string
		want       string // "!" => conflict
	}{
		// Independent insertions.
		{"abcdef", "abXcdef", "abcdeYf", "abXcdeYf"},
		// Independent deletions.
		{"abcdef", "acdef", "abcdf", "acdf"},
		// Colocated insertions (X first).
		{"abcdef", "abcXdef", "abcYdef", "abcXYdef"},
		// Colocated identical insertions (coalesced).
		{"abcdef", "abcXdef", "abcXdef", "abcXdef"},
		// Colocated insertions with common prefix (X first).
		// TODO(adonovan): would "abcXYdef" be better?
		//  i.e. should we dissect the insertions?
		{"abcdef", "abcXdef", "abcXYdef", "abcXXYdef"},
		// Mix of identical and independent insertions (X first).
		{"abcdef", "aIbcdXef", "aIbcdYef", "aIbcdXYef"},
		// Independent deletions.
		{"abcdef", "def", "abc", ""},
		// Overlapping deletions: conflict.
		{"abcdef", "adef", "abef", "!"},
		// Overlapping deletions with distinct insertions, X first.
		{"abcdef", "abXef", "abcYf", "!"},
		// Overlapping deletions with distinct insertions, Y first.
		{"abcdef", "abcXf", "abYef", "!"},
		// Overlapping deletions with common insertions.
		{"abcdef", "abXef", "abcXf", "!"},
		// Trailing insertions in X (observe X bias).
		{"abcdef", "aXbXcXdXeXfX", "aYbcdef", "aXYbXcXdXeXfX"},
		// Trailing insertions in Y (observe X bias).
		{"abcdef", "aXbcdef", "aYbYcYdYeYfY", "aXYbYcYdYeYfY"},
	} {
		dx := diff.Strings(test.base, test.x)
		dy := diff.Strings(test.base, test.y)
		got := "!" // conflict
		if dz, ok := diff.Merge(dx, dy); ok {
			var err error
			got, err = diff.Apply(test.base, dz)
			if err != nil {
				t.Errorf("Merge(%q, %q, %q) produced invalid edits %v: %v", test.base, test.x, test.y, dz, err)
				continue
			}
		}
		if test.want != got {
			t.Errorf("base=%q x=%q y=%q: got %q, want %q", test.base, test.x, test.y, got, test.want)
		}
	}
}
