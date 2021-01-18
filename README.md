# ensure
A balanced test framework for Go 1.14+.

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
  - [Table Driven Testing](#table-driven-testing)

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
func TestBasicExample(t *testing.T) {
  ensure := ensure.New(t)
  ...

  // Methods can be called on ensure, for example, Run:
  ensure.Run("my subtest", func(ensure ensurepkg.Ensure) {
    ...

    // To ensure a value is correct, use ensure as a function:
    ensure("abc").Equals("abc")
    ensure(produceError()).IsError(expectedError)
    ensure(doNotProduceError()).IsNotError()
    ensure(true).IsTrue()
    ensure(false).IsFalse()
    ensure("").IsEmpty()

    // Failing a test directly:
    ensure.Failf("Something went wrong, and we stop the test immediately")
  })
}
```

### Table Driven Testing
```go
func TestTableDrivenExample(t *testing.T) {
  ensure := ensure.New(t)

  table := []struct {
    Name    string
    Input   string
    IsEmpty bool
  }{
    {
      Name:    "with non empty input",
      Input:   "my string",
      IsEmpty: false,
    },
    {
      Name:    "with empty input",
      Input:   "",
      IsEmpty: true,
    },
  }

  ensure.RunTableByIndex(table, func(ensure Ensure, i int) {
    entry := table[i]

    isEmpty := strs.IsEmpty(entry.Input)
    ensure(isEmpty).Equals(entry.IsEmpty)
  })
}
