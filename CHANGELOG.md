## v0.28.2 (Fri, 11 Jul 2025 10:37:31 UTC)
- fix: panicking dumper for a special case where value is a nil error interface.
- style: remove unreachable code.
- doc: update documentation.

## v0.28.1 (Sun, 06 Jul 2025 20:36:47 UTC)
- Fix the memfs package name.

## v0.28.0 (Sun, 06 Jul 2025 20:25:45 UTC)
- Add check.Cap and assert.Cap functions.
- Change how test kits are organized. Now kits will be put in packages of "kit" package and be organized by "topic". The change also introduces changes to ErrReader and ErrWriter which were split into ErrReader, ErrReadCloser, ErrReadSeeker, ErrReadSeekCloser and ErrWriter, ErrWriteCloser.
- Add kit/memfs package with File struct which mimics `os.File` and implements `fs.File` interface among other useful I/O interfaces.
- Update documentation and change to better struct names.

## v0.27.0 (Thu, 03 Jul 2025 13:29:55 UTC)
- Add tstkit.SHA1Reader, tstkit.SHA1File helpers.

## v0.26.0 (Thu, 03 Jul 2025 13:11:26 UTC)
- Add tstkit.ReadAllFromStart helper.

## v0.25.0 (Thu, 03 Jul 2025 12:47:08 UTC)
- Add tstkit.ErrReader and tstkit.ErrWriter helpers.
- Update documentation.

## v0.24.1 (Mon, 30 Jun 2025 12:36:40 UTC)
- Fix mock generation when a method has a variadic "any" argument.

## v0.24.0 (Fri, 27 Jun 2025 13:38:10 UTC)
- The check.Delta, check.DeltaSlice, check.Epsilon, check.EpsilonSlice use >= instead of > check.

## v0.23.0 (Fri, 27 Jun 2025 09:28:33 UTC)
- Fix check.Epsilon and check.EpsilonSlice to check relative error not delta.
- Add check.Delta, check.DeltaSlice and its equivalents in the assert package.

## v0.22.0 (Thu, 26 Jun 2025 16:14:51 UTC)
- Add check.Increasing and check.Decreasing and their equivalents in assert package.
- Remove already resolved TODO.
- Add check.NotIncreasing and check.NotDecreasing and their equivalents in assert package.
- Add check.Greater and its equivalent in the assert package.
- Add check.Smaller and its equivalent in the assert package.
- Add check.GreaterOrEqual and its equivalent in the assert package.
- Add check.SmallerOrEqual and its equivalent in the assert package.

## v0.21.0 (Thu, 26 Jun 2025 10:25:01 UTC)
- Add check.EpsilonSlice and assert.EpsilonSlice functions to test all slice values are within given epsilon of each other respectively.

## v0.20.2 (Tue, 24 Jun 2025 09:13:01 UTC)
- Update golden file doc.
- Improve the mock error message when at least one call to the mocked method is expected.

## v0.20.1 (Mon, 23 Jun 2025 14:30:35 UTC)
- Fix mock generation where the interface method had a single variadic argument.

## v0.20.0 (Mon, 23 Jun 2025 10:55:47 UTC)
- Add the ability to generate mocks for interfaces which use parametrized types.

## v0.19.1 (Sat, 21 Jun 2025 07:28:20 UTC)
- Fix recursive structure checks by adding reflect.Type to the visited structure.

## v0.19.0 (Fri, 20 Jun 2025 13:28:01 UTC)
- Add dump.Dump.DiffValue which behaves the same way as dump.Dump.Duff but works on reflect.Value instances.
- Add dump.ValNotNil constant representing not nil values.
- Fix custom checkers when the checked type is an unexported struct field. Refactor deep equality test.

## v0.18.2 (Tue, 17 Jun 2025 14:54:44 UTC)
- Add support for string-based timezone comparisons and enhance error handling in check.Zone.

## v0.18.1 (Sun, 15 Jun 2025 09:43:54 UTC)
- Implement stack overflow prevention when comparing nested types.

## v0.18.0 (Sat, 14 Jun 2025 09:13:44 UTC)
- Use core.Value helper to access values in check.Equal.
- Rename check.Check to check.Checker.
- Implement WithinChecker partial application helper and add missing tests.

## v0.17.1 (Fri, 13 Jun 2025 14:50:05 UTC)
- Fix panic when comparing not exported struct fields which are pointers.

## v0.17.0 (Thu, 12 Jun 2025 13:54:48 UTC)
- Add WithZone option to support timezone adjustments in date comparisons.

## v0.16.0 (Thu, 12 Jun 2025 12:57:14 UTC)
- Add core.IsSimpleType helper.
- Add WithCmpBaseTypes option for comparing values with the same base type.

## v0.15.0 (Thu, 12 Jun 2025 10:34:11 UTC)
- Add the ability to define dumpers for a type globally.

## v0.14.2 (Tue, 10 Jun 2025 07:28:48 UTC)
- Fix mocker file name generation when the interface name has all capital letters.

## v0.14.1 (Mon, 09 Jun 2025 14:42:23 UTC)
- Simplify mock error messages. In many instances we are not able to display all the arguments and returns properly, especially in cases when matchers or specifically formated functions are used as return arguments.

## v0.14.0 (Sun, 08 Jun 2025 17:31:52 UTC)
- Fix linting errors.
- The goldy golden files can be Go text templates.

## v0.13.0 (Fri, 06 Jun 2025 13:13:44 UTC)
- Dump values of not exported fields, add special case for type error.

## v0.12.0 (Thu, 05 Jun 2025 13:11:42 UTC)
- Code style.
- Expose dumpers implemented in the dumper package, improve test coverage and code documentation.
- Make check.Equal work in the very similar way to reflect.DeepEqual when checking not exported fields.

## v0.11.0 (Sat, 31 May 2025 15:27:15 UTC)
- Better error messages, code style.
- Rename notice.Notice.Trail() to notice.Notice.SetTrail() to match notice.Notice.SetHeader().
- Rename notice.Notice.{SetData, GetData}() to notice.Notice.{MetaSet, MetaLookup}() to match future metadata interface.
- Improve check.Zone error message.
- Add helper method goldy.Goldy.SetContent setting golden file content from a string.
- Add notice.Pad helper function.
- The core.IsNil can now detect if the nil is wrapped nil.
- Rename Mock.SetData to Mock.MetaSetAll, Mock.GetData to Mock.MetaAll.
- Improve check.NoError error messages, improve check package error messages.
- Prefer "nil" instead of "<nil>" in error messages.
- Code style.
- Improve readability and ease of use in notice.Notice when displaying and creating multi-message instances. Improve value dumps in error messages. Update documentation.

## v0.10.3 (Fri, 23 May 2025 13:54:12 UTC)
- Improve check.Zone error message.

## v0.10.2 (Fri, 23 May 2025 12:14:43 UTC)
- The diff field in check.Within error is always in relation to "want" time. Improve documentation.

## v0.10.1 (Fri, 23 May 2025 11:25:05 UTC)
- Code style.

## v0.10.0 (Fri, 23 May 2025 09:20:29 UTC)
- Update tester.Spy documentation and wording.
- Add tester.Spy.ExamineLog method.

## v0.9.0 (Thu, 22 May 2025 14:29:15 UTC)
- Add clocks to tstkit package.
- Update workflow.
- Improve dumping structs with multi-line string fields.
- Export dump special values.
- Bring to the project code from https://github.com/golang/tools/tree/master/internal/diff and clean it up with some customizations.
- Implement diff in check and dump packages.

## v0.8.0 (Sun, 18 May 2025 14:23:21 UTC)
- Add tstkit package.
- Update documentation and code style.

## v0.7.1 (Sat, 17 May 2025 10:01:03 UTC)
- Code style / documentation.
- Figure out paths on github.
- Figure out paths on GitHub, update GitHub workflows.
- Update documentation. Add TODOs for GitHub Actions.

## v0.7.0 (Fri, 16 May 2025 13:09:51 UTC)
- Fix typo in the Dumper field name.
- Improve dumping long multi line strings. Update documentation.
- Update assertion functions documentation.
- Improve documentation in the tester README.md file.
- Improve documentation in the tester README.md file.
- Add the ability to register global type checkers, add global logger to provide information about interesting events in the test log.
- Update / add the custom type checker and global type checker documentation.
- Simplify the way the panicking assertions and checks are tested, improve documentation and test log messages.
- Code style.
- By convention, all golden files should have gld extension.
- Add new functionalities to the goldy package.
- Refactor check.Options.Trail management and improve code consistency. Export methods on check.Options.Trail so external modules can use them in custom checkers.
- Refactor goldy. Now, to open a file use Open function, the New function was added to create new golden files.
- All check and assert functions must follow the same pattern of arguments first want value then have value.
- Add mocker package.

## v0.6.0 (Fri, 18 Apr 2025 11:16:17 UTC)
- Regenerate assert documentation TOC.
- Link to specific package's README.md file from main README.md.
- Add skip field, element or key documentation.
- Regenerate assert documentation TOC.
- Add ability to set metadata on the notice.Notice instance, update documentation and tests.
- Move core.Len logic to inside check.Len. Clean up code.
- Rename core.DidPanic to core.WillPanic.
- Code style / documentation.
- Code style / documentation.
- Simplify the check.WillPanic returns to return two arguments instead of three.
- Update check.WillPanic documentation.
- Add link to blog.
- The notice package must be independent so it was wrong to use ErrAssert as a base error - renamed it to ErrNotice, also Notice.Wrap now just sets base error instead of wrapping it with current base error.
- The notice.Indent does not indent empty lines and does not trim the input string.
- Update method documentation.
- Use a slice to keep instance of Rows (not strings like before) in wanted order.
- Remove redundant code.
- Add ability in notice package to force row value to start at the next line, and add more documentation.
- Implement Notice.Unwrap.
- Update TOC for notice package.
- Move joined errors decorator from not exported type / helper in check package to exported notice.Join helper so other packages can use it.
- Strings shorter than 200 characters will be dumped as Flat. Added new dump.Dumper option WithFlatStrings to control the lengths of "flat strings".
- Mock package.
- Update README for mock package.
- Update README for mock package.
- Add SPDX (Software Package Data Exchange) identifiers to mock package files.
- Add core.T and core.Spy helper for testing functions in the internal directory which take *testing.T instances.
- Remove TODO.
- Refactor internal packages. Use core.Spy to test functions using *testing.T instance.
- Update Spy documentation.
- Add code.Spy.Log method. Code style improvements.
- Add goldy package.
- Add a requirement for a conscious decision to skip private fields when comparing nested structures.

## v0.5.0 (Thu, 03 Apr 2025 14:50:14 UTC)
- Add internal package diff.
- Add assert and check package.
- Update documentation and examples. Add ability to set dump.Config configuration in check.Options.
- Add check.ErrorIs and assert.ErrorIs.
- Add check.ErrorAs and assert.ErrorAs.
- Add check.ErrorEqual and assert.ErrorEqual.
- Add check.ErrorContain and assert.ErrorContain.
- Add check.Regexp and assert.Regexp.
- Add check.ErrorRegexp and assert.ErrorRegexp.
- Create files for specific check and assert topics.
- Refactor notice path - now it's called trail.
- Add check.FileExist, check.NoFileExist, check.DirExist, check.NoDirExist and assert.NoFileExist, assert.DirExist, assert.NoDirExist.
- Add check.Empty, check.NotEmpty and assert.Empty, assert.NotEmpty.
- Add check.Zero, check.NotZero and assert.Zero, assert.NotZero.
- Add internal helper Same which checks if two generic pointers point to the same memory.
- Add check.Same, check.NotSame and assert.Same, assert.NotSame.
- Add check.Len and assert.Len.
- Add check.True, check.False and assert.True, assert.False.
- Add check.Contain, check.NotContain and assert.Contain, assert.NotContain.
- Add check.Count, check.NotCount and assert.Count, assert.NotCount.
- Add check.Panic, check.NotPanic, check.PanicContain, check.PanicMsg and assert.Panic, assert.NotPanic, assert.PanicContain, assert.PanicMsg.
- Add check.FileContain and assert.FileContain.
- Add check.SameType and assert.SameType.
- Add check.Has, check.HasNo and assert.Has, assert.HasNo.
- Add check.HasKey, check.HasNoKey and assert.HasKey, assert.HasNoKey.
- Add check.HasKeyValue and assert.HasKeyValue.
- Add check.ExitCode and assert.ExitCode.
- Add check.SliceSubset and assert.SliceSubset.
- Change check.Options structure and add default time parsing option.
- Add FormatDates helper.
- Add helper getting time.Time from values of type any, where any may be time.Time, string, int, int64.
- Improve check.getTime error messages.
- Add check.TimeEqual and assert.TimeEqual.
- Add missing test cases to notice.Notice.
- Add check.TimeLoc and assert.TimeLoc. Build options only when needed.
- Add helper check.getDut getting time.Duration from values of type any, where any may be time.Duration, string, int, int64. Rename check and assertion TimeLoc to Zone and TimeEqual to Time.
- Add check.TimeExact and assert.TimeExact.
- Add check.Within and assert.Within.
- Better time check and assert function documentation.
- More readable log messages when asserting / checking comparing dates.
- Remove dead code.
- Update check.getDur helper. Code cleanup. Better error messages.
- Add check.Before and assert.Before.
- Add check.After and assert.After.
- Add check.EqualOrAfter and assert.EqualOrAfter.
- Add check.EqualOrBefore and assert.EqualOrBefore.
- Add ability to configure package-wide defaults.
- Update dump options and tests to use explicit "With" prefixes. Add ability to configure package-wide defaults.
- Add support for configurable recent duration in check package.
- Update documentation, add injectable clock to check.Options.
- Add check.Recent and assert.Recent.
- Better names.
- Improve test coverage.
- Add check.ChannelWillClose and assert.ChannelWillClose.
- Update dump package with better user interface and fixed dump format.
- Code cleanup and style.
- Indent with spaces when dumping values and when rendering notice messages.
- Improve notice package message formating.
- Remove internal diff package.
- Remove internal playground package.
- Update dump package documentation.
- Improve multiple notice messages rendering with the same header.
- Remove internal diff package.
- Add check.Equal which recursively compares two values.
- Fix code linting errors.
- Add assert.Equal.
- Export dump.Printer code printer.
- Refactor dump.Dump - separate dump.Config struct is not needed.
- Rename module.
- Update documentation.
- Use Option slice for equalError helper.
- Use custom byte dumper for equality errors.
- Add check.NotEqual and assert.NotEqual.
- Cleanup.
- Introduce internal core package with helpers used by higher level abstractions.
- Update behaviour for core.DidPanic and simplify affirm.Panic helper.
- Add check.JSON and assert.JSON.
- Add check.MapSubset and assert.MapSubset.
- Add check.MapsSubset and assert.MapsSubset.
- Rename check.SameType and assert.SameType to check.Type and assert.Type.
- Add check.Epsilon and assert.Epsilon.
- Add check.Fields and assert.Fields.
- Fix assert.Epsilon.
- The assert package documentation.
- The assert package documentation.

## v0.4.0 (Tue, 18 Mar 2025 18:20:26 UTC)
- Add Go Report Card and Go Doc icons / links.
- Add Go Test badge.
- Update affirm.Panic helper to handle errors, and other types that might be passed to panic.
- Add internal tstkit package with simple golden file reader.
- Add dump package.
- Update copyright lines to use SPDX-FileCopyrightText.

## v0.3.1 (Sat, 15 Mar 2025 12:00:58 UTC)
- Add work in progress disclaimer and more documentation.

## v0.3.0 (Sat, 15 Mar 2025 11:48:58 UTC)
- Create go.yml.
- Add package must with helper functions which panic on error. Add `affirm.DeeoEqual` helper.
- Add notice package and update documentation.

## v0.2.0 (Fri, 14 Mar 2025 19:26:12 UTC)
- Update documentation.

## v0.1.0 (Fri, 14 Mar 2025 17:37:29 UTC)
- Initial commit.

