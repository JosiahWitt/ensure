//go:build go1.25

package ensuring

import (
	"github.com/JosiahWitt/ensure/internal/plugins/all"
	"github.com/JosiahWitt/ensure/internal/tablerunner"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

// RunSync runs fn as a subtest called name, using [testing.T.Run].
// It executes the provided callback within [synctest.Test], making
// it easier to test concurrent code or code that uses [time.Sleep].
// See the [synctest] docs for more info.
func (e E) RunSync(name string, fn func(ensure E)) {
	c := e(nil)
	c.t.Helper()
	c.markRun()

	c.ctx.Run(name, func(ctx testctx.Context) {
		t := ctx.T()
		t.Helper()

		syncable := ctx.(testctx.SyncableContext)
		syncable.Sync(func(ctx testctx.Context) {
			t := ctx.T()
			t.Helper()
			ensure := wrap(t)
			fn(ensure)
		})
	})
}

// RunTableByIndexSync runs the table which is a slice (or array) of structs.
// It behaves identically to [E.RunTableByIndex], but it executes each table
// entry within [synctest.Test], making it easier to test concurrent code or
// code that uses [time.Sleep]. See the [synctest] docs for more info.
//
// See [E.RunTableByIndex] for more info on table driven testing.
func (e E) RunTableByIndexSync(table interface{}, fn func(ensure E, i int)) {
	c := e(nil)
	c.t.Helper()
	c.markRun()

	bt, err := tablerunner.BuildTable(table, all.TablePlugins())
	if err != nil {
		c.t.Fatalf(err.Error())
		return
	}

	bt.Run(c.ctx, func(ctx testctx.Context, i int) {
		t := ctx.T()
		t.Helper()

		syncable := ctx.(testctx.SyncableContext)
		syncable.Sync(func(ctx testctx.Context) {
			t := ctx.T()
			t.Helper()
			ensure := wrap(t)
			fn(ensure, i)
		})
	})
}
