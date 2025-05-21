// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/internal/diff"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/goldy"
)

func Test_method_generate(t *testing.T) {
	t.Run("without args and without returns", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_method/without_args_without_rets.gld"
		met := &method{
			name: "Method",
			args: nil,
			rets: nil,
		}

		// --- When ---
		have := met.generate("MyMock")

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})

	t.Run("with single arg without returns", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_method/with_arg_without_rets.gld"
		met := &method{
			name: "Method",
			args: []argument{
				{name: "a", typ: "int"},
			},
			rets: nil,
		}

		// --- When ---
		have := met.generate("MyMock")

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})

	t.Run("without args with returns", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_method/without_arg_with_ret.gld"
		met := &method{
			name: "Method",
			args: nil,
			rets: []argument{
				{typ: "error"},
			},
		}

		// --- When ---
		have := met.generate("MyMock")

		// --- Then ---
		want := goldy.Open(t, gfp).String()

		// TODO(rz):
		edits := diff.Strings(want, have)
		u, err := diff.ToUnified("have", "want", want, edits, 1)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(u)

		// assert.Equal(t, want, have)
	})
}

func Test_method_generateOn(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_on_method/without_args.gld"
		met := &method{
			name: "Method",
			args: nil,
		}

		// --- When ---
		have := met.generateOn("MyMock")

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})

	t.Run("with args", func(t *testing.T) {
		// --- Given ---
		gfp := "testdata/golden_on_method/with_args.gld"
		met := &method{
			name: "Method",
			args: []argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "any"},
			},
		}

		// --- When ---
		have := met.generateOn("MyMock")

		// --- Then ---
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
	})
}

func Test_method_genReceiver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		met := &method{name: "Method0"}

		// --- When ---
		have := met.genReceiver("TypeMock")

		// --- Then ---
		assert.Equal(t, "(_mck *TypeMock)", have)
	})

	t.Run("empty type and name", func(t *testing.T) {
		// --- Given ---
		met := &method{name: "Method0"}

		// --- When ---
		have := met.genReceiver("")

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("empty type and empty name", func(t *testing.T) {
		// --- Given ---
		met := &method{name: ""}

		// --- When ---
		have := met.genReceiver("")

		// --- Then ---
		assert.Equal(t, "", have)
	})
}

func Test_method_genArgs_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		want string
	}{
		{
			"nil args",
			nil,
			"()",
		},
		{
			"no args",
			[]argument{},
			"()",
		},
		{
			"one named argument",
			[]argument{
				{name: "a", typ: "int"},
			},
			"(a int)",
		},
		{
			"two named arguments",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
			},
			"(a int, b bool)",
		},
		{
			"one argument without a name",
			[]argument{
				{name: "", typ: "int"},
			},
			"(_a0 int)",
		},
		{
			"two arguments without a name",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "bool"},
			},
			"(_a0 int, _a1 bool)",
		},
		{
			"variadic",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "...bool"},
			},
			"(_a0 int, _a1 ...bool)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{args: tc.args}

			// --- When ---
			have := met.genArgs()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genAnyArgs_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		want string
	}{
		{
			"nil args",
			nil,
			"()",
		},
		{
			"no args",
			[]argument{},
			"()",
		},
		{
			"one named argument",
			[]argument{
				{name: "a", typ: "int"},
			},
			"(a any)",
		},
		{
			"two named arguments",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
			},
			"(a any, b any)",
		},
		{
			"one argument without a name",
			[]argument{
				{name: "", typ: "int"},
			},
			"(_a0 any)",
		},
		{
			"two arguments without a name",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "bool"},
			},
			"(_a0 any, _a1 any)",
		},
		{
			"variadic",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "...bool"},
			},
			"(_a0 any, _a1 ...any)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{args: tc.args}

			// --- When ---
			have := met.genAnyArgs()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_argNames_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		want []string
	}{
		{
			"nil args",
			nil,
			nil,
		},
		{
			"no args",
			[]argument{},
			nil,
		},
		{
			"one named argument",
			[]argument{
				{name: "a", typ: "int"},
			},
			[]string{"a"},
		},
		{
			"two named arguments",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
			},
			[]string{"a", "b"},
		},
		{
			"one argument without a name",
			[]argument{
				{name: "", typ: "int"},
			},
			[]string{"_a0"},
		},
		{
			"two arguments without a name",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "bool"},
			},
			[]string{"_a0", "_a1"},
		},
		{
			"two arguments without a name",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "bool"},
			},
			[]string{"_a0", "_a1"},
		},
		{
			"variadic",
			[]argument{
				{name: "", typ: "...int"},
				{name: "", typ: "bool"},
			},
			[]string{"_a0...", "_a1"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{args: tc.args}

			// --- When ---
			have := met.argNames()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genArgTypes_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		want string
	}{
		{
			"nil args",
			nil,
			"()",
		},
		{
			"no args",
			[]argument{},
			"()",
		},
		{
			"one named argument",
			[]argument{
				{name: "a", typ: "int"},
			},
			"(int)",
		},
		{
			"two named arguments",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
			},
			"(int, bool)",
		},
		{
			"one argument without a name",
			[]argument{
				{name: "", typ: "int"},
			},
			"(int)",
		},
		{
			"two arguments without a name",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "bool"},
			},
			"(int, bool)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{args: tc.args}

			// --- When ---
			have := met.genArgTypes()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genRets_tabular(t *testing.T) {
	tt := []struct {
		testN string

		rets []argument
		want string
	}{
		{"nil returns", nil, ""},
		{"no returns", []argument{}, ""},
		{
			"single return without a name",
			[]argument{
				{name: "", typ: "error"},
			},
			"error",
		},
		{
			"single return with name",
			[]argument{
				{name: "err", typ: "error"},
			},
			"error",
		},
		{
			"two returns without names",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			"(int, error)",
		},
		{
			"two returns with names",
			[]argument{
				{name: "err", typ: "int"},
				{name: "i", typ: "error"},
			},
			"(int, error)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{rets: tc.rets}

			// --- When ---
			have := met.genRets()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genSig_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ   string
		name  string
		args  []argument
		rets  []argument
		names bool
		want  string
	}{
		{
			"no arguments no returns",
			"TypeMock",
			"Method",
			nil,
			nil,
			true,
			"func (_mck *TypeMock) Method()",
		},
		{
			"no arguments one return",
			"TypeMock",
			"Method",
			nil,
			[]argument{
				{name: "", typ: "error"},
			},
			true,
			"func (_mck *TypeMock) Method() error",
		},
		{
			"no arguments two returns",
			"TypeMock",
			"Method",
			nil,
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			true,
			"func (_mck *TypeMock) Method() (int, error)",
		},
		{
			"one argument two returns",
			"TypeMock",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			true,
			"func (_mck *TypeMock) Method(a int) (int, error)",
		},
		{
			"one argument no returns",
			"TypeMock",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			nil,
			true,
			"func (_mck *TypeMock) Method(a int)",
		},
		{
			"no receiver type",
			"",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			nil,
			true,
			"func Method(a int)",
		},
		{
			"no receiver type and no argument types",
			"",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "float64"},
			},
			nil,
			false,
			"func Method(int, float64)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{
				name: tc.name,
				args: tc.args,
				rets: tc.rets,
			}

			// --- When ---
			have := met.genSig(tc.typ, tc.names)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genOnSig_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  string
		name string
		args []argument
		rets []argument
		want string
	}{
		{
			"no arguments no returns",
			"TypeMock",
			"Method",
			nil,
			nil,
			"func (_mck *TypeMock) OnMethod() *mock.Call",
		},
		{
			"no arguments one return",
			"TypeMock",
			"Method",
			nil,
			[]argument{
				{name: "", typ: "error"},
			},
			"func (_mck *TypeMock) OnMethod() *mock.Call",
		},
		{
			"no arguments two returns",
			"TypeMock",
			"Method",
			nil,
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			"func (_mck *TypeMock) OnMethod() *mock.Call",
		},
		{
			"one argument two returns",
			"TypeMock",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			"func (_mck *TypeMock) OnMethod(a any) *mock.Call",
		},
		{
			"one argument no returns",
			"TypeMock",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			nil,
			"func (_mck *TypeMock) OnMethod(a any) *mock.Call",
		},
		{
			"variadic",
			"TypeMock",
			"Method",
			[]argument{
				{name: "a", typ: "...int"},
			},
			nil,
			"func (_mck *TypeMock) OnMethod(a ...any) *mock.Call",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{
				name: tc.name,
				args: tc.args,
				rets: tc.rets,
			}

			// --- When ---
			have := met.genOnSig(tc.typ)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genArgSlice_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		want string
	}{
		{
			"no arguments",
			nil,
			"\tvar _args []any\n",
		},
		{
			"one argument",
			[]argument{
				{name: "a", typ: "int"},
			},
			"\t_args := []any{a}\n",
		},
		{
			"two arguments",
			[]argument{
				{name: "", typ: "int"},
				{name: "b", typ: "string"},
			},
			"\t_args := []any{_a0, b}\n",
		},
		{
			"single variadic argument",
			[]argument{
				{name: "a", typ: "...int"},
			},
			"" +
				"\tvar _args []any\n" +
				"\tfor _, _elem := range a {\n" +
				"\t\t_args = append(_args, _elem)\n" +
				"\t}\n",
		},
		{
			"argument and variadic argument",
			[]argument{
				{name: "a", typ: "string"},
				{name: "b", typ: "...int"},
			},
			"" +
				"\t_args := []any{a}\n" +
				"\tfor _, _elem := range b {\n" +
				"\t\t_args = append(_args, _elem)\n" +
				"\t}\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{
				args: tc.args,
			}

			// --- When ---
			have := met.genArgSlice()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genCalled_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		rets []argument
		want string
	}{
		{
			"no arguments no returns",
			nil,
			nil,
			"\tvar _args []any\n\t_mck.Called(_args...)\n",
		},
		{
			"no arguments with returns",
			nil,
			[]argument{
				{name: "", typ: "error"},
			},
			"\tvar _args []any\n\t_rets := _mck.Called(_args...)\n",
		},
		{
			"one argument with returns",
			[]argument{
				{name: "a", typ: "int"},
			},
			[]argument{
				{name: "", typ: "int"},
			},
			"\t_args := []any{a}\n\t_rets := _mck.Called(_args...)\n",
		},
		{
			"one argument no returns",
			[]argument{
				{name: "a", typ: "int"},
			},
			nil,
			"\t_args := []any{a}\n\t_mck.Called(_args...)\n",
		},
		{
			"two arguments no one return",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "int"},
			},
			[]argument{
				{name: "", typ: "error"},
			},
			"\t_args := []any{a, b}\n\t_rets := _mck.Called(_args...)\n",
		},
		{
			"one variadic argument no returns",
			[]argument{
				{name: "a", typ: "...int"},
			},
			nil,
			"" +
				"\tvar _args []any\n" +
				"\tfor _, _elem := range a {\n" +
				"\t\t_args = append(_args, _elem)\n" +
				"\t}\n" +
				"\t_mck.Called(_args...)\n",
		},
		{
			"multiple arguments with variadic argument no returns",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
				{name: "c", typ: "...int"},
			},
			nil,
			"" +
				"\t_args := []any{a, b}\n" +
				"\tfor _, _elem := range c {\n" +
				"\t\t_args = append(_args, _elem)\n" +
				"\t}\n" +
				"\t_mck.Called(_args...)\n",
		},
		{
			"multiple arguments with variadic argument with returns",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
				{name: "c", typ: "...int"},
			},
			[]argument{
				{name: "", typ: "error"},
			},
			"" +
				"\t_args := []any{a, b}\n" +
				"\tfor _, _elem := range c {\n" +
				"\t\t_args = append(_args, _elem)\n" +
				"\t}\n" +
				"\t_rets := _mck.Called(_args...)\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{
				args: tc.args,
				rets: tc.rets,
			}

			// --- When ---
			have := met.genCalled()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genRetCheck_tabular(t *testing.T) {
	tt := []struct {
		testN string

		rets []argument
		want string
	}{
		{
			"no returns",
			nil,
			"",
		},
		{
			"one return",
			[]argument{
				{name: "", typ: "error"},
			},
			"\tif len(_rets) != 1 {\n" +
				"\t\t_mck.t.Fatal(\"the number of mocked method " +
				"returns does not match\")\n" +
				"\t}\n",
		},
		{
			"two returns",
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			"\tif len(_rets) != 2 {\n" +
				"\t\t_mck.t.Fatal(\"the number of mocked method " +
				"returns does not match\")\n" +
				"\t}\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{rets: tc.rets}

			// --- When ---
			have := met.genRetCheck()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genReturn_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		rets []argument
		want string
	}{
		{
			"no returns",
			nil,
			nil,
			"",
		},
		{
			"one return",
			nil,
			[]argument{
				{name: "", typ: "error"},
			},
			"\treturn _r0\n",
		},
		{
			"two returns",
			nil,
			[]argument{
				{name: "", typ: "int"},
				{name: "", typ: "error"},
			},
			"\treturn _r0, _r1\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{
				args: tc.args,
				rets: tc.rets,
			}

			// --- When ---
			have := met.genReturn()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_imports(t *testing.T) {
	t.Run("no imports", func(t *testing.T) {
		// --- Given ---
		met := &method{}

		// --- When ---
		have := met.imports()

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("arg imports", func(t *testing.T) {
		// --- Given ---
		met := &method{
			args: []argument{
				{
					name: "a",
					pks: []*gopkg{
						{pkgName: "a0", pkgPath: "a0_path"},
						{pkgName: "a1", pkgPath: "a1_path"},
					},
				},
				{
					name: "b",
					pks: []*gopkg{
						{pkgName: "a2", pkgPath: "a2_path"},
						{pkgName: "a3", pkgPath: "a3_path"},
					},
				},
			},
		}

		// --- When ---
		have := met.imports()

		// --- Then ---
		want := []*gopkg{
			{pkgName: "a0", pkgPath: "a0_path"},
			{pkgName: "a1", pkgPath: "a1_path"},
			{pkgName: "a2", pkgPath: "a2_path"},
			{pkgName: "a3", pkgPath: "a3_path"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("ret imports", func(t *testing.T) {
		// --- Given ---
		met := &method{
			rets: []argument{
				{
					name: "a",
					pks: []*gopkg{
						{pkgName: "r0", pkgPath: "r0_path"},
						{pkgName: "r1", pkgPath: "r1_path"},
					},
				},
				{
					name: "b",
					pks: []*gopkg{
						{pkgName: "r2", pkgPath: "r2_path"},
						{pkgName: "r3", pkgPath: "r3_path"},
					},
				},
			},
		}

		// --- When ---
		have := met.imports()

		// --- Then ---
		want := []*gopkg{
			{pkgName: "r0", pkgPath: "r0_path"},
			{pkgName: "r1", pkgPath: "r1_path"},
			{pkgName: "r2", pkgPath: "r2_path"},
			{pkgName: "r3", pkgPath: "r3_path"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("arg and ret imports", func(t *testing.T) {
		// --- Given ---
		met := &method{
			args: []argument{
				{
					name: "a",
					pks: []*gopkg{
						{pkgName: "a0", pkgPath: "a0_path"},
					},
				},
				{
					name: "b",
					pks: []*gopkg{
						{pkgName: "a1", pkgPath: "a1_path"},
					},
				},
			},
			rets: []argument{
				{
					name: "a",
					pks: []*gopkg{
						{pkgName: "r0", pkgPath: "r0_path"},
					},
				},
				{
					name: "b",
					pks: []*gopkg{
						{pkgName: "r1", pkgPath: "r1_path"},
					},
				},
			},
		}

		// --- When ---
		have := met.imports()

		// --- Then ---
		want := []*gopkg{
			{pkgName: "a0", pkgPath: "a0_path"},
			{pkgName: "a1", pkgPath: "a1_path"},
			{pkgName: "r0", pkgPath: "r0_path"},
			{pkgName: "r1", pkgPath: "r1_path"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("duplicates removed", func(t *testing.T) {
		// --- Given ---
		met := &method{
			args: []argument{
				{
					name: "a",
					pks: []*gopkg{
						{pkgName: "a0", pkgPath: "a0_path"},
					},
				},
				{
					name: "b",
					pks: []*gopkg{
						{pkgName: "a0", pkgPath: "a0_path"},
					},
				},
			},
			rets: []argument{
				{
					name: "a",
					pks: []*gopkg{
						{pkgName: "r0", pkgPath: "r0_path"},
					},
				},
				{
					name: "b",
					pks: []*gopkg{
						{pkgName: "r0", pkgPath: "r0_path"},
					},
				},
			},
		}

		// --- When ---
		have := met.imports()

		// --- Then ---
		want := []*gopkg{
			{pkgName: "a0", pkgPath: "a0_path"},
			{pkgName: "r0", pkgPath: "r0_path"},
		}
		assert.Equal(t, want, have)
	})
}

func Test_method_isVariadic(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		// --- Given ---
		met := &method{
			args: []argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "bool"},
			},
		}

		// --- When ---
		have := met.isVariadic()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("variadic", func(t *testing.T) {
		// --- Given ---
		met := &method{
			args: []argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "...bool"},
			},
		}

		// --- When ---
		have := met.isVariadic()

		// --- Then ---
		assert.True(t, have)
	})
}
