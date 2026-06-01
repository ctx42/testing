// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// action represents [Spy] method call.
type action int

// Spy types of method calls.
const (
	mockedCall       action = iota // Call to one of the methods [Spy] mocks.
	mockedFailedCall               // Call to [Spy.Failed] method.
	expectCall                     // Call to one of the Spy.Expect* methods.
	assertCall                     // Call to one of the Spy.Assert* methods.
	closeCall                      // Call to [Spy.Close] method.
)

// Strategy controls how [Spy.ExpectLog] matches log messages produced by
// the helper under test.
//
// See the package README section on examining log messages for usage with
// ExpectLog and related methods.
type Strategy string

// Log matching strategies for [Spy.ExpectLog].
const (
	// Equal requires the produced log message to be identical to the
	// expected string (after formatting).
	Equal Strategy = "equal"

	// Contains requires the produced log message to contain the expected
	// substring.
	Contains Strategy = "contains"

	// NotContains requires the produced log message to NOT contain the
	// expected substring.
	NotContains Strategy = "not-contains"

	// Regexp requires the produced log message to match the expected
	// regular expression.
	Regexp Strategy = "regexp"
)

// find implements log matching strategies.
type find struct {
	strategy Strategy // Match strategy.
	want     string   // Expected string.
}

// match returns true if "want" can be found in "have" using strategy.
func (ent find) match(have string) bool {
	switch ent.strategy {
	case Equal:
		return ent.want == have
	case Regexp:
		return regexp.MustCompile(ent.want).MatchString(have)
	case Contains:
		return strings.Contains(have, ent.want)
	case NotContains:
		return !strings.Contains(have, ent.want)
	default:
		return ent.want == have
	}
}

// Spy usage errors.
const (
	errInvalidUsage             = "invalid Spy usage"
	errMockOnNotClosed          = "mocked call on not closed Spy is not allowed"
	errExpectOnClosed           = "expectation on closed Spy is not allowed"
	errExpectOnFinished         = "expectation on finished Spy is not allowed"
	errActionOnFinished         = "action on finished Spy is not allowed"
	errDoubleClose              = "calling Close twice is not allowed"
	errDoubleFinish             = "calling Finish twice is not allowed"
	errCloseOnFinished          = "close on finished Spy is not allowed"
	errAssertOnNotClosed        = "assertion on not closed Spy is not allowed"
	errAssertOnNotFinished      = "assertion on not finished Spy is not allowed"
	errIgnoreLogsAfterExpectLog = "calling IgnoreLogs after ExpectLog* is not allowed"
	errExpectLogAfterIgnoreLogs = "calling ExpectLog* after IgnoreLogs is not allowed"
)

// FailNowMsg is the exact panic message produced when [Spy.FailNow] (or
// Fatal/Fatalf) is called on the spy. Tests that intentionally exercise
// FailNow paths can match against this constant.
const FailNowMsg = "FailNow was called directly"

// Spy is a test double implementing [T] (and therefore a subset of
// [testing.TB]) that records calls and lets you set precise expectations.
//
// It is the standard way to test custom assertion helpers and other code
// that takes a test manager as an argument. Typical usage (the pattern
// used throughout this module, including in pkg/goldy tests):
//
// Happy path:
//
//	tspy := New(t)
//	tspy.Close()          // no more expectations
//	have := MyHelper(tspy, ...)
//	affirm.False(t, tspy.Failed())
//
// Error path:
//
//	tspy := New(t)
//	tspy.ExpectError()
//	tspy.ExpectLogEqual("expected message with %s", "detail")
//	tspy.Close()
//	have := MyHelper(tspy, badInput)
//	affirm.Nil(t, have)
//
// After Close, call [Spy.AssertExpectations] explicitly if you did not
// let the automatic cleanup handle it. See the package [README] for the
// complete guide to Expect*, log strategies, cleanups, TempDir, etc.
type Spy struct {
	// Set to true if requirement for number of calls to the mocked Helper
	// method was explicitly set.
	helperCntSet bool

	// Expected number of calls to the mocked Helper method. By default, set to
	// -1 which means the mocked method must be called at least once.
	wantHelperCnt int

	// Actual number of calls to the mocked Helper method made by the HUT.
	haveHelperCnt int

	// Expected number of calls to the mocked TempDir method.
	wantTempDirCnt int

	// Actual directories returned from the mocked TempDir method.
	haveTempDirs []string

	// Expected number of calls to mocked Name method.
	wantNamesCnt int

	// Actual number of calls to the mocked Name method made by the HUT.
	haveNamesCnt int

	// Environment variables expected to be set by the HUT.
	wantEnv map[string]string

	// Environment variables actually set by the HUT.
	haveEnv map[string]string

	// When true, no more expectations can be added to the Spy.
	closed bool

	// True when the test has finished. When the test is finished, calls all the
	// Spy methods (except Failed) will panic.
	finished bool

	// Expected outcome of running HUT. When set to true, it does not matter
	// which of the Error* or Fatal* methods were called by the HUT.
	wantFailed bool

	// True when we expect HUT to call one of the Error* methods at least once.
	wantError bool

	// True when HUT called one of the Error* methods at least once.
	haveError bool

	// True when we expect HUT to call one of the Fatal* methods at least once.
	wantFatal bool

	// True when HUT called one of the Fatal* methods at least once.
	haveFatal bool

	// Expected test skip status when running HUT.
	wantSkipped bool

	// Actual test skip status when running HUT.
	haveSkipped bool

	// When true, the Spy panicked due to misuse.
	panicked bool

	// Expected number of cleanup functions set by HUT.
	wantCleanupsCnt int

	// Actual number of cleanup functions set by HUT.
	hadCleanupsCnt int

	// Cleanup functions which will be run before running assertions.
	haveCleanups []func()

	// Messages sent to the actual test runner (the one received in New
	// function).
	savedMgs []string

	// Log messages expected to be printed by the HUT.
	wantLogMgs []find

	// Actual messages sent to the mocked test runner Log and Logf methods.
	haveLogMgs []string

	// When set to true, it will not trigger an assertion error if haveLogMgs
	// is not empty and wantLogMgs is empty.
	ignoreLog bool

	// Test runner which we use for reporting errors when Spy expectations
	// do not match the actual HUT behavior. It is also used to for TempDir
	// and Setenv methods.
	tt *testing.T

	// True when the Finish method is running.
	runningFinish bool

	// Context returned by Context method.
	ctx context.Context

	// When context was retrieved using Context method this is set to a
	// function canceling it. When set, it will run right before functions
	// are registered via the Cleanup method.
	cancelCtx context.CancelFunc

	// Guards the above fields.
	mx sync.Mutex
}

// New creates a new [Spy] bound to the real test runner tt.
//
// tt is used for:
//   - proxying TempDir, Setenv, Name, and Context
//   - reporting failures when expectations are not met
//
// New registers a [testing.TB.Cleanup] that automatically calls Finish and
// AssertExpectations when the test ends.
//
// The standard lifecycle (used by this module and recommended for all users):
//
//	tspy := New(t)           // create
//	tspy.ExpectError()       // set expectations (0 or more)
//	tspy.ExpectLogEqual(...)
//	tspy.Close()             // signal no more expectations
//
//	// exercise the helper under test with tspy
//	have := MyHelper(tspy, input)
//
//	// either let the cleanup assert, or call explicitly:
//	tspy.AssertExpectations()
//
// The optional expectHelpers argument is a convenience for the common case
// of wanting an exact number of Helper calls; see [Spy.ExpectHelpers].
func New(tt *testing.T, expectHelpers ...int) *Spy {
	tt.Helper()
	spy := &Spy{tt: tt, wantHelperCnt: -1}
	cu := func() {
		tt.Helper()
		spy.mx.Lock()
		if !spy.finished {
			spy.mx.Unlock()
			spy.Finish()
			spy.mx.Lock()
		}
		defer spy.mx.Unlock()
		spy.assertExpectations()
		spy.tt = nil
	}
	tt.Cleanup(cu)
	if len(expectHelpers) > 0 {
		spy.ExpectHelpers(expectHelpers[0])
	}
	return spy
}

// ExpectCleanups declares how many times the helper under test must call
// Cleanup. Use 0 to assert that no cleanups were registered.
func (spy *Spy) ExpectCleanups(cnt int) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	spy.wantCleanupsCnt = cnt
	return spy
}

// Cleanup registers a function to be called when the test and all its subtests
// complete. The registered function is always called at the end of the test.
func (spy *Spy) Cleanup(f func()) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedCall)
	spy.hadCleanupsCnt++
	spy.haveCleanups = append(spy.haveCleanups, f)
}

// ExpectError declares that the helper under test must call Error or Errorf
// at least once. Mutually exclusive with [Spy.ExpectFail].
func (spy *Spy) ExpectError() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantFailed {
		spy.panicked = true
		panic("cannot use ExpectError and ExpectFail at the same time")
	}
	spy.wantError = true
	return spy
}

func (spy *Spy) Error(args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.log(args...)
	spy.haveError = true
}

func (spy *Spy) Errorf(format string, args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.logf(format, args...)
	spy.haveError = true
}

// ExpectFatal declares that the helper under test must call Fatal or Fatalf
// at least once. Mutually exclusive with [Spy.ExpectFail].
func (spy *Spy) ExpectFatal() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantFailed {
		spy.panicked = true
		panic("cannot use ExpectFatal and ExpectFail at the same time")
	}
	spy.wantFatal = true
	return spy
}

func (spy *Spy) Fatal(args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.log(args...)
	spy.failNow()
}

func (spy *Spy) Fatalf(format string, args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.logf(format, args...)
	spy.failNow()
}

func (spy *Spy) FailNow() {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.failNow()
}

func (spy *Spy) failNow() {
	spy.tt.Helper()
	spy.checkState(mockedCall)
	spy.haveFatal = true
	panic(FailNowMsg)
}

// Failed reports whether the helper under test called any Error*, Fatal*,
// or FailNow method. Note that returning false does not mean all [Spy]
// expectations were satisfied — the spy may still have unmet call-count or
// log expectations.
func (spy *Spy) Failed() bool {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedFailedCall)
	return spy.haveFatal || spy.haveError
}

// ExpectHelpers sets the exact number of times the helper under test must
// call Helper. The value -1 (the default) means "at least once".
//
// Must be called at most once and before Close. Panics on invalid cnt (< -1)
// or multiple calls.
func (spy *Spy) ExpectHelpers(cnt int) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	if spy.helperCntSet {
		spy.panicked = true
		panic("ExpectHelpers may be called only once")
	}
	if cnt < -1 {
		spy.panicked = true
		panic("ExpectHelpers cnt must be greater or equal to minus one")
	}
	spy.checkState(expectCall)
	spy.wantHelperCnt = cnt
	spy.helperCntSet = true
	return spy
}

func (spy *Spy) Helper() {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedCall)
	spy.haveHelperCnt++
}

// ExpectSetenv declares that the helper must call Setenv with exactly this
// key and value.
func (spy *Spy) ExpectSetenv(key, value string) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantEnv == nil {
		spy.wantEnv = make(map[string]string)
	}
	spy.wantEnv[key] = value
	return spy
}

func (spy *Spy) Setenv(key, value string) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedCall)
	if spy.haveEnv == nil {
		spy.haveEnv = make(map[string]string)
	}
	spy.haveEnv[key] = value
	spy.tt.Setenv(key, value)
}

// ExpectSkipped declares that the helper under test must call Skip.
func (spy *Spy) ExpectSkipped() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	spy.wantSkipped = true
	return spy
}

func (spy *Spy) Skip(args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.log(args...)
	spy.haveSkipped = true
}

// ExpectTempDir declares how many times TempDir must be called. -1 means
// "any number of times" (including zero).
func (spy *Spy) ExpectTempDir(cnt int) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	spy.wantTempDirCnt = cnt
	return spy
}

// GetTempDir returns the path of the Nth TempDir call (0-based). Requires
// that ExpectTempDir was called first; otherwise it reports an error.
func (spy *Spy) GetTempDir(idx int) string {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	if spy.wantTempDirCnt == 0 {
		spy.tError("ExpectTempDir method must be called before GetTempDir")
		return ""
	}
	if idx >= len(spy.haveTempDirs) {
		format := "temp directory with index %d does not exist"
		spy.tErrorf(format, idx)
		return ""
	}
	return spy.haveTempDirs[idx]
}

func (spy *Spy) TempDir() string {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedCall)
	pth := spy.tt.TempDir()
	spy.haveTempDirs = append(spy.haveTempDirs, pth)
	spy.haveCleanups = append(spy.haveCleanups, func() { _ = os.RemoveAll(pth) })
	return pth
}

func (spy *Spy) Context() context.Context {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedCall)
	parent := spy.tt.Context()
	if parent == nil {
		parent = context.Background()
	}
	if spy.ctx == nil {
		spy.ctx, spy.cancelCtx = context.WithCancel(parent)
	}
	return spy.ctx
}

// IgnoreLogs tells the Spy to stop requiring that every log message be
// accounted for by an ExpectLog* call. Must not be called after any
// ExpectLog* method. See the "Ignore Log Messages" section in the README.
func (spy *Spy) IgnoreLogs() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if len(spy.wantLogMgs) > 0 {
		spy.panicked = true
		panic(errIgnoreLogsAfterExpectLog)
	}
	spy.ignoreLog = true
	return spy
}

// ExpectLog declares that the helper under test must produce a log message
// (via Log or Logf) that matches the given strategy and formatted string.
// Panics if [Spy.IgnoreLogs] was already called.
func (spy *Spy) ExpectLog(matcher Strategy, msg string, args ...any) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.ignoreLog {
		spy.panicked = true
		panic(errExpectLogAfterIgnoreLogs)
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	if msg == "" {
		return spy
	}
	ent := find{
		strategy: matcher,
		want:     msg,
	}
	spy.wantLogMgs = append(spy.wantLogMgs, ent)
	return spy
}

// ExpectLogEqual is a convenience for ExpectLog([Equal], ...). The message
// must match exactly after formatting.
func (spy *Spy) ExpectLogEqual(format string, args ...any) *Spy {
	return spy.ExpectLog(Equal, format, args...)
}

// ExpectLogContain is a convenience for ExpectLog([Contains], ...).
func (spy *Spy) ExpectLogContain(format string, args ...any) *Spy {
	return spy.ExpectLog(Contains, format, args...)
}

// ExpectLogNotContain is a convenience for ExpectLog([NotContains], ...).
func (spy *Spy) ExpectLogNotContain(format string, args ...any) *Spy {
	return spy.ExpectLog(NotContains, format, args...)
}

func (spy *Spy) Log(args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.log(args...)
}

func (spy *Spy) log(args ...any) {
	spy.tt.Helper()
	spy.checkState(mockedCall)
	msg := fmt.Sprintln(args...)
	if msg != "" {
		msg = msg[:len(msg)-1]
	}
	spy.haveLogMgs = append(spy.haveLogMgs, msg)
}

func (spy *Spy) Logf(format string, args ...any) {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.logf(format, args...)
}

func (spy *Spy) logf(format string, args ...any) {
	spy.tt.Helper()
	spy.checkState(mockedCall)
	msg := fmt.Sprintf(format, args...)
	spy.haveLogMgs = append(spy.haveLogMgs, msg)
}

// ExamineLog returns the concatenation of all messages logged so far by
// the helper under test (via Log/Logf). Useful for custom assertions when
// the built-in ExpectLog* matchers are not sufficient.
func (spy *Spy) ExamineLog() string {
	spy.tt.Helper()
	spy.mx.Lock()
	defer spy.mx.Unlock()
	return strings.Join(spy.haveLogMgs, "\n")
}

// ExpectNames declares how many times Name must be called on the spy.
func (spy *Spy) ExpectNames(cnt int) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	spy.wantNamesCnt = cnt
	return spy
}

func (spy *Spy) Name() string {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedCall)
	spy.haveNamesCnt++
	return spy.tt.Name()
}

// ExpectFail declares that the helper under test must call at least one of
// the error or fatal methods (Error*, Fatal*, or FailNow). This is the
// broadest "I expect the helper to report a problem" expectation.
// Mutually exclusive with ExpectError and ExpectFatal.
func (spy *Spy) ExpectFail() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantFatal {
		spy.panicked = true
		panic("cannot use ExpectFail and ExpectFatal at the same time")
	}
	if spy.wantError {
		spy.panicked = true
		panic("cannot use ExpectFail and ExpectError at the same time")
	}
	spy.wantFailed = true
	return spy
}

// Close marks the end of the expectation setup phase. After Close, no more
// Expect* calls are allowed. This is the required last step before passing
// the Spy to the helper under test in the patterns shown in [New].
//
// Returns the Spy for chaining: tspy := New(t).ExpectError().Close()
func (spy *Spy) Close() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(closeCall)
	spy.closed = true
	return spy
}

// Finish runs all registered cleanups and marks the Spy as finished. After
// Finish, most methods will panic if called (see per-method docs). Finish
// is called automatically by the cleanup registered in [New]; you normally
// only call it explicitly in advanced scenarios.
func (spy *Spy) Finish() *Spy {
	spy.mx.Lock()
	spy.tt.Helper()
	if spy.runningFinish || spy.finished {
		spy.panicked = true
		spy.mx.Unlock()
		panic(errDoubleFinish)
	}
	spy.runningFinish = true
	spy.mx.Unlock()

	spy.runCleanups()

	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.runningFinish = false
	spy.finished = true
	return spy
}

// AssertExpectations verifies that every expectation set via the Expect*
// methods was satisfied by the helper under test. It returns true on
// success. Failures are reported via the underlying *testing.T.
//
// In normal usage you do not need to call this explicitly — the cleanup
// registered by [New] does it for you when the test ends.
func (spy *Spy) AssertExpectations() bool {
	spy.tt.Helper()
	spy.mx.Lock()
	if !spy.finished {
		spy.mx.Unlock()
		spy.Finish()
		spy.mx.Lock()
	}
	defer spy.mx.Unlock()
	return spy.assertExpectations()
}

// assertExpectations has the same logic as [Spy.AssertExpectations] but
// assumes the caller has acquired the lock.
func (spy *Spy) assertExpectations() bool {
	spy.tt.Helper()
	spy.checkState(assertCall)
	if spy.panicked {
		spy.tErrorf(errInvalidUsage)
		return false
	}

	ok := spy.assertSkipped()
	if spy.wantFailed {
		if ret := spy.assertFailed(); ok {
			ok = ret
		}
	} else {
		if ret := spy.assertError(); ok {
			ok = ret
		}
		if ret := spy.assertFatal(); ok {
			ok = ret
		}
	}

	ret := spy.assertHelperCalls(spy.wantHelperCnt, spy.haveHelperCnt)
	if ok {
		ok = ret
	}
	ret = spy.checkCallCnt("Cleanup", spy.wantCleanupsCnt, spy.hadCleanupsCnt)
	if ok {
		ok = ret
	}
	ret = spy.checkCallMaybeCnt(
		"TempDir",
		spy.wantTempDirCnt,
		len(spy.haveTempDirs),
	)
	if ok {
		ok = ret
	}
	ret = spy.assertLogMgs(spy.wantLogMgs, spy.haveLogMgs)
	if ok {
		ok = ret
	}
	ret = spy.checkCallCnt("Name", spy.wantNamesCnt, spy.haveNamesCnt)
	if ok {
		ok = ret
	}
	ret = spy.assertEnv(spy.wantEnv, spy.haveEnv)
	if ok {
		ok = ret
	}
	return ok
}

// assertSkipped asserts HUT reacted according to expectation set by
// [Spy.ExpectSkipped] method.
func (spy *Spy) assertSkipped() bool {
	spy.tt.Helper()
	if spy.wantSkipped == spy.haveSkipped {
		return true
	}
	msg := "expected HUT to mark test as skipped:\n" +
		"\twant: %v\n" +
		"\thave: %v"
	spy.tErrorf(msg, spy.wantSkipped, spy.haveSkipped)
	return false
}

// assertFailed asserts HUT reacted according to the expectation set by
// [Spy.ExpectFail] method. If the [Spy.ExpectFail] method was not called, this
// method will always return true.
func (spy *Spy) assertFailed() bool {
	spy.tt.Helper()
	if spy.wantFailed {
		if spy.haveError || spy.haveFatal {
			return true
		}
		spy.tError("expected HUT to call the t.Error* or t.Fatal* methods")
		return false
	}
	return true
}

// assertError asserts HUT reacted according to the expectation set by
// [Spy.ExpectError] method.
func (spy *Spy) assertError() bool {
	spy.tt.Helper()
	if spy.wantError == spy.haveError {
		return true
	}
	msg := "expected HUT not to call any of the t.Error* methods"
	if spy.wantError {
		msg = "expected HUT to call any of the t.Error* methods"
	}
	spy.tError(msg)
	return false
}

// assertFatal asserts HUT reacted according to expectation set by
// [Spy.ExpectFatal] method.
func (spy *Spy) assertFatal() bool {
	spy.tt.Helper()
	if spy.wantFatal == spy.haveFatal {
		return true
	}
	msg := "expected HUT not to call any of the t.Fatal* methods"
	if spy.wantFatal {
		msg = "expected HUT to call any of the t.Fatal* methods"
	}
	spy.tError(msg)
	return false
}

// checkCallCnt checks method with the name was called expected number of times.
func (spy *Spy) checkCallCnt(name string, want, have int) bool {
	spy.tt.Helper()
	ok := true
	if want != have {
		format := "expected %s to be called N times:\n" +
			"\twant: %d\n" +
			"\thave: %d"
		spy.tErrorf(format, name, want, have)
		ok = false
	}
	return ok
}

// checkCallMaybeCnt checks method with the name was called expected number of
// times. If want is equal to -1, the return is always true, otherwise want and
// have is checked for equality.
func (spy *Spy) checkCallMaybeCnt(name string, want, have int) bool {
	spy.tt.Helper()
	if want == -1 {
		return true
	}
	ok := true
	if want != have {
		format := "expected %s to be called N times:\n" +
			"\twant: %d\n" +
			"\thave: %d"
		spy.tErrorf(format, name, want, have)
		ok = false
	}
	return ok
}

// assertHelperCalls asserts HUT reacted according to expectation set by
// [ExpectHelpers] method.
func (spy *Spy) assertHelperCalls(want, have int) bool {
	spy.tt.Helper()
	ok := true
	if (want == -1 && have == 0) || (want >= 0 && want != have) {
		wantN := ">= 1"
		if want >= 0 {
			wantN = strconv.Itoa(want)
		}
		format := "expected Helper to be called N times:\n" +
			"\twant: %s\n" +
			"\thave: %d"
		spy.tErrorf(format, wantN, have)
		ok = false
	}
	return ok
}

// assertLogMgs asserts HUT reacted according to expectation set by
// [Spy.ExpectLog], [Spy.ExpectLogContain], [Spy.ExpectLogEqual] methods.
func (spy *Spy) assertLogMgs(wants []find, haves []string) bool {
	spy.tt.Helper()
	haveStr := strings.Join(haves, "\n")
	if haveStr != "" && len(wants) == 0 {
		if spy.ignoreLog {
			return true
		}
		format := "expected HUT to log no messages but got:\n" +
			"\thave: %q"
		spy.tErrorf(format, haveStr)
		return false
	}
	ok := true
	for idx, want := range wants {
		if !want.match(haveStr) {
			format := "expected HUT to log message %d:\n" +
				"\tmatcher: %s\n" +
				"\t   want: %q\n" +
				"\t   have: %q"
			spy.tErrorf(format, idx, want.strategy, want.want, haveStr)
			ok = false
		}
	}
	return ok
}

// assertEnv asserts HUT reacted according to the expectation set by
// [Spy.ExpectSetenv] method.
func (spy *Spy) assertEnv(wants, haves map[string]string) bool {
	spy.tt.Helper()

	ok := true
	wantKeys := make([]string, 0, len(wants))
	for wantKey := range wants {
		wantKeys = append(wantKeys, wantKey)
	}
	sort.Strings(wantKeys)

	for _, wantKey := range wantKeys {
		wantVal := wants[wantKey]
		if haveVal, exists := haves[wantKey]; exists {
			if wantVal != haveVal {
				format := "expected HUT to set environment variable:\n" +
					"\t  want key: %q\n" +
					"\twant value: %q\n" +
					"\thave value: %q"
				spy.tErrorf(format, wantKey, wantVal, haveVal)
				ok = false
			}
		} else {
			format := "expected HUT to set environment variable:\n" +
				"\t  want key: %q\n" +
				"\twant value: %q"
			spy.tErrorf(format, wantKey, wantVal)
			ok = false
		}
	}

	if len(wants) < len(haves) {
		haveKeys := make([]string, 0, len(wants))
		for wantKey := range haves {
			haveKeys = append(haveKeys, wantKey)
		}
		sort.Strings(haveKeys)

		for _, haveKey := range haveKeys {
			haveVal := haves[haveKey]
			if _, exists := wants[haveKey]; !exists {
				format := "expected HUT not to set environment variable:\n" +
					"\t  have key: %q\n" +
					"\thave value: %q"
				spy.tErrorf(format, haveKey, haveVal)
				ok = false
			}
		}
	}
	return ok
}

// runCleanups runs registered cleanups.
func (spy *Spy) runCleanups() int {
	cnt := 0

	if spy.cancelCtx != nil {
		spy.cancelCtx()
	}

	for {
		spy.mx.Lock()
		var cleanup func()
		if len(spy.haveCleanups) > 0 {
			last := len(spy.haveCleanups) - 1
			cleanup = spy.haveCleanups[last]
			spy.haveCleanups = spy.haveCleanups[:last]
		}
		if cleanup == nil {
			spy.mx.Unlock()
			break
		}
		spy.mx.Unlock()
		cleanup()
		cnt++
	}
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.ctx = nil
	spy.cancelCtx = nil
	spy.haveCleanups = spy.haveCleanups[:0]
	return cnt
}

// checkState check if requested action is valid for current Spy state.
//
// nolint:cyclop
func (spy *Spy) checkState(action action) {
	spy.tt.Helper()
	switch action {
	case expectCall:
		if spy.finished {
			spy.tError(errExpectOnFinished)
			spy.panicked = true
			panic(errExpectOnFinished)
		}

		if spy.closed {
			spy.tError(errExpectOnClosed)
			spy.panicked = true
			panic(errExpectOnClosed)
		}

	case closeCall:
		if spy.finished {
			spy.tError(errCloseOnFinished)
			spy.panicked = true
			panic(errCloseOnFinished)
		}

		if spy.closed {
			spy.tError(errDoubleClose)
			spy.panicked = true
			panic(errDoubleClose)
		}

	case mockedCall:
		if spy.finished {
			spy.tError(errActionOnFinished)
			spy.panicked = true
			panic(errActionOnFinished)
		}

		if !spy.closed {
			spy.tError(errMockOnNotClosed)
			spy.panicked = true
			panic(errMockOnNotClosed)
		}

	case assertCall:
		if !spy.finished {
			spy.tError(errAssertOnNotFinished)
			spy.panicked = true
			panic(errAssertOnNotFinished)
		}

		if !spy.closed {
			spy.tError(errAssertOnNotClosed)
			spy.panicked = true
			panic(errAssertOnNotClosed)
		}

	case mockedFailedCall:
		if !spy.closed {
			spy.tError(errMockOnNotClosed)
			spy.panicked = true
			panic(errMockOnNotClosed)
		}
	}
}

// tError saves messages send to T.Error.
func (spy *Spy) tError(args ...any) {
	spy.tt.Helper()
	msg := fmt.Sprint(args...)
	spy.savedMgs = append(spy.savedMgs, msg)
	spy.tt.Error(msg)
}

// tErrorf saves messages send to T.Errorf.
func (spy *Spy) tErrorf(format string, args ...any) {
	spy.tt.Helper()
	msg := fmt.Sprintf(format, args...)
	spy.savedMgs = append(spy.savedMgs, msg)
	spy.tt.Error(msg)
}
