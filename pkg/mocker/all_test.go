// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

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
	check.RegisterTypeChecker(file{}, fileCheck)
}

// gopkgCheck is a custom checker, matching [check.Check] signature, comparing
// two instances of gopkg.
func gopkgCheck(want, have any, opts ...check.Option) error {
	ops := check.DefaultOptions(opts...)
	if err := check.Type(gopkg{}, have, check.WithOptions(ops)); err != nil {
		return err
	}
	w, h := want.(gopkg), have.(gopkg)

	fName := check.FieldName(ops, "gopkg")
	ers := []error{
		check.Equal(w.alias, h.alias, fName("alias")),
		check.Equal(w.pkgName, h.pkgName, fName("pkgName")),
		check.Equal(w.pkgPath, h.pkgPath, fName("pkgPath")),
		check.Equal(w.pkgDir, h.pkgDir, fName("pkgDir")),
		check.Equal(w.modName, h.modName, fName("modName")),
		check.Equal(w.modPath, h.modPath, fName("modPath")),
		check.Equal(w.modDir, h.modDir, fName("modDir")),
		check.Equal(w.wd, h.wd, fName("wd")),
		check.Equal(w.resolved, h.resolved, fName("resolved")),
		check.Equal(w.files, h.files, fName("files")),
		check.Fields(11, w, fName("{field count}")),
	}
	return notice.Join(ers...)
}

// argumentCheck is a custom checker, matching [check.Check] signature,
// comparing two instances of the type argument.
func argumentCheck(want, have any, opts ...check.Option) error {
	ops := check.DefaultOptions(opts...)
	if err := check.Type(argument{}, have, check.WithOptions(ops)); err != nil {
		return err
	}
	w, h := want.(argument), have.(argument)

	fName := check.FieldName(ops, "argument")
	ers := []error{
		check.Equal(w.name, h.name, fName("name")),
		check.Equal(w.typ, h.typ, fName("typ")),
		check.Equal(w.pks, h.pks, fName("pks")),
		check.Fields(3, w, fName("{field count}")),
	}
	return notice.Join(ers...)
}

// methodCheck is a custom checker, matching [check.Check] signature, comparing
// two instances of the type method.
func methodCheck(want, have any, opts ...check.Option) error {
	ops := check.DefaultOptions(opts...)
	if err := check.Type(method{}, have, check.WithOptions(ops)); err != nil {
		return err
	}
	w, h := want.(method), have.(method)

	fName := check.FieldName(ops, "method")
	ers := []error{
		check.Equal(w.name, h.name, fName("name")),
		check.Equal(w.args, h.args, fName("args")),
		check.Equal(w.rets, h.rets, fName("rets")),
		check.Fields(3, w, fName("{field count}")),
	}
	return notice.Join(ers...)
}

// fileCheck is a custom checker, matching [check.Check] signature, comparing
// two instances of the type file.
func fileCheck(want, have any, opts ...check.Option) error {
	ops := check.DefaultOptions(opts...)
	if err := check.Type(file{}, have, check.WithOptions(ops)); err != nil {
		return err
	}
	w, h := want.(file), have.(file)

	fName := check.FieldName(ops, "file")
	ers := []error{
		check.Equal(w.path, h.path, fName("path")),
		check.Equal(w.pks, h.pks, fName("pks")),
		check.Fields(4, w, fName("{field count}")),
	}
	return notice.Join(ers...)
}

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
