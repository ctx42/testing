package check

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_Count(t *testing.T) {
	t.Run("error - unsupported what type", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Count(1, 123, "ab cd ef", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected argument \"what\" to be string got int:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - unsupported where type", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Count(1, "ab", 123, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"unsupported \"where\" type: int:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Count(2, "a", "abc abc anc", opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected string to contain substrings:\n" +
			"       trail: type.field\n" +
			"  want count: 2\n" +
			"  have count: 3\n" +
			"        what: \"a\"\n" +
			"       where: \"abc abc anc\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Count_success_tabular(t *testing.T) {
	tt := []struct {
		testN string

		count int
		what  any
		where any
	}{
		{"one", 1, "ab", "ab cd ef"},
		{"multiple", 2, "ab", "ab cd ab"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Count(tc.count, tc.what, tc.where)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_Count_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		wantCnt int
		haveCnt int
		what    any
		where   any
	}{
		{"not existing", 1, 0, "gh", "ab cd ef"},
		{"existing with wrong count", 2, 1, "ab", "ab cd ef"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Count(tc.wantCnt, tc.what, tc.where)

			// --- Then ---
			affirm.NotNil(t, err)
			wMsg := "" +
				"expected string to contain substrings:\n" +
				"  want count: %d\n" +
				"  have count: %d\n" +
				"        what: %q\n" +
				"       where: %q"
			wMsg = fmt.Sprintf(wMsg, tc.wantCnt, tc.haveCnt, tc.what, tc.where)
			affirm.Equal(t, wMsg, err.Error())
		})
	}
}

func Test_SameType(t *testing.T) {
	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := SameType(42, 4.2, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected same types:\n" +
			"  trail: type.field\n" +
			"   want: int\n" +
			"   have: float64"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_SameType_success_tabular(t *testing.T) {
	var ptr *types.TPtr
	var itf types.TItf
	itf = &types.TPtr{}

	tt := []struct {
		testN string

		val0 any
		val1 any
	}{
		{"int", 42, 44},
		{"float64", 42.0, 44.0},
		{"bool", true, false},
		{"nil ptr 0", ptr, &types.TPtr{}},
		{"nil ptr 1", &types.TPtr{}, ptr},
		{"nil itf 0", itf, &types.TPtr{}},
		{"nil itf 1", &types.TPtr{}, itf},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := SameType(tc.val0, tc.val1)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_SameType_error_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val0 any
		val1 any
		wMsg string
	}{
		{
			"different types",
			42,
			42.0,
			"expected same types:\n  want: int\n  have: float64",
		},
		{
			"different ptr types",
			&types.TPtr{},
			&types.TVal{},
			"expected same types:\n  want: *types.TPtr\n  have: *types.TVal",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := SameType(tc.val0, tc.val1)

			// --- Then ---
			affirm.NotNil(t, err)
			affirm.Equal(t, tc.wMsg, err.Error())
		})
	}
}

func Test_Type(t *testing.T) {
	t.Run("assert type int", func(t *testing.T) {
		// --- Given ---
		var target int

		// --- When ---
		err := Type(&target, 42)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, 42, target)
	})

	t.Run("assert type string", func(t *testing.T) {
		// --- Given ---
		var target string

		// --- When ---
		err := Type(&target, "abc")

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, "abc", target)
	})

	t.Run("assert type struct", func(t *testing.T) {
		// --- Given ---
		var target *types.TPrv
		h := types.TPrv{Pub: 42}.SetInt(44)

		// --- When ---
		err := Type(&target, &h)

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, true, reflect.DeepEqual(target, &h))
	})

	t.Run("error - target must be a pointer", func(t *testing.T) {
		// --- Given ---
		var target int

		// --- When ---
		err := Type(target, 42)

		// --- Then ---
		affirm.Equal(t, "expected target to be a non-nil pointer", err.Error())
	})

	t.Run("error - target must be a non-nil pointer", func(t *testing.T) {
		// --- When ---
		err := Type(nil, 42)

		// --- Then ---
		affirm.Equal(t, "expected target to be a non-nil pointer", err.Error())
	})

	t.Run("error - cannot assert type", func(t *testing.T) {
		// --- Given ---
		var target *types.TPrv
		src := types.TIntStr{}

		// --- When ---
		err := Type(&target, &src)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected type to be assignable to the target:\n" +
			"  target: *types.TPrv\n" +
			"     src: *types.TIntStr"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - cannot assert type", func(t *testing.T) {
		// --- Given ---
		target := &types.TPrv{}
		src := types.TIntStr{}

		// --- When ---
		err := Type(target, &src)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected type to be assignable to the target:\n" +
			"  target: types.TPrv\n" +
			"     src: *types.TIntStr"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Fields(t *testing.T) {
	t.Run("zero fields", func(t *testing.T) {
		// --- Given ---
		s := struct{}{}

		// --- When ---
		err := Fields(0, s)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("value object", func(t *testing.T) {
		// --- When ---
		err := Fields(7, types.TA{})

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("pointer to object", func(t *testing.T) {
		// --- When ---
		err := Fields(7, &types.TA{})

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Fields(1, &types.TA{}, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected struct to have number of fields:\n" +
			"  trail: type.field\n" +
			"   want: 1\n" +
			"   have: 7"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not struct", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Fields(1, 1, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected struct type:\n" +
			"     trail: type.field\n" +
			"  got type: int"
		affirm.Equal(t, wMsg, err.Error())
	})
}
