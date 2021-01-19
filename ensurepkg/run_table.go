package ensurepkg

import (
	"reflect"

	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
)

type (
	erkTableInvalid struct{ erk.DefaultKind }
	erkEntriesGroup struct{ erk.DefaultKind }
	erkEntryInvalid struct{ erk.DefaultKind }
)

var (
	errInvalidTableType     = erk.New(erkTableInvalid{}, "Expected a slice or array for the table, got {{type .table}}")
	errInvalidEntryType     = erk.New(erkTableInvalid{}, "Expected entry in table to be a struct, got {{type .entry}}")
	errMissingNameField     = erk.New(erkTableInvalid{}, "Name field does not exist on struct in table")
	errInvalidNameFieldType = erk.New(erkTableInvalid{}, "Name field in struct in table is not a string")

	errEntriesInvalid     = erk.New(erkEntriesGroup{}, "Errors encountered while building table:")
	errEntryMissingName   = erk.New(erkEntryInvalid{}, "table[{{.index}}]: Name not set for item")
	errEntryDuplicateName = erk.New(erkEntryInvalid{}, "table[{{.index}}]: duplicate Name found; first occurrence was table[{{.firstIndex}}].Name: {{.name}}")
)

type tableEntry struct {
	name string
}

// RunTableByIndex runs the table which is a slice (or array) of structs.
// The struct must have a "Name" field which is a unique string describing each test.
// The fn is executed for each entry, with a scoped ensure instance and an index for an entry in the table.
//
// For example:
//  table := []struct {
//    Name    string
//    Input   string
//    IsEmpty bool
//  }{
//    {
//      Name:    "with non empty input",
//      Input:   "my string",
//      IsEmpty: false,
//    },
//    {
//      Name:    "with empty input",
//      Input:   "",
//      IsEmpty: true,
//    },
//  }
//
//  ensure.RunTableByIndex(table, func(ensure Ensure, i int) {
//    entry := table[i]
//
//    isEmpty := strs.IsEmpty(entry.Input)
//    ensure(isEmpty).Equals(entry.IsEmpty)
//  })
func (e Ensure) RunTableByIndex(table interface{}, fn func(ensure Ensure, i int)) {
	entries, err := buildTable(table)
	if err != nil {
		c := e(nil)
		c.t.Helper()
		c.markRun()
		c.t.Fatalf(err.Error())
	}

	for i, entry := range entries {
		i := i // Pin range variable

		c := e(nil)
		c.t.Helper()
		c.run(entry.name, func(ensure Ensure) {
			fn(ensure, i)
		})
	}
}

func buildTable(table interface{}) ([]*tableEntry, error) {
	val := reflect.ValueOf(table)
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return nil, erk.WithParam(errInvalidTableType, "table", table)
	}

	entries := make([]*tableEntry, 0, val.Len())
	errGroup := erg.NewAs(errEntriesInvalid)
	for i := 0; i < val.Len(); i++ {
		entry, err := buildTableEntry(val.Index(i))
		if err != nil {
			if erk.IsKind(err, erkEntryInvalid{}) {
				errGroup = erg.Append(errGroup, erk.WithParam(err, "index", i))
				continue
			}

			return nil, err
		}

		entries = append(entries, entry)
	}

	existingNames := make(map[string]int)
	for index, entry := range entries {
		if firstIndex, ok := existingNames[entry.name]; ok {
			errGroup = erg.Append(errGroup, erk.WithParams(errEntryDuplicateName, erk.Params{
				"index":      index,
				"firstIndex": firstIndex,
				"name":       entry.name,
			}))
			continue
		}

		existingNames[entry.name] = index
	}

	if erg.Any(errGroup) {
		return nil, errGroup
	}

	return entries, nil
}

func buildTableEntry(rawEntry reflect.Value) (*tableEntry, error) {
	if rawEntry.Kind() != reflect.Struct {
		return nil, erk.WithParam(errInvalidEntryType, "entry", rawEntry.Interface())
	}

	name, err := extractTableEntryName(rawEntry)
	if err != nil {
		return nil, err
	}

	return &tableEntry{
		name: name,
	}, nil
}

func extractTableEntryName(rawEntry reflect.Value) (string, error) {
	entryName := rawEntry.FieldByName("Name")
	zeroValue := reflect.Value{}
	if entryName == zeroValue {
		return "", errMissingNameField
	}

	if entryName.Kind() != reflect.String {
		return "", errInvalidNameFieldType
	}

	name := entryName.String()
	if name == "" {
		return "", errEntryMissingName
	}

	return name, nil
}
