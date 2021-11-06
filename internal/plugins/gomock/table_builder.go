package gomock

import "reflect"

type TableBuilder struct {
	mocksField          *mocksFieldConfig
	setupMocksSignature setupMocksSignature
}

type mocksFieldConfig struct {
	tags map[string]tableEntryMockTag
}

type tableEntryMockTag string

type setupMocksSignature int

const (
	setupMocksSignatureMissing setupMocksSignature = iota
	setupMocksSignatureOnlyMocks
)

func (b *TableBuilder) BeforeTableEntry(tableEntry reflect.Value) {
	if b.mocksField != nil {
		buildMocksField(tableEntry)
	}
}

func buildMocksField(tableEntry reflect.Value) {
	mocks := tableEntry.FieldByName(MocksFieldName)
	mocks.Set(reflect.New(mocks.Type().Elem()))
}
