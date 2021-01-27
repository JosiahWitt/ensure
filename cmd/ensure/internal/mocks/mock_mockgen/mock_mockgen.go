// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen (interfaces: GeneratorIface)

// Package mock_mockgen is a generated GoMock package.
package mock_mockgen

import (
	ensurefile "github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGeneratorIface is a mock of GeneratorIface interface
type MockGeneratorIface struct {
	ctrl     *gomock.Controller
	recorder *MockGeneratorIfaceMockRecorder
}

// MockGeneratorIfaceMockRecorder is the mock recorder for MockGeneratorIface
type MockGeneratorIfaceMockRecorder struct {
	mock *MockGeneratorIface
}

// NewMockGeneratorIface creates a new mock instance
func NewMockGeneratorIface(ctrl *gomock.Controller) *MockGeneratorIface {
	mock := &MockGeneratorIface{ctrl: ctrl}
	mock.recorder = &MockGeneratorIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGeneratorIface) EXPECT() *MockGeneratorIfaceMockRecorder {
	return m.recorder
}

// GenerateMocks mocks base method
func (m *MockGeneratorIface) GenerateMocks(arg0 *ensurefile.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateMocks", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateMocks indicates an expected call of GenerateMocks
func (mr *MockGeneratorIfaceMockRecorder) GenerateMocks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateMocks", reflect.TypeOf((*MockGeneratorIface)(nil).GenerateMocks), arg0)
}

// NEW creates a MockGeneratorIface.
func (*MockGeneratorIface) NEW(ctrl *gomock.Controller) *MockGeneratorIface {
	return NewMockGeneratorIface(ctrl)
}
