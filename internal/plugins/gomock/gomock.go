package gomock

import (
	"reflect"

	"github.com/JosiahWitt/erk"
)

const (
	MocksFieldName      = "Mocks"
	SetupMocksFieldName = "SetupMocks"
)

type (
	erkTableInvalid struct{ erk.DefaultKind }
)

var (
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
)

type EnsurePluginGoMock struct{}

func New() *EnsurePluginGoMock {
	return &EnsurePluginGoMock{}
}

func (p *EnsurePluginGoMock) ForTableEntryType(tableEntry reflect.Type) (*TableBuilder, error) {
	hasMocksField, err := validateMocksField(tableEntry)
	if err != nil {
		return nil, err
	}

	setupMocksSignature, err := validateSetupMocksField(tableEntry)
	if err != nil {
		return nil, err
	}

	return &TableBuilder{
		mocksField:          hasMocksField,
		setupMocksSignature: setupMocksSignature,
	}, nil
}
