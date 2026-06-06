package selfkit

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithArgs(t *testing.T) {
	// --- Given ---
	args := []string{"arg0", "arg1"}
	se := &Self{}

	// --- When ---
	WithArgs(args)(se)

	// --- Then ---
	assert.Equal(t, []string{"arg0", "arg1"}, se.args)
}

func Test_New(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		// --- When ---
		se := New(WithArgs([]string{}))

		// --- Then ---
		assert.Equal(t, "", se.toStdout)
		assert.Equal(t, "", se.toStderr)
		assert.Equal(t, "", se.printEnv)
		assert.Equal(t, "", se.printArgs)
		assert.False(t, se.printToStderr)
		assert.False(t, se.noWrap)
		assert.Equal(t, -1, se.exitCode)
	})

	t.Run("unknown arg", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--unknown", "abc"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.Equal(t, "", se.toStdout)
		assert.Equal(t, "", se.toStderr)
		assert.Equal(t, "", se.printEnv)
		assert.Equal(t, "", se.printArgs)
		assert.False(t, se.printToStderr)
		assert.False(t, se.noWrap)
		assert.Equal(t, -1, se.exitCode)
	})

	t.Run("set toStdout", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--toStdout", "abc"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.Equal(t, "abc", se.toStdout)
	})

	t.Run("set toStderr", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--toStderr", "abc"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.Equal(t, "abc", se.toStderr)
	})

	t.Run("set printEnv", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--printEnv", "abc"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.Equal(t, "abc", se.printEnv)
	})

	t.Run("set printToStderr", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--printToStderr"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.True(t, se.printToStderr)
	})

	t.Run("set noWrap", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--noWrap"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.True(t, se.noWrap)
	})

	t.Run("set printArgs", func(t *testing.T) {
		// --- Given ---
		args := []string{"program", "--printArgs", "abc"}

		// --- When ---
		se := New(WithArgs(args))

		// --- Then ---
		assert.Equal(t, "abc", se.printArgs)
	})
}

func Test_Self_Run(t *testing.T) {
	t.Run("no flags runs tests", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--unknown", "abc", "--other", "xyz", "-on"}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.True(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print to stdout", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--toStdout", "abc"}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "|sout: abc|", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print to stdout without wrap", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--noWrap", "--toStdout", " abc "}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, " abc ", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print to stderr", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--toStderr", "abc"}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, "|eout: abc|", eout.String())
	})

	t.Run("print to stderr without wrap", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--noWrap", "--toStderr", " abc "}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, " abc ", eout.String())
	})

	t.Run("print env", func(t *testing.T) {
		// --- Given ---
		kv := fmt.Sprintf("kv%d", time.Now().UnixNano())
		t.Setenv(kv, kv)
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--printEnv", kv}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "|env: "+kv+"|", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print env without wrap", func(t *testing.T) {
		// --- Given ---
		kv := fmt.Sprintf("kv%d", time.Now().UnixNano())
		t.Setenv(kv, kv)
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--noWrap", "--printEnv", kv}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, kv, sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print env to stderr", func(t *testing.T) {
		// --- Given ---
		kv := fmt.Sprintf("kv%d", time.Now().UnixNano())
		t.Setenv(kv, kv)
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--printToStderr", "--printEnv", kv}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, "|env: "+kv+"|", eout.String())
	})

	t.Run("print args", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--printArgs", "arg0", "arg1"}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "|args: arg0 arg1|", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print args without wrap", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--noWrap", "--printArgs", " arg0", "arg1 "}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, " arg0 arg1 ", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("print args to stderr", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{
			"program",
			"--printToStderr",
			"--printArgs",
			"arg0",
			"arg1",
		}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, "|args: arg0 arg1|", eout.String())
	})

	t.Run("set exit code", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--exitCode", "123"}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 123, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("set exit code zero", func(t *testing.T) {
		// --- Given ---
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{"program", "--exitCode", "0"}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "", sout.String())
		assert.Equal(t, "", eout.String())
	})

	t.Run("do all", func(t *testing.T) {
		// --- Given ---
		kv := fmt.Sprintf("kv%d", time.Now().UnixNano())
		t.Setenv(kv, kv)
		sout, eout := &bytes.Buffer{}, &bytes.Buffer{}
		args := []string{
			"program",
			"--printEnv", kv,
			"--toStdout", "abc",
			"--toStderr", "xyz",
			"--printArgs",
			"arg0", "arg1",
		}

		// --- When ---
		runTests, exitCode := New(WithArgs(args)).Run(sout, eout)

		// --- Then ---
		wantSout := "|sout: abc||env: " + kv + "||args: arg0 arg1|"
		assert.False(t, runTests)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, wantSout, sout.String())
		assert.Equal(t, "|eout: xyz|", eout.String())
	})
}
