package iokit

// Option represents an option function.
type Option func(*Options)

// WithReadErr is an [Option] setting custom read error.
func WithReadErr(err error) Option {
	return func(opts *Options) { opts.errRead = err }
}

// WithSeekErr is an [option] setting custom seek error.
func WithSeekErr(err error) Option {
	return func(opts *Options) { opts.errSeek = err }
}

// WithWriteErr is an [option] setting custom write error.
func WithWriteErr(err error) Option {
	return func(opts *Options) { opts.errWrite = err }
}

// WithCloseErr is an [option] setting custom close error.
func WithCloseErr(err error) Option {
	return func(opts *Options) { opts.errClose = err }
}

// Options represent options used by iokit tools.
type Options struct {
	errRead  error // Read error.
	errSeek  error // Seek error.
	errWrite error // Write error.
	errClose error // Close error.
}

// defaultOptions returns default options.
func defaultOptions() *Options {
	return &Options{
		errRead:  ErrRead,
		errClose: nil,
		errWrite: ErrWrite,
	}
}
