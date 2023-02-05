package ensurer

import "github.com/JosiahWitt/ensure/internal/testctx"

// Run fn as a subtest called name.
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
