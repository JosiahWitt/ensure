// Package ensurepkg contains the implementation for the ensure test framework.
// Use ensure.New to create a new instance of Ensure.
package ensurepkg

import (
	"runtime"
	"strings"
	"testing"
)

// T implements a subset of methods on testing.T.
// More methods may be added to T with a minor ensure release.
type T interface {
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

func wrap(t T) Ensure {
	return func(actual interface{}) *Chain {
		c := &Chain{
			t:      t,
			actual: actual,
		}

		t.Helper()
		t.Cleanup(func() {
			if !c.wasRun {
				t.Helper()
				t.Fatalf("Found ensure(<actual>) without chained assertion.")
			}
		})

		return c
	}
}
