package ensurepkg

import (
	"testing"

	"github.com/JosiahWitt/ensure/internal/testctx"
)

// Run fn as a subtest called name.
func (e Ensure) Run(name string, fn func(ensure Ensure)) {
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

func (c *Chain) deprecatedRun(name string, fn func(ensure Ensure)) {
	c.t.Helper()
	c.markRun()

	c.t.Run(name, func(t *testing.T) {
		t.Helper()
		ensure := wrap(t)
		fn(ensure)
	})
}
