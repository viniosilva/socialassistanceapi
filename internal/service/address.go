package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/store"
)

type AddressService struct {
	store store.AddressStore
}

func NewAddressService(store store.AddressStore) *AddressService {
	return &AddressService{store}
}

func (impl *AddressService) FindAll(ctx context.Context) (AddressesResponse, error) {
	addresses, err := impl.store.FindAll(ctx)
	if err != nil {
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

func (impl *AddressService) FindOneById(ctx context.Context, addressID int) (AddressResponse, error) {
	address, err := impl.store.FindOneById(ctx, addressID)
	if err != nil || address == nil {
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

func (impl *AddressService) Create(ctx context.Context, dto AddressDto) (AddressResponse, error) {
	address, err := impl.store.Create(ctx, model.Address{
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

func (impl *AddressService) Update(ctx context.Context, addressID int, dto AddressDto) (AddressResponse, error) {
	address, err := impl.store.Update(ctx, model.Address{
		ID:           addressID,
		Country:      dto.Country,
		State:        dto.State,
		City:         dto.City,
		Neighborhood: dto.Neighborhood,
		Street:       dto.Street,
		Number:       dto.Number,
		Complement:   dto.Complement,
		Zipcode:      dto.Zipcode,
	})
	if err != nil || address == nil {
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

func (impl *AddressService) Delete(ctx context.Context, addressID int) error {
	return impl.store.Delete(ctx, addressID)
}
