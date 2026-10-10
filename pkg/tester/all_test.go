// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
)

// assertSpyHasMsgs asserts spy has saved a "want" number of messages. On
// failure, it marks the test as failed, returns false, but continues execution.
func assertSpyHasMsgs(t T, spy *Spy, want int) bool {
	t.Helper()
	have := len(spy.savedMsgs)
	if want != have {
		var msg strings.Builder
		format := "" +
			"expected Spy to have saved messages:\n" +
			"\twant: %d\n" +
			"\thave: %d"
		_, _ = fmt.Fprintf(&msg, format, want, have)
		if have > 0 {
			_, _ = msg.WriteString("\n\tmessages:\n")
			for idx, sm := range spy.savedMsgs {
				_, _ = fmt.Fprintf(&msg, "\t\t%d: %q\n", idx, sm)
			}
		}
		t.Error(msg.String())
		return false
	}
	return true
}

func Test_assertSpyHasMsgs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := &Spy{savedMsgs: []string{"abc", "def"}}

		// --- When ---
		have := assertSpyHasMsgs(ti, spy, 2)

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, ti.Failed())
	})

	t.Run("error - count differs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := &Spy{savedMsgs: []string{"abc", "def"}}

		// --- When ---
		have := assertSpyHasMsgs(ti, spy, 3)

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, true, ti.Failed())
	})
}

// assertSpyHasMsg asserts Spy has saved a message with given idx equal to
// "want". On failure, it marks the test as failed, returns false, but
// continues execution.
func assertSpyHasMsg(t T, spy *Spy, idx int, want string) bool {
	t.Helper()
	if idx >= len(spy.savedMsgs) {
		t.Errorf("expected Spy saved message with index %d to exist", idx)
		return false
	}
	have := spy.savedMsgs[idx]
	if want != have {
		format := "" +
			"expected Spy saved message with index %d to be:\n" +
			"\twant: %q\n" +
			"\thave: %q"
		t.Errorf(format, idx, want, have)
		return false
	}
	return true
}

func Test_assertSpyHasMsg(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := &Spy{savedMsgs: []string{"abc", "def"}}

		// --- When ---
		have := assertSpyHasMsg(ti, spy, 1, "def")

		// --- Then ---
		affirm.Equal(t, true, have)
		affirm.Equal(t, false, ti.Failed())
	})

	t.Run("error - message differs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := &Spy{savedMsgs: []string{"abc", "def"}}

		// --- When ---
		have := assertSpyHasMsg(ti, spy, 1, "xyz")

		// --- Then ---
		affirm.Equal(t, false, have)
		affirm.Equal(t, true, ti.Failed())
	})
}

// isolateEnv sets keys to an empty value on t so they are restored when t
// ends. A Spy bound to a zero testing.T calls its Setenv, whose restore
// never runs.
func isolateEnv(t T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		t.Setenv(key, "")
	}
}

// removeTempDirs removes, when t ends, the parent directories testing creates
// for the temporary directories returned by spy. A Spy bound to a zero
// testing.T calls its TempDir, whose cleanup never runs.
func removeTempDirs(t T, spy *Spy) {
	t.Helper()
	t.Cleanup(func() {
		for _, dir := range spy.haveTempDirs {
			_ = os.RemoveAll(filepath.Dir(dir))
		}
	})
}
