package memfs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"syscall"
)

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// ErrOutOfBounds is returned for invalid offsets.
var ErrOutOfBounds = errors.New("offset out of bounds")

// errWriteAtInAppendMode returned when [File.WriteAt] is used with
// [os.O_APPEND] flag. It has the same message as the error with the same name
// in the os package.
var errWriteAtInAppendMode = errors.New("os: invalid use of WriteAt on file " +
	"opened with O_APPEND")

// WithFileOffset is a [File] constructor function option setting the offset.
func WithFileOffset(off int) func(*File) {
	return func(fil *File) { fil.off = off }
}

// WithFileAppend is a [File] constructor function option setting the offset to
// the end of the file. This option must be the last option on the option list.
//
// When [File.Truncate] is used, the offset will be set to the end of the [File].
func WithFileAppend(fil *File) { fil.flag |= os.O_APPEND }

// WithFileFlag is a [File] constructor function option setting flags. Flags
// are the same as for [os.OpenFile].
//
// Currently only [os.O_APPEND] flag is supported.
func WithFileFlag(flag int) func(*File) {
	return func(fil *File) { fil.flag = flag }
}

// WithFileName is a [File] constructor function option setting file name.
func WithFileName(name string) func(*File) {
	return func(fil *File) { fil.name = name }
}

// Compile time checks.
var (
	_ io.Seeker       = &File{}
	_ io.Reader       = &File{}
	_ io.ReaderAt     = &File{}
	_ io.Closer       = &File{}
	_ io.ReaderFrom   = &File{}
	_ io.Writer       = &File{}
	_ io.WriterAt     = &File{}
	_ io.StringWriter = &File{}
	_ io.WriterTo     = &File{}
	_ fs.File         = &File{}
)

// A File is a variable-sized buffer of bytes. Its zero value is an empty
// buffer ready to use.
type File struct {
	flag int    // Instance flags.
	off  int    // Current offset for read and write operations.
	name string // Optional file name.
	buf  []byte // Underlying buffer.
}

// NewFile returns a new instance of [File]. The difference between using this
// function and using the zero value is that this function will initialize the
// buffer with capacity of [bytes.MinRead]. It will panic with [ErrOutOfBounds]
// if the [WithFileOffset] option sets offset to a negative number or greater
// than [bytes.MinRead].
func NewFile(opts ...func(buffer *File)) *File {
	return FileWith(make([]byte, 0, bytes.MinRead), opts...)
}

// FileWith creates a new instance of [File] initialized with content. The
// created instance takes ownership of the content slice, and the caller must
// not use it after passing it to this function. FileWith is intended to
// prepare the File instance to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that, the content
// slice should have the desired capacity but a length of zero. It will panic
// with [ErrOutOfBounds] if the [WithFileOffset] option sets the offset to a
// negative number or beyond sata slice length.
func FileWith(content []byte, opts ...func(*File)) *File {
	fil := &File{buf: content}
	for _, opt := range opts {
		opt(fil)
	}
	if fil.off < 0 || fil.off > len(fil.buf) {
		panic(ErrOutOfBounds)
	}
	return fil
}

// Stat returns information about the in-memory file where the name is empty
// (unless, [WithFileName] option was used), the size is the length of the
// underlying buffer, mode is always 0444, modification time is always zero
// value time and [fs.FileInfo.Sys] always returns nil.
func (fil *File) Stat() (fs.FileInfo, error) {
	fi := FileInfo{
		name: fil.name,
		size: int64(len(fil.buf)),
	}
	return fi, nil
}

// Release releases ownership of the underlying buffer, the caller should not
// use this instance after this call.
func (fil *File) Release() []byte {
	buf := fil.buf
	fil.off = 0
	fil.buf = nil
	return buf
}

// Write writes the contents of p to the underlying buffer at the current
// offset, growing the buffer as needed. The return value n is the length of p;
// err is always nil.
func (fil *File) Write(p []byte) (n int, err error) {
	return fil.write(p), nil
}

// WriteByte writes a byte c to the underlying buffer at the current offset.
func (fil *File) WriteByte(c byte) error {
	fil.write([]byte{c})
	return nil
}

// WriteAt writes len(p) bytes to the underlying buffer starting at the current
// offset. It returns the number of bytes written; err is always nil. It does
// not change the offset.
func (fil *File) WriteAt(p []byte, off int64) (n int, err error) {
	if fil.flag&os.O_APPEND != 0 {
		return 0, errWriteAtInAppendMode
	}

	prev := fil.off
	c := cap(fil.buf)
	pl := len(p)

	// Handle writing beyond capacity.
	if int(off)+pl > c {
		fil.off = c // So tryGrowByReslice returns false.
		fil.grow(int(off) + pl - len(fil.buf))
		fil.buf = fil.buf[:int(off)+pl]
	}

	fil.off = int(off)
	n = fil.write(p)
	fil.off = prev
	return n, nil
}

// WriteTo writes data to w starting at the current offset until there are no
// more bytes to write or when an error occurs. The int64 return value is the
// number of bytes written. When an error occurred during the operation, it is
// also returned.
func (fil *File) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(fil.buf[fil.off:])
	fil.off += n
	return int64(n), err
}

// WriteString writes string s to the buffer at the current offset.
func (fil *File) WriteString(s string) (int, error) {
	return fil.Write([]byte(s)) // nolint: gocritic
}

// write writes p at the current offset.
func (fil *File) write(p []byte) int {
	if fil.flag&os.O_APPEND != 0 {
		fil.off = len(fil.buf)
	}
	l := len(fil.buf)
	fil.grow(len(p))
	n := copy(fil.buf[fil.off:], p)
	fil.off += n
	if fil.off > l {
		l = fil.off
	}
	fil.buf = fil.buf[:l]
	return n
}

// Read reads the next len(p) bytes from the buffer at the current offset or
// until the buffer is drained. The return value is the number of bytes read.
// If the buffer has no data to return, err is [io.EOF] (unless len(p) is zero);
// otherwise it is nil.
func (fil *File) Read(p []byte) (int, error) {
	// Nothing more to read.
	if len(p) > 0 && fil.off >= len(fil.buf) {
		return 0, io.EOF
	}
	n := copy(p, fil.buf[fil.off:])
	fil.off += n
	return n, nil
}

// ReadByte reads and returns the next byte from the buffer at the current
// offset or returns an error. If ReadByte returns an error, no input byte was
// consumed, and the returned byte value is undefined.
func (fil *File) ReadByte() (byte, error) {
	// Nothing more to read.
	if fil.off >= len(fil.buf) {
		return 0, io.EOF
	}
	v := fil.buf[fil.off]
	fil.off++
	return v, nil
}

// ReadAt reads len(p) bytes from the buffer at the current offset. It returns
// the number of bytes read and the error, if any. ReadAt always returns a
// non-nil error when n < len(p). It does not change the offset.
func (fil *File) ReadAt(p []byte, off int64) (int, error) {
	prev := fil.off
	defer func() { fil.off = prev }()
	fil.off = int(off)
	n, err := fil.Read(p)
	if err != nil {
		return n, err
	}
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// ReadFrom reads data from r until EOF and appends it to the buffer at the
// current offset, growing the buffer as needed. The return value is the number
// of bytes read. Any error except [io.EOF] encountered during the read is also
// returned. If the buffer becomes too large, ReadFrom will panic with
// [bytes.ErrTooLarge].
func (fil *File) ReadFrom(r io.Reader) (int64, error) {
	var err error
	var n, total int

	if fil.flag&os.O_APPEND != 0 {
		fil.off = len(fil.buf)
	}

	for {

		// Length before growing the buffer.
		l := len(fil.buf)

		// Make sure we can fit [bytes.MinRead] between the current offset and
		// the new buffer length.
		fil.grow(bytes.MinRead)

		// We will use bytes between l and cap(fil.buf) as a temporary scratch
		// space for reading from r and then slide read bytes to place. We have
		// to do it this way because [io.Read] documentation says that: "Even
		// if Read returns n < len(p), it may use all of p as scratch space
		// during the call." so we can't pass our buffer to Read because it
		// might change parts of it not involved in the read operation.
		tmp := fil.buf[l:cap(fil.buf)]
		n, err = r.Read(tmp)

		if l != fil.off {
			// Move bytes from temporary area to correct place.
			copy(fil.buf[fil.off:], tmp[:n])
			if n < len(tmp) {
				// Clean up any garbage the reader might put in there. We want
				// to keep all bytes between len and cap as zeros.
				zeroOutSlice(tmp[n:])
			}
		}

		fil.off += n
		total += n

		if fil.off > l {
			l = fil.off
		}

		// Set proper buffer length.
		fil.buf = fil.buf[:l]

		if err != nil {
			break
		}
	}

	// The [io.EOF] is not an error.
	if err == io.EOF {
		err = nil
	}

	return int64(total), err
}

// String returns string representation of the buffer starting at the current
// offset. Calling this method is considered as reading the buffer and advances
// offset to the end of the buffer.
func (fil *File) String() string {
	s := string(fil.buf[fil.off:])
	fil.off = len(fil.buf)
	return s
}

// Seek sets the offset for the next Read or Write on the buffer to the offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error (only if calculated offset < 0).
// Returns non-nil error of type [fs.PathError] where [fs.PathError.Path] field
// is an empty string unless [WithFileName] option was used during instance
// creation.
func (fil *File) Seek(offset int64, whence int) (int64, error) {
	var off int
	switch whence {
	case io.SeekStart:
		off = int(offset)
	case io.SeekCurrent:
		off = fil.off + int(offset)
	case io.SeekEnd:
		off = len(fil.buf) + int(offset)
	}

	if off < 0 {
		return 0, &fs.PathError{Op: "seek", Path: fil.name, Err: syscall.EINVAL}
	}
	fil.off = off

	return int64(fil.off), nil
}

// SeekStart is a convenience method setting the buffer's offset to zero and
// returning the value it had before the method was called.
func (fil *File) SeekStart() int64 {
	prev := fil.off
	fil.off = 0
	return int64(prev)
}

// SeekEnd is a convenience method setting the buffer's offset to the buffer
// length and returning the value it had before the method was called.
func (fil *File) SeekEnd() int64 {
	prev := fil.off
	fil.off = len(fil.buf)
	return int64(prev)
}

// Truncate changes the size of the buffer discarding bytes at the offsets
// greater than size. It does not change the offset unless a [WithFileAppend]
// option was used, then it sets the offset to the end of the buffer. Returns
// an error only when size is negative. The error is of type [fs.PathError]
// where [fs.PathError.Path] field is an empty string unless [WithFileName]
// option was used during instance creation.
func (fil *File) Truncate(size int64) error {
	if size < 0 {
		return &os.PathError{
			Op:   "truncate",
			Path: fil.name,
			Err:  syscall.EINVAL,
		}
	}

	prev := fil.off
	l := len(fil.buf)
	c := cap(fil.buf)

	switch {
	case int(size) == l:
		// Nothing to do.

	case int(size) == c:
		// Reslice.
		fil.buf = fil.buf[:size]

	case int(size) > l && int(size) < c:
		// Truncate between len and cap.
		fil.buf = fil.buf[:size]

	case int(size) > c:
		// Truncate beyond the cap.
		fil.off = c // So tryGrowByReslice returns false.
		fil.grow(int(size) - l)
		fil.buf = fil.buf[:int(size)]

	default:
		// Reduce the size of the buffer.
		zeroOutSlice(fil.buf[size:])
		fil.buf = fil.buf[:size]
	}

	fil.off = prev

	return nil
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation. If n is negative, Grow will panic. If the
// buffer can't grow, it will panic with [bytes.ErrTooLarge].
func (fil *File) Grow(n int) {
	if n < 0 {
		panic("fskit.File.Grow: negative count")
	}

	l := len(fil.buf)
	if l+n <= cap(fil.buf) {
		return
	}

	// Allocate bigger buffer.
	tmp := makeSlice(l + n)
	copy(tmp, fil.buf)
	fil.buf = tmp
	fil.buf = fil.buf[:l]
}

// grow grows the buffer's capacity to guarantee space for n more bytes. In
// other words, it makes sure there are n bytes between the current offset and
// the buffer capacity. It's worth noting that after calling this method the
// len(b.buf) changes. If the buffer can't grow, it will panic with
// [bytes.ErrTooLarge].
func (fil *File) grow(n int) {
	// Try to grow by a reslice.
	if ok := fil.tryGrowByReslice(n); ok {
		return
	}
	if fil.buf == nil && n <= smallBufferSize {
		fil.buf = make([]byte, n, smallBufferSize)
		return
	}
	// Allocate bigger buffer.
	tmp := makeSlice(cap(fil.buf)*2 + n) // cap(b.buf) may be zero.
	copy(tmp, fil.buf)
	fil.buf = tmp
}

// tryGrowByReslice is an inlineable version of [File.grow] for the fast-case
// where the internal buffer only needs to be resliced. It returns whether it
// succeeded.
func (fil *File) tryGrowByReslice(n int) bool {
	// No need to do anything if there is enough space between the current
	// offset and the length of the buffer.
	if n <= len(fil.buf)-fil.off {
		return true
	}

	if n <= cap(fil.buf)-fil.off {
		fil.buf = fil.buf[:fil.off+n]
		return true
	}
	return false
}

// Offset returns the current offset.
func (fil *File) Offset() int {
	return fil.off
}

// Len returns the buffer length.
func (fil *File) Len() int {
	return len(fil.buf)
}

// Cap returns the buffer capacity, that is, the total space allocated for the
// buffer's data.
func (fil *File) Cap() int {
	return cap(fil.buf)
}

// Close sets offset to zero and zero put the buffer. It always returns
// nil error.
func (fil *File) Close() error {
	if fil == nil {
		return nil
	}
	fil.off = 0
	zeroOutSlice(fil.buf[0:len(fil.buf)])
	fil.buf = fil.buf[:0]
	return nil
}

// zeroOutSlice zeroes out the byte slice.
func zeroOutSlice(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with [bytes.ErrTooLarge].
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(bytes.ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
