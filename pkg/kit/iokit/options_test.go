// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package iokit

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithReadErr(t *testing.T) {
	// --- Given ---
	err := errors.New("test error")
	opts := &Options{}

	// --- When ---
	WithReadErr(err)(opts)

	// --- Then ---
	assert.Same(t, err, opts.errRead)
}

func Test_WithSeekErr(t *testing.T) {
	// --- Given ---
	err := errors.New("test error")
	opts := &Options{}

	// --- When ---
	WithSeekErr(err)(opts)

	// --- Then ---
	assert.Same(t, err, opts.errSeek)
}

func Test_WithWriteErr(t *testing.T) {
	// --- Given ---
	err := errors.New("test error")
	opts := &Options{}

	// --- When ---
	WithWriteErr(err)(opts)

	// --- Then ---
	assert.Same(t, err, opts.errWrite)
}

func Test_WithCloseErr(t *testing.T) {
	// --- Given ---
	err := errors.New("test error")
	opts := &Options{}

	// --- When ---
	WithCloseErr(err)(opts)

	// --- Then ---
	assert.Same(t, err, opts.errClose)
}

func Test_defaultOptions(t *testing.T) {
	// --- When ---
	have := defaultOptions()

	// --- Then ---
	want := &Options{
		errRead:  ErrRead,
		errWrite: ErrWrite,
	}
	assert.Equal(t, want, have)
}
