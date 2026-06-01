// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package iokit

// Option configures error-injecting readers and writers (see With*Err functions).
type Option func(*Options)

// WithReadErr returns an [Option] that sets a custom error to be returned
// from read operations once the byte limit is reached.
func WithReadErr(err error) Option {
	return func(opts *Options) { opts.errRead = err }
}

// WithSeekErr returns an [Option] that sets a custom error to be returned
// from seek operations.
func WithSeekErr(err error) Option {
	return func(opts *Options) { opts.errSeek = err }
}

// WithWriteErr returns an [Option] that sets a custom error to be returned
// from write operations once the byte limit is reached.
func WithWriteErr(err error) Option {
	return func(opts *Options) { opts.errWrite = err }
}

// WithCloseErr returns an [Option] that sets a custom error to be returned
// from Close (the underlying Close is still called).
func WithCloseErr(err error) Option {
	return func(opts *Options) { opts.errClose = err }
}

// Options holds the configurable error values for iokit error readers and
// writers. Use the With*Err functions to customize them.
type Options struct {
	errRead  error
	errSeek  error
	errWrite error
	errClose error
}

// defaultOptions returns the default Options used by ErrReader/ErrWriter
// constructors (ErrRead for reads, ErrWrite for writes, no close error).
func defaultOptions() *Options {
	return &Options{
		errRead:  ErrRead,
		errClose: nil,
		errWrite: ErrWrite,
	}
}
