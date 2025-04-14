// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Interface_Name(t *testing.T) {
	// --- Given ---
	itf := &Interface{name: "MyInterface"}

	// --- When ---
	have := itf.Name()

	// --- Then ---
	assert.Equal(t, "MyInterface", have)
}
