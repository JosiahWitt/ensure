// Package tablerunner is an internal package that runs table-driven tests.
package tablerunner

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

// BuildTable prepares the table to be run, surfacing any errors encountered while parsing the table.
func BuildTable(rawTable interface{}, tablePlugins []plugins.TablePlugin) (*BuiltTable, error) {
	tableVal := reflect.ValueOf(rawTable)
	if tableVal.Kind() != reflect.Array && tableVal.Kind() != reflect.Slice {
		return nil, stringerr.Newf("Expected a slice or array for the table, got %T", rawTable)
	}

	tableType := tableVal.Type()
	entryType, isPointer, err := unpackStruct(tableType)
	if err != nil {
		return nil, err
	}

	defaultPlugins := []plugins.TablePlugin{&namePlugin{}}
	allPlugins := append(defaultPlugins, tablePlugins...) //nolint:gocritic
	entryPlugins := make([]plugins.TableEntryPlugin, 0, len(allPlugins))
	pluginErrs := []error{}

	for _, plugin := range allPlugins {
		entryPlugin, err := plugin.ParseEntryType(entryType)
		if err != nil {
			pluginErrs = append(pluginErrs, err)
			continue
		}

		entryPlugins = append(entryPlugins, entryPlugin)
	}

	if len(pluginErrs) > 0 {
		return nil, stringerr.NewGroup("Errors parsing table", pluginErrs)
	}

	return &BuiltTable{
		tableVal:  tableVal,
		tableType: tableType,
		isPointer: isPointer,

		plugins: entryPlugins,
	}, nil
}

func unpackStruct(tableType reflect.Type) (reflect.Type, bool, error) {
	tableTypeElem := tableType.Elem()

	if tableTypeElem.Kind() == reflect.Struct {
		return tableTypeElem, false, nil
	}

	if tableTypeElem.Kind() == reflect.Ptr && tableTypeElem.Elem().Kind() == reflect.Struct {
		return tableTypeElem.Elem(), true, nil
	}

	return nil, false, stringerr.Newf("Expected entry in table to be a struct or a pointer to a struct, got %s", tableTypeElem)
}
