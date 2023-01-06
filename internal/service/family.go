package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

//go:generate mockgen -destination ../../mock/family_service_mock.go -package mock . FamilyService
type FamilyService interface {
	FindAll(ctx context.Context) (FamiliesResponse, error)
	FindOneById(ctx context.Context, familyID int) (FamilyResponse, error)
	Create(ctx context.Context, dto FamilyCreateDto) (FamilyResponse, error)
	Update(ctx context.Context, dto FamilyUpdateDto) error
	Delete(ctx context.Context, familyID int) error
}

type FamilyServiceImpl struct {
	FamilyRepository repository.FamilyRepository
}

func (impl *FamilyServiceImpl) FindAll(ctx context.Context) (FamiliesResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.find_all"})

	data, err := impl.FamilyRepository.FindAll(ctx)
	if err != nil {
		log.Error(err.Error())
		return FamiliesResponse{}, err
	}

	res := []Family{}
	for _, a := range data {
		res = append(res, Family{
			ID:           a.ID,
			CreatedAt:    a.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:    a.UpdatedAt.Format("2006-01-02T15:04:05"),
			Country:      a.Country,
			State:        a.State,
			City:         a.City,
			Neighborhood: a.Neighborhood,
			Street:       a.Street,
			Number:       a.Number,
			Complement:   a.Complement,
			Zipcode:      a.Zipcode,
		})
	}

	return FamiliesResponse{Data: res}, nil
}

func (impl *FamilyServiceImpl) FindOneById(ctx context.Context, familyID int) (FamilyResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.find_one_by_id"})

	data, err := impl.FamilyRepository.FindOneById(ctx, familyID)
	if err != nil || data == nil {
		log.Error(err.Error())
		return FamilyResponse{}, err
	}

	return FamilyResponse{
		Data: &Family{
			ID:           data.ID,
			CreatedAt:    data.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:    data.UpdatedAt.Format("2006-01-02T15:04:05"),
			Country:      data.Country,
			State:        data.State,
			City:         data.City,
			Neighborhood: data.Neighborhood,
			Street:       data.Street,
			Number:       data.Number,
			Complement:   data.Complement,
			Zipcode:      data.Zipcode,
		},
	}, nil
}

func (impl *FamilyServiceImpl) Create(ctx context.Context, dto FamilyCreateDto) (FamilyResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.create"})

	data, err := impl.FamilyRepository.Create(ctx, model.Family{
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
		return FamilyResponse{}, err
	}

	return FamilyResponse{
		Data: &Family{
			ID:           data.ID,
			CreatedAt:    data.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:    data.UpdatedAt.Format("2006-01-02T15:04:05"),
			Country:      data.Country,
			State:        data.State,
			City:         data.City,
			Neighborhood: data.Neighborhood,
			Street:       data.Street,
			Number:       data.Number,
			Complement:   data.Complement,
			Zipcode:      data.Zipcode,
		},
	}, nil
}

func (impl *FamilyServiceImpl) Update(ctx context.Context, dto FamilyUpdateDto) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.family.update"})

	if err := impl.FamilyRepository.Update(ctx, model.Family{
		ID:           dto.ID,
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
