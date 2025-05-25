// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"strings"
)

// Indent indents lines with n number of runes. Lines are indented only if
// there are more than one line.
func Indent(n int, r rune, lns string) string {
	if lns == "" {
		return ""
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

// Pad left pads the string with spaces to the given length.
func Pad(str string, length int) string {
	l := len(str)
	if length > l {
		return strings.Repeat(" ", length-l) + str
	}
	return str
}

// TrialCmp is a comparison function for sorting Notice instances by their
// Trail values. It returns:
//
//	-1 "a" should come before "b"
//	 0 equal
//	 1 "a" should come after "b"
func TrialCmp(a, b *Notice) int {
	if a.Trail < b.Trail {
		return -1
	} else if a.Trail > b.Trail {
		return 1
	}
	return 0
}

// SortNotices sorts a doubly linked list of Notice instances starting at head,
// ordering nodes by their values in ascending order. It modifies the list
// in-place by updating prev and next pointers. The "cmp" function takes two
// [Notice] instances and returns -1 if "a" should come before "b", 0 if equal,
// or 1 if "a" should come after "b". If the head is nil or the list has one
// node, it returns the unchanged head. Returns the tail of the sorted list so
// it can be used directly with [Join] or [Node.Chain] to add more nodes.
func SortNotices(head *Notice, cmp func(a, b *Notice) int) *Notice {
	if head == nil || head.next == nil {
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

	next:
		for node := sorted; node != nil; node = node.next {
			switch result := cmp(current, node); result {
			case -1:
				break next
			case 0:
				found = node
			case 1:
				found = node
			}
		}

		// Detach node.
		if current.next != nil {
			current.next.prev = current.prev
		}

		switch {
		case found == nil:
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
