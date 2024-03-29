// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/api (interfaces: PersonApi)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPersonApi is a mock of PersonApi interface.
type MockPersonApi struct {
	ctrl     *gomock.Controller
	recorder *MockPersonApiMockRecorder
}

// MockPersonApiMockRecorder is the mock recorder for MockPersonApi.
type MockPersonApiMockRecorder struct {
	mock *MockPersonApi
}

// NewMockPersonApi creates a new mock instance.
func NewMockPersonApi(ctrl *gomock.Controller) *MockPersonApi {
	mock := &MockPersonApi{ctrl: ctrl}
	mock.recorder = &MockPersonApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonApi) EXPECT() *MockPersonApiMockRecorder {
	return m.recorder
}

// Configure mocks base method.
func (m *MockPersonApi) Configure() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Configure")
}

// Configure indicates an expected call of Configure.
func (mr *MockPersonApiMockRecorder) Configure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockPersonApi)(nil).Configure))
}
