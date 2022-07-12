// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/api (interfaces: ResourceApi)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResourceApi is a mock of ResourceApi interface.
type MockResourceApi struct {
	ctrl     *gomock.Controller
	recorder *MockResourceApiMockRecorder
}

// MockResourceApiMockRecorder is the mock recorder for MockResourceApi.
type MockResourceApiMockRecorder struct {
	mock *MockResourceApi
}

// NewMockResourceApi creates a new mock instance.
func NewMockResourceApi(ctrl *gomock.Controller) *MockResourceApi {
	mock := &MockResourceApi{ctrl: ctrl}
	mock.recorder = &MockResourceApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceApi) EXPECT() *MockResourceApiMockRecorder {
	return m.recorder
}

// Configure mocks base method.
func (m *MockResourceApi) Configure() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Configure")
}

// Configure indicates an expected call of Configure.
func (mr *MockResourceApiMockRecorder) Configure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockResourceApi)(nil).Configure))
}
