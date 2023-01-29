// Package setupmocks provides a plugin that runs SetupMocks for the provided Mocks.
package setupmocks

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/id"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

// New creates a new [TablePlugin].
func New() *TablePlugin {
	return &TablePlugin{}
}

// TablePlugin runs SetupMocks for the provided Mocks. This plugin should come after mocks are initialized.
type TablePlugin struct{}

var _ plugins.TablePlugin = &TablePlugin{}

// ParseEntryType is called during the first pass of plugin initialization.
// It is responsible for making sure the types are as expected.
func (t *TablePlugin) ParseEntryType(entryType reflect.Type) (plugins.TableEntryHooks, error) {
	h := &TableEntryHooks{}

	mocksStruct, hasMocks := entryType.FieldByName(id.Mocks)

	setupMocksFunc, hasSetupMocks := entryType.FieldByName(id.SetupMocks)
	if hasSetupMocks {
		if !hasMocks {
			return nil, stringerr.Newf("%s field must be set on the table to use %s", id.Mocks, id.SetupMocks)
		}

		if err := validateSetupMocksFieldType(&setupMocksFunc, &mocksStruct); err != nil {
			return nil, err
		}

		h.hasSetupMocks = true
	}

	return h, nil
}

func validateSetupMocksFieldType(setupMocksFunc, mocksStruct *reflect.StructField) error {
	t := setupMocksFunc.Type

	generateError := func() error {
		return stringerr.Newf("expected %s field to be a func(%v), got: %v", id.SetupMocks, mocksStruct.Type, t)
	}

	if t.Kind() != reflect.Func {
		return generateError()
	}

	invalidIns := t.NumIn() != 1 || t.In(0) != mocksStruct.Type
	invalidOuts := t.NumOut() != 0

	if invalidIns || invalidOuts {
		return generateError()
	}

	return nil
}

// TableEntryHooks exposes the before and after hooks for each entry in the table.
type TableEntryHooks struct {
	plugins.NoopAfterEntry

	hasSetupMocks bool
}

var _ plugins.TableEntryHooks = &TableEntryHooks{}

// BeforeEntry is called before the test is run for the table entry.
// It calls the SetupMocks function with the Mocks field as input.
func (h *TableEntryHooks) BeforeEntry(ctx *testctx.Context, entryValue reflect.Value, i int) error {
	if !h.hasSetupMocks {
		return nil
	}

	v := reflect.Indirect(entryValue)
	setupMocksFunc := v.FieldByName(id.SetupMocks)

	if setupMocksFunc.IsZero() {
		return nil
	}

	mocksField := v.FieldByName(id.Mocks)
	setupMocksFunc.Call([]reflect.Value{mocksField})

	return nil
}
