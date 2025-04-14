// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// Option represents [NewAction] option.
type Option func(*Action)

// WithSrcDir is [Option] setting directory with an interface to mock.
func WithSrcDir(dir string) Option {
	return func(act *Action) { act.Src.Dir = dir }
}

// WithSrcSpec is [Option] setting import spec with the interface to mock.
func WithSrcSpec(spec string) Option {
	return func(mck *Action) { mck.Src = mck.Src.SetSpec(spec) }
}

// WithDstDir is [Option] setting the destination directory for the mock.
func WithDstDir(dir string) Option {
	return func(act *Action) { act.Dst.Dir = dir }
}

// WithDstSpec is [Option] setting the destination import spec for the mock.
func WithDstSpec(spec string) Option {
	return func(act *Action) { act.Dst = act.Dst.SetSpec(spec) }
}

// WithDstName is [Option] setting the type name implementing the interface.
func WithDstName(name string) Option {
	return func(act *Action) { act.DstName = name }
}

// WithDstFilename is [Option] setting the destination filename for the mock.
func WithDstFilename(filename string) Option {
	return func(act *Action) { act.Filename = filename }
}

// WithOutput is [Option] setting the destination buffer for the generated mock.
// It overrides the [WithDstFilename] option. If buffer implements [io.Close]
// interface, it will be called.
func WithOutput(buf io.Writer) Option {
	return func(act *Action) { act.Out = buf }
}

// WithOnHelpers is [Option] turning on "OnXXX" helper methods generation for
// every interface method being mocked.
func WithOnHelpers(mck *Action) { mck.onHelpers = true }

// Action represents one mocking action.
type Action struct {
	SrcName   string    // Name of the interface to mock.
	Src       Import    // Package the interface is defined.
	DstName   string    // The type name implementing the interface.
	Dst       Import    // Package to create the mock in.
	Out       io.Writer // Destination to write generated mock to.
	Filename  string    // Absolute path to the destination file.
	wd        string    // Working directory.
	onHelpers bool      // Generate "OnXXX" helper methods.
}

// NewAction creates a new [Action] with given options.
func NewAction(name string, opts ...Option) (*Action, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	act := &Action{SrcName: name, wd: wd}
	for _, opt := range opts {
		opt(act)
	}
	if err = act.setup(); err != nil {
		return nil, err
	}
	return act, nil
}

// setup completes and validates the fields of the Action based on its partial
// configuration. It populates missing or default values for required fields and
// checks that the configuration is valid for the action’s intended operation.
func (act *Action) setup() (err error) {
	if act.Src.IsZero() {
		act.Src.Dir = act.wd
	}

	if act.Dst.IsZero() {
		act.Dst.Dir = act.wd
	}

	if act.Src, err = settleImport(act.Src); err != nil {
		return err
	}

	if act.Dst, err = settleImport(act.Dst); err != nil {
		return err
	}

	var dstFile = act.Filename
	if dstFile == "" {
		dstFile = toLowerSnakeCase(act.SrcName) + "_mock.go"
	}
	act.Filename = path.Join(act.Dst.Dir, dstFile)
	if act.Filename, err = filepath.Abs(act.Filename); err != nil {
		return err
	}

	// If a destination type for the mock was not set, it will be based on
	// the name of the interface being mocked.
	if act.DstName == "" {
		act.DstName = act.SrcName + "Mock"
	}

	return nil
}
