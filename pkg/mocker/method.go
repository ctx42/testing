// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// maxLine is the widest a generated line may be before a signature is
// broken one argument per line.
const maxLine = 80

// method represents an interface method.
type method struct {
	name string     // Name of the method.
	args []argument // Zero or more method arguments.
	rets []argument // Zero or more method return values.
}

// generate generates the start of the method's code: the signature, the
// Helper call, the call to Called, and the check of the returns count.
//
// Example:
//
//	func (_mck *CaseMock) Method13(tim mt.Time) error {
//		_mck.t.Helper()
//		_args := []any{tim}
//		_rets := _mck.Called(_args...)
//		if len(_rets) != 1 {
//			_mck.t.Fatal("the number of mocked method returns does not match")
//		}
func (met *method) generate(recType string) string {
	head := "func " + met.genReceiver(recType) + " " + met.name + "("
	tail := ")"
	if rets := met.genRets(); rets != "" {
		tail += " " + rets
	}
	code := met.genDoc()
	code += wrapSig(head, met.genArgList(), tail)
	code += "\t_mck.t.Helper()\n"
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
	format := "// On%s expects a call of [%s.%s].\n"
	if len(met.args) > 0 {
		format = "// On%s expects a call of [%s.%s] with the arguments.\n"
	}
	code := fmt.Sprintf(format, met.name, typ, met.name)
	head := "func " + met.genReceiver(typ) + " On" + met.name + "("
	code += wrapSig(head, met.genAnyArgList(), ") *mock.Call")
	code += "\t_mck.t.Helper()\n"
	code += met.genArgSlice()
	code += fmt.Sprintf("\treturn _mck.On(%q, _args...)\n", met.name)
	code += "}"
	return code
}

// genDoc generates the method's godoc line.
//
// Examples:
//
//	// Method records a call of Method.
//	// Method records a call of Method and returns its set results.
func (met *method) genDoc() string {
	if len(met.rets) == 0 {
		return fmt.Sprintf("// %s records a call of %[1]s.\n", met.name)
	}
	format := "// %s records a call of %[1]s and returns its set results.\n"
	return fmt.Sprintf(format, met.name)
}

// wrapSig returns the signature made of the head up to the opening
// parenthesis, the args, and the tail from the closing parenthesis,
// opening a method's body: on one line when it fits [maxLine], else with
// one argument per line and a blank line opening the body.
func wrapSig(head string, args []string, tail string) string {
	sig := head + strings.Join(args, ", ") + tail
	if utf8.RuneCountInString(sig+" {") <= maxLine || len(args) == 0 {
		return sig + " {\n"
	}
	return head + "\n\t" + strings.Join(args, ",\n\t") + ",\n" +
		tail + " {\n\n"
}

// genArgList returns the method's arguments, each with its name and type.
//
// Example:
//
//	[]string{"x int", "y ...bool"}
func (met *method) genArgList() []string {
	out := make([]string, 0, len(met.args))
	for i, arg := range met.args {
		out = append(out, arg.genName(i)+" "+arg.typ)
	}
	return out
}

// genAnyArgList returns the method's arguments with "any" types, as
// [method.genAnyArgs] writes them.
//
// Example:
//
//	[]string{"x any", "y ...any"}
func (met *method) genAnyArgList() []string {
	out := make([]string, 0, len(met.args))
	for i, arg := range met.args {
		typ := "any"
		if arg.isVariadic() {
			typ = "...any"
		}
		out = append(out, arg.genName(i)+" "+typ)
	}
	return out
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
	var code strings.Builder
	_, _ = code.WriteString("(")
	for i, arg := range met.args {
		if i > 0 {
			_, _ = code.WriteString(", ")
		}
		_, _ = code.WriteString(arg.genName(i) + " " + arg.typ)
	}
	_, _ = code.WriteString(")")
	return code.String()
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
	var code strings.Builder
	_, _ = code.WriteString("(")
	for i, arg := range met.args {
		if i > 0 {
			_, _ = code.WriteString(", ")
		}
		_, _ = code.WriteString(arg.genName(i) + " ")
		if arg.isVariadic() {
			_, _ = code.WriteString("...")
		}
		_, _ = code.WriteString("any")
	}
	_, _ = code.WriteString(")")
	return code.String()
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
	var code strings.Builder
	rc := len(met.rets)
	for i, ret := range met.rets {
		if rc > 1 && i == 0 {
			_, _ = code.WriteString("(")
		}
		if i > 0 {
			_, _ = code.WriteString(", ")
		}
		_, _ = code.WriteString(ret.typ)
		if rc > 1 && i == rc-1 {
			_, _ = code.WriteString(")")
		}
	}
	return code.String()
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
//	for _, _elem := range b {
//		_args = append(_args, _elem)
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
		idx := max(len(met.args)-1, 0)
		name := met.args[idx].genName(idx)
		code += genAppendFromTo("_args", name)
	}
	return code
}

// genCalled generates code calling the mock "_mck.Called" method.
//
// Example:
//
//	_args := []any{a, b}
//	for _, _elem := range c {
//		_args = append(_args, _elem)
//	}
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

// genReturnBody generates code for mocked method return arguments, one
// [method.genSingleRet] block each.
//
// Example:
//
//	var _r0 pkg.Tx
//	switch _v := _rets.Get(0).(type) {
//	case func() pkg.Tx:
//		_r0 = _v()
//	case pkg.Tx:
//		_r0 = _v
//	case nil:
//	default:
//		_mck.t.Fatalf("Method return 0 is %T, not pkg.Tx", _v)
//	}
//	var _r1 error
//	switch _v := _rets.Get(1).(type) {
//	...
//	}
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

// genSingleRet generates code for mocked method single return argument: a
// function of the method's arguments is called, a value of the type taken,
// a nil left as the zero value, and any other value fails the test.
//
// Example:
//
//	var _r0 pkg.Tx
//	switch _v := _rets.Get(0).(type) {
//	case func(pkg.Tx) pkg.Tx:
//		_r0 = _v(a)
//	case pkg.Tx:
//		_r0 = _v
//	case nil:
//	default:
//		_mck.t.Fatalf("Method return 0 is %T, not pkg.Tx", _v)
//	}
func (met *method) genSingleRet(idx int, ret argument, indent int) []string {
	ind := strings.Repeat("\t", indent)
	argNames := strings.Join(met.argNames(), ", ")
	codeFnTyp := fmt.Sprintf("func%s %s", met.genArgTypes(), ret.typ)
	lines := []string{
		fmt.Sprintf("%svar _r%d %s", ind, idx, ret.typ),
		fmt.Sprintf("%sswitch _v := _rets.Get(%d).(type) {", ind, idx),
		fmt.Sprintf("%scase %s:", ind, codeFnTyp),
		fmt.Sprintf("%s\t_r%d = _v(%s)", ind, idx, argNames),
	}
	if ret.typ == "any" {
		return append(lines,
			ind+"default:",
			fmt.Sprintf("%s\t_r%d = _v", ind, idx),
			ind+"}",
		)
	}
	// The message is a Go string literal holding a format: quote it, and
	// keep a "%" of the type text from reading as a verb.
	typ := strings.ReplaceAll(ret.typ, "%", "%%")
	msg := strconv.Quote(fmt.Sprintf("%s return %d is %%T, not %s",
		met.name, idx, typ))
	return append(lines,
		fmt.Sprintf("%scase %s:", ind, ret.typ),
		fmt.Sprintf("%s\t_r%d = _v", ind, idx),
		ind+"case nil:",
		ind+"default:",
		fmt.Sprintf("%s\t_mck.t.Fatalf(%s, _v)", ind, msg),
		ind+"}",
	)
}

// genReturn generates code for the method return line. Returns an empty string
// if there are no return arguments.
//
// Example:
//
//	return _r0
//	return _r0, _r1
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

// imports returns a list of imports needed by the method's arguments and
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

// isVariadic returns true when the method has variadic arguments.
func (met *method) isVariadic() bool {
	for _, arg := range met.args {
		if arg.isVariadic() {
			return true
		}
	}
	return false
}
