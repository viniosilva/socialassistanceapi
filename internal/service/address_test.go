package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		prepareMock func(mockAddressRepository *mock.MockAddressRepository)
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
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().FindAll(gomock.Any()).Return([]model.Address{{
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
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().FindAll(gomock.Any()).Return([]model.Address{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockAddressRepository := mock.NewMockAddressRepository(ctrl)
			cs.prepareMock(mockAddressRepository)

			impl := &service.AddressServiceImpl{AddressRepository: mockAddressRepository}

			// when
			res, err := impl.FindAll(ctx)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
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
		prepareMock    func(mockAddressRepository *mock.MockAddressRepository)
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
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(&model.Address{
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
			expectedErr:    &exception.NotFoundException{Err: fmt.Errorf("address 1 not found")},
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().FindOneById(gomock.Any(), 1).
					Return(nil, &exception.NotFoundException{Err: fmt.Errorf("address 1 not found")})
			},
		},
		"should throw error": {
			inputAddressID: 1,
			expectedErr:    fmt.Errorf("error"),
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockAddressRepository := mock.NewMockAddressRepository(ctrl)
			cs.prepareMock(mockAddressRepository)

			impl := &service.AddressServiceImpl{AddressRepository: mockAddressRepository}

			// when
			res, err := impl.FindOneById(ctx, cs.inputAddressID)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestAddressServiceCreate(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputDto    service.AddressCreateDto
		expectedRes service.AddressResponse
		expectedErr error
		prepareMock func(mockAddressRepository *mock.MockAddressRepository)
	}{
		"should create address": {
			inputDto: service.AddressCreateDto{
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
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Create(gomock.Any(), model.Address{
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}).Return(&model.Address{
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
			inputDto: service.AddressCreateDto{
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
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Create(gomock.Any(), model.Address{
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockAddressRepository := mock.NewMockAddressRepository(ctrl)
			cs.prepareMock(mockAddressRepository)

			impl := &service.AddressServiceImpl{AddressRepository: mockAddressRepository}

			// when
			res, err := impl.Create(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestAddressServiceUpdate(t *testing.T) {
	cases := map[string]struct {
		inputDto    service.AddressUpdateDto
		expectedErr error
		prepareMock func(mockAddressRepository *mock.MockAddressRepository)
	}{
		"should update address": {
			inputDto: service.AddressUpdateDto{
				ID:           1,
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			},
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Update(gomock.Any(), model.Address{
					ID:           1,
					Country:      "BR",
					State:        "RS",
					City:         "Porto Alegre",
					Neighborhood: "Hípica",
					Street:       "R. J",
					Number:       "1",
					Zipcode:      "91755450",
				}).Return(nil)
			},
		},
		"should return empty when address not exists": {
			inputDto:    service.AddressUpdateDto{ID: 1},
			expectedErr: &exception.NotFoundException{Err: fmt.Errorf("address 1 not found")},
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Update(gomock.Any(), model.Address{ID: 1}).
					Return(&exception.NotFoundException{Err: fmt.Errorf("address 1 not found")})
			},
		},
		"should throw error": {
			inputDto: service.AddressUpdateDto{
				ID:           1,
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Update(gomock.Any(), model.Address{
					ID:           1,
					Country:      "BR",
					State:        "RS",
					City:         "Porto Alegre",
					Neighborhood: "Hípica",
					Street:       "R. J",
					Number:       "1",
					Zipcode:      "91755450",
				}).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockAddressRepository := mock.NewMockAddressRepository(ctrl)
			cs.prepareMock(mockAddressRepository)

			impl := &service.AddressServiceImpl{AddressRepository: mockAddressRepository}

			// when
			err := impl.Update(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestAddressServiceDelete(t *testing.T) {
	cases := map[string]struct {
		inputAddressID int
		expectedErr    error
		prepareMock    func(mockAddressRepository *mock.MockAddressRepository)
	}{
		"should delete address": {
			inputAddressID: 1,
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
		},
		"should throw error": {
			inputAddressID: 1,
			expectedErr:    fmt.Errorf("error"),
			prepareMock: func(mockAddressRepository *mock.MockAddressRepository) {
				mockAddressRepository.EXPECT().Delete(gomock.Any(), 1).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockAddressRepository := mock.NewMockAddressRepository(ctrl)
			cs.prepareMock(mockAddressRepository)

			impl := &service.AddressServiceImpl{AddressRepository: mockAddressRepository}

			// when
			err := impl.Delete(ctx, cs.inputAddressID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
