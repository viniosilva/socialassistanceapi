package service

import (
	"context"

	"github.com/sirupsen/logrus"
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
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.donate_resource.donate"})

	if err := impl.DonateResourceRepository.Donate(ctx, dto.ResourceID, dto.AddressID, dto.Quantity); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (impl *DonateResourceServiceImpl) Return(ctx context.Context, resourceID int) error {
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.donate_resource.return"})

	if err := impl.DonateResourceRepository.Return(ctx, resourceID); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
