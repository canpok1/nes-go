// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/domain/pad.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	gomock "github.com/golang/mock/gomock"
	domain "nes-go/pkg/domain"
	reflect "reflect"
)

// MockPad is a mock of Pad interface
type MockPad struct {
	ctrl     *gomock.Controller
	recorder *MockPadMockRecorder
}

// MockPadMockRecorder is the mock recorder for MockPad
type MockPadMockRecorder struct {
	mock *MockPad
}

// NewMockPad creates a new mock instance
func NewMockPad(ctrl *gomock.Controller) *MockPad {
	mock := &MockPad{ctrl: ctrl}
	mock.recorder = &MockPadMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPad) EXPECT() *MockPadMockRecorder {
	return m.recorder
}

// Load mocks base method
func (m *MockPad) Load() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load")
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load
func (mr *MockPadMockRecorder) Load() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockPad)(nil).Load))
}

// IsPressed mocks base method
func (m *MockPad) IsPressed(arg0 domain.ButtonType) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPressed", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPressed indicates an expected call of IsPressed
func (mr *MockPadMockRecorder) IsPressed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPressed", reflect.TypeOf((*MockPad)(nil).IsPressed), arg0)
}
