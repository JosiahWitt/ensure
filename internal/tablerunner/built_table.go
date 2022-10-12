package tablerunner

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

// BuiltTable contains the output from building a table via [BuildTable].
type BuiltTable struct {
	tableVal  reflect.Value
	tableType reflect.Type
	isPointer bool

	plugins []plugins.TableEntryPlugin
}

// Run executes each entry in the table. It uses runEntry to handle wrapping each entry in its own testing.T Run block.
// All plugins are run before and after each entry.
func (bt *BuiltTable) Run(outerT T, runEntry func(name string, callback func(T, func(int)))) {
	outerT.Helper()

	for i := 0; i < bt.tableVal.Len(); i++ {
		fieldVal := bt.tableVal.Index(i)

		if bt.isPointer {
			fieldVal = fieldVal.Elem()
		}

		name := fieldVal.FieldByName(nameField).String()
		runEntry(name, func(t T, callback func(int)) {
			t.Helper()

			entryPlugins, err := bt.buildPlugins(fieldVal, i)
			if err != nil {
				t.Fatalf(err.Error())
				return
			}

			if err := runEntryPlugins(entryPlugins, plugins.TableEntryHooks.BeforeEntry); err != nil {
				t.Fatalf(err.Error())
				return
			}

			callback(i)

			if err := runEntryPlugins(entryPlugins, plugins.TableEntryHooks.AfterEntry); err != nil {
				t.Fatalf(err.Error())
				return
			}
		})
	}
}

func (bt *BuiltTable) buildPlugins(fieldVal reflect.Value, i int) ([]plugins.TableEntryHooks, error) {
	entryPlugins := make([]plugins.TableEntryHooks, 0, len(bt.plugins))
	errs := []error{}

	for _, plugin := range bt.plugins {
		entryPlugin, err := plugin.ParseEntryValue(fieldVal, i)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		entryPlugins = append(entryPlugins, entryPlugin)
	}

	if len(errs) > 0 {
		return nil, stringerr.NewGroup("Errors parsing table entry", errs)
	}

	return entryPlugins, nil
}

func runEntryPlugins(plugins []plugins.TableEntryHooks, run func(plugins.TableEntryHooks) error) error {
	errs := []error{}

	for _, plugin := range plugins {
		if err := run(plugin); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return stringerr.NewGroup("Errors running plugins", errs)
	}

	return nil
}
