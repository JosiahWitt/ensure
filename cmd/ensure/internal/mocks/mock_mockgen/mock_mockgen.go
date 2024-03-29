// Code generated by `ensure mocks generate`. DO NOT EDIT.
// Source: github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen (interfaces: Generator)

// Package mock_mockgen is a generated GoMock package.
package mock_mockgen

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockGenerator is a mock of the Generator interface in github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen.
type MockGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockGeneratorMockRecorder
}

// MockGeneratorMockRecorder is the mock recorder for MockGenerator.
type MockGeneratorMockRecorder struct {
	mock *MockGenerator
}

// NewMockGenerator creates a new mock instance.
func NewMockGenerator(ctrl *gomock.Controller) *MockGenerator {
	mock := &MockGenerator{ctrl: ctrl}
	mock.recorder = &MockGeneratorMockRecorder{mock}
	return mock
}

// NEW creates a MockGenerator. This method is used internally by ensure.
func (*MockGenerator) NEW(ctrl *gomock.Controller) *MockGenerator {
	return NewMockGenerator(ctrl)
}

// EXPECT returns a struct that allows setting up expectations.
func (m *MockGenerator) EXPECT() *MockGeneratorMockRecorder {
	return m.recorder
}

// GenerateMocks mocks GenerateMocks on Generator.
func (m *MockGenerator) GenerateMocks(_pkgs []*ifacereader.Package, _imports *uniqpkg.UniquePackagePaths) ([]*mockgen.PackageMock, error) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_pkgs, _imports}
	ret := m.ctrl.Call(m, "GenerateMocks", inputs...)
	ret0, _ := ret[0].([]*mockgen.PackageMock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateMocks sets up expectations for calls to GenerateMocks.
// Calling this method multiple times allows expecting multiple calls to GenerateMocks with a variety of parameters.
//
// Inputs:
//
//	pkgs []*ifacereader.Package
//	imports *uniqpkg.UniquePackagePaths
//
// Outputs:
//
//	[]*mockgen.PackageMock
//	error
func (mr *MockGeneratorMockRecorder) GenerateMocks(_pkgs interface{}, _imports interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_pkgs, _imports}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateMocks", reflect.TypeOf((*MockGenerator)(nil).GenerateMocks), inputs...)
}
