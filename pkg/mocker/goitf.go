// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// Imports added to every generated mock.
const (
	selfImp   = "github.com/ctx42/testing/pkg/mock"
	testerImp = "github.com/ctx42/testing/pkg/tester"
)

// goitf represents an interface.
type goitf struct {
	name    string    // The interface name.
	methods []*method // The interface methods.
}

// find returns the interface method by the name, or [ErrUnkMet] if not found.
func (itf *goitf) find(name string) (*method, error) {
	for _, met := range itf.methods {
		if name == met.name {
			return met, nil
		}
	}
	return nil, ErrUnkMet
}

// generate generates code for the interface mock. When onHelpers is true, the
// OnXXX helper methods are also generated.
func (itf *goitf) generate(recType string, onHelpers bool) string {
	var code string
	for i, met := range itf.methods {
		code += met.generate(recType)
		if onHelpers {
			code += "\n\n" + met.generateOn(recType)
		}
		if i < len(itf.methods)-1 {
			code += "\n\n"
		}
	}
	return code
}

// imports returns unique imports used by all the interface methods in
// arguments and return values.
func (itf *goitf) imports() []*gopkg {
	var imps []*gopkg
	for _, met := range itf.methods {
		imps = addUniquePackage(imps, met.imports()...)
	}
	return imps
}

// genImports generates and returns Go code representing interface imports.
func (itf *goitf) genImports() string {
	mckImp := &gopkg{pkgName: assumedPackageName(selfImp), pkgPath: selfImp}
	tstImp := &gopkg{pkgName: assumedPackageName(testerImp), pkgPath: testerImp}
	imps := append(itf.imports(), mckImp, tstImp)
	return genImports(imps)
}
