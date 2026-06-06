// Package randkit provides cryptographically random test helpers for
// generating strings, file names, identifiers, and other test fixtures.
// All functions use [crypto/rand] and panic if randomness is unavailable.
package randkit

import (
	"crypto/rand"
	"math/big"
	"path/filepath"
	"strings"
)

// Uppercase is the list of uppercase letters.
const Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Lowercase is the list of lowercase letters.
const Lowercase = "abcdefghijklmnopqrstuvwxyz"

// Letters list of lowercase and uppercase letters.
const Letters = Lowercase + Uppercase

// Digits is a list of digits 0 to 9.
const Digits = "0123456789"

// WithChars is a [Str] option setting the list of characters to use when
// generating random strings. All passed strings are concatenated in the order
// they were passed.
func WithChars(list ...string) func(*options) {
	return func(opts *options) {
		opts.chars = strings.Join(list, "")
	}
}

// WithPrefix is a [Str] option setting prefix for generated strings.
func WithPrefix(prefix string) func(*options) {
	return func(opt *options) { opt.prefix = prefix }
}

// WithSuffix is a [Str] option setting suffix for generated strings.
func WithSuffix(suffix string) func(*options) {
	return func(opt *options) { opt.suffix = suffix }
}

// WithExt is a [Str] option alias for [WithSuffix].
func WithExt(ext string) func(*options) { return WithSuffix(ext) }

// WithLen is a [Str] option setting the length for the generated string.
func WithLen(n int) func(*options) {
	return func(opt *options) { opt.n = n }
}

// options represents options for generating random strings.
type options struct {
	chars  string // List of characters to use for random string.
	n      int    // Length of the random string.
	prefix string // Random string prefix.
	suffix string // Random string suffix.
}

// Str returns a random string based on provided options. When no options are
// given, the generated string will be 10 characters long containing only
// letters.
func Str(opts ...func(*options)) string {
	def := options{
		chars:  Letters,
		n:      10,
		prefix: "",
		suffix: "",
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def.prefix + randStr(def.chars, def.n) + def.suffix
}

// FileName returns a random file name. By default, the file name is 7
// letters long [a-zA-Z] with the prefix "file-" and extension ".txt".
func FileName(dir string, opts ...func(*options)) string {
	def := options{
		chars:  Letters,
		n:      7,
		prefix: "file-",
		suffix: ".txt",
	}
	for _, opt := range opts {
		opt(&def)
	}
	name := def.prefix + randStr(def.chars, def.n) + def.suffix
	return filepath.Join(dir, name)
}

// Int generates a random integer in the range [1, max].
func Int(maximum int) int {
	return intn(maximum) + 1
}

// Password returns an n-character random password drawn from letters
// [a-zA-Z] and digits [0-9]. No special characters are included.
func Password(n int) string {
	return Str(WithChars(Letters, Digits), WithLen(n))
}

// randStr returns n random characters drawn from chars.
func randStr(chars string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[intn(len(chars))]
	}
	return string(b)
}

// intn returns a uniform random value in [0, max). It panics if max <= 0
// or when a random number could not be generated.
func intn(maximum int) int {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(maximum)))
	if err != nil {
		panic(err)
	}
	return int(val.Int64())
}
