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

	//nolint:thelper // We already tag it as a test helper
	c.t.Run(name, func(t *testing.T) {
		ensure := wrap(t)
		fn(ensure)
	})
}
