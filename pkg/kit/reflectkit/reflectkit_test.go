package reflectkit

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_GetField(t *testing.T) {
	s := &struct {
		F int `url:"f"`
	}{}

	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := GetField(tspy, s, "F")

		// --- Then ---
		assert.NotNil(t, have)
	})

	t.Run("error - empty field name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogEqual("the struct field name must not be empty")
		tspy.Close()

		// --- When ---
		have := GetField(tspy, s, "")

		// --- Then ---
		assert.Zero(t, have)
	})

	t.Run("error - not existing field", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "struct `*struct { F int \"url:\\\"f\\\"\" }` has no field `A`"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := GetField(tspy, s, "A")

		// --- Then ---
		assert.Zero(t, have)
	})
}

func Test_GetValue(t *testing.T) {
	s := &struct{ F int }{}

	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		// --- When ---
		have := GetValue(tspy, s, "F")

		// --- Then ---
		assert.True(t, have.IsValid())
		have.Set(reflect.ValueOf(1))
		assert.Equal(t, 1, s.F)
	})

	t.Run("error - empty field name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogEqual("the struct field name must not be empty")
		tspy.Close()

		// --- When ---
		have := GetValue(tspy, s, "")

		// --- Then ---
		assert.False(t, have.IsValid())
	})

	t.Run("error - not existing field", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogEqual("cannot get value for `*struct { F int }.A`")
		tspy.Close()

		// --- When ---
		have := GetValue(tspy, s, "A")

		// --- Then ---
		assert.False(t, have.IsValid())
	})

	t.Run("error - not a struct", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogEqual("type `int` is not struct")
		tspy.Close()

		// --- When ---
		have := GetValue(tspy, 1, "A")

		// --- Then ---
		assert.False(t, have.IsValid())
	})
}
