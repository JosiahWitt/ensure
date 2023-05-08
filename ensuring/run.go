package ensuring

import "github.com/JosiahWitt/ensure/internal/testctx"

// Run runs fn as a subtest called name, using [testing.T.Run].
func (e E) Run(name string, fn func(ensure E)) {
	c := e(nil)
	c.t.Helper()
	c.markRun()

	c.ctx.Run(name, func(ctx testctx.Context) {
		t := ctx.T()
		t.Helper()
		ensure := wrap(t)
		fn(ensure)
	})
}

// RunParallel runs fn as a subtest called name, using [testing.T.Run].
// It calls [testing.T.Parallel] just before calling fn, causing fn to
// be run in parallel with (and only with) other parallel tests. See
// the [testing.T.Parallel] docs for more info.
func (e E) RunParallel(name string, fn func(ensure E)) {
	c := e(nil)
	c.t.Helper()
	c.markRun()

	c.ctx.Run(name, func(ctx testctx.Context) {
		t := ctx.T()
		t.Helper()
		ensure := wrap(t)
		t.Parallel()
		fn(ensure)
	})
}
