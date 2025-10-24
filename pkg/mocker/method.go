// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"fmt"
	"strings"
)

// method represents an interface method.
type method struct {
	name string     // Name of the method.
	args []argument // Zero or more method arguments.
	rets []argument // Zero or more method return values.
}

// generate generates code representing the method.
//
// Example:
//
//	func (_mck *CaseMock) Method13(tim mt.Time) error {
//		_rets := _mck.Called(tim)
//		return _rets.Error(0)
//	}
func (met *method) generate(recType string) string {
	code := met.genSig(recType, true)
	code += " {\n\t_mck.t.Helper()\n"
	code += met.genCalled()
	code += met.genRetCheck()
	retBody := met.genReturnBody(1)
	if retBody != "" {
		code += "\n" + retBody + "\n"
		code += met.genReturn()
	}
	code += "}"
	return code
}

// generateOn generates code for the method's "OnXXX" helper.
func (met *method) generateOn(typ string) string {
	code := met.genOnSig(typ)
	code += " {\n\t_mck.t.Helper()\n"
	code += met.genArgSlice()
	code += fmt.Sprintf("\treturn _mck.On(%q, _args...)\n", met.name)
	code += "}"
	return code
}

// genReceiver generates code for the method's receiver where "typ" represents
// the receiver type. Returns an empty string if typ is empty.
//
// Example:
//
//	(_mck *SimpleMock)
func (met *method) genReceiver(typ string) string {
	if typ == "" {
		return ""
	}
	return fmt.Sprintf("(_mck *%s)", typ)
}

// genArgs generates code representing method's arguments wrapped in
// parentheses.
//
// Examples:
//
//	()
//	(x int)
//	(x int, y bool)
func (met *method) genArgs() string {
	code := "("
	for i, arg := range met.args {
		if i > 0 {
			code += ", "
		}
		code += arg.genArg(i)
	}
	code += ")"
	return code
}

// genAnyArgs generates code representing method's arguments wrapped in
// parentheses where all types are replaced with "any" to allow usage of
// [mock.ArgMatcher]. Used in "OnXXX" helper methods.
//
// Examples:
//
//	()
//	(x any)
//	(x any, y any)
//	(x any, y any, z ...any)
func (met *method) genAnyArgs() string {
	code := "("
	for i, arg := range met.args {
		if i > 0 {
			code += ", "
		}
		code += arg.genName(i) + " "
		if arg.isVariadic() {
			code += "..."
		}
		code += "any"
	}
	code += ")"
	return code
}

// argNames returns method's argument names.
//
// Examples:
//
//	a
//	a, b
//	a, b, c
//	a, b...
func (met *method) argNames() []string {
	var names []string
	for i, arg := range met.args {
		name := arg.genName(i)
		if arg.isVariadic() {
			name += "..."
		}
		names = append(names, name)
	}
	return names
}

// genArgTypes returns method's argument types.
//
// Examples:
//
//	()
//	(int)
//	(int, bool)
//	(int, bool, *Concrete)
func (met *method) genArgTypes() string {
	var names []string
	for _, arg := range met.args {
		names = append(names, arg.typ)
	}
	return "(" + strings.Join(names, ", ") + ")"
}

// genRets generates code representing method's return arguments. Returns an
// empty string if the method has no return arguments.
//
// Examples:
//
//	error
//	(int, error)
func (met *method) genRets() string {
	var code string
	rc := len(met.rets)
	for i, ret := range met.rets {
		if rc > 1 && i == 0 {
			code += "("
		}
		if i > 0 {
			code += ", "
		}
		code += ret.getType()
		if rc > 1 && i == rc-1 {
			code += ")"
		}
	}
	return code
}

// genSig generates code for a method signature. If recType is non-empty, the
// method receiver is generated with the given type. If argNames is true,
// the function arguments are generated with names.
//
// Examples:
//
//	func (_mck *TypeMock) Method()
//	func (_mck *TypeMock) Method(a int) (int, error)
//	func (_mck *TypeMock) Method(int) (int, error)
//	func Method(int) (int, error)
func (met *method) genSig(recType string, argNames bool) string {
	code := "func"
	rcv := met.genReceiver(recType)
	if rcv != "" {
		code += " " + rcv
	}
	if met.name != "" {
		code += " " + met.name
	}
	if argNames {
		code += met.genArgs()
	} else {
		code += met.genArgTypes()
	}
	retCode := met.genRets()
	if retCode != "" {
		code += " " + retCode
	}
	return code
}

// genOnSig generates code for the method's "OnXXX" helper signature where
// method's argument types are replaced with "any" to allow usage of
// [mock.ArgMatcher], and the return value of [mock.Call].
//
// Examples:
//
//	func (_mck *TypeMock) OnMethod() *mock.Call
//	func (_mck *TypeMock) OnMethod(a ...any) *mock.Call
func (met *method) genOnSig(typ string) string {
	code := "func"
	rcv := met.genReceiver(typ)
	if rcv != "" {
		code += " " + rcv
	}
	if met.name != "" {
		code += " On" + met.name
	}
	code += met.genAnyArgs()
	code += " *mock.Call"
	return code
}

// genArgSlice generates code for building an "[]any" slice with all the
// method's arguments.
//
// Examples:
//
//	_args := []any{a}
//	for _, _arg := range b {
//	  _args = append(_args, _arg)
//	}
func (met *method) genArgSlice() string {
	code := "\tvar _args []any\n"
	args := met.args
	if met.isVariadic() {
		// Remove variadic argument (always the last one).
		args = met.args[0 : len(met.args)-1]
	}
	for i, arg := range args {
		if i == 0 {
			code = "\t_args := []any{"
		}
		code += arg.genName(i)
		if i == len(args)-1 {
			code += "}\n"
		} else {
			code += ", "
		}
	}
	if met.isVariadic() {
		idx := len(met.args) - 1
		if idx < 0 {
			idx = 0
		}
		name := met.args[idx].genName(idx)
		code += genAppendFromTo("_args", name)
	}
	return code
}

// genCalled generates code calling the mock "_mck.Called" method.
//
// Example:
//
//	 _args := []any{a, b}
//	 for _, _arg := range c {
//	     _args = append(_args, _arg)
//	 }
//	_rets := _mck.Called(_args...)
func (met *method) genCalled() string {
	code := met.genArgSlice()
	if len(met.rets) > 0 {
		code += "\t_rets := _mck.Called(_args...)\n"
	} else {
		code += "\t_mck.Called(_args...)\n"
	}
	return code
}

// genRetCheck generates code checking the number of return arguments matches
// the mocked function arguments count.
func (met *method) genRetCheck() string {
	var code string
	if len(met.rets) > 0 {
		msg := "the number of mocked method returns does not match"
		code += fmt.Sprintf("\tif len(_rets) != %d {\n", len(met.rets))
		code += fmt.Sprintf("\t\t_mck.t.Fatal(%q)\n", msg)
		code += "\t}\n"
	}
	return code
}

// genReturnBody generates code for mocked method return arguments.
//
// Example:
//
//	var _r0 pkg.Tx
//	if _r := _rets.Get(0); _r != nil {
//		_r0 = _r.(pkg.Tx)
//	}
//	_r1 := _rets.Error(1)
func (met *method) genReturnBody(indent int) string {
	if len(met.rets) == 0 {
		return ""
	}
	lines := make([]string, 0, len(met.rets))
	for i, ret := range met.rets {
		single := met.genSingleRet(i, ret, indent)
		lines = append(lines, single...)
	}
	return strings.Join(lines, "\n")
}

// genSingleRet generates code for mocked method single return argument.
//
// Example:
//
//	var _r0 pkg.Tx
//	if _rFn, ok := _rets.Get(0).(func(pkg.Tx) pkg.Tx); ok {
//		_r0 = _rFn(a)
//	} else if _r := _rets.Get(0).(pkg.Tx); _r != nil {
//		_r0 = _r.(pkg.Tx)
//	}
func (met *method) genSingleRet(idx int, ret argument, indent int) []string {
	ind := strings.Repeat("\t", indent)
	lines := make([]string, 0, 5)

	// Declare variable.
	lines = append(lines, fmt.Sprintf("%svar _r%d %s", ind, idx, ret.typ))

	// Check if Get returns a function.
	codeFnTyp := fmt.Sprintf("func%s %s", met.genArgTypes(), ret.typ)
	code := fmt.Sprintf(
		"%sif _rFn, ok := _rets.Get(%d).(%s); ok {",
		ind,
		idx,
		codeFnTyp,
	)
	lines = append(lines, code)

	// If function call.
	argNames := strings.Join(met.argNames(), ", ")
	code = fmt.Sprintf("%s\t_r%d = _rFn(%s)", ind, idx, argNames)
	lines = append(lines, code)

	// Else cast.
	code = fmt.Sprintf("%s} else if _r := _rets.Get(%d); _r != nil {", ind, idx)
	lines = append(lines, code)
	if ret.typ == "any" {
		code = fmt.Sprintf("%s\t_r%d = _r", ind, idx)
	} else {
		code = fmt.Sprintf("%s\t_r%d = _r.(%s)", ind, idx, ret.typ)
	}
	lines = append(lines, code, ind+"}")

	return lines
}

// genReturn generates code for the method return line. Returns an empty string
// if there are no return arguments.
//
// Example:
//
//	return _rets.Error(0)
//	return _rets.Int(0), _rets.Error(1)
func (met *method) genReturn() string {
	if len(met.rets) == 0 {
		return ""
	}
	code := make([]string, 0, len(met.rets))
	for i := range met.rets {
		code = append(code, fmt.Sprintf("_r%d", i))
	}
	return "\treturn " + strings.Join(code, ", ") + "\n"
}

// imports return a list of imports needed by the method's arguments and
// returns.
func (met *method) imports() []*gopkg {
	var imps []*gopkg
	for _, arg := range met.args {
		imps = addUniquePackage(imps, arg.pks...)
	}
	for _, ret := range met.rets {
		imps = addUniquePackage(imps, ret.pks...)
	}
	return imps
}

// variadic returns true when method has variadic arguments.
func (met *method) isVariadic() bool {
	for _, arg := range met.args {
		if arg.isVariadic() {
			return true
		}
	}
	return false
}
