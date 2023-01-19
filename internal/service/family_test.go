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

func Test_FamilyService_FindAll(t *testing.T) {
	cases := map[string]struct {
		inputLimit    int
		inputOffset   int
		expectedRes   []model.Family
		expectedTotal int
		expectedErr   error
		prepareMock   func(mockFamilyRepository *mock.MockFamilyRepository)
	}{
		"should return families list": {
			inputLimit:  10,
			inputOffset: 0,
			expectedRes: []model.Family{{
				ID:           1,
				CreatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
				Name:         "Sauro",
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
			expectedTotal: 1,
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().FindAll(gomock.Any(), 10, 0).Return([]model.Family{{
					ID:           1,
					CreatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
					UpdatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
					Name:         "Sauro",
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}}, nil)
				mockFamilyRepository.EXPECT().Count(gomock.Any()).Return(1, nil)
			},
		},
		"should return empty families list": {
			inputLimit:    10,
			inputOffset:   0,
			expectedRes:   []model.Family{},
			expectedTotal: 0,
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().FindAll(gomock.Any(), 10, 0).Return([]model.Family{}, nil)
			},
		},
		"should throw error when FindAll": {
			inputLimit:    10,
			inputOffset:   0,
			expectedErr:   fmt.Errorf("error"),
			expectedTotal: 0,
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().FindAll(gomock.Any(), 10, 0).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockFamilyRepository := mock.NewMockFamilyRepository(ctrl)
			cs.prepareMock(mockFamilyRepository)

			impl := &service.FamilyServiceImpl{FamilyRepository: mockFamilyRepository}

			// when
			res, total, err := impl.FindAll(ctx, cs.inputLimit, cs.inputOffset)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedTotal, total)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FamilyService_FindOneByID(t *testing.T) {
	cases := map[string]struct {
		inputFamilyID int
		expectedRes   *model.Family
		expectedErr   error
		prepareMock   func(mockFamilyRepository *mock.MockFamilyRepository)
	}{
		"should return family when exists": {
			inputFamilyID: 1,
			expectedRes: &model.Family{
				ID:           1,
				CreatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
				Name:         "Sauro",
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(&model.Family{
					ID:           1,
					CreatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
					UpdatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
					Name:         "Sauro",
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
		"should return empty when family not exists": {
			inputFamilyID: 1,
			expectedErr:   &exception.NotFoundException{Err: fmt.Errorf("family 1 not found")},
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().FindOneById(gomock.Any(), 1).
					Return(nil, &exception.NotFoundException{Err: fmt.Errorf("family 1 not found")})
			},
		},
		"should throw error": {
			inputFamilyID: 1,
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockFamilyRepository := mock.NewMockFamilyRepository(ctrl)
			cs.prepareMock(mockFamilyRepository)

			impl := &service.FamilyServiceImpl{FamilyRepository: mockFamilyRepository}

			// when
			res, err := impl.FindOneById(ctx, cs.inputFamilyID)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FamilyService_Create(t *testing.T) {
	cases := map[string]struct {
		inputDto    service.FamilyCreateDto
		expectedRes *model.Family
		expectedErr error
		prepareMock func(mockFamilyRepository *mock.MockFamilyRepository)
	}{
		"should create family": {
			inputDto: service.FamilyCreateDto{
				Name:         "Sauro",
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedRes: &model.Family{
				ID:           1,
				CreatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
				Name:         "Sauro",
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Create(gomock.Any(), model.Family{
					Name:         "Sauro",
					Country:      "BR",
					State:        "SP",
					City:         "São Paulo",
					Neighborhood: "Pq. Novo Mundo",
					Street:       "R. Sd. Teodoro Francisco Ribeiro",
					Number:       "1",
					Complement:   "1",
					Zipcode:      "02180110",
				}).Return(&model.Family{
					ID:           1,
					CreatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
					UpdatedAt:    time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC),
					Name:         "Sauro",
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
			inputDto: service.FamilyCreateDto{
				Name:         "Sauro",
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
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Create(gomock.Any(), model.Family{
					Name:         "Sauro",
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

			mockFamilyRepository := mock.NewMockFamilyRepository(ctrl)
			cs.prepareMock(mockFamilyRepository)

			impl := &service.FamilyServiceImpl{FamilyRepository: mockFamilyRepository}

			// when
			res, err := impl.Create(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FamilyService_Update(t *testing.T) {
	cases := map[string]struct {
		inputDto    service.FamilyUpdateDto
		expectedErr error
		prepareMock func(mockFamilyRepository *mock.MockFamilyRepository)
	}{
		"should update family": {
			inputDto: service.FamilyUpdateDto{
				ID:           1,
				Name:         "Sauro",
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			},
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Update(gomock.Any(), model.Family{
					ID:           1,
					Name:         "Sauro",
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
		"should return empty when family not exists": {
			inputDto:    service.FamilyUpdateDto{ID: 1},
			expectedErr: &exception.NotFoundException{Err: fmt.Errorf("family 1 not found")},
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Update(gomock.Any(), model.Family{ID: 1}).
					Return(&exception.NotFoundException{Err: fmt.Errorf("family 1 not found")})
			},
		},
		"should throw error": {
			inputDto: service.FamilyUpdateDto{
				ID:           1,
				Name:         "Sauro",
				Country:      "BR",
				State:        "RS",
				City:         "Porto Alegre",
				Neighborhood: "Hípica",
				Street:       "R. J",
				Number:       "1",
				Zipcode:      "91755450",
			},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Update(gomock.Any(), model.Family{
					ID:           1,
					Name:         "Sauro",
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

			mockFamilyRepository := mock.NewMockFamilyRepository(ctrl)
			cs.prepareMock(mockFamilyRepository)

			impl := &service.FamilyServiceImpl{FamilyRepository: mockFamilyRepository}

			// when
			err := impl.Update(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FamilyService_Delete(t *testing.T) {
	cases := map[string]struct {
		inputFamilyID int
		expectedErr   error
		prepareMock   func(mockFamilyRepository *mock.MockFamilyRepository)
	}{
		"should delete family": {
			inputFamilyID: 1,
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
		},
		"should throw error": {
			inputFamilyID: 1,
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mockFamilyRepository *mock.MockFamilyRepository) {
				mockFamilyRepository.EXPECT().Delete(gomock.Any(), 1).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockFamilyRepository := mock.NewMockFamilyRepository(ctrl)
			cs.prepareMock(mockFamilyRepository)

			impl := &service.FamilyServiceImpl{FamilyRepository: mockFamilyRepository}

			// when
			err := impl.Delete(ctx, cs.inputFamilyID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
