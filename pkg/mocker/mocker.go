// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"errors"
)

// TODO(rz): Test interface mocks generated from interfaces using types imported
//  with dot import.

// Sentinel errors.
var (
	// ErrInvPkg is returned when a directory or an import path does not point
	// to a valid Go package.
	//
	// This error occurs in the following cases:
	//   - The provided path cannot be resolved to a valid import path.
	//   - The import path cannot be resolved to a valid directory.
	ErrInvPkg = errors.New("invalid package")

	// ErrUnkMet is returned when an interface method cannot be found.
	ErrUnkMet = errors.New("method not found")

	// ErrUnkPkg is returned when a package cannot be found.
	ErrUnkPkg = errors.New("package not found")
)
