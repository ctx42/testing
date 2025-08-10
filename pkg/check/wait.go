package check

import (
	"context"
	"time"

	"github.com/ctx42/testing/pkg/notice"
)

// Wait waits for "fn" to return true but no longer then given timeout. By
// default, calls to "fn" are throttled with a default throttle set in
// [Options.WaitThrottle] - use [WithWaitThrottle] option to change it.
//
// The "timeout" may represent duration in the form of a string, int, int64 or
// [time.Duration].
func Wait(within any, fn func() bool, opts ...any) error {
	ops := DefaultOptions(opts...)

	dur, durStr, _, err := getDur(within, opts...)
	if err != nil {
		return notice.From(err, "within")
	}

	ctx, cxl := context.WithTimeout(context.Background(), dur)
	defer cxl()

	sleep := time.NewTimer(ops.WaitThrottle)
	defer sleep.Stop()

	for {
		sleep.Reset(ops.WaitThrottle)
		select {
		case <-ctx.Done():
			if !sleep.Stop() {
				<-sleep.C
			}
			msg := notice.New("expected function to return true").
				Append("within", "%s", durStr).
				Append("throttle", ops.WaitThrottle.String())
			return AddRows(ops, msg)

		default:
			if fn() {
				return nil
			}
			<-sleep.C
		}
	}
}
