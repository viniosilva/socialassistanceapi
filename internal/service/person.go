package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

//go:generate mockgen -destination ../../mock/person_service_mock.go -package mock . PersonService
type PersonService interface {
	FindAll(ctx context.Context) (PersonsResponse, error)
	FindOneById(ctx context.Context, personID int) (PersonResponse, error)
	Create(ctx context.Context, dto PersonCreateDto) (PersonResponse, error)
	Update(ctx context.Context, dto PersonUpdateDto) error
	Delete(ctx context.Context, personID int) error
}

type PersonServiceImpl struct {
	PersonRepository repository.PersonRepository
}

func (impl *PersonServiceImpl) FindAll(ctx context.Context) (PersonsResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.find_all"})

	persons, err := impl.PersonRepository.FindAll(ctx)
	if err != nil {
		log.Error(err.Error())
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

func (impl *PersonServiceImpl) FindOneById(ctx context.Context, personID int) (PersonResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.find_one_by_id"})

	person, err := impl.PersonRepository.FindOneById(ctx, personID)
	if err != nil || person == nil {
		log.Error(err.Error())
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

func (impl *PersonServiceImpl) Create(ctx context.Context, dto PersonCreateDto) (PersonResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.create"})

	person, err := impl.PersonRepository.Create(ctx, model.Person{
		AddressID: dto.AddressID,
		Name:      dto.Name,
	})
	if err != nil {
		log.Error(err.Error())
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

func (impl *PersonServiceImpl) Update(ctx context.Context, dto PersonUpdateDto) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.update"})

	if err := impl.PersonRepository.Update(ctx, model.Person{
		ID:        dto.ID,
		AddressID: dto.AddressID,
		Name:      dto.Name,
	}); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (impl *PersonServiceImpl) Delete(ctx context.Context, personID int) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.delete"})

	if err := impl.PersonRepository.Delete(ctx, personID); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
