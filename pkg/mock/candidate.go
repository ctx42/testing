// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mock

// candidate implements logic to answer a question which Call is a better match.
type candidate struct {
	call    *Call    // A Call instance we are comparing.
	diff    []string // Argument diff.
	diffCnt int      // Number of arguments not matching.
}

// betterThan returns true if this candidate is better than the other one. The
// other is better if more of its arguments match and the method name is the
// same.
func (can candidate) betterThan(other candidate) bool {
	this := can.call
	if this == nil {
		return false
	}
	if other.call == nil {
		return true
	}
	if this.Method != other.call.Method {
		return true
	}
	if can.diffCnt > other.diffCnt {
		return false
	}
	return true
}
