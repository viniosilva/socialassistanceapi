// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/repository (interfaces: FamilyRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/socialassistanceapi/internal/model"
)

// MockFamilyRepository is a mock of FamilyRepository interface.
type MockFamilyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFamilyRepositoryMockRecorder
}

// MockFamilyRepositoryMockRecorder is the mock recorder for MockFamilyRepository.
type MockFamilyRepositoryMockRecorder struct {
	mock *MockFamilyRepository
}

// NewMockFamilyRepository creates a new mock instance.
func NewMockFamilyRepository(ctrl *gomock.Controller) *MockFamilyRepository {
	mock := &MockFamilyRepository{ctrl: ctrl}
	mock.recorder = &MockFamilyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFamilyRepository) EXPECT() *MockFamilyRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockFamilyRepository) Create(arg0 context.Context, arg1 model.Family) (*model.Family, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Family)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockFamilyRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFamilyRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockFamilyRepository) Delete(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockFamilyRepositoryMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFamilyRepository)(nil).Delete), arg0, arg1)
}

// FindAll mocks base method.
func (m *MockFamilyRepository) FindAll(arg0 context.Context) ([]model.Family, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].([]model.Family)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockFamilyRepositoryMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockFamilyRepository)(nil).FindAll), arg0)
}

// FindOneById mocks base method.
func (m *MockFamilyRepository) FindOneById(arg0 context.Context, arg1 int) (*model.Family, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneById", arg0, arg1)
	ret0, _ := ret[0].(*model.Family)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneById indicates an expected call of FindOneById.
func (mr *MockFamilyRepositoryMockRecorder) FindOneById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneById", reflect.TypeOf((*MockFamilyRepository)(nil).FindOneById), arg0, arg1)
}

// Update mocks base method.
func (m *MockFamilyRepository) Update(arg0 context.Context, arg1 model.Family) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockFamilyRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockFamilyRepository)(nil).Update), arg0, arg1)
}
