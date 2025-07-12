// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"path/filepath"
	"strings"
)

// gopkg represents Go package.
type gopkg struct {
	// Alias used when generating code for the package import line.
	alias string

	// Is the package name as used in code (e.g., "http", "exec").
	pkgName string

	// Is the package import path (e.g., "net/http", "os/exec").
	pkgPath string

	// Is the absolute path to the package source directory.
	pkgDir string

	// Is the name of the module containing the package, as used in code (e.g.,
	// "http" for "net/http", "exec" for "os/exec").
	modName string

	// Is the full import path of the module containing the package (e.g.,
	// "net/http", "github.com/user/mod").
	modPath string

	// Is the absolute path to the root directory of the module containing the
	// package.
	modDir string

	// Is the absolute path to the working directory in which the `go list`
	// command is executed.
	wd string

	// Indicates the package has been resolved (found).
	resolved bool

	// Source files in the package.
	files []*file

	// Manages sources and their positions.
	fset *token.FileSet
}

// newPkg returns a new instance of gopkg for the given working directory
// (expects the path to be an absolute path) and optional import path. If the
// import path is set, the pkgName field will be set as well. You need to call
// resolve method to get all the fields set. Uses the provided working
// directory to run `go list`.
//
// Examples:
//
//	newPkg("/module/dir", "")
//	newPkg("/module/dir", "github.com/pkg/name")
func newPkg(dir, pth string) *gopkg {
	return &gopkg{
		pkgName: assumedPackageName(pth),
		pkgPath: pth,
		wd:      dir,
	}
}

// resolve finds the package and module it belongs to based on fields set on the
// instance. The minimum information needed is the package directory or its
// import path.
func (pkg *gopkg) resolve() error {
	if pkg.resolved {
		return nil
	}
	if err := pkg.getPkgInfo(); err == nil {
		pkg.resolved = true
		return nil
	}
	if err := pkg.getModInfo(); err != nil {
		return err
	}
	pkg.resolved = true
	return nil
}

// getPkgInfo uses `go list` to retrieve package information.
func (pkg *gopkg) getPkgInfo() (err error) {
	var out []byte
	args := []string{"list", "-json"}
	if id := pkg.id(); id != "" {
		args = append(args, id)
	}
	if pkg.wd == "" {
		pkg.wd = pkg.pkgDir
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = pkg.wd
	if out, err = cmd.Output(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return fmt.Errorf("%w: %s", ErrUnkPkg, ee.Stderr)
		}
		return fmt.Errorf("%w: %w", ErrUnkPkg, err)
	}
	p := struct {
		ImportPath string
		Dir        string
		Name       string
		Module     struct {
			Path string
			Dir  string
		}
	}{}
	if err = json.Unmarshal(out, &p); err != nil { // nolint: musttag
		return fmt.Errorf("%w: %w", ErrUnkPkg, err)
	}

	pkg.modName = assumedPackageName(p.Module.Path)
	pkg.modPath = p.Module.Path
	pkg.modDir = p.Module.Dir

	pkg.pkgName = p.Name
	pkg.pkgPath = p.ImportPath
	pkg.pkgDir = p.Dir

	return nil
}

// getModInfo retrieves module and package information for the package.
//
// nolint: gocognit, cyclop
func (pkg *gopkg) getModInfo() (err error) {
	var out []byte
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = pkg.wd
	if out, err = cmd.Output(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return fmt.Errorf("%w: %s", ErrUnkPkg, ee.Stderr)
		}
		return fmt.Errorf("%w: %w", ErrUnkPkg, err)
	}
	mod := struct {
		Path string
		Dir  string
	}{}

	var i int
	dec := json.NewDecoder(bytes.NewReader(out))
	for dec.More() {
		i++
		if err = dec.Decode(&mod); err != nil { // nolint: musttag
			return fmt.Errorf("%w: %w", ErrUnkPkg, err)
		}
		// The module in the working directory is the first on the list.
		if i == 1 {
			pkg.modName = assumedPackageName(mod.Path)
			pkg.modPath = mod.Path
			pkg.modDir = mod.Dir

			if pkg.pkgPath == "" {
				if sub, err := filepath.Rel(mod.Dir, pkg.wd); err == nil {
					if sub == "." {
						pkg.pkgName = assumedPackageName(mod.Dir)
					} else {
						pkg.pkgName = assumedPackageName(sub)
					}
					pkg.pkgPath = filepath.Join(mod.Path, sub)
					pkg.pkgDir = filepath.Join(mod.Dir, sub)
					break
				}
			} else if strings.HasPrefix(pkg.pkgPath, pkg.modPath) {
				if sub, err := filepath.Rel(mod.Path, pkg.pkgPath); err == nil {
					pkg.pkgName = assumedPackageName(pkg.pkgPath)
					pkg.pkgPath = filepath.Join(mod.Path, sub)
					pkg.pkgDir = filepath.Join(mod.Dir, sub)
					break
				}
			}
			continue
		}

		if mod.Path == pkg.pkgPath {
			pkg.pkgName = assumedPackageName(pkg.pkgPath)
			pkg.pkgDir = mod.Dir
			break
		}

		if strings.HasPrefix(pkg.pkgPath, mod.Path+"/") {
			if sub, err := filepath.Rel(mod.Path, pkg.pkgPath); err == nil {
				pkg.pkgName = assumedPackageName(pkg.pkgPath)
				pkg.pkgDir = filepath.Join(mod.Dir, sub)
				break
			}
		}

	}

	if !pkg.isValid() {
		return fmt.Errorf("%w: %s", ErrUnkPkg, pkg.id())
	}
	return nil
}

// parse reads and processes source files (excluding test files) in the package.
// For efficiency, it parses files only once, even if called multiple times.
func (pkg *gopkg) parse() error {
	if err := pkg.resolve(); err != nil {
		return err
	}
	if pkg.files != nil {
		return nil
	}
	names, err := findSources(pkg.pkgDir)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnkPkg, err)
	}
	pkg.fset = token.NewFileSet()
	pkg.files = make([]*file, 0, len(names))
	for _, name := range names {
		astFil, err := parser.ParseFile(pkg.fset, name, nil, 0)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrAstParse, err)
		}
		pkg.files = append(pkg.files, &file{path: name, ast: astFil})
	}
	return nil
}

// findType locates a type declaration named `name` in the package. It returns
// the containing file, the type's AST node, and nil error on success.
func (pkg *gopkg) findType(name string) (*file, *ast.TypeSpec, error) {
	if err := pkg.parse(); err != nil {
		return nil, nil, err
	}
	for _, fil := range pkg.files {
		if typ := fil.findType(name); typ != nil {
			return fil, typ, nil
		}
	}
	return nil, nil, fmt.Errorf("%w: %s", ErrUnkType, name)
}

// findItf locates an interface type declaration named `name` in the package.
// It returns the containing file, the interface's AST node, and nil error if
// the named type is an interface.
func (pkg *gopkg) findItf(name string) (*file, *ast.InterfaceType, error) {
	fil, typ, err := pkg.findType(name)
	if err != nil {
		return nil, nil, err
	}
	if itf, ok := typ.Type.(*ast.InterfaceType); ok {
		return fil, itf, nil
	}
	return nil, nil, fmt.Errorf("%w: %s is not an interface", ErrUnkItf, name)
}

// isValid returns true if the package has all the required fields set.
func (pkg *gopkg) isValid() bool {
	return pkg.pkgName != "" && pkg.pkgPath != "" && pkg.pkgDir != "" &&
		pkg.modName != "" && pkg.modPath != "" && pkg.modDir != "" &&
		pkg.wd != ""
}

// id returns package identification which is either import path or package
// source directory. Used in errors to identify problematic package.
func (pkg *gopkg) id() string {
	if pkg.pkgPath != "" {
		return pkg.pkgPath
	}
	return pkg.pkgDir
}

// equal checks if two packages are equivalent for the purpose of reuse.
//
// It returns true if the provided package has the same identity as the current
// package, allowing the caller to decide whether to resolve and parse the
// other package or reuse this one. Both packages must be non-nil.
func (pkg *gopkg) equal(other *gopkg) bool {
	if pkg.pkgPath != "" && pkg.pkgPath == other.pkgPath {
		return true
	}
	if pkg.pkgDir != "" && pkg.pkgDir == other.pkgDir {
		return true
	}
	return false
}

// from copies all fields (except alias) from the `other` to this one.
func (pkg *gopkg) from(other *gopkg) {
	pkg.pkgName = other.pkgName
	pkg.pkgPath = other.pkgPath
	pkg.pkgDir = other.pkgDir
	pkg.modName = other.modName
	pkg.modPath = other.modPath
	pkg.modDir = other.modDir
	pkg.wd = other.wd
	pkg.resolved = other.resolved
	pkg.files = other.files
	pkg.fset = other.fset
}

// setAlias sets package [gopkg.alias] and [gopkg.pkgName] to given alias.
func (pkg *gopkg) setAlias(alias string) {
	pkg.alias = alias
	if alias == "" {
		pkg.pkgName = assumedPackageName(pkg.pkgPath)
	} else {
		pkg.pkgName = alias
	}
}

// genImport generates code for the package import line. The dot alias is never
// used.
//
// Example:
//
//	"fmt"
//	"github.com/user/project/pkg/package"
//	alias "github.com/user/project/pkg/package"
func (pkg *gopkg) genImport() string {
	if pkg.pkgPath == "" {
		return ""
	}
	code := `"` + pkg.pkgPath + `"`
	if pkg.alias != "" && pkg.alias != "." {
		code = pkg.alias + " " + code
	}
	return code
}
