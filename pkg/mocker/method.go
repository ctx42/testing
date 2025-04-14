// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// Method represents an interface method function argument type.
type Method struct {
	name string     // Name which may be empty for argument types.
	args []Argument // Zero or more method arguments.
	rets []Argument // Zero or more method return values.
}

func NewMethod(name string, args, rets []Argument) *Method {
	// TODO(rz): test this.
	// TODO(rz): document this.
	return &Method{
		name: name,
		args: args,
		rets: rets,
	}
}
