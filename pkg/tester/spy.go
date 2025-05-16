// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

// Strategy is strategy of matching logs produced by Helper Under Test (HUT).
type Strategy string

// Log matching strategies.
const (
	// Equal is a strategy where messages logged by HUT are matched exactly.
	Equal Strategy = "equal"

	// Contains is a strategy where messages logged by HUT contain a string.
	Contains Strategy = "contains"

	// NotContains is a strategy where messages logged by HUT don't contain a
	// string.
	NotContains Strategy = "not-contains"

	// Regexp is a strategy where messages logged by the HUT match regular
	// expression.
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

// FailNowMsg represents a message the [Spy.FailNow] method uses in panic.
const FailNowMsg = "FailNow was called directly"

// Spy is a spy for [tester.T] interface.
//
// Creating test helpers is an integral part of comprehensive testing, but
// those helpers in turn also need to be tested to make sure assertions made by
// them are implemented correctly. The Spy is a tool that makes testing such
// helpers very easy.
//
// Pass instance of Spy to the Helper Under Test (HUT) and assert the expected
// behaviour using Spy.Expect* methods.
type Spy struct {
	// Set to true if requirement for number of calls to mocked Helper method
	// was explicitly set.
	helperCntSet bool

	// Expected number of calls to mocked Helper method. By default, set to -1
	// which means the mocked method must be called at least once.
	wantHelperCnt int

	// Actual number of calls to mocked Helper method made by the HUT.
	haveHelperCnt int

	// Expected number of calls to mocked TempDir method.
	wantTempDirCnt int

	// Actual directories returned from mocked TempDir method.
	haveTempDirs []string

	// Expected number of calls to mocked Name method.
	wantNamesCnt int

	// Actual number of calls to mocked Name method made by the HUT.
	haveNamesCnt int

	// Environment variables expected to be set by the HUT.
	wantEnv map[string]string

	// Environment variables actually set by the HUT.
	haveEnv map[string]string

	// When true no more expectations can be added to the Spy.
	closed bool

	// True when the test have finished. When test is finished calls all the
	// Spy methods (except Failed) will panic.
	finished bool

	// Expected outcome of running HUT. When set to true it does not matter
	// which of the Error* or Fatal* methods were called by the HUT.
	wantFailed bool

	// True when we expect HUT to call at least once one of the Error* methods.
	wantError bool

	// True when HUT called the at least once one of the Error* methods.
	haveError bool

	// True when we expect HUT to call at least once one of the Fatal* methods.
	wantFatal bool

	// True when HUT called at least once one of the Fatal* methods.
	haveFatal bool

	// Expected test skip status when running HUT.
	wantSkipped bool

	// Actual test skip status when running HUT.
	haveSkipped bool

	// When true the Spy panicked due to misuse.
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

	// When set to true it will not trigger assertion error if haveLogMgs is
	// not empty and wantLogMgs is empty.
	ignoreLog bool

	// Test runner which we use for reporting errors when Spy expectations
	// do not match the actual HUT behaviour. It is also used to for TempDir
	// and Setenv methods.
	tt *testing.T

	// True when Finish method is running.
	runningFinish bool

	// Context returned by Context method.
	ctx context.Context

	// When context was retrieved using Context method this is set to a
	// function canceling it. When set it will run right before functions
	// registered via Cleanup method.
	cancelCtx context.CancelFunc

	// Guards the above fields.
	mx sync.Mutex
}

// New returns new instance of [Spy] which implements [T] interface. The tt
// argument is used to proxy calls to [testing.T.TempDir], [testing.T.Setenv]
// and [testing.T.Context] as well as to report errors when the [Spy]
// expectations are not met by Helper Under Test (HUT). The constructor
// function adds a cleanup function to tt which calls [Spy.Finish] and
// [Spy.AssertExpectations] methods to determine if the tt should be failed.
//
// The call to New should be followed by zero or more calls to Spy.Expect*
// methods and finished with call to [Spy.Close] method:
//
//	tspy := New(t)
//	tspy.ExpectError()
//	// tspy.ExpectXXX()
//	tspy.Close()
//
// To assert expectations manually call [Spy.AssertExpectations] or it will be
// called automatically when test (tt) finishes.
//
// Full example:
//
//	t.Run("closes file at the end of test", func(t *testing.T) {
//		// --- Given ---
//		tspy := tester.New(t).ExpectCleanups(1).Close()
//
//		// --- When ---
//		fil := OpenFile(tspy, "testdata/file.txt")
//
//		// --- Then ---
//		tspy.AssertExpectations()
//		assert.ErrorIs(t, os.ErrClosed, fil.Close())
//	})
//
// If the optional argument expectHelpers is provided the [Spy.ExpectHelpers]
// will be called with it. See [Spy.ExpectHelpers] method documentation for
// details.
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

// ExpectCleanups sets number of expected calls to [Spy.Cleanup] method.
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

// ExpectError sets expectation that HUT should call at least once one of the
// [Spy.Error] or [Spy.Errorf] methods.
func (spy *Spy) ExpectError() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantFailed {
		spy.panicked = true
		panic("cannot use ExpectError and ExpectFail in the same time")
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

// ExpectFatal sets expectation that HUT should call at least once one of the
// [Spy.Fatal] or [Spy.Fatalf] methods.
func (spy *Spy) ExpectFatal() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantFailed {
		spy.panicked = true
		panic("cannot use ExpectFatal and ExpectFail in the same time")
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

// Failed reports whether the HUT called any of the [Spy.Error], [Spy.Errorf],
// [Spy.Fatal], [Spy.Fatalf] or FailNow methods. It's worth noting that this
// method returning false DOES NOT mean the Spy expectations were met. The HUT
// may haven never called the methods listed previously but the spy itself
// didn't met expectations.
func (spy *Spy) Failed() bool {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(mockedFailedCall)
	return spy.haveFatal || spy.haveError
}

// ExpectHelpers sets expectation how many times HUT should call [Spy.Helper]
// method. The value -1 means the [Spy.Helper] method must be run at least once.
//
// Method will panic if the cnt value is less than -1 or the method is called
// more than once.
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

// ExpectSetenv sets expectation that given environment variable is set by the
// HUT.
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

// ExpectSkipped sets expectation that HUT will skip the test.
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

// ExpectTempDir sets expectation the HUT should call [Spy.TempDir] cnt number
// of times. If cnt is -1 the [Spy.TempDir] method can be called any number of
// times.
func (spy *Spy) ExpectTempDir(cnt int) *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	spy.wantTempDirCnt = cnt
	return spy
}

// GetTempDir returns Nth (zero indexed) temporary directory path returned by
// [Spy.TempDir] method. It will fail the test if the [Spy.TempDir] method was
// never called or index of the directory is invalid.
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

// IgnoreLogs instruct Spy to ignore checking logged messages. Method will
// panic if any of the Spy.ExpectLog* methods were already called.
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

// ExpectLog sets expectation the HUT should call one of the [Spy.Log] or
// [Spy.Logf] methods with given message. The expected message is constructed
// using format and args arguments, which are the same as in [fmt.Sprintf]. The
// matcher strategy is used to match the message.
//
// Method call will panic if [Spy.IgnoreLogs] was called before.
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

// ExpectLogEqual sets expectation the HUT should call one of the Log* methods.
// The expected message is constructed using format and args arguments which
// are the same as in [fmt.Sprintf]. The [Equal] strategy is used to match
// messages.
//
// Method call will panic if [Spy.IgnoreLogs] was called before.
func (spy *Spy) ExpectLogEqual(format string, args ...any) *Spy {
	return spy.ExpectLog(Equal, format, args...)
}

// ExpectLogContain sets expectation the HUT should call one of the Log*
// methods. The expected message is constructed using format and args arguments
// which are the same as in [fmt.Sprintf]. The [Contains] strategy is used to
// match log message.
//
// Method call will panic if [Spy.IgnoreLogs] was called before.
func (spy *Spy) ExpectLogContain(format string, args ...any) *Spy {
	return spy.ExpectLog(Contains, format, args...)
}

// ExpectLogNotContain sets expectation the HUT should call one of the Log*
// methods. The expected message is constructed using format and args arguments
// which are the same as in [fmt.Sprintf]. The [NotContains] strategy is used
// to match log messages.
//
// Method call will panic if [Spy.IgnoreLogs] was called before.
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

// ExpectedNames sets expectation the HUT should call Name cnt number of times.
func (spy *Spy) ExpectedNames(cnt int) *Spy {
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

// ExpectFail sets expectation the HUT should call one of the Fatal* or Error*
// methods.
func (spy *Spy) ExpectFail() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(expectCall)
	if spy.wantFatal {
		spy.panicked = true
		panic("cannot use ExpectFail and ExpectFatal in the same time")
	}
	if spy.wantError {
		spy.panicked = true
		panic("cannot use ExpectFail and ExpectError in the same time")
	}
	spy.wantFailed = true
	return spy
}

// Close closes the instance. You cannot add any expectations to closed
// instance.
func (spy *Spy) Close() *Spy {
	spy.mx.Lock()
	defer spy.mx.Unlock()
	spy.tt.Helper()
	spy.checkState(closeCall)
	spy.closed = true
	return spy
}

// Finish marks the end of the test. It can be called by hand, or it's called
// automatically by a cleanup function as described in [New]. After Finish is
// called most of the Spy methods will panic when called - check specific
// method documentation for details.
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

// AssertExpectations asserts all expectations and returns true on success,
// false otherwise. Each failed expectation is logged using tt instance.
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

// assertExpectations has the same logic as AssertExpectations but
// assumes lock has been acquired by the caller.
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
	ret = spy.checkCallMaybeCnt("TempDir", spy.wantTempDirCnt, len(spy.haveTempDirs))
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
// [ExpectSkipped] method.
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

// assertFailed asserts HUT reacted according to expectation set by [ExpectFail]
// method. If [ExpectFail] method was not called this method will always return
// true.
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

// assertError asserts HUT reacted according to expectation set by [ExpectError]
// method.
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

// assertFatal asserts HUT reacted according to expectation set by [ExpectFatal]
// method.
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

// checkCallCnt checks method with name was called expected number of times.
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

// checkCallMaybeCnt checks method with name was called expected number of
// times. If want is equal to -1 the return is always true, otherwise want and
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
// [ExpectLog], [ExpectLogContain], [ExpectLogEqual] methods.
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

// assertEnv asserts HUT reacted according to expectation set by
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
