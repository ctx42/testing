// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package notice

import (
	"strings"
)

// Indent prepends n copies of rune r to each line in lns.
// Single-line input is returned unchanged.
func Indent(n int, r rune, lns string) string {
	if lns == "" || n <= 0 {
		return lns
	}
	rows := strings.Split(lns, "\n")
	if len(rows) == 1 {
		return lns
	}
	for i, lin := range rows {
		var ind string
		if lin != "" {
			ind = strings.Repeat(string(r), n)
		}
		rows[i] = ind + lin
	}
	return strings.Join(rows, "\n")
}

// Pad left-pads str with spaces until it reaches at least the given length.
func Pad(str string, length int) string {
	l := len(str)
	if length > l {
		return strings.Repeat(" ", length-l) + str
	}
	return str
}

// TrailCmp compares two notices by their Trail strings.
// Returns -1 if x < y, 0 if equal, +1 if x > y.
// Suitable for use with [SortNotices] or slices.SortFunc.
func TrailCmp(x, y *Notice) int {
	if x.Trail < y.Trail {
		return -1
	}
	if x.Trail > y.Trail {
		return 1
	}
	return 0
}

// SortNotices sorts, in place and stably, the chain head belongs to using
// cmp. The cmp function returns a negative number if a should precede b,
// zero if they are equal, or a positive number if a should follow b.
//
// It mutates prev/next pointers. If head is nil or its chain has one node it
// is returned unchanged. The returned value is the tail of the sorted chain
// (useful for further [Join] or [Notice.Chain] calls).
func SortNotices(head *Notice, cmp func(a, b *Notice) int) *Notice {
	if head == nil {
		return nil
	}
	head = head.Head()
	if head.next == nil {
		return head
	}

	var found *Notice

	// We detach the head from the list by removing the pointer to it from the
	// next node and removing the pointer to the next node from itself.
	head.next.prev = nil
	current := head.next
	sorted := head
	sorted.next = nil
	teil := sorted

	for current != nil {
		found = nil
		next := current.next

		for node := sorted; node != nil; node = node.next {
			if cmp(current, node) < 0 {
				break
			}
			found = node
		}

		// Detach node.
		if current.next != nil {
			current.next.prev = current.prev
		}

		switch found {
		case nil:
			// Insert at the beginning.
			current.prev = nil
			current.next = sorted
			sorted.prev = current

			sorted = current

		default:
			// Insert after found.
			current.next = found.next
			if found.next != nil {
				found.next.prev = current
			}

			found.next = current
			current.prev = found

			//goland:noinspection GoDirectComparisonOfErrors
			if found == teil {
				teil = current
			}
		}

		current = next // Move to the next unsorted node.
	}

	return teil
}
