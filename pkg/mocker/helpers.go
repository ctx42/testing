package mocker

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// assumedPackageName returns the assumed package name of an import path.
// It does this using only string parsing of the import path. It picks the last
// element of the path that does not look like a major version and then picks
// the valid identifier at the start of that element.
//
// Copied from: https://github.com/golang/tools/blob/a318c19ff2fd8d6aae74e36fe7e1a8b8afef3bf7/internal/imports/fix.go#L1233
//
// Example:
//
//	github.com/user/project/pkg/package -> package
//	github.com/user/project/pkg/go_package -> go_package
//	github.com/user/project/pkg/go-package-abc -> abc
func assumedPackageName(pth string) string {
	notIdentifier := func(ch rune) bool {
		return !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
			'0' <= ch && ch <= '9' ||
			ch == '_' ||
			ch >= utf8.RuneSelf &&
				(unicode.IsLetter(ch) || unicode.IsDigit(ch)))
	}

	base := path.Base(pth)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			dir := path.Dir(pth)
			if dir != "." {
				base = path.Base(dir)
			}
		}
	}
	parts := strings.Split(base, "-")
	base = parts[len(parts)-1]
	if i := strings.IndexFunc(base, notIdentifier); i >= 0 {
		base = base[:i]
	}
	return base
}

// genAppendFromTo generates code that appends all elements from the "from"
// slice to the "to" slice. It returns a string containing Go code with a
// for-range loop that iterates over the "from" slice and appends each element
// to the "to" slice.
//
// For example, `genAppendFromTo("dst", "src")` call returns:
//
//	for _, _elem := range src {
//	    dst = append(dst, _elem)
//	}
func genAppendFromTo(to, from string) string {
	code := fmt.Sprintf("\tfor _, _elem := range %s {\n", from)
	code += fmt.Sprintf("\t\t%[1]s = append(%[1]s, _elem)\n", to)
	code += "\t}\n"
	return code
}

// addUniquePackage appends a package to the dst slice only if it's not already
// present. Packages are considered equal if their import paths are equal.
func addUniquePackage(dst []*gopkg, src ...*gopkg) []*gopkg {
next:
	for _, imp := range src {
		for _, have := range dst {
			if have.pkgPath == imp.pkgPath {
				continue next
			}
		}
		dst = append(dst, imp)
	}
	return dst
}

// sortImports organizes import statements into three sorted groups: standard
// library imports are sorted and listed first, followed by other sorted
// imports, and finally sorted dot imports.
//
// Example output:
//
//	"fmt"
//	"net/http"
//	mt "time"
//
//	"github.com/tst/pkga"
//	"github.com/tst/pkgb"
//	"github.com/tst/pkgc"
//
//	. "github.com/tst/pkgd"
func sortImports(ips []*gopkg) []*gopkg {
	fn := func(i, j int) bool {
		iStd := strings.Index(ips[i].pkgPath, ".") == -1
		jStd := strings.Index(ips[j].pkgPath, ".") == -1
		if iStd != jStd {
			return iStd == true
		}
		iAn := assumedPackageName(ips[i].pkgPath)
		jAn := assumedPackageName(ips[j].pkgPath)
		return iAn < jAn
	}
	sort.SliceStable(ips, fn)
	for i := 1; i < len(ips); i++ {
		prevStd := strings.Index(ips[i-1].pkgPath, ".") == -1
		currStd := strings.Index(ips[i].pkgPath, ".") == -1
		if ips[i-1].pkgPath != "" && prevStd && !currStd {
			ips = slices.Insert(ips, i, &gopkg{})
			continue
		}

		prevDot := ips[i-1].alias == "."
		currDot := ips[i].alias == "."
		if ips[i-1].pkgPath != "" && !prevDot && currDot {
			ips = slices.Insert(ips, i, &gopkg{})
			continue
		}
	}
	return ips
}

// genImports generates and returns Go code representing interface imports.
func genImports(imps []*gopkg) string {
	if len(imps) == 0 {
		return ""
	}
	ips := sortImports(imps)

	buf := &bytes.Buffer{}
	buf.WriteString("import (\n")
	for i, imp := range ips {
		str := imp.GoString()
		if str == "" && i > 0 {
			buf.WriteString("\n")
		} else {
			buf.WriteString("\t")
			buf.WriteString(str)
			buf.WriteString("\n")
		}
	}
	buf.WriteString(")")
	return buf.String()
}

// toLowerSnakeCase converts camel case to lowercase snake case.
func toLowerSnakeCase(camel string) string {
	var runes = make([]rune, 0, len(camel)+10)
	const lower = 1
	const upper = 2

	var prev int
	for i := 0; i < len(camel); i++ {
		curr := lower
		r := rune(camel[i])
		if unicode.IsUpper(r) {
			curr = upper
		}
		if prev == lower && curr == upper {
			runes = append(runes, '_', r)
		} else {
			runes = append(runes, r)
		}
		prev = curr
	}
	return strings.ToLower(string(runes))
}

// findSources returns a list of paths to all Go source files (excluding test
// files) in the specified directory. It does not recurse into subdirectories.
// The returned paths are absolute.
func findSources(dir string) ([]string, error) {
	var err error
	if dir, err = filepath.Abs(dir); err != nil {
		return nil, err
	}

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	ets, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var sources []string
	for _, entry := range ets {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		sources = append(sources, filepath.Join(dir, name))
	}
	sort.Strings(sources)
	return sources, nil
}
