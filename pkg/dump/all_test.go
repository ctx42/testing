// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package dump

// Node is a tree node used to test dumping nested pointer slices.
type Node struct {
	Value    int
	Children []*Node
}

// nodeTree returns a three-level [Node] tree.
func nodeTree() *Node {
	return &Node{
		Value: 1,
		Children: []*Node{
			{Value: 2},
			{Value: 3, Children: []*Node{{Value: 4}}},
		},
	}
}
