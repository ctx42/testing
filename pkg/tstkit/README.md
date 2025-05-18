<!-- TOC -->
* [The `tstkit` Package](#the-tstkit-package)
  * [The `Buffer` Type](#the-buffer-type)
    * [WetBuffer](#wetbuffer)
    * [DryBuffer](#drybuffer)
  * [Clocks](#clocks)
<!-- TOC -->

# The `tstkit` Package

The `tstkit` package provides a set of reusable test helpers to streamline Go
testing. Instead of rewriting common test utilities, this package offers a
curated collection. Its goal is to balance simplicity and functionality, 
focusing on practical tools generic.

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

## Clocks

The `tstkit` implements four functions with the same signature as `time.Time`
which can be used to inject deterministic clocks. 

- `ClockStartingAt` - returns current time with given offset.  
- `ClockFixed` - always returns the same time.
- `ClockDeterministic` - returns time advanced by given duration no mather how fast you call it. 
- `TikTak` - like `ClockDeterministic` with duration set to 1 second.