package ensurepkg

import (
	"reflect"

	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
	"github.com/golang/mock/gomock"
)

type (
	erkTableInvalid struct{ erk.DefaultKind }
	erkEntriesGroup struct{ erk.DefaultKind }
	erkEntryInvalid struct{ erk.DefaultKind }
)

var (
	errInvalidTableType      = erk.New(erkTableInvalid{}, "Expected a slice or array for the table, got {{type .table}}")
	errInvalidEntryType      = erk.New(erkTableInvalid{}, "Expected entry in table to be a struct, got {{type .entry}}")
	errMissingNameField      = erk.New(erkTableInvalid{}, "Name field does not exist on struct in table")
	errInvalidNameFieldType  = erk.New(erkTableInvalid{}, "Name field in struct in table is not a string")
	errMocksNotStructPointer = erk.New(erkTableInvalid{}, "Mocks field should be a pointer to a struct, got {{type .mocksField}}")

	errEntriesInvalid     = erk.New(erkEntriesGroup{}, "Errors encountered while building table:")
	errEntryMissingName   = erk.New(erkEntryInvalid{}, "table[{{.index}}]: Name not set for item")
	errEntryDuplicateName = erk.New(erkEntryInvalid{}, "table[{{.index}}]: duplicate Name found; first occurrence was table[{{.firstIndex}}].Name: {{.name}}")
)

type tableEntry struct {
	index int
	name  string

	rawEntry  reflect.Value
	mockSetup func(c *Chain)
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

	for _, entry := range entries {
		entry := entry // Pin range variable

		c := e(nil)
		c.t.Helper()
		c.run(entry.name, func(ensure Ensure) {
			c := ensure(nil)
			c.t.Helper()
			c.markRun()

			entry.mockSetup(c)
			fn(ensure, entry.index)
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

		entry.index = i
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

	entry := &tableEntry{
		rawEntry:  rawEntry,
		mockSetup: func(c *Chain) {}, // Initialize so it is always safe to call
	}

	name, err := entry.extractTableEntryName()
	if err != nil {
		return nil, err
	}
	entry.name = name

	if err := entry.prepareMocks(); err != nil {
		return nil, err
	}

	return entry, nil
}

func (entry *tableEntry) extractTableEntryName() (string, error) {
	entryName, ok := entry.fieldByName("Name")
	if !ok {
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

func (entry *tableEntry) fieldByName(name string) (reflect.Value, bool) {
	field := entry.rawEntry.FieldByName(name)
	zeroValue := reflect.Value{}
	return field, field != zeroValue
}

func (entry *tableEntry) prepareMocks() error {
	return entry.prepareMocksStruct()
}

func (entry *tableEntry) prepareMocksStruct() error {
	entryMocks, ok := entry.fieldByName("Mocks")
	if !ok {
		return nil
	}

	if entryMocks.Kind() != reflect.Ptr || entryMocks.Type().Elem().Kind() != reflect.Struct {
		return erk.WithParams(errMocksNotStructPointer, erk.Params{"mocksField": entryMocks.Interface()})
	}

	// Create new Mocks struct
	entryMocks.Set(reflect.New(entryMocks.Type().Elem()))

	// Ensure everything is correct during preparation, so we can report errors early
	mockEntries := []reflect.Value{}
	controllerType := reflect.TypeOf(&gomock.Controller{})
	for i := 0; i < entryMocks.Elem().NumField(); i++ {
		mockEntry := entryMocks.Elem().Field(i)

		// Mocks should have a NEW method, to allow creating the mock
		newMethod := mockEntry.MethodByName("NEW")
		zeroValue := reflect.Value{}
		if newMethod == zeroValue {
			continue
		}

		// NEW signature should be:
		//  func (m *MockXYZ) NEW(ctrl *gomock.Controller) *MockXYZ { ... }
		newMethodType := newMethod.Type()
		if newMethodType.NumIn() != 1 || newMethodType.In(0) != controllerType {
			continue
		}

		if newMethodType.NumOut() != 1 || newMethodType.Out(0) != mockEntry.Type() {
			continue
		}

		mockEntries = append(mockEntries, mockEntry)
	}

	// At this point, everything should be correct, so we can blindly execute without worrying about types
	entry.mockSetup = func(c *Chain) {
		gomockCtrl := reflect.ValueOf(c.gomockController())

		for _, mockEntry := range mockEntries {
			returns := mockEntry.MethodByName("NEW").Call([]reflect.Value{gomockCtrl})
			mockEntry.Set(returns[0])
		}
	}

	return nil
}
