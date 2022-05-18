// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/store (interfaces: CustomerStore)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/socialassistanceapi/internal/model"
)

// MockCustomerStore is a mock of CustomerStore interface.
type MockCustomerStore struct {
	ctrl     *gomock.Controller
	recorder *MockCustomerStoreMockRecorder
}

// MockCustomerStoreMockRecorder is the mock recorder for MockCustomerStore.
type MockCustomerStoreMockRecorder struct {
	mock *MockCustomerStore
}

// NewMockCustomerStore creates a new mock instance.
func NewMockCustomerStore(ctrl *gomock.Controller) *MockCustomerStore {
	mock := &MockCustomerStore{ctrl: ctrl}
	mock.recorder = &MockCustomerStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomerStore) EXPECT() *MockCustomerStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCustomerStore) Create(arg0 context.Context, arg1 model.Customer) (*model.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCustomerStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCustomerStore)(nil).Create), arg0, arg1)
}

// FindAll mocks base method.
func (m *MockCustomerStore) FindAll(arg0 context.Context) ([]model.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].([]model.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockCustomerStoreMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockCustomerStore)(nil).FindAll), arg0)
}

// FindOneById mocks base method.
func (m *MockCustomerStore) FindOneById(arg0 context.Context, arg1 int) (*model.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneById", arg0, arg1)
	ret0, _ := ret[0].(*model.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneById indicates an expected call of FindOneById.
func (mr *MockCustomerStoreMockRecorder) FindOneById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneById", reflect.TypeOf((*MockCustomerStore)(nil).FindOneById), arg0, arg1)
}

// Update mocks base method.
func (m *MockCustomerStore) Update(arg0 context.Context, arg1 model.Customer) (*model.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*model.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCustomerStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCustomerStore)(nil).Update), arg0, arg1)
}