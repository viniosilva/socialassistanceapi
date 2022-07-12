package service_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/mock"
)

func TestAddressServiceFindAll(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		expectedRes service.AddressesResponse
		expectedErr error
		prepareMock func(mock *mock.MockAddressStore)
	}{
		"should return addresses list": {
			expectedRes: service.AddressesResponse{Data: []service.Address{{
				ID:           1,
				CreatedAt:    DATE,
				UpdatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}}},
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Address{{
					ID:           1,
					CreatedAt:    DATETIME,
					UpdatedAt:    DATETIME,
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}}, nil)
			},
		},
		"should return empty addresses list": {
			expectedRes: service.AddressesResponse{Data: []service.Address{}},
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Address{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockAddressStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewAddressService(storeMock)

			// when
			res, err := impl.FindAll(ctx)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("AddressService.FindAll() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("AddressService.FindAll() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestAddressServiceFindOneByID(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputAddressID int
		expectedRes    service.AddressResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockAddressStore)
	}{
		"should return address when exists": {
			inputAddressID: 1,
			expectedRes: service.AddressResponse{Data: &service.Address{
				ID:           1,
				CreatedAt:    DATE,
				UpdatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(&model.Address{
					ID:           1,
					CreatedAt:    DATETIME,
					UpdatedAt:    DATETIME,
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}, nil)
			},
		},
		"should return empty when address not exists": {
			inputAddressID: 1,
			expectedErr:    exception.NewNotFoundException("resource"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).
					Return(nil, exception.NewNotFoundException("resource"))
			},
		},
		"should throw error": {
			inputAddressID: 1,
			expectedErr:    fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockAddressStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewAddressService(storeMock)

			// when
			res, err := impl.FindOneById(ctx, cs.inputAddressID)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("AddressService.FindOneById() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("AddressService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestAddressServiceCreate(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputAddress service.AddressDto
		expectedRes  service.AddressResponse
		expectedErr  error
		prepareMock  func(mock *mock.MockAddressStore)
	}{
		"should create address": {
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedRes: service.AddressResponse{Data: &service.Address{
				ID:           1,
				CreatedAt:    DATE,
				UpdatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Address{
					ID:           1,
					CreatedAt:    DATETIME,
					UpdatedAt:    DATETIME,
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}, nil)
			},
		},
		"should throw error": {
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockAddressStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewAddressService(storeMock)

			// when
			res, err := impl.Create(ctx, cs.inputAddress)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("AddressService.FindOneById() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("AddressService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestAddressServiceUpdate(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputAddressID int
		inputAddress   service.AddressDto
		expectedRes    service.AddressResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockAddressStore)
	}{
		"should update address": {
			inputAddressID: 1,
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			},
			expectedRes: service.AddressResponse{Data: &service.Address{
				ID:           1,
				CreatedAt:    DATE,
				UpdatedAt:    DATE,
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			}},
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&model.Address{
					ID:           1,
					CreatedAt:    DATETIME,
					UpdatedAt:    DATETIME,
					Country:      "BR",
					State:        "RS",
					City:         "Porto Alegre",
					Neighborhood: "Hípica",
					Street:       "R. J",
					Number:       "1",
					Zipcode:      "91755450",
				}, nil)
			},
		},
		"should return empty when address not exists": {
			inputAddressID: 1,
			expectedErr:    exception.NewNotFoundException("resource"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(nil, exception.NewNotFoundException("resource"))
			},
		},
		"should throw error": {
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockAddressStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewAddressService(storeMock)

			// when
			res, err := impl.Update(ctx, cs.inputAddressID, cs.inputAddress)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("AddressService.Update() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("AddressService.Update() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestAddressServiceDelete(t *testing.T) {
	cases := map[string]struct {
		inputAddressID int
		expectedErr    error
		prepareMock    func(mock *mock.MockAddressStore)
	}{
		"should delete address": {
			inputAddressID: 1,
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		"should throw error": {
			inputAddressID: 1,
			expectedErr:    fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockAddressStore) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockAddressStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewAddressService(storeMock)

			// when
			err := impl.Delete(ctx, cs.inputAddressID)

			// then
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("AddressService.Delete() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}
