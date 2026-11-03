// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import (
	"strings"
	"testing"

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
		assert.Equal(t, goldy.Open(t, gfp).String(), have)
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

	t.Run("wrapped", func(t *testing.T) {
		// --- Given ---
		one := strings.Repeat("a", 30)
		two := strings.Repeat("b", 30)
		met := &method{
			name: "Long",
			args: []argument{{name: one, typ: "int"}, {name: two, typ: "bool"}},
		}

		// --- When ---
		have := met.generateOn("M")

		// --- Then ---
		want := "" +
			"// OnLong expects a call of [M.Long] with the arguments.\n" +
			"func (_mck *M) OnLong(\n" +
			"\t" + one + " any,\n" +
			"\t" + two + " any,\n" +
			") *mock.Call {\n" +
			"\n" +
			"\t_mck.t.Helper()\n" +
			"\t_args := []any{" + one + ", " + two + "}\n" +
			"\treturn _mck.On(\"Long\", _args...)\n" +
			"}"
		assert.Equal(t, want, have)
	})

	t.Run("wrapped doc", func(t *testing.T) {
		// --- Given ---
		name := strings.Repeat("N", 20)
		met := &method{name: name, args: []argument{{name: "a", typ: "int"}}}

		// --- When ---
		have := met.generateOn("M")

		// --- Then ---
		want := "" +
			"// On" + name + " expects a call of [M." + name + "] with the\n" +
			"// arguments.\n" +
			"func (_mck *M) On" + name + "(a any) *mock.Call {\n" +
			"\t_mck.t.Helper()\n" +
			"\t_args := []any{a}\n" +
			"\treturn _mck.On(\"" + name + "\", _args...)\n" +
			"}"
		assert.Equal(t, want, have)
	})
}

func Test_method_genDoc(t *testing.T) {
	t.Run("no returns", func(t *testing.T) {
		// --- Given ---
		met := &method{name: "Close"}

		// --- When ---
		have := met.genDoc()

		// --- Then ---
		assert.Equal(t, "// Close records a call of Close.\n", have)
	})

	t.Run("with returns", func(t *testing.T) {
		// --- Given ---
		met := &method{name: "Count", rets: []argument{{typ: "int"}}}

		// --- When ---
		have := met.genDoc()

		// --- Then ---
		want := "" +
			"// Count records a call of Count and returns its set results.\n"
		assert.Equal(t, want, have)
	})

	t.Run("wrapped", func(t *testing.T) {
		// --- Given ---
		name := strings.Repeat("N", 20)
		met := &method{name: name, rets: []argument{{typ: "int"}}}

		// --- When ---
		have := met.genDoc()

		// --- Then ---
		want := "" +
			"// " + name + " records a call of " + name + " and returns its\n" +
			"// set results.\n"
		assert.Equal(t, want, have)
	})
}

func Test_wrapSig_tabular(t *testing.T) {
	head := "func (_mck *M) Method("
	arg80 := strings.Repeat("a", 51) + " int" // Signature with " {": 80.
	arg81 := strings.Repeat("a", 52) + " int" // Signature with " {": 81.
	long := strings.Repeat("M", 70)

	tt := []struct {
		testN string

		head string
		args []string
		tail string
		want string
	}{
		{
			"fits at 80 columns",
			head,
			[]string{arg80},
			")",
			head + arg80 + ") {\n",
		},
		{
			"wraps at 81 columns",
			head,
			[]string{arg81},
			")",
			head + "\n\t" + arg81 + ",\n) {\n\n",
		},
		{
			"long without arguments",
			"func (_mck *M) " + long + "(",
			nil,
			") int",
			"func (_mck *M) " + long + "() int {\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := wrapSig(tc.head, tc.args, tc.tail)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_wrapComment_tabular(t *testing.T) {
	word75 := strings.Repeat("a", 75)
	word76 := strings.Repeat("a", 76)
	long := strings.Repeat("b", 90)

	tt := []struct {
		testN string

		text string
		want string
	}{
		{"short", "Close records a call.", "// Close records a call.\n"},
		{"fits at 80 columns", "a " + word75, "// a " + word75 + "\n"},
		{"wraps at 81 columns", "a " + word76, "// a\n// " + word76 + "\n"},
		{
			"word over the limit",
			"a " + long + " c",
			"// a\n// " + long + "\n// c\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := wrapComment(tc.text)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genArgList(t *testing.T) {
	// --- Given ---
	met := &method{
		args: []argument{
			{name: "x", typ: "int"},
			{name: "", typ: "string"},
			{name: "y", typ: "...bool"},
		},
	}

	// --- When ---
	have := met.genArgList()

	// --- Then ---
	assert.Equal(t, []string{"x int", "_a1 string", "y ...bool"}, have)
}

func Test_method_genAnyArgList(t *testing.T) {
	// --- Given ---
	met := &method{
		args: []argument{
			{name: "x", typ: "int"},
			{name: "y", typ: "...bool"},
		},
	}

	// --- When ---
	have := met.genAnyArgList()

	// --- Then ---
	assert.Equal(t, []string{"x any", "y ...any"}, have)
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

func Test_method_argNames_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args []argument
		want []string
	}{
		{"nil args", nil, nil},
		{"no args", []argument{}, nil},
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
		{"nil args", nil, "()"},
		{"no args", []argument{}, "()"},
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
			"func (_mck *TypeMock) Method(int) (int, error)",
		},
		{
			"one argument no returns",
			"TypeMock",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			nil,
			"func (_mck *TypeMock) Method(int)",
		},
		{
			"no receiver type",
			"",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
			},
			nil,
			"func Method(int)",
		},
		{
			"no receiver type two arguments",
			"",
			"Method",
			[]argument{
				{name: "a", typ: "int"},
				{name: "b", typ: "float64"},
			},
			nil,
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
			have := met.genSig(tc.typ)

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
		{"no arguments", nil, "\tvar _args []any\n"},
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
		{"no returns", nil, ""},
		{
			"one return",
			[]argument{
				{name: "", typ: "error"},
			},
			"" +
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
			"" +
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

func Test_method_genReturnBody(t *testing.T) {
	t.Run("no returns", func(t *testing.T) {
		// --- Given ---
		met := &method{}

		// --- When ---
		have := met.genReturnBody(1)

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("one return", func(t *testing.T) {
		// --- Given ---
		met := &method{name: "Count", rets: []argument{{typ: "int"}}}

		// --- When ---
		have := met.genReturnBody(1)

		// --- Then ---
		want := "" +
			"\tvar _r0 int\n" +
			"\tswitch _v := _rets.Get(0).(type) {\n" +
			"\tcase func() int:\n" +
			"\t\t_r0 = _v()\n" +
			"\tcase int:\n" +
			"\t\t_r0 = _v\n" +
			"\tcase nil:\n" +
			"\tdefault:\n" +
			"\t\t_mck.t.Fatalf(\"Count return 0 is %T, not int\", _v)\n" +
			"\t}"
		assert.Equal(t, want, have)
	})
}

func Test_method_genSingleRet_tabular(t *testing.T) {
	// Widths below are of the "Method" return 1 lines at indent 0.
	fit := "pkg." + strings.Repeat("T", 24)    // Fatalf call: 80 columns.
	over := fit + "T"                          // Fatalf call: 81 columns.
	msgFit := "pkg." + strings.Repeat("T", 38) // Fatalf message: 80 columns.
	msgOver := msgFit + "T"                    // Fatalf message: 81 columns.
	argFit := "pkg." + strings.Repeat("A", 60) // Func case: 80 columns.
	argOver := argFit + "A"                    // Func case: 81 columns.
	long := "pkg." + strings.Repeat("T", 30)

	tt := []struct {
		testN string

		args   []argument
		ret    argument
		indent int
		want   []string
	}{
		{
			"only return arg",
			[]argument{},
			argument{name: "", typ: "int"},
			0,
			[]string{
				"var _r1 int",
				"switch _v := _rets.Get(1).(type) {",
				"case func() int:",
				"\t_r1 = _v()",
				"case int:",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(\"Method return 1 is %T, not int\", _v)",
				"}",
			},
		},
		{
			"arg and ret arg",
			[]argument{
				{name: "a", typ: "int"},
			},
			argument{name: "", typ: "int"},
			0,
			[]string{
				"var _r1 int",
				"switch _v := _rets.Get(1).(type) {",
				"case func(int) int:",
				"\t_r1 = _v(a)",
				"case int:",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(\"Method return 1 is %T, not int\", _v)",
				"}",
			},
		},
		{
			"percent and quote in type",
			[]argument{},
			argument{name: "", typ: `[len("a%")]byte`},
			0,
			[]string{
				`var _r1 [len("a%")]byte`,
				"switch _v := _rets.Get(1).(type) {",
				`case func() [len("a%")]byte:`,
				"\t_r1 = _v()",
				`case [len("a%")]byte:`,
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(" +
					`"Method return 1 is %T, not [len(\"a%%\")]byte", _v)`,
				"}",
			},
		},
		{
			"fatal call fits at 80 columns",
			[]argument{},
			argument{name: "", typ: fit},
			0,
			[]string{
				"var _r1 " + fit,
				"switch _v := _rets.Get(1).(type) {",
				"case func() " + fit + ":",
				"\t_r1 = _v()",
				"case " + fit + ":",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(\"Method return 1 is %T, not " +
					fit + "\", _v)",
				"}",
			},
		},
		{
			"fatal call wraps at 81 columns",
			[]argument{},
			argument{name: "", typ: over},
			0,
			[]string{
				"var _r1 " + over,
				"switch _v := _rets.Get(1).(type) {",
				"case func() " + over + ":",
				"\t_r1 = _v()",
				"case " + over + ":",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not " + over + "\",",
				"\t\t_v,",
				"\t)",
				"}",
			},
		},
		{
			"fatal message fits at 80 columns",
			[]argument{},
			argument{name: "", typ: msgFit},
			0,
			[]string{
				"var _r1 " + msgFit,
				"switch _v := _rets.Get(1).(type) {",
				"case func() " + msgFit + ":",
				"\t_r1 = _v()",
				"case " + msgFit + ":",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not " + msgFit + "\",",
				"\t\t_v,",
				"\t)",
				"}",
			},
		},
		{
			"fatal message splits at 81 columns",
			[]argument{},
			argument{name: "", typ: msgOver},
			0,
			[]string{
				"var _r1 " + msgOver,
				"switch _v := _rets.Get(1).(type) {",
				"case func() " + msgOver + ":",
				"\t_r1 = _v()",
				"case " + msgOver + ":",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not \"+",
				"\t\t\t\"" + msgOver + "\",",
				"\t\t_v,",
				"\t)",
				"}",
			},
		},
		{
			"func case fits at 80 columns",
			[]argument{{name: "a", typ: argFit}},
			argument{name: "", typ: "int"},
			0,
			[]string{
				"var _r1 int",
				"switch _v := _rets.Get(1).(type) {",
				"case func(" + argFit + ") int:",
				"\t_r1 = _v(a)",
				"case int:",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(\"Method return 1 is %T, not int\", _v)",
				"}",
			},
		},
		{
			"func case wraps at 81 columns",
			[]argument{{name: "a", typ: argOver}},
			argument{name: "", typ: "int"},
			0,
			[]string{
				"var _r1 int",
				"switch _v := _rets.Get(1).(type) {",
				"case func(",
				"\t" + argOver + ",",
				") int:",
				"\t_r1 = _v(a)",
				"case int:",
				"\t_r1 = _v",
				"case nil:",
				"default:",
				"\t_mck.t.Fatalf(\"Method return 1 is %T, not int\", _v)",
				"}",
			},
		},
		{
			"wrapped fatal call",
			[]argument{},
			argument{name: "", typ: long},
			1,
			[]string{
				"\tvar _r1 " + long,
				"\tswitch _v := _rets.Get(1).(type) {",
				"\tcase func() " + long + ":",
				"\t\t_r1 = _v()",
				"\tcase " + long + ":",
				"\t\t_r1 = _v",
				"\tcase nil:",
				"\tdefault:",
				"\t\t_mck.t.Fatalf(",
				"\t\t\t\"Method return 1 is %T, not " + long + "\",",
				"\t\t\t_v,",
				"\t\t)",
				"\t}",
			},
		},
		{
			"ret argument is of any type",
			[]argument{},
			argument{name: "", typ: "any"},
			0,
			[]string{
				"var _r1 any",
				"switch _v := _rets.Get(1).(type) {",
				"case func() any:",
				"\t_r1 = _v()",
				"default:",
				"\t_r1 = _v",
				"}",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{name: "Method", args: tc.args}

			// --- When ---
			have := met.genSingleRet(1, tc.ret, tc.indent)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genFuncCase_tabular(t *testing.T) {
	argFit := "pkg." + strings.Repeat("A", 60) // Case: 80 columns.
	argOver := argFit + "A"                    // Case: 81 columns.
	long := "pkg." + strings.Repeat("R", 80)

	tt := []struct {
		testN string

		args []argument
		ret  string
		want []string
	}{
		{
			"fits at 80 columns",
			[]argument{{name: "a", typ: argFit}},
			"int",
			[]string{"case func(" + argFit + ") int:"},
		},
		{
			"wraps at 81 columns",
			[]argument{{name: "a", typ: argOver}},
			"int",
			[]string{"case func(", "\t" + argOver + ",", ") int:"},
		},
		{
			"wraps variadic",
			[]argument{{name: "a", typ: argOver}, {name: "b", typ: "...int"}},
			"int",
			[]string{"case func(", "\t" + argOver + ",", "\t...int,", ") int:"},
		},
		{
			"long without arguments",
			nil,
			long,
			[]string{"case func() " + long + ":"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{name: "Method", args: tc.args}

			// --- When ---
			have := met.genFuncCase(argument{typ: tc.ret}, "")

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_method_genFatalf_tabular(t *testing.T) {
	// Widths below are of the "Method" return 1 lines indented one tab.
	callFit := "pkg." + strings.Repeat("T", 24) // Call: 80 columns.
	callOver := callFit + "T"                   // Call: 81 columns.
	msgFit := "pkg." + strings.Repeat("T", 38)  // Message: 80 columns.
	msgOver := msgFit + "T"                     // Message: 81 columns.

	tt := []struct {
		testN string

		typ  string
		want []string
	}{
		{
			"call fits at 80 columns",
			callFit,
			[]string{
				"\t_mck.t.Fatalf(\"Method return 1 is %T, not " + callFit +
					"\", _v)",
			},
		},
		{
			"call breaks at 81 columns",
			callOver,
			[]string{
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not " + callOver + "\",",
				"\t\t_v,",
				"\t)",
			},
		},
		{
			"message fits at 80 columns",
			msgFit,
			[]string{
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not " + msgFit + "\",",
				"\t\t_v,",
				"\t)",
			},
		},
		{
			"message splits at 81 columns",
			msgOver,
			[]string{
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not \"+",
				"\t\t\t\"" + msgOver + "\",",
				"\t\t_v,",
				"\t)",
			},
		},
		{
			"percent and quote in split type",
			`[len("a%")]` + msgOver,
			[]string{
				"\t_mck.t.Fatalf(",
				"\t\t\"Method return 1 is %T, not \"+",
				"\t\t\t" + `"[len(\"a%%\")]` + msgOver + "\",",
				"\t\t_v,",
				"\t)",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			met := &method{name: "Method"}

			// --- When ---
			have := met.genFatalf(1, argument{typ: tc.typ}, "\t")

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
		{"no returns", nil, nil, ""},
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
