package tablerunner

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

// BuiltTable contains the output from building a table via [BuildTable].
type BuiltTable struct {
	tableVal  reflect.Value
	tableType reflect.Type
	isPointer bool

	entryHooks []plugins.TableEntryHooks
}

// Run executes each entry in the table. It uses runEntry to handle wrapping each entry in its own testing.T Run block.
// All plugins are run before and after each entry.
func (bt *BuiltTable) Run(outerT testctx.T, runEntry func(name string, callback func(*testctx.Context, func(int)))) {
	outerT.Helper()

	for i := 0; i < bt.tableVal.Len(); i++ {
		fieldVal := bt.tableVal.Index(i)

		if bt.isPointer {
			fieldVal = fieldVal.Elem()
		}

		name := fieldVal.FieldByName(nameField).String()
		runEntry(name, func(ctx *testctx.Context, callback func(int)) {
			ctx.T.Helper()

			if err := bt.runEntryHooks(ctx, fieldVal, i, plugins.TableEntryHooks.BeforeEntry); err != nil {
				ctx.T.Fatalf(err.Error())
				return
			}

			callback(i)

			if err := bt.runEntryHooks(ctx, fieldVal, i, plugins.TableEntryHooks.AfterEntry); err != nil {
				ctx.T.Fatalf(err.Error())
				return
			}
		})
	}
}

type runEntryHook func(entryHooks plugins.TableEntryHooks, ctx *testctx.Context, entryValue reflect.Value, i int) error

func (bt *BuiltTable) runEntryHooks(ctx *testctx.Context, entryValue reflect.Value, i int, run runEntryHook) error {
	errs := []error{}

	for _, hook := range bt.entryHooks {
		if err := run(hook, ctx, entryValue, i); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return stringerr.NewGroup("Errors running plugins", errs)
	}

	return nil
}
