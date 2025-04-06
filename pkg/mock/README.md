# Ctx42 Testing Module 🚀

Welcome to the Ctx42 Testing Module, an open-source toolkit designed to
supercharge your Go testing experience! This growing collection of libraries,
crafted from years of Go development expertise, empowers you to write reliable,
readable, and efficient tests with ease.

[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/testing)](https://goreportcard.com/report/github.com/ctx42/testing)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/testing)
![Tests](https://github.com/ctx42/testing/actions/workflows/go.yml/badge.svg?branch=master)

Dive into my blog posts for a behind-the-scenes look at how this module is
evolving. I’d love your feedback—star the repo, share your thoughts, or
contribute to shape the future of Go testing!

<!-- TOC -->
<!-- TOC -->

# Why Ctx42 Testing?

The Ctx42 Testing Module is your go-to solution for streamlined Go testing.
Whether you’re debugging a complex codebase or writing your first unit test,
this module offers intuitive tools to ensure your code is rock-solid. Built for
developers by a passionate Go enthusiast, it’s designed to make testing fast,
fun, and frustration-free.

## Simple and Lightweight
S
ay goodbye to bloated dependencies! Ctx42 Testing Module is dependency-free,
keeping your project lean and conflict-free. Enjoy:

- A **fluent**, **chainable API** for effortless test writing.
- **Clear**, **descriptive error messages** to pinpoint issues instantly.
- **Comprehensive documentation** packed with practical examples to get you up to speed.

Focus on coding, not wrestling with complex setups.

## Modular and Flexible

Mix and match packages to suit your project’s needs:

- Write clean assertions with `assert`.
- Craft powerful mocks with mock and `mocker`.
- Keep tests concise with `tstkit`.

The extensible design lets you build custom helpers, ensuring the module grows
with your requirements. No bloat, just the tools you need.

## Get Started

Install the module in seconds:

```shell
go get github.com/ctx42/testing
```

Ready to explore? Check out the package READMEs and `examples_test.go` files for
hands-on demos.

## Packages

#### Core Testing Packages

Power your test cases with these essentials:

- [assert](pkg/assert/README.md) Robust assertion toolkit for confident testing.
- [check](pkg/check/README.md) Equality checks that power `assert`.
- [mock](pkg/mock) Primitives for crafting interface mocks with ease.
- [must](pkg/must/README.md) Panic-on-error helpers for streamlined tests.

#### Infrastructure Packages

Build your own testing tools with these utilities:

- [dump](pkg/dump/README.md) Configurable renderer for any type to string.
- [notice](pkg/notice/README.md) Create polished assertion messages.
- [tester](pkg/tester/README.md) Utilities for testing your custom helpers.

Each package includes detailed docs and examples—click the links to dive in!

---

# Join the Journey!

This project is a work in progress, and your input can make it even better. Try
it out, open an issue, or submit a PR to help shape a testing toolkit the Go
community loves. Let’s build something awesome together!
