<!-- TOC -->
* [The `iokit` package](#the-iokit-package)
  * [The `Buffer` Type](#the-buffer-type)
    * [WetBuffer](#wetbuffer)
    * [DryBuffer](#drybuffer)
  * [Error writers and readers](#error-writers-and-readers)
<!-- TOC -->

# The `iokit` package

The `iokit` package provides I/O and buffer related helpers. 

## The `Buffer` Type

The `Buffer` type, defined in `buffer.go`, is a thread-safe wrapper around
`bytes.Buffer`. It supports three kinds of behavior for test cleanup:

### WetBuffer

A `WetBuffer` uses `Buffer` and ensures it's written to and its contents are examined during the test.

```go
func TestAction(t *testing.T) {
    // --- Given ---
    buf := tstkit.WetBuffer(t, "wet-buffer")

    // --- When ---
    Action(buf) // Writes to buf.

    // --- Then ---
    // Fails if Action doesn't write to buf or if buf.String() is not called.
    assert.Equal(t, "expected output", buf.String())   
}
```

To skip the examination requirement:

```go
buf := tstkit.WetBuffer(t, "wet-buffer").SkipExamine()
buf.WriteString("data") // No failure for unexamined content
```

### DryBuffer

A DryBuffer ensures the buffer remains empty.

```go
func TestAction(t *testing.T) {
    // --- Given ---
    buf := tstkit.DryBuffer(t, "dry-buffer")

    // --- When ---
    DoSomething(buf) // Must not write to buf.

    // --- Then ---
    // Fails if DoSomething writes to buf.
}
```

## Error writers and readers

Package provides helpers for controlling when and how the most important I/O 
interfaces return an error.  

- `ErrReader` - control when and what error an `io.Reader` returns.
- `ErrReadCloser` - control when and what error an `io.ReadCloser` returns.
- `ErrReadSeeker` - control when and what error an `io.ReadSeeker` returns.
- `ErrReadSeekCloser` - control when and what error an `io.ReadSeekCloser` returns.
- `ErrWriter` - control when and what error an `io.Writer` returns.
- `ErrWriteCloser` - control when and what error an `io.WriteCloser` returns.
