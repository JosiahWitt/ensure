package ensurepkg

import (
	"github.com/JosiahWitt/ensure/internal/plugins/all"
	"github.com/JosiahWitt/ensure/internal/tablerunner"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

// RunTableByIndex runs the table which is a slice (or array) of structs.
// The struct must have a "Name" field which is a unique string describing each test.
// The fn is executed for each entry, with a scoped ensure instance and an index for an entry in the table.
//
// For example:
//
//	table := []struct {
//	  Name    string
//	  Input   string
//	  IsEmpty bool
//	}{
//	  {
//	    Name:    "with non empty input",
//	    Input:   "my string",
//	    IsEmpty: false,
//	  },
//	  {
//	    Name:    "with empty input",
//	    Input:   "",
//	    IsEmpty: true,
//	  },
//	}
//
//	ensure.RunTableByIndex(table, func(ensure Ensure, i int) {
//	  entry := table[i]
//
//	  isEmpty := strs.IsEmpty(entry.Input)
//	  ensure(isEmpty).Equals(entry.IsEmpty)
//	})
//
// Support for mocks is also included.
// Please see the README for an example.
func (e Ensure) RunTableByIndex(table interface{}, fn func(ensure Ensure, i int)) {
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
		ensure := wrap(t)

		fn(ensure, i)
	})
}
