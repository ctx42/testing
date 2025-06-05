// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"os"
	"path/filepath"
)

// modCache returns the absolute path to the Go module cache joined with dir.
// It uses the GOMODCACHE environment variable with fallback to GOPATH.
func modCache(dir string) string {
	mod := os.Getenv("GOMODCACHE")
	if mod == "" {
		mod = os.Getenv("GOPATH")
		if mod == "" {
			mod = filepath.Join(os.Getenv("HOME"), "go", "pkg", "mod")
		} else {
			mod = filepath.Join(mod, "pkg", "mod")
		}
	}
	return filepath.Join(mod, dir)
}
