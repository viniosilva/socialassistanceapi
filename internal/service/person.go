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

	data, err := impl.PersonRepository.FindAll(ctx)
	if err != nil {
		log.Error(err.Error())
		return PersonsResponse{}, err
	}

	res := []Person{}
	for _, d := range data {
		res = append(res, Person{
			ID:        d.ID,
			CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05"),
			FamilyID:  d.FamilyID,
			Name:      d.Name,
		})
	}

	return PersonsResponse{Data: res}, nil
}

func (impl *PersonServiceImpl) FindOneById(ctx context.Context, personID int) (PersonResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.find_one_by_id"})

	data, err := impl.PersonRepository.FindOneById(ctx, personID)
	if err != nil || data == nil {
		log.Error(err.Error())
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        data.ID,
			CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05"),
			FamilyID:  data.FamilyID,
			Name:      data.Name,
		},
	}, nil
}

func (impl *PersonServiceImpl) Create(ctx context.Context, dto PersonCreateDto) (PersonResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.create"})

	data, err := impl.PersonRepository.Create(ctx, model.Person{
		FamilyID: dto.FamilyID,
		Name:     dto.Name,
	})
	if err != nil {
		log.Error(err.Error())
		return PersonResponse{}, err
	}

	return PersonResponse{
		Data: &Person{
			ID:        data.ID,
			CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05"),
			FamilyID:  data.FamilyID,
			Name:      data.Name,
		},
	}, nil
}

func (impl *PersonServiceImpl) Update(ctx context.Context, dto PersonUpdateDto) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.person.update"})

	if err := impl.PersonRepository.Update(ctx, model.Person{
		ID:       dto.ID,
		FamilyID: dto.FamilyID,
		Name:     dto.Name,
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
