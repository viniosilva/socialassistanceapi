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

func Test_ResourceService_FindAll(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		expectedRes service.ResourcesResponse
		expectedErr error
		prepareMock func(mockResourceRepository *mock.MockResourceRepository)
	}{
		"should return resource list": {
			expectedRes: service.ResourcesResponse{Data: []service.Resource{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().FindAll(gomock.Any()).Return([]model.Resource{{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}}, nil)
			},
		},
		"should return empty resource list": {
			expectedRes: service.ResourcesResponse{Data: []service.Resource{}},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().FindAll(gomock.Any()).Return([]model.Resource{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockResourceRepository := mock.NewMockResourceRepository(ctrl)
			cs.prepareMock(mockResourceRepository)

			impl := &service.ResourceServiceImpl{ResourceRepository: mockResourceRepository}

			// when
			res, err := impl.FindAll(ctx)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_ResourceService_FindOneByID(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputResourceID int
		expectedRes     service.ResourceResponse
		expectedErr     error
		prepareMock     func(mockResourceRepository *mock.MockResourceRepository)
	}{
		"should return resource when exists": {
			inputResourceID: 1,
			expectedRes:     service.ResourceResponse{Data: &service.Resource{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().FindOneById(gomock.Any(), 1).
					Return(&model.Resource{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should return when resource not exists": {
			inputResourceID: 1,
			expectedErr:     &exception.NotFoundException{Err: fmt.Errorf("resource 1 not found")},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().FindOneById(gomock.Any(), 1).
					Return(nil, &exception.NotFoundException{Err: fmt.Errorf("resource 1 not found")})
			},
		},
		"should throw error": {
			inputResourceID: 1,
			expectedErr:     fmt.Errorf("error"),
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockResourceRepository := mock.NewMockResourceRepository(ctrl)
			cs.prepareMock(mockResourceRepository)

			impl := &service.ResourceServiceImpl{ResourceRepository: mockResourceRepository}

			// when
			res, err := impl.FindOneById(ctx, cs.inputResourceID)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_ResourceService_Create(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputDto    service.CreateResourceDto
		expectedRes service.ResourceResponse
		expectedErr error
		prepareMock func(mockResourceRepository *mock.MockResourceRepository)
	}{
		"should create resource": {
			inputDto: service.CreateResourceDto{Name: "Test", Measurement: "Kg"},
			expectedRes: service.ResourceResponse{Data: &service.Resource{
				ID:          1,
				CreatedAt:   DATE,
				UpdatedAt:   DATE,
				Name:        "Test",
				Measurement: "Kg",
			}},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().Create(gomock.Any(), model.Resource{Name: "Test", Measurement: "Kg"}).
					Return(&model.Resource{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test", Measurement: "Kg"}, nil)
			},
		},
		"should create resource when quantity is greater than 0": {
			inputDto: service.CreateResourceDto{Name: "Test", Measurement: "Kg", Quantity: 1},
			expectedRes: service.ResourceResponse{Data: &service.Resource{
				ID:          1,
				CreatedAt:   DATE,
				UpdatedAt:   DATE,
				Name:        "Test",
				Measurement: "Kg",
				Quantity:    1,
			}},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().Create(gomock.Any(), model.Resource{Name: "Test", Measurement: "Kg", Quantity: 1}).
					Return(&model.Resource{
						ID:          1,
						CreatedAt:   DATETIME,
						UpdatedAt:   DATETIME,
						Name:        "Test",
						Measurement: "Kg",
						Quantity:    1,
					}, nil)
			},
		},
		"should throw error": {
			inputDto:    service.CreateResourceDto{},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().Create(gomock.Any(), model.Resource{}).Return(nil, fmt.Errorf("error"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockResourceRepository := mock.NewMockResourceRepository(ctrl)
			cs.prepareMock(mockResourceRepository)

			impl := &service.ResourceServiceImpl{ResourceRepository: mockResourceRepository}

			// when
			res, err := impl.Create(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_ResourceService_Update(t *testing.T) {
	cases := map[string]struct {
		inputDto    service.UpdateResourceDto
		expectedErr error
		prepareMock func(mockResourceRepository *mock.MockResourceRepository)
	}{
		"should update resource": {
			inputDto: service.UpdateResourceDto{ID: 1, Name: "Test"},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().Update(gomock.Any(), model.Resource{ID: 1, Name: "Test"}).Return(nil)
			},
		},
		"should return empty when resource not exists": {
			inputDto:    service.UpdateResourceDto{ID: 1},
			expectedErr: &exception.NotFoundException{Err: fmt.Errorf("resource 1 not found")},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().Update(gomock.Any(), model.Resource{ID: 1}).
					Return(&exception.NotFoundException{Err: fmt.Errorf("resource 1 not found")})
			},
		},
		"should throw error": {
			inputDto:    service.UpdateResourceDto{ID: 1},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().Update(gomock.Any(), model.Resource{ID: 1}).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockResourceRepository := mock.NewMockResourceRepository(ctrl)
			cs.prepareMock(mockResourceRepository)

			impl := &service.ResourceServiceImpl{ResourceRepository: mockResourceRepository}

			// when
			err := impl.Update(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_ResourceService_UpdateQuantity(t *testing.T) {
	cases := map[string]struct {
		inputResourceID int
		inputDto        service.UpdateResourceQuantityDto
		expectedErr     error
		prepareMock     func(mockResourceRepository *mock.MockResourceRepository)
	}{
		"should update resource": {
			inputResourceID: 1,
			inputDto:        service.UpdateResourceQuantityDto{Quantity: 2},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().UpdateQuantity(gomock.Any(), 1, 2.0).Return(nil)
			},
		},
		"should return empty when resource not exists": {
			inputResourceID: 1,
			inputDto:        service.UpdateResourceQuantityDto{Quantity: 2},
			expectedErr:     &exception.NotFoundException{Err: fmt.Errorf("resource 1 not found")},
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().UpdateQuantity(gomock.Any(), 1, 2.0).
					Return(&exception.NotFoundException{Err: fmt.Errorf("resource 1 not found")})
			},
		},
		"should throw error": {
			inputResourceID: 1,
			inputDto:        service.UpdateResourceQuantityDto{Quantity: 2},
			expectedErr:     fmt.Errorf("error"),
			prepareMock: func(mockResourceRepository *mock.MockResourceRepository) {
				mockResourceRepository.EXPECT().UpdateQuantity(gomock.Any(), 1, 2.0).
					Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockResourceRepository := mock.NewMockResourceRepository(ctrl)
			cs.prepareMock(mockResourceRepository)

			impl := &service.ResourceServiceImpl{ResourceRepository: mockResourceRepository}

			// when
			err := impl.UpdateQuantity(ctx, cs.inputResourceID, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
