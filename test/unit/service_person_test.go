package unit

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/mock"
)

func TestPersonServiceFindAll(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		expectedPeople service.PeopleResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockPersonStore)
	}{
		"should return person list": {
			expectedPeople: service.PeopleResponse{Data: []service.Person{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Person{{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}}, nil)
			},
		},
		"should return empty person list": {
			expectedPeople: service.PeopleResponse{Data: []service.Person{}},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]model.Person{}, nil)
			},
		},
		"should throw error": {
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockPersonStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewPersonService(storeMock)

			// when
			people, err := impl.FindAll(ctx)

			// then
			if !reflect.DeepEqual(people, cs.expectedPeople) {
				t.Errorf("PersonService.FindAll() = %v, expected %v", people, cs.expectedPeople)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("PersonService.FindAll() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestPersonServiceFindOneByID(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputPersonID  int
		expectedPerson service.PersonResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockPersonStore)
	}{
		"should return person when exists": {
			inputPersonID:  1,
			expectedPerson: service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(&model.Person{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should return empty when person not exists": {
			inputPersonID:  1,
			expectedPerson: service.PersonResponse{},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		"should throw error": {
			inputPersonID: 1,
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().FindOneById(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockPersonStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewPersonService(storeMock)

			// when
			person, err := impl.FindOneById(ctx, cs.inputPersonID)

			// then
			if !reflect.DeepEqual(person, cs.expectedPerson) {
				t.Errorf("PersonService.FindOneById() = %v, expected %v", person, cs.expectedPerson)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("PersonService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestPersonServiceCreate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputPerson    service.PersonDto
		expectedPerson service.PersonResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockPersonStore)
	}{
		"should create person": {
			inputPerson:    service.PersonDto{Name: "Test"},
			expectedPerson: service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Person{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test"}, nil)
			},
		},
		"should throw error": {
			inputPerson: service.PersonDto{Name: "Test"},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockPersonStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewPersonService(storeMock)

			// when
			person, err := impl.Create(ctx, cs.inputPerson)

			// then
			if !reflect.DeepEqual(person, cs.expectedPerson) {
				t.Errorf("PersonService.FindOneById() = %v, expected %v", person, cs.expectedPerson)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("PersonService.FindOneById() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestPersonServiceUpdate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"
	DATETIME := time.Date(2000, 1, 1, 12, 3, 0, 0, time.UTC)

	cases := map[string]struct {
		inputPersonID  int
		inputPerson    service.PersonDto
		expectedPerson service.PersonResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockPersonStore)
	}{
		"should update person": {
			inputPersonID:  1,
			inputPerson:    service.PersonDto{Name: "Test update"},
			expectedPerson: service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test update"}},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&model.Person{ID: 1, CreatedAt: DATETIME, UpdatedAt: DATETIME, Name: "Test update"}, nil)
			},
		},
		"should return empty when person not exists": {
			inputPersonID:  1,
			expectedPerson: service.PersonResponse{},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		"should throw error": {
			inputPerson: service.PersonDto{Name: "Test"},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockPersonStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewPersonService(storeMock)

			// when
			person, err := impl.Update(ctx, cs.inputPersonID, cs.inputPerson)

			// then
			if !reflect.DeepEqual(person, cs.expectedPerson) {
				t.Errorf("PersonService.Update() = %v, expected %v", person, cs.expectedPerson)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("PersonService.Update() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}

func TestPersonServiceDelete(t *testing.T) {
	cases := map[string]struct {
		inputPersonID  int
		expectedPerson service.PersonResponse
		expectedErr    error
		prepareMock    func(mock *mock.MockPersonStore)
	}{
		"should delete person": {
			inputPersonID:  1,
			expectedPerson: service.PersonResponse{},
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		"should throw error": {
			inputPersonID: 1,
			expectedErr:   fmt.Errorf("error"),
			prepareMock: func(mock *mock.MockPersonStore) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			storeMock := mock.NewMockPersonStore(ctrl)
			cs.prepareMock(storeMock)

			impl := service.NewPersonService(storeMock)

			// when
			person, err := impl.Delete(ctx, cs.inputPersonID)

			// then
			if !reflect.DeepEqual(person, cs.expectedPerson) {
				t.Errorf("PersonService.Delete() = %v, expected %v", person, cs.expectedPerson)
			}
			if err != nil && err.Error() != cs.expectedErr.Error() {
				t.Errorf("PersonService.Delete() error = %v, expected %v", err, cs.expectedErr)
			}
		})
	}
}
