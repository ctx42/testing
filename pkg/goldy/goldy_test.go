package goldy

import (
	"strings"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/internal/core"
)

func Test_Text(t *testing.T) {
	t.Run("success case 1", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Text(tspy, "testdata/text_case1.txt")

		// --- Then ---
		affirm.Equal(t, "Content #1.\nContent #2.", have)
	})

	t.Run("success case 2", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Text(tspy, "testdata/text_case2.txt")

		// --- Then ---
		affirm.Equal(t, "Content #1.\nContent #2.\n", have)
	})

	t.Run("success case 3", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy()

		// --- When ---
		have := Text(tspy, "testdata/text_case3.txt")

		// --- Then ---
		affirm.Equal(t, "Content #1.\nContent #2.\n\n", have)
	})

	t.Run("error no marker", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := Text(tspy, "testdata/text_no_marker.txt")

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "golden file is missing \"---\" marker"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})

	t.Run("not existing file", func(t *testing.T) {
		// --- Given ---
		tspy := core.NewSpy().Capture()

		// --- When ---
		have := Text(tspy, "testdata/not-existing.txt")

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.Equal(t, true, tspy.Failed())
		wMsg := "no such file or directory"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
		wMsg = "testdata/not-existing.txt"
		affirm.Equal(t, true, strings.Contains(tspy.Log(), wMsg))
	})
}
