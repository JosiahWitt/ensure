// Package testhelper provides helpers for testing the plugins.
package testhelper

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
)

// MockData contains data required to build a mock.
type MockData struct {
	Path     string
	Optional bool

	Mock   interface{}
	Values []interface{}
}

// BuildMocks builds all the provided mocks into [mocks.All].
func BuildMocks(mockData []*MockData) *mocks.All {
	m := &mocks.All{}

	for _, md := range mockData {
		m.AddMock(md.Path, md.Optional, reflect.TypeOf(md.Mock))
	}

	for mockIdx, mock := range m.Slice() {
		for valIdx, val := range mockData[mockIdx].Values {
			mock.SetValueByEntryIndex(valIdx, reflect.ValueOf(val))
		}
	}

	return m
}
