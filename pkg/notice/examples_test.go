// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package notice_test

import (
	"errors"
	"fmt"

	"github.com/ctx42/testing/pkg/notice"
)

func ExampleNew() {
	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz")

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   want: abc
	//   have: xyz
}

func ExampleNew_formatedHeader() {
	msg := notice.New("expected %s to be equal", "values").
		Want("%s", "abc").
		Have("%s", "xyz")

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   want: abc
	//   have: xyz
}

func ExampleFrom() {
	var err error
	err = notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz")

	msg := notice.From(err, "optional prefix").
		Append("my", "%s", "value")

	fmt.Println(msg)
	// Output:
	// [optional prefix] expected values to be equal:
	//   want: abc
	//   have: xyz
	//     my: value
}

func ExampleNotice_SetHeader() {
	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz")

	_ = msg.SetHeader("some other %s", "header")

	fmt.Println(msg)
	// Output:
	// some other header:
	//   want: abc
	//   have: xyz
}

func ExampleNotice_Append() {
	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		Append("name", "%d", 5)

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   want: abc
	//   have: xyz
	//   name: 5
}

func ExampleNotice_Append_multiLine() {
	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "x\ny\nz").
		Append("name", "%d", 5)

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   want: abc
	//   have:
	//         x
	//         y
	//         z
	//   name: 5
}

func ExampleNotice_Append_forceNexLine() {
	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("\n%s", "xyz").
		Append("name", "%d", 5)

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   want: abc
	//   have:
	//         xyz
	//   name: 5
}

func ExampleNotice_AppendRow() {
	row0 := notice.NewRow("number", "%d", 5)
	row1 := notice.NewRow("string", "%s", "abc")

	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		AppendRow(row0, row1)

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//     want: abc
	//     have: xyz
	//   number: 5
	//   string: abc
}

func ExampleNotice_Prepend() {
	msg := notice.New("expected values to be equal").
		SetTrail("type.field").
		Want("%s", "abc").
		Have("%s", "xyz").
		Prepend("name", "%d", 5)

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   trail: type.field
	//    name: 5
	//    want: abc
	//    have: xyz
}

func ExampleNotice_SetTrail() {
	msg := notice.New("expected values to be equal").
		SetTrail("type.field").
		Want("%s", "abc").
		Have("%s", "xyz")

	fmt.Println(msg)
	// Output:
	// expected values to be equal:
	//   trail: type.field
	//    want: abc
	//    have: xyz
}

func ExampleNotice_Wrap() {
	ErrMy := errors.New("my error")

	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		Wrap(ErrMy)

	is := errors.Is(msg, ErrMy)
	fmt.Println(is)
	// Output: true
}

func ExampleNotice_Unwrap() {
	ErrMy := errors.New("my error")

	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		Wrap(ErrMy)

	fmt.Println(msg.Unwrap())
	// Output: my error
}

func ExampleNotice_MetaSet() {
	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		MetaSet("key", 123)

	fmt.Println(msg.MetaLookup("key"))
	// Output: 123 true
}

func ExampleIndent() {
	lines := notice.Indent(4, ' ', "line1\nline2\nline3")

	fmt.Println(lines)
	// Output:
	//     line1
	//     line2
	//     line3
}
