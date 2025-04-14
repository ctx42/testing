// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_WithSrc(t *testing.T) {
	// --- Given ---
	cfg := &Config{}

	// --- When ---
	WithSrc("pth")(cfg)

	// --- Then ---
	assert.Equal(t, "pth", cfg.srcDirOrImp)
}

func Test_WithTgt(t *testing.T) {
	// --- Given ---
	cfg := &Config{}

	// --- When ---
	WithTgt("dir")(cfg)

	// --- Then ---
	assert.Equal(t, "dir", cfg.tgtDirOrImp)
}

func Test_WithTgtName(t *testing.T) {
	// --- Given ---
	cfg := &Config{}

	// --- When ---
	WithTgtName("name")(cfg)

	// --- Then ---
	assert.Equal(t, "name", cfg.tgtName)
}

func Test_WithTgtFilename(t *testing.T) {
	// --- Given ---
	cfg := &Config{}

	// --- When ---
	WithTgtFilename("file.go")(cfg)

	// --- Then ---
	assert.Equal(t, "file.go", cfg.tgtFilename)
}

func Test_WithTgtOutput(t *testing.T) {
	// --- Given ---
	buf := &bytes.Buffer{}
	cfg := &Config{}

	// --- When ---
	WithTgtOutput(buf)(cfg)

	// --- Then ---
	assert.Same(t, buf, cfg.tgtOut)
}

func Test_WithTgtOnHelpers(t *testing.T) {
	// --- Given ---
	cfg := &Config{}

	// --- When ---
	WithTgtOnHelpers(cfg)

	// --- Then ---
	assert.True(t, cfg.onHelpers)
}

func Test_newConfig(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())

		// --- When ---
		have, err := newConfig("TstItf")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "TstItf", have.srcName)
		wSrcImp := &gopkg{
			pkgName: "mocker",
			pkgPath: "github.com/ctx42/testing/pkg/mocker",
			pkgDir:  wd,
			modName: "testing",
			modPath: "github.com/ctx42/testing",
			modDir:  filepath.Join(wd, "../.."),
			wd:      wd,
		}
		assert.Equal(t, wSrcImp, have.srcPkg)

		assert.Equal(t, "TstItfMock", have.tgtName)
		wTgtImp := &gopkg{
			pkgName: "mocker",
			pkgPath: "github.com/ctx42/testing/pkg/mocker",
			pkgDir:  wd,
			modName: "testing",
			modPath: "github.com/ctx42/testing",
			modDir:  filepath.Join(wd, "../.."),
			wd:      wd,
		}
		assert.Equal(t, wTgtImp, have.tgtPkg)
		assert.Equal(t, "tst_itf_mock.go", have.tgtFilename)
		assert.Nil(t, have.tgtOut)
		assert.False(t, have.onHelpers)
	})

	t.Run("with a custom target name", func(t *testing.T) {
		// --- Given ---
		opt := WithTgtName("MyMock")

		// --- When ---
		have, err := newConfig("TstItf", opt)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "MyMock", have.tgtName)
		assert.Equal(t, "my_mock.go", have.tgtFilename)
	})

	t.Run("with a custom target filename", func(t *testing.T) {
		// --- Given ---
		opt := WithTgtFilename("my_super_mock.go")

		// --- When ---
		have, err := newConfig("TstItf", opt)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "TstItfMock", have.tgtName)
		assert.Equal(t, "my_super_mock.go", have.tgtFilename)
	})

	t.Run("error interface name is required", func(t *testing.T) {
		// --- When ---
		_, err := newConfig("")

		// --- Then ---
		assert.ErrorEqual(t, "interface name for mocking is required", err)
	})

	t.Run("error invalid source", func(t *testing.T) {
		// --- When ---
		_, err := newConfig("TstItf", WithSrc("!!!"))

		// --- Then ---
		assert.ErrorIs(t, err, ErrUnkPkg)
	})

	t.Run("error invalid target", func(t *testing.T) {
		// --- When ---
		_, err := newConfig("TstItf", WithTgt("!!!"))

		// --- Then ---
		assert.ErrorIs(t, err, ErrUnkPkg)
	})

	t.Run("error cannot set filename and output", func(t *testing.T) {
		// --- When ---
		buf := &bytes.Buffer{}
		opts := []Option{WithTgtFilename("file.go"), WithTgtOutput(buf)}

		// --- When ---
		_, err := newConfig("TstItf", opts...)

		// --- Then ---
		assert.ErrorContain(t, "cannot use both", err)
	})
}
