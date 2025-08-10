// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/types"
)

func Test_Panic(t *testing.T) {
	t.Run("panicked", func(t *testing.T) {
		// --- When ---
		err := Panic(func() { panic("test") })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("not panicked", func(t *testing.T) {
		// --- When ---
		err := Panic(func() {})

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "func should panic"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Panic(func() {}, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "func should panic:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_NoPanic(t *testing.T) {
	t.Run("panicked", func(t *testing.T) {
		// --- When ---
		err := NoPanic(func() { panic("test") })

		// --- Then ---
		affirm.NotNil(t, err)

		hMsg := err.Error()
		affirm.Equal(t, true, strings.Contains(hMsg, "func should not panic"))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic value: test"))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic stack:"))
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := NoPanic(func() { panic("test") }, opt)

		// --- Then ---
		affirm.NotNil(t, err)

		hMsg := err.Error()
		affirm.Equal(t, true, strings.Contains(hMsg, "func should not panic"))
		affirm.Equal(t, true, strings.Contains(hMsg, "        trail: type.field"))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic value: test"))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic stack:"))
	})

	t.Run("no panic", func(t *testing.T) {
		// --- When ---
		err := NoPanic(func() {})

		// --- Then ---
		affirm.Nil(t, err)
	})
}

func Test_PanicContain(t *testing.T) {
	t.Run("panic message contains", func(t *testing.T) {
		// --- When ---
		err := PanicContain("def", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("panic message does not contain", func(t *testing.T) {
		// --- When ---
		err := PanicContain("xyz", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "func should panic with string containing:"
		hMsg := err.Error()
		affirm.Equal(t, true, strings.Contains(hMsg, wMsg))
		affirm.Equal(t, true, strings.Contains(hMsg, "substring: \"xyz\""))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic value: abc def ghi"))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic stack:"))
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := PanicContain("xyz", func() { panic("abc def ghi") }, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "func should panic with string containing:"
		hMsg := err.Error()
		affirm.Equal(t, true, strings.Contains(hMsg, wMsg))
		affirm.Equal(t, true, strings.Contains(hMsg, "trail: type.field"))
		affirm.Equal(t, true, strings.Contains(hMsg, "substring: \"xyz\""))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic value: abc def ghi"))
		affirm.Equal(t, true, strings.Contains(hMsg, "panic stack:"))
	})

	t.Run("panics with error type", func(t *testing.T) {
		// --- When ---
		err := PanicContain("def", func() { panic(errors.New("abc def ghi")) })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("panics with other type", func(t *testing.T) {
		// --- Given ---
		v := types.TInt{V: 42}

		// --- When ---
		err := PanicContain("{42}", func() { panic(v) })

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("does not panic", func(t *testing.T) {
		// --- When ---
		err := PanicContain("xyz", func() {})

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "func should panic"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_PanicMsg(t *testing.T) {
	t.Run("panics", func(t *testing.T) {
		// --- When ---
		msg, err := PanicMsg(func() { panic("abc def ghi") })

		// --- Then ---
		affirm.Nil(t, err)
		if msg == nil {
			t.Error("expected PanicMsg to return non-nil value")
			return
		}
		affirm.Equal(t, "abc def ghi", *msg)
	})

	t.Run("no panic", func(t *testing.T) {
		// --- When ---
		msg, err := PanicMsg(func() {})

		// --- Then ---
		affirm.NotNil(t, err)
		if msg != nil {
			t.Error("expected PanicMsg to return nil value")
		}
		wMsg := "func should panic"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("additional message rows added", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		msg, err := PanicMsg(func() {}, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		if msg != nil {
			t.Error("expected PanicMsg to return nil value")
		}
		wMsg := "func should panic:\n" +
			"  trail: type.field"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("panics with error type", func(t *testing.T) {
		// --- When ---
		msg, err := PanicMsg(func() { panic(errors.New("abc")) })

		// --- Then ---
		affirm.Nil(t, err)
		affirm.NotNil(t, msg)
		affirm.Equal(t, "abc", *msg)
	})

	t.Run("panics with other type", func(t *testing.T) {
		// --- Given ---
		v := types.TInt{V: 42}

		// --- When ---
		msg, err := PanicMsg(func() { panic(v) })

		// --- Then ---
		affirm.Nil(t, err)
		affirm.NotNil(t, msg)
		affirm.Equal(t, "{42}", *msg)
	})
}
