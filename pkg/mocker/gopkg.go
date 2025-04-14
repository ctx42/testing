package mocker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO(rz): how will the alias be used? What if it is modified?

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

	// Is the working directory in which the `go list` command is executed.
	wd string

	// Source files in the package.
	files []file

	// Manages sources and their positions.
	fset *token.FileSet
}

// newPkg gathers information about the Go package and the module it belongs to
// based on the specified directory.
//
// Uses the provided directory as the working directory to run `go list`. The
// path may be empty.
//
// Examples:
//
//	newPkg(".", "")
//	newPkg("module/dir", "")
//	newPkg("module/dir", "github.com/pkg/name")
//	newPkg("/module/dir", "github.com/pkg/name")
func newPkg(dir, pth string) (_ *gopkg, err error) {
	if dir, err = filepath.Abs(dir); err != nil {
		return nil, err
	}
	pkg := &gopkg{wd: dir, pkgPath: pth}
	if err = pkg.getPkgInfo(); err == nil {
		return pkg, nil
	}
	if err = pkg.getModInfo(); err != nil {
		return nil, err
	}
	return pkg, nil
}

// getPkgInfo uses `go list` to retrieve package information for the specified
// import path. The wd parameter specifies the working directory for the
// `go list` command and must be a path within the module containing the import
// path.
func (pkg *gopkg) getPkgInfo() (err error) {
	var out []byte
	args := []string{"list", "-json"}
	if pkg.pkgPath != "" {
		args = append(args, pkg.pkgPath)
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = pkg.wd
	if out, err = cmd.Output(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return fmt.Errorf("%w: %s", ErrInvPkg, ee.Stderr)
		}
		return fmt.Errorf("%w: %s", ErrInvPkg, err)
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
	if err = json.Unmarshal(out, &p); err != nil {
		return fmt.Errorf("%w: %w", ErrInvPkg, err)
	}

	pkg.modName = assumedPackageName(p.Module.Path)
	pkg.modPath = p.Module.Path
	pkg.modDir = p.Module.Dir

	pkg.pkgName = p.Name
	pkg.pkgPath = p.ImportPath
	pkg.pkgDir = p.Dir

	return nil
}

// getModInfo retrieves module information for the specified directory using
// the `go list -m` command. The dir parameter must be a path within the target
// module.
func (pkg *gopkg) getModInfo() (err error) {
	var out []byte
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = pkg.wd
	if out, err = cmd.Output(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return fmt.Errorf("%w: %s", ErrInvPkg, ee.Stderr)
		}
		return fmt.Errorf("%w: %s", ErrInvPkg, err)
	}
	mod := struct {
		Path string
		Dir  string
	}{}

	var i int
	dec := json.NewDecoder(bytes.NewReader(out))
	for dec.More() {
		i++
		if err = dec.Decode(&mod); err != nil {
			return fmt.Errorf("%w: %w", ErrInvPkg, err)
		}
		// The module in working directory is the first on the list.
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
					return nil
				}
			} else if strings.HasPrefix(pkg.pkgPath, pkg.modPath) {
				if sub, err := filepath.Rel(mod.Path, pkg.pkgPath); err == nil {
					pkg.pkgName = assumedPackageName(pkg.pkgPath)
					pkg.pkgPath = filepath.Join(mod.Path, sub)
					pkg.pkgDir = filepath.Join(mod.Dir, sub)
					return nil
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

func (pkg *gopkg) parse() error {
	if pkg.files != nil {
		return nil
	}
	// TODO(rz):
	return nil
}

// isDot returns true if the instance represents a dot import.
func (pkg *gopkg) isDot() bool { return pkg.alias == "." }

// isValid returns true if package has all the required fields set.
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

// GoString generates and returns code representing the import line for the
// package. If the pkgPath field is empty, it returns an empty string.
//
// Example:
//
//	"fmt"
//	"github.com/user/project/pkg/package"
//	alias "github.com/user/project/pkg/package"
//	. "github.com/user/project/pkg/package"
func (pkg *gopkg) GoString() string {
	if pkg.pkgPath == "" {
		return ""
	}
	code := `"` + pkg.pkgPath + `"`
	if pkg.alias != "" {
		code = pkg.alias + " " + code
	}
	return code
}
