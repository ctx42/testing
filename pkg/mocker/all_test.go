// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"go/ast"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/tester"
)

const (
	// tstModImp is the import path used in the test module's code.
	tstModImp = "github.com/ctx42/tst-int/pkg/mocker/first"

	// tstModImp is the path used to cache the module used by the test module.
	tstModImpCached = "github.com/ctx42/tst-int@v0.1.1/pkg/mocker/first"
)

// createTestModule creates an example Go module in a temporary directory with
// imports from "github.com/ctx42/ctx-int-tst" repository. It returns an
// absolute path to the created module, which is removed automatically when the
// test ends.
func createTestModule(t tester.T) string {
	t.Helper()

	// tstModGo is "go.mod" content.
	const tstModGo = "" +
		"module github.com/ctx42/tst-project\n\n" +
		"go 1.21\n\n" +
		"require github.com/ctx42/tst-int v0.1.1 // indirect\n"

	// tstModSumGo is "go.sum" content.
	const tstModSumGo = "" +
		"github.com/ctx42/tst-int v0.1.1 " +
		"h1:r1h7BNp3A+NKy4ltWzNnrgD1fjm1pVFnSGxEnGMYRp8=\n" +
		"github.com/ctx42/tst-int v0.1.1/go.mod " +
		"h1:+aAusX6/kK+nYscB9k+ZylzlG1K7Da7RATdXuOKcq54=\n"

	// tstModCodeGo is "project.go" content.
	const tstModCodeGo = "" +
		"package project\n\n" +
		"import \"" + tstModImp + "\"\n\n" +
		"var T first.Medium\n\n" +
		"func fn() {}\n"

	dir := filepath.Join(t.TempDir(), "project")
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Error(err)
		return ""
	}

	// Create go source files in the project directory.
	dst := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(dst, []byte(tstModGo), 0600); err != nil {
		t.Error(err)
		return ""
	}
	dst = filepath.Join(dir, "go.sum")
	if err := os.WriteFile(dst, []byte(tstModSumGo), 0600); err != nil {
		t.Error(err)
		return ""
	}
	dst = filepath.Join(dir, "project.go")
	if err := os.WriteFile(dst, []byte(tstModCodeGo), 0600); err != nil {
		t.Error(err)
		return ""
	}
	return dir
}

// keys returns sorted string map keys.
func keys[T any](m map[string]T) []string {
	var ks []string
	for key := range m {
		ks = append(ks, key)
	}
	sort.Strings(ks)
	return ks
}

// findItf looks up an interface with the given name in the package specified
// by the import spec. It returns the AST file containing the interface
// definition and the AST node representing the interface. If the interface is
// not found or an error occurs (e.g., invalid import path or syntax error in
// the package), it marks the test as failed using the provided [tester.T] and
// returns nil for both the file and interface type.
func findItf(t tester.T, name, spec string) (*ast.File, *ast.InterfaceType) {
	t.Helper()
	imp, err := settleImport(NewImport(spec))
	if err != nil {
		t.Error(err)
		return nil, nil
	}
	pkg, err := NewPackage(imp)
	if err != nil {
		t.Error(err)
		return nil, nil
	}
	return pkg.findItf(name)
}

// ================================= TESTS =====================================

func Test_keys(t *testing.T) {
	// --- Given ---
	m := map[string]struct{}{"G": {}, "C": {}, "R": {}}

	// --- When ---
	have := keys(m)

	// --- Then ---
	assert.Equal(t, []string{"C", "G", "R"}, have)
}
