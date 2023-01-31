// Code generated by `ensure mocks generate`. DO NOT EDIT.
// Source: github.com/JosiahWitt/ensure/internal/testctx (interfaces: T, Context)

// Package mock_testctx is a generated GoMock package.
package mock_testctx

import (
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

// MockT is a mock of the T interface in github.com/JosiahWitt/ensure/internal/testctx.
type MockT struct {
	ctrl     *gomock.Controller
	recorder *MockTMockRecorder
}

// MockTMockRecorder is the mock recorder for MockT.
type MockTMockRecorder struct {
	mock *MockT
}

// NewMockT creates a new mock instance.
func NewMockT(ctrl *gomock.Controller) *MockT {
	mock := &MockT{ctrl: ctrl}
	mock.recorder = &MockTMockRecorder{mock}
	return mock
}

// NEW creates a MockT. This method is used internally by ensure.
func (*MockT) NEW(ctrl *gomock.Controller) *MockT {
	return NewMockT(ctrl)
}

// EXPECT returns a struct that allows setting up expectations.
func (m *MockT) EXPECT() *MockTMockRecorder {
	return m.recorder
}

// Cleanup mocks Cleanup on T.
func (m *MockT) Cleanup(_arg0 func()) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_arg0}
	ret := m.ctrl.Call(m, "Cleanup", inputs...)
	var _ = ret // Unused, since there are no returns
	return
}

// Cleanup sets up expectations for calls to Cleanup.
// Calling this method multiple times allows expecting multiple calls to Cleanup with a variety of parameters.
//
// Inputs:
//
//	func()
//
// Outputs:
//
//	none
func (mr *MockTMockRecorder) Cleanup(_arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_arg0}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockT)(nil).Cleanup), inputs...)
}

// Errorf mocks Errorf on T.
func (m *MockT) Errorf(_format string, _args ...interface{}) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_format}
	for _, variadicInput := range _args {
		inputs = append(inputs, variadicInput)
	}
	ret := m.ctrl.Call(m, "Errorf", inputs...)
	var _ = ret // Unused, since there are no returns
	return
}

// Errorf sets up expectations for calls to Errorf.
// Calling this method multiple times allows expecting multiple calls to Errorf with a variety of parameters.
//
// Inputs:
//
//	format string
//	args ...interface{}
//
// Outputs:
//
//	none
func (mr *MockTMockRecorder) Errorf(_format interface{}, _args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_format}
	for _, variadicInput := range _args {
		inputs = append(inputs, variadicInput)
	}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockT)(nil).Errorf), inputs...)
}

// Fatalf mocks Fatalf on T.
func (m *MockT) Fatalf(_format string, _args ...interface{}) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_format}
	for _, variadicInput := range _args {
		inputs = append(inputs, variadicInput)
	}
	ret := m.ctrl.Call(m, "Fatalf", inputs...)
	var _ = ret // Unused, since there are no returns
	return
}

// Fatalf sets up expectations for calls to Fatalf.
// Calling this method multiple times allows expecting multiple calls to Fatalf with a variety of parameters.
//
// Inputs:
//
//	format string
//	args ...interface{}
//
// Outputs:
//
//	none
func (mr *MockTMockRecorder) Fatalf(_format interface{}, _args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_format}
	for _, variadicInput := range _args {
		inputs = append(inputs, variadicInput)
	}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatalf", reflect.TypeOf((*MockT)(nil).Fatalf), inputs...)
}

// Helper mocks Helper on T.
func (m *MockT) Helper() {
	m.ctrl.T.Helper()
	inputs := []interface{}{}
	ret := m.ctrl.Call(m, "Helper", inputs...)
	var _ = ret // Unused, since there are no returns
	return
}

// Helper sets up expectations for calls to Helper.
// Calling this method multiple times allows expecting multiple calls to Helper with a variety of parameters.
//
// Inputs:
//
//	none
//
// Outputs:
//
//	none
func (mr *MockTMockRecorder) Helper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper", reflect.TypeOf((*MockT)(nil).Helper), inputs...)
}

// Logf mocks Logf on T.
func (m *MockT) Logf(_format string, _args ...interface{}) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_format}
	for _, variadicInput := range _args {
		inputs = append(inputs, variadicInput)
	}
	ret := m.ctrl.Call(m, "Logf", inputs...)
	var _ = ret // Unused, since there are no returns
	return
}

// Logf sets up expectations for calls to Logf.
// Calling this method multiple times allows expecting multiple calls to Logf with a variety of parameters.
//
// Inputs:
//
//	format string
//	args ...interface{}
//
// Outputs:
//
//	none
func (mr *MockTMockRecorder) Logf(_format interface{}, _args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_format}
	for _, variadicInput := range _args {
		inputs = append(inputs, variadicInput)
	}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logf", reflect.TypeOf((*MockT)(nil).Logf), inputs...)
}

// Run mocks Run on T.
func (m *MockT) Run(_name string, _f func(t *testing.T)) bool {
	m.ctrl.T.Helper()
	inputs := []interface{}{_name, _f}
	ret := m.ctrl.Call(m, "Run", inputs...)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Run sets up expectations for calls to Run.
// Calling this method multiple times allows expecting multiple calls to Run with a variety of parameters.
//
// Inputs:
//
//	name string
//	f func(t *testing.T)
//
// Outputs:
//
//	bool
func (mr *MockTMockRecorder) Run(_name interface{}, _f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_name, _f}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockT)(nil).Run), inputs...)
}

// MockContext is a mock of the Context interface in github.com/JosiahWitt/ensure/internal/testctx.
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext.
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance.
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// NEW creates a MockContext. This method is used internally by ensure.
func (*MockContext) NEW(ctrl *gomock.Controller) *MockContext {
	return NewMockContext(ctrl)
}

// EXPECT returns a struct that allows setting up expectations.
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// GoMockController mocks GoMockController on Context.
func (m *MockContext) GoMockController() *gomock.Controller {
	m.ctrl.T.Helper()
	inputs := []interface{}{}
	ret := m.ctrl.Call(m, "GoMockController", inputs...)
	ret0, _ := ret[0].(*gomock.Controller)
	return ret0
}

// GoMockController sets up expectations for calls to GoMockController.
// Calling this method multiple times allows expecting multiple calls to GoMockController with a variety of parameters.
//
// Inputs:
//
//	none
//
// Outputs:
//
//	*gomock.Controller
func (mr *MockContextMockRecorder) GoMockController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoMockController", reflect.TypeOf((*MockContext)(nil).GoMockController), inputs...)
}

// Run mocks Run on Context.
func (m *MockContext) Run(_name string, _fn func(testctx.Context)) {
	m.ctrl.T.Helper()
	inputs := []interface{}{_name, _fn}
	ret := m.ctrl.Call(m, "Run", inputs...)
	var _ = ret // Unused, since there are no returns
	return
}

// Run sets up expectations for calls to Run.
// Calling this method multiple times allows expecting multiple calls to Run with a variety of parameters.
//
// Inputs:
//
//	name string
//	fn func(testctx.Context)
//
// Outputs:
//
//	none
func (mr *MockContextMockRecorder) Run(_name interface{}, _fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{_name, _fn}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockContext)(nil).Run), inputs...)
}

// T mocks T on Context.
func (m *MockContext) T() testctx.T {
	m.ctrl.T.Helper()
	inputs := []interface{}{}
	ret := m.ctrl.Call(m, "T", inputs...)
	ret0, _ := ret[0].(testctx.T)
	return ret0
}

// T sets up expectations for calls to T.
// Calling this method multiple times allows expecting multiple calls to T with a variety of parameters.
//
// Inputs:
//
//	none
//
// Outputs:
//
//	testctx.T
func (mr *MockContextMockRecorder) T() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	inputs := []interface{}{}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "T", reflect.TypeOf((*MockContext)(nil).T), inputs...)
}
