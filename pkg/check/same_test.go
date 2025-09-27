// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Same(t *testing.T) {
	t.Run("pointers", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TPtr{Val: "0"}

		// --- When ---
		err := Same(ptr0, ptr0)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - want is value", func(t *testing.T) {
		// --- Given ---
		want := testcases.TPtr{Val: "0"}
		have := &testcases.TPtr{Val: "0"}

		// --- When ---
		err := Same(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same pointers:\n" +
			"  want: %%!p(testcases.TPtr={0}) testcases.TPtr{Val:\"0\"}\n" +
			"  have: %p &testcases.TPtr{Val:\"0\"}"
		wMsg = fmt.Sprintf(wMsg, have)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - have is value", func(t *testing.T) {
		// --- Given ---
		want := &testcases.TPtr{Val: "0"}
		have := testcases.TPtr{Val: "0"}

		// --- When ---
		err := Same(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same pointers:\n" +
			"  want: %p &testcases.TPtr{Val:\"0\"}\n" +
			"  have: %%!p(testcases.TPtr={0}) testcases.TPtr{Val:\"0\"}"
		wMsg = fmt.Sprintf(wMsg, want)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - not same pointers", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TPtr{Val: "0"}
		ptr1 := &testcases.TPtr{Val: "1"}

		// --- When ---
		err := Same(ptr0, ptr1)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same pointers:\n" +
			"  want: %p &testcases.TPtr{Val:\"0\"}\n" +
			"  have: %p &testcases.TPtr{Val:\"1\"}"
		wMsg = fmt.Sprintf(wMsg, ptr0, ptr1)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TPtr{Val: "0"}
		ptr1 := &testcases.TPtr{Val: "1"}

		opt := WithTrail("type.field")

		// --- When ---
		err := Same(ptr0, ptr1, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected same pointers:\n" +
			"  trail: type.field\n" +
			"   want: %p &testcases.TPtr{Val:\"0\"}\n" +
			"   have: %p &testcases.TPtr{Val:\"1\"}"
		wMsg = fmt.Sprintf(wMsg, ptr0, ptr1)
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Same_tabular(t *testing.T) {
	ptr0 := &testcases.TPtr{Val: "0"}
	ptr1 := &testcases.TPtr{Val: "1"}
	var itfPtr0, itfPtr1 testcases.TItf
	itfPtr0, itfPtr1 = &testcases.TPtr{Val: "0"}, &testcases.TPtr{Val: "1"}

	tt := []struct {
		testN string

		p0   any
		p1   any
		same bool
	}{
		{"same ptr", ptr0, ptr0, true},
		{"not same ptr", ptr0, ptr1, false},
		{"itf ptr", itfPtr0, itfPtr0, true},

		{"not same itf ptr", itfPtr0, itfPtr1, false},
		{"not same val", testcases.TVal{}, testcases.TVal{}, false},
		{"not same ptr different types", &testcases.TPtr{}, &testcases.TVal{}, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := Same(tc.p0, tc.p1)

			// --- Then ---
			if tc.same && have != nil {
				format := "expected same values:\n  want: %#v\n  have: %#v"
				t.Errorf(format, tc.p0, tc.p1)
			}
		})
	}
}

func Test_NotSame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TPtr{Val: "0"}
		ptr1 := &testcases.TPtr{Val: "0"}

		// --- When ---
		err := NotSame(ptr0, ptr1)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TPtr{Val: "0"}

		// --- When ---
		err := NotSame(ptr0, ptr0)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected different pointers:\n" +
			"  want: %p &testcases.TPtr{Val:\"0\"}\n" +
			"  have: %p &testcases.TPtr{Val:\"0\"}"
		wMsg = fmt.Sprintf(wMsg, ptr0, ptr0)
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		ptr0 := &testcases.TPtr{Val: "0"}

		opt := WithTrail("type.field")

		// --- When ---
		err := NotSame(ptr0, ptr0, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected different pointers:\n" +
			"  trail: type.field\n" +
			"   want: %p &testcases.TPtr{Val:\"0\"}\n" +
			"   have: %p &testcases.TPtr{Val:\"0\"}"
		wMsg = fmt.Sprintf(wMsg, ptr0, ptr0)
		affirm.Equal(t, wMsg, err.Error())
	})
}
