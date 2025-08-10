package assert

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Wait waits for "fn" to return true but no longer then given timeout. By
// default, calls to "fn" are throttled with a default throttle set in
// [check.Options.WaitThrottle] - use [check.WithWaitThrottle] option to change
// it. Returns true when the function returns true within given timeout,
// otherwise marks the test as failed, writes an error message to the test log
// and returns false.
//
// The "timeout" may represent duration in the form of a string, int, int64 or
// [time.Duration].
func Wait(t tester.T, timeout string, fn func() bool, opts ...any) bool {
	t.Helper()
	if e := check.Wait(timeout, fn, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
