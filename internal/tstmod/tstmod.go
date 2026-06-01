// Package tstmod provides helpers for creating temporary Go modules
// with controlled external dependencies. It is used internally to test
// the mocker package's module and package discovery logic.
package tstmod

import (
	"os"
	"path/filepath"

	"github.com/ctx42/testing/pkg/tester"
)

// Module represents a test module used in tests.
//
// Uses:
//   - github.com/ctx42/tst-a@v0.1.0
//   - github.com/ctx42/tst-b@v0.1.0
//   - github.com/ctx42/tst-a@v0.2.0
//   - github.com/ctx42/tst-b@v0.2.0
type Module struct {
	Version string   // Module version.
	Dir     string   // Absolute path to the module root directory.
	t       tester.T // Test manager.

	// ExternalDirs maps module path@version to the absolute directory
	// containing that module's source (created via replace directives).
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

	base := t.TempDir()
	mod := &Module{
		Version:      version,
		Dir:          filepath.Join(base, "project"),
		t:            t,
		ExternalDirs: make(map[string]string),
	}

	mod.CreateDir(mod.Dir)
	mod.CreateDir("pkg/empty")
	mod.CreateDir("pkg/mercury")

	// Create external test modules (tst-a, tst-b) so that getModInfo
	// tests do not depend on the global module cache. This makes the
	// tests hermetic and pass both locally and on GitHub Actions.
	mod.createExternalModules(base)

	mod.WriteFile("go.mod", mod.goMod())
	mod.WriteFile("go.sum", mod.goSum())
	mod.WriteFile("project.go", mod.projectGo())
	mod.WriteFile("pkg/mercury/mercury.go", "package mercury")

	return mod
}

// WriteFile writes a file, rooted at the module's root directory, with the
// given name and content. Calls t.Fatal on error.
func (mod *Module) WriteFile(pth, content string) string {
	mod.t.Helper()
	pth = filepath.Join(mod.Dir, pth)
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
	mod.t.Helper()
	return filepath.Join(append([]string{mod.Dir}, elems...)...)
}

// goMod returns "go.mod" file content for the test project based on the
// project version. Panics if the version is unknown.
func (mod *Module) goMod() string {
	mod.t.Helper()

	var code string
	switch mod.Version {
	case "v1":
		code = "" +
			"module github.com/ctx42/tst-project\n\n" +
			"go 1.24.0\n\n" +
			"require github.com/ctx42/tst-b v0.1.0\n\n" +
			"require github.com/ctx42/tst-a v0.1.0 // indirect\n\n" +
			"replace github.com/ctx42/tst-a v0.1.0 => ../tst-a\n" +
			"replace github.com/ctx42/tst-b v0.1.0 => ../tst-b\n"

	case "v2":
		code = "" +
			"module github.com/ctx42/tst-project\n\n" +
			"go 1.24.0\n\n" +
			"require github.com/ctx42/tst-b v0.2.0\n\n" +
			"require github.com/ctx42/tst-a v0.2.0 // indirect\n\n" +
			"replace github.com/ctx42/tst-a v0.2.0 => ../tst-a\n" +
			"replace github.com/ctx42/tst-b v0.2.0 => ../tst-b\n"

	default:
		mod.t.Fatalf("unknown test project version: %s", mod.Version)
	}

	return code
}

// goSum returns "go.sum" file content for the test project based on the
// project version. Panics if the version is unknown.
func (mod *Module) goSum() string {
	mod.t.Helper()

	var code string
	switch mod.Version {
	case "v1":
		code = "" +
			"github.com/ctx42/tst-a v0.1.0 h1:XyxFm6pNY+vWXQp2pXjBHbWSEDFbWMwyU0x8OMmmYgk=\n" +
			"github.com/ctx42/tst-a v0.1.0/go.mod h1:sdk0IipiroBzJ93xBUDEYmGg9jE/pIKKwGgOxAS8M9g=\n" +
			"github.com/ctx42/tst-b v0.1.0 h1:jWHJOcnj9mmkKAszIL9cW+Rn+EWP2ccKaPPdTKCM0G0=\n" +
			"github.com/ctx42/tst-b v0.1.0/go.mod h1:fWWfo/LxuXFxKqCF3JA8Xp7AU2DwRmrLi+OB+Zy5bN0=\n"

	case "v2":
		code = "" +
			"github.com/ctx42/tst-a v0.2.0 h1:htAY7tEaalz2nVreCpP3hm8m0Bs7S3AGL5Mh7WSf6ls=\n" +
			"github.com/ctx42/tst-a v0.2.0/go.mod h1:sdk0IipiroBzJ93xBUDEYmGg9jE/pIKKwGgOxAS8M9g=\n" +
			"github.com/ctx42/tst-b v0.2.0 h1:NcGrkt9nplwE9R2WkcxQHWA00+Kq3S0PuubZV7WnOao=\n" +
			"github.com/ctx42/tst-b v0.2.0/go.mod h1:QYSBV5cYU5SOYsefdpuky2h5UzzJlARHYKfC7jOyTBI=\n"

	default:
		mod.t.Fatalf("unknown test project version: %s", mod.Version)
	}

	return code
}

// projectGo returns "project.go" file content for the test project based on
// the project version. Panics if the version is unknown.
func (mod *Module) projectGo() string {
	mod.t.Helper()

	var code string
	switch mod.Version {
	case "v1":
		code = "" +
			"package project\n\n" +
			"import \"github.com/ctx42/tst-b/pkg/mocker/first\"\n\n" +
			"type Project interface {\n" +
			"\tName()\n" +
			"\tfirst.First\n" +
			"}\n"

	case "v2":
		code = "" +
			"package project\n\n" +
			"import \"github.com/ctx42/tst-b/pkg/mocker/first\"\n\n" +
			"type Project interface {\n" +
			"\tName()\n" +
			"\tfirst.First\n" +
			"}\n"

	default:
		mod.t.Fatalf("unknown test project version: %s", mod.Version)
	}

	return code
}

// createExternalModules creates minimal source trees for tst-a and tst-b
// inside the same temporary directory as the main module and records their
// locations so that getModInfo tests do not depend on the global module cache.
func (mod *Module) createExternalModules(baseDir string) {
	var versions []struct {
		path    string
		version string
	}

	switch mod.Version {
	case "v1":
		versions = []struct {
			path    string
			version string
		}{
			{"github.com/ctx42/tst-a", "v0.1.0"},
			{"github.com/ctx42/tst-b", "v0.1.0"},
		}
	case "v2":
		versions = []struct {
			path    string
			version string
		}{
			{"github.com/ctx42/tst-a", "v0.2.0"},
			{"github.com/ctx42/tst-b", "v0.2.0"},
		}
	default:
		return
	}

	for _, v := range versions {
		dir := filepath.Join(baseDir, filepath.Base(v.path))
		mod.CreateDir(dir)
		mod.writeExternalGoMod(dir, v.path)

		key := v.path + "@" + v.version
		mod.ExternalDirs[key] = dir
	}
}

// writeExternalGoMod writes a minimal go.mod for an external test module.
func (mod *Module) writeExternalGoMod(dir, modulePath string) {
	content := "module " + modulePath + "\n\ngo 1.24.0\n"
	pth := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(pth, []byte(content), 0600); err != nil {
		mod.t.Fatal(err)
	}
}
