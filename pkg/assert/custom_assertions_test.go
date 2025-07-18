package assert_test

import (
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"

	"github.com/ctx42/testing/pkg/assert"
)

func AssertTA(t tester.T, want, have *types.TA, opts ...check.Option) bool {
	t.Helper()

	ops := check.DefaultOptions(opts...)
	fName := check.FieldName(ops, "TA")

	if !assert.NotNil(t, want, fName("")) {
		return false
	}
	assert.Equal(t, want.Int, have.Int, fName("Int"))
	assert.Equal(t, want.Str, have.Str, fName("Str"))
	assert.Equal(t, want.TAp, have.TAp, fName("TAp"))

	return t.Failed()
}

func Test_AssertTA(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		w := &types.TA{Int: 1, Str: "a", TAp: &types.TA{Int: 2, Str: "b"}}
		h := &types.TA{Int: 1, Str: "a", TAp: &types.TA{Int: 2, Str: "b"}}

		// --- When ---
		AssertTA(tspy, w, h)
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
		AssertTA(tspy, w, h)
	})

	t.Run("both want and have are nil", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		wMsg := "expected non-nil value:\n  trail: TA"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		affirm.Panic(t, func() { AssertTA(tspy, nil, nil) })
	})
}
