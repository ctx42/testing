// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

func ExampleNew_formattedHeader() {
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
	err := notice.New("expected values to be equal").
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

func ExampleNotice_Append_forceNextLine() {
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
	errMy := errors.New("my error")

	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		Wrap(errMy)

	is := errors.Is(msg, errMy)
	fmt.Println(is)
	// Output: true
}

func ExampleNotice_Unwrap() {
	errMy := errors.New("my error")

	msg := notice.New("expected values to be equal").
		Want("%s", "abc").
		Have("%s", "xyz").
		Wrap(errMy)

	fmt.Println(msg.Unwrap())
	// Output: my error
}

func ExampleNotice_Is() {
	first := notice.New("first failure")
	joined := notice.Join(first, notice.New("second failure"))

	fmt.Println(errors.Is(joined, first))
	fmt.Println(errors.Is(joined, notice.ErrNotice))
	// Output:
	// true
	// true
}

func ExampleNotice_Chain() {
	first := notice.New("first failure")
	second := notice.New("second failure").Chain(first)

	fmt.Println(second.Prev().HeaderText)
	fmt.Println(first.Next().HeaderText)
	// Output:
	// first failure
	// second failure
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

func ExampleJoin() {
	err0 := notice.New("expected values to be equal").
		Want("%d", 1).
		Have("%d", 2)
	err1 := notice.New("expected value to be true").Have("%v", false)

	err := notice.Join(err0, err1)

	fmt.Println(err)
	// Output:
	// multiple expectations violated:
	//   error: expected values to be equal
	//    want: 1
	//    have: 2
	//       ---
	//   error: expected value to be true
	//    have: false
}

func ExampleSortNotices() {
	msgB := notice.New("second").SetTrail("B")
	msgA := notice.New("first").SetTrail("A")
	_ = notice.Join(msgB, msgA)

	tail := notice.SortNotices(msgB, notice.TrailCmp)

	for _, msg := range tail.All() {
		fmt.Println(msg.Trail, msg.HeaderText)
	}
	// Output:
	// A first
	// B second
}
