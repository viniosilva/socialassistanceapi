package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

type ResourceService interface {
	FindAll(ctx context.Context) (ResourcesResponse, error)
	FindOneById(ctx context.Context, resourceID int) (ResourceResponse, error)
	Create(ctx context.Context, dto CreateResourceDto) (ResourceResponse, error)
	Update(ctx context.Context, resourceID int, dto UpdateResourceDto) error
	UpdateQuantity(ctx context.Context, resourceID int, dto UpdateResourceQuantityDto) error
}

type ResourceServiceImpl struct {
	ResourceRepository repository.ResourceRepository
}

func (impl *ResourceServiceImpl) FindAll(ctx context.Context) (ResourcesResponse, error) {
	resources, err := impl.ResourceRepository.FindAll(ctx)
	if err != nil {
		return ResourcesResponse{}, err
	}

	res := []Resource{}
	for _, resource := range resources {
		res = append(res, Resource{
			ID:          resource.ID,
			CreatedAt:   resource.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:   resource.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:        resource.Name,
			Amount:      resource.Amount,
			Measurement: resource.Measurement,
			Quantity:    resource.Quantity,
		})
	}

	return ResourcesResponse{Data: res}, nil
}

func (impl *ResourceServiceImpl) FindOneById(ctx context.Context, resourceID int) (ResourceResponse, error) {
	resource, err := impl.ResourceRepository.FindOneById(ctx, resourceID)
	if err != nil || resource == nil {
		return ResourceResponse{}, err
	}

	return ResourceResponse{
		Data: &Resource{
			ID:          resource.ID,
			CreatedAt:   resource.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:   resource.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:        resource.Name,
			Amount:      resource.Amount,
			Measurement: resource.Measurement,
			Quantity:    resource.Quantity,
		},
	}, nil
}

func (impl *ResourceServiceImpl) Create(ctx context.Context, dto CreateResourceDto) (ResourceResponse, error) {
	resource, err := impl.ResourceRepository.Create(ctx, model.Resource{
		Name:        dto.Name,
		Amount:      dto.Amount,
		Measurement: dto.Measurement,
		Quantity:    dto.Quantity,
	})
	if err != nil {
		return ResourceResponse{}, err
	}

	return ResourceResponse{
		Data: &Resource{
			ID:          resource.ID,
			CreatedAt:   resource.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:   resource.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:        resource.Name,
			Amount:      resource.Amount,
			Measurement: resource.Measurement,
			Quantity:    resource.Quantity,
		},
	}, nil
}

func (impl *ResourceServiceImpl) Update(ctx context.Context, resourceID int, dto UpdateResourceDto) error {
	return impl.ResourceRepository.Update(ctx, model.Resource{
		ID:          resourceID,
		Name:        dto.Name,
		Amount:      dto.Amount,
		Measurement: dto.Measurement,
	})
}

func (impl *ResourceServiceImpl) UpdateQuantity(ctx context.Context, resourceID int, dto UpdateResourceQuantityDto) error {
	return impl.ResourceRepository.UpdateQuantity(ctx, resourceID, dto.Quantity)
}
