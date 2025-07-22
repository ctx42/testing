package check

import (
	"testing"
	"time"

	"github.com/ctx42/testing/internal/affirm"
)

func Test_Wait(t *testing.T) {
	t.Run("fn called only once", func(t *testing.T) {
		// --- Given ---
		var cnt int
		fn := func() bool { cnt++; return true }

		// --- When ---
		err := Wait("100ms", fn)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, 1, cnt)
	})

	t.Run("error - wait timeout", func(t *testing.T) {
		// --- Given ---
		var cnt int
		fn := func() bool { cnt++; return false }
		start := time.Now()

		// --- When ---
		err := Wait("100ms", fn)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected function to return true:\n" +
			"    within: 100ms\n" +
			"  throttle: 10ms"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, cnt >= 9 && cnt <= 10)
		affirm.Equal(t, true, 100*time.Millisecond <= time.Since(start))
	})

	t.Run("error - wait given timeout with throttle", func(t *testing.T) {
		// --- Given ---
		var cnt int
		fn := func() bool { cnt++; return false }
		start := time.Now()
		throttle := WithWaitThrottle(5 * time.Millisecond)

		// --- When ---
		err := Wait("101ms", fn, throttle)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected function to return true:\n" +
			"    within: 101ms\n" +
			"  throttle: 5ms"
		affirm.Equal(t, wMsg, err.Error())
		affirm.Equal(t, true, cnt >= 19 && cnt <= 20)
		affirm.Equal(t, true, 100*time.Millisecond <= time.Since(start))
	})

	t.Run("error - trail set", func(t *testing.T) {
		// --- Given ---
		fn := func() bool { return false }
		opt := WithTrail("type.field")

		// --- When ---
		err := Wait("100ms", fn, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected function to return true:\n" +
			"     trail: type.field\n" +
			"    within: 100ms\n" +
			"  throttle: 10ms"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("with throttle", func(t *testing.T) {
		// --- Given ---
		var cnt int
		fn := func() bool {
			cnt++
			return true
		}
		throttle := WithWaitThrottle(100 * time.Millisecond)

		// --- When ---
		err := Wait("1s", fn, throttle)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, 1, cnt)
	})

	t.Run("error - invalid timeout", func(t *testing.T) {
		// --- When ---
		err := Wait("abc", nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "[within] failed to parse duration:\n  value: abc"
		affirm.Equal(t, wMsg, err.Error())
	})
}
