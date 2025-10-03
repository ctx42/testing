# Contributing to CTX42 Testing

Thank you for your interest in contributing to CTX42 Testing! This guide will help you understand our development process, coding standards, and how to submit contributions.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Standards](#code-standards)
- [Project Structure](#project-structure)
- [Testing Guidelines](#testing-guidelines)
- [Documentation Standards](#documentation-standards)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- A GitHub account

### Quick Start

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/testing.git
   cd testing
   ```
3. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. Make your changes and test them
5. Submit a pull request

## Development Setup

### Project Philosophy

CTX42 Testing is designed with these core principles:

- **Zero Dependencies**: We maintain no external dependencies
- **Modular Design**: Each package has a focused responsibility
- **Excellent Developer Experience**: Clear APIs, comprehensive documentation
- **Self-Testing**: Test helpers that test themselves

### Building and Testing

```bash
# Run all tests
go test ./...

# Run tests with race detection (as in CI)
go test -v -race ./...

# Run tests for a specific package
go test ./pkg/assert

# Generate mocks (if needed)
go generate ./...
```

## Code Standards

### Formatting and Style

- **Indentation**: 4 spaces (configured in `.editorconfig`)
- **Line Length**: 80 characters maximum
- **Go Formatting**: Use `go fmt` and `go vet`
- **File Headers**: All files must include SPDX license headers

### File Header Template

```go
// SPDX-License-Identifier: MIT
//
// Copyright (c) 2025 Rafal Zajac <rzajac@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
```

### Naming Conventions

- **Packages**: Short, lowercase names (e.g., `assert`, `check`, `mock`)
- **Functions**: CamelCase with descriptive names
- **Variables**: CamelCase, avoid abbreviations
- **Constants**: CamelCase or ALL_CAPS for exported constants

## Project Structure

### Main Packages (`pkg/`)

- **`assert/`**: Main assertion toolkit - primary user-facing package
- **`check/`**: Equality checking toolkit used by assert
- **`mock/`**: Mock creation primitives
- **`mocker/`**: Interface mock generator
- **`dump/`**: Value dumping/rendering utilities
- **`goldy/`**: Golden file testing support
- **`kit/`**: Test helpers and utilities
- **`must/`**: Panic-on-error test helpers
- **`notice/`**: Formatted error message utilities
- **`tester/`**: Test helper testing utilities

### Internal Packages (`internal/`)

- **`core/`**: Foundation utilities (nil checking, panics)
- **`diff/`**: Text diffing utilities
- **`affirm/`**: Low-level assertion helpers
- **`cases/`**: Test case utilities
- **`constraints/`**: Type constraints
- **`tstmod/`**: Test module utilities
- **`types/`**: Type utilities

### Adding New Packages

When adding new packages:

1. Create package in appropriate directory (`pkg/` for user-facing, `internal/` for internal)
2. Add comprehensive `README.md` with usage examples
3. Include thorough tests and examples
4. Update main project documentation if needed

## Testing Guidelines

### Test File Organization

- **`*_test.go`**: Standard unit tests
- **`examples_test.go`**: Usage examples (shows up in Go documentation)
- **`all_test.go`**: Comprehensive test suites
- **`*.gld`**: Golden files for expected outputs

### Testing Patterns

#### Self-Testing Pattern

Use the `tester` package to test test helpers:

```go
func TestAssertEqual(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // Given
        tspy := tester.New(t)
        defer tspy.Close()
        
        // When
        result := assert.Equal(tspy, "expected", "expected")
        
        // Then
        if !result {
            t.Error("expected assertion to pass")
        }
        tspy.AssertExpectations()
    })
}
```

#### Golden File Testing

Use golden files for complex output testing:

```go
func TestComplexOutput(t *testing.T) {
    goldy := goldy.New(t)
    
    // Generate output
    output := generateComplexOutput()
    
    // Compare with golden file
    goldy.AssertGolden(output)
}
```

### Test Requirements

- All new features must include tests
- Test coverage should be comprehensive
- Include both positive and negative test cases
- Test edge cases and error conditions
- Use meaningful test names that describe the scenario

## Documentation Standards

### Package Documentation

Each package must have:

1. **README.md** with:
   - Table of Contents
   - Usage examples
   - API documentation
   - Advanced usage patterns

2. **Inline Go documentation**:
   - Package-level documentation
   - Function-level documentation with examples
   - Clear parameter and return value descriptions

### Documentation Example

```go
// Package assert provides a comprehensive assertion toolkit for Go testing.
// It offers fluent APIs for common testing scenarios while maintaining
// excellent error messages and zero external dependencies.
package assert

// Equal asserts that two values are equal using deep comparison.
// It returns true if the assertion passes, false otherwise.
//
// Example:
//   assert.Equal(t, "expected", actual)
//   assert.Equal(t, 42, calculate())
func Equal(t testing.TB, expected, actual any) bool {
    // Implementation
}
```

## Submitting Changes

### Pull Request Process

1. **Create a feature branch** from `master`
2. **Make your changes** following the guidelines above
3. **Add tests** for any new functionality
4. **Update documentation** if needed
5. **Run tests** to ensure everything passes
6. **Submit a pull request** with:
   - Clear description of changes
   - Link to any related issues
   - Test results

### Pull Request Guidelines

- Keep changes focused and atomic
- Write clear, descriptive commit messages
- Include tests for new functionality
- Update documentation as needed
- Follow the existing code style
- Ensure CI passes

### Commit Message Format

```
type(scope): brief description

Longer description if needed, explaining the motivation
for the change and how it addresses the issue.

Closes #123
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Release Process

### Version Management

- Version is stored in the `VER` file
- Follow semantic versioning (MAJOR.MINOR.PATCH)
- Update `CHANGELOG.md` with each release
- Tag releases with `git tag vX.Y.Z`

### Changelog Format

```markdown
## [0.28.1] - 2025-01-XX

### Added
- New feature descriptions

### Changed
- Modifications to existing features

### Fixed
- Bug fixes

### Removed
- Deprecated features
```

## Design Principles

When contributing, keep these principles in mind:

1. **Zero Dependencies**: Never add external dependencies
2. **Modular Design**: Keep packages focused on single responsibilities
3. **Fluent APIs**: Design chainable, readable interfaces
4. **Rich Error Messages**: Provide helpful, descriptive error messages
5. **Performance**: Consider performance implications of changes
6. **Backward Compatibility**: Avoid breaking changes when possible

## Getting Help

- Check existing [Issues](https://github.com/ctx42/testing/issues)
- Review the [documentation](./README.md)
- Look at existing code for patterns and examples
- Ask questions in issue discussions

## Code of Conduct

We are committed to providing a welcoming and inspiring community for all. Please be respectful and constructive in all interactions.

---

Thank you for contributing to CTX42 Testing! Your contributions help make Go testing better for everyone.