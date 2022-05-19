package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/store"
)

type PersonService struct {
	store store.PersonStore
}

func NewPersonService(store store.PersonStore) *PersonService {
	return &PersonService{store}
}

func (impl *PersonService) FindAll(ctx context.Context) (PeopleResponse, error) {
	people, err := impl.store.FindAll(ctx)
	if err != nil {
		return PeopleResponse{}, err
	}

	res := []Person{}
	for _, c := range people {
		res = append(res, Person{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      c.Name,
		})
	}

	return PeopleResponse{Data: res}, nil
}

func (impl *PersonService) FindOneById(ctx context.Context, personID int) (PersonResponse, error) {
	person, err := impl.store.FindOneById(ctx, personID)
	if err != nil || person == nil {
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        person.ID,
			CreatedAt: person.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: person.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      person.Name,
		},
	}, nil
}

func (impl *PersonService) Create(ctx context.Context, dto PersonDto) (PersonResponse, error) {
	person, err := impl.store.Create(ctx, model.Person{Name: dto.Name})
	if err != nil {
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        person.ID,
			CreatedAt: person.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: person.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      person.Name,
		},
	}, nil
}

func (impl *PersonService) Update(ctx context.Context, personID int, dto PersonDto) (PersonResponse, error) {
	person, err := impl.store.Update(ctx, model.Person{
		ID:   personID,
		Name: dto.Name,
	})
	if err != nil || person == nil {
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        person.ID,
			CreatedAt: person.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: person.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      person.Name,
		},
	}, nil
}

func (impl *PersonService) Delete(ctx context.Context, personID int) (PersonResponse, error) {
	err := impl.store.Delete(ctx, personID)
	return PersonResponse{}, err
}
