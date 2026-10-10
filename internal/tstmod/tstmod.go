// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package tstmod provides helpers for creating temporary Go modules
// with controlled external dependencies. It is used internally to test
// the mocker package's module and package discovery logic.
package tstmod

import (
	"os"
	"path/filepath"

	"github.com/ctx42/testing/pkg/tester"
)

// Paths of the external modules the test module depends on.
const (
	tstA = "github.com/ctx42/tst-a"
	tstB = "github.com/ctx42/tst-b"
)

// variant describes the external module version a test module variant uses
// and the go.sum hashes of the published modules. The external modules are
// replaced by local directories, so Go never checks the hashes; they keep the
// go.sum file realistic.
type variant struct {
	version string // Version of both tst-a and tst-b.
	sumA    string // The h1 hash of tst-a.
	sumAMod string // The h1 hash of tst-a go.mod.
	sumB    string // The h1 hash of tst-b.
	sumBMod string // The h1 hash of tst-b go.mod.
}

// variants maps a test module variant to its definition.
var variants = map[string]variant{
	"v1": {
		version: "v0.1.0",
		sumA:    "h1:XyxFm6pNY+vWXQp2pXjBHbWSEDFbWMwyU0x8OMmmYgk=",
		sumAMod: "h1:sdk0IipiroBzJ93xBUDEYmGg9jE/pIKKwGgOxAS8M9g=",
		sumB:    "h1:jWHJOcnj9mmkKAszIL9cW+Rn+EWP2ccKaPPdTKCM0G0=",
		sumBMod: "h1:fWWfo/LxuXFxKqCF3JA8Xp7AU2DwRmrLi+OB+Zy5bN0=",
	},
	"v2": {
		version: "v0.2.0",
		sumA:    "h1:htAY7tEaalz2nVreCpP3hm8m0Bs7S3AGL5Mh7WSf6ls=",
		sumAMod: "h1:sdk0IipiroBzJ93xBUDEYmGg9jE/pIKKwGgOxAS8M9g=",
		sumB:    "h1:NcGrkt9nplwE9R2WkcxQHWA00+Kq3S0PuubZV7WnOao=",
		sumBMod: "h1:QYSBV5cYU5SOYsefdpuky2h5UzzJlARHYKfC7jOyTBI=",
	},
}

// Module is a temporary Go module depending on the tst-a and tst-b modules
// in the version of its variant. Build it with [New]; the zero value is not
// usable.
type Module struct {
	Version string   // Module variant: "v1" or "v2".
	Dir     string   // Absolute path to the module root directory.
	t       tester.T // Test manager.

	// ExternalDirs maps module path@version to the absolute directory the
	// module is replaced with in the go.mod file.
	ExternalDirs map[string]string
}

// New creates a new module in a temporary directory.
//
// The version parameter must be one of the following:
//
//   - "v1" - uses tst-a@v0.1.0 and tst-b@v0.1.0 modules.
//   - "v2" - uses tst-a@v0.2.0 and tst-b@v0.2.0 modules.
func New(t tester.T, version string) *Module {
	t.Helper()
	vrt, ok := variants[version]
	if !ok {
		t.Fatalf("unknown test project version: %s", version)
		return nil
	}

	base := t.TempDir()
	mod := &Module{
		Version:      version,
		Dir:          filepath.Join(base, "project"),
		t:            t,
		ExternalDirs: make(map[string]string),
	}

	mod.CreateDir("pkg/empty")

	// Create external test modules (tst-a, tst-b) so that getModInfo
	// tests do not depend on the global module cache. This makes the
	// tests hermetic and pass both locally and on GitHub Actions.
	mod.createExternalModules(base, vrt.version)

	mod.WriteFile("go.mod", goMod(vrt.version))
	mod.WriteFile("go.sum", vrt.goSum())
	mod.WriteFile("project.go", projectGo)
	mod.WriteFile("pkg/mercury/mercury.go", "package mercury")

	return mod
}

// WriteFile writes a file with the given name and content, creating missing
// parent directories, and returns its path. A relative pth is rooted at the
// module's root directory; an absolute one is used as-is. Calls t.Fatal on
// error.
func (mod *Module) WriteFile(pth, content string) string {
	mod.t.Helper()
	if !filepath.IsAbs(pth) {
		pth = filepath.Join(mod.Dir, pth)
	}
	if err := os.MkdirAll(filepath.Dir(pth), 0700); err != nil {
		mod.t.Fatal(err)
	}
	if err := os.WriteFile(pth, []byte(content), 0600); err != nil {
		mod.t.Fatal(err)
	}
	return pth
}

// CreateDir creates a new directory rooted at the project root directory.
// If pth is absolute, it is used as-is.
func (mod *Module) CreateDir(pth string) string {
	mod.t.Helper()
	if !filepath.IsAbs(pth) {
		pth = filepath.Join(mod.Dir, pth)
	}
	if err := os.MkdirAll(pth, 0700); err != nil {
		mod.t.Fatal(err)
	}
	return pth
}

// Path returns a path described by elements rooted at the test module root.
func (mod *Module) Path(elems ...string) string {
	return filepath.Join(append([]string{mod.Dir}, elems...)...)
}

// projectGo is the content of the test module's project.go file.
const projectGo = "" +
	"package project\n\n" +
	"import \"github.com/ctx42/tst-b/pkg/mocker/first\"\n\n" +
	"type Project interface {\n" +
	"\tName()\n" +
	"\tfirst.First\n" +
	"}\n"

// goMod returns the go.mod file content of the test module using the given
// version of the external modules.
func goMod(version string) string {
	return "" +
		"module github.com/ctx42/tst-project\n\n" +
		"go 1.24.0\n\n" +
		"require " + tstB + " " + version + "\n\n" +
		"require " + tstA + " " + version + " // indirect\n\n" +
		"replace " + tstA + " " + version + " => ../tst-a\n" +
		"replace " + tstB + " " + version + " => ../tst-b\n"
}

// goSum returns the go.sum file content of the test module variant.
func (vrt variant) goSum() string {
	ver := vrt.version
	return "" +
		tstA + " " + ver + " " + vrt.sumA + "\n" +
		tstA + " " + ver + "/go.mod " + vrt.sumAMod + "\n" +
		tstB + " " + ver + " " + vrt.sumB + "\n" +
		tstB + " " + ver + "/go.mod " + vrt.sumBMod + "\n"
}

// createExternalModules creates the tst-a and tst-b modules inside baseDir
// and records their locations, so that getModInfo tests do not depend on the
// global module cache. Each module has a go.mod file; tst-b also requires
// tst-a and provides the pkg/mocker/first package that project.go imports.
func (mod *Module) createExternalModules(baseDir, version string) {
	mod.t.Helper()
	dirA := mod.CreateDir(filepath.Join(baseDir, "tst-a"))
	goModA := "module " + tstA + "\n\ngo 1.24.0\n"
	mod.WriteFile(filepath.Join(dirA, "go.mod"), goModA)
	mod.ExternalDirs[tstA+"@"+version] = dirA

	dirB := mod.CreateDir(filepath.Join(baseDir, "tst-b"))
	goModB := "" +
		"module " + tstB + "\n\n" +
		"go 1.24.0\n\n" +
		"require " + tstA + " " + version + "\n"
	mod.WriteFile(filepath.Join(dirB, "go.mod"), goModB)
	first := "package first\n\ntype First interface{ First() }\n"
	mod.WriteFile(filepath.Join(dirB, "pkg/mocker/first/first.go"), first)
	mod.ExternalDirs[tstB+"@"+version] = dirB
}
