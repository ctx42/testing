// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_Error(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := Error(errors.New("e0"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := Error(nil)

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "expected non-nil error", err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Error(nil, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected non-nil error:\n  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NoError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := NoError(nil)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := NoError(errors.New("e0"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error to be nil:\n" +
			"  want: nil\n" +
			"  have: \"e0\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NoError(errors.New("e0"), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error to be nil:\n" +
			"  trail: type.field\n" +
			"   want: nil\n" +
			"   have: \"e0\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil interface", func(t *testing.T) {
		// --- Given ---
		var e *types.TPtr

		// --- When ---
		err := NoError(e)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error to be nil:\n" +
			"  want: nil\n" +
			"  have: *types.TPtr"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		var e *types.TPtr
		opt := WithTrail("type.field")

		// --- When ---
		err := NoError(e, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error to be nil:\n" +
			"  trail: type.field\n" +
			"   want: nil\n" +
			"   have: *types.TPtr"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("empty non-nil error", func(t *testing.T) {
		// --- Given ---
		e := &types.TPtr{}

		// --- When ---
		err := NoError(e)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error to be nil:\n" +
			"  want: nil\n" +
			"  have: \"\""
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_ErrorIs(t *testing.T) {
	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		have := ErrorIs(errors.New("err0"), errors.New("err0"), opt)

		// --- Then ---
		affirm.NotNil(t, have)
		wMsg := "expected error to have a target in its tree:\n" +
			"  trail: type.field\n" +
			"   want: (*errors.errorString) err0\n" +
			"   have: (*errors.errorString) err0"
		affirm.Equal(t, wMsg, have.Error())
	})
}

func Test_ErrorIs_success_tabular(t *testing.T) {
	err0 := errors.New("err0")
	err1 := errors.New("err1")
	err2 := fmt.Errorf("wrap: %w %w", err0, err1)

	tt := []struct {
		testN string

		err    error
		target error
	}{
		{"nil nil", nil, nil},
		{"same error", err0, err0},
		{"related 0", err2, err0},
		{"related 1", err2, err1},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := ErrorIs(tc.target, tc.err)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_ErrorIs_error_tabular(t *testing.T) {
	err0 := errors.New("err0")
	err1 := errors.New("err1")
	err2 := fmt.Errorf("wrap: %w %w", err0, err1)

	tt := []struct {
		testN string

		have     error
		haveType string
		haveStr  string
		want     error
		wantType string
		wantStr  string
	}{
		{"err nil", err0, "*errors.errorString", "err0", nil, "<nil>", "<nil>"},
		{"nil err", nil, "<nil>", "<nil>", err0, "*errors.errorString", "err0"},
		{"not related", err0, "*errors.errorString", "err0", err1, "*errors.errorString", "err1"},
		{"not related 0", err0, "*errors.errorString", "err0", err2, "*fmt.wrapErrors", "wrap: err0 err1"},
		{"not related 1", err1, "*errors.errorString", "err1", err2, "*fmt.wrapErrors", "wrap: err0 err1"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			const format = "expected error to have a target in its tree:\n" +
				"  want: (%s) %s\n" +
				"  have: (%s) %s"

			wantLog := fmt.Sprintf(
				format,
				tc.wantType,
				tc.wantStr,
				tc.haveType,
				tc.haveStr,
			)

			// --- When ---
			err := ErrorIs(tc.want, tc.have)

			// --- Then ---
			affirm.NotNil(t, err)
			affirm.Equal(t, wantLog, err.Error())
		})
	}
}

func Test_ErrorAs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		var target *types.TPtr

		// --- When ---
		err := ErrorAs(&target, &types.TPtr{Val: "A"})

		// --- Then ---
		affirm.Nil(t, err)
		affirm.Equal(t, "A", target.Val)
	})

	t.Run("nil", func(t *testing.T) {
		// --- Given ---
		var target *types.TPtr

		// --- When ---
		err := ErrorAs(target, nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected non-nil error"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("error - target is a value", func(t *testing.T) {
		// --- Given ---
		var target types.TVal

		// --- When ---
		err := ErrorAs(&target, &types.TPtr{Val: "A"})

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected error to have a target in its tree:\n" +
			"  target: *types.TVal\n" +
			"   error: *types.TPtr"
		affirm.Equal(t, wMsg, err.Error())

		affirm.Equal(t, "", target.Val)
	})

	t.Run("error - target is a pointer", func(t *testing.T) {
		// --- Given ---
		var target *types.TPtr

		// --- When ---
		err := ErrorAs(&target, types.TVal{Val: "A"})

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"expected error to have a target in its tree:\n" +
			"  target: *types.TPtr\n" +
			"   error: types.TVal"
		affirm.Equal(t, wMsg, err.Error())

		affirm.Nil(t, target)
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")
		var target *types.TPtr

		// --- When ---
		err := ErrorAs(&target, errors.New("msg"), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error to have a target in its tree:\n" +
			"   trail: type.field\n" +
			"  target: *types.TPtr\n" +
			"   error: *errors.errorString"
		affirm.Equal(t, wMsg, err.Error())

		affirm.Nil(t, target)
	})
}

func Test_ErrorEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := ErrorEqual("e0", errors.New("e0"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("with format", func(t *testing.T) {
		// --- When ---
		err := ErrorEqual("e00", errors.New("e00"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("with percent in want", func(t *testing.T) {
		// --- When ---
		err := ErrorEqual("e0%", errors.New("e0%"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := ErrorEqual("e1", errors.New("e0"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to be:\n" +
			"  want: \"e1\"\n" +
			"  have: \"e0\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := ErrorEqual("e1", errors.New("e0"), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to be:\n" +
			"  trail: type.field\n" +
			"   want: \"e1\"\n" +
			"   have: \"e0\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		err := ErrorEqual("e1", nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to be:\n" +
			"  want: \"e1\"\n" +
			"  have: <nil>"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_ErrorContain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- When ---
		err := ErrorContain("def", errors.New("abc def ghi"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := ErrorContain("xyz", errors.New("abc def ghi"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to contain:\n" +
			"  want: \"xyz\"\n" +
			"  have: \"abc def ghi\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := ErrorContain("xyz", errors.New("abc def ghi"), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to contain:\n" +
			"  trail: type.field\n" +
			"   want: \"xyz\"\n" +
			"   have: \"abc def ghi\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		err := ErrorContain("xyz", nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error not to be nil:\n" +
			"  want: <non-nil>\n" +
			"  have: <nil>"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := ErrorContain("xyz", nil, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error not to be nil:\n" +
			"  trail: type.field\n" +
			"   want: <non-nil>\n" +
			"   have: <nil>"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil interface", func(t *testing.T) {
		// --- Given ---
		var e *types.TPtr

		// --- When ---
		err := ErrorContain("abc", e)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error not to be nil:\n" +
			"  want: <non-nil>\n" +
			"  have: *types.TPtr"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_ErrorRegexp(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		// --- When ---
		err := ErrorRegexp("def", errors.New("abc def ghi"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("regex", func(t *testing.T) {
		// --- Given ---
		rx := regexp.MustCompile(".* def .*")

		// --- When ---
		err := ErrorRegexp(rx, errors.New("abc def ghi"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("regex", func(t *testing.T) {
		// --- When ---
		err := ErrorRegexp("ghi$", errors.New("abc def ghi"))

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		err := ErrorRegexp("^xyz", errors.New("abc def ghi"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to match regexp:\n" +
			"  regexp: ^xyz\n" +
			"    have: \"abc def ghi\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := ErrorRegexp("^xyz", errors.New("abc def ghi"), opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to match regexp:\n" +
			"   trail: type.field\n" +
			"  regexp: ^xyz\n" +
			"    have: \"abc def ghi\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("invalid regexp", func(t *testing.T) {
		// --- When ---
		err := ErrorRegexp("[a-z", errors.New("abc def ghi"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected the error message to match regexp:\n" +
			"  error: \"error parsing regexp: missing closing ]: `[a-z`\""
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		err := ErrorRegexp("xyz", nil)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error not to be nil:\n" +
			"  want: <non-nil>\n" +
			"  have: <nil>"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := ErrorRegexp("^ab", nil, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error not to be nil:\n" +
			"  trail: type.field\n" +
			"   want: <non-nil>\n" +
			"   have: <nil>"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil interface", func(t *testing.T) {
		// --- Given ---
		var e *types.TPtr

		// --- When ---
		err := ErrorRegexp("^ab", e)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected error not to be nil:\n" +
			"  want: <non-nil>\n" +
			"  have: *types.TPtr"
		affirm.Equal(t, wMsg, err.Error())
	})
}
