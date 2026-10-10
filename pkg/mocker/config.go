// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Option represents a configuration option for [Mocker] or the package-level
// [Generate] function.
type Option func(*Config)

// WithSrc sets the source directory or import path where the interface to
// mock is defined. Required when the interface is not in the current package.
func WithSrc(dirOrImp string) Option {
	return func(cfg *Config) { cfg.srcDirOrImp = dirOrImp }
}

// WithTgt sets the target directory or import path for the generated mock
// file. Defaults to the current package.
func WithTgt(dirOrImp string) Option {
	return func(cfg *Config) { cfg.tgtDirOrImp = dirOrImp }
}

// WithTgtName sets a custom name for the generated mock type.
// Defaults to "<Interface>Mock".
func WithTgtName(name string) Option {
	return func(cfg *Config) { cfg.tgtName = name }
}

// WithTgtFilename sets a custom filename for the generated mock.
// Ignored if [WithTgtOutput] is used.
func WithTgtFilename(filename string) Option {
	return func(cfg *Config) { cfg.tgtFilename = filename }
}

// WithTgtOnHelpers enables generation of "OnXXX" helper methods on the mock
// (in addition to the standard recorder methods).
func WithTgtOnHelpers(cfg *Config) { cfg.onHelpers = true }

// WithTgtOutput configures a custom writer for the generated mock output.
// Takes precedence over [WithTgtFilename]. If the writer implements
// [io.Closer], it will be closed after writing.
func WithTgtOutput(w io.Writer) Option {
	return func(cfg *Config) { cfg.tgtOut = w }
}

// WithTesterAlias sets the alias for the tester import in the generated file.
// An empty alias selects "_tester". Without this option the import is not
// aliased and the package name "tester" is used. The alias must be a valid
// identifier other than "_".
func WithTesterAlias(alias string) Option {
	if alias == "" {
		alias = "_tester"
	}
	return func(cfg *Config) { cfg.testerAlias = alias }
}

// Config holds the configuration for generating a mock.
// It is usually created internally via options rather than directly by users.
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

// newConfig creates a validated configuration for mocking the given
// interface. Use the Option functions to customize source/target locations,
// output destination, mock name, etc.
//
// This is an internal helper; users normally use [Generate] or
// [Mocker.Generate] with options.
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

	// Check the options before running go list to resolve the packages.
	if cfg.tgtOut != nil && cfg.tgtFilename != "" {
		msg := "cannot use both WithTgtOutput and WithTgtFilename options"
		return Config{}, errors.New(msg)
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

	if alias := cfg.testerAlias; alias != "" &&
		(!token.IsIdentifier(alias) || alias == "_") {

		return Config{}, fmt.Errorf("invalid tester alias: %q", alias)
	}

	if cfg.tgtOut == nil {
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

// create opens the target file for writing when no output writer is set.
// Call only after the mock source is fully generated so a prior failure cannot
// wipe an existing file.
func (cfg Config) create() (Config, error) {
	if cfg.tgtOut != nil {
		return cfg, nil
	}
	if !filepath.IsAbs(cfg.tgtFilename) {
		format := "target file path is not absolute: %s"
		return cfg, fmt.Errorf(format, cfg.tgtFilename)
	}
	fName := cfg.tgtFilename
	fMode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	// G304: output path comes from trusted mocker configuration.
	out, err := os.OpenFile(fName, fMode, 0644) //nolint:gosec
	if err != nil {
		return cfg, err
	}
	cfg.tgtOut = out
	return cfg, nil
}
