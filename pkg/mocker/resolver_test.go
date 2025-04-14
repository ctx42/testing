// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_resolver(t *testing.T) {
	t.Run("package already in the cache", func(t *testing.T) {
		// --- Given ---
		res := &resolver{
			cache: []*gopkg{
				{
					pkgPath: "github.com/ctx42/testing/pkg/mocker",
					wd:      "/dir",
				},
			},
		}
		pkg := &gopkg{pkgPath: "github.com/ctx42/testing/pkg/mocker"}

		// --- When ---
		err := res.resolve(pkg)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, res.cache[0], pkg)
		assert.NotSame(t, res.cache[0], pkg)
	})

	t.Run("package not in the cache", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		res := &resolver{}
		pkg := &gopkg{
			pkgPath: "github.com/ctx42/testing/pkg/mocker",
			wd:      wd,
		}

		// --- When ---
		err := res.resolve(pkg)

		// --- Then ---
		assert.NoError(t, err)
		want := &gopkg{
			alias:    "",
			pkgName:  "mocker",
			pkgPath:  "github.com/ctx42/testing/pkg/mocker",
			pkgDir:   wd,
			modName:  "testing",
			modPath:  "github.com/ctx42/testing",
			modDir:   filepath.Join(wd, "../.."),
			wd:       wd,
			resolved: true,
		}
		assert.Equal(t, want, pkg)
		assert.Same(t, res.cache[0], pkg)
	})

	t.Run("already resolved are returned right away", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		res := &resolver{}
		pkg := &gopkg{
			pkgPath:  "github.com/ctx42/testing/pkg/mocker",
			resolved: true,
			wd:       wd,
		}

		// --- When ---
		err := res.resolve(pkg)

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("error package not found", func(t *testing.T) {
		// --- Given ---
		res := &resolver{}
		pkg := &gopkg{
			pkgPath: "github.com/ctx42/testing/pkg/abc",
			wd:      "/dir",
		}

		// --- When ---
		err := res.resolve(pkg)

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
	})
}
