// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump_test

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/dump"
)

func ExampleDump_Any() {
	val := types.TA{
		Dur: 3,
		Int: 42,
		Loc: types.WAW,
		Str: "abc",
		Tim: time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
		TAp: nil,
	}

	have := dump.New().Any(val)

	fmt.Println(have)
	// Output:
	// {
	//   Int: 42,
	//   Str: "abc",
	//   Tim: "2000-01-02T03:04:05Z",
	//   Dur: "3ns",
	//   Loc: "Europe/Warsaw",
	//   TAp: nil,
	//   private: 0,
	// }
}

func ExampleDump_Any_flatCompact() {
	val := map[string]any{
		"int": 42,
		"loc": types.WAW,
		"nil": nil,
	}

	have := dump.New(dump.WithFlat).Any(val)

	fmt.Println(have)
	// Output:
	// map[string]any{"int": 42, "loc": "Europe/Warsaw", "nil": nil}
}

func ExampleDump_Any_customTimeFormat() {
	val := map[time.Time]int{time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC): 42}

	have := dump.New(dump.WithFlat, dump.WithTimeFormat(time.Kitchen)).Any(val)

	fmt.Println(have)
	// Output:
	// map[time.Time]int{"3:04AM": 42}
}

func ExampleDump_Any_customDumper() {
	var i int
	dumper := func(dmp dump.Dump, lvl int, val reflect.Value) string {
		switch val.Kind() {
		case reflect.Int:
			return fmt.Sprintf("%X", val.Int())
		default:
			panic("unexpected kind")
		}
	}
	opts := []dump.Option{
		dump.WithFlat,
		dump.WithCompact,
		dump.WithDumper(i, dumper),
	}

	have := dump.New(opts...).Any(42)

	fmt.Println(have)
	// Output:
	// 2A
}

func ExampleDump_Any_recursive() {
	type Node struct {
		Value    int
		Children []*Node
	}

	val := &Node{
		Value: 1,
		Children: []*Node{
			{
				Value:    2,
				Children: nil,
			},
			{
				Value: 3,
				Children: []*Node{
					{
						Value:    4,
						Children: nil,
					},
				},
			},
		},
	}

	have := dump.New().Any(val)
	fmt.Println(have)
	// Output:
	// {
	//   Value: 1,
	//   Children: []*dump_test.Node{
	//     {
	//       Value: 2,
	//       Children: nil,
	//     },
	//     {
	//       Value: 3,
	//       Children: []*dump_test.Node{
	//         {
	//           Value: 4,
	//           Children: nil,
	//         },
	//       },
	//     },
	//   },
	// }
}
