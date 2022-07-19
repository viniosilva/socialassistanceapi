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

func (impl *PersonService) FindAll(ctx context.Context) (PersonsResponse, error) {
	persons, err := impl.store.FindAll(ctx)
	if err != nil {
		return PersonsResponse{}, err
	}

	res := []Person{}
	for _, p := range persons {
		res = append(res, Person{
			ID:        p.ID,
			CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05"),
			AddressID: p.AddressID,
			Name:      p.Name,
		})
	}

	return PersonsResponse{Data: res}, nil
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
			AddressID: person.AddressID,
			Name:      person.Name,
		},
	}, nil
}

func (impl *PersonService) Create(ctx context.Context, dto PersonDto) (PersonResponse, error) {
	person, err := impl.store.Create(ctx, model.Person{
		AddressID: dto.AddressID,
		Name:      dto.Name,
	})
	if err != nil {
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        person.ID,
			CreatedAt: person.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: person.UpdatedAt.Format("2006-01-02T15:04:05"),
			AddressID: person.AddressID,
			Name:      person.Name,
		},
	}, nil
}

func (impl *PersonService) Update(ctx context.Context, personID int, dto PersonDto) (PersonResponse, error) {
	person, err := impl.store.Update(ctx, model.Person{
		ID:        personID,
		AddressID: dto.AddressID,
		Name:      dto.Name,
	})
	if err != nil || person == nil {
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        person.ID,
			CreatedAt: person.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: person.UpdatedAt.Format("2006-01-02T15:04:05"),
			AddressID: person.AddressID,
			Name:      person.Name,
		},
	}, nil
}

func (impl *PersonService) Delete(ctx context.Context, personID int) error {
	return impl.store.Delete(ctx, personID)
}
