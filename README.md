# ensure
A balanced test framework for Go 1.13+.

[![Documentation](https://pkg.go.dev/badge/github.com/JosiahWitt/ensure)](https://pkg.go.dev/github.com/JosiahWitt/ensure)
[![CI](https://github.com/JosiahWitt/ensure/workflows/CI/badge.svg)](https://github.com/JosiahWitt/ensure/actions?query=branch%3Amaster+workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/JosiahWitt/ensure)](https://goreportcard.com/report/github.com/JosiahWitt/ensure)
[![codecov](https://codecov.io/gh/JosiahWitt/ensure/branch/master/graph/badge.svg)](https://codecov.io/gh/JosiahWitt/ensure)

## Table of Contents
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Install](#install)
- [About](#about)
- [Overview](#overview)
- [Examples](#examples)
  - [Basic Testing](#basic-testing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## Install
```bash
$ go get github.com/JosiahWitt/ensure
```


## About
Ensure supports Go 1.13+ error comparisons (via [`errors.Is`](https://pkg.go.dev/errors?tab=doc#Is)), and provides easy to read diffs (via [`deep.Equal`](https://pkg.go.dev/github.com/go-test/deep#Equal)).

Ensure was partially inspired by the [`is`](https://github.com/matryer/is) testing mini-framework.


## Overview

Creating a test instance starts by calling:
```go
ensure := ensure.New(t)
```

Then, `ensure` can be used as a function to asset a value is correct, using the pattern `ensure(<actual>).<Method>(<expected>)`. Methods can also be called on `ensure`, using the pattern `ensure.<Method>()`.

## Examples

### Basic Testing
```go
func TestMyFunction(t *testing.T) {
  ensure := ensure.New(t)
  ...

  t.Run("my subtest", func(t *testing.T) {
    ensure := ensure.New(t) // This is using the shadowed version of ensure, and can easily be refactored
    ...

    ensure("abc").Equals("abc") // To ensure a value is correct, use ensure as a function
    ensure.Fail() // Methods can be called directly on ensure
  })
}
```
