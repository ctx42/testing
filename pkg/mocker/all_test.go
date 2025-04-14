package mocker

import (
	"os"
	"path/filepath"

	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/notice"
)

func init() {
	check.RegisterTypeChecker(gopkg{}, gopkgCheck)
	check.RegisterTypeChecker(argument{}, argumentCheck)
	check.RegisterTypeChecker(method{}, methodCheck)
}

// TODO(rz): it would be cool if we could log to the test log from the checkers.

// gopkgCheck is a custom checker, matching [check.Check] signature, comparing
// two instances of gopkg.
func gopkgCheck(want, have any, opts ...check.Option) error {
	if err := check.Type(gopkg{}, have, opts...); err != nil {
		return err
	}
	w, h := want.(gopkg), have.(gopkg)

	Name := check.Name(opts)

	ers := []error{
		check.Equal(w.alias, h.alias, Name("alias")),
		check.Equal(w.pkgName, h.pkgName, Name("pkgName")),
		check.Equal(w.pkgPath, h.pkgPath, Name("pkgPath")),
		check.Equal(w.pkgDir, h.pkgDir, Name("pkgDir")),
		check.Equal(w.modName, h.modName, Name("modName")),
		check.Equal(w.modPath, h.modPath, Name("modPath")),
		check.Equal(w.modDir, h.modDir, Name("modDir")),
		check.Equal(w.wd, h.wd, Name("wd")),
		check.Fields(9, w, Name("{gopkg field count}")),
	}

	return notice.Join(ers...)
}

// argumentCheck is a custom checker, matching [check.Check] signature,
// comparing two instances of type argument.
func argumentCheck(want, have any, opts ...check.Option) error {
	if err := check.Type(argument{}, have, opts...); err != nil {
		return err
	}
	w, h := want.(argument), have.(argument)

	Name := check.Name(opts)

	ers := []error{
		check.Equal(w.name, h.name, Name("name")),
		check.Equal(w.typ, h.typ, Name("typ")),
		check.Equal(w.pks, h.pks, Name("pks")),
		check.Fields(3, w, Name("{argument field count}")),
	}

	return notice.Join(ers...)
}

// methodCheck is a custom checker, matching [check.Check] signature, comparing
// two instances of type method.
func methodCheck(want, have any, opts ...check.Option) error {
	if err := check.Type(method{}, have, opts...); err != nil {
		return err
	}
	w, h := want.(method), have.(method)

	Name := check.Name(opts)

	ers := []error{
		check.Equal(w.name, h.name, Name("name")),
		check.Equal(w.args, h.args, Name("args")),
		check.Equal(w.rets, h.rets, Name("rets")),
		check.Fields(3, w, Name("{method field count}")),
	}

	return notice.Join(ers...)
}

// goModCache returns the absolute path to the Go module cache joined with dir.
// It uses the GOMODCACHE environment variable, falling back to GOPATH if unset.
func goModCache(dir string) string {
	mod := os.Getenv("GOMODCACHE")
	if mod == "" {
		mod = filepath.Join(os.Getenv("GOPATH"), "pkg/mod")
	}
	return filepath.Join(mod, dir)
}
