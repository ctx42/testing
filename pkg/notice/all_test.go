// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice

import (
	"fmt"
)

// FWD is a test helper function printing notice chain in the forward direction
// starting at the given node. If the given node has no previous node, the
// leading pipe will be printed.
//
// Example:
//
//	| hdr0 (A) -> hdr1 (B) # The given node doesn't have the previous node.
//	hdr0 (A) -> hdr1 (B) # The given node has the previous node.
func FWD(node *Notice) string {
	var have string
	if node != nil && node.prev == nil {
		have += "| "
	}
	for node != nil {
		have += fmt.Sprintf("%s (%s)", node.Header, node.Trail)
		node = node.next
		if node != nil {
			have += " -> "
		}
	}
	return have
}

// REV is a test helper function printing notice chain in the backward
// direction starting at the given node. If the given node has no next node,
// the leading pipe will be printed.
//
// Example:
//
//	| hdr0 (A) -> hdr1 (B) # The given node doesn't have the next node.
//	hdr0 (A) -> hdr1 (B) # The given node has the next node.
func REV(node *Notice) string {
	var have string
	if node != nil && node.next == nil {
		have += "| "
	}
	for node != nil {
		have += fmt.Sprintf("%s (%s)", node.Header, node.Trail)
		node = node.prev
		if node != nil {
			have += " -> "
		}
	}
	return have
}
