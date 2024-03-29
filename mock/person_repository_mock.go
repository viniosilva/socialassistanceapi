// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/repository (interfaces: PersonRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/socialassistanceapi/internal/model"
)

// MockPersonRepository is a mock of PersonRepository interface.
type MockPersonRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPersonRepositoryMockRecorder
}

// MockPersonRepositoryMockRecorder is the mock recorder for MockPersonRepository.
type MockPersonRepositoryMockRecorder struct {
	mock *MockPersonRepository
}

// NewMockPersonRepository creates a new mock instance.
func NewMockPersonRepository(ctrl *gomock.Controller) *MockPersonRepository {
	mock := &MockPersonRepository{ctrl: ctrl}
	mock.recorder = &MockPersonRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonRepository) EXPECT() *MockPersonRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPersonRepository) Create(arg0 context.Context, arg1 model.Person) (*model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPersonRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPersonRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockPersonRepository) Delete(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPersonRepositoryMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPersonRepository)(nil).Delete), arg0, arg1)
}

// FindAll mocks base method.
func (m *MockPersonRepository) FindAll(arg0 context.Context) ([]model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].([]model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockPersonRepositoryMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockPersonRepository)(nil).FindAll), arg0)
}

// FindOneById mocks base method.
func (m *MockPersonRepository) FindOneById(arg0 context.Context, arg1 int) (*model.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneById", arg0, arg1)
	ret0, _ := ret[0].(*model.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneById indicates an expected call of FindOneById.
func (mr *MockPersonRepositoryMockRecorder) FindOneById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneById", reflect.TypeOf((*MockPersonRepository)(nil).FindOneById), arg0, arg1)
}

// Update mocks base method.
func (m *MockPersonRepository) Update(arg0 context.Context, arg1 model.Person) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPersonRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPersonRepository)(nil).Update), arg0, arg1)
}
