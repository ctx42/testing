package assert

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_wait(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		fn := func() bool { return true }

		// --- When ---
		have := Wait(tspy, "100ms", fn)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		fn := func() bool { return false }

		// --- When ---
		have := Wait(tspy, "100ms", fn)

		// --- Then ---
		affirm.Equal(t, false, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		fn := func() bool { return false }
		opt := check.WithTrail("type.field")

		// --- When ---
		have := Wait(tspy, "100ms", fn, opt)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
