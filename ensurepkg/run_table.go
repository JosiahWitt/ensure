package ensurepkg

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

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
	errInvalidTableType     = erk.New(erkTableInvalid{}, "Expected a slice or array for the table, got {{type .table}}")
	errInvalidEntryType     = erk.New(erkTableInvalid{}, "Expected entry in table to be a struct, got {{type .entry}}")
	errMissingNameField     = erk.New(erkTableInvalid{}, "Name field does not exist on struct in table")
	errInvalidNameFieldType = erk.New(erkTableInvalid{}, "Name field in struct in table is not a string")

	errMocksNotStructPointer      = erk.New(erkTableInvalid{}, "Mocks field should be a pointer to a struct, got {{type .mocksField}}")
	errMocksEntryNotStructPointer = erk.New(erkTableInvalid{}, "Mocks.{{.mocksFieldName}} should be a pointer to a struct, got {{type .mockEntry}}")
	errMocksEmbeddedNotStruct     = erk.New(erkTableInvalid{}, "Mocks.{{.mocksFieldName}} should be an embedded struct with no pointers, got {{type .mockEntry}}")
	errMocksNEWMissing            = erk.New(erkTableInvalid{},
		"\nMocks.{{.mocksFieldName}} is missing the NEW method. Expected:\n\tfunc ({{type .expectedReturn}}) NEW(*gomock.Controller) {{type .expectedReturn}}"+
			"\nPlease ensure you generated the mocks using the `ensure mocks generate` command.",
	)
	errMocksNEWInvalidSignature = erk.New(erkTableInvalid{},
		"\nMocks.{{.mocksFieldName}}.NEW has this method signature:\n\t{{type .actualMethod}}\nExpected:\n\tfunc(*gomock.Controller) {{type .expectedReturn}}",
	)
	errMocksDuplicatesFound = erk.New(erkTableInvalid{}, "Found multiple mocks with type '{{type .duplicate}}'; only one mock of each type is allowed")

	errSetupMocksWithoutMocks     = erk.New(erkTableInvalid{}, "SetupMocks field requires the Mocks field")
	errSetupMocksInvalidSignature = erk.New(erkTableInvalid{},
		"\nSetupMocks has this function signature:\n\t{{type .actualSetupMocks}}\nExpected:\n\tfunc({{type .expectedMockParam}})",
	)

	errSubjectNotStructPointer           = erk.New(erkTableInvalid{}, "Subject field should be a pointer to a struct, got {{type .subjectField}}")
	errSubjectFieldMatchingMultipleMocks = erk.New(erkTableInvalid{},
		"Subject.{{.subjectFieldName}} matches multiple mocks; only one mock should exist for each interface: {{.matchedMockTypes}}",
	)

	errEntriesInvalid     = erk.New(erkEntriesGroup{}, "Errors encountered while building table:")
	errEntryMissingName   = erk.New(erkEntryInvalid{}, "table[{{.index}}]: Name not set for item")
	errEntryDuplicateName = erk.New(erkEntryInvalid{}, "table[{{.index}}]: duplicate Name found; first occurrence was table[{{.firstIndex}}].Name: {{.name}}")
)

type tableEntryMockTag string

type tableEntryMock struct {
	fieldName string
	value     reflect.Value
	tag       tableEntryMockTag
}

type tableEntry struct {
	index int
	name  string

	rawEntry   reflect.Value
	setupFuncs []func(c *Chain)

	mocks map[reflect.Type]*tableEntryMock

	// These should be the same for every entry.
	//
	// TODO: Refactor this file so it does more single time pre-computation,
	// and allows raising issues on the table level.
	// See: https://github.com/JosiahWitt/ensure/issues/23
	tableLevelWarnings []string
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
//
// Support for mocks is also included.
// Please see the README for an example.
func (e Ensure) RunTableByIndex(table interface{}, fn func(ensure Ensure, i int)) {
	entries, err := buildTable(table)
	if err != nil {
		c := e(nil)
		c.t.Helper()
		c.markRun()
		c.t.Fatalf(err.Error())
	}

	if len(entries) > 0 && len(entries[0].tableLevelWarnings) > 0 {
		warnings := strings.Join(entries[0].tableLevelWarnings, "\n - ")

		c := e(nil)
		c.markRun()
		c.t.Helper()
		c.t.Logf(
			"\n\n⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️\n\n"+
				"WARNINGS:\n - %s\n\n"+
				"These may or may not be the cause of a problem. If you recently changed an interface, make sure to rerun `ensure mocks generate`.\n\n"+
				"⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️\n\n",

			warnings,
		)
	}

	for _, entry := range entries {
		entry := entry // Pin range variable

		c := e(nil)
		c.t.Helper()
		c.run(entry.name, func(ensure Ensure) {
			c := ensure(nil)
			c.t.Helper()
			c.markRun()

			for _, setupFunc := range entry.setupFuncs {
				setupFunc(c)
			}

			fn(ensure, entry.index)
		})
	}
}

func buildTable(table interface{}) ([]*tableEntry, error) {
	tableReflect := reflect.ValueOf(table)
	if tableReflect.Kind() != reflect.Array && tableReflect.Kind() != reflect.Slice {
		return nil, erk.WithParam(errInvalidTableType, "table", table)
	}

	entries := make([]*tableEntry, 0, tableReflect.Len())
	errGroup := erg.NewAs(errEntriesInvalid)
	for i := 0; i < tableReflect.Len(); i++ {
		entry, err := buildTableEntry(tableReflect.Index(i))
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

	//nolint:exhaustivestruct
	entry := &tableEntry{
		rawEntry: rawEntry,
		mocks:    make(map[reflect.Type]*tableEntryMock),
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
	if err := entry.prepareMocksStruct(); err != nil {
		return err
	}

	if err := entry.prepareSetupMocks(); err != nil {
		return err
	}

	return entry.prepareSubjectStruct()
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

	return entry.fillMocksStruct(entryMocks.Elem())
}

func (entry *tableEntry) fillMocksStruct(entryMocks reflect.Value) error {
	for i := 0; i < entryMocks.NumField(); i++ {
		mockEntry := entryMocks.Field(i)
		mockEntryType := entryMocks.Type().Field(i)
		mockFieldName := mockEntryType.Name

		// Skip unexported fields
		if mockEntryType.PkgPath != "" {
			continue
		}

		// Support embedded structs
		if mockEntryType.Anonymous {
			if mockEntry.Kind() != reflect.Struct {
				return erk.WithParams(errMocksEmbeddedNotStruct, erk.Params{
					"mocksFieldName": mockFieldName,
					"mockEntry":      mockEntry.Interface(),
				})
			}

			if err := entry.fillMocksStruct(mockEntry); err != nil {
				return err
			}

			continue
		}

		if err := entry.prepareMock(mockFieldName, mockEntry); err != nil {
			return err
		}

		entry.mocks[mockEntry.Type()].tag = tableEntryMockTag(mockEntryType.Tag.Get("ensure"))
	}

	return nil
}

func (entry *tableEntry) prepareMock(mockFieldName string, mockEntry reflect.Value) error {
	if mockEntry.Kind() != reflect.Ptr {
		return erk.WithParams(errMocksEntryNotStructPointer, erk.Params{
			"mocksFieldName": mockFieldName,
			"mockEntry":      mockEntry.Interface(),
		})
	}

	// Mocks should have a NEW method, to allow creating the mock
	newMethod := mockEntry.MethodByName("NEW")
	if newMethod == (reflect.Value{}) {
		return erk.WithParams(errMocksNEWMissing, erk.Params{
			"mocksFieldName": mockFieldName,
			"expectedReturn": mockEntry.Interface(),
		})
	}

	// NEW signature should be one of:
	//  func (m *MockXYZ) NEW(ctrl *gomock.Controller) *MockXYZ { ... }
	//  func (m *MockXYZ) NEW() *MockXYZ { ... }
	newMethodType := newMethod.Type()
	controllerType := reflect.TypeOf(&gomock.Controller{}) //nolint:exhaustivestruct
	isInvalidParam := newMethodType.NumIn() > 1 || (newMethodType.NumIn() == 1 && newMethodType.In(0) != controllerType)
	isInvalidReturn := newMethodType.NumOut() != 1 || newMethodType.Out(0) != mockEntry.Type()
	if isInvalidParam || isInvalidReturn {
		return erk.WithParams(errMocksNEWInvalidSignature, erk.Params{
			"mocksFieldName": mockFieldName,
			"actualMethod":   newMethod.Interface(),
			"expectedReturn": mockEntry.Interface(),
		})
	}

	// Save placeholder for mock entry
	if _, alreadyExists := entry.mocks[mockEntry.Type()]; alreadyExists {
		return erk.WithParams(errMocksDuplicatesFound, erk.Params{
			"duplicate": mockEntry.Interface(),
		})
	}

	//nolint:exhaustivestruct
	entry.mocks[mockEntry.Type()] = &tableEntryMock{
		fieldName: mockFieldName,
	}

	// At this point, everything should be correct, so we can blindly execute without worrying about types
	entry.setupFuncs = append(entry.setupFuncs, func(c *Chain) {
		gomockCtrl := reflect.ValueOf(c.gomockController())

		input := []reflect.Value{}
		if newMethodType.NumIn() == 1 {
			input = append(input, gomockCtrl)
		}

		returns := newMethod.Call(input)
		mockInstance := returns[0]
		mockEntry.Set(mockInstance)
		entry.mocks[mockEntry.Type()].value = mockInstance
	})

	return nil
}

func (entry *tableEntry) prepareSetupMocks() error {
	setupMocks, ok := entry.fieldByName("SetupMocks")
	if !ok {
		return nil
	}

	mocks, ok := entry.fieldByName("Mocks")
	if !ok {
		return errSetupMocksWithoutMocks
	}

	if setupMocks.IsNil() {
		return nil
	}

	if setupMocks.Type().NumIn() != 1 || setupMocks.Type().In(0) != mocks.Type() || setupMocks.Type().NumOut() != 0 {
		return erk.WithParams(errSetupMocksInvalidSignature, erk.Params{
			"expectedMockParam": mocks.Interface(),
			"actualSetupMocks":  setupMocks.Interface(),
		})
	}

	// At this point, everything should be correct, so we can blindly execute without worrying about types
	entry.setupFuncs = append(entry.setupFuncs, func(c *Chain) {
		setupMocks.Call([]reflect.Value{mocks})
	})

	return nil
}

//nolint:funlen,cyclop // TODO: Refactor (https://github.com/JosiahWitt/ensure/issues/23)
func (entry *tableEntry) prepareSubjectStruct() error {
	entrySubject, ok := entry.fieldByName("Subject")
	if !ok {
		return nil
	}

	if entrySubject.Kind() != reflect.Ptr || entrySubject.Type().Elem().Kind() != reflect.Struct {
		return erk.WithParams(errSubjectNotStructPointer, erk.Params{"subjectField": entrySubject.Interface()})
	}

	// Create new Subject struct
	entrySubject.Set(reflect.New(entrySubject.Type().Elem()))

	// Ensure everything is correct during preparation, so we can report errors early
	matchedMocks := map[reflect.Type]bool{}
	for i := 0; i < entrySubject.Elem().NumField(); i++ {
		subjectEntry := entrySubject.Elem().Field(i)
		subjectFieldName := entrySubject.Elem().Type().Field(i).Name

		if subjectEntry.Kind() != reflect.Interface {
			continue
		}

		var interfaceMatches []reflect.Type
		for mockType := range entry.mocks {
			if mockType.Implements(subjectEntry.Type()) {
				interfaceMatches = append(interfaceMatches, mockType)
				matchedMocks[mockType] = true
			}
		}

		if len(interfaceMatches) == 0 {
			entry.tableLevelWarnings = append(entry.tableLevelWarnings,
				fmt.Sprintf(
					"No mocks matched '%s', the interface for Subject.%s",
					subjectEntry.Type().String(),
					subjectFieldName,
				),
			)

			continue
		}

		if len(interfaceMatches) > 1 {
			matchedMockTypes := []string{}
			for _, interfaceMatch := range interfaceMatches {
				matchedMockTypes = append(matchedMockTypes, interfaceMatch.String())
			}
			sort.Strings(matchedMockTypes)

			return erk.WithParams(errSubjectFieldMatchingMultipleMocks, erk.Params{
				"subjectFieldName": subjectFieldName,
				"matchedMockTypes": strings.Join(matchedMockTypes, ", "),
			})
		}

		// At this point, everything should be correct, so we can blindly execute without worrying about types
		interfaceMatch := interfaceMatches[0]
		entry.setupFuncs = append(entry.setupFuncs, func(c *Chain) {
			subjectEntry.Set(entry.mocks[interfaceMatch].value)
		})
	}

	// Add warnings for unused mocks
	for mockType, mock := range entry.mocks {
		if _, ok := matchedMocks[mockType]; !ok && !mock.tag.isIgnoreUnused() {
			entry.tableLevelWarnings = append(entry.tableLevelWarnings,
				fmt.Sprintf(
					"Mocks.%s (type %s) did not match any interfaces in the Subject",
					mock.fieldName,
					mockType.String(),
				),
			)
		}
	}

	return nil
}

func (tag tableEntryMockTag) isIgnoreUnused() bool {
	return tag == "ignoreunused"
}
