// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_Mocker_Run(t *testing.T) {
	t.Run("todo", func(t *testing.T) {
		// --- Given ---
		buf := &bytes.Buffer{}
		opts := []Option{
			WithSrcDir("testdata/cases"),
			WithDstDir("testdata/cases"),
			WithOutput(buf),
		}
		act := must.Value(NewAction("Cases", opts...))
		mck := NewMocker()

		// --- When ---
		err := mck.Run(act)

		// --- Then ---
		assert.NoError(t, err)
		fmt.Println(buf.String())

		// TODO(rz): finish this test
		// TODO(rz): name this test.
	})
}
