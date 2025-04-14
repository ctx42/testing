package tstmod

import (
	"os"
	"path/filepath"

	"github.com/ctx42/testing/pkg/tester"
)

// TODO(rz): test this. Add list dir assertion?

// Module represents a test module used in tests.
//
// Uses:
//   - github.com/ctx42/tst-a@v0.1.0
//   - github.com/ctx42/tst-b@v0.1.0
//   - github.com/ctx42/tst-a@v0.2.0
//   - github.com/ctx42/tst-b@v0.2.0
type Module struct {
	Version string   // Project version.
	Dir     string   // The project's absolute directory path.
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

	prj := &Module{
		Version: version,
		Dir:     filepath.Join(t.TempDir(), "project"),
		t:       t,
	}

	if err := os.Mkdir(prj.Dir, 0777); err != nil {
		t.Fatal(err)
	}

	// Create go source files in the project directory.
	dst := filepath.Join(prj.Dir, "go.mod")
	if err := os.WriteFile(dst, prj.goMod(), 0600); err != nil {
		t.Fatal(err)
	}

	dst = filepath.Join(prj.Dir, "go.sum")
	if err := os.WriteFile(dst, prj.goSum(), 0600); err != nil {
		t.Fatal(err)
	}

	dst = filepath.Join(prj.Dir, "project.go")
	if err := os.WriteFile(dst, prj.projectGo(), 0600); err != nil {
		t.Fatal(err)
	}

	mer := filepath.Join(prj.Dir, "pkg/mercury")
	if err := os.MkdirAll(mer, 0777); err != nil {
		t.Fatal(err)
	}

	dst = filepath.Join(mer, "mercury.go")
	if err := os.WriteFile(dst, []byte("package mercury"), 0600); err != nil {
		t.Fatal(err)
	}

	empty := filepath.Join(prj.Dir, "pkg/empty")
	if err := os.MkdirAll(empty, 0777); err != nil {
		t.Fatal(err)
	}

	return prj
}

// goMod returns "go.mod" file content for the test project based on the
// project version. Panics if the version is unknown.
func (prj *Module) goMod() []byte {
	prj.t.Helper()

	var code string
	switch prj.Version {
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
		prj.t.Fatalf("unknown test project version: %s", prj.Version)
	}

	return []byte(code)
}

// goSum returns "go.sum" file content for the test project based on the
// project version. Panics if the version is unknown.
func (prj *Module) goSum() []byte {
	prj.t.Helper()

	var code string
	switch prj.Version {
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
		prj.t.Fatalf("unknown test project version: %s", prj.Version)
	}
	return []byte(code)
}

// projectGo returns "project.go" file content for the test project based on
// the project version. Panics if the version is unknown.
func (prj *Module) projectGo() []byte {
	prj.t.Helper()

	var code string
	switch prj.Version {
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
		prj.t.Fatalf("unknown test project version: %s", prj.Version)
	}

	return []byte(code)
}
