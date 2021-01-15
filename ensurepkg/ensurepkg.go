// Package ensurepkg contains the implementation for the ensure test framework.
// Use ensure.New to create a new instance of Ensure.
package ensurepkg

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// T implements a subset of methods on testing.T.
// More methods may be added to T with a minor ensure release.
type T interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
	Fail()
	Helper()
}

// Ensure the actual value is correct using Chain.
// Ensure also has methods that can be called directly.
type Ensure func(actual interface{}) Chain

// Chain assetions to the ensure function call.
type Chain struct {
	t      T
	actual interface{}
}

// InternalCreateDoNotCallDirectly should NOT be called directly.
// Instead use `ensure := ensure.New(t)` to allow for easy test refactoring.
func InternalCreateDoNotCallDirectly(t T) Ensure {
	const validWrapperFilePathSuffix = "/ensure.go"

	_, callerFilePath, callerLineNumber, ok := runtime.Caller(1)
	if !ok {
		panic("Can't get caller from runtime")
	}

	if !strings.HasSuffix(callerFilePath, validWrapperFilePathSuffix) {
		panic(fmt.Sprintf("Do not call ensurepkg.New directly. Instead use ensure.New. Called ensurepkg.New from: %v:%v", callerFilePath, callerLineNumber))
	}

	return wrap(t)
}

// Fail the test directly.
func (e Ensure) Fail() {
	c := e(nil)
	c.t.Helper()
	c.t.Fail()
}

// T exposes the test context provided to ensure.New(t).
func (e Ensure) T() T {
	return e(nil).t
}

func wrap(t T) Ensure {
	return func(actual interface{}) Chain {
		return Chain{
			t:      t,
			actual: actual,
		}
	}
}
