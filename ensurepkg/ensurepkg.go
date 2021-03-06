// Package ensurepkg contains the implementation for the ensure test framework.
// Use ensure.New to create a new instance of Ensure.
package ensurepkg

import (
	"runtime"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

// T implements a subset of methods on testing.T.
// More methods may be added to T with a minor ensure release.
type T interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
	Helper()
	Cleanup(func())
}

// Ensure the actual value is correct using Chain.
// Ensure also has methods that can be called directly.
type Ensure func(actual interface{}) *Chain

// Chain assetions to the ensure function call.
type Chain struct {
	t      T
	actual interface{}
	wasRun bool

	memoized *memoized
}

type memoized struct {
	goMockController *gomock.Controller
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
func (e Ensure) T() T {
	c := e(nil)
	c.markRun()
	return c.t
}

// GoMockController exposes a GoMock Controller scoped to the current test context.
// Learn more about GoMock here: https://github.com/golang/mock
func (e Ensure) GoMockController() *gomock.Controller {
	c := e(nil)
	c.markRun()
	return c.gomockController()
}

func wrap(t T) Ensure {
	memoized := &memoized{}

	return func(actual interface{}) *Chain {
		c := &Chain{
			t:        t,
			actual:   actual,
			memoized: memoized,
		}

		// Cleanup should never call Fatalf, otherwise panics are hidden, and
		// the Fatal message is displayed instead, which is really tricky for debugging.
		t.Helper()
		t.Cleanup(func() {
			if c.memoized.goMockController != nil {
				c.memoized.goMockController.Finish()
			}

			if !c.wasRun {
				t.Helper()
				t.Errorf("Found ensure(<actual>) without chained assertion.")
			}
		})

		return c
	}
}

func (c *Chain) gomockController() *gomock.Controller {
	if c.memoized.goMockController != nil {
		return c.memoized.goMockController
	}

	c.memoized.goMockController = gomock.NewController(c.t)
	return c.memoized.goMockController
}
