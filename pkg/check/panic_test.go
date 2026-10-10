// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"regexp"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Panic(t *testing.T) {
	t.Run("panicked", func(t *testing.T) {
		// --- When ---
		err := Panic(func() { panic("test") })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - not panicked", func(t *testing.T) {
		// --- When ---
		err := Panic(func() {})

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "func should panic", err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := Panic(func() {}, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "" +
			"func should panic:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NoPanic(t *testing.T) {
	t.Run("error - panicked", func(t *testing.T) {
		// --- When ---
		err := NoPanic(func() { panic("test") })

		// --- Then ---
		affirm.NotNil(t, err)
		wRx := `(?s)^func should not panic:.*panic value: test.*panic stack:`
		affirm.Equal(t, true, regexp.MustCompile(wRx).MatchString(err.Error()))
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		err := NoPanic(func() { panic("test") }, WithTrail("type.field"))

		// --- Then ---
		affirm.NotNil(t, err)
		wRx := "" +
			`(?s)^func should not panic:\n {8}trail: type\.field\n` +
			`.*panic value: test.*panic stack:`
		affirm.Equal(t, true, regexp.MustCompile(wRx).MatchString(err.Error()))
	})

	t.Run("no panic", func(t *testing.T) {
		// --- When ---
		err := NoPanic(func() {})

		// --- Then ---
		affirm.Nil(t, err)
	})
}

func Test_PanicContain(t *testing.T) {
	t.Run("error - no panic with trail", func(t *testing.T) {
		// --- When ---
		err := PanicContain("abc", func() {}, WithTrail("type.field"))

		// --- Then ---
		wMsg := "" +
			"func should panic:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("panic message contains", func(t *testing.T) {
		// --- When ---
		err := PanicContain("def", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - panic message does not contain", func(t *testing.T) {
		// --- When ---
		err := PanicContain("xyz", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.NotNil(t, err)
		wRx := "" +
			`(?s)^func should panic with string containing:.*` +
			`substring: "xyz".*panic value: abc def ghi.*panic stack:`
		affirm.Equal(t, true, regexp.MustCompile(wRx).MatchString(err.Error()))
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := PanicContain("xyz", func() { panic("abc def ghi") }, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wRx := "" +
			`(?s)^func should panic with string containing:.*` +
			`trail: type\.field.*substring: "xyz".*` +
			`panic value: abc def ghi.*panic stack:`
		affirm.Equal(t, true, regexp.MustCompile(wRx).MatchString(err.Error()))
	})

	t.Run("panics with error type", func(t *testing.T) {
		// --- When ---
		err := PanicContain("def", func() { panic(errors.New("abc def ghi")) })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("panics with other type", func(t *testing.T) {
		// --- Given ---
		v := testcases.TInt{V: 42}

		// --- When ---
		err := PanicContain("{42}", func() { panic(v) })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error - does not panic", func(t *testing.T) {
		// --- When ---
		err := PanicContain("xyz", func() {})

		// --- Then ---
		affirm.NotNil(t, err)
		affirm.Equal(t, "func should panic", err.Error())
	})
}

func Test_PanicMsg(t *testing.T) {
	t.Run("panics", func(t *testing.T) {
		// --- When ---
		have, err := PanicMsg(func() { panic("abc def ghi") })

		// --- Then ---
		affirm.Nil(t, err)
		affirm.NotNil(t, have)
		affirm.Equal(t, "abc def ghi", *have)
	})

	t.Run("error - no panic", func(t *testing.T) {
		// --- When ---
		have, err := PanicMsg(func() {})

		// --- Then ---
		affirm.Nil(t, have)
		affirm.Equal(t, "func should panic", err.Error())
	})

	t.Run("error - with trail", func(t *testing.T) {
		// --- When ---
		have, err := PanicMsg(func() {}, WithTrail("type.field"))

		// --- Then ---
		affirm.Nil(t, have)
		wMsg := "" +
			"func should panic:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("panics with error type", func(t *testing.T) {
		// --- When ---
		have, err := PanicMsg(func() { panic(errors.New("abc")) })

		// --- Then ---
		affirm.Nil(t, err)
		affirm.NotNil(t, have)
		affirm.Equal(t, "abc", *have)
	})

	t.Run("panics with other type", func(t *testing.T) {
		// --- Given ---
		v := testcases.TInt{V: 42}

		// --- When ---
		have, err := PanicMsg(func() { panic(v) })

		// --- Then ---
		affirm.Nil(t, err)
		affirm.NotNil(t, have)
		affirm.Equal(t, "{42}", *have)
	})
}
