package ensurepkg

import (
	"testing"
)

// Run fn as a subtest called name.
func (e Ensure) Run(name string, fn func(ensure Ensure)) {
	c := e(nil)
	c.t.Helper()
	c.run(name, fn)
}

func (c *Chain) run(name string, fn func(ensure Ensure)) {
	c.t.Helper()
	c.markRun()

	c.t.Run(name, func(t *testing.T) {
		t.Helper()
		ensure := wrap(t)
		fn(ensure)
	})
}
