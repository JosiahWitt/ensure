// Package mocks provides a plugin that initializes Mocks in the test entry.
package mocks

import (
	"fmt"
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/id"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/iterate"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

//nolint:gochecknoglobals // This is only used internally for comparison.
var goMockControllerType = reflect.TypeOf(&gomock.Controller{})

// New uses the collection of mocks, initializing and populating it for each test entry.
func New(m *mocks.All) *TablePlugin {
	return &TablePlugin{mocks: m}
}

// TablePlugin uses the collection of mocks, initializing and populating it for each test entry.
// This plugin should come before any steps that require mocks.
type TablePlugin struct {
	mocks *mocks.All
}

var _ plugins.TablePlugin = &TablePlugin{}

// ParseEntryType is called during the first pass of plugin initialization.
// It is responsible for making sure the types are as expected.
func (t *TablePlugin) ParseEntryType(entryType reflect.Type) (plugins.TableEntryPlugin, error) {
	p := &TableEntryPlugin{}

	mocksStruct, ok := entryType.FieldByName(id.Mocks)
	if ok {
		if err := validateMocksFieldType(&mocksStruct); err != nil {
			return nil, err
		}

		mockFields, structFieldsResult, err := t.parseMocks(&mocksStruct)
		if err != nil {
			return nil, err
		}

		p.hasMocks = true
		p.mockFields = mockFields
		p.structFields = structFieldsResult
	}

	return p, nil
}

func validateMocksFieldType(mocksStruct *reflect.StructField) error {
	t := mocksStruct.Type

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return stringerr.Newf("expected %s field to be a pointer to a struct, got %s", id.Mocks, t)
	}

	return nil
}

//nolint:funlen // It seems clearer to keep this in a single method.
func (t *TablePlugin) parseMocks(mocksStruct *reflect.StructField) (map[string]*mockField, *iterate.StructFieldsResult, error) {
	mockFields := map[string]*mockField{}

	structFieldsResult, errs := iterate.StructFields(id.Mocks, mocksStruct.Type, func(mocksFieldPath string, mocksField *reflect.StructField) []error {
		tag, err := parseTag(&mocksField.Tag)
		if err != nil {
			return []error{stringerr.Newf("%s: %v", mocksFieldPath, err)}
		}

		if tag.ignore {
			return nil
		}

		newMethod, hasNew := mocksField.Type.MethodByName(id.NEW)
		if !hasNew {
			return []error{stringerr.Newf(
				"%s (%v) is missing a %s method, which is required for automatically initializing the mock. To ignore %[1]s, add the %[4]s tag.",
				mocksFieldPath,
				mocksField.Type,
				id.NEW,
				id.ExampleIgnore,
			)}
		}

		// NEW signature should be one of:
		//  func (m *MockXYZ) NEW(ctrl *gomock.Controller) *MockXYZ { ... }
		//  func (m *MockXYZ) NEW() *MockXYZ { ... }

		numIn := newMethod.Type.NumIn() - 1 // We subtract one, since it includes the method receiver, which we don't specifically care about
		needsGoMockController := numIn == 1
		invalidIns := numIn > 1 || (needsGoMockController && newMethod.Type.In(1) != goMockControllerType)

		numOut := newMethod.Type.NumOut()
		invalidOuts := numOut != 1 || newMethod.Type.Out(0) != mocksField.Type

		if invalidIns || invalidOuts {
			return []error{stringerr.NewBlock(
				fmt.Sprintf("%s (%v) must have a %s method matching one of the following signatures", mocksFieldPath, mocksField.Type, id.NEW),
				[]error{
					stringerr.Newf("func (m %[1]v) NEW() %[1]v { ... }", mocksField.Type),
					stringerr.Newf("func (m %[1]v) NEW(ctrl *gomock.Controller) %[1]v { ... }", mocksField.Type),
				},
				fmt.Sprintf("To ignore %s, add the %s tag.", mocksFieldPath, id.ExampleIgnore),
			)}
		}

		mock := t.mocks.AddMock(mocksFieldPath, tag.optional, mocksField.Type)

		mockFields[mocksFieldPath] = &mockField{
			mock: mock,

			needsGoMockController: needsGoMockController,
		}

		return nil
	})

	if len(errs) != 0 {
		return nil, nil, stringerr.NewGroup(fmt.Sprintf("Unable to build %s field", id.Mocks), errs)
	}

	return mockFields, structFieldsResult, nil
}

type tag struct {
	ignore   bool
	optional bool
}

func parseTag(structTag *reflect.StructTag) (*tag, error) {
	t, ok := structTag.Lookup(id.Ensure)
	if !ok {
		return &tag{}, nil
	}

	switch t {
	case id.IgnoreUnused:
		return &tag{optional: true}, nil
	case id.Ignore:
		return &tag{ignore: true}, nil
	default:
		return nil, stringerr.Newf("Only %s or %s tags are supported, got: `%s:\"%s\"`", id.ExampleIgnore, id.ExampleIgnoreUnused, id.Ensure, t)
	}
}

type shared struct {
	hasMocks     bool
	mockFields   map[string]*mockField
	structFields *iterate.StructFieldsResult
}

type mockField struct {
	mock *mocks.Mock

	needsGoMockController bool
}

// TableEntryPlugin is called for each entry in the table.
type TableEntryPlugin struct {
	shared
}

var _ plugins.TableEntryPlugin = &TableEntryPlugin{}

// ParseEntryValue parses the value associated with the entry.
func (p *TableEntryPlugin) ParseEntryValue(entryValue reflect.Value, i int) (plugins.TableEntryHooks, error) {
	return &TableEntryHooks{
		shared: p.shared,
		value:  entryValue.FieldByName(id.Mocks),
		index:  i,
	}, nil
}

// TableEntryHooks exposes the before and after hooks for each entry in the table.
type TableEntryHooks struct {
	shared
	value reflect.Value
	index int
}

var _ plugins.TableEntryHooks = &TableEntryHooks{}

// BeforeEntry is called before the test is run for the table entry.
// It initializes the Mocks struct and calls NEW for each of the mocks.
func (h *TableEntryHooks) BeforeEntry(t *testctx.Context) error {
	if !h.hasMocks {
		return nil
	}

	// This is only populated the first time it is needed
	var goMockController reflect.Value

	h.structFields.InitializeStruct(h.value, func(fieldPath string, field reflect.Value) {
		mockField, ok := h.mockFields[fieldPath]
		if !ok {
			return
		}

		var ins []reflect.Value
		if mockField.needsGoMockController {
			if !goMockController.IsValid() {
				goMockController = reflect.ValueOf(t.GoMockController())
			}

			ins = append(ins, goMockController)
		}

		newMethod := field.MethodByName(id.NEW)
		outs := newMethod.Call(ins)
		mock := outs[0]

		field.Set(mock)
		mockField.mock.SetValueByEntryIndex(h.index, mock)
	})

	return nil
}

// AfterEntry is called after the test is run for the table entry.
func (*TableEntryHooks) AfterEntry(*testctx.Context) error { return nil }
