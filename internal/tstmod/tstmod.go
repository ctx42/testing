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
}

// New creates a new module in a temporary directory.
//
// The version parameter must be one of the following:
//
//   - "v1" - uses tst-a@v0.1.0 and tst-b@v0.1.0 modules.
//   - "v2" - uses tst-a@v0.2.0 and tst-b@v0.2.0 modules.
func New(t tester.T, version string) *Module {
	t.Helper()

	mod := &Module{
		Version: version,
		Dir:     filepath.Join(t.TempDir(), "project"),
		t:       t,
	}

	mod.CreateDir(mod.Dir)
	mod.CreateDir("pkg/empty")
	mod.CreateDir("pkg/mercury")
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
func (mod *Module) CreateDir(pth string) string {
	mod.t.Helper()
	pth = filepath.Join(mod.Dir, pth)
	if err := os.MkdirAll(pth, 0777); err != nil {
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
			"require github.com/ctx42/tst-a v0.1.0 // indirect\n"

	case "v2":
		code = "" +
			"module github.com/ctx42/tst-project\n\n" +
			"go 1.24.0\n\n" +
			"require github.com/ctx42/tst-b v0.2.0\n\n" +
			"require github.com/ctx42/tst-a v0.2.0 // indirect\n"

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
