package mock

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/assert"
)

func Test_callStack(t *testing.T) {
	// --- When ---
	have := callStack()

	// --- Then ---
	assert.Len(t, 2, have)
	assert.Contain(t, "/pkg/mock/helpers.go:", have[0])

	ln := have[0][strings.Index(have[0], ":")+1:]
	_, err := strconv.Atoi(ln)
	assert.Nil(t, err)
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

func Test_formatArgs_tabular(t *testing.T) {
	tt := []struct {
		testN string

		args Arguments
		want string
	}{
		{"nil arguments", nil, ""},
		{"single simple argument", []any{1}, "0: 1"},
		// TODO(rz): see how dumper is dumping multi line strings in assert.Equal.
		{
			"multiple simple arguments",
			[]any{1, "abc", 2.2},
			"" +
				"0: 1\n" +
				"1: \"abc\"\n" +
				"2: 2.2",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := formatArgs(tc.args)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_formatStack(t *testing.T) {
	t.Run("stack without padding", func(t *testing.T) {
		// --- Given ---
		stack := []string{"/dir/file.go:1", "/dir/file.go:2"}

		// --- When ---
		have := formatStack(stack, 0)

		// --- Then ---
		want := "" +
			"/dir/file.go:1\n" +
			"/dir/file.go:2"
		assert.Equal(t, want, have)
	})

	t.Run("stack with padding", func(t *testing.T) {
		// --- Given ---
		stack := []string{"/dir/file.go:1", "/dir/file.go:2"}

		// --- When ---
		have := formatStack(stack, 1)

		// --- Then ---
		want := "" +
			"\t/dir/file.go:1\n" +
			"\t/dir/file.go:2"
		assert.Equal(t, want, have)
	})

	t.Run("empty", func(t *testing.T) {
		// --- When ---
		have := formatStack([]string{}, 0)

		// --- Then ---
		assert.Equal(t, "", have)
	})
}

func Test_formatDiff(t *testing.T) {
	t.Run("format", func(t *testing.T) {
		// --- Given ---
		diff := []string{"AAA", "BBB", "CCC"}

		// --- When ---
		have := formatDiff(diff, 0)

		// --- Then ---
		assert.Equal(t, "AAA\nBBB\nCCC\n", have)
	})

	t.Run("padding", func(t *testing.T) {
		// --- Given ---
		diff := []string{"AAA", "BBB", "CCC"}

		// --- When ---
		have := formatDiff(diff, 1)

		// --- Then ---
		assert.Equal(t, "\tAAA\n\tBBB\n\tCCC\n", have)
	})

	t.Run("empty", func(t *testing.T) {
		// --- Given ---
		diff := make([]string, 0)

		// --- When ---
		have := formatDiff(diff, 1)

		// --- Then ---
		assert.Equal(t, "", have)
	})
}

func Test_isTest_tabular(t *testing.T) {
	tt := []struct {
		name   string
		prefix string
		want   bool
	}{
		{"", "", true},
		{"TestAbc", "Test", true},
		{"Test_Abc", "Test", true},
		{"Test_abc", "Test", true},
		{"BenchmarkAbc", "Benchmark", true},
		{"Benchmark_Abc", "Benchmark", true},
		{"Benchmark_abc", "Benchmark", true},
		{"ExampleAbc", "Example", true},
		{"Example_Abc", "Example", true},
		{"Example_abc", "Example", true},
		{"Test", "Test", true},
		{"Benchmark", "Benchmark", true},
		{"Example", "Example", true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, isTest(tc.name, tc.prefix))
		})
	}
}

func Test_methodName(t *testing.T) {
	// --- Given ---
	ptr := &types.TPtr{}
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
			"1",
			[]string{"aaa", "bb", "c"},
			[]string{"111", "222", "333"},
			[]string{
				"aaa 111",
				"bb  222",
				"c   333",
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
