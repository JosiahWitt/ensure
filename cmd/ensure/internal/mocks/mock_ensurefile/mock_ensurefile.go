// Code generated by `ensure mocks generate`. DO NOT EDIT.
// Source: github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile (interfaces: LoaderIface)

// Package mock_ensurefile is a generated GoMock package.
package mock_ensurefile

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockLoaderIface is a mock of the LoaderIface interface in github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile.
type MockLoaderIface struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderIfaceMockRecorder
}

// MockLoaderIfaceMockRecorder is the mock recorder for MockLoaderIface.
type MockLoaderIfaceMockRecorder struct {
	mock *MockLoaderIface
}

// NewMockLoaderIface creates a new mock instance.
func NewMockLoaderIface(ctrl *gomock.Controller) *MockLoaderIface {
	mock := &MockLoaderIface{ctrl: ctrl}
	mock.recorder = &MockLoaderIfaceMockRecorder{mock}
	return mock
}

// NEW creates a MockLoaderIface. This method is used internally by ensure.
func (*MockLoaderIface) NEW(ctrl *gomock.Controller) *MockLoaderIface {
	return NewMockLoaderIface(ctrl)
}

// EXPECT returns a struct that allows setting up expectations.
func (m *MockLoaderIface) EXPECT() *MockLoaderIfaceMockRecorder {
	return m.recorder
}

// LoadConfig mocks LoadConfig on LoaderIface.
func (m *MockLoaderIface) LoadConfig(_pwd string) (*ensurefile.Config, error) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_pwd}
	ret := m.ctrl.Call(m, "LoadConfig", inputs...)
	ret0, _ := ret[0].(*ensurefile.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadConfig sets up expectations for calls to LoadConfig.
// Calling this method multiple times allows expecting multiple calls to LoadConfig with a variety of parameters.
//
// Inputs:
//
//  pwd string
//
// Outputs:
//
//  *ensurefile.Config
//  error
func (mr *MockLoaderIfaceMockRecorder) LoadConfig(_pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_pwd}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockLoaderIface)(nil).LoadConfig), inputs...)
}
