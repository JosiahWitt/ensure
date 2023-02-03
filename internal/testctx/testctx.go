// Package testctx provides a context containing scoped test helpers.
package testctx

import (
	"testing"

	"github.com/golang/mock/gomock"
)

// T is a minimal implementation of [testing.T] that may expand whenever a new method is needed.
type T interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
	Helper()
	Cleanup(func())
}

var _ T = &testing.T{}

// Context contains scoped test helpers.
type Context interface {
	T() T
	Run(name string, fn func(Context))
	GoMockController() *gomock.Controller
}

type baseContext struct {
	t T

	goMockController *gomock.Controller
}

var _ Context = &baseContext{}

// New creates a new [Context].
func New(t T) Context {
	return &baseContext{t: t}
}

// T returns the currently in-scope [T].
func (ctx *baseContext) T() T {
	return ctx.t
}

// Run wraps the [testing.T] Run method, making it mockable.
func (ctx *baseContext) Run(name string, fn func(Context)) {
	ctx.t.Helper()

	ctx.t.Run(name, func(t *testing.T) {
		t.Helper()
		wrappedCtx := New(t)
		fn(wrappedCtx)
	})
}

// GoMockController returns the [gomock.Controller] relating to the in-scope [T].
// It memoizes the value for subsequent calls.
func (ctx *baseContext) GoMockController() *gomock.Controller {
	ctx.t.Helper()

	if ctx.goMockController == nil {
		ctx.goMockController = gomock.NewController(ctx.t)
		ctx.t.Cleanup(ctx.goMockController.Finish)
	}

	return ctx.goMockController
}
