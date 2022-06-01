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

func (impl *ResourceService) FindAll(ctx context.Context) (ResourceResponse, error) {
	resource, err := impl.store.FindAll(ctx)
	if err != nil {
		return ResourceResponse{}, err
	}

	res := []Resource{}
	for _, c := range resource {
		res = append(res, Resource{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      c.Name,
		})
	}

	return ResourceResponse{Data: res}, nil
}

func (impl *ResourceService) FindOneById(ctx context.Context, resourceID int) (ResourceResponse, error) {
	resource, err := impl.store.FindOneById(ctx, resourceID)
	if err != nil || resource == nil {
		return ResourceResponse{}, err
	}

	return ResourceResponse{
		Data: &Resource{
			ID:        resource.ID,
			CreatedAt: resource.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: resource.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      resource.Name,
		},
	}, nil
}

func (impl *ResourceService) Create(ctx context.Context, dto ResourceDto) (ResourceResponse, error) {
	resource, err := impl.store.Create(ctx, model.Resource{Name: dto.Name})
	if err != nil {
		return ResourceResponse{}, err
	}

	return ResourceResponse{
		Data: &Resource{
			ID:        resource.ID,
			CreatedAt: resource.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: resource.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      resource.Name,
		},
	}, nil
}

func (impl *ResourceService) Update(ctx context.Context, ResourceID int, dto ResourceDto) (ResourceResponse, error) {
	resource, err := impl.store.Update(ctx, model.Resource{
		ID:   ResourceID,
		Name: dto.Name,
	})
	if err != nil || resource == nil {
		return ResourceResponse{}, err
	}

	return ResourceResponse{
		Data: &Resource{
			ID:        resource.ID,
			CreatedAt: resource.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt: resource.UpdatedAt.Format("2006-01-02T15:04:05"),
			Name:      resource.Name,
		},
	}, nil
}

func (impl *ResourceService) Delete(ctx context.Context, resourceID int) (ResourceResponse, error) {
	err := impl.store.Delete(ctx, resourceID)
	return ResourceResponse{}, err
}