// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

import (
	"errors"
	"go/ast"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// WithTesterAlias sets the alias for the "github.com/ctx42/testing/pkg/tester"
// package. When the alias is set to empty string, it will use "_tester" as the
// alias.
func WithTesterAlias(alias string) Option {
	if alias == "" {
		alias = "_tester"
	}
	return func(cfg *Config) { cfg.testerAlias = alias }
}

// Config represents the configuration for the mocker.
type Config struct {
	srcName     string // Name of the interface to mock.
	srcDirOrImp string // Directory or import path the interface is defined in.
	srcPkg      *gopkg // Source package (based on srcDirOrImp field).
	srcFile     *file  // File with the interface to mock.

	srcItf *ast.InterfaceType // Interface to mock.

	tgtName     string    // Custom name for the mock type.
	tgtDirOrImp string    // Target directory or import path.
	tgtFilename string    // Custom filename for the mock.
	tgtOut      io.Writer // Target to write generated mock to.
	tgtPkg      *gopkg    // Destination package (based on tgtDirOrImp field).

	onHelpers   bool   // Generate "OnXXX" helper methods.
	testerAlias string // Alias for the CTX42 tester package.
}

// newConfig creates a new configuration for the interface with the provided
// name to be mocked. By default, the interface is expected to be in the
// current working directory, and the mock with name based on the source
// interface name is written to the file in the current working directory. Use
// options to change the default behavior.
//
// The function validates and sets fields and does not create the output file.
// Use the create method to create the output file.
func newConfig(name string, opts ...Option) (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	if name == "" {
		return Config{}, errors.New("interface name is required for mocking")
	}

	cfg := Config{srcName: name}
	for _, opt := range opts {
		opt(&cfg)
	}

	var srcWd string
	srcWd, cfg.srcDirOrImp = detectDirOrImp(wd, cfg.srcDirOrImp)
	cfg.srcPkg = newPkg(srcWd, cfg.srcDirOrImp)
	if err = cfg.srcPkg.resolve(); err != nil {
		return Config{}, err
	}

	var tgtWd string
	tgtWd, cfg.tgtDirOrImp = detectDirOrImp(wd, cfg.tgtDirOrImp)
	cfg.tgtPkg = newPkg(tgtWd, cfg.tgtDirOrImp)
	if err = cfg.tgtPkg.resolve(); err != nil {
		return Config{}, err
	}

	if cfg.tgtName == "" {
		cfg.tgtName = cfg.srcName + "Mock"
	}

	if cfg.tgtOut != nil {
		if cfg.tgtFilename != "" {
			msg := "cannot use both WithTgtOutput and WithTgtFilename options"
			return Config{}, errors.New(msg)
		}
	} else {
		if cfg.tgtFilename == "" {
			var tmp string
			if strings.HasSuffix(cfg.tgtName, "Mock") {
				tmp = cfg.tgtName[:len(cfg.tgtName)-4]
				cfg.tgtFilename = toLowerSnakeCase(tmp) + "_mock.go"
			} else {
				cfg.tgtFilename = toLowerSnakeCase(cfg.tgtName) + ".go"
			}
		}
		if !filepath.IsAbs(cfg.tgtFilename) {
			cfg.tgtFilename = filepath.Join(cfg.tgtPkg.pkgDir, cfg.tgtFilename)
		}
	}

	return cfg, nil
}

// create creates the target file if needed.
func (cfg Config) create() (Config, bool, error) {
	if cfg.tgtOut == nil && filepath.IsAbs(cfg.tgtFilename) {
		fName := cfg.tgtFilename
		fMode := os.O_RDWR | os.O_CREATE | os.O_TRUNC
		file, err := os.OpenFile(fName, fMode, 0664)
		if err != nil {
			return cfg, false, err
		}
		cfg.tgtOut = file
		_, err = file.WriteString("package " + cfg.tgtPkg.pkgName)
		if err != nil {
			_ = os.Remove(fName)
			return cfg, false, err
		}
		return cfg, true, nil
	}
	return cfg, false, nil
}
