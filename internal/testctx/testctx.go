// Package testctx provides a context containing scoped test helpers.
package testctx

import "github.com/golang/mock/gomock"

// Context contains scoped test helpers.
type Context struct {
	T T

	GoMockController func() *gomock.Controller
}

// T is a minimal implementation of [testing.T] that may expand whenever a new method is needed.
type T interface {
	Fatalf(format string, args ...interface{})
	Helper()
}
