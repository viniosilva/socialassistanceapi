package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

type DonateResourceService interface {
	Donate(ctx context.Context, dto DonateResourceDonateDto) error
	Return(ctx context.Context, resourceID int) error
}

type DonateResourceServiceImpl struct {
	DonateResourceRepository repository.DonateResourceRepository
}

func (impl *DonateResourceServiceImpl) Donate(ctx context.Context, dto DonateResourceDonateDto) error {
	if err := impl.DonateResourceRepository.Donate(ctx, dto.ResourceID, dto.AddressID, dto.Quantity); err != nil {
		return err
	}

	return nil
}

func (impl *DonateResourceServiceImpl) Return(ctx context.Context, resourceID int) error {
	if err := impl.DonateResourceRepository.Return(ctx, resourceID); err != nil {
		return err
	}

	return nil
}
