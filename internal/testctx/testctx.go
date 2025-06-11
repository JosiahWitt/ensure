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
	Cleanup(f func())
	Parallel()
}

var _ T = &testing.T{}

// WrapEnsure is a function that returns the [ensuring.E] for the provided [T].
// It returns an interface instead of the concrete type to avoid an import cycle.
type WrapEnsure func(T) interface{}

// Context contains scoped test helpers.
type Context interface {
	// T returns the currently in-scope [T].
	T() T

	// Run wraps the [testing.T] Run method, making it mockable.
	Run(name string, fn func(Context))

	// GoMockController returns the [gomock.Controller] relating to the in-scope [T].
	// It memoizes the value for subsequent calls.
	GoMockController() *gomock.Controller

	// Ensure returns the [ensuring.E] relating to the in-scope [T].
	// It returns an interface instead of the concrete type to avoid an import cycle.
	// It memoizes the value for subsequent calls.
	Ensure() interface{}
}

// SyncableContext extends [Context] with the ability to use [synctest.Test] in Go 1.25+.
type SyncableContext interface {
	Context

	// Sync wraps the Go 1.25+ [synctest.Test] function, making it mockable.
	Sync(fn func(Context))
}

type baseContext struct {
	t T

	goMockController *gomock.Controller
	wrapEnsure       WrapEnsure
	ensure           interface{}
}

var _ Context = &baseContext{}

// New creates a new [Context].
func New(t T, wrapEnsure WrapEnsure) Context {
	return &baseContext{t: t, wrapEnsure: wrapEnsure}
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
		wrappedCtx := New(t, ctx.wrapEnsure)
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

// Ensure returns the [ensuring.E] relating to the in-scope [T].
// It returns an interface instead of the concrete type to avoid an import cycle.
// It memoizes the value for subsequent calls.
func (ctx *baseContext) Ensure() interface{} {
	ctx.t.Helper()

	if ctx.ensure == nil {
		ctx.ensure = ctx.wrapEnsure(ctx.t)
	}

	return ctx.ensure
}
