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

func TestResourceServiceFindAll(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		expectedRes service.ResourcesResponse
		expectedErr error
		prepareMock func(mock *mock.MockResourceStore)
	}{
		"should return resource list": {
			expectedRes: service.ResourcesResponse{Data: []service.Resource{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Resource{{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}}, nil)
			},
		},
		"should return empty resource list": {
			expectedRes: service.ResourcesResponse{Data: []service.Resource{}},
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Resource{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockResourceStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewResourceService(storeMock)

			// when
			res, err := impl.FindAll(ctx)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("ResourceService.FindAll() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("ResourceService.FindAll() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestResourceServiceFindOneByID(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputResourceID int
		expectedRes     service.ResourceResponse
		expectedErr     error
		prepareMock     func(mock *mock.MockResourceStore)
	}{
		"should return resource when exists": {
			inputResourceID: 1,
			expectedRes:     service.ResourceResponse{Data: &service.Resource{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(&model.Resource{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should return when resource not exists": {
			inputResourceID: 1,
			expectedErr:     exception.NewNotFoundException("resource"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).
					Return(nil, exception.NewNotFoundException("resource"))
			},
		},
		"should throw error": {
			inputResourceID: 1,
			expectedErr:     fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().FindOneById(gomock.All(), gomock.All()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockResourceStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewResourceService(storeMock)

			// when
			res, err := impl.FindOneById(ctx, cs.inputResourceID)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("Resource.Service.FindOneById() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("ResourceService.FindOneById() error = %v, epected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestResourceServiceCreate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputResource service.ResourceDto
		expectedRes   service.ResourceResponse
		expectedErr   error
		prepareMock   func(mock *mock.MockResourceStore)
	}{
		"should create resource": {
			inputResource: service.ResourceDto{Name: "Teste"},
			expectedRes:   service.ResourceResponse{Data: &service.Resource{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Create(gomock.All(), gomock.All()).Return(&model.Resource{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should throw error": {
			inputResource: service.ResourceDto{},
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Create(gomock.All(), gomock.All()).Return(nil, fmt.Errorf("error"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockResourceStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewResourceService(storeMock)

			// when
			res, err := impl.Create(ctx, cs.inputResource)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("ResourceService.FindOneById() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("ResourceService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestResourceServiceUpdate(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputResourceID int
		inputResource   service.ResourceUpdateDto
		expectedRes     service.ResourceResponse
		expectedErr     error
		prepareMock     func(mock *mock.MockResourceStore)
	}{
		"should update resource": {
			inputResourceID: 1,
			inputResource:   service.ResourceUpdateDto{Name: "Test"},
			expectedRes: service.ResourceResponse{Data: &service.Resource{
				ID: 1, CreatedAt: DATE, UpdatedAt: DATE,
				Name:        "Test",
				Amount:      1,
				Measurement: "Kg",
			}},
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Update(gomock.All(), gomock.All()).Return(&model.Resource{
					ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME,
					Name:        "Test",
					Amount:      1,
					Measurement: "Kg",
				}, nil)
			},
		},
		"should return empty when resource not exists": {
			inputResourceID: 1,
			expectedErr:     exception.NewNotFoundException("resource"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Update(gomock.All(), gomock.All()).
					Return(nil, exception.NewNotFoundException("resource"))
			},
		},
		"should throw error": {
			inputResource: service.ResourceUpdateDto{Name: "Test"},
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockResourceStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewResourceService(storeMock)

			// when
			res, err := impl.Update(ctx, cs.inputResourceID, cs.inputResource)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("ResourceService.Update() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("ResourceService.Update() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestResourceServiceTransferAmount(t *testing.T) {
	DATE := "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputResourceID int
		inputResource   service.ResourceUpdateDto
		expectedRes     service.ResourceResponse
		expectedErr     error
		prepareMock     func(mock *mock.MockResourceStore)
	}{
		"should update resource": {
			inputResourceID: 1,
			inputResource:   service.ResourceUpdateDto{Name: "Test"},
			expectedRes: service.ResourceResponse{Data: &service.Resource{
				ID: 1, CreatedAt: DATE, UpdatedAt: DATE,
				Name:        "Test",
				Amount:      1,
				Measurement: "Kg",
			}},
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Update(gomock.All(), gomock.All()).Return(&model.Resource{
					ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME,
					Name:        "Test",
					Amount:      1,
					Measurement: "Kg",
				}, nil)
			},
		},
		"should return empty when resource not exists": {
			inputResourceID: 1,
			expectedErr:     exception.NewNotFoundException("resource"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Update(gomock.All(), gomock.All()).
					Return(nil, exception.NewNotFoundException("resource"))
			},
		},
		"should throw error": {
			inputResource: service.ResourceUpdateDto{Name: "Test"},
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockResourceStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockResourceStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewResourceService(storeMock)

			// when
			res, err := impl.Update(ctx, cs.inputResourceID, cs.inputResource)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("ResourceService.Update() = %v, expected %v", res, cs.expectedRes)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("ResourceService.Update() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}
