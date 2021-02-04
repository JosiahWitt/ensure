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
  - [Library](#library)
  - [CLI](#cli)
- [About](#about)
- [Overview](#overview)
- [Configuring the CLI](#configuring-the-cli)
- [Examples](#examples)
  - [Basic Testing](#basic-testing)
  - [Table Driven Testing](#table-driven-testing)
  - [Table Driven Testing with Mocks](#table-driven-testing-with-mocks)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## Install
### Library
```bash
$ go get github.com/JosiahWitt/ensure
```

### CLI
```bash
# Before Go 1.16
$ go get github.com/JosiahWitt/ensure-cli/cmd/ensure
$ go get github.com/golang/mock/mockgen

# After Go 1.16
$ go install github.com/JosiahWitt/ensure-cli/cmd/ensure
$ go install github.com/golang/mock/mockgen
```

Note: [`mockgen`](https://github.com/golang/mock) is required to generate mocks.


## About
Ensure supports Go 1.13+ error comparisons (using [`errors.Is`](https://pkg.go.dev/errors?tab=doc#Is)), and provides easy to read diffs (using [`deep.Equal`](https://pkg.go.dev/github.com/go-test/deep#Equal)).
Ensure also [supports mocks](#table-driven-testing-with-mocks) using [GoMock](https://github.com/golang/mock).

Ensure was partially inspired by the [`is`](https://github.com/matryer/is) testing mini-framework.


## Overview

Creating a test instance starts by calling:
```go
ensure := ensure.New(t)
```

Then, `ensure` can be used as a function to asset a value is correct, using the pattern `ensure(<actual>).<Method>(<expected>)`. Methods can also be called on `ensure`, using the pattern `ensure.<Method>()`.


## Configuring the CLI
The `ensure` CLI is configured using a `.ensure.yml` file which is located in the root of your Go Module (next to the `go.mod` file).
Source code for the `ensure` CLI is located in the [`ensure-cli` repo](https://github.com/JosiahWitt/ensure-cli).

Here is an example `.ensure.yml` file:

```yaml
mocks:
  # Used as the directory path relative to the root of the module
  # for any interfaces that are not within internal directories.
  # Optional, defaults to "internal/mocks".
  primaryDestination: internal/mocks

  # Used as the directory path relative to internal directories within the project.
  # Optional, defaults to "mocks".
  internalDestination: mocks

  # Packages with interfaces for which to generate mocks
  packages:
    - path: github.com/my/app/some/pkg
      interfaces: [Iface1, Iface2]
```


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
```

### Table Driven Testing with Mocks
Mocks can be generated by running `ensure generate mocks`, which wraps [GoMock](https://github.com/golang/mock).
To install the `ensure` CLI, see the [Install section](#install).

```go
// db/db.go
type DB interface {
  Put(id string, data interface{}) error
  ...
}

// user/user.go
type UserStorage struct {
  DB db.DB
  ...
}

type User struct {
  ID    string
  Name  string
  ...
}

func (s *UserStorage) Save(ctx context.Context, u *User) error { ... }

// user/user_test.go
func TestTableDrivenMocksExample(t *testing.T) {
  ensure := ensure.New(t)

  type Mocks struct {
    mocksets.DefaultMocks
    DB *mock_db.MockDB // Mock of the db.DB interface generated by `ensure generate mocks`
  }

  table := []struct {
    Name          string
    Input         *user.User
    ExpectedError error

    Mocks      *Mocks            // Mocks to automatically initialize
    SetupMocks func(*Mocks)      // Optional function to allow for mock setup
    Subject    *user.UserStorage // Optional subject containing interfaces with which to assign the mocks
  }{
    {
      Name:    "with valid user",
      Input:   &user.User{
        ID:   "my-id",
        Name: "Mary",
      },
      SetupMocks: func(m *Mocks) {
        m.DB.EXPECT().Put("my-id", &user.User{
          ID:   "my-id",
          Name: "Mary",
        })
      },
    },
    {
      Name:    "with missing ID",
      Input:   &user.User{
        ID:   "",
        Name: "Mary",
      },
      SetupMocks: func(m *Mocks) {
        m.DB.EXPECT().Put("", &user.User{
          ID:   "",
          Name: "Mary",
        }).Return(errors.New("missing ID"))
      },
      ExpectedError: user.ErrSavingUser,
    },
  }

  ensure.RunTableByIndex(table, func(ensure Ensure, i int) {
    entry := table[i]

    err := entry.Subject.Save(entry.Mocks.Context, entry.Input)
    ensure(err).IsError(entry.ExpectedError)
  })
}

// mocksets/mocksets.go
type DefaultMocks struct {
  // Tag suppresses warning when it isn't used in the Subject
  Context *mockctx.MockContext `ensure:"ignoreunused"`
}

// mockctx/mockctx.go
type MockContext struct { context.Context }

// NEW method allows creating a MockContext, and is automatically called by ensure.
func (*MockContext) NEW() *MockContext {
  return &MockContext{Context: context.Background()}
}
```
