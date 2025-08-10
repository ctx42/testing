package check_test

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"

	"github.com/ctx42/testing/pkg/check"
)

// AssertTA is custom assertion example.
func AssertTA(t tester.T, want, have *types.TA, opts ...any) bool {
	t.Helper()

	ops := check.DefaultOptions(opts...)
	fName := check.FieldName(ops, "TA")

	if e := check.NotNil(want, fName("")); e != nil {
		t.Error(notice.From(e).Append("argument", "%s", "want"))
		return false
	}
	if e := check.NotNil(have, fName("")); e != nil {
		t.Error(notice.From(e).Append("argument", "%s", "have"))
		return false
	}
	if e := check.Equal(want.Int, have.Int, fName("Int")); e != nil {
		t.Error(e)
		return false
	}
	if e := check.Equal(want.Str, have.Str, fName("Str")); e != nil {
		t.Error(e)
		return false
	}
	if e := check.Equal(want.TAp, have.TAp, fName("TAp")); e != nil {
		t.Error(e)
		return false
	}
	return true
}

func Test_AssertTA(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		w := &types.TA{Int: 1, Str: "a", TAp: &types.TA{Int: 2, Str: "b"}}
		h := &types.TA{Int: 1, Str: "a", TAp: &types.TA{Int: 2, Str: "b"}}

		// --- When ---
		have := AssertTA(tspy, w, h)

		// --- Then ---
		affirm.Equal(t, true, have)
	})

	t.Run("not equal field", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"expected values to be equal:\n" +
			"  trail: TA.TAp.Int\n" +
			"   want: 2\n" +
			"   have: 3"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		w := &types.TA{Int: 1, Str: "a", TAp: &types.TA{Int: 2, Str: "b"}}
		h := &types.TA{Int: 1, Str: "a", TAp: &types.TA{Int: 3, Str: "b"}}

		// --- When ---
		have := AssertTA(tspy, w, h)
		// --- Then ---

		affirm.Equal(t, false, have)
	})

	t.Run("both want and have are nil", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"expected non-nil value:\n" +
			"     trail: TA\n" +
			"  argument: want"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertTA(tspy, nil, nil)

		// --- Then ---
		affirm.Equal(t, false, have)
	})
}
