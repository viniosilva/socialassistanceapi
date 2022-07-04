package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/store"
)

type ResourceService struct {
	store store.ResourceStore
}

func NewResourceService(store store.ResourceStore) *ResourceService {
	return &ResourceService{store}
}

func (impl *ResourceService) FindAll(ctx context.Context) (ResourcesResponse, error) {
	resource, err := impl.store.FindAll(ctx)
	if err != nil {
		return ResourcesResponse{}, err
	}

	res := []Resource{}
	for _, c := range resource {
		res = append(res, Resource{
			ID:          c.ID,
			CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:        c.Name,
			Amount:      c.Amount,
			Measurement: c.Measurement,
		})
	}

	return ResourcesResponse{Data: res}, nil
}

func (impl *ResourceService) FindOneById(ctx context.Context, resourceID int) (ResourceResponse, error) {
	resource, err := impl.store.FindOneById(ctx, resourceID)
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
		},
	}, nil
}

func (impl *ResourceService) Create(ctx context.Context, dto ResourceDto) (ResourceResponse, error) {
	resource, err := impl.store.Create(ctx, model.Resource{
		Name:        dto.Name,
		Amount:      dto.Amount,
		Measurement: dto.Measurement,
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
		},
	}, nil
}

func (impl *ResourceService) Update(ctx context.Context, ResourceID int, dto ResourceDto) (ResourceResponse, error) {
	resource, err := impl.store.Update(ctx, model.Resource{
		ID:          ResourceID,
		Name:        dto.Name,
		Measurement: dto.Measurement,
	})
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
		},
	}, nil
}
