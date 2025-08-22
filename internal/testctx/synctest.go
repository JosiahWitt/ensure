//go:build go1.25

package testctx

import (
	"testing"
	"testing/synctest"
)

type syncable interface {
	sync(t T, fn func(t *testing.T))
}

type baseSync struct{}

func (baseSync) sync(t T, fn func(t *testing.T)) {
	synctest.Test(t.(*testing.T), fn) //nolint:forcetypeassert // In practice, this will always be true.
}

//nolint:gochecknoglobals // This is a non-exported global variable so it can be tested.
var sync syncable = baseSync{}

// Sync wraps the Go 1.25+ [synctest.Test] function, making it mockable.
func (ctx *baseContext) Sync(fn func(Context)) {
	ctx.t.Helper()

	sync.sync(ctx.t, func(t *testing.T) {
		t.Helper()
		wrappedCtx := New(t, ctx.wrapEnsure)
		fn(wrappedCtx)
	})
}
