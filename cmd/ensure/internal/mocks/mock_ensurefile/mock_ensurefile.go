// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile (interfaces: LoaderIface)

// Package mock_ensurefile is a generated GoMock package.
package mock_ensurefile

import (
	ensurefile "github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLoaderIface is a mock of LoaderIface interface
type MockLoaderIface struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderIfaceMockRecorder
}

// MockLoaderIfaceMockRecorder is the mock recorder for MockLoaderIface
type MockLoaderIfaceMockRecorder struct {
	mock *MockLoaderIface
}

// NewMockLoaderIface creates a new mock instance
func NewMockLoaderIface(ctrl *gomock.Controller) *MockLoaderIface {
	mock := &MockLoaderIface{ctrl: ctrl}
	mock.recorder = &MockLoaderIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLoaderIface) EXPECT() *MockLoaderIfaceMockRecorder {
	return m.recorder
}

// LoadConfig mocks base method
func (m *MockLoaderIface) LoadConfig(arg0 string) (*ensurefile.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfig", arg0)
	ret0, _ := ret[0].(*ensurefile.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadConfig indicates an expected call of LoadConfig
func (mr *MockLoaderIfaceMockRecorder) LoadConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockLoaderIface)(nil).LoadConfig), arg0)
}

// NEW creates a MockLoaderIface.
func (*MockLoaderIface) NEW(ctrl *gomock.Controller) *MockLoaderIface {
	return NewMockLoaderIface(ctrl)
}