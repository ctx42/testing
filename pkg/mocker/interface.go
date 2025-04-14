// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// Interface represents an interface.
type Interface struct {
	name    string    // The interface name.
	methods []*Method // The interface methods.
}

// NewInterface creates a new [Interface] from a given AST file and node.
func NewInterface(name string, methods ...*Method) *Interface {
	return &Interface{name: name, methods: methods}
}

// Name returns the name of the interface.
func (itf *Interface) Name() string { return itf.name }
