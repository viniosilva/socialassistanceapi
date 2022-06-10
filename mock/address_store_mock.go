// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/store (interfaces: AddressStore)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/socialassistanceapi/internal/model"
)

// MockAddressStore is a mock of AddressStore interface.
type MockAddressStore struct {
	ctrl     *gomock.Controller
	recorder *MockAddressStoreMockRecorder
}

// MockAddressStoreMockRecorder is the mock recorder for MockAddressStore.
type MockAddressStoreMockRecorder struct {
	mock *MockAddressStore
}

// NewMockAddressStore creates a new mock instance.
func NewMockAddressStore(ctrl *gomock.Controller) *MockAddressStore {
	mock := &MockAddressStore{ctrl: ctrl}
	mock.recorder = &MockAddressStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAddressStore) EXPECT() *MockAddressStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAddressStore) Create(arg0 context.Context, arg1 model.Address) (*model.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAddressStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAddressStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockAddressStore) Delete(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAddressStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAddressStore)(nil).Delete), arg0, arg1)
}

// FindAll mocks base method.
func (m *MockAddressStore) FindAll(arg0 context.Context) ([]model.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].([]model.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockAddressStoreMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockAddressStore)(nil).FindAll), arg0)
}

// FindOneById mocks base method.
func (m *MockAddressStore) FindOneById(arg0 context.Context, arg1 int) (*model.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneById", arg0, arg1)
	ret0, _ := ret[0].(*model.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneById indicates an expected call of FindOneById.
func (mr *MockAddressStoreMockRecorder) FindOneById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneById", reflect.TypeOf((*MockAddressStore)(nil).FindOneById), arg0, arg1)
}

// Update mocks base method.
func (m *MockAddressStore) Update(arg0 context.Context, arg1 model.Address) (*model.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*model.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAddressStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAddressStore)(nil).Update), arg0, arg1)
}