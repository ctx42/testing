<!-- TOC -->
* [The `kit` Package](#the-kit-package)
  * [Top-level helpers](#top-level-helpers)
  * [Sub-packages](#sub-packages)
<!-- TOC -->

# The `kit` Package

The `kit` package is the curated home for focused testing utilities in the
CTX42 testing module. Rather than duplicating common test helpers across
projects, kit collects practical, well-documented tools that integrate
naturally with `tester.T` and the assertion packages.

Its goal is to provide high-quality building blocks that are simple to use
yet powerful enough for real test scenarios.

## Top-level helpers

In addition to the sub-packages, `kit` provides a small number of
standalone helpers at the package level:

- `SHA1Reader` / `SHA1File` — convenient wrappers for computing SHA-1
  hashes (they panic on error, which is the expected behavior in tests).
- [AddGlobalCleanup] / [RunGlobalCleanups] — a global post-test cleanup
  mechanism intended for use from `TestMain`. See the godoc for important
  warnings about global mutable state and recommended usage patterns.

## Sub-packages

- [iokit](iokit/README.md) — I/O and buffer-related helpers, including
  thread-safe buffers with automatic cleanup checks and error-injecting
  readers/writers.
- [timekit](timekit/README.md) — Controllable and deterministic clocks
  (fixed, starting-at, and tick-based) for testing time-dependent code
  without relying on the real system clock.
- [reflectkit](reflectkit/reflectkit.go) — Lightweight reflection
  utilities, primarily for safe struct field inspection during tests.
- [randkit](randkit/random.go) — Cryptographically random test helpers
  for generating strings, file names, integers, and passwords via
  `crypto/rand`.

See the individual sub-package READMEs and godoc for detailed usage,
examples, and cross-references.

