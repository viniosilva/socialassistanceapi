// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/store (interfaces: HealthStore)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHealthStore is a mock of HealthStore interface.
type MockHealthStore struct {
	ctrl     *gomock.Controller
	recorder *MockHealthStoreMockRecorder
}

// MockHealthStoreMockRecorder is the mock recorder for MockHealthStore.
type MockHealthStoreMockRecorder struct {
	mock *MockHealthStore
}

// NewMockHealthStore creates a new mock instance.
func NewMockHealthStore(ctrl *gomock.Controller) *MockHealthStore {
	mock := &MockHealthStore{ctrl: ctrl}
	mock.recorder = &MockHealthStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthStore) EXPECT() *MockHealthStoreMockRecorder {
	return m.recorder
}

// Health mocks base method.
func (m *MockHealthStore) Health(arg0 context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockHealthStoreMockRecorder) Health(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockHealthStore)(nil).Health), arg0)
}
