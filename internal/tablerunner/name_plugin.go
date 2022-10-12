package tablerunner

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

const nameField = "Name"

type namePlugin struct{}

func (*namePlugin) ParseEntryType(entryType reflect.Type) (plugins.TableEntryPlugin, error) {
	name, ok := entryType.FieldByName(nameField)
	if !ok {
		return nil, stringerr.Newf("Required Name field does not exist on struct in table")
	}

	if name.Type.Kind() != reflect.String {
		return nil, stringerr.Newf("Required Name field in struct in table is not a string")
	}

	return &nameEntryPlugin{existingNames: make(map[string]int)}, nil
}

type nameEntryPlugin struct {
	existingNames map[string]int
}

func (p *nameEntryPlugin) ParseEntryValue(entryValue reflect.Value, i int) (plugins.TableEntryHooks, error) {
	name := entryValue.FieldByName(nameField).String()
	if name == "" {
		return nil, stringerr.Newf("table[%d].Name is empty", i)
	}

	if initialIdx, ok := p.existingNames[name]; ok {
		return nil, stringerr.Newf("table[%d].Name duplicates table[%d].Name: %s", i, initialIdx, name)
	}

	p.existingNames[name] = i

	return &nameEntryPluginHooks{}, nil
}

type nameEntryPluginHooks struct{}

func (*nameEntryPluginHooks) BeforeEntry() error { return nil }
func (*nameEntryPluginHooks) AfterEntry() error  { return nil }
