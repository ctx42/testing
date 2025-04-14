// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"errors"
	"io"
	"os"
)

// Option represents a mocker option.
type Option func(*Config)

// WithSrc sets the source directory or import path (package) where the
// interface to mock is defined.
func WithSrc(dirOrImp string) Option {
	return func(cfg *Config) { cfg.srcDirOrImp = dirOrImp }
}

// WithTgt sets the target directory or import path (package) where the
// interface to mock should be created.
func WithTgt(dirOrImp string) Option {
	return func(cfg *Config) { cfg.tgtDirOrImp = dirOrImp }
}

// WithTgtName sets the interface mock type name.
func WithTgtName(name string) Option {
	return func(cfg *Config) { cfg.tgtName = name }
}

// WithTgtFilename sets the filename to write the generated interface mock to.
func WithTgtFilename(filename string) Option {
	return func(cfg *Config) { cfg.tgtFilename = filename }
}

// WithTgtOnHelpers turns on "OnXXX" helper methods generation.
func WithTgtOnHelpers(cfg *Config) { cfg.onHelpers = true }

// WithTgtOutput configures the writer for the generated interface output. It
// takes precedence over the [WithTgtFilename] option. If the provided writer
// implements [io.Closer], its Close method will be called after writing.
func WithTgtOutput(w io.Writer) Option {
	return func(cfg *Config) { cfg.tgtOut = w }
}

// Config represents the configuration for the mocker.
type Config struct {
	srcName     string // Name of the interface to mock.
	srcDirOrImp string // Directory or import path the interface is defined in.
	srcPkg      *gopkg // Source package (based on srcDirOrImp field).

	tgtName     string    // Custom name for the mock type.
	tgtDirOrImp string    // Target directory or import path.
	tgtFilename string    // Custom filename for the mock.
	tgtOut      io.Writer // Target to write generated mock to.
	tgtPkg      *gopkg    // Destination package (based on tgtDirOrImp field).

	onHelpers bool // Generate "OnXXX" helper methods.
}

// newConfig creates a new configuration for the interface with the provided
// name to be mocked. By default, the interface is expected to be in the
// current working directory, and the mock with name based on the source
// interface name is written to the current working directory. Use options to
// change the default behavior.
func newConfig(name string, opts ...Option) (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	if name == "" {
		return Config{}, errors.New("interface name for mocking is required")
	}

	cfg := Config{
		srcName: name,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.srcPkg, err = newPkg(wd, cfg.srcDirOrImp); err != nil {
		return Config{}, err
	}
	if cfg.tgtPkg, err = newPkg(wd, cfg.tgtDirOrImp); err != nil {
		return Config{}, err
	}
	if cfg.tgtName == "" {
		cfg.tgtName = cfg.srcName + "Mock"
	}
	if cfg.tgtOut != nil && cfg.tgtFilename != "" {
		msg := "cannot use both WithTgtOutput and WithTgtFilename options"
		return Config{}, errors.New(msg)
	}
	if cfg.tgtOut == nil && cfg.tgtFilename == "" {
		cfg.tgtFilename = toLowerSnakeCase(cfg.tgtName) + ".go"
	}
	return cfg, nil
}
