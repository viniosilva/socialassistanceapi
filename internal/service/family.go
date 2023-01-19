package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

//go:generate mockgen -destination ../../mock/family_service_mock.go -package mock . FamilyService
type FamilyService interface {
	FindAll(ctx context.Context, limit, offset int) ([]model.Family, int, error)
	FindOneById(ctx context.Context, familyID int) (*model.Family, error)
	Create(ctx context.Context, dto FamilyCreateDto) (*model.Family, error)
	Update(ctx context.Context, dto FamilyUpdateDto) error
	Delete(ctx context.Context, familyID int) error
}

type FamilyServiceImpl struct {
	FamilyRepository repository.FamilyRepository
}

func (impl *FamilyServiceImpl) FindAll(ctx context.Context, limit, offset int) ([]model.Family, int, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.find_all"})

	data, err := impl.FamilyRepository.FindAll(ctx, limit, offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}

	total := 0
	if len(data) > 0 {
		total, err = impl.FamilyRepository.Count(ctx)
		if err != nil {
			log.Error(err.Error())
			return nil, 0, err
		}
	}

	return data, total, nil
}

func (impl *FamilyServiceImpl) FindOneById(ctx context.Context, familyID int) (*model.Family, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.find_one_by_id"})

	data, err := impl.FamilyRepository.FindOneById(ctx, familyID)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return data, nil
}

func (impl *FamilyServiceImpl) Create(ctx context.Context, dto FamilyCreateDto) (*model.Family, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.create"})

	data, err := impl.FamilyRepository.Create(ctx, model.Family{
		Name:         dto.Name,
		Country:      dto.Country,
		State:        dto.State,
		City:         dto.City,
		Neighborhood: dto.Neighborhood,
		Street:       dto.Street,
		Number:       dto.Number,
		Complement:   dto.Complement,
		Zipcode:      dto.Zipcode,
	})

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return data, nil
}

func (impl *FamilyServiceImpl) Update(ctx context.Context, dto FamilyUpdateDto) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.update"})

	if err := impl.FamilyRepository.Update(ctx, model.Family{
		ID:           dto.ID,
		Name:         dto.Name,
		Country:      dto.Country,
		State:        dto.State,
		City:         dto.City,
		Neighborhood: dto.Neighborhood,
		Street:       dto.Street,
		Number:       dto.Number,
		Complement:   dto.Complement,
		Zipcode:      dto.Zipcode,
	}); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (impl *FamilyServiceImpl) Delete(ctx context.Context, familyID int) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.delete"})

	if err := impl.FamilyRepository.Delete(ctx, familyID); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
