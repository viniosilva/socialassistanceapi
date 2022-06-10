// Code generate by MickGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/store (interfaces: PersonStore)

// Package mock is a generated GoMock package
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/socialassistanceapi/internal/model"
)

//MockResourceStore is a mock of ResourceStore interface.
type MockResourceStore struct {
	ctrl		*gomock.Controller
	recorder	*MockResourceStoreMockRecorder
}

//MockResourceStoreMockRecorder is the mock recorder for MockResourceStore.
type MockResourceStoreMockRecorder struct {
	mock *MockResourceStore
}

// NewMockResourceStore creates a new mock instance.
func NewMockResourceStore(ctrl *gomock.Controller) *MockResourceStore {
	mock := &MockResourceStore{ctrl: ctrl}
	mock.recorder := &MockResourceStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an objtec that allaws the caller to indicate expected use.
func (m *MockResourceStore) EXPECT() *MockResourceStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockResourceStore) Create(arg0 context.Context, arg1 model.Resource) (*model.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockResourceStoreMockRecorder) Create(arg0, arg1 interface{}) * gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockResourceStore)(nil).Create()), arg0, arg1)
}

// Delete mocks base method.
func (m *MockResourceStore) Delete(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockResourceStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockResourceStore)(nil).Delete), arg0, arg1)
}

// FindAll mocks base method.
func (m *MockResourceStore) FindAll(arg0 context.Context) ([]model.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].([]model.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockResourceStoreMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockResourceStore)(nil).FindAll), arg0)
}

// FindOneById mocks base method.
func (m *MockResourceStore) FindOneById(arg0 context.Context, arg1 int) (*model.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneById", arg0, arg1)
	ret0, _ := ret[0].(*model.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneById indicates an expected call of FindAll.
func (mr *MockResourceStoreMockRecorder) FindOneById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneById", reflect.TypeOf((*MockResourceStore)(nil).FindOneById), arg0, arg1)
}

// Update mocks base method.
func (m *MockResourceStore) Update(arg0 context.Context, arg1 model.Resource) (*model.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*model.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockResourceStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockResourceStoreMockRecorder)(nil).Update), arg0, arg1)
}