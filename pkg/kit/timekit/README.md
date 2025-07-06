<!-- TOC -->
* [The `timekit` package](#the-timekit-package)
  * [Clocks](#clocks)
<!-- TOC -->

# The `timekit` package

The `timekit` package provides `time.Time` related helpers. 

## Clocks

Clock functions have the same signature as `time.Time` which can be used to
inject deterministic clocks.

- `ClockStartingAt` - returns current time with given offset.
- `ClockFixed` - always returns the same time.
- `ClockDeterministic` - returns time advanced by given duration no mather how fast you call it.
- `TikTak` - like `ClockDeterministic` with duration set to 1 second.
