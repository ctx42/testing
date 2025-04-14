// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"go/ast"
)

// task represents a single task for generating an interface mock.
type task struct {
	Config                    // Embedded configuration for the task.
	pkg    *gopkg             // Package containing the interface to mock.
	file   file               // File with the interface to mock.
	itf    *ast.InterfaceType // Interface to mock.
}
