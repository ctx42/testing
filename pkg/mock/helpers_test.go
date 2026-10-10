// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock

import (
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_callStack(t *testing.T) {
	t.Run("frames from other directories", func(t *testing.T) {
		// --- Given ---
		var have []string
		less := func(i, j int) bool {
			if have == nil {
				have = callStack()
			}
			return i < j
		}

		// --- When ---
		sort.Slice([]int{2, 1, 3}, less)

		// --- Then ---
		for _, frame := range have {
			pth := frame[:strings.LastIndex(frame, ":")]
			_, err := os.Stat(pth)
			assert.NoError(t, err)
		}
	})

	// --- When ---
	have := callStack()

	// --- Then ---
	assert.Len(t, 2, have)
	assert.Contain(t, "/pkg/mock/helpers.go:", have[0])
	ln := have[0][strings.Index(have[0], ":")+1:]
	_, err := strconv.Atoi(ln)
	assert.NoError(t, err)
}

func Test_formatMethod_tabular(t *testing.T) {
	tt := []struct {
		testN string

		method string
		args   Arguments
		rets   Arguments
		want   string
	}{
		{"no args no rets", "Method", nil, nil, "Method()"},
		{"with one arg no rets", "Method", []any{1}, nil, "Method(int)"},
		{
			"with one arg and one ret",
			"Method",
			[]any{1},
			[]any{1},
			"Method(int) int",
		},
		{
			"with two args and one ret",
			"Method",
			[]any{1, 2},
			[]any{1},
			"Method(int, int) int",
		},
		{
			"with two args and two ret",
			"Method",
			[]any{1, 2},
			[]any{1, 2},
			"Method(int, int) (int, int)",
		},
		{
			"with no args and one ret",
			"Method",
			nil,
			[]any{1},
			"Method() int",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := formatMethod(tc.method, tc.args, tc.rets)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_isTestName_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name   string
		prefix string
		want   bool
	}{
		{"empty empty", "", "", true},
		{"TestAbc Test", "TestAbc", "Test", true},
		{"Test_Abc Test", "Test_Abc", "Test", true},
		{"Test_abc Test", "Test_abc", "Test", true},
		{"BenchmarkAbc Benchmark", "BenchmarkAbc", "Benchmark", true},
		{"Benchmark_Abc Benchmark", "Benchmark_Abc", "Benchmark", true},
		{"Benchmark_abc Benchmark", "Benchmark_abc", "Benchmark", true},
		{"ExampleAbc Example", "ExampleAbc", "Example", true},
		{"Example_Abc Example", "Example_Abc", "Example", true},
		{"Example_abc Example", "Example_abc", "Example", true},
		{"Test Test", "Test", "Test", true},
		{"Benchmark Benchmark", "Benchmark", "Benchmark", true},
		{"Example Example", "Example", "Example", true},
		{"TesticularCancer Test", "TesticularCancer", "Test", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Equal(t, tc.want, isTestName(tc.name, tc.prefix))
		})
	}
}

func Test_methodName(t *testing.T) {
	// --- Given ---
	ptr := &testcases.TPtr{}
	met := reflect.ValueOf(ptr.AAA)

	// --- When ---
	have := methodName(met)

	// --- Then ---
	assert.Equal(t, "AAA", have)
}

func Test_twoColumns_tabular(t *testing.T) {
	tt := []struct {
		testN string

		col1 []string
		col2 []string
		want []string
	}{
		{
			"three rows",
			[]string{"aaa", "bb", "c"},
			[]string{"111", "222", "333"},
			[]string{
				"aaa -> 111",
				"bb  -> 222",
				"c   -> 333",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			twoColumns(tc.col1, tc.col2)

			// --- Then ---
			assert.Equal(t, tc.want, tc.col1)
		})
	}
}
