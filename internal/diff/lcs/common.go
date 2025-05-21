// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcs

import (
	"log"
	"sort"
)

// lcs is the Longest Common Sequence (LCS).
type lcs []diag

// A diag is a piece of the edit graph where A[X+i] == B[Y+i], for 0<=i<Len.
// All computed diagonals are parts of the longest common subsequence.
type diag struct {
	X, Y int
	Len  int
}

// sort sorts in place, by the lowest X, and if tied, inversely by Len.
func (l lcs) sort() lcs {
	sort.Slice(l, func(i, j int) bool {
		if l[i].X != l[j].X {
			return l[i].X < l[j].X
		}
		return l[i].Len > l[j].Len
	})
	return l
}

// validate validates that the elements of the lcs do not overlap (can only
// happen when the two-sided algorithm ends early) expects the lcs to be sorted.
func (l lcs) valid() bool {
	for i := 1; i < len(l); i++ {
		if l[i-1].X+l[i-1].Len > l[i].X {
			return false
		}
		if l[i-1].Y+l[i-1].Len > l[i].Y {
			return false
		}
	}
	return true
}

// repair overlapping lcs only called if two-sided stops early
func (l lcs) fix() lcs {
	// from the set of diagonals in l, find a maximal non-conflicting set
	// this problem may be NP-complete, but we use a greedy heuristic,
	// which is quadratic, but with a better data structure, could be D log D.
	// independent is not enough: {0,3,1} and {3,0,2} can't both occur in an
	// lcs which has to have monotone x and y
	if len(l) == 0 {
		return nil
	}
	sort.Slice(l, func(i, j int) bool { return l[i].Len > l[j].Len })
	tmp := make(lcs, 0, len(l))
	tmp = append(tmp, l[0])
	for i := 1; i < len(l); i++ {
		var dir direction
		nxt := l[i]
		for _, in := range tmp {
			if dir, nxt = overlap(in, nxt); dir == empty || dir == bad {
				break
			}
		}
		if nxt.Len > 0 && dir != bad {
			tmp = append(tmp, nxt)
		}
	}
	tmp.sort()
	if false && !tmp.valid() { // Debug checking.
		log.Fatalf("here %d", len(tmp))
	}
	return tmp
}

type direction int

const (
	empty    direction = iota // Diag is empty (so not in lcs).
	leftdown                  // Proposed acceptably to the left and below.
	rightup                   // The proposed diag is acceptably to the right and above.
	bad                       // The proposed diag is inconsistent with the lcs so far.
)

// overlap trims the proposed diag prop so it doesn't overlap with the existing
// diag that has already been added to the lcs.
//
// nolint: cyclop
func overlap(exist, prop diag) (direction, diag) {
	if prop.X <= exist.X && exist.X < prop.X+prop.Len {
		// Remove the end of the prop where it overlaps with the X end of exist.
		delta := prop.X + prop.Len - exist.X
		prop.Len -= delta
		if prop.Len <= 0 {
			return empty, prop
		}
	}
	if exist.X <= prop.X && prop.X < exist.X+exist.Len {
		// Remove the beginning of prop where overlaps with exist.
		delta := exist.X + exist.Len - prop.X
		prop.Len -= delta
		if prop.Len <= 0 {
			return empty, prop
		}
		prop.X += delta
		prop.Y += delta
	}
	if prop.Y <= exist.Y && exist.Y < prop.Y+prop.Len {
		// Remove the end of the prop that overlaps (in Y) with exist.
		delta := prop.Y + prop.Len - exist.Y
		prop.Len -= delta
		if prop.Len <= 0 {
			return empty, prop
		}
	}
	if exist.Y <= prop.Y && prop.Y < exist.Y+exist.Len {
		// Remove the beginning of the prop that overlaps with exist.
		delta := exist.Y + exist.Len - prop.Y
		prop.Len -= delta
		if prop.Len <= 0 {
			return empty, prop
		}
		prop.X += delta // No test reaches this code.
		prop.Y += delta
	}
	if prop.X+prop.Len <= exist.X && prop.Y+prop.Len <= exist.Y {
		return leftdown, prop
	}
	if exist.X+exist.Len <= prop.X && exist.Y+exist.Len <= prop.Y {
		return rightup, prop
	}
	// Prop can't be in an lcs that contains exist
	return bad, prop
}

// manipulating Diag and lcs

// prepend a diagonal (x,y)-(x+1,y+1) segment either to an empty lcs
// or to its first Diag. prepend is only called to extend diagonals in the
// backward direction.
func (l lcs) prepend(x, y int) lcs {
	if len(l) > 0 {
		d := &l[0]
		if int(d.X) == x+1 && int(d.Y) == y+1 {
			// extend the diagonal down and to the left
			d.X, d.Y = int(x), int(y)
			d.Len++
			return l
		}
	}

	r := diag{X: int(x), Y: int(y), Len: 1}
	return append([]diag{r}, l...)
}

// append appends a diagonal or extends the existing one by adding the edge
// (x,y)-(x+1.y+1). append is only called to extend diagonals in the forward
// direction.
func (l lcs) append(x, y int) lcs {
	if len(l) > 0 {
		last := &l[len(l)-1]
		// Expand last element if adjoining.
		if last.X+last.Len == x && last.Y+last.Len == y {
			last.Len++
			return l
		}
	}
	return append(l, diag{X: x, Y: y, Len: 1})
}

// enforce constraint on d, k.
func ok(d, k int) bool {
	return d >= 0 && -d <= k && k <= d
}
