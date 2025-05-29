package notice

import (
	"fmt"
)

func FWD(head *Notice) string {
	// TODO(rz):
	var have string
	if head != nil && head.prev == nil {
		have += "| "
	}
	for head != nil {
		have += fmt.Sprintf("%s (%s)", head.Header, head.Trail)
		head = head.next
		if head != nil {
			have += " -> "
		}
	}
	return have
}

func REV(head *Notice) string {
	// TODO(rz):
	var have string
	if head != nil && head.next == nil {
		have += "| "
	}
	for head != nil {
		have += fmt.Sprintf("%s (%s)", head.Header, head.Trail)
		head = head.prev
		if head != nil {
			have += " -> "
		}
	}
	return have
}
