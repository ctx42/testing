// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/must"
)

func Test_WithSrcDir(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithSrcDir("pth")(act)

	// --- Then ---
	want := &Action{Src: Import{Dir: "pth"}}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_WithSrcSpec(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithSrcSpec("github.com/ctx42/testing")(act)

	// --- Then ---
	want := &Action{
		Src: Import{
			Name: "testing",
			Spec: "github.com/ctx42/testing",
		},
	}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_WithDstDir(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithDstDir("dir")(act)

	// --- Then ---
	want := &Action{Dst: Import{Dir: "dir"}}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_WithDstSpec(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithDstSpec("github.com/ctx42/testing")(act)

	// --- Then ---
	want := &Action{
		Dst: Import{
			Name: "testing",
			Spec: "github.com/ctx42/testing",
		},
	}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_WithDstName(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithDstName("name")(act)

	// --- Then ---
	want := &Action{DstName: "name"}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_WithDstFilename(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithDstFilename("name")(act)

	// --- Then ---
	want := &Action{Filename: "name"}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_WithOutput(t *testing.T) {
	// --- Given ---
	buf := &bytes.Buffer{}
	act := &Action{}

	// --- When ---
	WithOutput(buf)(act)

	// --- Then ---
	want := &Action{Out: buf}
	assert.Equal(t, want, act, check.WithSkipUnexported)
	assert.Same(t, buf, act.Out)
}

func Test_WithOnHelpers(t *testing.T) {
	// --- Given ---
	act := &Action{}

	// --- When ---
	WithOnHelpers(act)

	// --- Then ---
	want := &Action{onHelpers: true}
	assert.Equal(t, want, act, check.WithSkipUnexported)
}

func Test_NewAction(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have, err := NewAction("MyItf")

		// --- Then ---
		assert.NoError(t, err)

		wd := must.Value(os.Getwd())
		want := &Action{
			SrcName: "MyItf",
			Src: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			DstName: "MyItfMock",
			Dst: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			Filename: filepath.Join(wd, "my_itf_mock.go"),
		}
		assert.Equal(t, want, have, check.WithSkipUnexported)
	})

	t.Run("with option", func(t *testing.T) {
		// --- When ---
		have, err := NewAction("MyItf", WithDstFilename("custom_mock.go"))

		// --- Then ---
		assert.NoError(t, err)

		wd := must.Value(os.Getwd())
		want := &Action{
			SrcName: "MyItf",
			Src: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			DstName: "MyItfMock",
			Dst: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			Filename: filepath.Join(wd, "custom_mock.go"),
		}
		assert.Equal(t, want, have, check.WithSkipUnexported)
	})

	t.Run("error", func(t *testing.T) {
		// --- When ---
		have, err := NewAction("MyItf", WithSrcSpec("!!!"))

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.Nil(t, have)
	})
}

func Test_Action_setup(t *testing.T) {
	t.Run("only itf name", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{
			SrcName:  "SomeItf",
			Src:      Import{},
			DstName:  "",
			Dst:      Import{},
			wd:       wd,
			Filename: "",
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.NoError(t, err)
		want := &Action{
			SrcName: "SomeItf",
			Src: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			DstName: "SomeItfMock",
			Dst: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			Filename: filepath.Join(wd, "some_itf_mock.go"),
		}
		assert.Equal(t, want, act, check.WithSkipUnexported)
	})

	t.Run("src dir wd, dst dir wd", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{
			SrcName:  "SomeItf",
			Src:      Import{Dir: wd},
			DstName:  "",
			Dst:      Import{Dir: wd},
			Filename: "",
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.NoError(t, err)
		want := &Action{
			SrcName: "SomeItf",
			Src: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			DstName: "SomeItfMock",
			Dst: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			Filename: filepath.Join(wd, "some_itf_mock.go"),
		}
		assert.Equal(t, want, act, check.WithSkipUnexported)
	})

	t.Run("src dir wd, dst spec", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{
			SrcName: "SomeItf",
			Src:     Import{Dir: wd},
			DstName: "",
			Dst: Import{
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/dest",
			},
			Filename: "",
			wd:       wd,
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.NoError(t, err)
		want := &Action{
			SrcName: "SomeItf",
			Src: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			DstName: "SomeItfMock",
			Dst: Import{
				Name: "dest",
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/dest",
				Dir:  filepath.Join(wd, "testdata/dest"),
			},
			Filename: filepath.Join(wd, "testdata/dest/some_itf_mock.go"),
		}
		assert.Equal(t, want, act, check.WithSkipUnexported)
	})

	t.Run("src spec, dst dir wd", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{
			SrcName: "SomeItf",
			Src: Import{
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/dest",
			},

			DstName:  "",
			Dst:      Import{},
			Filename: "",
			wd:       wd,
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.NoError(t, err)
		want := &Action{
			SrcName: "SomeItf",
			Src: Import{
				Name: "dest",
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/dest",
				Dir:  filepath.Join(wd, "testdata/dest"),
			},
			DstName: "SomeItfMock",
			Dst: Import{
				Name: "mocker",
				Spec: "github.com/ctx42/testing/pkg/mocker",
				Dir:  wd,
			},
			Filename: filepath.Join(wd, "some_itf_mock.go"),
		}
		assert.Equal(t, want, act, check.WithSkipUnexported)
	})

	t.Run("custom dest filename", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{
			Src: Import{
				Spec: "github.com/ctx42/testing/pkg/mocker",
			},

			Dst: Import{
				Spec: "github.com/ctx42/testing/pkg/mocker/testdata/dest",
			},
			Filename: "custom.go",
			wd:       wd,
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.NoError(t, err)
		want := filepath.Join(wd, "testdata/dest/custom.go")
		assert.Equal(t, want, act.Filename)
	})

	t.Run("custom dest type name", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{SrcName: "SomeItf", DstName: "CustomMock", wd: wd}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "CustomMock", act.DstName)
	})

	t.Run("error src spec", func(t *testing.T) {
		// --- Given ---
		act := &Action{
			Src: Import{Spec: "github.com/ctx42/testing/pkg/mocker/abc"},
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.ErrorContain(t, "no required module provides package", err)
	})

	t.Run("error dst spec", func(t *testing.T) {
		// --- Given ---
		wd := must.Value(os.Getwd())
		act := &Action{
			Dst: Import{Spec: "github.com/ctx42/testing/pkg/mocker/abc"},
			wd:  wd,
		}

		// --- When ---
		err := act.setup()

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvSpec)
		assert.ErrorContain(t, "no required module provides package", err)
	})
}
