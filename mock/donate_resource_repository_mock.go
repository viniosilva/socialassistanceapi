// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/socialassistanceapi/internal/repository (interfaces: DonateResourceRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDonateResourceRepository is a mock of DonateResourceRepository interface.
type MockDonateResourceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDonateResourceRepositoryMockRecorder
}

// MockDonateResourceRepositoryMockRecorder is the mock recorder for MockDonateResourceRepository.
type MockDonateResourceRepositoryMockRecorder struct {
	mock *MockDonateResourceRepository
}

// NewMockDonateResourceRepository creates a new mock instance.
func NewMockDonateResourceRepository(ctrl *gomock.Controller) *MockDonateResourceRepository {
	mock := &MockDonateResourceRepository{ctrl: ctrl}
	mock.recorder = &MockDonateResourceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDonateResourceRepository) EXPECT() *MockDonateResourceRepositoryMockRecorder {
	return m.recorder
}

// Donate mocks base method.
func (m *MockDonateResourceRepository) Donate(arg0 context.Context, arg1, arg2 int, arg3 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Donate", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Donate indicates an expected call of Donate.
func (mr *MockDonateResourceRepositoryMockRecorder) Donate(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Donate", reflect.TypeOf((*MockDonateResourceRepository)(nil).Donate), arg0, arg1, arg2, arg3)
}

// Return mocks base method.
func (m *MockDonateResourceRepository) Return(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Return", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Return indicates an expected call of Return.
func (mr *MockDonateResourceRepositoryMockRecorder) Return(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Return", reflect.TypeOf((*MockDonateResourceRepository)(nil).Return), arg0, arg1)
}
