package unit

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/mock"
)

func TestCustomerServiceFindAll(t *testing.T) {
	cases := map[string]struct {
		expectedCustomers service.CustomersResponse
		expectedErr       error
		prepareMock       func(mock *mock.MockCustomerStore)
	}{
		"should return customer list": {
			expectedCustomers: service.CustomersResponse{Data: []service.Customer{{ID: 1, Name: "Test"}}},
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Customer{{ID: 1, Name: "Test"}}, nil)
			},
		},
		"should return empty customer list": {
			expectedCustomers: service.CustomersResponse{Data: []service.Customer{}},
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Customer{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockCustomerStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewCustomerService(storeMock)

			// when
			customers, err := impl.FindAll(ctx)

			// then
			if !reflect.DeepEqual(customers, cs.expectedCustomers) {
				t.Errorf("CustomerService.FindAll() = %v, expected %v", customers, cs.expectedCustomers)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("CustomerService.FindAll() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestCustomerServiceFindOneByID(t *testing.T) {
	cases := map[string]struct {
		inputCustomerID  int
		expectedCustomer service.CustomerResponse
		expectedErr      error
		prepareMock      func(mock *mock.MockCustomerStore)
	}{
		"should return customer when exists": {
			inputCustomerID:  1,
			expectedCustomer: service.CustomerResponse{Data: &service.Customer{ID: 1, Name: "Test"}},
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(&model.Customer{ID: 1, Name: "Test"}, nil)
			},
		},
		"should return empty when customer not exists": {
			inputCustomerID:  1,
			expectedCustomer: service.CustomerResponse{},
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		"should throw error": {
			inputCustomerID: 1,
			expectedErr:     fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockCustomerStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewCustomerService(storeMock)

			// when
			customer, err := impl.FindOneById(ctx, cs.inputCustomerID)

			// then
			if !reflect.DeepEqual(customer, cs.expectedCustomer) {
				t.Errorf("CustomerService.FindOneById() = %v, expected %v", customer, cs.expectedCustomer)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("CustomerService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestCustomerServiceCreate(t *testing.T) {
	cases := map[string]struct {
		inputCustomer    service.CreateCustomerDto
		expectedCustomer service.CustomerResponse
		expectedErr      error
		prepareMock      func(mock *mock.MockCustomerStore)
	}{
		"should create customer": {
			inputCustomer:    service.CreateCustomerDto{Name: "Test"},
			expectedCustomer: service.CustomerResponse{Data: &service.Customer{ID: 1, Name: "Test"}},
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Customer{ID: 1, Name: "Test"}, nil)
			},
		},
		"should throw error": {
			inputCustomer: service.CreateCustomerDto{Name: "Test"},
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockCustomerStore) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockCustomerStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewCustomerService(storeMock)

			// when
			customer, err := impl.Create(ctx, cs.inputCustomer)

			// then
			if !reflect.DeepEqual(customer, cs.expectedCustomer) {
				t.Errorf("CustomerService.FindOneById() = %v, expected %v", customer, cs.expectedCustomer)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("CustomerService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}
