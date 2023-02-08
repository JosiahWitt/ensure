// Package setupmocks provides a plugin that runs SetupMocks for the provided Mocks.
package setupmocks

import (
	"fmt"
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/id"
	"github.com/JosiahWitt/ensure/internal/reflectensure"
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

		funcType, err := parseSetupMocksField(&setupMocksFunc, &mocksStruct)
		if err != nil {
			return nil, err
		}

		h.hasSetupMocks = true
		h.funcType = funcType
	}

	return h, nil
}

func parseSetupMocksField(setupMocksFunc, mocksStruct *reflect.StructField) (funcType, error) {
	t := setupMocksFunc.Type

	generateError := func() error {
		return stringerr.NewBlock(
			fmt.Sprintf("expected %s field to be one of the following", id.SetupMocks),
			[]error{
				stringerr.Newf("func(m %v)", mocksStruct.Type),
				stringerr.Newf("func(m %v, %s %s)", mocksStruct.Type, id.Ensure, id.EnsuringE),
			},
			fmt.Sprintf("Got: %v", t),
		)
	}

	if t.Kind() != reflect.Func {
		return 0, generateError()
	}

	validDefaultIns := t.NumIn() == 1 && t.In(0) == mocksStruct.Type
	validEnsureIns := t.NumIn() == 2 && t.In(0) == mocksStruct.Type && reflectensure.IsEnsuringE(t.In(1))
	validIns := validDefaultIns || validEnsureIns

	validOuts := t.NumOut() == 0

	if !validIns || !validOuts {
		return 0, generateError()
	}

	switch {
	case validEnsureIns:
		return funcTypeEnsure, nil
	default:
		return funcTypeDefault, nil
	}
}

type funcType int

const (
	funcTypeDefault funcType = iota
	funcTypeEnsure
)

// TableEntryHooks exposes the before and after hooks for each entry in the table.
type TableEntryHooks struct {
	plugins.NoopAfterEntry

	hasSetupMocks bool
	funcType      funcType
}

var _ plugins.TableEntryHooks = &TableEntryHooks{}

// BeforeEntry is called before the test is run for the table entry.
// It calls the SetupMocks function with the Mocks field as input.
func (h *TableEntryHooks) BeforeEntry(ctx testctx.Context, entryValue reflect.Value, i int) error {
	if !h.hasSetupMocks {
		return nil
	}

	v := reflect.Indirect(entryValue)
	setupMocksFunc := v.FieldByName(id.SetupMocks)

	if setupMocksFunc.IsZero() {
		return nil
	}

	mocksField := v.FieldByName(id.Mocks)

	var ins []reflect.Value
	switch h.funcType {
	case funcTypeDefault:
		ins = []reflect.Value{mocksField}
	case funcTypeEnsure:
		ins = []reflect.Value{mocksField, reflect.ValueOf(ctx.Ensure())}
	}

	setupMocksFunc.Call(ins)

	return nil
}
