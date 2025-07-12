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

func Test_WithTesterAlias(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		cfg := &Config{}

		// --- When ---
		WithTesterAlias("alias")(cfg)

		// --- Then ---
		assert.Equal(t, "alias", cfg.testerAlias)
	})

	t.Run("set empty string", func(t *testing.T) {
		// --- Given ---
		cfg := &Config{}

		// --- When ---
		WithTesterAlias("")(cfg)

		// --- Then ---
		assert.Equal(t, "_tester", cfg.testerAlias)
	})
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
			pkgName:  "mocker",
			pkgPath:  "github.com/ctx42/testing/pkg/mocker",
			pkgDir:   wd,
			modName:  "testing",
			modPath:  "github.com/ctx42/testing",
			modDir:   filepath.Join(wd, "../.."),
			wd:       wd,
			resolved: true,
		}
		assert.Equal(t, wSrcImp, have.srcPkg)
		assert.Nil(t, have.srcFile)
		assert.Nil(t, have.srcItf)

		assert.Equal(t, "TstItfMock", have.tgtName)
		wTgtImp := &gopkg{
			pkgName:  "mocker",
			pkgPath:  "github.com/ctx42/testing/pkg/mocker",
			pkgDir:   wd,
			modName:  "testing",
			modPath:  "github.com/ctx42/testing",
			modDir:   filepath.Join(wd, "../.."),
			wd:       wd,
			resolved: true,
		}
		assert.Equal(t, wTgtImp, have.tgtPkg)
		wFilePth := filepath.Join(wTgtImp.pkgDir, "tst_itf_mock.go")
		assert.Equal(t, wFilePth, have.tgtFilename)
		assert.Nil(t, have.tgtOut)
		assert.False(t, have.onHelpers)
		assert.Empty(t, have.testerAlias)
	})

	t.Run("with a source directory", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())

		// --- When ---
		have, err := newConfig("Case00", WithSrc("testdata/cases"))

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "Case00", have.srcName)
		wSrcImp := &gopkg{
			pkgName:  "cases",
			pkgPath:  "github.com/ctx42/testing/pkg/mocker/testdata/cases",
			pkgDir:   filepath.Join(wd, "testdata/cases"),
			modName:  "testing",
			modPath:  "github.com/ctx42/testing",
			modDir:   filepath.Join(wd, "../.."),
			wd:       filepath.Join(wd, "testdata/cases"),
			resolved: true,
		}
		assert.Equal(t, wSrcImp, have.srcPkg)
		assert.Nil(t, have.srcFile)
		assert.Nil(t, have.srcItf)

		assert.Equal(t, "Case00Mock", have.tgtName)
		wTgtImp := &gopkg{
			pkgName:  "mocker",
			pkgPath:  "github.com/ctx42/testing/pkg/mocker",
			pkgDir:   wd,
			modName:  "testing",
			modPath:  "github.com/ctx42/testing",
			modDir:   filepath.Join(wd, "../.."),
			wd:       wd,
			resolved: true,
		}
		assert.Equal(t, wTgtImp, have.tgtPkg)
		assert.Equal(t, filepath.Join(wd, "case00_mock.go"), have.tgtFilename)
		assert.Nil(t, have.tgtOut)
		assert.False(t, have.onHelpers)
		assert.Empty(t, have.testerAlias)
	})

	t.Run("with a custom target name", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		opt := WithTgtName("MyMock")

		// --- When ---
		have, err := newConfig("TstItf", opt)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "MyMock", have.tgtName)
		assert.Equal(t, filepath.Join(wd, "my_mock.go"), have.tgtFilename)
	})

	t.Run("interface name has all capital letters", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())

		// --- When ---
		have, err := newConfig("DB")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "DBMock", have.tgtName)
		assert.Equal(t, filepath.Join(wd, "db_mock.go"), have.tgtFilename)
	})

	t.Run("with a custom target filename", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		opt := WithTgtFilename("my_super_mock.go")

		// --- When ---
		have, err := newConfig("TstItf", opt)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "TstItfMock", have.tgtName)
		assert.Equal(t, filepath.Join(wd, "my_super_mock.go"), have.tgtFilename)
	})

	t.Run("error - interface name is required", func(t *testing.T) {
		// --- When ---
		_, err := newConfig("")

		// --- Then ---
		assert.ErrorEqual(t, "interface name is required for mocking", err)
	})

	t.Run("error - invalid source", func(t *testing.T) {
		// --- When ---
		_, err := newConfig("TstItf", WithSrc("!!!"))

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
	})

	t.Run("error - invalid target", func(t *testing.T) {
		// --- When ---
		_, err := newConfig("TstItf", WithTgt("!!!"))

		// --- Then ---
		assert.ErrorIs(t, ErrUnkPkg, err)
	})

	t.Run("error - cannot set filename and output", func(t *testing.T) {
		// --- When ---
		buf := &bytes.Buffer{}
		opts := []Option{WithTgtFilename("file.go"), WithTgtOutput(buf)}

		// --- When ---
		_, err := newConfig("TstItf", opts...)

		// --- Then ---
		assert.ErrorContain(t, "cannot use both", err)
	})
}

func Test_Config_create(t *testing.T) {
	t.Run("buffer already set", func(t *testing.T) {
		// --- Given ---
		buf := &bytes.Buffer{}
		pth := filepath.Join(t.TempDir(), "_target_.go")
		cfg := &Config{
			tgtOut:      buf,
			tgtFilename: pth,
		}

		// --- When ---
		have, err := cfg.create()

		// --- Then ---
		assert.NoError(t, err)
		assert.NoFileExist(t, pth)
		assert.Same(t, buf, have.tgtOut)
	})

	t.Run("path is not absolute", func(t *testing.T) {
		// --- Given ---
		buf := &bytes.Buffer{}
		cfg := &Config{
			tgtOut:      buf,
			tgtFilename: "_target_.go",
		}

		// --- When ---
		have, err := cfg.create()

		// --- Then ---
		assert.NoError(t, err)
		assert.NoFileExist(t, "_target_.go")
		assert.Same(t, buf, have.tgtOut)
	})

	t.Run("file created", func(t *testing.T) {
		// --- Given ---
		pth := filepath.Join(t.TempDir(), "_target_.go")
		cfg := &Config{tgtFilename: pth}

		// --- When ---
		have, err := cfg.create()

		// --- Then ---
		assert.NoError(t, err)
		assert.FileExist(t, pth)
		assert.NotNil(t, have.tgtOut)
		assert.Type(t, &os.File{}, have.tgtOut)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		pth := filepath.Join(t.TempDir(), "dir", "_target_.go")
		cfg := &Config{tgtFilename: pth}

		// --- When ---
		have, err := cfg.create()

		// --- Then ---
		assert.ErrorContain(t, "no such file or directory", err)
		assert.Nil(t, have.tgtOut)
	})
}
