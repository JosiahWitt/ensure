package gomock

import (
	"reflect"

	"github.com/JosiahWitt/erk"
	"github.com/golang/mock/gomock"
)

const (
	EnsureTagName     = "ensure"
	NewMockMethodName = "NEW"
)

func validateMocksField(tableEntry reflect.Type) (*mocksFieldConfig, error) {
	mocks, ok := tableEntry.FieldByName(MocksFieldName)
	if !ok {
		return nil, nil
	}

	if mocks.Type.Kind() != reflect.Ptr || mocks.Type.Elem().Kind() != reflect.Struct {
		return nil, erk.WithParams(errMocksNotStructPointer, erk.Params{"mocksField": mocks.Type.String()})
	}

	tags, err := validateMocksFieldStruct(mocks.Type.Elem())
	if err != nil {
		return nil, err
	}

	return &mocksFieldConfig{
		tags: tags,
	}, nil
}

func validateMocksFieldStruct(mocksField reflect.Type) (map[string]tableEntryMockTag, error) {
	tags := make(map[string]tableEntryMockTag)

	for i := 0; i < mocksField.NumField(); i++ {
		mockEntry := mocksField.Field(i)
		mockName := mockEntry.Name

		// Skip unexported fields
		if mockEntry.PkgPath != "" {
			continue
		}

		// Support embedded structs
		if mockEntry.Anonymous {
			if mockEntry.Type.Kind() != reflect.Struct {
				return nil, erk.WithParams(errMocksEmbeddedNotStruct, erk.Params{
					"mocksFieldName": mockName,
					"mockEntry":      mockEntry.Type.String(),
				})
			}

			tags, err := validateMocksFieldStruct(mockEntry.Type)
			if err != nil {
				return nil, err
			}

			for k, v := range tags {
				tags[k] = v
			}

			continue
		}

		if err := validateMocksFieldEntry(mockName, mockEntry.Type); err != nil {
			return nil, err
		}

		tags[mockName] = tableEntryMockTag(mockEntry.Tag.Get(EnsureTagName))
	}

	return tags, nil
}

// TODO: Return configuration about the NEW method
func validateMocksFieldEntry(mockFieldName string, mockEntry reflect.Type) error {
	if mockEntry.Kind() != reflect.Ptr {
		return erk.WithParams(errMocksEntryNotStructPointer, erk.Params{
			"mocksFieldName": mockFieldName,
			"mockEntry":      mockEntry.String(),
		})
	}

	// Mocks should have a NEW method, to allow creating the mock
	newMethod, ok := mockEntry.MethodByName(NewMockMethodName)
	if ok {
		return erk.WithParams(errMocksNEWMissing, erk.Params{
			"mocksFieldName": mockFieldName,
			"expectedReturn": mockEntry.String(),
		})
	}

	// NEW signature should be one of:
	//  func (m *MockXYZ) NEW(ctrl *gomock.Controller) *MockXYZ { ... }
	//  func (m *MockXYZ) NEW() *MockXYZ { ... }
	newMethodType := newMethod.Type
	controllerType := reflect.TypeOf(&gomock.Controller{})
	isInvalidParam := newMethodType.NumIn() > 1 || (newMethodType.NumIn() == 1 && newMethodType.In(0) != controllerType)
	isInvalidReturn := newMethodType.NumOut() != 1 || newMethodType.Out(0) != mockEntry
	if isInvalidParam || isInvalidReturn {
		return erk.WithParams(errMocksNEWInvalidSignature, erk.Params{
			"mocksFieldName": mockFieldName,
			"actualMethod":   newMethodType.String(),
			"expectedReturn": mockEntry.String(),
		})
	}

	return nil
}

func validateSetupMocksField(tableEntry reflect.Type) (setupMocksSignature, error) {
	setupMocks, ok := tableEntry.FieldByName(SetupMocksFieldName)
	if !ok {
		return setupMocksSignatureMissing, nil
	}

	mocks, ok := tableEntry.FieldByName(MocksFieldName)
	if !ok {
		return setupMocksSignatureMissing, errSetupMocksWithoutMocks
	}

	if setupMocks.Type.NumIn() == 1 && setupMocks.Type.In(0) == mocks.Type && setupMocks.Type.NumOut() == 0 {
		return setupMocksSignatureOnlyMocks, nil
	}

	return setupMocksSignatureMissing, erk.WithParams(errSetupMocksInvalidSignature, erk.Params{
		"expectedMockParam": mocks.Type.String(),
		"actualSetupMocks":  setupMocks.Type.String(),
	})
}
