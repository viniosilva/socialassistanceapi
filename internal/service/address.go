package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

//go:generate mockgen -destination ../../mock/address_service_mock.go -package mock . AddressService
type AddressService interface {
	FindAll(ctx context.Context) (AddressesResponse, error)
	FindOneById(ctx context.Context, addressID int) (AddressResponse, error)
	Create(ctx context.Context, dto AddressCreateDto) (AddressResponse, error)
	Update(ctx context.Context, dto AddressUpdateDto) error
	Delete(ctx context.Context, addressID int) error
}

type AddressServiceImpl struct {
	AddressRepository repository.AddressRepository
}

func (impl *AddressServiceImpl) FindAll(ctx context.Context) (AddressesResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.address.find_all"})

	addresses, err := impl.AddressRepository.FindAll(ctx)
	if err != nil {
		log.Error(err.Error())
		return AddressesResponse{}, err
	}

	res := []Address{}
	for _, a := range addresses {
		res = append(res, Address{
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

	return AddressesResponse{Data: res}, nil
}

func (impl *AddressServiceImpl) FindOneById(ctx context.Context, addressID int) (AddressResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.address.find_one_by_id"})

	address, err := impl.AddressRepository.FindOneById(ctx, addressID)
	if err != nil || address == nil {
		log.Error(err.Error())
		return AddressResponse{}, err
	}

	return AddressResponse{
		Data: &Address{
			ID:           address.ID,
			CreatedAt:    address.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:    address.UpdatedAt.Format("2006-01-02T15:04:05"),
			Country:      address.Country,
			State:        address.State,
			City:         address.City,
			Neighborhood: address.Neighborhood,
			Street:       address.Street,
			Number:       address.Number,
			Complement:   address.Complement,
			Zipcode:      address.Zipcode,
		},
	}, nil
}

func (impl *AddressServiceImpl) Create(ctx context.Context, dto AddressCreateDto) (AddressResponse, error) {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.address.create"})

	address, err := impl.AddressRepository.Create(ctx, model.Address{
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
		return AddressResponse{}, err
	}

	return AddressResponse{
		Data: &Address{
			ID:           address.ID,
			CreatedAt:    address.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:    address.UpdatedAt.Format("2006-01-02T15:04:05"),
			Country:      address.Country,
			State:        address.State,
			City:         address.City,
			Neighborhood: address.Neighborhood,
			Street:       address.Street,
			Number:       address.Number,
			Complement:   address.Complement,
			Zipcode:      address.Zipcode,
		},
	}, nil
}

func (impl *AddressServiceImpl) Update(ctx context.Context, dto AddressUpdateDto) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.address.update"})

	if err := impl.AddressRepository.Update(ctx, model.Address{
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

func (impl *AddressServiceImpl) Delete(ctx context.Context, addressID int) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.address.delete"})

	if err := impl.AddressRepository.Delete(ctx, addressID); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
