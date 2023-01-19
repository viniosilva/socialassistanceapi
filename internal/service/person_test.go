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

func Test_PersonService_FindAll(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		expectedRes service.PersonsResponse
		expectedErr error
		prepareMock func(mockPersonRepository *mock.MockPersonRepository)
	}{
		"should return persons list": {
			expectedRes: service.PersonsResponse{Data: []service.Person{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().FindAll(gomock.Any()).Return([]model.Person{{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}}, nil)
			},
		},
		"should return empty persons list": {
			expectedRes: service.PersonsResponse{Data: []service.Person{}},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().FindAll(gomock.Any()).Return([]model.Person{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockPersonRepository := mock.NewMockPersonRepository(ctrl)
			cs.prepareMock(mockPersonRepository)

			impl := &service.PersonServiceImpl{PersonRepository: mockPersonRepository}

			// when
			res, err := impl.FindAll(ctx)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PersonService_FindOneByID(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputPersonID int
		expectedRes   service.PersonResponse
		expectedErr   error
		prepareMock   func(mockPersonRepository *mock.MockPersonRepository)
	}{
		"should return person when exists": {
			inputPersonID: 1,
			expectedRes:   service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(&model.Person{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should return empty when person not exists": {
			inputPersonID: 1,
			expectedErr:   &exception.NotFoundException{Err: fmt.Errorf("person 1 not found")},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().FindOneById(gomock.Any(), 1).
					Return(nil, &exception.NotFoundException{Err: fmt.Errorf("person 1 not found")})
			},
		},
		"should throw error": {
			inputPersonID: 1,
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().FindOneById(gomock.Any(), 1).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockPersonRepository := mock.NewMockPersonRepository(ctrl)
			cs.prepareMock(mockPersonRepository)

			impl := &service.PersonServiceImpl{PersonRepository: mockPersonRepository}

			// when
			res, err := impl.FindOneById(ctx, cs.inputPersonID)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PersonService_Create(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputDto    service.PersonCreateDto
		expectedRes service.PersonResponse
		expectedErr error
		prepareMock func(mockPersonRepository *mock.MockPersonRepository)
	}{
		"should create person": {
			inputDto:    service.PersonCreateDto{Name: "Test"},
			expectedRes: service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Create(gomock.Any(), model.Person{Name: "Test"}).
					Return(&model.Person{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should throw error": {
			inputDto:    service.PersonCreateDto{Name: "Test"},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Create(gomock.Any(), model.Person{Name: "Test"}).
					Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockPersonRepository := mock.NewMockPersonRepository(ctrl)
			cs.prepareMock(mockPersonRepository)

			impl := &service.PersonServiceImpl{PersonRepository: mockPersonRepository}

			// when
			res, err := impl.Create(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PersonService_Update(t *testing.T) {
	cases := map[string]struct {
		inputDto    service.PersonUpdateDto
		expectedErr error
		prepareMock func(mockPersonRepository *mock.MockPersonRepository)
	}{
		"should update person": {
			inputDto: service.PersonUpdateDto{ID: 1, Name: "Test update"},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Update(gomock.Any(), model.Person{ID: 1, Name: "Test update"}).Return(nil)
			},
		},
		"should return empty when person not exists": {
			inputDto:    service.PersonUpdateDto{ID: 1},
			expectedErr: &exception.NotFoundException{Err: fmt.Errorf("person 1 not found")},
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Update(gomock.Any(), model.Person{ID: 1}).
					Return(&exception.NotFoundException{Err: fmt.Errorf("person 1 not found")})
			},
		},
		"should throw error": {
			inputDto:    service.PersonUpdateDto{ID: 1, Name: "Test update"},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Update(gomock.Any(), model.Person{ID: 1, Name: "Test update"}).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockPersonRepository := mock.NewMockPersonRepository(ctrl)
			cs.prepareMock(mockPersonRepository)

			impl := &service.PersonServiceImpl{PersonRepository: mockPersonRepository}

			// when
			err := impl.Update(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PersonService_Delete(t *testing.T) {
	cases := map[string]struct {
		inputPersonID int
		expectedErr   error
		prepareMock   func(mockPersonRepository *mock.MockPersonRepository)
	}{
		"should delete person": {
			inputPersonID: 1,
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
		},
		"should throw error": {
			inputPersonID: 1,
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mockPersonRepository *mock.MockPersonRepository) {
				mockPersonRepository.EXPECT().Delete(gomock.Any(), 1).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockPersonRepository := mock.NewMockPersonRepository(ctrl)
			cs.prepareMock(mockPersonRepository)

			impl := &service.PersonServiceImpl{PersonRepository: mockPersonRepository}

			// when
			err := impl.Delete(ctx, cs.inputPersonID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
