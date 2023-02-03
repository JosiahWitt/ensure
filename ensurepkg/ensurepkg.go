// Package ensurepkg contains the implementation for the ensure test framework.
// Use ensure.New to create a new instance of Ensure.
package ensurepkg

import (
	"runtime"
	"strings"
	"testing"

	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

//nolint:gochecknoglobals // This is stored as a variable so we can override it for tests in init_test.go.
var newTestContext = testctx.New

// T implements a subset of methods on testing.T.
// More methods may be added to T with a minor ensure release.
type T = testctx.T

// Ensure the actual value is correct using Chain.
// Ensure also has methods that can be called directly.
type Ensure func(actual interface{}) *Chain

// Chain assertions to the ensure function call.
type Chain struct {
	t      testctx.T
	ctx    testctx.Context
	actual interface{}
	wasRun bool
}

// InternalCreateDoNotCallDirectly should NOT be called directly.
// Instead use `ensure := ensure.New(t)` to allow for easy test refactoring.
func InternalCreateDoNotCallDirectly(t T) Ensure {
	const validWrapperFilePathSuffix = "/ensure.go"

	_, callerFilePath, _, ok := runtime.Caller(1)
	if !ok {
		t.Helper()
		t.Fatalf("Can't get caller from runtime")
	}

	if !strings.HasSuffix(callerFilePath, validWrapperFilePathSuffix) {
		t.Helper()
		t.Fatalf("Do not call `ensurepkg.InternalCreateDoNotCallDirectly(t)` directly. Instead use `ensure.New(t)`.")
	}

	return wrap(t)
}

// New creates an instance of ensure with the provided testing context.
//
// This allows the `ensure` package to be shadowed by the `ensure` variable,
// while still allowing new instances of ensure to be created.
func (e Ensure) New(t T) Ensure {
	return wrap(t)
}

// Failf fails the test immediately with a formatted message.
// The formatted message follows the same format as the fmt package.
func (e Ensure) Failf(format string, args ...interface{}) {
	c := e(nil)
	c.t.Helper()
	c.markRun()
	c.t.Fatalf(format, args...)
}

// T exposes the test context provided to ensure.New(t).
//
// If an instance of *testing.T was not provided to ensure.New(t), this method cannot be used.
// The test will fail immediately.
func (e Ensure) T() *testing.T {
	c := e(nil)
	c.markRun()

	t, ok := c.t.(*testing.T)
	if !ok {
		c.t.Helper()
		c.t.Fatalf("An instance of *testing.T was not provided to ensure.New(t), thus T() cannot be used.")
	}

	return t
}

// GoMockController exposes a GoMock Controller scoped to the current test context.
// Learn more about GoMock here: https://github.com/golang/mock
func (e Ensure) GoMockController() *gomock.Controller {
	c := e(nil)
	c.markRun()
	return c.ctx.GoMockController()
}

func wrap(t T) Ensure {
	// Created outside the callback, so the same context is used across ensure calls
	ctx := newTestContext(t)

	return func(actual interface{}) *Chain {
		c := &Chain{
			t:      t,
			ctx:    ctx,
			actual: actual,
			wasRun: false,
		}

		// Cleanup should never call Fatalf, otherwise panics are hidden, and
		// the Fatal message is displayed instead, which is really tricky for debugging.
		t.Helper()
		t.Cleanup(func() {
			if !c.wasRun {
				t.Helper()
				t.Errorf("Found ensure(<actual>) without chained assertion.")
			}
		})

		return c
	}
}
