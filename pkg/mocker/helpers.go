// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"bytes"
	"fmt"
	"go/ast"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// specDir returns the absolute path to a module or package directory based on
// an import spec.
//
// It executes the `go list` command with the provided import spec to resolve
// the directory path. The wd argument specifies the working directory for the
// `go list` command, which is necessary to:
//
//   - Establish the context for module resolution, as `go list` uses the
//     working directory to locate the nearest go.mod file.
//   - Support relative import paths or module lookups in projects with
//     multiple modules or vendored dependencies.
//   - Ensure correct resolution when the import specification is relative to
//     the module or GOPATH in the working directory.
//
// Example:
//
//	github.com/user/project -> /path/to/project
//	github.com/user/project/pkg/package -> /path/to/project/pkg/package
func specDir(wd, impSpec string) (string, error) {
	eout := &bytes.Buffer{}
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", impSpec)
	cmd.Stderr = eout
	cmd.Dir = wd
	out, err := cmd.Output()
	if err != nil {
		msg := fmtCmdError(eout.String())
		return "", fmt.Errorf("%w: %s", ErrInvSpec, msg)
	}
	return strings.TrimSpace(string(out)), nil
}

// dirToSpec returns the import path for a Go package or module based on the
// provided directory path.
//
// It executes the `go list` command in the specified directory (dir) to
// resolve the import path. The wd argument specifies the directory containing
// the Go package or module, which is necessary to:
//
//   - Provide the context for module resolution, as `go list` uses the
//     directory to locate the nearest go.mod file or GOPATH structure.
//   - Ensure accurate import path resolution for packages within a module or
//     GOPATH, especially in multi-module projects or complex directory layouts.
//   - Support resolution when the directory is not the current working
//     directory of the process.
//
// Example:
//
//	/path/to/project -> github.com/user/project
//	/path/to/project/pkg/package -> github.com/user/project/pkg/package
func dirToSpec(dir string) (string, error) {
	eout := &bytes.Buffer{}
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}")
	cmd.Stderr = eout
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		msg := fmtCmdError(eout.String())
		return "", fmt.Errorf("%w: %s", ErrInvSpec, msg)
	}
	return strings.TrimSpace(string(out)), nil
}

// settleImport resolves and completes the fields of an [Import] struct based
// on the provided values. It populates missing fields using the provided ones,
// treating a zero-value [Import] (empty Spec and Dir) as referring to the
// working directory (wd). It may return partially modified [Import] in case of
// an error.
func settleImport(imp Import) (Import, error) {
	if imp.IsZero() {
		return imp, ErrInvImport
	}

	var err error
	if imp.Spec != "" {
		if imp.Dir, err = specDir(imp.Dir, imp.Spec); err != nil {
			return imp, err
		}
	} else {
		if imp.Spec, err = dirToSpec(imp.Dir); err != nil {
			return imp, err
		}
	}
	imp.Name = assumedPackageName(imp.Spec)
	if imp.Dir, err = filepath.Abs(imp.Dir); err != nil {
		return imp, err
	}
	return imp, nil
}

// findSources returns a list of paths to all Go source files (excluding test
// files) in the specified directory. It does not recurse into subdirectories.
// The returned paths are absolute.
func findSources(dir string) ([]string, error) {
	var err error
	if dir, err = filepath.Abs(dir); err != nil {
		return nil, err
	}

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	ets, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var sources []string
	for _, entry := range ets {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		sources = append(sources, filepath.Join(dir, name))
	}
	sort.Strings(sources)
	return sources, nil
}

// assumedPackageName returns the assumed package name of an import path.
// It does this using only string parsing of the import path. It picks the last
// element of the path that does not look like a major version and then picks
// the valid identifier off the start of that element. It is used to determine
// if a local rename should be added to an import for clarity.
//
// Copied from: https://github.com/golang/tools/blob/a318c19ff2fd8d6aae74e36fe7e1a8b8afef3bf7/internal/imports/fix.go#L1233
//
// Example:
//
//	github.com/user/project/pkg/package -> package
//	github.com/user/project/pkg/go_package -> go_package
//	github.com/user/project/pkg/go-package-abc -> abc
func assumedPackageName(importPath string) string {
	notIdentifier := func(ch rune) bool {
		return !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
			'0' <= ch && ch <= '9' ||
			ch == '_' ||
			ch >= utf8.RuneSelf &&
				(unicode.IsLetter(ch) || unicode.IsDigit(ch)))
	}

	base := path.Base(importPath)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			dir := path.Dir(importPath)
			if dir != "." {
				base = path.Base(dir)
			}
		}
	}
	parts := strings.Split(base, "-")
	base = parts[len(parts)-1]
	if i := strings.IndexFunc(base, notIdentifier); i >= 0 {
		base = base[:i]
	}
	return base
}

// toLowerSnakeCase converts camel case to lowercase snake case.
func toLowerSnakeCase(camel string) string {
	var runes = make([]rune, 0, len(camel)+10)
	const lower = 1
	const upper = 2

	var prev int
	for i := 0; i < len(camel); i++ {
		curr := lower
		r := rune(camel[i])
		if unicode.IsUpper(r) {
			curr = upper
		}
		if prev == lower && curr == upper {
			runes = append(runes, '_', r)
		} else {
			runes = append(runes, r)
		}
		prev = curr
	}
	return strings.ToLower(string(runes))
}

// space is a regular expression matching one or more whitespace characters.
var space = regexp.MustCompile(`\s+`)

// fmtCmdError formats the error output from a command by making it a single
// line string without duplicate whitespaces.
func fmtCmdError(out string) string {
	out = space.ReplaceAllString(out, " ")
	return strings.TrimSpace(out)
}

// astMethods returns the interface methods as a slice. If the interface has no
// methods, it returns an [ErrNoMethods].
func astMethods(itf *ast.InterfaceType) ([]*ast.Field, error) {
	// TODO(rz): this was moved to Mocker. Remove.
	if itf.Methods == nil || len(itf.Methods.List) == 0 {
		return nil, ErrNoMethods
	}
	return itf.Methods.List, nil
}

// findMethod finds a method in the slice by name. It returns nil if method is
// not found.
func findMethod(mts []*Method, name string) *Method {
	// TODO(rz): test this.
	// TODO(rz): document this.
	for _, met := range mts {
		if met.name == name {
			return met
		}
	}
	return nil
}

// builtinTypes list of builtin types.
var builtinTypes = map[string]bool{
	// TODO(rz): how do we get this list?
	"bool":       true,
	"byte":       true,
	"complex128": true,
	"complex64":  true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"int8":       true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uint8":      true,
	"uintptr":    true,
	"any":        true,
}

// isBuiltinType returns true if a type represented by "typ" is a builtin type.
func isBuiltinType(typ string) bool {
	// TODO(rz): test this.
	// TODO(rz): document this.
	_, ok := builtinTypes[typ]
	return ok
}
