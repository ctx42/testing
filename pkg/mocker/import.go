// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// Import represents Go import spec.
//
// Examples:
//
//	"github.com/user/project/pkg/package"
//	alias "github.com/user/project/pkg/package"
//	. "github.com/user/project/pkg/package"
type Import struct {
	Alias string // Alias (set to "." for dot imports).
	Name  string // Name used in expressions (always set based on spec).
	Spec  string // Import itself.
	Dir   string // Absolute directory path to the code import represents.
}

// NewImport returns a new instance of [Import] with the given import spec.
func NewImport(spec string) Import {
	return Import{}.SetSpec(spec)
}

// SetSpec sets import spec.
func (imp Import) SetSpec(spec string) Import {
	imp.Name = assumedPackageName(spec)
	imp.Spec = spec
	return imp
}

// SetAlias sets import alias.
func (imp Import) SetAlias(alias string) Import {
	imp.Alias = alias
	return imp
}

// SetDir sets the directory path to the code import represents.
func (imp Import) SetDir(dir string) Import {
	imp.Dir = dir
	return imp
}

// IsDot returns true if the [Import] represents a dot import.
func (imp Import) IsDot() bool { return imp.Alias == "." }

// IsZero returns true when the [Import] is considered to hold zero value.
func (imp Import) IsZero() bool { return imp.Spec == "" && imp.Dir == "" }

// GoString generates and returns Go code representing the import line.
// If the spec is empty, it returns an empty string.
//
// Example:
//
//	"fmt"
//	"github.com/user/project/pkg/package"
//	alias "github.com/user/project/pkg/package"
//	. "github.com/user/project/pkg/package"
func (imp Import) GoString() string {
	if imp.Spec == "" {
		return ""
	}
	code := `"` + imp.Spec + `"`
	if imp.Alias != "" {
		code = imp.Alias + " " + code
	}
	return code
}
