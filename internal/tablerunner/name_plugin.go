package tablerunner

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

const nameField = "Name"

type namePlugin struct{}

func (*namePlugin) ParseEntryType(entryType reflect.Type) (plugins.TableEntryHooks, error) {
	name, ok := entryType.FieldByName(nameField)
	if !ok {
		return nil, stringerr.Newf("Required Name field does not exist on struct in table")
	}

	if name.Type.Kind() != reflect.String {
		return nil, stringerr.Newf("Required Name field in struct in table is not a string")
	}

	return &nameEntryHooks{existingNames: make(map[string]int)}, nil
}

type nameEntryHooks struct {
	plugins.NoopAfterEntry

	existingNames map[string]int
}

func (h *nameEntryHooks) BeforeEntry(ctx *testctx.Context, entryValue reflect.Value, i int) error {
	name := entryValue.FieldByName(nameField).String()
	if name == "" {
		return stringerr.Newf("table[%d].Name is empty", i)
	}

	if initialIdx, ok := h.existingNames[name]; ok {
		return stringerr.Newf("table[%d].Name duplicates table[%d].Name: %s", i, initialIdx, name)
	}

	h.existingNames[name] = i

	return nil
}
